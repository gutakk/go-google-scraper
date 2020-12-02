package models

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/gutakk/go-google-scraper/db"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserDBTestSuite struct {
	suite.Suite
	userID   uint
	email    string
	password string
}

func (s *UserDBTestSuite) SetupTest() {
	testDB, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
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

func (s *UserDBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
}

func TestUserDBTestSuite(t *testing.T) {
	suite.Run(t, new(UserDBTestSuite))
}

func (s *UserDBTestSuite) TestSaveUserWithValidParams() {
	db.GetDB().Exec("DELETE FROM users")
	err := SaveUser(s.email, s.password)
	assert.Equal(s.T(), nil, err)

	user := &User{}
	db.GetDB().First(user)
	assert.Equal(s.T(), s.email, user.Email)
}

func (s *UserDBTestSuite) TestSaveUserWithDuplicateEmail() {
	err := SaveUser(s.email, s.password)
	assert.NotEqual(s.T(), nil, err)

	result := db.GetDB().First(&User{})
	assert.Equal(s.T(), int64(1), result.RowsAffected)
}

func (s *UserDBTestSuite) TestSaveUserWithEmptyStringEmail() {
	db.GetDB().Exec("DELETE FROM users")
	err := SaveUser("", "password")
	assert.Equal(s.T(), "Email or password cannot be blank", err.Error())

	result := db.GetDB().First(&User{})
	assert.Equal(s.T(), int64(0), result.RowsAffected)
}

func (s *UserDBTestSuite) TestSaveUserWithEmptyStringPassword() {
	db.GetDB().Exec("DELETE FROM users")
	err := SaveUser("email@email.com", "")
	assert.Equal(s.T(), "Email or password cannot be blank", err.Error())

	result := db.GetDB().First(&User{})
	assert.Equal(s.T(), int64(0), result.RowsAffected)
}

func (s *UserDBTestSuite) TestFindUserByConditionWithValidEmail() {
	user, err := FindUserBy(&User{Email: s.email})

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.email, user.Email)
}

func (s *UserDBTestSuite) TestFindUserByConditionWithInvalidEmail() {
	user, err := FindUserBy(&User{Email: "test"})

	assert.NotEqual(s.T(), nil, err)
	assert.Equal(s.T(), &User{}, user)
}

func (s *UserDBTestSuite) TestFindUserByIDWithValidID() {
	user, err := FindUserByID(s.userID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), s.email, user.Email)
}

func (s *UserDBTestSuite) TestFindUserByIDWithInvalidID() {
	user, err := FindUserByID("testID")

	assert.NotEqual(s.T(), nil, err)
	assert.Equal(s.T(), &User{}, user)
}

func TestHashPassword(t *testing.T) {
	hashedPassword, _ := hashPassword("password")
	result := bcrypt.CompareHashAndPassword(hashedPassword, []byte("password"))

	assert.Equal(t, nil, result)
}

func TestValidatePasswordWithValidPassword(t *testing.T) {
	hashedPassword, _ := hashPassword("password")
	result := ValidatePassword(string(hashedPassword), "password")

	assert.Equal(t, nil, result)
}

func TestValidatePasswordWithInvalidPassword(t *testing.T) {
	hashedPassword, _ := hashPassword("password")
	result := ValidatePassword(string(hashedPassword), "drowssap")

	assert.NotEqual(t, nil, result)
}
