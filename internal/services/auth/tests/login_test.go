package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestLoginUser(t *testing.T) {
	testEndpoint := "/v1/auth/login"
	testMethod := http.MethodPost

	appItems, err := app.NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	appItems.App.Store = testutils.NewTestStore(t, appItems.DB)

	mux := appItems.App.Mount()
	ctx := context.Background()

	createTestUser := func(activate bool) testutils.TestUserData {
		testUserData := testutils.NewTestUserData(activate)
		_, _, err := testUserData.CreateTestUser(ctx, appItems.App.Store)
		if err != nil {
			t.Fatal(err)
		}
		return testUserData
	}

	createDeactivatedUser := func() testutils.TestUserData {
		testUserData := testutils.NewTestUserData(true)
		_, err := testUserData.CreateDeactivatedTestUser(ctx, appItems.App.Store)
		if err != nil {
			t.Fatal(err)
		}
		return testUserData
	}

	t.Run("should log in a user", func(t *testing.T) {
		testUserData := createTestUser(true)
		data := testutils.TestRequestData{"email": testUserData.Email, "password": testUserData.Password}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		token, ok := data["token"]
		assert.True(t, ok)
		_, err = appItems.App.Authenticator.ValidateToken(token.(string))
		assert.Nil(t, err)
	})

	t.Run("should not log in with no credentials provided", func(t *testing.T) {
		data := testutils.TestRequestData{"email": "", "password": ""}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		var responseMessage = "Invalid request body"
		var requiredFieldMessage = "Field is required"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, requiredFieldMessage, errorMessages["email"])
		assert.Equal(t, requiredFieldMessage, errorMessages["password"])
	})

	t.Run("should not log in if account doesn't exist", func(t *testing.T) {
		email, _ := testutils.GenerateEmailAndUsername()
		data := testutils.TestRequestData{"email": email, "password": testutils.TestPassword}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		var responseMessage = "Invalid credentials"
		assert.Equal(t, responseMessage, response.GetMessage())
	})

	t.Run("should not log in if account is not active", func(t *testing.T) {
		testUserData := createTestUser(false)
		data := testutils.TestRequestData{"email": testUserData.Email, "password": testUserData.Password}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())

		var responseMessage = "Please verify your email address to proceed."
		assert.Equal(t, responseMessage, response.GetMessage())
	})

	t.Run("should not log in with invalid password", func(t *testing.T) {
		testUserData := createTestUser(true)
		data := testutils.TestRequestData{"email": testUserData.Email, "password": "wrong_password"}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		var responseMessage = "Invalid credentials"
		assert.Equal(t, responseMessage, response.GetMessage())
	})

	t.Run("should not log in a deactivated user", func(t *testing.T) {
		testUserData := createDeactivatedUser()
		data := testutils.TestRequestData{"email": testUserData.Email, "password": testUserData.Password}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		var responseMessage = "Invalid credentials"
		assert.Equal(t, responseMessage, response.GetMessage())
	})
}
