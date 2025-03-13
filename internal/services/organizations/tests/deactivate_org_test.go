package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/utils/testutils"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/assert"
)

func TestDeactivateAccount(t *testing.T) {
	testEndpoint := "/v1/organizations"
	testMethod := http.MethodPatch

	appItems, err := app.NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	appItems.App.Store = testutils.NewTestStore(t, appItems.DB)

	mux := appItems.App.Mount()
	ctx := context.Background()

	createTestUser := func(activate bool, role models.UserRole) *models.User {
		testUserData := testutils.NewTestUserData(activate)
		user, userProfile, err := testUserData.CreateTestUser(ctx, appItems.App.Store, role)
		if err != nil {
			t.Fatal(err)
		}

		user.UserProfile = userProfile
		return user
	}

	createTestOrg := func(isActive bool, role *models.Role, userID int64) *models.Organization {
		org, err := testutils.CreateTestOrganization(ctx, appItems.App.Store, isActive, role, userID)
		if err != nil {
			t.Fatal(err)
		}

		return org
	}

	generateRole := func(permission string) *models.Role {
		role := &models.Role{Name: faker.Username(options.WithGenerateUniqueValues(true))}

		switch permission {
		case "valid":
			role.Permissions = []string{internal.OrgDelete}
		case "invalid":
			role.Permissions = []string{internal.OrgDeactivate, internal.OrgUpdate}
		default:
			role.Permissions = internal.Permissions
		}

		return role
	}

	generateToken := func(ID int64, isValid bool) string {
		token, err := testutils.GenerateTesAuthToken(appItems.App.Authenticator, appItems.App.Config.AuthConfig, isValid, ID)
		if err != nil {
			t.Fatal(err)
		}

		return token
	}

	t.Run("should deactivate organization", func(t *testing.T) {
		testAdminUser := createTestUser(true, models.UserAdminRole)
		testUser := createTestUser(true, models.UserClientRole)
		org := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testAdminUser.ID, true)}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, "Done", response.GetMessage())
	})

	t.Run("should not deactivate organization not authenticated", func(t *testing.T) {
		_ = createTestUser(true, models.UserAdminRole)
		testUser := createTestUser(true, models.UserClientRole)
		org := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
		assert.Equal(t, "unauthorized", response.GetMessage())
	})

	t.Run("should not deactivate organization non staff user", func(t *testing.T) {
		_ = createTestUser(true, models.UserAdminRole)
		testUser := createTestUser(true, models.UserClientRole)
		org := createTestOrg(true, generateRole("all"), testUser.UserProfile.ID)

		endpoint := fmt.Sprintf("%s/%d", testEndpoint, org.ID)
		headers := testutils.TestRequestHeaders{"Authorization": "Bearer " + generateToken(testUser.ID, true)}

		response, err := testutils.RunTestRequest(mux, testMethod, endpoint, headers, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusForbidden, response.StatusCode())
		assert.Equal(t, "forbidden", response.GetMessage())
	})
}
