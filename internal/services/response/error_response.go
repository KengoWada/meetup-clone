package response

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/logger"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

var log = logger.Get()

type H map[string]any

type Errors map[string]string

type ErrorResponse struct {
	Message string `json:"message"`
	Errors  Errors `json:"errors,omitempty"`
}

func InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
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

func BadRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error, response ErrorResponse) {
	reqIDRaw := middleware.GetReqID(r.Context())
	log.Warn().
		Str("requestID", reqIDRaw).
		Str("method", r.Method).
		Str("url", r.URL.Path).
		Err(errors.Wrap(err, "bad request")).
		Msg("Bad Request")

	utils.WriteJSON(w, http.StatusBadRequest, response)
}
