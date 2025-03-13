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

func TestCreateOrgRole(t *testing.T) {
	testEndpoint := func(orgID int64) string {
		return fmt.Sprintf("/v1/organizations/%d/roles", orgID)
	}
	testMethod := http.MethodPost

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
		role := &models.Role{Name: faker.Username(options.WithGenerateUniqueValues(true))}

		switch permission {
		case "valid":
			role.Permissions = []string{internal.RoleCreate}
		case "invalid":
			role.Permissions = []string{internal.RoleDelete, internal.RoleUpdate}
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

	t.Run("should create org role", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"name":        faker.Username(options.WithGenerateUniqueValues(true)),
			"description": "Simple Description",
			"permissions": []string{internal.RoleCreate, internal.RoleUpdate},
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(org.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, response.StatusCode())
		assert.Equal(t, "Done", response.GetMessage())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, payload["name"], data["name"])
		assert.Equal(t, payload["description"], data["description"])
		assert.ElementsMatch(t, payload["permissions"], data["permissions"])
	})

	t.Run("should not create role unknown field", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		unknownField := "unknownField"
		payload := testutils.TestRequestData{
			"name":        faker.Username(options.WithGenerateUniqueValues(true)),
			"description": "Simple Description",
			"permissions": []string{internal.RoleCreate, internal.RoleUpdate},
			unknownField:  "fakeData",
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(org.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Unknown field in request", response.GetMessage())

		errMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		assert.Equal(t, "unknown field", errMessages[unknownField])
	})

	t.Run("should not create org role with invalid data", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"name":        faker.Username(options.WithGenerateUniqueValues(true)),
			"description": "Simple Description",
			"permissions": []string{internal.RoleCreate, internal.RoleUpdate, internal.RoleUpdate},
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(org.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid request body", response.GetMessage())

		errMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		assert.Equal(t, "Duplicate permissions are not allowed", errMessages["permissions"])
	})

	t.Run("should not create org role with duplicate name", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"name":        faker.Username(options.WithGenerateUniqueValues(true)),
			"description": "Simple Description",
			"permissions": []string{internal.RoleCreate, internal.RoleUpdate},
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(org.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, response.StatusCode())
		assert.Equal(t, "Done", response.GetMessage())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, payload["name"], data["name"])
		assert.Equal(t, payload["description"], data["description"])
		assert.ElementsMatch(t, payload["permissions"], data["permissions"])

		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint(org.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid request body", response.GetMessage())

		errMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		assert.Equal(t, "name already exists", errMessages["name"])
	})
}
