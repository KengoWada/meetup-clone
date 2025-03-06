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

func TestUpdateUserProfile(t *testing.T) {
	testEndpoint := "/v1/profiles"
	testMethod := http.MethodPut

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

	t.Run("should update user details", func(t *testing.T) {
		testUser := createTestUser(true)
		token := generateToken(testUser.ID, true)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + token}

		email, username := testutils.GenerateEmailAndUsername()
		payload := testutils.TestRequestData{
			"email":       email,
			"username":    username,
			"profilePic":  testUser.UserProfile.ProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		respEmail, ok := data["email"]
		assert.True(t, ok)
		assert.Equal(t, email, respEmail)

		respUsername, ok := data["username"]
		assert.True(t, ok)
		assert.Equal(t, username, respUsername)
	})

	t.Run("should not update user with no auth token", func(t *testing.T) {
		testUser := createTestUser(true)

		email, username := testutils.GenerateEmailAndUsername()
		payload := testutils.TestRequestData{
			"email":       email,
			"username":    username,
			"profilePic":  testUser.UserProfile.ProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
		assert.Equal(t, "unauthorized", response.GetMessage())
	})

	t.Run("should not update user extra fields in request", func(t *testing.T) {
		testUser := createTestUser(true)
		token := generateToken(testUser.ID, true)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + token}

		unknownField := "unknownField"
		email, username := testutils.GenerateEmailAndUsername()
		payload := testutils.TestRequestData{
			"email":       email,
			"username":    username,
			"profilePic":  testUser.UserProfile.ProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
			unknownField:  "fakeData",
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		const responseMessage = "Unknown field in request"
		const unknownFieldErr = "unknown field"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, unknownFieldErr, errorMessages[unknownField])
	})

	t.Run("should not update user with invalid fields", func(t *testing.T) {
		testUser := createTestUser(true)
		token := generateToken(testUser.ID, true)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + token}

		payload := testutils.TestRequestData{
			"email":       "fakeemail.com",
			"username":    testUser.UserProfile.Username,
			"profilePic":  testUser.UserProfile.ProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		const responseMessage = "Invalid request body"
		var errorMessage = "Invalid email address provided"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["email"])

		// Invalid username
		payload = testutils.TestRequestData{
			"email":       testUser.Email,
			"username":    "we",
			"profilePic":  testUser.UserProfile.ProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
		}

		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		errorMessage = "Username must have at least 3 characters"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["username"])

		// Invalid profile pic url
		payload = testutils.TestRequestData{
			"email":       testUser.Email,
			"username":    testUser.UserProfile.Username,
			"profilePic":  "/home/local/img.png",
			"dateOfBirth": testutils.GenerateDate(),
		}

		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		errorMessage = "Invalid URL format"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["profilePic"])

		// Invalid date fornat
		payload = testutils.TestRequestData{
			"email":       testUser.Email,
			"username":    testUser.UserProfile.Username,
			"profilePic":  testUser.UserProfile.ProfilePic,
			"dateOfBirth": "23/01/1432",
		}

		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		errorMessage = "Invalid date format. mm/dd/yyyy"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["dateOfBirth"])
	})
}
