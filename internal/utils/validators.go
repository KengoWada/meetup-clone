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

// ValidatePayload validates the provided payload using custom validation rules. It returns a
// map of field names to error messages for any fields that fail validation, along with an error
// if validation fails or if there's an issue processing the payload.
//
// Parameters:
//   - payload: The data to be validated, typically a struct or map that contains the fields to
//     be checked according to custom validation rules.
//   - errorMessages: A map of custom error messages for specific fields that fail validation.
//     This map will be used to generate specific error messages for the fields that don't meet
//     the validation criteria.
//
// Returns:
//   - A map where keys are the names of fields that failed validation and the values are
//     the corresponding error messages. This map is empty if no fields fail validation.
//   - An error if there is an issue during validation (e.g., if the payload is malformed,
//     or if a validation function fails).
//
// Example:
//
//	payload := struct {
//	  Name string `validate:"required"`
//	  Age  int    `validate:"min=18"`
//	}
//	errorMessages := map[string]string{"Name": "Name is required", "Age": "Age must be at least 18"}
//	result, err := ValidatePayload(payload, errorMessages)
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

// generateErrorMessages generates a map of field-specific error messages based on the provided
// error and custom error messages. It processes the validation errors and returns a map where
// the keys are the field names and the values are the error messages.
//
// Parameters:
//   - err: The error object that contains the validation errors. This error is from validate
//     that includes field-level error details.
//   - errorMessages: A map of custom error messages for specific fields. This map is used to
//     customize the error message for each field when validation fails.
//
// Returns:
//   - A map where keys are the names of fields that failed validation, and the values are
//     the corresponding error messages.
//   - An error if there is an issue processing the input error or generating the error messages.
//
// Example:
//
//	err := validate.Struct(payload) // returns validation error
//	errorMessages := map[string]string{"Name": "Name is required", "Age": "Age must be at least 18"}
//	result, err := generateErrorMessages(err, errorMessages)
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

// firstLetterToLower converts the first letter of a string to lowercase, leaving
// the rest of the string unchanged. If the string is empty, it returns the string as is.
//
// Parameters:
//   - s: The string whose first letter needs to be converted to lowercase.
//
// Returns:
//   - The string with the first letter converted to lowercase. If the string is empty,
//     the original string is returned unchanged.
//
// Example:
//
//	firstLetterToLower("Hello") // returns "hello"
//	firstLetterToLower("world") // returns "world"
//	firstLetterToLower("")      // returns ""
func firstLetterToLower(s string) string {
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
