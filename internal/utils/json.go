package utils

import (
	"encoding/json"
	"net/http"
)

// ReadJSON reads the JSON request body from the HTTP request and unmarshals it into
// the provided data object. It ensures that no unknown fields are included in the
// request body, returning an error if such fields are present.
//
// Parameters:
//   - w: The HTTP response writer, which can be used to send a response if an error occurs.
//   - r: The HTTP request containing the JSON body to be read.
//   - data: The object where the JSON body will be unmarshaled.
//
// Returns:
//   - An error if the request body contains unknown fields or any other unmarshalling issue.
//
// Errors:
//   - Returns an error if unknown fields are encountered in the JSON body or if unmarshalling fails.
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578
	http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

// WriteJSON writes the given data as a JSON response with the specified HTTP status code.
// It serializes the provided data and sends it back to the client. The response content
// type is set to application/json.
//
// Parameters:
//   - w: The HTTP response writer used to send the response.
//   - status: The HTTP status code to be returned with the response.
//   - data: The data to be serialized to JSON and included in the response body.
//
// Returns:
//   - An error if there is a problem with serializing the data to JSON or writing the response.
//
// Errors:
//   - Returns an error if JSON serialization fails or if writing to the response fails.
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
