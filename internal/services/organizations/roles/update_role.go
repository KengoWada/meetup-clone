package roles

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/logger"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
	"github.com/go-chi/chi/v5"
)

type updateRolePayload struct {
	Name        utils.TrimString `json:"name" validate:"required,max=100"`
	Description utils.TrimString `json:"description" validate:"required"`
	Permissions []string         `json:"permissions" validate:"required,gt=0,unique,dive,is_permission"`
}

// UpdateOrganizationRole godoc
//
//	@Summary		Update an organization role
//	@Description	Update an organization role
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int					true	"orgID to associate role to"
//	@Param			roleID	path		int					true	"roleID to update"
//	@Param			payload	body		updateRolePayload	true	"update organization role payload"
//	@Success		200		{object}	models.SimpleRole	"role successfully updated"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		403		{object}	response.DocsErrorResponseForbidden
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID}/roles/{roleID} [put]
func (h *Handler) updateRole(w http.ResponseWriter, r *http.Request) {
	var payload updateRolePayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		response.ErrorResponseInvalidJSON(w, r, err)
		return
	}

	if errorMessages, err := validate.ValidatePayload(payload, createRolePayloadErrors); err != nil {
		switch err {
		case validate.ErrFailedValidation:
			errorResponse := response.NewValidationErrorResponse(errorMessages)
			response.ErrorResponseBadRequest(w, r, err, errorResponse)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	ctx := r.Context()
	orgID, _ := strconv.ParseInt(chi.URLParam(r, "orgID"), 10, 64)
	role, _ := ctx.Value(internal.RoleCtx).(*models.Role)

	nameExists, err := roleNameExists(ctx, h.store, string(payload.Name), orgID, role.ID)
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

	role.Name = string(payload.Name)
	role.Description = string(payload.Description)
	role.Permissions = payload.Permissions

	if err := h.store.Roles.Update(ctx, role); err != nil {
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
		if err := h.cacheStore.Roles.Delete(role.ID); err != nil {
			logger.ErrLoggerCache(r, err)
		}
	}

	simpleRole := models.SimpleRole{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
	}
	response.SuccessResponseOK(w, "", simpleRole)
}

func (h *Handler) deleteRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	role, _ := ctx.Value(internal.RoleCtx).(*models.Role)

	fields := []string{"role_id", "org_id"}
	values := []any{role.ID, role.OrganizationID}
	_, err := h.store.OrganizationMembers.Get(ctx, false, fields, values)
	if err != nil && err != store.ErrNotFound {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if err != store.ErrNotFound {
		err := errors.New("role has active users attached to it")
		errorResponse := response.ErrorResponse{Message: "Role is assigned to active users. Please reassign them before deleting."}
		response.ErrorResponseBadRequest(w, r, err, errorResponse)
		return
	}

	if err := h.store.Roles.SoftDelete(ctx, role); err != nil {
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
		if err := h.cacheStore.Roles.Delete(role.ID); err != nil {
			logger.ErrLoggerCache(r, err)
		}
	}

	response.SuccessResponseOK(w, "Done", nil)
}

func roleNameExists(ctx context.Context, appStore store.Store, name string, orgID, roleID int64) (bool, error) {
	fields := []string{"name", "org_id"}
	values := []any{name, orgID}
	dbRole, err := appStore.Roles.Get(ctx, false, fields, values)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return false, nil
		default:
			return false, err
		}
	}

	if dbRole.ID == roleID {
		return false, nil
	}

	return true, nil
}
