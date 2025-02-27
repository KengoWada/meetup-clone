package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KengoWada/meetup-clone/internal/models"
)

var (
	ErrDuplicateEmail    = errors.New("an account is already attached to that email address")
	ErrDuplicateUsername = errors.New("username is already taken")
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *models.User, userProfile *models.UserProfile) error {
	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.createUser(ctx, tx, user); err != nil {
			return err
		}

		userProfile.UserID = user.ID
		if err := s.createUserProfile(ctx, tx, userProfile); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET is_active = 't', version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING version, is_active, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.ID, user.Version).Scan(&user.Version, &user.IsActive, &user.UpdatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) createUser(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `
		INSERT INTO users(email, password, role)
		VALUES($1, $2, $3)
		RETURNING id, is_active, version, created_at, updated_at, deleted_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(
		&user.ID,
		&user.IsActive,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) createUserProfile(ctx context.Context, tx *sql.Tx, userProfile *models.UserProfile) error {
	query := `
		INSERT INTO user_profiles(username, profile_pic, date_of_birth, user_id)
		VALUES($1, $2, $3, $4)
		RETURNING id, version, created_at, updated_at, deleted_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		userProfile.Username,
		userProfile.ProfilePic,
		userProfile.DateOfBirth,
		userProfile.UserID,
	).Scan(
		&userProfile.ID,
		&userProfile.Version,
		&userProfile.CreatedAt,
		&userProfile.UpdatedAt,
		&userProfile.DeletedAt,
	)

	if err != nil {
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "user_profiles_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}
