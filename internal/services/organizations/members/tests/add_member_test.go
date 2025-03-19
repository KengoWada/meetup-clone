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

func TestAddMember(t *testing.T) {
	testEndpoint := func(orgID int64) string {
		return fmt.Sprintf("/v1/organizations/%d/members", orgID)
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

	createTestRole := func(isDeleted bool, orgID int64) *models.Role {
		permissions := []string{internal.MemberRemove}
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
			role.Permissions = []string{internal.MemberAdd}
		case "invalid":
			role.Permissions = []string{internal.MemberRoleUpdate, internal.MemberRemove}
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

	t.Run("should invite organization member", func(t *testing.T) {
		testUser := createTestUser(true)
		invitedTestUser := createTestUser(true)
		org := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, org.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"roleId": testRole.ID,
			"email":  invitedTestUser.Email,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(org.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, response.StatusCode())
		assert.Equal(t, "Invite sent", response.GetMessage())
	})

	t.Run("should not invite user not authenticated", func(t *testing.T) {
		testUser := createTestUser(true)
		invitedTestUser := createTestUser(true)
		org := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, org.ID)

		payload := testutils.TestRequestData{
			"roleId": testRole.ID,
			"email":  invitedTestUser.Email,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(org.ID), nil, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
		assert.Equal(t, "unauthorized", response.GetMessage())
	})

	t.Run("should not invite member invalid permissions", func(t *testing.T) {
		testUser := createTestUser(true)
		invitedTestUser := createTestUser(true)
		org := createTestOrg(true, generateRole("invalid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, org.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"roleId": testRole.ID,
			"email":  invitedTestUser.Email,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(org.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})

	t.Run("should not invite member when not group member", func(t *testing.T) {
		testUser := createTestUser(true)
		testUserTwo := createTestUser(true)
		invitedTestUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		createTestOrg(true, generateRole("valid"), testUserTwo.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUserTwo.ID, true)}
		payload := testutils.TestRequestData{
			"roleId": testRole.ID,
			"email":  invitedTestUser.Email,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})

	t.Run("should not invite member invalid email", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"roleId": testRole.ID,
			"email":  "fakemail.com",
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid request body", response.GetMessage())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "Invalid email address provided", errorMessages["email"])
	})

	t.Run("should return invite sent for email not attached to user", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)
		email, _ := testutils.GenerateEmailAndUsername()

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"roleId": testRole.ID,
			"email":  email,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, response.StatusCode())
		assert.Equal(t, "Invite sent", response.GetMessage())
	})

	t.Run("should not invite user invalid role ID", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		invitedTestUser := createTestUser(true)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"roleId": -12,
			"email":  invitedTestUser.Email,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid request body", response.GetMessage())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "invalid role id", errorMessages["roleId"])
	})

	t.Run("should not invite member with unknown field", func(t *testing.T) {
		testUser := createTestUser(true)
		testOrg := createTestOrg(true, generateRole("valid"), testUser.UserProfile.ID)
		testRole := createTestRole(false, testOrg.ID)
		invitedTestUser := createTestUser(true)
		unknownField := "unknownField"

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"roleId":     testRole.ID,
			"email":      invitedTestUser.Email,
			unknownField: "fakeData",
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testOrg.ID), headers, payload)
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
}
