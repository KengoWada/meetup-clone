package auth

import (
	"github.com/KengoWada/meetup-clone/internal/validate"
)

var (
	registerUserPayloadErrors = validate.FieldErrorMessages{
		"email":       validate.TagErrorsEmail,
		"password":    validate.TagErrorsPassword,
		"username":    validate.TagErrorsUsername,
		"profilePic":  validate.TagErrorsURL,
		"dateOfBirth": validate.TagErrorsDOB,
	}

	resendVerificationEmailPayloadErrors = validate.FieldErrorMessages{
		"email": validate.TagErrorsEmail,
	}

	resetUserPasswordPayloadErrors = validate.FieldErrorMessages{
		"password": validate.TagErrorsPassword,
	}

	passwordResetRequestPayloadErrors = validate.FieldErrorMessages{
		"email": validate.TagErrorsEmail,
	}
)
