package response

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/utils"
)

// SuccessResponseCreated returns a success response with a status of HTTP 201 (Created).
// It includes a message indicating the successful creation and optional data
// associated with the created resource. The function sends the response to
// the client with the provided message and data.
func SuccessResponseCreated(w http.ResponseWriter, message string, data any) {
	response := SuccessResponse{Message: message, Data: data}
	utils.WriteJSON(w, http.StatusCreated, response)
}

// SuccessResponseOK returns a success response with a status of HTTP 200 (OK).
// It includes a message indicating the successful operation and optional data
// associated with the operation's result. The function sends the response to
// the client with the provided message and data.
func SuccessResponseOK(w http.ResponseWriter, message string, data any) {
	response := SuccessResponse{Message: message, Data: data}
	utils.WriteJSON(w, http.StatusOK, response)
}
