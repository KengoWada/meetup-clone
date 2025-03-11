package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/lib/pq"
)

var ErrDuplicateOrgName = errors.New("organization name is already taken")

type OrganizationStore struct {
	db *sql.DB
}

func (s *OrganizationStore) Create(ctx context.Context, organization *models.Organization, role *models.Role, member *models.OrganizationMember) error {
	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := createOrganization(ctx, tx, organization); err != nil {
			return err
		}

		role.OrganizationID = organization.ID
		if err := createRole(ctx, tx, role); err != nil {
			return err
		}

		member.OrganizationID = organization.ID
		member.RoleID = role.ID
		if err := createOrganizationMember(ctx, tx, member); err != nil {
			return err
		}

		return nil
	})
}

func (s *OrganizationStore) Get(ctx context.Context, isDeleted bool, fields []string, values []any) (*models.Organization, error) {
	query := fmt.Sprintf("SELECT * FROM organizations WHERE %s", generateQueryConditions(isDeleted, fields))

	var organization = models.Organization{}
	err := s.db.QueryRowContext(ctx, query, values...).Scan(
		&organization.ID,
		&organization.Name,
		&organization.Description,
		&organization.ProfilePic,
		&organization.IsActive,
		&organization.Version,
		&organization.CreatedAt,
		&organization.UpdatedAt,
		&organization.DeletedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &organization, nil
}

func (s *OrganizationStore) Update(ctx context.Context, organization *models.Organization) error {
	query := `
		UPDATE organizations
		SET name = $1, description = $2, profile_pic = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		organization.Name,
		organization.Description,
		organization.ProfilePic,
		organization.ID,
		organization.Version,
	).Scan(
		&organization.Version,
		&organization.UpdatedAt,
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

func createOrganization(ctx context.Context, tx *sql.Tx, organization *models.Organization) error {
	query := `
		INSERT INTO organizations(name, description, profile_pic)
		VALUES($1, $2, $3)
		RETURNING id, is_active, version, created_at, updated_at, deleted_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, organization.Name, organization.Description, organization.ProfilePic).Scan(
		&organization.ID,
		&organization.IsActive,
		&organization.Version,
		&organization.CreatedAt,
		&organization.UpdatedAt,
		&organization.DeletedAt,
	)

	if err != nil {
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "organizations_name_key"`:
			return ErrDuplicateOrgName
		default:
			return err
		}
	}

	return nil
}

func createRole(ctx context.Context, tx *sql.Tx, role *models.Role) error {
	query := `
		INSERT INTO roles(name, description, org_id, permissions)
		VALUES($1, $2, $3, $4)
		RETURNING id, version, created_at, updated_at, deleted_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return tx.QueryRowContext(ctx, query, role.Name, role.Description, role.OrganizationID, pq.Array(role.Permissions)).Scan(
		&role.ID,
		&role.Version,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
	)
}

func createOrganizationMember(ctx context.Context, tx *sql.Tx, member *models.OrganizationMember) error {
	query := `
		INSERT INTO organization_members(org_id, user_id, role_id)
		VALUES($1, $2, $3)
		RETURNING id, version, created_at, updated_at, deleted_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return tx.QueryRowContext(ctx, query, member.OrganizationID, member.UserProfileID, member.RoleID).Scan(
		&member.ID,
		&member.Version,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.DeletedAt,
	)
}
