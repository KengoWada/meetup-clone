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

func TestResetUserPassword(t *testing.T) {
	testEndpoint := "/v1/auth/reset-password"
	testMethod := http.MethodPost

	appItems, err := app.NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	appItems.App.Store = testutils.NewTestStore(t, appItems.DB)

	mux := appItems.App.Mount()
	ctx := context.Background()

	createTestUserAndSetPasswordResetToken := func(activate bool, expiredToken bool) (testutils.TestUserData, string) {
		testUserData := testutils.NewTestUserData(activate)
		user, _, err := testUserData.CreateTestUser(ctx, appItems.App.Store, models.UserClientRole)
		if err != nil {
			t.Fatal(err)
		}

		timeNow := time.Now().UTC()
		if expiredToken {
			timeNow = timeNow.Add(-time.Hour)
		}

		token, err := utils.GenerateTestToken(
			testUserData.Email,
			[]byte(appItems.App.Config.SecretKey),
			timeNow.Format(internal.DateTimeFormat),
		)
		if err != nil {
			t.Fatal(err)
		}

		user.PasswordResetToken = token
		err = appItems.App.Store.Users.SetPasswordResetToken(ctx, user)
		if err != nil {
			t.Fatal(err)
		}

		return testUserData, token
	}

	t.Run("should reset the users password", func(t *testing.T) {
		testUserData, token := createTestUserAndSetPasswordResetToken(true, false)

		data := testutils.TestRequestData{"token": token, "password": testUserData.Password}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, "Password successfully updated", response.GetMessage())
	})

	t.Run("should not reset password unknown field", func(t *testing.T) {
		testUserData, token := createTestUserAndSetPasswordResetToken(true, false)

		const unknownField = "fakeField"
		data := testutils.TestRequestData{
			"token":      token,
			"password":   testUserData.Password,
			unknownField: "unknown field",
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

	t.Run("should not reset password invalid request body", func(t *testing.T) {
		testUserData, token := createTestUserAndSetPasswordResetToken(true, false)

		invalidPasswordError := "Password must contain a number, lower case character, upper case character and one of the special symbols(including space) !@#$%^&*()-_+=,.?|\\/<>[]{}"
		testInvalidData := []map[string]string{
			{
				"token":        "",
				"password":     testUserData.Password,
				"invalidField": "token",
				"errorMessage": "Field is required",
			},
			{
				"token":        token,
				"password":     "shortpass",
				"invalidField": "password",
				"errorMessage": "Password must have at least 10 characters",
			},
			{
				"token":        token,
				"password":     "60ksenB!PZcp*gYucryJmsfsky@A%jtr4$$1hLXzD@^Xcavuj6$*3iKb^YVdRFVynvGprXalw",
				"invalidField": "password",
				"errorMessage": "Password must have at most 72 characters",
			},
			{
				"token":        token,
				"password":     "n3w_p@ssw0rd", // missing uppercase character
				"invalidField": "password",
				"errorMessage": invalidPasswordError,
			},
			{
				"token":        token,
				"password":     "N3W_P@SSW0RD", // missing lowercase character
				"invalidField": "password",
				"errorMessage": invalidPasswordError,
			},
			{
				"token":        token,
				"password":     "NeW_P@SSWoRD", // missing number character
				"invalidField": "password",
				"errorMessage": invalidPasswordError,
			},
			{
				"token":        token,
				"password":     "N3WPaSSW0RD", // missing special character
				"invalidField": "password",
				"errorMessage": invalidPasswordError,
			},
		}

		for _, invalidData := range testInvalidData {
			data := testutils.TestRequestData{
				"token":    invalidData["token"],
				"password": invalidData["password"],
			}

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

			field := invalidData["invalidField"]
			assert.Equal(t, invalidData["errorMessage"], errorMessages[field])
		}
	})

	t.Run("should not reset password for expired token", func(t *testing.T) {
		testUserData, token := createTestUserAndSetPasswordResetToken(true, true)

		data := testutils.TestRequestData{"token": token, "password": testUserData.Password}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())
		assert.Equal(t, "Password reset token has exipred", response.GetMessage())
	})

	t.Run("should not reset password for email not in db", func(t *testing.T) {
		email, _ := testutils.GenerateEmailAndUsername()
		token, err := utils.GenerateToken(email, []byte(appItems.App.Config.SecretKey))
		if err != nil {
			t.Fatal(err)
		}

		data := testutils.TestRequestData{"token": token, "password": testutils.TestPassword}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Password reset token is invalid", response.GetMessage())
	})

	t.Run("should not reset password for un verified email", func(t *testing.T) {
		testUserData, token := createTestUserAndSetPasswordResetToken(false, false)

		data := testutils.TestRequestData{"token": token, "password": testUserData.Password}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Password reset token is invalid", response.GetMessage())
	})

	t.Run("should not reset password if the token is not the same as password reset token", func(t *testing.T) {
		testUserData, _ := createTestUserAndSetPasswordResetToken(true, false)
		token, err := utils.GenerateToken(testUserData.Email, []byte(appItems.App.Config.SecretKey))
		if err != nil {
			t.Fatal(err)
		}

		data := testutils.TestRequestData{"token": token, "password": testUserData.Password}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		assert.Equal(t, "Password reset token is invalid", response.GetMessage())
	})
}
