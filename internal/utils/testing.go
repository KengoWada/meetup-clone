package utils

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/KengoWada/meetup-clone/internal/store"
	_ "github.com/lib/pq"
)

func NewTestStore(t *testing.T) store.Store {
	t.Helper()

	dbAddr := GetString("TEST_DB_ADDR", "")
	db := sql.OpenDB(txdb.New("postgres", dbAddr))
	if db == nil {
		t.Fatal("failed to connect to test db")
	}

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
