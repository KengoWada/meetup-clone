package members

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
)

type inviteMemberPayload struct {
	Email utils.TrimString `json:"email" validate:"required,is_email"`
}

func (h *Handler) inviteMember(w http.ResponseWriter, r *http.Request) {
	var payload inviteMemberPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		response.ErrorResponseInvalidJSON(w, r, err)
		return
	}

	errorMessages, err := validate.ValidatePayload(payload, validate.FieldErrorMessages{"email": validate.TagErrorsEmail})
	if err != nil {
		errorResponse := response.NewValidationErrorResponse(errorMessages)
		response.ErrorResponseBadRequest(w, r, err, errorResponse)
		return
	}
}
