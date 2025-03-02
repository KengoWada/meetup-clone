package auth

import "github.com/KengoWada/meetup-clone/internal/utils"

var (
	emailErrors = utils.TagErrorMessages{"email": "Invalid email address provided"}

	registerUserPayloadErrors = utils.FieldErrorMessages{
		"email": emailErrors,
		"password": {
			"min":         "Password must have at least 10 characters",
			"max":         "Password must have at most 72 characters",
			"is_password": "Password must contain a number, lower case character, upper case character and one of the special symbols(including space) !@#$%^&*()-_+=,.?|\\/<>[]{}",
		},
		"username": {
			"min": "Username must have at least 3 characters",
			"max": "Username must have at most 100 characters",
		},
		"profilePic":  {"http_url": "Invalid URL format"},
		"dateOfBirth": {"is_date": "Invalid date format. mm/dd/yyyy"},
	}

	resendVerificationEmailPayloadErrors = utils.FieldErrorMessages{
		"email": emailErrors,
	}
)
