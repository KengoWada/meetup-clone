package organizations

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
)

type updateOrganizationPyload struct {
	Name        utils.TrimString `json:"name" validate:"required,max=100,is_org_name"`
	Description utils.TrimString `json:"description" validate:"required"`
	ProfilePic  utils.TrimString `json:"profilePic" validate:"required,http_url"`
}

// UpdateOrganization godoc
//
//	@Summary		Update an organization
//	@Description	Update an organization
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int							true	"orgID to deactivate"
//	@Param			payload	body		updateOrganizationPyload	true	"update organization payload"
//	@Success		200		{object}	orgResponse					"organization successfully updated"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		403		{object}	response.DocsErrorResponseForbidden
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID} [put]
func (h *Handler) updateOrganization(w http.ResponseWriter, r *http.Request) {
	var payload updateOrganizationPyload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		response.ErrorResponseInvalidJSON(w, r, err)
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

	ctx := r.Context()
	organization, _ := ctx.Value(internal.OrgCtx).(*models.Organization)

	organization.Name = string(payload.Name)
	organization.Description = string(payload.Description)
	organization.ProfilePic = string(payload.ProfilePic)

	err = h.store.Organizations.Update(ctx, organization)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			res := response.ErrorResponse{Message: "Try again later"}
			response.ErrorResponseBadRequest(w, r, err, res)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	if cfg.CacheConfig.Enabled {
		if err := h.cacheStore.Organizations.Delete(organization.ID); err != nil {
			// TODO: log cache error
		}
	}

	orgData := orgResponse{
		ID:          organization.ID,
		Name:        organization.Name,
		Description: organization.Description,
		ProfilePic:  organization.ProfilePic,
		CreatedAt:   organization.CreatedAt,
	}
	response.SuccessResponseOK(w, "", orgData)
}
