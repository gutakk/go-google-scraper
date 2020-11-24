package models

import (
	"os"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/tests"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestHashPassword(t *testing.T) {
	hashedPassword, _ := hashPassword("password")
	result := bcrypt.CompareHashAndPassword(hashedPassword, []byte("password"))

	assert.Equal(t, nil, result)
}

type DBTestSuite struct {
	suite.Suite
	email    string
	password string
}

func (s *DBTestSuite) SetupTest() {
	testDB, _ := gorm.Open(postgres.Open(tests.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return testDB
	}

	_ = db.GetDB().AutoMigrate(&User{})

	s.email = faker.Email()
	s.password = faker.Password()
}

func (s *DBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (s *DBTestSuite) TestSaveUserWithValidParams() {
	err := SaveUser(s.email, s.password)
	assert.Equal(s.T(), nil, err)

	user := &User{}
	db.GetDB().First(user)
	assert.Equal(s.T(), s.email, user.Email)
}

func (s *DBTestSuite) TestSaveUserWithDuplicateEmail() {
	db.GetDB().Create(&User{Email: s.email, Password: s.password})

	err := SaveUser(s.email, s.password)
	assert.NotEqual(s.T(), nil, err)
}
