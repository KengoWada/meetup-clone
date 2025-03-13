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

var ErrDuplicateOrgName = errors.New("organization name is already taken")

type OrganizationStore struct {
	db *sql.DB
}

func (s *OrganizationStore) Create(ctx context.Context, organization *models.Organization, role *models.Role, member *models.OrganizationMember) error {
	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := createOrganizationTx(ctx, tx, organization); err != nil {
			return err
		}

		role.OrganizationID = organization.ID
		if err := createRoleTx(ctx, tx, role); err != nil {
			return err
		}

		member.OrganizationID = organization.ID
		member.RoleID = role.ID
		if err := createOrgMemberTx(ctx, tx, member); err != nil {
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

func (s *OrganizationStore) Deactivate(ctx context.Context, organization *models.Organization) error {
	query := `
		UPDATE organizations
		SET is_active = 'f', version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING is_active, version, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		organization.ID,
		organization.Version,
	).Scan(
		&organization.IsActive,
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

func (s *OrganizationStore) SoftDelete(ctx context.Context, organization *models.Organization) error {
	query := `
		UPDATE organizations
		SET deleted_at = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version, updated_at, deleted_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	values := []any{time.Now().UTC().Format(internal.DateTimeFormat), organization.ID, organization.Version}
	err := s.db.QueryRowContext(ctx, query, values...).Scan(
		&organization.Version,
		&organization.UpdatedAt,
		&organization.DeletedAt,
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

func createOrganizationTx(ctx context.Context, tx *sql.Tx, organization *models.Organization) error {
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
