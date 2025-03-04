package profiles

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
)

type userProfile struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	ProfilePic  string `json:"profilePic"`
	Role        string `json:"role" example:"client"`
	DateOfBirth string `json:"dateOfBirth" example:"mm/dd/yyyy"`
}

// GetPersonalProfile godoc
//
//	@Summary		Get a users details based on token provided
//	@Description	Get a users details based on token provided
//	@Tags			profiles
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	userProfile
//	@Failure		401	{object}	response.DocsErrorResponseUnauthorized
//	@Failure		500	{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/profiles/ [get]
func (h *Handler) getPersonalProfile(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(internal.UserCtx).(*models.User)

	userDetails := userProfile{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.UserProfile.Username,
		ProfilePic:  user.UserProfile.ProfilePic,
		Role:        string(user.Role),
		DateOfBirth: user.UserProfile.DateOfBirth,
	}

	response.SuccessResponseOK(w, "", userDetails)
}
