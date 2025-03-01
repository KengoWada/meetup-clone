package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestLoginUser(t *testing.T) {
	appItems, err := app.NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	appItems.App.Store = testutils.NewTestStore(t, appItems.DB)

	mux := appItems.App.Mount()
	ctx := context.Background()

	createTestUser := func(activate bool) testutils.CreateTestUserData {
		testUserData := testutils.NewCreateTestUserData(activate)
		_, _, err := testUserData.CreateTestUser(ctx, appItems.App.Store)
		if err != nil {
			t.Fatal(err)
		}
		return testUserData
	}

	t.Run("should log in a user", func(t *testing.T) {
		testUserData := createTestUser(true)
		data := H{"email": testUserData.Email, "password": testUserData.Password}

		rr, response, err := loginUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, rr.Code)

		data, ok := response["data"].(map[string]any)
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		token, ok := data["token"]
		assert.True(t, ok)
		_, err = appItems.App.Authenticator.ValidateToken(token.(string))
		assert.Nil(t, err)
	})

	t.Run("should not log in with no credentials provided", func(t *testing.T) {
		data := H{"email": "", "password": ""}

		rr, response, err := loginUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		errorMessages, ok := response["errors"].(map[string]any)
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		var responseMessage = "Invalid request body"
		var requiredFieldMessage = "Field is required"

		assert.Equal(t, responseMessage, response["message"])
		assert.Equal(t, requiredFieldMessage, errorMessages["email"])
		assert.Equal(t, requiredFieldMessage, errorMessages["password"])
	})

	t.Run("should not log in if account doesn't exist", func(t *testing.T) {
		email, _ := testutils.GenerateEmailAndUsername()
		data := H{"email": email, "password": testutils.TestPassword}

		rr, response, err := loginUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var responseMessage = "Invalid credentials"
		assert.Equal(t, responseMessage, response["message"])
	})

	t.Run("should not log in if account is not active", func(t *testing.T) {
		testUserData := createTestUser(false)
		data := H{"email": testUserData.Email, "password": testUserData.Password}

		rr, response, err := loginUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var responseMessage = "Invalid credentials"
		assert.Equal(t, responseMessage, response["message"])
	})

	t.Run("should not log in with invalid password", func(t *testing.T) {
		testUserData := createTestUser(true)
		data := H{"email": testUserData.Email, "password": "wrong_password"}

		rr, response, err := loginUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var responseMessage = "Invalid credentials"
		assert.Equal(t, responseMessage, response["message"])
	})
}

func loginUserHelper(data H, mux http.Handler) (*httptest.ResponseRecorder, H, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	rr := testutils.ExecuteTestRequest(req, mux)
	resp := H{}

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		return nil, nil, err
	}

	return rr, resp, nil
}
