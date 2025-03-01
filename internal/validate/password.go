package validate

import (
	"slices"
	"unicode"

	"github.com/go-playground/validator/v10"
)

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
