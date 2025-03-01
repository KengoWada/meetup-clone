package validate

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// dateValidator is a custom validation function that checks if the date
// provided matches the required format of `mm/dd/yyyy`. It returns true
// if the date matches the format, and false otherwise.
func dateValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse("01/02/2006", fl.Field().String())
	return err == nil
}
