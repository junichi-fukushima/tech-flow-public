// articleのバリデーション/リクエスト関連処理
package dto

import (
	"time"
)

// 定数の定義
const (
	defaultArticleLimit         int = 10
	defaultArticleOffset        int = 0
	DefaultTrendingArticleLimit int = 30
)

// ArticleRequest はリクエストパラメータの構造体
type ArticleRequest struct {
	Category string `json:"category"`
	Tag      string `json:"tag"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
	Offset   int    `json:"offset" validate:"min=0"`
	Keyword  string `json:"keyword"`
}

func (r *ArticleRequest) SetDefaults() {
	if r.Limit <= 0 {
		r.Limit = defaultArticleLimit
	}
	if r.Offset < 0 {
		r.Offset = defaultArticleOffset
	}
}

// Article represents the structure of an article
type ArticleResponse struct {
	ID          int64     `json:"id"`
	Feed        string    `json:"feed"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description *string   `json:"description"`
	PubDate     time.Time `json:"pub_date"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Meta represents the metadata for the response
type Meta struct {
	Total          int    `json:"total"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	RankingEventID string `json:"ranking_event_id"`
}

// Response represents the structure of the API response
type Response struct {
	Data []ArticleResponse `json:"data"`
	Meta Meta              `json:"meta"`
}
