package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
)

type activateUserPayload struct {
	Token string `json:"token"`
}

type resendVerificationEmailPayload struct {
	Email string `json:"email" validate:"required,email"`
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

	timedToken, err := utils.ValidateToken(payload.Token, []byte(cfg.SecretKey), time.Minute*30)
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

// ResendVerificationEmail godoc
//
//	@Summary		Resend verification email to user
//	@Description	Resend verification email to user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		resendVerificationEmailPayload		true	"resend verification email payload"
//	@Success		200		{object}	response.DocsResponseMessageOnly	"email sent if account exists"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security
//	@Router	/auth/resend-verification-email [post]
func (h *Handler) resendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	var payload resendVerificationEmailPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "json: unknown field") {
			response.ErrorResponseUnknownField(w, r, err)
			return
		}
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if errResponse, err := validate.ValidatePayload(payload, resendVerificationEmailPayloadErrors); err != nil {
		switch err {
		case validate.ErrFailedValidation:
			errorMessage := response.NewValidationErrorResponse(errResponse)
			response.ErrorResponseBadRequest(w, r, err, errorMessage)

		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	const responseMessage = "Email has been sent"

	user, err := h.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			response.SuccessResponseOK(w, responseMessage, nil)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	if user.IsActivated() || user.IsDeactivated() {
		response.SuccessResponseOK(w, responseMessage, nil)
		return
	}

	token, err := utils.GenerateToken(user.Email, []byte(cfg.SecretKey))
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}
	// TODO: Send email to activate account.
	fmt.Printf("'%s'\n", token)

	response.SuccessResponseOK(w, responseMessage, nil)
}
