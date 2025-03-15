package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRole(t *testing.T) {
	testEndpoint := func(orgID, roleID int64) string {
		return fmt.Sprintf("/v1/organizations/%d/roles/%d", orgID, roleID)
	}
	testMethod := http.MethodDelete

	appItems, err := app.NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	appItems.App.Store = testutils.NewTestStore(t, appItems.DB)

	mux := appItems.App.Mount()
	ctx := context.Background()

	createTestUser := func(activate bool) *models.User {
		testUserData := testutils.NewTestUserData(activate)
		user, userProfile, err := testUserData.CreateTestUser(ctx, appItems.App.Store, models.UserClientRole)
		if err != nil {
			t.Fatal(err)
		}

		user.UserProfile = userProfile
		return user
	}

	createTestOrg := func(isActive bool, role *models.Role, userID int64) *models.Organization {
		org, err := testutils.CreateTestOrganization(ctx, appItems.App.Store, isActive, role, userID)
		if err != nil {
			t.Fatal(err)
		}

		return org
	}

	createTestRole := func(isDeleted bool, orgID int64) *models.Role {
		permissions := []string{internal.RoleCreate, internal.RoleUpdate}
		role, err := testutils.CreateTestRole(ctx, appItems.App.Store, isDeleted, orgID, permissions)
		if err != nil {
			t.Fatal(err)
		}

		return role
	}

	addMemberToOrg := func(count int, orgID, roleID int64) {
		for range count {
			user := createTestUser(true)
			_, err := testutils.CreateTestOrganizationMember(ctx, appItems.App.Store, orgID, roleID, user.UserProfile.ID)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	generateRole := func(permission string) *models.Role {
		role := &models.Role{Name: faker.Username(options.WithGenerateUniqueValues(true))}

		switch permission {
		case "valid":
			role.Permissions = []string{internal.RoleDelete}
		case "invalid":
			role.Permissions = []string{internal.RoleUpdate, internal.RoleCreate}
		default:
			role.Permissions = internal.Permissions
		}

		return role
	}

	generateToken := func(ID int64, isValid bool) string {
		token, err := testutils.GenerateTesAuthToken(appItems.App.Authenticator, appItems.App.Config.AuthConfig, isValid, ID)
		if err != nil {
			t.Fatal(err)
		}

		return token
	}

	t.Run("should delete user role", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID, testRole.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, "Done", response.GetMessage())
	})

	t.Run("should not delete role attached to active members", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)
		addMemberToOrg(3, testOrg.ID, testRole.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID, testRole.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Role is assigned to active users. Please reassign them before deleting.", response.GetMessage())
	})

	t.Run("should not delete role with invalid permissions", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("invalid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID, testRole.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})
}
