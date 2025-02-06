package auth

import (
	"encoding/json"
	"time"
)

type User struct {
	ID                    string
	SessionToken          string
	HasFavoriteCategories bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type UserMetadataEvent struct {
	ID        string
	Timestamp time.Time
	Fields    *json.RawMessage
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
