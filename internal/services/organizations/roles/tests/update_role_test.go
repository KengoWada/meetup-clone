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

func TestUpdateRole(t *testing.T) {
	testEndpoint := func(orgID, roleID int64) string {
		return fmt.Sprintf("/v1/organizations/%d/roles/%d", orgID, roleID)
	}
	testMethod := http.MethodPut

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

	generateRole := func(permission string) *models.Role {
		role := &models.Role{Name: faker.Username(options.WithGenerateUniqueValues(true))}

		switch permission {
		case "valid":
			role.Permissions = []string{internal.RoleUpdate}
		case "invalid":
			role.Permissions = []string{internal.RoleDelete, internal.RoleCreate}
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

	t.Run("should update role", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newRoleName := faker.Username(options.WithGenerateUniqueValues(true))
		payload := testutils.TestRequestData{
			"name":        newRoleName,
			"description": newRoleName + " description",
			"permissions": []string{internal.RoleCreate, internal.RoleDelete},
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID, testRole.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response data to map")
		}

		assert.Equal(t, payload["name"], data["name"])
		assert.Equal(t, payload["description"], data["description"])
		assert.ElementsMatch(t, payload["permissions"], data["permissions"])
	})

	t.Run("should not update role with invalid permission", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("invalid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newRoleName := faker.Username(options.WithGenerateUniqueValues(true))
		payload := testutils.TestRequestData{
			"name":        newRoleName,
			"description": newRoleName + " description",
			"permissions": []string{internal.RoleCreate, internal.RoleDelete},
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID, testRole.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})

	t.Run("should not update role unknown field", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newRoleName := faker.Username(options.WithGenerateUniqueValues(true))
		unknownField := "unknownField"
		payload := testutils.TestRequestData{
			"name":        newRoleName,
			"description": newRoleName + " description",
			"permissions": []string{internal.RoleCreate, internal.RoleDelete},
			unknownField:  "fakeData",
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID, testRole.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Unknown field in request", response.GetMessage())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "unknown field", errorMessages[unknownField])
	})

	t.Run("should not update role with invalid data", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newRoleName := faker.Username(options.WithGenerateUniqueValues(true))
		payload := testutils.TestRequestData{
			"name":        newRoleName,
			"description": newRoleName + " description",
			"permissions": []string{internal.RoleCreate, ""},
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID, testRole.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid request body", response.GetMessage())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "Invalid permission sent", errorMessages["permissions"])
	})

	t.Run("should not update role with duplicate name", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)
		testRoleTwo := createTestRole(false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"name":        testRoleTwo.Name,
			"description": testRoleTwo.Name + " description",
			"permissions": []string{internal.RoleCreate},
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID, testRole.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid request body", response.GetMessage())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "name already exists", errorMessages["name"])
	})
}
