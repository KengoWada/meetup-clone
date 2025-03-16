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

func TestPasswordResetRequest(t *testing.T) {
	testEndpoint := "/v1/auth/password-reset-request"
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

	t.Run("should send password reset email", func(t *testing.T) {
		testUserData := createTestUser(true)

		data := testutils.TestRequestData{"email": testUserData.Email}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, "Email has been sent.", response.GetMessage())
	})

	t.Run("should not send password reset email unknown field", func(t *testing.T) {
		testUserData := createTestUser(true)
		const unknownField = "fakeField"

		data := testutils.TestRequestData{"email": testUserData.Email, unknownField: "some data"}
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

	t.Run("should not send password reset email invalid request", func(t *testing.T) {
		data := testutils.TestRequestData{"email": ""}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid request body", response.GetMessage())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "Field is required", errorMessages["email"])
	})

	t.Run("should not send password reset email email not in db", func(t *testing.T) {
		email, _ := testutils.GenerateEmailAndUsername()

		data := testutils.TestRequestData{"email": email}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, "Email has been sent.", response.GetMessage())
	})

	t.Run("should not send password reset email if user deactivated or not activated", func(t *testing.T) {
		var testUsers = []testutils.TestUserData{
			createTestUser(false),   // Inactive user
			createDeactivatedUser(), // Deactivated user
		}

		for _, testUser := range testUsers {
			data := testutils.TestRequestData{"email": testUser.Email}
			response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusOK, response.StatusCode())
			assert.Equal(t, "Email has been sent.", response.GetMessage())
		}
	})
}
