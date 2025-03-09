package organizations

import "github.com/KengoWada/meetup-clone/internal/validate"

var (
	createOrganizationPayloadErrors = validate.FieldErrorMessages{
		"name": validate.TagErrorMessages{
			"max":         "Organization name should be at most 100 character",
			"is_org_name": "Organization name can only include alphanumeric characters and spaces",
		},
		"description": validate.TagErrorMessages{},
		"profilePic":  validate.TagErrorsURL,
	}
)
