package validate

import (
	"errors"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type TagErrorMessages map[string]string

type FieldErrorMessages map[string]TagErrorMessages

var (
	ErrFailedValidation = errors.New("payload failed validation")

	TagErrorsEmail = TagErrorMessages{"email": "Invalid email address provided"}

	TagErrorsPassword = TagErrorMessages{
		"min":         "Password must have at least 10 characters",
		"max":         "Password must have at most 72 characters",
		"is_password": "Password must contain a number, lower case character, upper case character and one of the special symbols(including space) !@#$%^&*()-_+=,.?|\\/<>[]{}",
	}

	TagErrorsDOB = TagErrorMessages{"is_date": "Invalid date format. mm/dd/yyyy"}

	TagErrorsUsername = TagErrorMessages{
		"min": "Username must have at least 3 characters",
		"max": "Username must have at most 100 characters",
	}

	TagErrorsURL = TagErrorMessages{"http_url": "Invalid URL format"}
)

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
	err := Validate.Struct(payload)
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
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
