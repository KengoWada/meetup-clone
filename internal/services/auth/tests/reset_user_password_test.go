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

	// createTestUser := func(activate bool) testutils.TestUserData {
	// 	testUserData := testutils.NewTestUserData(activate)
	// 	_, _, err := testUserData.CreateTestUser(ctx, appItems.App.Store)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	return testUserData
	// }

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

		const responseMessage = "Password successfully updated"
		assert.Equal(t, responseMessage, response.GetMessage())
	})

	t.Run("should not reset password unknown field", func(t *testing.T) {
		testUserData, token := createTestUserAndSetPasswordResetToken(true, false)

		const unknownField = "fakeField"
		data := testutils.TestRequestData{"token": token, "password": testUserData.Password, unknownField: "unknown field"}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		const responseMessage = "Unknown field in request"
		const unknownFieldMessage = "unknown field"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, unknownFieldMessage, errorMessages[unknownField])
	})

	t.Run("should not reset password invalid request body", func(t *testing.T) {
		testUserData, token := createTestUserAndSetPasswordResetToken(true, false)

		// No token sent
		data := testutils.TestRequestData{"token": "", "password": testUserData.Password}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok := response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		const responseMessage = "Invalid request body"
		var errorMessage = "Field is required"

		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["token"])

		// Password length less than 10
		data["token"] = token
		data["password"] = "shortpass"
		response, err = testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		errorMessages, ok = response.GetErrorMessages()
		if !ok {
			t.Fatal("failed to convert response errors to map")
		}

		errorMessage = "Password must have at least 10 characters"
		assert.Equal(t, responseMessage, response.GetMessage())
		assert.Equal(t, errorMessage, errorMessages["password"])

		// Password length greater than 72
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

	t.Run("should not reset password for expired token", func(t *testing.T) {
		testUserData, token := createTestUserAndSetPasswordResetToken(true, true)

		data := testutils.TestRequestData{"token": token, "password": testUserData.Password}
		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode())

		const responseMessage = "Password reset token has exipred"
		assert.Equal(t, responseMessage, response.GetMessage())
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

		const responseMessage = "Password reset token is invalid"
		assert.Equal(t, responseMessage, response.GetMessage())
	})

	t.Run("should not reset password for un verified email", func(t *testing.T) {
		testUserData, token := createTestUserAndSetPasswordResetToken(false, false)
		data := testutils.TestRequestData{"token": token, "password": testUserData.Password}

		response, err := testutils.RunTestRequest(mux, testMethod, testEndpoint, nil, data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode())

		const responseMessage = "Password reset token is invalid"
		assert.Equal(t, responseMessage, response.GetMessage())
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

		const responseMessage = "Password reset token is invalid"
		assert.Equal(t, responseMessage, response.GetMessage())
	})
}
