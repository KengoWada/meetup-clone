package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
)

var (
	ErrDuplicateEmail    = errors.New("an account is already attached to that email address")
	ErrDuplicateUsername = errors.New("username is already taken")
)

// UserStore provides methods for interacting with the database related to
// user operations, such as creating, updating, and retrieving users.
// It encapsulates the database connection and contains methods to perform
// CRUD operations on the user data.
type UserStore struct {
	db *sql.DB
}

// Create creates a new user and their associated profile in the database.
// It takes the user and userProfile as input, and inserts the corresponding
// records into the appropriate tables. If the operation is successful,
// it returns nil; otherwise, it returns an error.
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

// Activate activates a user account in the database. It updates the user's
// status to active and sets the activation timestamp. If the operation is
// successful, it returns nil; otherwise, it returns an error.
func (s *UserStore) Activate(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET is_active = 't', activated_at = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version, is_active, activated_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	timeNow := time.Now().UTC().Format(internal.DateTimeFormat)
	err := s.db.QueryRowContext(
		ctx,
		query,
		&timeNow,
		user.ID,
		user.Version,
	).Scan(
		&user.Version,
		&user.IsActive,
		&user.ActivatedAt,
		&user.UpdatedAt,
	)
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

// Deactivate deactivates a user by setting the IsActive field to false and
// optionally updating the ActivatedAt field to indicate when the deactivation occurred.
// This method assumes the user exists in the database.
//
// Parameters:
//   - ctx (context.Context): The context for managing cancellation and deadlines.
//   - user (*models.User): The user object to be deactivated.
//
// Returns:
//   - error: An error if the operation fails, or nil if the deactivation was successful.
func (s *UserStore) Deactivate(ctx context.Context, user *models.User) error {
	if user.IsDeactivated() {
		return nil
	}

	if user.IsActivated() {
		return s.deactivateActiveUser(ctx, user)
	}

	timeNow := time.Now().UTC().Format(internal.DateTimeFormat)
	user.ActivatedAt = &timeNow
	return s.deactivateInActiveUser(ctx, user)
}

// GetByEmail retrieves a user from the database by their email address.
func (s *UserStore) GetByID(ctx context.Context, ID int) (*models.User, error) {
	query := `
		SELECT u.*, up.* FROM users u
		INNER JOIN user_profiles up
		ON u.id = up.user_id
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := models.User{UserProfile: &models.UserProfile{}}
	err := s.db.
		QueryRowContext(ctx, query, ID).
		Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.IsActive,
			&user.ActivatedAt,
			&user.Role,
			&user.PasswordResetToken,
			&user.Version,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.UserProfile.ID,
			&user.UserProfile.Username,
			&user.UserProfile.ProfilePic,
			&user.UserProfile.DateOfBirth,
			&user.UserProfile.UserID,
			&user.UserProfile.Version,
			&user.UserProfile.CreatedAt,
			&user.UserProfile.UpdatedAt,
			&user.UserProfile.DeletedAt,
		)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// GetByEmail retrieves a user from the database by their email address.
// It returns the user object if a matching user is found, or nil and an error
// if no user is found or if an issue occurs during the query.
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT * FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := models.User{}
	err := s.db.
		QueryRowContext(ctx, query, email).
		Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.IsActive,
			&user.ActivatedAt,
			&user.Role,
			&user.PasswordResetToken,
			&user.Version,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// ResetPassword updates a user's password in the database.
//
// This function resets the password for the given user and persists the change to the database.
// It expects that the password has already been validated and hashed before being passed in.
//
// Parameters:
//   - ctx (context.Context): The context for managing request deadlines and cancellations.
//   - user (*models.User): A pointer to the user model containing the updated password.
//
// Returns:
//   - error: An error if updating the password in the database fails, or nil on success.
func (s *UserStore) ResetPassword(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET password = $1, password_reset_token = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.Password, user.PasswordResetToken, user.ID, user.Version).Scan(&user.Version, &user.UpdatedAt)
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

// SetPasswordResetToken sets a password reset token for a user.
//
// This function sets the secure token, associates it with the given user.
//
// Parameters:
//   - ctx (context.Context): The context for managing request deadlines and cancellations.
//   - user (*models.User): A pointer to the user model for which the reset token is generated.
//
// Returns:
//   - error: An error if generating the token or updating the user record fails, or nil on success.
func (s *UserStore) SetPasswordResetToken(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET password_reset_token = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.PasswordResetToken, user.ID, user.Version).Scan(&user.Version, &user.UpdatedAt)
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

func (s *UserStore) UpdateUserDetails(ctx context.Context, user *models.User) error {
	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		err := s.updateUser(ctx, tx, user)
		if err != nil {
			return err
		}

		err = s.updateUserProfile(ctx, tx, user.UserProfile)
		if err != nil {
			return err
		}

		return nil
	})
}

// createUser creates a new user in the database within an active transaction.
// It inserts the user record into the appropriate table and returns an error
// if the operation fails. The method ensures the operation is performed
// within the context of the provided transaction.
func (s *UserStore) createUser(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `
		INSERT INTO users(email, password, role)
		VALUES($1, $2, $3)
		RETURNING id, is_active, activated_at, version, created_at, updated_at, deleted_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(
		&user.ID,
		&user.IsActive,
		&user.ActivatedAt,
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

// createUserProfile creates a new user profile in the database after the
// user ID has been obtained. It inserts the profile record into the appropriate
// table and links it to the user. This operation is performed within the context
// of the provided transaction to ensure atomicity.
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

// deactivateActiveUser is a private method that deactivates a user if the user is currently active.
// It sets the IsActive field to false.
// This method is intended to be used when it's confirmed that the user is already active.
//
// Parameters:
//   - ctx (context.Context): The context for managing cancellation and deadlines.
//   - user (*models.User): The user object to be deactivated.
//
// Returns:
//   - error: An error if the operation fails, or nil if the deactivation was successful.
func (s *UserStore) deactivateActiveUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET is_active = 'f', version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING version, is_active, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Version,
	).Scan(
		&user.Version,
		&user.IsActive,
		&user.UpdatedAt,
	)
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

// deactivateInActiveUser is a private method that deactivates a user who is not activated.
// It updates the user's `ActivatedAt` field as it was previously `nil`, and sets the `IsActive` field to `false`.
// This method is used to deactivate users who are not yet activated.
//
// Parameters:
//   - ctx (context.Context): The context for managing cancellation and deadlines.
//   - user (*models.User): The user object to be deactivated, whose `ActivatedAt` field will be updated if it was previously `nil`.
//
// Returns:
//   - error: An error if the operation fails, or nil if the deactivation was successful.
func (s *UserStore) deactivateInActiveUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET is_active = 'f', activated_at = $1 version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version, is_active, activated_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.ActivatedAt,
		user.ID,
		user.Version,
	).Scan(
		&user.Version,
		&user.IsActive,
		&user.ActivatedAt,
		&user.UpdatedAt,
	)
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

func (s *UserStore) updateUser(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `
		UPDATE users SET email = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING email, version, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Email, user.ID, user.Version).Scan(&user.Email, &user.Version, &user.UpdatedAt)
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

func (s *UserStore) updateUserProfile(ctx context.Context, tx *sql.Tx, userProfile *models.UserProfile) error {
	query := `
		UPDATE user_profiles
		SET username = $1, profile_pic = $2, date_of_birth = $3, version = version + 1
		WHERE user_id = $4 AND version = $5
		RETURNING username, profile_pic, date_of_birth, version, updated_at
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
		userProfile.Version,
	).Scan(
		&userProfile.Username,
		&userProfile.ProfilePic,
		&userProfile.DateOfBirth,
		&userProfile.Version,
		&userProfile.UpdatedAt,
	)
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

func (s *UserStore) SoftDeleteUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET deleted_at = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING deleted_at, updated_at, version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		time.Now().UTC().Format(internal.DateTimeFormat),
		user.ID,
		user.Version,
	).Scan(
		&user.DeletedAt,
		&user.UpdatedAt,
		&user.Version,
	)

	if err != nil {
		fmt.Println(err)
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) GetByIDIcludeDeleted(ctx context.Context, ID int) (*models.User, error) {
	query := `
		SELECT u.*, up.* FROM users u
		INNER JOIN user_profiles up
		ON u.id = up.user_id
		WHERE u.id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := models.User{UserProfile: &models.UserProfile{}}
	err := s.db.
		QueryRowContext(ctx, query, ID).
		Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.IsActive,
			&user.ActivatedAt,
			&user.Role,
			&user.PasswordResetToken,
			&user.Version,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.UserProfile.ID,
			&user.UserProfile.Username,
			&user.UserProfile.ProfilePic,
			&user.UserProfile.DateOfBirth,
			&user.UserProfile.UserID,
			&user.UserProfile.Version,
			&user.UserProfile.CreatedAt,
			&user.UserProfile.UpdatedAt,
			&user.UserProfile.DeletedAt,
		)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
