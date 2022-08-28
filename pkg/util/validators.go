package util

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateStruct(i any) error {
	if validate == nil {
		validate = validator.New()
	}

	return validate.Struct(i)
}

func ValidateVar(field any, tag string) error {
	if validate == nil {
		validate = validator.New()
	}

	return validate.Var(field, tag)
}
