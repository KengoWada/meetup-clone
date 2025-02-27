package utils

import (
	"errors"
	"unicode"

	"github.com/KengoWada/meetup-clone/internal/validate"
	"github.com/go-playground/validator/v10"
)

type TagErrorMessages map[string]string
type FieldErrorMessages map[string]TagErrorMessages

var ErrFailedValidation = errors.New("payload failed validation")

func ValidatePayload(payload any, errorMessages FieldErrorMessages) (map[string]string, error) {
	err := validate.Validate.Struct(payload)
	if err == nil {
		return nil, nil
	}

	errResponse, err := generateErrorMessages(err, errorMessages)
	if err != nil {
		return nil, err
	}

	return errResponse, ErrFailedValidation
}

func generateErrorMessages(err error, errorMessages FieldErrorMessages) (map[string]string, error) {
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
