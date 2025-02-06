package dto

import (
	"encoding/json"
	"time"
)

type RecommendationMetadataEventRequest struct {
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	User      string          `json:"user"`
	Items     json.RawMessage `json:"items"`
}

func (c *RecommendationMetadataEventRequest) FromRecommendationMetadataEvent(uuID string, userID string, items []byte) *RecommendationMetadataEventRequest {
	return &RecommendationMetadataEventRequest{
		ID:        uuID,
		Timestamp: time.Now(),
		User:      userID,
		Items:     items,
	}
}
