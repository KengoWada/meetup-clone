package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
)

type resetUserPasswordPayload struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=10,max=72,is_password"`
}

type passwordResetRequestPayload struct {
	Email string `json:"email" validate:"required"`
}

// ResetUserPassword godoc
//
//	@Summary		Reset a users password
//	@Description	Reset a users password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		resetUserPasswordPayload	true	"reset user password payload"
//	@Success		200		{object}	response.DocsResponseMessageOnly
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		422		{object}	response.DocsResponseMessageOnly
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security
//	@Router	/auth/reset-password [post]
func (h *Handler) resetUserPassword(w http.ResponseWriter, r *http.Request) {
	var payload resetUserPasswordPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "json: unknown field") {
			response.ErrorResponseUnknownField(w, r, err)
			return
		}
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if errResponse, err := utils.ValidatePayload(payload, resetUserPasswordPayloadErrors); err != nil {
		switch err {
		case utils.ErrFailedValidation:
			errorMessage := response.NewValidationErrorResponse(errResponse)
			response.ErrorResponseBadRequest(w, r, err, errorMessage)

		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	timedToken, err := utils.ValidateToken(payload.Token, []byte(h.config.SecretKey), time.Minute*30)
	if err != nil {
		switch err {
		case utils.ErrExpiredToken:
			errorMessage := response.ErrorResponse{Message: "Password reset token has exipred"}
			response.ErrorResponseUnprocessableEntity(w, r, err, errorMessage)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	ctx := r.Context()
	errorMessage := response.ErrorResponse{Message: "Password reset token is invalid"}

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

	if user.IsDeactivated() || !user.IsActive || payload.Token != user.PasswordResetToken {
		response.ErrorResponseBadRequest(w, r, err, errorMessage)
		return
	}

	passwordHash, err := utils.GeneratePasswordHash(payload.Password)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	user.Password = passwordHash
	user.PasswordResetToken = ""
	err = h.store.Users.ResetPassword(ctx, user)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	response.SuccessResponseOK(w, "Password successfully updated", nil)
}

// PasswordResetRequest godoc
//
//	@Summary		Request to reset a users password
//	@Description	Request to reset a users password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		passwordResetRequestPayload	true	"password reset request payload"
//	@Success		200		{object}	response.DocsResponseMessageOnly
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security
//	@Router	/auth/password-reset-request [post]
func (h *Handler) passwordResetRequest(w http.ResponseWriter, r *http.Request) {
	var payload passwordResetRequestPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "json: unknown field") {
			response.ErrorResponseUnknownField(w, r, err)
			return
		}
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if errResponse, err := utils.ValidatePayload(payload, passwordResetRequestPayloadErrors); err != nil {
		switch err {
		case utils.ErrFailedValidation:
			errorMessage := response.NewValidationErrorResponse(errResponse)
			response.ErrorResponseBadRequest(w, r, err, errorMessage)

		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	const message = "Email has been sent."
	var ctx = r.Context()

	user, err := h.store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			response.SuccessResponseOK(w, message, nil)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	if user.IsDeactivated() || !user.IsActive {
		response.SuccessResponseOK(w, message, nil)
		return
	}

	token, err := utils.GenerateToken(user.Email, []byte(h.config.SecretKey))
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	user.PasswordResetToken = token
	err = h.store.Users.SetPasswordResetToken(ctx, user)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	// TODO: send email to user

	response.SuccessResponseOK(w, message, nil)
}
