package article

import (
	"time"

	"github.com/junichi-fukushima/tech-flow/backend/domain/category"
	"github.com/junichi-fukushima/tech-flow/backend/domain/feed"
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
)

type Article struct {
	ID          *int64
	FeedID      int64
	Feed        feed.Feed
	CategoryID  int64
	Category    category.Category
	Title       string
	Link        string
	Description *string
	PubDate     time.Time
	GUID        string
	ImageUrl    *string
	Tags        []*tag.Tag `gorm:"many2many:article_tags;joinForeignKey:ArticleID;joinReferences:TagID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewArticle(article Article) *Article {
	return &article
}
