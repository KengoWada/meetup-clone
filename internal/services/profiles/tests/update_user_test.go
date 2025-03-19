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

	generateData := func() testutils.TestRequestData {
		email, username := testutils.GenerateEmailAndUsername()
		return testutils.TestRequestData{
			"email":       email,
			"username":    username,
			"profilePic":  testutils.TestProfilePic,
			"dateOfBirth": testutils.GenerateDate(),
		}
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

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := generateData()

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())

		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, payload["email"], data["email"])
		assert.Equal(t, payload["username"], data["username"])
	})

	t.Run("should not update user when not authenticated", func(t *testing.T) {
		createTestUser(true)

		payload := generateData()
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, payload)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
		assert.Equal(t, "unauthorized", response.GetMessage())
	})

	t.Run("should not update user extra fields in request", func(t *testing.T) {
		testUser := createTestUser(true)

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		unknownField := "unknownField"
		payload := generateData()
		payload[unknownField] = "fakeData"

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
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

	t.Run("should not update user with invalid fields", func(t *testing.T) {
		testUser := createTestUser(true)

		invalidTestData := []map[string]string{
			{
				"field":        "email",
				"value":        "fakeemail.com",
				"errorMessage": "Invalid email address provided",
			},
			{
				"field":        "username",
				"value":        "we",
				"errorMessage": "Username must have at least 3 characters",
			},
			{
				"field":        "profilePic",
				"value":        "/home/local/img.png",
				"errorMessage": "Invalid URL format",
			},
			{
				"field":        "dateOfBirth",
				"value":        "23/01/1432",
				"errorMessage": "Invalid date format. mm/dd/yyyy",
			},
		}

		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		for _, invalidData := range invalidTestData {
			payload := generateData()
			field := invalidData["field"]
			payload[field] = invalidData["value"]

			response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, http.StatusBadRequest, response.StatusCode())
			assert.Equal(t, "Invalid request body", response.GetMessage())

			errorMessages, ok := response.GetErrorMessages()
			if !ok {
				t.Fatal("failed to convert response errors to map")
			}

			assert.Equal(t, invalidData["errorMessage"], errorMessages[field])
		}
	})
}
