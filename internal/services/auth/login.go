package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

var errInvalidPassword = errors.New("invalid password provided")

type loginUserPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var payload loginUserPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "json: unknown field") {
			response.UnknownFieldErrorResponse(w, r, err)
			return
		}
		response.InternalServerErrorResponse(w, r, err)
		return
	}

	if errResponse, err := utils.ValidatePayload(payload, utils.FieldErrorMessages{}); err != nil {
		switch err {
		case utils.ErrFailedValidation:
			errorMessage := response.ErrorResponse{
				Message: "Invalid request body",
				Errors:  errResponse,
			}
			response.BadRequestErrorResponse(w, r, err, errorMessage)
		default:
			response.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	errorMessage := response.ErrorResponse{Message: "Invalid credentials"}

	user, err := h.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			response.BadRequestErrorResponse(w, r, err, errorMessage)
		default:
			response.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	if user.IsDeactivated() {
		err := errors.New("log in attempt on deactivated account")
		response.BadRequestErrorResponse(w, r, err, errorMessage)
		return
	}

	if !user.IsActive {
		err := errors.New("email not verified")
		errorMessage := response.ErrorResponse{Message: "Please verify your email address to proceed."}
		response.UnprocessableEntityErrorResponse(w, r, err, errorMessage)
		return
	}

	ok, err := utils.ComparePasswordAndHash(payload.Password, user.Password)
	if err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}

	if !ok {
		response.BadRequestErrorResponse(w, r, errInvalidPassword, errorMessage)
		return
	}

	exp := time.Hour * time.Duration(h.config.AuthConfig.Exp)
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.config.AuthConfig.Issuer,
		"aud": h.config.AuthConfig.Audience,
	}

	token, err := h.authenticator.GenerateToken(claims)
	if err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}

	resp := response.SuccessResponse{Data: map[string]string{"token": token}}
	utils.WriteJSON(w, http.StatusOK, resp)
}
