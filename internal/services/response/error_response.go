package response

import (
	"net/http"
	"strings"

	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

func ErrorResponseInternalServerErr(w http.ResponseWriter, r *http.Request, err error) {
	reqID := middleware.GetReqID(r.Context())
	log.Error().
		Str("requestID", reqID).
		Str("method", r.Method).
		Str("url", r.URL.Path).
		Err(errors.Wrap(err, "internal server error")).
		Msg("Internal Server Error")

	response := ErrorResponse{Message: "internal server error"}
	utils.WriteJSON(w, http.StatusInternalServerError, response)
}

func ErrorResponseBadRequest(w http.ResponseWriter, r *http.Request, err error, response ErrorResponse) {
	reqIDRaw := middleware.GetReqID(r.Context())
	log.Warn().
		Str("requestID", reqIDRaw).
		Str("method", r.Method).
		Str("url", r.URL.Path).
		Err(errors.Wrap(err, "bad request")).
		Msg("Bad Request")

	utils.WriteJSON(w, http.StatusBadRequest, response)
}

func ErrorResponseUnprocessableEntity(w http.ResponseWriter, r *http.Request, err error, response ErrorResponse) {
	reqIDRaw := middleware.GetReqID(r.Context())
	log.Info().
		Str("requestID", reqIDRaw).
		Str("method", r.Method).
		Str("url", r.URL.Path).
		Err(errors.Wrap(err, "unprocessable entity")).
		Msg("Unprocessable Entity")

	utils.WriteJSON(w, http.StatusUnprocessableEntity, response)
}

func ErrorResponseUnknownField(w http.ResponseWriter, r *http.Request, err error) {
	items := strings.Split(err.Error(), " ")
	fieldName := strings.ReplaceAll(items[len(items)-1], `"`, "")
	errorResponse := ErrorResponse{
		Message: "Unknown field in request",
		Errors:  ErrorsResponse{fieldName: "unknown field"},
	}

	ErrorResponseBadRequest(w, r, err, errorResponse)
}
