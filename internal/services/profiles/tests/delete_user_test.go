package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	testEndpoint := "/v1/profiles"
	testMethod := http.MethodDelete

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

	generateToken := func(ID int64, isValid bool) string {
		token, err := testutils.GenerateTesAuthToken(appItems.App.Authenticator, appItems.App.Config.AuthConfig, isValid, ID)
		if err != nil {
			t.Fatal(err)
		}

		return token
	}

	t.Run("should delete user", func(t *testing.T) {
		testUser := createTestUser(true)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		fields, values := []string{"id"}, []any{testUser.ID}
		_, err = appItems.App.Store.Users.GetWithProfile(ctx, false, fields, values)
		assert.NotNil(t, err)
		assert.Equal(t, store.ErrNotFound, err)
	})

	t.Run("should not delete user invalid token", func(t *testing.T) {
		testUser := createTestUser(true)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, false)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode())

		fields, values := []string{"id"}, []any{testUser.ID}
		user, err := appItems.App.Store.Users.GetWithProfile(ctx, false, fields, values)
		assert.Nil(t, err)
		assert.Nil(t, user.DeletedAt)
	})
}
