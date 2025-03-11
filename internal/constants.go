package internal

import "time"

type userKey string
type orgKey string

const (
	DateTimeFormat = time.RFC3339
	// mm/dd/yyyy
	DateFormat = "01/02/2006"

	UserCtx userKey = "user"
	OrgCtx  orgKey  = "organization"

	// Event permissions
	EventCreate  = "create_event"
	EventPublish = "publish_event"
	EventUpdate  = "update_event"
	EventCancel  = "cancel_event"
	EventDelete  = "delete_event"

	// Member permissions
	MemberAdd        = "add_member"
	MemberRemove     = "remove_member"
	MemberRoleUpdate = "update_members_role"

	// Role permissions
	RoleCreate = "create_role"
	RoleUpdate = "update_role"
	RoleDelete = "delete_role"

	// Organization permissions
	OrgUpdate     = "update_org"
	OrgDeactivate = "deactivate_org"
	OrgDelete     = "delete_org"
)

var (
	Permissions = []string{
		EventCreate, EventPublish, EventUpdate,
		EventCancel, EventDelete,
		MemberAdd, MemberRemove, MemberRoleUpdate,
		RoleCreate, RoleUpdate, RoleDelete,
		OrgUpdate, OrgDeactivate, OrgDelete,
	}

	PermissionsMap = map[string][]string{
		"events":        {EventCreate, EventPublish, EventUpdate, EventCancel, EventDelete},
		"members":       {MemberAdd, MemberRemove, MemberRoleUpdate},
		"roles":         {RoleCreate, RoleUpdate, RoleDelete},
		"organizations": {OrgUpdate, OrgDeactivate, OrgDelete},
	}
)
