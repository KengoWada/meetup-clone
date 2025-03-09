package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrganization(t *testing.T) {
	testEndpoint := "/v1/organizations"
	testMethod := http.MethodPost

	appItems, err := app.NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	appItems.App.Store = testutils.NewTestStore(t, appItems.DB)

	mux := appItems.App.Mount()
	ctx := context.Background()

	createTestUser := func(activate bool) *models.User {
		testUserData := testutils.NewTestUserData(activate)
		user, _, err := testUserData.CreateTestUser(ctx, appItems.App.Store, models.UserClientRole)
		if err != nil {
			t.Fatal(err)
		}
		return user
	}

	generateToken := func(ID int64, isValid bool) string {
		token, err := testutils.GenerateTesAuthToken(appItems.App.Authenticator, appItems.App.Config.AuthConfig, isValid, ID)
		if err != nil {
			t.Fatal(err)
		}

		return token
	}

	t.Run("should create organization", func(t *testing.T) {
		testUser := createTestUser(true)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"name":        faker.Username(options.WithGenerateUniqueValues(true)),
			"description": "Simple Description",
			"profilePic":  testutils.TestProfilePic,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, response.StatusCode())
		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, payload["name"], data["name"])
		assert.Equal(t, payload["description"], data["description"])
		assert.Equal(t, payload["profilePic"], data["profilePic"])
	})

	t.Run("should not create organization not authenticated", func(t *testing.T) {
		payload := testutils.TestRequestData{
			"name":        faker.Username(options.WithGenerateUniqueValues(true)),
			"description": "Simple Description",
			"profilePic":  testutils.TestProfilePic,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
		assert.Equal(t, "unauthorized", response.GetMessage())
	})

	t.Run("should not create same organization name", func(t *testing.T) {
		testUser := createTestUser(true)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"name":        faker.Username(options.WithGenerateUniqueValues(true)),
			"description": "Simple Description",
			"profilePic":  testutils.TestProfilePic,
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, response.StatusCode())
		data, ok := response.GetData()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		assert.Equal(t, payload["name"], data["name"])
		assert.Equal(t, payload["description"], data["description"])
		assert.Equal(t, payload["profilePic"], data["profilePic"])

		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Invalid request body", response.GetMessage())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}
		assert.Equal(t, "organization name is already taken", errorMessages["name"])
	})

	t.Run("should not create organization invalid body", func(t *testing.T) {
		testUser := createTestUser(true)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"name":        faker.Username(options.WithGenerateUniqueValues(true)),
			"description": "",
			"profilePic":  testutils.TestProfilePic,
		}

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
		assert.Equal(t, "Field is required", errorMessages["description"])
	})

	t.Run("should not create organization with unknown field", func(t *testing.T) {
		testUser := createTestUser(true)
		unknownField := "unknownField"
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}
		payload := testutils.TestRequestData{
			"name":        faker.Username(options.WithGenerateUniqueValues(true)),
			"description": "simpleDesctiption",
			"profilePic":  testutils.TestProfilePic,
			unknownField:  "fakedata",
		}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, headers, payload)
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
}
