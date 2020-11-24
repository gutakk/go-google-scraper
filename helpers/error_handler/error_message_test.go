package helpers

import (
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

func TestValidationErrorToTextForEmailTag(t *testing.T) {
	fieldError := &mockFieldError{tag: "email"}
	result := ValidationErrorToText(fieldError)

	assert.Equal(t, "Invalid email format", result)
}

func TestValidationErrorToTextForEqfieldTag(t *testing.T) {
	fieldError := &mockFieldError{tag: "eqfield"}
	result := ValidationErrorToText(fieldError)

	assert.Equal(t, "Passwords do not match", result)
}

func TestValidationErrorToTextForMinTag(t *testing.T) {
	fieldError := &mockFieldError{
		tag:   "min",
		field: "password",
		param: "6",
	}
	result := ValidationErrorToText(fieldError)

	assert.Equal(t, "password must be longer than 6", result)
}

func TestValidationErrorToTextForRequiredTag(t *testing.T) {
	fieldError := &mockFieldError{
		tag:   "required",
		field: "password",
	}
	result := ValidationErrorToText(fieldError)

	assert.Equal(t, "password is required", result)
}

func TestValidationErrorToTextForDefaultCaseTag(t *testing.T) {
	fieldError := &mockFieldError{
		tag:   "test",
		field: "password",
	}
	result := ValidationErrorToText(fieldError)

	assert.Equal(t, "password is not valid", result)
}

func TestDatabaseErrorToTextForDuplicateEmail(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505"}
	result := DatabaseErrorToText(pgErr)

	assert.Equal(t, "Email already exists", result)
}

func TestDatabaseErrorToTextForDefaultCode(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:     "23506",
		Severity: "ERROR",
		Message:  "Test",
	}
	result := DatabaseErrorToText(pgErr)

	assert.Equal(t, "ERROR: Test (SQLSTATE 23506)", result)
}
