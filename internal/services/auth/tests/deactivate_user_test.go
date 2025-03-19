package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDeactivateUser(t *testing.T) {
	testMethod := http.MethodPatch
	testEndpoint := func(ID int64) string {
		return fmt.Sprintf("/v1/auth/users/%d/deactivate", ID)
	}

	appItems, err := app.NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	appItems.App.Store = testutils.NewTestStore(t, appItems.DB)

	mux := appItems.App.Mount()
	ctx := context.Background()

	createTestUser := func(activate bool, role models.UserRole) *models.User {
		testUserData := testutils.NewTestUserData(activate)
		user, userProfile, err := testUserData.CreateTestUser(ctx, appItems.App.Store, role)
		if err != nil {
			t.Fatal(err)
		}
		user.UserProfile = userProfile
		return user
	}

	createDeactivatedUser := func(role models.UserRole) *models.User {
		testUserData := testutils.NewTestUserData(true)
		user, err := testUserData.CreateDeactivatedTestUser(ctx, appItems.App.Store, role)
		if err != nil {
			t.Fatal(err)
		}
		return user
	}

	generateToken := func(ID int64, isValid bool) string {
		token, err := testutils.GenerateTesAuthToken(appItems.App.Authenticator, appItems.App.Config.AuthConfig, isValid, ID)
		if err != nil {
			t.Fatal(err)
		}

		return token
	}

	t.Run("should deactivate user", func(t *testing.T) {
		adminTestUser := createTestUser(true, models.UserAdminRole)
		testUser := createTestUser(true, models.UserClientRole)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(adminTestUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testUser.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, "Done", response.GetMessage())
	})

	t.Run("should not deactivate user non staff user", func(t *testing.T) {
		clientTestUser := createTestUser(true, models.UserClientRole)
		testUser := createTestUser(true, models.UserClientRole)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(clientTestUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testUser.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})

	t.Run("should not deactivate user invalid user ID", func(t *testing.T) {
		staffTestUser := createTestUser(true, models.UserStaffRole)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(staffTestUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(0), headers, nil)
		if err != nil {
			t.Fatal(err)
		}

		const errorMessage = "Invalid user ID"
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, errorMessage, response.GetMessage())

		// string id
		endpoint := "/v1/auth/users/someID/deactivate"
		response, err = testutils.RunTestRequest(mux, testMethod, endpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, errorMessage, response.GetMessage())
	})

	t.Run("should not deactivate deactivated user", func(t *testing.T) {
		staffTestUser := createTestUser(true, models.UserStaffRole)
		testUser := createDeactivatedUser(models.UserClientRole)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(staffTestUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint(testUser.ID), headers, nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "User is already deactivated", response.GetMessage())
	})
}
