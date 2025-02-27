package utils

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KengoWada/meetup-clone/internal/store"
)

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

func ExecuteRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}
