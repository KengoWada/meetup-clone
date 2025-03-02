// Package store contains the logic for interacting with the database,
// including SQL queries and operations that map data between the models
// and the database. It provides the functionality to retrieve, insert,
// update, and delete data from the database.
package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/KengoWada/meetup-clone/internal/models"
)

const QueryTimeoutDuration = time.Second * 5

var ErrNotFound = errors.New("item not found")

// Store defines the interfaces and methods for interacting with various
// data models in the database. It includes methods for creating, reading,
// updating, and deleting data.
type Store struct {
	Users interface {
		Create(context.Context, *models.User, *models.UserProfile) error
		Activate(context.Context, *models.User) error
		Deactivate(context.Context, *models.User) error
		GetByEmail(ctx context.Context, email string) (*models.User, error)
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Users: &UserStore{db},
	}
}

// WithTx runs a set of queries within a database transaction. It begins a
// transaction, executes the provided function with the transaction, and
// commits or rolls back the transaction depending on the result. The function
// ensures that all queries within the transaction are executed atomically.
//
// The function fn is passed a pointer to an active transaction, and should
// return an error if any issues arise during query execution. If fn returns
// no error, the transaction is committed; otherwise, it is rolled back.
func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
