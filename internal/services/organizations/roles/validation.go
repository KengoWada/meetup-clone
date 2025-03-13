package roles

import "github.com/KengoWada/meetup-clone/internal/validate"

var (
	createRolePayloadErrors = validate.FieldErrorMessages{
		"name": validate.TagErrorMessages{
			"max": "Organization name should be at most 100 character",
		},
		"description": validate.TagErrorMessages{},
		"permissions": validate.TagErrorMessages{
			"is_permission": "Invalid permission sent",
			"unique":        "Duplicate permissions are not allowed",
		},
	}
)
