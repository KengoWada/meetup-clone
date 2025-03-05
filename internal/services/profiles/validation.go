package profiles

import "github.com/KengoWada/meetup-clone/internal/validate"

var (
	updateUserPayloadErrors = validate.FieldErrorMessages{
		"email":       validate.TagErrorsEmail,
		"password":    validate.TagErrorsPassword,
		"username":    validate.TagErrorsUsername,
		"profilePic":  validate.TagErrorsURL,
		"dateOfBirth": validate.TagErrorsDOB,
	}
)
