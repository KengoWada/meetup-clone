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

func TestUpdateOrganizationDetails(t *testing.T) {
	testEndpoint := "/v1/organizations"
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

	generateRole := func(permission string) *models.Role {
		role := &models.Role{Name: faker.Username(options.WithGenerateUniqueValues(true))}

		switch permission {
		case "valid":
			role.Permissions = []string{internal.OrgUpdate}
		case "invalid":
			role.Permissions = []string{internal.OrgDeactivate, internal.OrgDelete}
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

	t.Run("should update organization details", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newOrgName := faker.Username(options.WithGenerateUniqueValues(true))
		payload := testutils.TestRequestData{
			"name":        newOrgName,
			"description": newOrgName + "description",
			"profilePic":  testutils.TestProfilePic,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, payload)
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
		assert.Equal(t, payload["profilePic"], data["profilePic"])
	})

	t.Run("should not update org not authenticated", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)
		newOrgName := faker.Username(options.WithGenerateUniqueValues(true))
		payload := testutils.TestRequestData{
			"name":        newOrgName,
			"description": newOrgName + "description",
			"profilePic":  testutils.TestProfilePic,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, nil, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
		assert.Equal(t, "unauthorized", response.GetMessage())
	})

	t.Run("should not update org with no permissions", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, generateRole("invalid"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newOrgName := faker.Username(options.WithGenerateUniqueValues(true))
		payload := testutils.TestRequestData{
			"name":        newOrgName,
			"description": newOrgName + "description",
			"profilePic":  testutils.TestProfilePic,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})

	t.Run("should not update org with unknown field in request", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newOrgName := faker.Username(options.WithGenerateUniqueValues(true))
		unknownField := "unknownField"
		payload := testutils.TestRequestData{
			"name":        newOrgName,
			"description": newOrgName + "description",
			"profilePic":  testutils.TestProfilePic,
			unknownField:  "fakeData",
		}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "unknown field", errMessages[unknownField])
	})

	t.Run("should not update org invalid request data", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newOrgName := "B@Dnam3"
		payload := testutils.TestRequestData{
			"name":        newOrgName,
			"description": newOrgName + "description",
			"profilePic":  testutils.TestProfilePic,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		const nameErr = "Organization name can only include alphanumeric characters and spaces"
		assert.Equal(t, nameErr, errMessages["name"])
	})

	t.Run("should not update org with invalid orgID", func(t *testing.T) {
		testUser := createTestUser(true)
		createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%s", testEndpoint, "some")
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newOrgName := faker.Username(options.WithGenerateUniqueValues(true))
		payload := testutils.TestRequestData{
			"name":        newOrgName,
			"description": newOrgName + "description",
			"profilePic":  testutils.TestProfilePic,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid organization ID", response.GetMessage())

		endpoint = fmt.Sprintf("%s/%d", testEndpoint, 0)
		response, err = testutils.RunTestRequest(mux, testMethod, endpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})

	t.Run("should not update org if org is deactivated", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(false, generateRole("all"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		newOrgName := faker.Username(options.WithGenerateUniqueValues(true))
		payload := testutils.TestRequestData{
			"name":        newOrgName,
			"description": newOrgName + "description",
			"profilePic":  testutils.TestProfilePic,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})
}
