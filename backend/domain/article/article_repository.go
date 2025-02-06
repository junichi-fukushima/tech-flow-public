package article

import (
	"github.com/junichi-fukushima/tech-flow/backend/domain/category"
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
)

// DBにつなく
type ArticleRepository interface {
	GetArticlesAll() ([]Article, error)
	GetArticlesByCategoryID(categoryID int) ([]Article, error)
	GetArticlesByTagID(tagID int) ([]Article, error)
	GetArticlesByArticleIDs(articleIDs []int) ([]Article, error)
	UpsertArticles(articles []Article) error
	GetArticlesByCategoryAndTag(categoryName, tagName string, limit, offset int, keyword string) ([]Article, error)
	GetArticles(limit, offset int) ([]Article, error)
	GetArticleIDsByGUIDs(guids []string) ([]int64, error)
	GetArticleByClickCount() ([]Article, error)
	GetArticleByFavoriteCategories(userID string) ([]Article, error)
}

// 外部API
type RestArticleRepository interface {
	// カテゴリとタグの判定
	GetTagAndCategoryByClaudeAPI(
		title string,
		description string,
		categories []category.Category,
		tags []*tag.Tag,
	) (category.Category, []*tag.Tag, error)
}
