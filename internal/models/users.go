package models

// Constants representing the different user roles in the application.
const (
	UserAdminRole  UserRole = "admin"  // Role for admin users with full privileges.
	UserStaffRole  UserRole = "staff"  // Role for staff users with limited privileges.
	UserClientRole UserRole = "client" // Role for client users with basic privileges.
)

// UserRole defines the type for user roles within the application.
type UserRole string

// User represents a user in the application, including their credentials,
// status, role, and profile information. It embeds the BaseModel to include
// common fields such as ID, version, and timestamps. Sensitive fields like
// password and password reset token are omitted from the JSON response.
type User struct {
	BaseModel                       // Embeds common fields like ID, version, and timestamps.
	Email              string       `json:"email"`                 // The user's email address.
	Password           string       `json:"password"`              // The user's password (omitted from JSON).
	IsActive           bool         `json:"isActive"`              // Whether the user is active (omitted from JSON).
	ActivatedAt        *string      `json:"activatedAt"`           // Timestamp of when the user was activated (omitted from JSON).
	Role               UserRole     `json:"role"`                  // The user's role (e.g., admin, staff, client).
	PasswordResetToken string       `json:"passwordResetToken"`    // Token used for password reset (omitted from JSON).
	UserProfile        *UserProfile `json:"userProfile,omitempty"` // The user's profile.
}

// UserProfile represents a user's profile information, including their
// username, profile picture, date of birth, and associated user ID. It
// embeds the BaseModel to include common fields such as ID, version, and
// timestamps. The User field links the profile to the corresponding
// User entity.
type UserProfile struct {
	BaseModel
	Username    string `json:"username"`
	ProfilePic  string `json:"profilePic"`
	DateOfBirth string `json:"dateOfBirth"`
	UserID      int64  `json:"userId"`
	User        *User  `json:"user,omitempty"`
}

// IsDeactivated checks if the user is deactivated. It returns true if the
// user is not active (IsActive is false) and the user has an activation timestamp
// (ActivatedAt is not nil). If either condition is not met, it returns false.
func (u User) IsDeactivated() bool {
	return !u.IsActive && u.ActivatedAt != nil
}

// IsActivated checks if the user is activated. It returns true if the
// user is active (IsActive is true) and the user has an activation timestamp
// (ActivatedAt is not nil). If either condition is not met, it returns false.
func (u User) IsActivated() bool {
	return u.IsActive && u.ActivatedAt != nil
}
