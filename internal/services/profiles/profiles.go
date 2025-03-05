package profiles

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/validate"
)

type userProfile struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	ProfilePic  string `json:"profilePic"`
	Role        string `json:"role" example:"client"`
	DateOfBirth string `json:"dateOfBirth" example:"mm/dd/yyyy"`
}

type updateUserDetailsPayload struct {
	Email       string `json:"email" validate:"required,email"`
	Username    string `json:"username" validate:"required,min=3,max=100"`
	ProfilePic  string `json:"profilePic" validate:"required,http_url"`
	DateOfBirth string `json:"dateOfBirth" validate:"required,is_date"`
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

// UpdateUserProfiles godoc
//
//	@Summary		Update a users profile details
//	@Description	Update a users profile details
//	@Tags			profiles
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	userProfile
//	@Failure		400	{object}	response.DocsResponseMessageOnly
//	@Failure		401	{object}	response.DocsErrorResponseUnauthorized
//	@Failure		500	{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/profiles/ [put]
func (h *Handler) updateUserProfile(w http.ResponseWriter, r *http.Request) {
	var payload updateUserDetailsPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "json: unknown field") {
			response.ErrorResponseUnknownField(w, r, err)
			return
		}
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if errResponse, err := validate.ValidatePayload(payload, updateUserPayloadErrors); err != nil {
		switch err {
		case validate.ErrFailedValidation:
			errorMessage := response.NewValidationErrorResponse(errResponse)
			response.ErrorResponseBadRequest(w, r, err, errorMessage)

		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	user, _ := r.Context().Value(internal.UserCtx).(*models.User)

	user.Email = payload.Email
	user.UserProfile.Username = payload.Username
	user.UserProfile.ProfilePic = payload.ProfilePic
	user.UserProfile.DateOfBirth = payload.DateOfBirth

	err = h.store.Users.UpdateUserDetails(r.Context(), user)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			// The reasoning behind this is that the version may be
			// off because an admin is running an update on the user as well
			// causing the query to return the no rows error.
			res := response.ErrorResponse{Message: "Try again later"}
			response.ErrorResponseBadRequest(w, r, err, res)
		default:
			fmt.Println(err)
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

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
