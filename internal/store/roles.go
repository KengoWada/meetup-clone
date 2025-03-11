package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/lib/pq"
)

type RoleStore struct {
	db *sql.DB
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
