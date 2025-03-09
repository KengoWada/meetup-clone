package validate

import (
	"regexp"
	"slices"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// dateValidator is a custom validation function that checks if the date
// provided matches the required format of `mm/dd/yyyy`. It returns true
// if the date matches the format, and false otherwise.
func dateValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse("01/02/2006", fl.Field().String())
	return err == nil
}

// passwordValidator is a custom validation function that checks if a password
// meets the following criteria:
// 1. Contains at least one uppercase letter.
// 2. Contains at least one lowercase letter.
// 3. Contains at least one number.
// 4. Contains at least one special character from a predefined set of symbols.
//
// It returns true if the password meets all the conditions, and false otherwise.
func passwordValidator(fl validator.FieldLevel) bool {
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

// orgName is a custom validator function for validating organization names.
// It checks if the field value contains only letters, numbers, and spaces.
func orgNameValidator(fl validator.FieldLevel) bool {
	r, _ := regexp.Compile(`^[A-Za-z0-9 ]+$`)
	return r.MatchString(fl.Field().String())
}
