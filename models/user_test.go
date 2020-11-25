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

func TestCheckPasswordWithValidPassword(t *testing.T) {
	hashedPassword, _ := hashPassword("password")
	result := CheckPassword(string(hashedPassword), "password")

	assert.Equal(t, nil, result)
}

func TestCheckPasswordWithInvalidPassword(t *testing.T) {
	hashedPassword, _ := hashPassword("password")
	result := CheckPassword(string(hashedPassword), "drowssap")

	assert.NotEqual(t, nil, result)
}

type DBTestSuite struct {
	suite.Suite
	userID   uint
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

	user := &User{Email: s.email, Password: s.password}
	db.GetDB().Create(user)

	s.userID = user.ID
}

func (s *DBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (s *DBTestSuite) TestSaveUserWithValidParams() {
	db.GetDB().Exec("DELETE FROM users")
	err := SaveUser(s.email, s.password)
	assert.Equal(s.T(), nil, err)

	user := &User{}
	db.GetDB().First(user)
	assert.Equal(s.T(), s.email, user.Email)
}

func (s *DBTestSuite) TestSaveUserWithDuplicateEmail() {
	err := SaveUser(s.email, s.password)
	assert.Equal(s.T(), nil, err)

	result := db.GetDB().First(&User{})
	assert.Equal(s.T(), int64(1), result.RowsAffected)
}

func (s *DBTestSuite) TestSaveUserWithEmptyStringEmail() {
	db.GetDB().Exec("DELETE FROM users")
	err := SaveUser("", "password")
	assert.Equal(s.T(), "Email or password cannot be blank", err.Error())

	result := db.GetDB().First(&User{})
	assert.Equal(s.T(), int64(0), result.RowsAffected)
}

func (s *DBTestSuite) TestSaveUserWithEmptyStringPassword() {
	db.GetDB().Exec("DELETE FROM users")
	err := SaveUser("email@email.com", "")
	assert.Equal(s.T(), "Email or password cannot be blank", err.Error())

	result := db.GetDB().First(&User{})
	assert.Equal(s.T(), int64(0), result.RowsAffected)
}

func (s *DBTestSuite) TestFindOneUserByConditionWithValidParams() {
	user, err := FindOneUserByCondition(&User{Email: s.email})

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.email, user.Email)
}

func (s *DBTestSuite) TestFindOneUserByConditionWithInvalidParams() {
	user, err := FindOneUserByCondition(&User{Email: "test"})

	assert.NotEqual(s.T(), nil, err)
	assert.Equal(s.T(), &User{}, user)
}

func (s *DBTestSuite) TestFindOneUserByIDWithValidID() {
	user, err := FindOneUserByID(s.userID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.email, user.Email)
}

func (s *DBTestSuite) TestFindOneUserByIDWithInvalidID() {
	user, err := FindOneUserByID("testID")

	assert.NotEqual(s.T(), nil, err)
	assert.Equal(s.T(), &User{}, user)
}
