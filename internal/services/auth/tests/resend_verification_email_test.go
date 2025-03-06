package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestResendVerificationEmail(t *testing.T) {
	testEndpoint := "/v1/auth/resend-verification-email"
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

	t.Run("should resend verification email", func(t *testing.T) {
		testUserData := createTestUser(false)
		data := testutils.TestRequestData{"email": testUserData.Email}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		responseMessage := "Email has been sent"
		assert.Equal(t, responseMessage, response.GetMessage())
	})

	t.Run("should return an error if the request has an unknown field", func(t *testing.T) {
		testUserData := createTestUser(false)

		const unknownField = "fakeField"
		data := testutils.TestRequestData{"email": testUserData.Email, unknownField: "random data :)"}
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

	t.Run("should return an error if email is invalid", func(t *testing.T) {
		data := testutils.TestRequestData{"email": "bademail.com"}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		const responseMesaage = "Invalid request body"
		var emailErrorMesaage = "Invalid email address provided"

		assert.Equal(t, responseMesaage, response.GetMessage())
		assert.Equal(t, emailErrorMesaage, errorMessages["email"])

		data["email"] = ""
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		emailErrorMesaage = "Field is required"

		assert.Equal(t, responseMesaage, response.GetMessage())
		assert.Equal(t, emailErrorMesaage, errorMessages["email"])
	})

	t.Run("should not alert user that email does not exist or user is already deactivated or activated", func(t *testing.T) {
		// Email has no account
		email, _ := testutils.GenerateEmailAndUsername()
		data := testutils.TestRequestData{"email": email}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		const responseMessage = "Email has been sent"
		assert.Equal(t, responseMessage, response.GetMessage())

		// Activated user
		testUserData := createTestUser(true)
		data["email"] = testUserData.Email

		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, responseMessage, response.GetMessage())

		// Deactivated user
		testUserData = createDeactivatedUser()
		data["email"] = testUserData.Email

		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, responseMessage, response.GetMessage())
	})
}
