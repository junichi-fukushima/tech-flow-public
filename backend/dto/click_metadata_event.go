package dto

import (
	"encoding/json"
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"time"
)

type ClickMetadataEventRequest struct {
	Event     string           `json:"event"`
	ID        string           `json:"id"`
	Ranking   string           `json:"ranking"`
	Timestamp time.Time        `json:"timestamp"`
	User      string           `json:"user"`
	Type      string           `json:"type"`
	Item      string           `json:"item"`
	Fields    *json.RawMessage `json:"fields"`
}

func (c *ClickMetadataEventRequest) FromClickMetadataEvent(interactionEvent *metaRank.InteractionEvent, userID string, articleID string) *ClickMetadataEventRequest {
	if interactionEvent.Fields == nil {
		// nilの場合、Fieldsを空配列で初期化する
		fields := json.RawMessage([]byte("[]"))
		interactionEvent.Fields = &fields
	}
	return &ClickMetadataEventRequest{
		Event:     "interaction",
		ID:        interactionEvent.ID,
		Ranking:   *interactionEvent.RankingEventID,
		Timestamp: interactionEvent.Timestamp,
		User:      userID,
		Type:      string(interactionEvent.EventType),
		Item:      articleID,
		Fields:    interactionEvent.Fields,
	}
}
