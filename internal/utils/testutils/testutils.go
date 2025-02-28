package testutils

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

func ExecuteRequest(r *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	return w
}
