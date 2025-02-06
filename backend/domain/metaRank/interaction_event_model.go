package metaRank

import (
	"encoding/json"
	"time"
)

type InteractionEvent struct {
	ID                  string
	Timestamp           time.Time
	Fields              *json.RawMessage
	UserMetadataEventID string
	RankingEventID      *string
	ItemMetadataEventID string
	EventType           EventType
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
