package roles

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
)

// GetOrganizationRole godoc
//
//	@Summary		Get an organization role
//	@Description	Get an organization role
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int					true	"orgID to associate role to"
//	@Param			roleID	path		int					true	"roleID to fetch"
//	@Success		200		{object}	models.SimpleRole	"role successfully fetched"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID}/roles/{roleID} [get]
func (h *Handler) getRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	role := ctx.Value(internal.RoleCtx).(*models.Role)

	simpleRole := models.SimpleRole{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
	}

	response.SuccessResponseOK(w, "", simpleRole)
}

// GetOrganizationRoles godoc
//
//	@Summary		Get an organization roles
//	@Description	Get an organization roles
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int					true	"id of the org whose roles to fetch"
//	@Success		200		{object}	[]models.SimpleRole	"org roles successfully fetched"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID}/roles [get]
func (h *Handler) getOrganizationRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization, _ := ctx.Value(internal.OrgCtx).(*models.Organization)

	roles, err := h.store.Roles.GetByOrgID(ctx, organization.ID)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	data := map[string]any{"roles": roles}
	response.SuccessResponseOK(w, "", data)
}

// GetOrganizationPermissions godoc
//
//	@Summary		Get an organization permissions
//	@Description	Get an organization permissions
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int					true	"id of the org whose permissions to fetch"
//	@Success		200		{object}	map[string][]string	"org roles successfully fetched"
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		403		{object}	response.DocsErrorResponseForbidden
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID}/roles/permissions [get]
func (h *Handler) getPermissions(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{"permissions": internal.PermissionsMap}
	response.SuccessResponseOK(w, "", data)
}
