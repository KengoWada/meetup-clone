package models

type Organization struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	ProfilePic  string `json:"profilePic"`
	IsActive    bool   `json:"isActive"`
}

type SimpleOrganization struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProfilePic  string `json:"profilePic"`
}

type Role struct {
	BaseModel
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	OrganizationID int64         `json:"organizationId"`
	Organization   *Organization `json:"organization"`
	Permissions    []string      `json:"permissions"`
}

type SimpleRole struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type OrganizationMember struct {
	BaseModel
	OrganizationID int64         `json:"organizationId"`
	Organization   *Organization `json:"organization"`
	UserProfileID  int64         `json:"userProfileId"`
	UserProfile    *UserProfile  `json:"userProfile"`
	RoleID         int64         `json:"roleId"`
	Role           *Role         `json:"role"`
}

type OrganizationInvite struct {
	BaseModel
	OrganizationID int64         `json:"organizationId"`
	Organization   *Organization `json:"organization"`
	UserProfileID  int64         `json:"userProfileId"`
	UserProfile    *UserProfile  `json:"userProfile"`
	RoleID         int64         `json:"roleId"`
	Role           *Role         `json:"role"`
	AcceptedAt     *string       `json:"acceptedAt"`
	DeclinedAt     *string       `json:"declinedAt"`
}
