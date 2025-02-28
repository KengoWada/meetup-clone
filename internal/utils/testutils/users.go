package testutils

import (
	"context"
	"fmt"
	"strings"

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

type CreateTestUserData struct {
	Email    string
	Username string
	Password string
	Activate bool
}

func NewCreateTestUserData(activate bool) CreateTestUserData {
	email, username := GenerateEmailAndUsername()

	return CreateTestUserData{
		Email:    email,
		Username: username,
		Password: TestPassword,
		Activate: activate,
	}
}

func (c CreateTestUserData) CreateTestUser(ctx context.Context, store store.Store) (*models.User, *models.UserProfile, error) {
	passwordHash, err := utils.GeneratePasswordHash(c.Password)
	if err != nil {
		return nil, nil, err
	}

	user := models.User{Email: c.Email, Password: passwordHash}
	userProfile := models.UserProfile{
		Username:    c.Username,
		ProfilePic:  TestProfilePic,
		DateOfBirth: GenerateDate(),
	}

	err = store.Users.Create(ctx, &user, &userProfile)
	if err != nil {
		return nil, nil, err
	}

	if c.Activate {
		err := store.Users.Activate(ctx, &user)
		if err != nil {
			return nil, nil, err
		}
	}

	return &user, &userProfile, nil
}

func GenerateDate() string {
	date := strings.Split(faker.Date(), "-")
	return fmt.Sprintf("%s/%s/%s", date[1], date[2], date[0])
}

func GenerateEmailAndUsername() (email string, username string) {
	name := faker.Username(options.WithGenerateUniqueValues(true))

	return name + "@clone.meetup.org", name
}
