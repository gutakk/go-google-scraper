package helpers

import (
	"fmt"

	"github.com/jackc/pgconn"
)

// interface from github.com/go-playground/validator/v10
type FieldError interface {
	Tag() string
	Field() string
	Param() string
}

func ValidationErrorToText(err FieldError) string {
	switch err.Tag() {
	case "email":
		return "Invalid email format"
	case "eqfield":
		return "Passwords do not match"
	case "min":
		return fmt.Sprintf("%s must be longer than %s", err.Field(), err.Param())
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	}
	return fmt.Sprintf("%s is not valid", err.Field())
}

// TODO: Improve later as this feel brittle
func DatabaseErrorToText(err error) string {
	pgErr := err.(*pgconn.PgError)

	switch pgErr.Code {
	case "23505":
		return "Email already exists"
	}
	return pgErr.Error()
}
