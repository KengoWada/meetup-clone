// Package models defines the data models that are mapped to the database.
// It contains structs representing the application's entities, which are used
// for data storage and retrieval from the database.
package models

// BaseModel defines common fields that can be embedded in other models.
// It includes fields for tracking the entity's ID, version, and timestamps
// for creation, updates, and deletion. These fields are meant to be reused
// across different models in the application.
type BaseModel struct {
	ID        int64   `json:"id"`
	Version   int64   `json:"-"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt *string `json:"-"`
	DeletedAt *string `json:"-"`
}
