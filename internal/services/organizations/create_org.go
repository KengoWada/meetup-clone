package organizations

import (
	"net/http"
	"strings"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
)

type createOrganizationPayload struct {
	Name        utils.TrimString `json:"name" validate:"required,max=100,is_org_name"`
	Description utils.TrimString `json:"description" validate:"required"`
	ProfilePic  utils.TrimString `json:"profilePic" validate:"required,http_url"`
}

type orgResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProfilePic  string `json:"profilePic"`
	CreatedAt   string `json:"createdAt"`
}

func (h *Handler) createOrganization(w http.ResponseWriter, r *http.Request) {
	var payload createOrganizationPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "json: unknown field") {
			response.ErrorResponseUnknownField(w, r, err)
			return
		}
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if errResponse, err := validate.ValidatePayload(payload, createOrganizationPayloadErrors); err != nil {
		switch err {
		case validate.ErrFailedValidation:
			errorMessage := response.NewValidationErrorResponse(errResponse)
			response.ErrorResponseBadRequest(w, r, err, errorMessage)

		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	user, _ := r.Context().Value(internal.UserCtx).(*models.User)

	organization := models.Organization{
		Name:        string(payload.Name),
		Description: string(payload.Description),
		ProfilePic:  string(payload.ProfilePic),
	}
	role := models.Role{
		Name:        "sudo",
		Description: "This is role has all permissions",
		Permissions: internal.Permissions,
	}
	member := models.OrganizationMember{
		UserProfileID: user.UserProfile.ID,
	}

	if err := h.store.Organizations.Create(r.Context(), &organization, &role, &member); err != nil {
		switch err {
		case store.ErrDuplicateOrgName:
			errorMessage := response.NewValidationErrorResponse(response.ErrorsResponse{"name": err.Error()})
			response.ErrorResponseBadRequest(w, r, err, errorMessage)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	orgData := orgResponse{
		ID:          organization.ID,
		Name:        organization.Name,
		Description: organization.Description,
		ProfilePic:  organization.ProfilePic,
		CreatedAt:   organization.CreatedAt,
	}

	response.SuccessResponseCreated(w, "", orgData)
}
