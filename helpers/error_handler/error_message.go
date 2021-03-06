package helpers

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

const (
	commonFieldError        = "%s is not valid"
	emailFormatError        = "invalid email format"
	emailDuplicateError     = "email already exists"
	invalidInputError       = "invalid input"
	KeywordNotFoundError    = "keyword not found"
	minError                = "%s must be longer than %s"
	passwordEqError         = "passwords do not match"
	requiredError           = "%s is required"
	somethingWentWrongError = "something went wrong, please try again"

	foreignKeyErrorCode   = "23503"
	duplicateErrorCode    = "23505"
	invalidInputErrorCode = "22P02"
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
	pgErr, isPgErr := err.(*pgconn.PgError)

	if isPgErr {
		switch pgErr.Code {
		case foreignKeyErrorCode:
			return errors.New(somethingWentWrongError)
		case duplicateErrorCode:
			return errors.New(emailDuplicateError)
		case invalidInputErrorCode:
			return errors.New(invalidInputError)
		}
		return errors.New(pgErr.Error())
	} else {
		switch err {
		case gorm.ErrRecordNotFound:
			return errors.New(KeywordNotFoundError)
		}
		return errors.New(somethingWentWrongError)
	}
}
