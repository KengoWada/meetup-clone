package models

type UserRole string

const (
	UserAdminRole  UserRole = "admin"
	UserStaffRole  UserRole = "staff"
	UserClientRole UserRole = "client"
)

type User struct {
	BaseModel
	Email              string       `json:"email"`
	Password           string       `json:"-"`
	IsActive           bool         `json:"-"`
	Role               UserRole     `json:"role"`
	PasswordResetToken string       `json:"-"`
	UserProfile        *UserProfile `json:"userProfile"`
}

type UserProfile struct {
	BaseModel
	Username    string `json:"username"`
	ProfilePic  string `json:"profilePic"`
	DateOfBirth string `json:"dateOfBirth"`
	UserID      int64  `json:"userId"`
	User        *User  `json:"user"`
}
