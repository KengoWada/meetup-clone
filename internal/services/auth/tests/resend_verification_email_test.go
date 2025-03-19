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
		assert.Equal(t, "Email has been sent", response.GetMessage())
	})

	t.Run("should return an error if the request has an unknown field", func(t *testing.T) {
		testUserData := createTestUser(false)

		const unknownField = "fakeField"
		data := testutils.TestRequestData{
			"email":      testUserData.Email,
			unknownField: "random data :)",
		}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
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

	t.Run("should return an error if email is invalid", func(t *testing.T) {
		data := testutils.TestRequestData{"email": "bademail.com"}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		const responseMesaage = "Invalid request body"
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, responseMesaage, response.GetMessage())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, responseMesaage, response.GetMessage())
		assert.Equal(t, "Invalid email address provided", errorMessages["email"])

		data["email"] = ""
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, responseMesaage, response.GetMessage())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "Field is required", errorMessages["email"])
	})

	t.Run("should not alert user that email does not exist or user is already deactivated or activated", func(t *testing.T) {
		activeTestUser := createTestUser(true)
		deactivatedTestUser := createDeactivatedUser()
		email, _ := testutils.GenerateEmailAndUsername()

		invalidEmails := []string{email, activeTestUser.Email, deactivatedTestUser.Email}
		for _, invalidEmail := range invalidEmails {
			data := testutils.TestRequestData{"email": invalidEmail}
			response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusOK, response.StatusCode())
			assert.Equal(t, "Email has been sent", response.GetMessage())
		}
	})
}
