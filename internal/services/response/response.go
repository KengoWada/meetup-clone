package response

import "github.com/KengoWada/meetup-clone/internal/logger"

const ValidationErrorMessage = "Invalid request body"

var log = logger.Get()

type Response map[string]any

type ErrorsResponse map[string]string

type SuccessResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string         `json:"message"`
	Errors  ErrorsResponse `json:"errors,omitempty"`
}

func NewValidationErrorResponse(errorsResponse ErrorsResponse) ErrorResponse {
	return ErrorResponse{Message: ValidationErrorMessage, Errors: errorsResponse}
}
