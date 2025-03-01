// Package response provides various structs and utility functions for generating
// consistent API responses, including success and error responses. These types
// are used to standardize how responses are sent back to the client, making the
// API easier to integrate with and consume. The package supports different
// response types such as success messages, validation errors, internal errors,
// and examples for Swagger documentation.
package response

import "github.com/KengoWada/meetup-clone/internal/logger"

// ValidationErrorMessage is a constant that represents the default error message
// used when the request body is invalid or fails validation.
const ValidationErrorMessage = "Invalid request body"

var log = logger.Get()

// Response represents a generic response structure that can hold any type of data.
// It is a map where the key is a string and the value can be of any type.
type Response map[string]any

// ErrorsResponse is a map that represents validation errors.
// The keys are the field names, and the values are the error messages.
type ErrorsResponse map[string]string

// SuccessResponse represents a successful API response.
// It contains a message and optional data that can be returned as part of the success.
type SuccessResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// ErrorResponse represents an error response in the API.
// It includes a message and optionally a set of validation errors.
type ErrorResponse struct {
	Message string         `json:"message"`          // The error message.
	Errors  ErrorsResponse `json:"errors,omitempty"` // Optional map of validation errors.
}

// NewValidationErrorResponse creates a new ErrorResponse with a validation error message
// and a map of validation errors.
func NewValidationErrorResponse(errorsResponse ErrorsResponse) ErrorResponse {
	return ErrorResponse{Message: ValidationErrorMessage, Errors: errorsResponse}
}
