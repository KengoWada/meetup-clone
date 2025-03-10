package testutils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

var (
	TestPassword   = "C0mpl3x_P@ssw0rD"
	TestProfilePic = "https://fake.link.org/pfp.png"
)

// TestUserData holds the data required to create a test user. It is typically used
// for setting up mock or test user data in unit tests or integration tests.
//
// Fields:
//   - Email: The email address of the test user. This field is used to create a user in
//     the system, and should be unique for each test case.
//   - Username: The username for the test user. It can be any string value used in the
//     user creation process.
//   - Password: The password for the test user. This should be a valid password string
//     used for testing login and authentication logic.
//   - Activate: A boolean flag indicating whether the user should be activated upon creation.
//     If true, the user will be marked as activated; otherwise, they will remain inactive.
//
// Example usage:
//
//	testData := TestUserData{
//	  Email:    "testuser@example.com",
//	  Username: "testuser",
//	  Password: "securePassword123",
//	  Activate: true,
//	}
type TestUserData struct {
	Email    string
	Username string
	Password string
	Activate bool
}

// NewTestUserData generates a new instance of TestUserData with a random email,
// username, a predefined test password, and an activation status based on the `activate` argument.
//
// This function utilizes the `GenerateEmailAndUsername` function to generate a unique email
// and username, and sets a default test password (`TestPassword`) while allowing the caller
// to specify whether the user should be activated or not.
//
// Parameters:
//   - activate (bool): A flag indicating whether the user should be activated or not.
//
// Returns:
//   - A `TestUserData` struct populated with the generated email, username, password,
//     and activation status.
//
// Example usage:
//
//	userData := NewTestUserData(true)
//	fmt.Println(userData.Email, userData.Username, userData.Activate)
func NewTestUserData(activate bool) TestUserData {
	email, username := GenerateEmailAndUsername()

	return TestUserData{
		Email:    email,
		Username: username,
		Password: TestPassword,
		Activate: activate,
	}
}

// CreateTestUser creates a test user and their associated profile in the database using
// the provided store.
// The method interacts with the `store` to save the user and user profile data,
// returning the created `User` and `UserProfile` models along with any errors encountered
// during the process.
//
// Parameters:
//   - ctx (context.Context): The context to associate with the database operation.
//   - store (store.Store): The store instance used to interact with the database to create the user and profile.
//
// Returns:
//   - *models.User: The created `User` model.
//   - *models.UserProfile: The created `UserProfile` model.
//   - error: Any error that occurred during the creation process, or nil if successful.
//
// Example usage:
//
//	userData := NewTestUserData(true)
//	user, profile, err := userData.CreateTestUser(ctx, store)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("User created:", user, "Profile created:", profile)
func (c TestUserData) CreateTestUser(ctx context.Context, appStore store.Store, role models.UserRole) (*models.User, *models.UserProfile, error) {
	passwordHash, err := utils.GeneratePasswordHash(c.Password)
	if err != nil {
		return nil, nil, err
	}

	user := models.User{Email: c.Email, Password: passwordHash, Role: role}
	userProfile := models.UserProfile{
		Username:    c.Username,
		ProfilePic:  TestProfilePic,
		DateOfBirth: GenerateDate(),
	}

	err = appStore.Users.Create(ctx, &user, &userProfile)
	if err != nil {
		return nil, nil, err
	}

	if c.Activate {
		timeNow := time.Now().UTC().Format(internal.DateTimeFormat)
		user.ActivatedAt = &timeNow

		err := appStore.Users.Activate(ctx, &user)
		if err != nil {
			return nil, nil, err
		}
	}

	return &user, &userProfile, nil
}

// CreateDeactivatedTestUser creates a new user with a deactivated state for testing purposes.
// The user is created with the specified details, and the `IsActive` field is set to `false`.
// The `ActivatedAt` field is set to `time.Now()` to simulate a deactivated user.
//
// Parameters:
//   - ctx (context.Context): The context for managing cancellation and deadlines.
//   - store (store.Store): The store object that handles the interaction with the database.
//
// Returns:
//   - *models.User: The created user object with the deactivated state.
//   - error: An error if the user creation fails, or nil if the creation is successful.
func (c TestUserData) CreateDeactivatedTestUser(ctx context.Context, store store.Store, role models.UserRole) (*models.User, error) {
	user, _, err := c.CreateTestUser(ctx, store, role)
	if err != nil {
		return nil, err
	}

	err = store.Users.Deactivate(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GenerateDate generates a random date in the format "mm/dd/yyyy" using the faker package.
// The function splits the generated date string (in the format "yyyy-mm-dd") and reformats
// it to "mm/dd/yyyy" to comply with the desired format.
//
// Returns:
//   - A string representing the randomly generated date in "mm/dd/yyyy" format.
//
// Example usage:
//
//	date := GenerateDate()
//	fmt.Println("Generated date:", date) // Output: "03/15/2025"
func GenerateDate() string {
	date := strings.Split(faker.Date(), "-")
	return fmt.Sprintf("%s/%s/%s", date[1], date[2], date[0])
}

// GenerateEmailAndUsername generates a unique email and username using the faker package.
// The function uses the faker library to create a random, unique username and then formats
// it into an email address.
//
// Returns:
//   - email: A randomly generated email address formatted as "<username>@clone.meetup.org".
//   - username: A randomly generated username. This is used as the local part of the email address.
//
// Example usage:
//
//	email, username := GenerateEmailAndUsername()
//	fmt.Println("Generated email:", email)
//	fmt.Println("Generated username:", username)
func GenerateEmailAndUsername() (email string, username string) {
	name := faker.Username(options.WithGenerateUniqueValues(true))

	return name + "@clone.meetup.org", name
}
