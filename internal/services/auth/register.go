package auth

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
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
	utils.ReadJSON(w, r, &payload)

	if err := validate.Validate.Struct(payload); err != nil {
		errResponse, err := utils.GenerateErrorMessages(err, registerUserPayloadErrors)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error."})
			return
		}
		utils.WriteJSON(w, http.StatusBadRequest, errResponse)
		return
	}

	passwordHash, err := utils.GeneratePasswordHash(payload.Password)
	if err != nil {
		// TODO: Add logging
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error."})
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

	err = h.store.Users.Create(r.Context(), user, userProfile)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		case store.ErrDuplicateUsername:
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error."})
		}
		return
	}

	// TODO: Send email to activate account.

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Done."})
}
