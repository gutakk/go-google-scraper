package controllers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createFormFile(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", "text/csv")
	return w.CreatePart(h)
}

func createMultipartPayload(filename string) (http.Header, *bytes.Buffer) {
	path := filename
	file, _ := os.Open(path)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := createFormFile(writer, "file", filepath.Base(path))
	_, _ = io.Copy(part, file)
	writer.Close()

	headers := http.Header{}
	headers.Set("Content-Type", writer.FormDataContentType())

	return headers, body
}

func performRequest(r http.Handler, method, path string, headers http.Header, payload *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, payload)
	req.Header = headers

	response := httptest.NewRecorder()

	r.ServeHTTP(response, req)

	return response
}

type KeywordDbTestSuite struct {
	suite.Suite
	engine *gin.Engine
	cookie string
}

func (s *KeywordDbTestSuite) SetupTest() {
	testDB, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return testDB
	}

	_ = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})

	s.engine = testConfig.GetRouter(true)
	new(LoginController).applyRoutes(EnsureGuestUserGroup(s.engine))
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(s.engine))

	email := faker.Email()
	password := faker.Password()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Email: email, Password: string(hashedPassword)}
	db.GetDB().Create(&user)

	formData := url.Values{}
	formData.Set("email", email)
	formData.Set("password", password)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	response := testHttp.PerformRequest(s.engine, "POST", "/login", headers, formData)

	s.cookie = response.Header().Get("Set-Cookie")
}

func (s *KeywordDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
}

func TestKeywordDbTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordDbTestSuite))
}

func (s *KeywordDbTestSuite) TestUploadKeywordWithValidParams() {
	headers, payload := createMultipartPayload("tests/csv/adword_keywords.csv")
	headers.Set("Cookie", s.cookie)

	response := performRequest(s.engine, "POST", "/keyword", headers, payload)

	p, err := ioutil.ReadAll(response.Body)
	isKeywordPage := err == nil && strings.Index(string(p), "<title>Keyword</title>") > 0
	pageNotice := err == nil && strings.Index(string(p), "CSV uploaded successfully") > 0

	assert.Equal(s.T(), http.StatusOK, response.Code)
	assert.Equal(s.T(), true, isKeywordPage)
	assert.Equal(s.T(), true, pageNotice)
}

func (s *KeywordDbTestSuite) TestUploadKeywordWithBlankPayload() {
	headers := http.Header{}
	headers.Set("Cookie", s.cookie)

	response := performRequest(s.engine, "POST", "/keyword", headers, &bytes.Buffer{})

	p, err := ioutil.ReadAll(response.Body)
	isKeywordPage := err == nil && strings.Index(string(p), "<title>Keyword</title>") > 0
	pageError := err == nil && strings.Index(string(p), "File is required") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, isKeywordPage)
	assert.Equal(s.T(), true, pageError)
}

func TestDisplayKeywordWithGuestUser(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	response := testHttp.PerformRequest(engine, "GET", "/keyword", nil, nil)

	assert.Equal(t, http.StatusFound, response.Code)
	assert.Equal(t, "/login", response.Header().Get("Location"))
}

func TestUploadKeywordWithGuestUser(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	response := testHttp.PerformRequest(engine, "POST", "/keyword", nil, nil)

	assert.Equal(t, http.StatusFound, response.Code)
	assert.Equal(t, "/login", response.Header().Get("Location"))
}

func TestDisplayKeywordWithAuthenticatedUser(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	// Cookie from login API Set-Cookie header
	headers := http.Header{}
	cookie := "go-google-scraper=MTYwNjQ2Mjk3MXxEdi1CQkFFQ180SUFBUkFCRUFBQUlmLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FFZFdsdWRBWUVBUDRFdFE9PXzl6APqAQw3gAQqlHoXMYrPpnqPFkEP8SRHJZEpl-_LDQ=="
	headers.Set("Cookie", cookie)

	response := testHttp.PerformRequest(engine, "GET", "/keyword", headers, nil)
	p, err := ioutil.ReadAll(response.Body)
	isKeywordPage := err == nil && strings.Index(string(p), "<title>Keyword</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, isKeywordPage)
}
