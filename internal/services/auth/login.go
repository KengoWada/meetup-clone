package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
	"github.com/golang-jwt/jwt/v5"
)

var (
	errInvalidPassword         = errors.New("invalid password provided")
	errEmailNotVerified        = errors.New("email not verified")
	errDeactivatedAccountLogin = errors.New("log in attempt on a deactivated account")
)

type loginUserPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginUser godoc
//
//	@Summary		Log in a user
//	@Description	Log in a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		loginUserPayload						true	"log in payload"
//	@Success		200		{object}	response.DocsSuccessResponseLoginUser	"user successfully logged in"
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security
//	@Router	/auth/login [post]
func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var payload loginUserPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		response.ErrorResponseInvalidJSON(w, r, err)
		return
	}

	if errResponse, err := validate.ValidatePayload(payload, validate.FieldErrorMessages{}); err != nil {
		switch err {
		case validate.ErrFailedValidation:
			errorMessage := response.NewValidationErrorResponse(errResponse)
			response.ErrorResponseBadRequest(w, r, err, errorMessage)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	errorMessage := response.ErrorResponse{Message: "Invalid credentials"}

	fields, values := []string{"email"}, []any{payload.Email}
	user, err := h.store.Users.Get(r.Context(), false, fields, values)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			response.ErrorResponseBadRequest(w, r, err, errorMessage)
		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	if user.IsDeactivated() {
		response.ErrorResponseBadRequest(w, r, errDeactivatedAccountLogin, errorMessage)
		return
	}

	if !user.IsActive {
		errorMessage := response.ErrorResponse{Message: "Please verify your email address to proceed."}
		response.ErrorResponseUnprocessableEntity(w, r, errEmailNotVerified, errorMessage)
		return
	}

	ok, err := utils.ComparePasswordAndHash(payload.Password, user.Password)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if !ok {
		response.ErrorResponseBadRequest(w, r, errInvalidPassword, errorMessage)
		return
	}

	exp := time.Hour * time.Duration(cfg.AuthConfig.Exp)
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": cfg.AuthConfig.Issuer,
		"aud": cfg.AuthConfig.Audience,
	}

	token, err := h.authenticator.GenerateToken(claims)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	data := response.Response{"token": token}
	response.SuccessResponseOK(w, "", data)
}
