package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/lib/pq"
)

const createRoleQuery = `
	INSERT INTO roles(name, description, org_id, permissions)
	VALUES($1, $2, $3, $4)
	RETURNING id, version, created_at, updated_at, deleted_at
`

type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) Create(ctx context.Context, role *models.Role) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	values := []any{role.Name, role.Description, role.OrganizationID, pq.Array(role.Permissions)}
	return s.db.QueryRowContext(ctx, createRoleQuery, values...).Scan(
		&role.ID,
		&role.Version,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
	)
}

func (s *RoleStore) Get(ctx context.Context, isDeleted bool, fields []string, values []any) (*models.Role, error) {
	var queryConditions []string
	for index, field := range fields {
		queryField := fmt.Sprintf("%s = $%d", field, index+1)
		queryConditions = append(queryConditions, queryField)
	}

	if !isDeleted {
		queryConditions = append(queryConditions, "deleted_at IS NULL")
	}

	query := fmt.Sprintf(
		"SELECT * FROM roles WHERE %s",
		strings.Join(queryConditions, " AND "),
	)
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var role models.Role
	err := s.db.QueryRowContext(ctx, query, values...).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.OrganizationID,
		pq.Array(&role.Permissions),
		&role.Version,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &role, nil
}

func (s *RoleStore) GetByOrgID(ctx context.Context, orgID int64) ([]*models.SimpleRole, error) {
	query := `
		SELECT id, name, description, permissions FROM roles
		WHERE org_id = $1 AND deleted_at IS NULL
		ORDER BY name ASC, created_at ASC
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.SimpleRole
	for rows.Next() {
		var role models.SimpleRole
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			pq.Array(&role.Permissions),
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, &role)
	}

	return roles, nil
}

func (s *RoleStore) Update(ctx context.Context, role *models.Role) error {
	query := `
		UPDATE roles
		SET name = $1, description = $2, permissions = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		role.Name,
		role.Description,
		pq.Array(role.Permissions),
		role.ID,
		role.Version,
	).Scan(
		&role.Version,
		&role.UpdatedAt,
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

func createRoleTx(ctx context.Context, tx *sql.Tx, role *models.Role) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	values := []any{role.Name, role.Description, role.OrganizationID, pq.Array(role.Permissions)}
	return tx.QueryRowContext(ctx, createRoleQuery, values...).Scan(
		&role.ID,
		&role.Version,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
	)
}
