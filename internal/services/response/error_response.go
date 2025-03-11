package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

// ErrorResponseInternalServerErr returns an internal server error response (HTTP 500)
// and logs the request details along with the provided error for error tracking.
// If the provided error is nil, no logging is performed, as this indicates a
// unique situation where the logger itself has failed. In all other cases,
// the error and request context are logged for debugging and monitoring purposes.
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

// ErrorResponseBadRequest returns a bad request error response (HTTP 400)
// and includes a detailed error message. It also logs the request context
// and error for debugging purposes. The provided response is included in
// the response body, which may contain additional information or validation errors.
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

// ErrorResponseUnprocessableEntity returns an unprocessable entity error response (HTTP 422)
// and includes a detailed error message. The provided response is sent to the user,
// which may contain additional information or validation errors. This function
// logs the request context and error for debugging purposes.
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

// ErrorResponseUnknownField returns a bad request error response (HTTP 400)
// when the user sends unknown fields in the request. It includes a message
// indicating the presence of unknown fields and logs the request context
// and error for debugging purposes.
func ErrorResponseUnknownField(w http.ResponseWriter, r *http.Request, err error) {
	items := strings.Split(err.Error(), " ")
	fieldName := strings.ReplaceAll(items[len(items)-1], `"`, "")
	errorResponse := ErrorResponse{
		Message: "Unknown field in request",
		Errors:  ErrorsResponse{fieldName: "unknown field"},
	}

	ErrorResponseBadRequest(w, r, err, errorResponse)
}

// ErrorResponseUnauthorized returns an unauthorized request error response (HTTP 401)
// when the user sends invalid credentials.
func ErrorResponseUnauthorized(w http.ResponseWriter, r *http.Request, err error) {
	reqIDRaw := middleware.GetReqID(r.Context())
	log.Warn().
		Str("requestID", reqIDRaw).
		Str("method", r.Method).
		Str("url", r.URL.Path).
		Err(errors.Wrap(err, "unauthorized")).
		Msg("Unauthorized")

	response := ErrorResponse{Message: "unauthorized"}
	utils.WriteJSON(w, http.StatusUnauthorized, response)
}

func ErrorResponseForbidden(w http.ResponseWriter, r *http.Request, err error) {
	reqIDRaw := middleware.GetReqID(r.Context())
	log.Warn().
		Str("requestID", reqIDRaw).
		Str("method", r.Method).
		Str("url", r.URL.Path).
		Err(errors.Wrap(err, "forbidden")).
		Msg("Forbidden")

	response := ErrorResponse{Message: "forbidden"}
	utils.WriteJSON(w, http.StatusForbidden, response)
}

func ErrorResponseInvalidJSON(w http.ResponseWriter, r *http.Request, err error) {
	var (
		status             int            = http.StatusBadRequest
		errorsResponse     ErrorsResponse = ErrorsResponse{}
		errorMessage       string
		syntaxError        *json.SyntaxError
		unmarshalTypeError *json.UnmarshalTypeError
		maxBytesError      *http.MaxBytesError
	)
	switch {
	case errors.As(err, &syntaxError):
		errorMessage = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
		errorMessage = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

	case errors.As(err, &unmarshalTypeError):
		errorMessage = fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		ErrorResponseUnknownField(w, r, err)
		return

	case errors.Is(err, io.EOF):
		errorMessage = "Request body must not be empty"

	case errors.As(err, &maxBytesError):
		errorMessage = fmt.Sprintf("Request body must not be larger than %d bytes", maxBytesError.Limit)

	default:
		status = http.StatusInternalServerError
	}

	if status == http.StatusInternalServerError {
		ErrorResponseInternalServerErr(w, r, err)
		return
	}

	errorResponse := ErrorResponse{Message: errorMessage, Errors: errorsResponse}
	ErrorResponseBadRequest(w, r, err, errorResponse)
}
