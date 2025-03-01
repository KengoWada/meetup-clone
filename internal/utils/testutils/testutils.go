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
	"database/sql"
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

// ExecuteTestRequest simulates an HTTP request to a specified handler and returns the response recorder.
// This is useful for testing HTTP endpoints in isolation without needing to spin up a full HTTP server.
//
// Parameters:
//   - r (*http.Request): The HTTP request to be executed, typically created using httptest.NewRequest().
//   - mux (http.Handler): The HTTP handler to which the request is routed. This is typically the router or multiplexer (mux) that handles routing of the request.
//
// Returns:
//   - *httptest.ResponseRecorder: A response recorder that captures the HTTP response generated by the handler.
//     This allows access to the response status code, headers, and body content for assertions in tests.
//
// Example usage:
//
//	func TestMyHandler(t *testing.T) {
//	    req := httptest.NewRequest("GET", "/my-endpoint", nil)
//	    mux := http.NewServeMux()
//	    mux.HandleFunc("/my-endpoint", myHandler) // Register the handler for the endpoint
//
//	    respRecorder := ExecuteTestRequest(req, mux)
//
//	    if respRecorder.Code != http.StatusOK {
//	        t.Errorf("Expected status OK, got %v", respRecorder.Code)
//	    }
//	}
func ExecuteTestRequest(r *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	return w
}
