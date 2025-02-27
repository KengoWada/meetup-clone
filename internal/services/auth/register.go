package auth

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
)

type registerUserPayload struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=10,max=72,is_password"`
	Username    string `json:"username" validate:"required,min=3,max=100"`
	ProfilePic  string `json:"profilePic" validate:"required,http_url"`
	DateOfBirth string `json:"dateOfBirth" validate:"required,is_date"`
}

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}

	if errResponse, err := utils.ValidatePayload(payload, registerUserPayloadErrors); err != nil {
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

	passwordHash, err := utils.GeneratePasswordHash(payload.Password)
	if err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}

	user := &models.User{
		Email:    payload.Email,
		Password: passwordHash,
		Role:     models.UserClientRole,
	}
	userProfile := &models.UserProfile{
		Username:    payload.Username,
		ProfilePic:  payload.ProfilePic,
		DateOfBirth: payload.DateOfBirth,
	}

	ctx := r.Context()

	err = h.store.Users.Create(ctx, user, userProfile)
	if err != nil {
		errorMessage := response.ErrorResponse{Message: "Invalid request body"}

		switch err {
		case store.ErrDuplicateEmail:
			errorMessage.Errors = response.Errors{"email": err.Error()}
			response.BadRequestErrorResponse(w, r, err, errorMessage)

		case store.ErrDuplicateUsername:
			errorMessage.Errors = response.Errors{"username": err.Error()}
			response.BadRequestErrorResponse(w, r, err, errorMessage)

		default:
			response.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	// TODO: Send email to activate account.

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Done."})
}
