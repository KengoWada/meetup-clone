package auth

import (
	"net/http"
	"strconv"

	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/go-chi/chi/v5"
)

// DeactivateUser godoc
//
//	@Summary		Deactivate a user
//	@Description	Deactivate a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int									true	"userID to deactivate"
//	@Success		200		{object}	response.DocsResponseMessageOnly	"user successfully deactivated"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		401		{object}	response.DocsErrorResponseUnauthorized
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/auth/users/{userID}/deactivate [patch]
func (h *Handler) deactivateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		errorMessage := response.ErrorResponse{Message: "Invalid user ID"}
		response.ErrorResponseBadRequest(w, r, err, errorMessage)
		return
	}

	user, err := h.store.Users.GetByIDIcludeDeleted(r.Context(), int(userID))
	if err != nil {
		errorMessage := response.ErrorResponse{Message: "Invalid user ID"}
		response.ErrorResponseBadRequest(w, r, err, errorMessage)
		return
	}

	if user.IsDeactivated() {
		errorMessage := response.ErrorResponse{Message: "User is already deactivated"}
		response.ErrorResponseBadRequest(w, r, err, errorMessage)
		return
	}

	err = h.store.Users.Deactivate(r.Context(), user)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			errorMessage := response.ErrorResponse{Message: "Try again later"}
			response.ErrorResponseBadRequest(w, r, err, errorMessage)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	if cfg.CacheConfig.Enabled {
		err := h.cacheStore.Users.Delete(userID)
		if err != nil {
			// TODO: Log error but don't return 500
		}
	}

	response.SuccessResponseOK(w, "Done", nil)
}
