package helpers

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

const (
	commonFieldError    = "%s is not valid"
	emailFormatError    = "Invalid email format"
	emailDuplicateError = "Email already exists"
	minError            = "%s must be longer than %s"
	passwordEqError     = "Passwords do not match"
	requiredError       = "%s is required"

	duplicateErrorCode = "23505"
)

// interface from github.com/go-playground/validator/v10
type FieldError interface {
	Tag() string
	Field() string
	Param() string
}

func ValidationErrorMessage(err FieldError) error {
	switch err.Tag() {
	case "email":
		return errors.New(emailFormatError)
	case "eqfield":
		return errors.New(passwordEqError)
	case "min":
		return fmt.Errorf(minError, err.Field(), err.Param())
	case "required":
		return fmt.Errorf(requiredError, err.Field())
	}
	return fmt.Errorf(commonFieldError, err.Field())
}

// TODO: Improve later as this feel brittle
func DatabaseErrorMessage(err error) error {
	pgErr := err.(*pgconn.PgError)

	switch pgErr.Code {
	case duplicateErrorCode:
		return errors.New(emailDuplicateError)
	}
	return errors.New(pgErr.Error())
}
