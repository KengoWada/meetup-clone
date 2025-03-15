package roles

import (
	"errors"
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
)

type createRolePayload struct {
	Name        utils.TrimString `json:"name" validate:"required,max=100"`
	Description utils.TrimString `json:"description" validate:"required"`
	Permissions []string         `json:"permissions" validate:"required,gt=0,unique,dive,is_permission"`
}

// CreateOrganizationRole godoc
//
//	@Summary		Create an organization role
//	@Description	Create an organization role
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int					true	"orgID to associate role to"
//	@Param			payload	body		createRolePayload	true	"create organization role payload"
//	@Success		201		{object}	models.SimpleRole	"organization role successfully created"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		403		{object}	response.DocsErrorResponseForbidden
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID}/roles [post]
func (h *Handler) createRole(w http.ResponseWriter, r *http.Request) {
	var payload createRolePayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		response.ErrorResponseInvalidJSON(w, r, err)
		return
	}

	errorMessages, err := validate.ValidatePayload(payload, createRolePayloadErrors)
	if err != nil {
		switch err {
		case validate.ErrFailedValidation:
			errorMessage := response.NewValidationErrorResponse(errorMessages)
			response.ErrorResponseBadRequest(w, r, err, errorMessage)

		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	ctx := r.Context()
	organization := ctx.Value(internal.OrgCtx).(*models.Organization)

	nameExists, err := roleNameExists(ctx, h.store, string(payload.Name), organization.ID, 0)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if nameExists {
		err := errors.New("duplicate organization role name")
		errorResponse := response.NewValidationErrorResponse(response.ErrorsResponse{"name": "name already exists"})
		response.ErrorResponseBadRequest(w, r, err, errorResponse)
		return
	}

	role := &models.Role{
		Name:           string(payload.Name),
		Description:    string(payload.Description),
		Permissions:    payload.Permissions,
		OrganizationID: organization.ID,
	}

	if err := h.store.Roles.Create(ctx, role); err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	simpleRole := models.SimpleRole{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
	}
	response.SuccessResponseCreated(w, "Done", simpleRole)
}
