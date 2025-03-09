package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

// TrimString is a custom string type that trims leading and trailing whitespace
// from JSON string inputs during deserialization.
//
// This type is useful for automatically cleaning up user inputs or any string fields
// that may accidentally include unnecessary spaces.
//
// Example:
//
//	jsonData := `{"name": "  John Doe  "}`
//	After unmarshaling, the name will be "John Doe" without extra spaces.
type TrimString string

// UnmarshalJSON implements the json.Unmarshaler interface for TrimString.
// It trims any leading and trailing whitespace from the input string during deserialization.
//
// Parameters:
//   - data ([]byte): The raw JSON data to be unmarshaled.
//
// Returns:
//   - error: An error if the input data cannot be unmarshaled as a string, otherwise nil.
//
// Example usage:
//
//	var name TrimString
//	err := json.Unmarshal([]byte(`"  Example Name  "`), &name)
//	fmt.Println(name) // Output: "Example Name"
func (t *TrimString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = TrimString(strings.TrimSpace(s))
	return nil
}

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
