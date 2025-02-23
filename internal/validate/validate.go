package validate

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	Validate.RegisterValidation("is_date", dateValidator)
	Validate.RegisterValidation("is_password", passwordValidator)
}
