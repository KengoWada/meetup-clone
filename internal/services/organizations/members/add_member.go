package members

import (
	"context"
	"errors"
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
)

type inviteMemberPayload struct {
	Email  utils.TrimString `json:"email" validate:"required,email"`
	RoleID int64            `json:"roleId" validate:"required"`
}

// InviteOrganizationMember godoc
//
//	@Summary		Invite an organization member
//	@Description	Invite an organization member
//	@Tags			members
//	@Accept			json
//	@Produce		json
//	@Param			orgID	path		int									true	"orgID to update"
//	@Param			payload	body		inviteMemberPayload					true	"invite organization member payload"
//	@Success		201		{object}	response.DocsResponseMessageOnly	"invite sent successfully"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		403		{object}	response.DocsErrorResponseForbidden
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgID}/members [post]
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

	ctx := r.Context()
	organization, _ := ctx.Value(internal.OrgCtx).(*models.Organization)

	exists, err := roleExists(ctx, h.store, payload.RoleID, organization.ID)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if !exists {
		err := errors.New("invalid role id")
		errorResponse := response.NewValidationErrorResponse(response.ErrorsResponse{"roleId": "invalid role id"})
		response.ErrorResponseBadRequest(w, r, err, errorResponse)
		return
	}

	fields := []string{"email", "is_active"}
	values := []any{string(payload.Email), true}
	user, err := h.store.Users.GetWithProfile(ctx, false, fields, values)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			response.SuccessResponseCreated(w, "Invite sent", nil)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	fields = []string{"user_id", "org_id"}
	values = []any{user.UserProfile.ID, organization.ID}
	_, err = h.store.OrganizationInvites.Get(ctx, false, fields, values)
	if err != nil && err != store.ErrNotFound {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if err == nil {
		response.SuccessResponseCreated(w, "Invite sent", nil)
		return
	}

	_, err = h.store.OrganizationMembers.Get(ctx, false, fields, values)
	if err != nil && err != store.ErrNotFound {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if err == nil {
		err := errors.New("tried to invite team member")
		errorResponse := response.ErrorResponse{Message: "User is already a member"}
		response.ErrorResponseBadRequest(w, r, err, errorResponse)
		return
	}

	invite := &models.OrganizationInvite{
		OrganizationID: organization.ID,
		UserProfileID:  user.UserProfile.ID,
		RoleID:         payload.RoleID,
	}
	if err := h.store.OrganizationInvites.Create(ctx, invite); err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	response.SuccessResponseCreated(w, "Invite sent", nil)
}

func roleExists(ctx context.Context, appStore store.Store, roleID, orgID int64) (bool, error) {
	fields := []string{"id", "org_id"}
	values := []any{roleID, orgID}
	_, err := appStore.Roles.Get(ctx, false, fields, values)
	if err != nil && err != store.ErrNotFound {
		return false, err
	}

	if err == nil {
		return true, nil
	}

	return false, nil
}
