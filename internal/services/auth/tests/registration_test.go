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

func TestUserRegistration(t *testing.T) {
	testEndpoint := "/v1/auth/register"
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

	generateRequestData := func() testutils.TestRequestData {
		email, username := testutils.GenerateEmailAndUsername()
		return testutils.TestRequestData{
			"email":       email,
			"password":    testutils.TestPassword,
			"username":    username,
			"profilePic":  testutils.TestProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
		}
	}

	t.Run("should create user", func(t *testing.T) {
		data := generateRequestData()
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, response.StatusCode())
		assert.Equal(t, "Done.", response.GetMessage())
	})

	t.Run("should not create user with same email twice", func(t *testing.T) {
		testUserData := createTestUser(true)
		_, newUsername := testutils.GenerateEmailAndUsername()

		data := generateRequestData()
		data["email"] = testUserData.Email
		data["username"] = newUsername

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

		assert.Equal(t, "an account is already attached to that email address", errorMessages["email"])
	})

	t.Run("should not create user with same username twice", func(t *testing.T) {
		testUserData := createTestUser(true)
		newEmail, _ := testutils.GenerateEmailAndUsername()

		data := generateRequestData()
		data["email"] = newEmail
		data["username"] = testUserData.Username

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

		assert.Equal(t, "username is already taken", errorMessages["username"])
	})

	t.Run("should not create user invalid date of birth", func(t *testing.T) {
		data := generateRequestData()
		data["dateOfBirth"] = "21/08/1997"

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

		assert.Equal(t, "Invalid date format. mm/dd/yyyy", errorMessages["dateOfBirth"])
	})

	t.Run("should not create user invalid password", func(t *testing.T) {
		invalidPasswordError := "Password must contain a number, lower case character, upper case character and one of the special symbols(including space) !@#$%^&*()-_+=,.?|\\/<>[]{}"
		invalidPasswords := []map[string]string{
			{
				"password":     "simple",
				"errorMessage": "Password must have at least 10 characters",
			},
			{
				"password":     "60ksenB!PZcp*gYucryJmsfsky@A%jtr4$$1hLXzD@^Xcavuj6$*3iKb^YVdRFVynvGprXalw",
				"errorMessage": "Password must have at most 72 characters",
			},
			{
				"password":     "n3w_p@ssw0rd", // Missing uppercase character
				"errorMessage": invalidPasswordError,
			},
			{
				"password":     "N3W_P@SSW0RD", // Missing lowercase character
				"errorMessage": invalidPasswordError,
			},
			{
				"password":     "NeW_P@SSWoRD", // Missing number character
				"errorMessage": invalidPasswordError,
			},
			{
				"password":     "N3WPaSSW0RD", // Missing special character
				"errorMessage": invalidPasswordError,
			},
		}

		data := generateRequestData()
		for _, invalidPassword := range invalidPasswords {
			data["password"] = invalidPassword["password"]
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

			assert.Equal(t, invalidPassword["errorMessage"], errorMessages["password"])
		}
	})

	t.Run("should not create user invalid email", func(t *testing.T) {
		data := generateRequestData()
		data["email"] = "sixemail.com"

		// Invalid email
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

		assert.Equal(t, "Invalid email address provided", errorMessages["email"])

		// Email field missing
		delete(data, "email")
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid request body", response.GetMessage())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, "Field is required", errorMessages["email"])
	})

	t.Run("should not create user with invalid profile pic", func(t *testing.T) {
		data := generateRequestData()
		data["profilePic"] = "/fake/image.png"

		// Invalid email
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

		assert.Equal(t, "Invalid URL format", errorMessages["profilePic"])
	})
}
