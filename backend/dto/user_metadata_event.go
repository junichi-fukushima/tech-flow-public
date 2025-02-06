package dto

import (
	"encoding/json"
	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
	"time"
)

type UserMetadataEventRequest struct {
	Event     string           `json:"event"`
	ID        string           `json:"id"`
	Timestamp time.Time        `json:"timestamp"`
	User      string           `json:"user"`
	Fields    *json.RawMessage `json:"fields"`
}

func (u *UserMetadataEventRequest) FromUserMetadataEvent(userMetadataEvent *auth.UserMetadataEvent) *UserMetadataEventRequest {
	if userMetadataEvent.Fields == nil {
		// nilの場合、Fieldsを空配列で初期化する
		fields := json.RawMessage([]byte("[]"))
		userMetadataEvent.Fields = &fields
	}
	return &UserMetadataEventRequest{
		Event:     "user",
		ID:        userMetadataEvent.ID,
		Timestamp: userMetadataEvent.Timestamp,
		User:      userMetadataEvent.UserID,
		Fields:    userMetadataEvent.Fields,
	}
}
