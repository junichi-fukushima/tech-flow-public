package usecase

import (
	"github.com/junichi-fukushima/tech-flow/backend/domain/article"
	"github.com/junichi-fukushima/tech-flow/backend/domain/category"
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
)

type ArticleUsecase interface {
	GetArticles(
		category string,
		tag string,
		limitParam int,
		offsetParam int,
		keyword string,
	) ([]article.Article, error)
	GetCategoryAndTag(
		articleTitle string,
		articleDescription *string,
		tags []*tag.Tag,
		categories []category.Category,
	) ([]*tag.Tag, category.Category, error)
	GetArticlesByArticleIDs(articleIDs []int) ([]article.Article, error)
	GetArticlesByClickCount() ([]article.Article, error)
	GetArticlesByFavoriteCategories(userID string) ([]article.Article, error)
}

type articleUsecase struct {
	articleRepository     article.ArticleRepository
	restArticleRepository article.RestArticleRepository
}

func NewArticleUsecase(
	articleRepository article.ArticleRepository,
	restArticleRepository article.RestArticleRepository,
) ArticleUsecase {
	return &articleUsecase{
		articleRepository:     articleRepository,
		restArticleRepository: restArticleRepository,
	}
}

// 記事一括取得
func (u *articleUsecase) GetArticles(
	category string,
	tag string,
	limitParam int,
	offsetParam int,
	keyword string,
) ([]article.Article, error) {
	return u.articleRepository.GetArticlesByCategoryAndTag(
		category,
		tag,
		int(limitParam),
		int(offsetParam),
		keyword,
	)
}

// 記事のカテゴリとタグを記事情報から取得する
func (u *articleUsecase) GetCategoryAndTag(
	articleTitle string,
	articleDescription *string,
	tags []*tag.Tag,
	categories []category.Category,
) ([]*tag.Tag, category.Category, error) {
	description := ""
	if articleDescription != nil {
		description = *articleDescription
	}

	matchedCategory, matchedTags, err := u.restArticleRepository.GetTagAndCategoryByClaudeAPI(
		articleTitle,
		description,
		categories,
		tags,
	)

	if err != nil {
		return matchedTags, matchedCategory, err
	}

	return matchedTags, matchedCategory, nil
}

func (u *articleUsecase) GetArticlesByArticleIDs(articleIDs []int) ([]article.Article, error) {
	return u.articleRepository.GetArticlesByArticleIDs(articleIDs)
}

func (u *articleUsecase) GetArticlesByClickCount() ([]article.Article, error) {
	return u.articleRepository.GetArticleByClickCount()
}

func (u *articleUsecase) GetArticlesByFavoriteCategories(userID string) ([]article.Article, error) {
	return u.articleRepository.GetArticleByFavoriteCategories(userID)
}
