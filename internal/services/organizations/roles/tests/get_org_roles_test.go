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

func TestGetOrgRoles(t *testing.T) {
	testEndpoint := func(orgID int64) string {
		return fmt.Sprintf("/v1/organizations/%d/roles", orgID)
	}
	testMethod := http.MethodGet

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

	createTestRoles := func(num int, isDeleted bool, orgID int64) []*models.Role {
		permissions := []string{internal.RoleCreate, internal.RoleUpdate}
		roles := make([]*models.Role, 0)
		for range num {
			role, err := testutils.CreateTestRole(ctx, appItems.App.Store, isDeleted, orgID, permissions)
			if err != nil {
				t.Fatal(err)
			}

			roles = append(roles, role)
		}

		return roles
	}

	generateRole := func(permission string) *models.Role {
		roleName := faker.Username(options.WithGenerateUniqueValues(true))
		role := &models.Role{
			Name:        roleName,
			Description: roleName + " description",
		}

		switch permission {
		case "valid":
			role.Permissions = []string{internal.RoleUpdate, internal.RoleDelete, internal.MemberAdd, internal.MemberRoleUpdate}
		case "invalid":
			role.Permissions = []string{internal.OrgUpdate}
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

	t.Run("should fetch organization roles", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)
		testRoles := createTestRoles(4, false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response data to map")
		}

		roles, ok := data["roles"].([]any)
		if !ok {
			t.Fatal("failed to convert response roles to slice of map")
		}

		// Includes the role created when starting an organization
		assert.Equal(t, 5, len(roles))

		roleIDs := make([]int64, 0)
		for _, role := range roles {
			roleData, ok := role.(map[string]any)
			if !ok {
				t.Fatal("failed to convert response data to map")
			}
			roleID := int64(roleData["id"].(float64))
			roleIDs = append(roleIDs, roleID)
		}

		for _, role := range testRoles {
			assert.Contains(t, roleIDs, role.ID)
		}
	})

	t.Run("should not fetch organization roles when not organization member", func(t *testing.T) {
		testUser := createTestUser(true)
		testUserTwo := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		createTestRoles(4, false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUserTwo.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})

	t.Run("should return only roles that are not deleted", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRoles := createTestRoles(4, false, testOrg.ID)
		createTestRoles(4, true, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		roles, ok := data["roles"].([]any)
		if !ok {
			t.Fatal("failed to convert response roles to slice of map")
		}

		// Includes the role created when starting an organization
		assert.Equal(t, 5, len(roles))

		roleIDs := make([]int64, 0)
		for _, role := range roles {
			roleData, ok := role.(map[string]any)
			if !ok {
				t.Fatal("failed to convert response data to map")
			}
			roleID := int64(roleData["id"].(float64))
			roleIDs = append(roleIDs, roleID)
		}

		for _, role := range testRoles {
			assert.Contains(t, roleIDs, role.ID)
		}
	})

	t.Run("should not fetch roles with invalid permissions", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("invalid"), testUser.UserProfile.ID)
		createTestRoles(4, false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})
}
