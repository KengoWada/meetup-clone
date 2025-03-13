package organizations

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/logger"
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
//	@Param			orgID	path		int							true	"orgID to update"
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
			logger.ErrLoggerCache(r, err)
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

// DeleteOrganization godoc
//
//	@Summary		Delete an organization
//	@Description	Delete an organization
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int										true	"orgID to delete"
//	@Success		200		{object}	response.DocsSuccessResponseDoneMessage	"organization successfully deleted"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		403		{object}	response.DocsErrorResponseForbidden
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID} [delete]
func (h *Handler) deleteOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization, _ := ctx.Value(internal.OrgCtx).(*models.Organization)

	err := h.store.Organizations.SoftDelete(ctx, organization)
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
			logger.ErrLoggerCache(r, err)
		}
	}

	response.SuccessResponseOK(w, "Done", nil)
}

// DeactivateOrganization godoc
//
//	@Summary		Deactivate an organization(staff or admin)
//	@Description	Deactivate an organization(staff or admin)
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int										true	"orgID to deactivate"
//	@Success		200		{object}	response.DocsSuccessResponseDoneMessage	"organization successfully deactivated"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		403		{object}	response.DocsErrorResponseForbidden
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID} [patch]
func (h *Handler) deactivateOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization, _ := ctx.Value(internal.OrgCtx).(*models.Organization)

	err := h.store.Organizations.Deactivate(ctx, organization)
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
			logger.ErrLoggerCache(r, err)
		}
	}

	response.SuccessResponseOK(w, "Done", nil)
}
