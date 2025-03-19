package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestActivateUser(t *testing.T) {
	testEndpoint := "/v1/auth/activate"
	testMethod := http.MethodPatch

	appItems, err := app.NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	appItems.App.Store = testutils.NewTestStore(t, appItems.DB)

	mux := appItems.App.Mount()
	ctx := context.Background()

	createTestUser := func(activate bool) testutils.TestUserData {
		testUserData := testutils.NewTestUserData(activate)
		_, _, err := testUserData.CreateTestUser(ctx, appItems.App.Store, models.UserClientRole)
		if err != nil {
			t.Fatal(err)
		}
		return testUserData
	}

	createDeactivatedUser := func() testutils.TestUserData {
		testUserData := testutils.NewTestUserData(true)
		_, err := testUserData.CreateDeactivatedTestUser(ctx, appItems.App.Store, models.UserClientRole)
		if err != nil {
			t.Fatal(err)
		}
		return testUserData
	}

	generateToken := func(email string, isValid bool) string {
		var createdAt string = time.Now().UTC().Format(internal.DateTimeFormat)
		if !isValid {
			createdAt = time.Now().Add(-time.Hour).UTC().Format(internal.DateTimeFormat)
		}

		token, err := utils.GenerateTestToken(email, []byte(appItems.App.Config.SecretKey), createdAt)
		if err != nil {
			t.Fatal("failed to generate test token to activate a user")
		}

		return token
	}

	t.Run("should activate user", func(t *testing.T) {
		testUserData := createTestUser(false)

		data := testutils.TestRequestData{"token": generateToken(testUserData.Email, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, "Email successfully verified", response.GetMessage())
	})

	t.Run("should not activate if the request has an unknown field", func(t *testing.T) {
		testUserData := createTestUser(false)

		const unknownField = "fakeField"
		data := testutils.TestRequestData{
			"token":      generateToken(testUserData.Email, true),
			unknownField: "random data :)",
		}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "Unknown field in request", response.GetMessage())
		assert.Equal(t, "unknown field", errorMessages[unknownField])
	})

	t.Run("should not activate if token is expired", func(t *testing.T) {
		testUserData := createTestUser(false)

		data := testutils.TestRequestData{"token": generateToken(testUserData.Email, false)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
		assert.Equal(t, "Activation token has exipred", response.GetMessage())
	})

	t.Run("should not activate if email is invalid", func(t *testing.T) {
		email, _ := testutils.GenerateEmailAndUsername()

		data := testutils.TestRequestData{"token": generateToken(email, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Activation token is invalid", response.GetMessage())
	})

	t.Run("should not activate already active user", func(t *testing.T) {
		testUserData := createTestUser(true)

		data := testutils.TestRequestData{"token": generateToken(testUserData.Email, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Activation token is invalid", response.GetMessage())
	})

	t.Run("should not activate deactivated user", func(t *testing.T) {
		testUserData := createDeactivatedUser()

		data := testutils.TestRequestData{"token": generateToken(testUserData.Email, true)}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Activation token is invalid", response.GetMessage())
	})
}
