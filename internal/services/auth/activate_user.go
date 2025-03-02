package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
)

type activateUserPayload struct {
	Token string `json:"token"`
}

// ActivateUser godoc
//
//	@Summary		Activate a user
//	@Description	Activate a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		activateUserPayload					true	"activate user payload"
//	@Success		200		{object}	response.DocsResponseMessageOnly	"user successfully activated"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		422		{object}	response.DocsResponseMessageOnly
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security
//	@Router	/auth/activate [patch]
func (h *Handler) activateUser(w http.ResponseWriter, r *http.Request) {
	var payload activateUserPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "json: unknown field") {
			response.ErrorResponseUnknownField(w, r, err)
			return
		}
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	timedToken, err := utils.ValidateToken(payload.Token, []byte(h.config.SecretKey), time.Minute*30)
	if err != nil {
		switch err {
		case utils.ErrExpiredToken:
			errorMessage := response.ErrorResponse{Message: "Activation token has exipred"}
			response.ErrorResponseUnprocessableEntity(w, r, err, errorMessage)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	ctx := r.Context()
	errorMessage := response.ErrorResponse{Message: "Activation token is invalid"}

	user, err := h.store.Users.GetByEmail(ctx, timedToken.Body)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			response.ErrorResponseBadRequest(w, r, err, errorMessage)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	if user.IsActivated() || user.IsDeactivated() {
		response.ErrorResponseBadRequest(w, r, err, errorMessage)
		return
	}

	err = h.store.Users.Activate(ctx, user)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	response.SuccessResponseOK(w, "Email successfully verified", nil)
}
