package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/KengoWada/meetup-clone/internal/models"
)

const createOrgMemberQuery = `
	INSERT INTO organization_members(org_id, user_id, role_id)
	VALUES($1, $2, $3)
	RETURNING id, version, created_at, updated_at, deleted_at
`

type OrganizationMembersStore struct {
	db *sql.DB
}

func (s *OrganizationMembersStore) Create(ctx context.Context, member *models.OrganizationMember) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	values := []any{member.OrganizationID, member.UserProfileID, member.RoleID}
	return s.db.QueryRowContext(ctx, createOrgMemberQuery, values...).Scan(
		&member.ID,
		&member.Version,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.DeletedAt,
	)
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func createOrgMemberTx(ctx context.Context, tx *sql.Tx, member *models.OrganizationMember) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	values := []any{member.OrganizationID, member.UserProfileID, member.RoleID}
	return tx.QueryRowContext(ctx, createOrgMemberQuery, values...).Scan(
		&member.ID,
		&member.Version,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.DeletedAt,
	)
}
