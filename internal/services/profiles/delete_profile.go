package profiles

import (
	"fmt"
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
)

// DeleteUserAccount godoc
//
//	@Summary		Delete a users account details(soft delete)
//	@Description	Delete a users profile details(soft delete)
//	@Tags			profiles
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.DocsResponseMessageOnly
//	@Failure		400	{object}	response.DocsResponseMessageOnly
//	@Failure		401	{object}	response.DocsErrorResponseUnauthorized
//	@Failure		500	{object}	response.DocsErrorResponseInternalServerErr
//	@Security		ApiKeyAuth
//	@Router			/profiles/ [delete]
func (h *Handler) deleteUserProfile(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(internal.UserCtx).(*models.User)

	err := h.store.Users.SoftDeleteUser(r.Context(), user)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			res := response.ErrorResponse{Message: "Try again later"}
			response.ErrorResponseBadRequest(w, r, err, res)
		default:
			fmt.Println(err)
			response.ErrorResponseInternalServerErr(w, r, err)
		}
		return
	}

	response.SuccessResponseOK(w, "Done", nil)
}
