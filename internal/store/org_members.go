package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/KengoWada/meetup-clone/internal/models"
)

type OrganizationMembersStore struct {
	db *sql.DB
}

func (s *OrganizationMembersStore) Get(ctx context.Context, isDeleted bool, fields []string, values []any) (*models.OrganizationMember, error) {
	var queryConditions []string
	for index, field := range fields {
		queryField := fmt.Sprintf("%s = $%d", field, index+1)
		queryConditions = append(queryConditions, queryField)
	}

	if !isDeleted {
		queryConditions = append(queryConditions, "deleted_at IS NULL")
	}

	query := fmt.Sprintf(
		"SELECT * FROM organization_members WHERE %s",
		strings.Join(queryConditions, " AND "),
	)

	var member = models.OrganizationMember{}
	err := s.db.QueryRowContext(ctx, query, values...).Scan(
		&member.ID,
		&member.OrganizationID,
		&member.UserProfileID,
		&member.RoleID,
		&member.Version,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.DeletedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &member, nil
}
