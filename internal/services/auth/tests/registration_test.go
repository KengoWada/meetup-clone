package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KengoWada/meetup-clone/internal/services/auth"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/stretchr/testify/assert"
)

type H map[string]any

func TestUserRegistration(t *testing.T) {
	store := utils.NewTestStore(t)
	authHandler := auth.NewHandler(store)
	mux := authHandler.RegisterRoutes()

	t.Run("should create user", func(t *testing.T) {
		data := H{
			"email":       "one@email.com",
			"password":    "C0mpl3x_P@ssw0rD",
			"username":    "username",
			"profilePic":  "https://github.com/image.png",
			"dateOfBirth": "08/21/1997",
		}
		rr, response, err := registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, response["message"], "Done.")
	})

	t.Run("should not create user with same email twice", func(t *testing.T) {
		data := H{
			"email":       "two@email.com",
			"password":    "C0mpl3x_P@ssw0rD",
			"username":    "username1",
			"profilePic":  "https://github.com/image.png",
			"dateOfBirth": "08/21/1997",
		}

		rr, response, err := registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, response["message"], "Done.")

		data["username"] = "username2"
		rr, response, err = registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		var errorMessage string = "an account is already attached to that email address"
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["error"])
	})

	t.Run("should not create user with same username twice", func(t *testing.T) {
		data := H{
			"email":       "three@email.com",
			"password":    "C0mpl3x_P@ssw0rD",
			"username":    "username3",
			"profilePic":  "https://github.com/image.png",
			"dateOfBirth": "08/21/1997",
		}

		rr, response, err := registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, response["message"], "Done.")

		data["email"] = "four@email.com"
		rr, response, err = registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		var errorMessage string = "username is already taken"
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["error"])
	})

	t.Run("should not create user invalid date of birth", func(t *testing.T) {
		data := H{
			"email":       "five@email.com",
			"password":    "C0mpl3x_P@ssw0rD",
			"username":    "username4",
			"profilePic":  "https://github.com/image.png",
			"dateOfBirth": "21/08/1997",
		}

		rr, response, err := registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		var errorMessage string = "Invalid date format. mm/dd/yyyy"
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["dateOfBirth"])
	})

	t.Run("should not create user invalid password", func(t *testing.T) {
		data := H{
			"email":       "five@email.com",
			"password":    "simple",
			"username":    "username5",
			"profilePic":  "https://github.com/image.png",
			"dateOfBirth": "21/08/1997",
		}

		// Password less than 10 characters
		rr, response, err := registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		var errorMessage string = "Password must have at least 10 characters"
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["password"])

		// Password more than 72 characters
		var veryLongPassword string
		for range 73 {
			veryLongPassword += "i"
		}

		data["password"] = veryLongPassword
		rr, response, err = registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		errorMessage = "Password must have at most 72 characters"
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["password"])

		// Password missing upper case character
		data["password"] = "n3w_p@ssw0rd"
		rr, response, err = registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		errorMessage = "Password must contain a number, lower case character, upper case character and one of the special symbols(including space) !@#$%^&*()-_+=,.?|\\/<>[]{}"
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["password"])

		// Password missing lower case character
		data["password"] = "N3W_P@SSW0RD"
		rr, response, err = registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["password"])

		// Password missing number character
		data["password"] = "NeW_P@SSWoRD"
		rr, response, err = registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["password"])

		// Password missing special character
		data["password"] = "N3WPaSSW0RD"
		rr, response, err = registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["password"])
	})

	t.Run("should not create user invalid email", func(t *testing.T) {
		data := H{
			"email":       "sixemail.com",
			"password":    "C0mpl3x_P@ssw0rD",
			"username":    "username6",
			"profilePic":  "https://github.com/image.png",
			"dateOfBirth": "21/08/1997",
		}

		// Invalid email
		rr, response, err := registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		var errorMessage string = "Invalid email address provided"
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["email"])

		// Email field missing
		delete(data, "email")
		rr, response, err = registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		errorMessage = "Field is required"
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["email"])
	})

	t.Run("should not create user with invalid profile pic", func(t *testing.T) {
		data := H{
			"email":       "sixemail.com",
			"password":    "C0mpl3x_P@ssw0rD",
			"username":    "username6",
			"profilePic":  "/github/image.png",
			"dateOfBirth": "21/08/1997",
		}

		// Invalid email
		rr, response, err := registerUserHelper(data, mux)
		if err != nil {
			t.Fatal(err)
		}

		var errorMessage string = "Invalid URL format"
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, errorMessage, response["profilePic"])
	})
}

func registerUserHelper(data H, mux http.Handler) (*httptest.ResponseRecorder, H, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	rr := utils.ExecuteRequest(req, mux)
	resp := H{}

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		return nil, nil, err
	}

	return rr, resp, nil
}
