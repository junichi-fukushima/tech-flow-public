package metaRank

import (
	"encoding/json"
	"github.com/junichi-fukushima/tech-flow/backend/domain/article"
	"time"
)

type ItemMetadataEvent struct {
	ID        string
	Timestamp time.Time
	Fields    *json.RawMessage
	ArticleID int64
	Article   article.Article
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewItemMetadataEvent(itemMetadataEvent ItemMetadataEvent) *ItemMetadataEvent {
	return &itemMetadataEvent
}
