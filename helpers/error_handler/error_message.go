package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
)

func ValidationErrorToText(err validator.FieldError) string {
	switch err.Tag() {
	case "min":
		return fmt.Sprintf("%s must be longer than %s", err.Field(), err.Param())
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return "Invalid email format"
	case "eqfield":
		return "Password not match"
	}
	return fmt.Sprintf("%s is not valid", err.Field())
}

func DatabaseErrorToText(err error) string {
	pgErr := err.(*pgconn.PgError)

	switch pgErr.Code {
	case "23505":
		return "Email already exists"
	}
	return pgErr.Error()
}
