package metaRank

import (
	"encoding/json"
	"time"
)

type RankingEvent struct {
	ID                  string
	Timestamp           time.Time
	Fields              *json.RawMessage
	UserMetadataEventID string
	Articles            json.RawMessage
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
