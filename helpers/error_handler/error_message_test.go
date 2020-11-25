package helpers

import (
	"errors"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"gopkg.in/go-playground/assert.v1"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

type mockFieldError struct {
	tag   string
	field string
	param string
}

func (m *mockFieldError) Tag() string {
	return m.tag
}

func (m *mockFieldError) Field() string {
	return m.field
}

func (m *mockFieldError) Param() string {
	return m.param
}

func TestValidationErrorMessageForEmailTag(t *testing.T) {
	fieldError := &mockFieldError{tag: "email"}
	result := ValidationErrorMessage(fieldError)

	assert.Equal(t, errors.New("Invalid email format"), result)
}

func TestValidationErrorMessageForEqfieldTag(t *testing.T) {
	fieldError := &mockFieldError{tag: "eqfield"}
	result := ValidationErrorMessage(fieldError)

	assert.Equal(t, errors.New("Passwords do not match"), result)
}

func TestValidationErrorMessageForMinTag(t *testing.T) {
	fieldError := &mockFieldError{
		tag:   "min",
		field: "password",
		param: "6",
	}
	result := ValidationErrorMessage(fieldError)

	assert.Equal(t, errors.New("password must be longer than 6"), result)
}

func TestValidationErrorMessageForRequiredTag(t *testing.T) {
	fieldError := &mockFieldError{
		tag:   "required",
		field: "password",
	}
	result := ValidationErrorMessage(fieldError)

	assert.Equal(t, errors.New("password is required"), result)
}

func TestValidationErrorMessageForDefaultCaseTag(t *testing.T) {
	fieldError := &mockFieldError{
		tag:   "test",
		field: "password",
	}
	result := ValidationErrorMessage(fieldError)

	assert.Equal(t, errors.New("password is not valid"), result)
}

func TestDatabaseErrorMessageForDuplicateEmail(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505"}
	result := DatabaseErrorMessage(pgErr)

	assert.Equal(t, errors.New("Email already exists"), result)
}

func TestDatabaseErrorMessageForDefaultCode(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:     "23506",
		Severity: "ERROR",
		Message:  "Test",
	}
	result := DatabaseErrorMessage(pgErr)

	assert.Equal(t, errors.New("ERROR: Test (SQLSTATE 23506)"), result)
}
