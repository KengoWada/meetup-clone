package auth

import (
	"net/http"
	"strings"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
)

type registerUserPayload struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=10,max=72,is_password"`
	Username    string `json:"username" validate:"required,min=3,max=100"`
	ProfilePic  string `json:"profilePic" validate:"required,http_url" example:"https://fake.link/img.png"`
	DateOfBirth string `json:"dateOfBirth" validate:"required,is_date" example:"mm/dd/yyyy"`
}

// RegisterUser godoc
//
//	@Summary		Register a user
//	@Description	Register a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		registerUserPayload	true	"register user payload"
//	@Success		201		{object}	response.DocsSuccessResponseRegisterUser
//	@Failure		400		{object}	response.DocsErrorResponse
//	@Failure		500		{object}	response.DocsErrorResponseInternalServerErr
//	@Security
//	@Router	/auth/register [post]
func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "json: unknown field") {
			response.ErrorResponseUnknownField(w, r, err)
			return
		}
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}

	if errResponse, err := utils.ValidatePayload(payload, registerUserPayloadErrors); err != nil {
		switch err {
		case utils.ErrFailedValidation:
			errorMessage := response.NewValidationErrorResponse(errResponse)
			response.ErrorResponseBadRequest(w, r, err, errorMessage)

		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	passwordHash, err := utils.GeneratePasswordHash(payload.Password)
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
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
		errorMessage := response.ErrorResponse{Message: response.ValidationErrorMessage}

		switch err {
		case store.ErrDuplicateEmail:
			errorMessage.Errors = response.ErrorsResponse{"email": err.Error()}
			response.ErrorResponseBadRequest(w, r, err, errorMessage)

		case store.ErrDuplicateUsername:
			errorMessage.Errors = response.ErrorsResponse{"username": err.Error()}
			response.ErrorResponseBadRequest(w, r, err, errorMessage)

		default:
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	_, err = utils.GenerateToken(user.Email, []byte(h.config.SecretKey))
	if err != nil {
		response.ErrorResponseInternalServerErr(w, r, err)
		return
	}
	// TODO: Send email to activate account.

	response.SuccessResponseCreated(w, "Done.", nil)
}
