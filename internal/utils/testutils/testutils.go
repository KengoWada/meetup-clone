// Package testutils contains helper functions and utilities to assist with running tests.
// These functions are designed to simplify test setup, teardown, and common tasks
// required during unit and integration testing in Go.
//
// It includes utilities for mocking, creating test data, and other commonly
// needed operations that can help make test code more readable, reusable, and maintainable.
// Example usage:
//
//	testutils.NewTestStore(t, db) // Set up a test store that rolls back all changes after a test run.
package testutils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KengoWada/meetup-clone/internal/store"
)

// NewTestStore sets up a test store for running tests that interacts with a real database.
// It ensures that all changes made during the test are rolled back after the test completes.
// This is achieved by using a database transaction, so any database modifications (inserts, updates, deletes)
// performed during the test will not persist beyond the test run.
//
// Parameters:
//   - t (testing.T): The test instance provided by the testing framework, used for managing test state.
//   - db (*sql.DB): The database connection to be used by the test store.
//
// Returns:
//   - store.Store: A store instance configured for testing. It uses transactions that will be rolled back
//     after the test completes.
//
// Example usage:
//
//	func TestCreateUser(t *testing.T) {
//	    testStore := NewTestStore(t, db)
//	    // Use testStore to interact with the database during the test
//	    // At the end of the test, all changes to the database will be rolled back
func NewTestStore(t *testing.T, db *sql.DB) store.Store {
	t.Helper()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	t.Cleanup(func() {
		tx.Rollback()
		db.Close()
	})

	return store.NewStore(db)
}

// TestRequestData represents the request payload for a test HTTP request.
// It uses a flexible map structure to allow arbitrary key-value pairs for building request bodies.
type TestRequestData map[string]any

// TestResponseData represents the response body of a test HTTP request.
// Like the request, it uses a map to handle dynamic response structures.
type TestResponseData map[string]any

// TestRequestHeaders represents a collection of HTTP headers for test requests.
// It is a map where the key is the header name and the value is the header value.
//
// Example:
//
//	headers := TestRequestHeaders{
//	  "Authorization": "Bearer some.jwt.token",
//	  "Content-Type":  "application/json",
//	}
type TestRequestHeaders map[string]string

// TestRequestResponse wraps the result of an HTTP test request.
// It includes the raw response recorder and the parsed response data.
//
// Fields:
//   - Response (*httptest.ResponseRecorder): Captures the HTTP response for status, headers, and body.
//   - ResponseData (TestResponseData): The decoded JSON response body.
type TestRequestResponse struct {
	Response     *httptest.ResponseRecorder
	ResponseData TestResponseData
}

// StatusCode returns the HTTP status code from the response.
//
// This method provides a convenient way to access the status code of the
// test HTTP response. If the response is nil, it returns nil to avoid panics.
//
// Returns:
//   - *int: A pointer to the status code, or nil if the response is not set.
//
// Example usage:
//
//	res, _ := RunTestRequest(router, "GET", "/users", nil, nil)
//	if statusCode := res.StatusCode(); statusCode != nil {
//	  fmt.Println(*statusCode)  // prints the HTTP status code
//	}
func (r *TestRequestResponse) StatusCode() int {
	if r.Response != nil {
		return r.Response.Code
	}
	return -1
}

// GetMessage extracts the "message" field from the response data.
//
// This method checks if the "message" field exists and is a string.
// If the field is missing or not a string, it returns an empty string.
//
// Returns:
//   - string: The message string from the response, or an empty string if not found or invalid.
//
// Example usage:
//
//	res, _ := RunTestRequest(router, "POST", "/login", nil, requestData)
//	fmt.Println(res.GetMessage()) // prints the message from the response, or an empty string
func (r *TestRequestResponse) GetMessage() string {
	message, ok := r.ResponseData["message"]
	if !ok {
		return ""
	}

	messageStr, ok := message.(string)
	if !ok {
		return ""
	}

	return messageStr
}

// GetErrorMessages extracts the "errors" field from the response data.
//
// This method returns the error messages as a map, along with a boolean indicating success.
// If the "errors" field isn't present or isn't a map, it returns false.
//
// Returns:
//   - map[string]any: The error messages from the response.
//   - bool: True if the "errors" field was found and is a map, false otherwise.
//
// Example usage:
//
//	res, _ := RunTestRequest(router, "POST", "/register", nil, requestData)
//	if errors, ok := res.GetErrorMessages(); ok {
//	  fmt.Println(errors) // prints the error messages
//	}
func (r *TestRequestResponse) GetErrorMessages() (map[string]any, bool) {
	errorMessages, ok := r.ResponseData["errors"].(map[string]any)
	return errorMessages, ok
}

// GetData extracts the "data" field from the response data.
//
// This method returns the response data as a map, along with a boolean indicating success.
// If the "data" field isn't present or isn't a map, it returns false.
//
// Returns:
//   - map[string]any: The data from the response.
//   - bool: True if the "data" field was found and is a map, false otherwise.
//
// Example usage:
//
//	res, _ := RunTestRequest(router, "GET", "/profile", nil, nil)
//	if data, ok := res.GetData(); ok {
//	  fmt.Println(data) // prints the response data
//	}
func (r *TestRequestResponse) GetData() (map[string]any, bool) {
	errorMessages, ok := r.ResponseData["data"].(map[string]any)
	return errorMessages, ok
}

// RunTestRequest performs a test HTTP request and returns the result.
//
// It marshals the provided request data to JSON, sets the appropriate content type,
// and executes the request against the given HTTP handler (mux). The function then
// captures the response and unmarshals the JSON response body for easier assertions.
//
// Parameters:
//   - mux (http.Handler): The HTTP handler to serve the request (e.g., your router).
//   - method (string): The HTTP method to use (e.g., "GET", "POST", "PUT", "DELETE").
//   - endpoint (string): The API endpoint to test (e.g., "/api/login").
//   - headers (TestRequestHeaders): The headers to include in the request (e.g., authorization, content type).
//   - data (TestRequestData): The request payload, which will be marshaled to JSON.
//
// Returns:
//   - *TestRequestResponse: Contains the HTTP response recorder and the parsed response data.
//   - error: An error if request creation, execution, or response decoding fails.
//
// Example usage:
//
//			res, err := RunTestRequest(
//	         router, "/login", "POST",
//	         TestRequestHeaders{"Authorization": "Bearer some.jwt.token"},
//		     TestRequestData{
//			   "email": "test@example.com",
//			   "password": "password123",
//		     })
//			if err != nil {
//			  log.Fatal(err)
//			}
//			fmt.Println(res.Response.Code)               // prints the HTTP status code
//			fmt.Println(res.ResponseData["message"])     // prints the response message
func RunTestRequest(mux http.Handler, method, endpoint string, headers TestRequestHeaders, data TestRequestData) (*TestRequestResponse, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		r.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	response := make(TestResponseData)
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}

	return &TestRequestResponse{Response: w, ResponseData: response}, nil
}
