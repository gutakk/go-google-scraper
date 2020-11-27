package helpers

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
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
		return errors.New("Invalid email format")
	case "eqfield":
		return errors.New("Passwords do not match")
	case "min":
		return fmt.Errorf("%s must be longer than %s", err.Field(), err.Param())
	case "required":
		return fmt.Errorf("%s is required", err.Field())
	}
	return fmt.Errorf("%s is not valid", err.Field())
}

// TODO: Improve later as this feel brittle
func DatabaseErrorMessage(err error) error {
	pgErr := err.(*pgconn.PgError)

	switch pgErr.Code {
	case "23505":
		return errors.New("Email already exists")
	}
	return errors.New(pgErr.Error())
}
