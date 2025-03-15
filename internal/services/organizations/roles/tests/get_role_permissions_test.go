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

func TestGetRolePermission(t *testing.T) {
	testEndpoint := func(orgID int64) string {
		return fmt.Sprintf("/v1/organizations/%d/roles/permissions", orgID)
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

	generateRole := func(permission string) *models.Role {
		roleName := faker.Username(options.WithGenerateUniqueValues(true))
		role := &models.Role{
			Name:        roleName,
			Description: roleName + " description",
		}

		switch permission {
		case "valid":
			role.Permissions = []string{internal.RoleCreate, internal.RoleUpdate}
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

	t.Run("should fetch role permissions", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)

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

		permissions, ok := data["permissions"].(map[string]any)
		if !ok {
			t.Fatal("failed to convert response permissions to map")
		}

		assert.ElementsMatch(t, internal.PermissionsMap["events"], permissions["events"])
		assert.ElementsMatch(t, internal.PermissionsMap["members"], permissions["members"])
		assert.ElementsMatch(t, internal.PermissionsMap["roles"], permissions["roles"])
		assert.ElementsMatch(t, internal.PermissionsMap["organizations"], permissions["organizations"])
	})

	t.Run("should not fetch permissions with invalid permissions", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("invalid"), testUser.UserProfile.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})
}
