package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KengoWada/meetup-clone/internal/models"
)

type OrganizationInviteStore struct {
	db *sql.DB
}

func (s *OrganizationInviteStore) Create(ctx context.Context, invite *models.OrganizationInvite) error {
	query := `
		INSERT INTO organization_invites(org_id, user_id, role_id)
		VALUES($1, $2, $3)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		invite.OrganizationID,
		invite.UserProfileID,
		invite.RoleID,
	).Scan(
		&invite.ID,
		&invite.CreatedAt,
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

func (s *OrganizationInviteStore) Get(ctx context.Context, isDeleted bool, fields []string, values []any) (*models.OrganizationInvite, error) {
	query := fmt.Sprintf("SELECT * FROM organization_invites WHERE %s ORDER BY created_at ASC", generateQueryConditions(isDeleted, fields))
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var invite = models.OrganizationInvite{}
	err := s.db.QueryRowContext(ctx, query, values...).Scan(
		&invite.ID,
		&invite.OrganizationID,
		&invite.UserProfileID,
		&invite.RoleID,
		&invite.AcceptedAt,
		&invite.DeclinedAt,
		&invite.Version,
		&invite.CreatedAt,
		&invite.UpdatedAt,
		&invite.DeletedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &invite, nil
}
