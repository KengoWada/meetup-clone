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

	t.Run("should create user", func(t *testing.T) {
		email, username := testutils.GenerateEmailAndUsername()
		data := testutils.TestRequestData{
			"email":       email,
			"password":    testutils.TestPassword,
			"username":    username,
			"profilePic":  testutils.TestProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
		}

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

		data := testutils.TestRequestData{
			"email":       testUserData.Email,
			"password":    testUserData.Password,
			"username":    newUsername,
			"profilePic":  testutils.TestProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
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

		var emailErrorMessage = "an account is already attached to that email address"
		var responseMessage = "Invalid request body"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, emailErrorMessage, errorMessages["email"])
	})

	t.Run("should not create user with same username twice", func(t *testing.T) {
		testUserData := createTestUser(true)
		newEmail, _ := testutils.GenerateEmailAndUsername()

		data := testutils.TestRequestData{
			"email":       newEmail,
			"password":    testUserData.Password,
			"username":    testUserData.Username,
			"profilePic":  testutils.TestProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
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

		var usernameErrorMessage = "username is already taken"
		var responseMessage = "Invalid request body"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, usernameErrorMessage, errorMessages["username"])
	})

	t.Run("should not create user invalid date of birth", func(t *testing.T) {
		email, username := testutils.GenerateEmailAndUsername()
		data := testutils.TestRequestData{
			"email":       email,
			"password":    testutils.TestPassword,
			"username":    username,
			"profilePic":  testutils.TestProfilePic,
			"dateOfBirth": "21/08/1997",
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

		var errorMessage string = "Invalid date format. mm/dd/yyyy"
		var responseMessage = "Invalid request body"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["dateOfBirth"])
	})

	t.Run("should not create user invalid password", func(t *testing.T) {
		email, username := testutils.GenerateEmailAndUsername()
		data := testutils.TestRequestData{
			"email":       email,
			"password":    "simple",
			"username":    username,
			"profilePic":  testutils.TestProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
		}

		// Password less than 10 characters
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		var errorMessage string = "Password must have at least 10 characters"
		var responseMessage = "Invalid request body"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["password"])

		// Password more than 72 characters
		var veryLongPassword string
		for range 73 {
			veryLongPassword += "i"
		}

		data["password"] = veryLongPassword
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		errorMessage = "Password must have at most 72 characters"
		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["password"])

		// Password missing upper case character
		data["password"] = "n3w_p@ssw0rd"
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		errorMessage = "Password must contain a number, lower case character, upper case character and one of the special symbols(including space) !@#$%^&*()-_+=,.?|\\/<>[]{}"
		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["password"])

		// Password missing lower case character
		data["password"] = "N3W_P@SSW0RD"
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["password"])

		// Password missing number character
		data["password"] = "NeW_P@SSWoRD"
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["password"])

		// Password missing special character
		data["password"] = "N3WPaSSW0RD"
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["password"])
	})

	t.Run("should not create user invalid email", func(t *testing.T) {
		_, username := testutils.GenerateEmailAndUsername()
		data := testutils.TestRequestData{
			"email":       "sixemail.com",
			"password":    testutils.TestPassword,
			"username":    username,
			"profilePic":  testutils.TestProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
		}

		// Invalid email
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		var errorMessage = "Invalid email address provided"
		var responseMessage = "Invalid request body"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["email"])

		// Email field missing
		delete(data, "email")
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		errorMessage = "Field is required"
		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["email"])
	})

	t.Run("should not create user with invalid profile pic", func(t *testing.T) {
		email, username := testutils.GenerateEmailAndUsername()
		data := testutils.TestRequestData{
			"email":       email,
			"password":    testutils.TestPassword,
			"username":    username,
			"profilePic":  "/fake/image.png",
			"dateOfBirth": testutils.GenerateDate(),
		}

		// Invalid email
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		var errorMessage = "Invalid URL format"
		var responseMessage = "Invalid request body"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["profilePic"])
	})
}
