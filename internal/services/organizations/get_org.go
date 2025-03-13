package organizations

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
)

// GetOrganization godoc
//
//	@Summary		Get an organization
//	@Description	Get an organization
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int			true	"orgID to fetch"
//	@Success		200		{object}	orgResponse	"organization successfully fetched"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		403		{object}	response.DocsErrorResponseForbidden
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID} [get]
func (h *Handler) getOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization, _ := ctx.Value(internal.OrgCtx).(*models.Organization)

	orgData := orgResponse{
		ID:          organization.ID,
		Name:        organization.Name,
		Description: organization.Description,
		ProfilePic:  organization.ProfilePic,
		CreatedAt:   organization.CreatedAt,
	}
	response.SuccessResponseOK(w, "", orgData)
}

// GetOrganizations godoc
//
//	@Summary		Get a users organizations
//	@Description	Get a users organizations
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.SimpleOrganization	"organizations successfully fetched"
//	@Failure		400	{object}	response.DocsErrorResponse
//	@Failure		401	{object}	response.DocsErrorResponseUnauthorized
//	@Failure		500	{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations [get]
func (h *Handler) getUsersOrganizations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, _ := ctx.Value(internal.UserCtx).(*models.User)

	organizations, err := h.store.Organizations.GetByUserID(ctx, user.UserProfile.ID)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	response.SuccessResponseOK(w, "", map[string]any{"organizations": organizations})
}
