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

	t.Run("should activate user", func(t *testing.T) {
		testUserData := createTestUser(false)
		token, err := utils.GenerateToken(testUserData.Email, []byte(appItems.App.Config.SecretKey))
		if err != nil {
			t.Fatal("failed to generate test token to activate a user")
		}

		data := testutils.TestRequestData{"token": token}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		const responseMessage = "Email successfully verified"
		assert.Equal(t, responseMessage, response.GetMessage())
	})

	t.Run("should not activate if the request has an unknown field", func(t *testing.T) {
		testUserData := createTestUser(false)
		token, err := utils.GenerateToken(testUserData.Email, []byte(appItems.App.Config.SecretKey))
		if err != nil {
			t.Fatal("failed to generate test token to activate a user")
		}

		const unknownField = "fakeField"
		data := testutils.TestRequestData{"token": token, unknownField: "random data :)"}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		const responseMesaage = "Unknown field in request"
		const unknownFieldMessage = "unknown field"

		assert.Equal(t, responseMesaage, response.GetMessage())
		assert.Equal(t, unknownFieldMessage, errorMessages[unknownField])
	})

	t.Run("should not activate if token is expired", func(t *testing.T) {
		testUserData := createTestUser(false)
		token, err := utils.GenerateTestToken(
			testUserData.Email,
			[]byte(appItems.App.Config.SecretKey),
			// get current time and subtract an hour
			time.Now().Add(-time.Hour).UTC().Format(internal.DateTimeFormat),
		)
		if err != nil {
			t.Fatal("failed to generate test token to activate a user")
		}

		data := testutils.TestRequestData{"token": token}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())

		const responseMesaage = "Activation token has exipred"
		assert.Equal(t, responseMesaage, response.GetMessage())
	})

	t.Run("should not activate if email is invalid", func(t *testing.T) {
		email, _ := testutils.GenerateEmailAndUsername()
		token, err := utils.GenerateToken(email, []byte(appItems.App.Config.SecretKey))
		if err != nil {
			t.Fatal("failed to generate test token to activate a user")
		}

		data := testutils.TestRequestData{"token": token}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		const responseMesaage = "Activation token is invalid"
		assert.Equal(t, responseMesaage, response.GetMessage())
	})

	t.Run("should not activate already active user", func(t *testing.T) {
		testUserData := createTestUser(true)
		token, err := utils.GenerateToken(testUserData.Email, []byte(appItems.App.Config.SecretKey))
		if err != nil {
			t.Fatal("failed to generate test token to activate a user")
		}

		data := testutils.TestRequestData{"token": token}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		const responseMesaage = "Activation token is invalid"
		assert.Equal(t, responseMesaage, response.GetMessage())
	})

	t.Run("should not activate deactivated user", func(t *testing.T) {
		testUserData := createDeactivatedUser()
		token, err := utils.GenerateToken(testUserData.Email, []byte(appItems.App.Config.SecretKey))
		if err != nil {
			t.Fatal("failed to generate test token to activate a user")
		}

		data := testutils.TestRequestData{"token": token}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		const responseMesaage = "Activation token is invalid"
		assert.Equal(t, responseMesaage, response.GetMessage())
	})
}
