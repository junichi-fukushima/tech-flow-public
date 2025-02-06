package dto

import (
	"encoding/json"
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"time"
)

type ImpMetadataEventRequest struct {
	Event     string          `json:"event"`
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	User      string          `json:"user"`
	Items     json.RawMessage `json:"items"`
}

func (c *ImpMetadataEventRequest) FromImpMetadataEvent(rankingEvent *metaRank.RankingEvent, userID string) *ImpMetadataEventRequest {
	if rankingEvent.Fields == nil {
		// nilの場合、Fieldsを空配列で初期化する
		fields := json.RawMessage([]byte("[]"))
		rankingEvent.Fields = &fields
	}
	return &ImpMetadataEventRequest{
		Event:     "ranking",
		ID:        rankingEvent.ID,
		Timestamp: rankingEvent.Timestamp,
		User:      userID,
		Items:     rankingEvent.Articles,
	}
}
