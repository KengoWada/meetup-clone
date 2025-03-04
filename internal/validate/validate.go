// Package validate contains custom validation functions that can be used
// with the validator package to perform complex or application-specific validations.
package validate

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once
	Validate *validator.Validate
)

func Get() *validator.Validate {
	once.Do(func() {
		Validate = validator.New(validator.WithRequiredStructEnabled())
		Validate.RegisterValidation("is_date", dateValidator)
		Validate.RegisterValidation("is_password", passwordValidator)
	})

	return Validate
}
