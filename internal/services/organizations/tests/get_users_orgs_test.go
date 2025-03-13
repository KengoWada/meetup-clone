package tests

import (
	"context"
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

func TestGetUsersOrganizations(t *testing.T) {
	testEndpoint := "/v1/organizations"
	testMethod := http.MethodGet

	appItems, err := app.NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	appItems.App.Store = testutils.NewTestStore(t, appItems.DB)

	mux := appItems.App.Mount()
	ctx := context.Background()
	var role = &models.Role{
		Name:        faker.Username(options.WithGenerateUniqueValues(true)),
		Permissions: []string{internal.OrgDeactivate, internal.OrgUpdate, internal.OrgDelete},
	}

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

	generateToken := func(ID int64, isValid bool) string {
		token, err := testutils.GenerateTesAuthToken(appItems.App.Authenticator, appItems.App.Config.AuthConfig, isValid, ID)
		if err != nil {
			t.Fatal(err)
		}

		return token
	}

	createOrgs := func(num int, activate bool, userID int64) []*models.Organization {
		var orgs []*models.Organization
		for range num {
			org := createTestOrg(activate, role, userID)
			orgs = append(orgs, org)
		}

		return orgs
	}

	t.Run("should get a users organizations", func(t *testing.T) {
		testUser := createTestUser(true)
		createOrgs(5, true, testUser.UserProfile.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response data to map")
		}

		organizations, ok := data["organizations"].([]any)
		if !ok {
			t.Fatal("failed to convert response data to slice")
		}

		assert.Equal(t, 5, len(organizations))
	})

	t.Run("should not fetch deactivated organizations", func(t *testing.T) {
		testUser := createTestUser(true)
		activeOrgs := createOrgs(5, true, testUser.UserProfile.ID)
		createOrgs(5, false, testUser.UserProfile.ID)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response data to map")
		}

		organizations, ok := data["organizations"].([]any)
		if !ok {
			t.Fatal("failed to convert response data to slice")
		}

		var activeOrgIDs []int64
		var responseOrgIDs []int64

		for _, org := range activeOrgs {
			activeOrgIDs = append(activeOrgIDs, org.ID)
		}

		for _, org := range organizations {
			orgData, ok := org.(map[string]any)
			if !ok {
				t.Fatal("failed to convert response data to map")
			}
			orgID := int64(orgData["id"].(float64))
			responseOrgIDs = append(responseOrgIDs, orgID)
		}

		assert.Equal(t, 5, len(organizations))
		assert.ElementsMatch(t, activeOrgIDs, responseOrgIDs)
	})

	t.Run("should return null for no orgs associated with the user", func(t *testing.T) {
		testUser := createTestUser(true)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response data to map")
		}
		assert.Nil(t, data["organizations"])
	})
}
