package testutils

import (
	"context"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
)

type CreateTestUserData struct {
	Email    string
	Username string
	Password string
	Activate bool
}

func CreateTestUser(ctx context.Context, store store.Store, data CreateTestUserData) (*models.User, *models.UserProfile, error) {
	passwordHash, err := utils.GeneratePasswordHash(data.Password)
	if err != nil {
		return nil, nil, err
	}

	user := models.User{Email: data.Email, Password: passwordHash}
	userProfile := models.UserProfile{
		Username:    data.Username,
		ProfilePic:  "https://fake.img.com/user.png",
		DateOfBirth: "01/23/1999",
	}

	err = store.Users.Create(ctx, &user, &userProfile)
	if err != nil {
		return nil, nil, err
	}

	if data.Activate {
		err := store.Users.Activate(ctx, &user)
		if err != nil {
			return nil, nil, err
		}
	}

	return &user, &userProfile, nil
}
