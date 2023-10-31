package utils

import (
	"github.com/go-playground/validator/v10"
	"unicode"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidateStruct(userInput interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(userInput)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func VerifyPassword(s string) bool {
	var number, upper, special bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			continue
		default:
			return false
		}
		if number && upper && special {
			return true
		}
	}
	return false
}
