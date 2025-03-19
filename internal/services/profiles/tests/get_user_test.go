package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetUserProfile(t *testing.T) {
	testEndpoint := "/v1/profiles"
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

	createDeactivatedUser := func() *models.User {
		testUserData := testutils.NewTestUserData(true)
		user, err := testUserData.CreateDeactivatedTestUser(ctx, appItems.App.Store, models.UserClientRole)
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

	t.Run("should get a users profile details", func(t *testing.T) {
		testUser := createTestUser(true)
		token := generateToken(testUser.ID, true)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + token}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		dobOriginal, err := time.Parse("01/02/2006", testUser.UserProfile.DateOfBirth)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, testUser.Email, data["email"])
		assert.Equal(t, testUser.UserProfile.Username, data["username"])
		assert.Equal(t, testUser.UserProfile.ProfilePic, data["profilePic"])
		assert.Equal(t, dobOriginal.Format(time.RFC3339), data["dateOfBirth"])
	})

	t.Run("should not fetch details for deactivated user", func(t *testing.T) {
		testUser := createDeactivatedUser()

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
		assert.Equal(t, "unauthorized", response.GetMessage())
	})
}
