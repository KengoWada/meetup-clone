package testutils

import (
	"context"
	"fmt"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

func CreateTestOrganization(ctx context.Context, appStore store.Store, isActive bool, role *models.Role, userID int64) (*models.Organization, error) {
	name := faker.Username(options.WithGenerateUniqueValues(true))
	organization := &models.Organization{
		Name:        name,
		Description: fmt.Sprintf("%s Description", name),
		ProfilePic:  TestProfilePic,
	}
	member := &models.OrganizationMember{UserProfileID: userID}

	err := appStore.Organizations.Create(ctx, organization, role, member)
	if err != nil {
		return nil, err
	}

	if !isActive {
		if err := appStore.Organizations.Deactivate(ctx, organization); err != nil {
			return nil, err
		}
	}

	return organization, nil
}

func CreateTestRole(ctx context.Context, appStore store.Store, isDeleted bool, orgID int64, permissions []string) (*models.Role, error) {
	name := faker.Username(options.WithGenerateUniqueValues(true))
	role := &models.Role{
		Name:           name,
		Description:    fmt.Sprintf("%s Description", name),
		OrganizationID: orgID,
		Permissions:    permissions,
	}

	if err := appStore.Roles.Create(ctx, role); err != nil {
		return nil, err
	}

	if isDeleted {
		if err := appStore.Roles.SoftDelete(ctx, role); err != nil {
			return nil, err
		}
	}

	return role, nil
}

func CreateTestOrganizationMember(ctx context.Context, appStore store.Store, orgID, roleID, userID int64) (*models.OrganizationMember, error) {
	member := &models.OrganizationMember{
		OrganizationID: orgID,
		RoleID:         roleID,
		UserProfileID:  userID,
	}

	if err := appStore.OrganizationMembers.Create(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}
