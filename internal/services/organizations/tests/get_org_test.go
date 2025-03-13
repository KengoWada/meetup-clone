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

func TestGetOrg(t *testing.T) {
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

	t.Run("should get an organization", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, role, testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response data to map")
		}

		assert.Equal(t, org.Name, data["name"])
		assert.Equal(t, org.Description, data["description"])
		assert.Equal(t, org.ProfilePic, data["profilePic"])
	})

	t.Run("should not get organization not authenticated", func(t *testing.T) {
		testUser := createTestUser(true)
		org := createTestOrg(true, role, testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
		assert.Equal(t, "unauthorized", response.GetMessage())
	})

	t.Run("should not get organization invalid id", func(t *testing.T) {
		testUser := createTestUser(true)
		_ = createTestOrg(true, role, testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%s", testEndpoint, "word")
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid organization ID", response.GetMessage())

		endpoint = fmt.Sprintf("%s/0", testEndpoint)
		response, err = testutils.RunTestRequest(mux, testMethod, endpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid organization ID", response.GetMessage())
	})
}
