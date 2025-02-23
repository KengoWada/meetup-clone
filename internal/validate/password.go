package validate

import (
	"slices"
	"unicode"

	"github.com/go-playground/validator/v10"
)

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
