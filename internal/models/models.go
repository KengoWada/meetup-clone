package models

type BaseModel struct {
	ID        int64   `json:"id"`
	Version   int64   `json:"-"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt *string `json:"-"`
	DeletedAt *string `json:"-"`
}
