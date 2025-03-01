package response

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/utils"
)

func SuccessResponseCreated(w http.ResponseWriter, message string, data any) {
	response := SuccessResponse{Message: message, Data: data}
	utils.WriteJSON(w, http.StatusCreated, response)
}

func SuccessResponseOK(w http.ResponseWriter, message string, data any) {
	response := SuccessResponse{Message: message, Data: data}
	utils.WriteJSON(w, http.StatusOK, response)
}
