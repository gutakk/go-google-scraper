package error_handler

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

	assert.Equal(t, "invalid email format", result.Error())
}

func TestValidationErrorMessageForEqfieldTag(t *testing.T) {
	fieldError := &mockFieldError{tag: "eqfield"}
	result := ValidationErrorMessage(fieldError)

	assert.Equal(t, "passwords do not match", result.Error())
}

func TestValidationErrorMessageForMinTag(t *testing.T) {
	fieldError := &mockFieldError{
		tag:   "min",
		field: "password",
		param: "6",
	}
	result := ValidationErrorMessage(fieldError)

	assert.Equal(t, "password must be longer than 6", result.Error())
}

func TestValidationErrorMessageForRequiredTag(t *testing.T) {
	fieldError := &mockFieldError{
		tag:   "required",
		field: "password",
	}
	result := ValidationErrorMessage(fieldError)

	assert.Equal(t, "password is required", result.Error())
}

func TestValidationErrorMessageForDefaultCaseTag(t *testing.T) {
	fieldError := &mockFieldError{
		tag:   "test",
		field: "password",
	}
	result := ValidationErrorMessage(fieldError)

	assert.Equal(t, "password is not valid", result.Error())
}

func TestDatabaseErrorMessageForDuplicateEmail(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505"}
	result := DatabaseErrorMessage(pgErr)

	assert.Equal(t, "email already exists", result.Error())
}

func TestDatabaseErrorMessageForDefaultCode(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:     "23506",
		Severity: "ERROR",
		Message:  "Test",
	}
	result := DatabaseErrorMessage(pgErr)

	assert.Equal(t, "ERROR: Test (SQLSTATE 23506)", result.Error())
}

func TestDatabaseErrorMessageForNonPgErrorType(t *testing.T) {
	err := errors.New("custom error")
	result := DatabaseErrorMessage(err)

	assert.Equal(t, "something went wrong, please try again", result.Error())
}
