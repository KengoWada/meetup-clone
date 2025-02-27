package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/KengoWada/meetup-clone/internal/models"
)

const (
	QueryTimeoutDuration = time.Second * 5
)

var ErrNotFound = errors.New("item not found")

type Store struct {
	Users interface {
		Create(context.Context, *models.User, *models.UserProfile) error
		Activate(context.Context, *models.User) error
		GetByEmail(ctx context.Context, email string) (*models.User, error)
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Users: &UserStore{db},
	}
}

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
