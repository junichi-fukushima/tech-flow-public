package metaRank

import (
	"encoding/json"
	"time"
)

type UserMetadataEvent struct {
	ID        string
	Timestamp time.Time
	Fields    *json.RawMessage
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
