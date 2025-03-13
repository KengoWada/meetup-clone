package roles

import (
	"errors"
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
)

type createRolePayload struct {
	Name        utils.TrimString `json:"name" validate:"required,max=100"`
	Description utils.TrimString `json:"description" validate:"required"`
	Permissions []string         `json:"permissions" validate:"required,is_permission,unique"`
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
//	@Success		201		{object}	SimpleRole			"organization role successfully created"
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

	fields := []string{"name", "org_id"}
	values := []any{string(payload.Name), organization.ID}
	_, err = h.store.Roles.Get(ctx, false, fields, values)
	if err != nil && err != store.ErrNotFound {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if err == nil {
		err := errors.New("duplicate role name for organization")
		errorMessage := response.NewValidationErrorResponse(
			response.ErrorsResponse{"name": "name already exists"})
		response.ErrorResponseBadRequest(w, r, err, errorMessage)
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
