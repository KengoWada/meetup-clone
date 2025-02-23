package validate

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func dateValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse("02/01/2006", fl.Field().String())
	return err != nil
}
