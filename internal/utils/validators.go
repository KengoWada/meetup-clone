package utils

import (
	"errors"
	"slices"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type TagErrorMessages map[string]string
type FieldErrorMessages map[string]TagErrorMessages

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	Validate.RegisterValidation("is_date", DateValidator)
	Validate.RegisterValidation("is_password", PasswordValidator)
}

func DateValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse("02/01/2006", fl.Field().String())
	return err != nil
}

func PasswordValidator(fl validator.FieldLevel) bool {
	var (
		hasNumber           = false
		hasUpperCase        = false
		hasLowerCase        = false
		hasSpecialCharacter = false
		specialCharacters   = []rune{
			'!', '@', '#', '$', '%', '^', '&',
			'*', '(', ')', '-', '_', '+', '=',
			',', '.', '?', ' ', '|', '\\', '/',
			'<', '>', '[', ']', '{', '}',
		}
	)

	for _, char := range fl.Field().String() {
		if unicode.IsUpper(char) {
			hasUpperCase = true
			continue
		}

		if unicode.IsLower(char) {
			hasLowerCase = true
			continue
		}

		if unicode.IsNumber(char) {
			hasNumber = true
			continue
		}

		if slices.Contains(specialCharacters, char) {
			hasSpecialCharacter = true
		}
	}

	return hasNumber && hasUpperCase && hasLowerCase && hasSpecialCharacter
}

func GenerateErrorMessages(err error, errorMessages FieldErrorMessages) (map[string]string, error) {
	var validateErrs validator.ValidationErrors
	var response = make(map[string]string)

	if errors.As(err, &validateErrs) {
		for _, err := range validateErrs {
			field := firstLetterToLower(err.Field())
			tag := err.ActualTag()

			if tag == "required" {
				response[field] = "Field is required"
				continue
			}

			message, ok := errorMessages[field][tag]
			if !ok {
				message = "Invalid field"
			}

			response[field] = message
		}
	} else {
		return nil, errors.New("internal server error")
	}

	return response, nil
}

func firstLetterToLower(s string) string {
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
