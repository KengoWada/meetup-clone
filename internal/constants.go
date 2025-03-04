package internal

import "time"

type userKey string

const (
	DateTimeFormat = time.RFC3339
	// mm/dd/yyyy
	DateFormat = "01/02/2006"

	UserCtx userKey = "user"
)
