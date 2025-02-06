package repository

import (
	"fmt"
	"time"

	"github.com/junichi-fukushima/tech-flow/backend/domain/article"
	"github.com/junichi-fukushima/tech-flow/backend/domain/category"
	"github.com/junichi-fukushima/tech-flow/backend/domain/feed"
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/gorm"
	gormclause "gorm.io/gorm/clause"
)

type Article struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	FeedID      int64     `gorm:"column:feed_id;not null;index"`
	Feed        Feed      `gorm:"foreignKey:FeedID;references:ID"`
	CategoryID  int       `gorm:"column:category_id;not null"`
	Category    Category  `gorm:"foreignKey:CategoryID;references:ID"`
	Title       string    `gorm:"column:title;type:varchar(255);not null"`
	Link        string    `gorm:"column:link;type:text;not null"`
	Description *string   `gorm:"column:description;type:text"`
	PubDate     time.Time `gorm:"column:pub_date;type:datetime;not null"`
	GUID        string    `gorm:"column:guid;type:varchar(255);not null;unique"`
	ImageUrl    *string   `gorm:"column:image_url;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Tags        []Tag     `gorm:"many2many:article_tags;joinForeignKey:ArticleID;joinReferences:TagID"`
}

type ArticleTag struct {
	ArticleID int64 `gorm:"column:article_id;primaryKey"`
	TagID     int64 `gorm:"column:tag_id;primaryKey"`
}

type GormArticleRepository struct {
	db           *gorm.SQLHandler
	feedRepo     feed.FeedRepository
	categoryRepo category.CategoryRepository
	tagRepo      tag.TagsRepository
}

func NewArticleRepository() (article.ArticleRepository, error) {
	sqlHandler, err := gorm.GetSQLHandler()
	if err != nil {
		return nil, err
	}
	// repository
	feedRepo, _ := NewFeedRepository()
	categoryRepo, _ := NewCategoryRepository()
	tagRepo, _ := NewTagsRepository()
	return &GormArticleRepository{
		db:           sqlHandler,
		feedRepo:     feedRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}, nil
}

func (repo *GormArticleRepository) GetArticlesAll() ([]article.Article, error) {
	return repo.getArticlesByCondition("1 = 1")
}

func (repo *GormArticleRepository) GetArticles(limit, offset int) ([]article.Article, error) {
	return repo.getArticlesWithLimitOffset("1 = 1", limit, offset)
}

func (repo *GormArticleRepository) GetArticlesByCategoryID(categoryID int) ([]article.Article, error) {
	return repo.getArticlesByCondition("category_id = ?", categoryID)
}

// カテゴリ名とタグ名を指定して記事を取得
func (repo *GormArticleRepository) GetArticlesByCategoryAndTag(categoryName, tagName string, limit, offset int, keyword string) ([]article.Article, error) {
	if keyword != "" {
		return repo.searchArticles(limit, offset, keyword)
	}
	if categoryName != "" && tagName != "" {
		return repo.getArticlesByCategoryAndTagName(categoryName, tagName, limit, offset)
	}
	if categoryName != "" {
		return repo.getArticlesByCategoryName(categoryName, limit, offset)
	}
	if tagName != "" {
		return repo.getArticlesByTagName(tagName, limit, offset)
	}
	return repo.GetArticles(limit, offset)
}

func (repo *GormArticleRepository) searchArticles(limit, offset int, keyword string, args ...interface{}) ([]article.Article, error) {
	var articles []Article
	err := repo.db.DB.
		Where("articles.title LIKE ? OR articles.description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("articles.created_at DESC"). // created_atの降順でソート
		Limit(limit).
		Offset(offset).
		Preload("Feed").
		Preload("Category").
		Preload("Tags").
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	var domainArticles []article.Article
	for _, article := range articles {
		domainArticles = append(domainArticles, article.toDomain())
	}
	return domainArticles, nil
}

func (repo *GormArticleRepository) getArticlesByCategoryName(categoryName string, limit, offset int, args ...interface{}) ([]article.Article, error) {
	var articles []Article
	err := repo.db.DB.
		Joins("JOIN categories ON articles.category_id = categories.id").
		Where("categories.name = ?", categoryName). // カテゴリ名で絞り込み
		Order("articles.pub_date DESC").            // pub_date
		Limit(limit).
		Offset(offset).
		Preload("Feed").
		Preload("Category").
		Preload("Tags").
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	var domainArticles []article.Article
	for _, article := range articles {
		domainArticles = append(domainArticles, article.toDomain())
	}
	return domainArticles, nil
}

func (repo *GormArticleRepository) getArticlesByTagName(tagName string, limit, offset int, args ...interface{}) ([]article.Article, error) {
	var articles []Article
	err := repo.db.DB.
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Joins("JOIN tags ON tags.id = article_tags.tag_id").
		Order("articles.pub_date DESC"). // pub_dateの降順でソート
		Where("tags.name = ?", tagName). // タグ名で絞り込み
		Limit(limit).
		Offset(offset).
		Preload("Feed").
		Preload("Category").
		Preload("Tags").
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	var domainArticles []article.Article
	for _, article := range articles {
		domainArticles = append(domainArticles, article.toDomain())
	}
	return domainArticles, nil
}

func (repo *GormArticleRepository) getArticlesByCategoryAndTagName(categoryName, tagName string, limit, offset int, args ...interface{}) ([]article.Article, error) {
	var articles []Article
	err := repo.db.DB.
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Joins("JOIN tags ON tags.id = article_tags.tag_id").
		Where("categories.name = ?", categoryName). // カテゴリ名で絞り込み
		Where("tags.name = ?", tagName).            // タグ名で絞り込み
		Order("articles.pub_date DESC").            // pub_dateの降順でソート
		Limit(limit).
		Offset(offset).
		Preload("Feed").
		Preload("Category").
		Preload("Tags").
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	var domainArticles []article.Article
	for _, article := range articles {
		domainArticles = append(domainArticles, article.toDomain())
	}

	return domainArticles, nil
}

func (repo *GormArticleRepository) GetArticlesByTagID(tagID int) ([]article.Article, error) {
	var articles []Article
	err := repo.db.DB.
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id = ?", tagID).
		Preload("Feed").
		Preload("Category").
		Preload("Tags").
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	var domainArticles []article.Article
	for _, article := range articles {
		domainArticles = append(domainArticles, article.toDomain())
	}

	return domainArticles, nil
}

func (repo *GormArticleRepository) GetArticlesByArticleIDs(articleIDs []int) ([]article.Article, error) {
	var articles []Article
	err := repo.db.DB.
		Preload("Feed").
		Where("id in ?", articleIDs).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	var domainArticles []article.Article
	for _, article := range articles {
		domainArticles = append(domainArticles, article.toDomain())
	}

	return domainArticles, nil
}

// limit, offsetを指定して記事を取得
func (repo *GormArticleRepository) getArticlesWithLimitOffset(condition string, limit, offset int, args ...interface{}) ([]article.Article, error) {
	var articles []Article
	err := repo.db.DB.
		Preload("Feed").
		Preload("Category"). // カテゴリをリレーションで取得
		Preload("Tags").     // タグをリレーションで取得
		Where(condition, args...).
		Order("articles.pub_date DESC"). // pub_dateの降順でソート
		Limit(limit).
		Offset(offset).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	// ドメインモデルに変換
	var domainArticles []article.Article
	for _, article := range articles {
		domainArticles = append(domainArticles, article.toDomain())
	}

	return domainArticles, nil
}

// 共通クエリロジック
func (repo *GormArticleRepository) getArticlesByCondition(condition string, args ...interface{}) ([]article.Article, error) {
	var articles []Article
	err := repo.db.DB.
		Preload("Feed").
		Preload("Category"). // カテゴリをリレーションで取得
		Preload("Tags").     // タグをリレーションで取得
		Where(condition, args...).Find(&articles).Error
	if err != nil {
		return nil, err
	}

	// ドメインモデルに変換
	var domainArticles []article.Article
	for _, article := range articles {
		domainArticles = append(domainArticles, article.toDomain())
	}

	return domainArticles, nil
}

func (repo *GormArticleRepository) UpsertArticles(articles []article.Article) error {
	if len(articles) == 0 {
		return fmt.Errorf("articles slice is empty")
	}

	var dbArticles []Article
	for _, domainArticle := range articles {
		articleEntity := (&Article{}).fromDomain(domainArticle)
		dbArticles = append(dbArticles, articleEntity)
	}

	err := repo.db.DB.Debug().Clauses(
		gormclause.OnConflict{
			Columns:   []gormclause.Column{{Name: "guid"}},
			DoUpdates: gormclause.AssignmentColumns([]string{"feed_id", "category_id", "title", "link", "description", "pub_date", "image_url", "updated_at"}),
		},
	).Save(&dbArticles).Error

	if err != nil {
		return fmt.Errorf("failed to upsert articles: %w", err)
	}

	return nil
}

func (repo *GormArticleRepository) GetArticleIDsByGUIDs(guids []string) ([]int64, error) {
	if len(guids) == 0 {
		return nil, fmt.Errorf("guids slice is empty")
	}

	var articleIDs []int64
	err := repo.db.DB.
		Model(&Article{}).
		Where("guid IN ?", guids).
		Pluck("id", &articleIDs).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to get article IDs: %w", err)
	}

	return articleIDs, nil
}

func (repo *GormArticleRepository) GetArticleByClickCount() ([]article.Article, error) {
	sql := `
SELECT
    a.*
FROM
    articles a
        JOIN
    item_metadata_events ime ON a.id = ime.article_id
        JOIN
    interaction_events ie ON ime.id = ie.item_metadata_event_id
WHERE
    ie.event_type = 'Click'
GROUP BY
    a.id
ORDER BY
    COUNT(*) DESC
LIMIT 20
`

	var articles []Article
	err := repo.db.DB.Raw(sql).Scan(&articles).Error
	if err != nil {
		return nil, err
	}

	var domainArticles []article.Article
	for _, a := range articles {
		domainArticles = append(domainArticles, a.toDomain())
	}
	return domainArticles, nil
}

// お気に入りカテゴリに関連する記事を取得
func (repo *GormArticleRepository) GetArticleByFavoriteCategories(userID string) ([]article.Article, error) {

	// user_categoriesから特定のuser_idに関連するcategory_idを取得
	queryUserCategories := repo.db.DB.Table("user_categories").
		Select("category_id").
		Where("user_id = ?", userID)

	// カテゴリIDの数を確認
	var categoryCount int64
	queryUserCategories.Count(&categoryCount)

	var articles []Article
	query := repo.db.DB.
		Preload("Feed").
		Preload("Category").
		Preload("Tags")

	// category_idが存在する場合のみ絞り込み
	if categoryCount > 0 {
		query = query.Where("category_id IN (?)", queryUserCategories)
	}

	err := query.
		Order("articles.pub_date DESC").
		Limit(300).
		Find(&articles).Error

	if err != nil {
		return nil, err
	}

	var domainArticles []article.Article
	for _, a := range articles {
		domainArticles = append(domainArticles, a.toDomain())
	}
	return domainArticles, nil
}

// アプリケーション側で使用するドメインモデルを返却する
func (a *Article) toDomain() article.Article {

	feed := feed.Feed{
		ID:    a.Feed.ID,
		Title: a.Feed.Title,
	}
	category := category.Category{
		ID:   a.Category.ID,
		Name: a.Category.Name,
	}

	var tags []*tag.Tag
	for _, t := range a.Tags {
		tags = append(tags, &tag.Tag{
			ID:         t.ID,
			Name:       t.Name,
			CategoryID: t.CategoryID,
		})
	}

	return article.Article{
		ID:          &a.ID,
		FeedID:      a.FeedID,
		Feed:        feed,
		Category:    category,
		Title:       a.Title,
		Link:        a.Link,
		Description: a.Description,
		PubDate:     a.PubDate,
		GUID:        a.GUID,
		ImageUrl:    a.ImageUrl,
		Tags:        tags,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func (a *Article) fromDomain(domainArticle article.Article) Article {
	var articleID int64
	if domainArticle.ID != nil {
		articleID = *domainArticle.ID
	} else {
		articleID = 0 // NOTE: 一旦致し方なくデフォルト値セットしておく。。
	}

	// 新しいタグを挿入
	var tags []Tag
	for _, t := range domainArticle.Tags {
		// FIXME: toDomainとfromFomainでnilポインタだったりそうだ
		if t != nil { // nil チェックを追加
			tagEntity := (&Tag{}).fromDomain(t)
			tags = append(tags, tagEntity)
		}
	}

	return Article{
		ID:          articleID,
		Feed:        Feed{ID: domainArticle.Feed.ID, Title: domainArticle.Feed.Title},
		CategoryID:  domainArticle.Category.ID,
		Title:       domainArticle.Title,
		Link:        domainArticle.Link,
		Description: domainArticle.Description,
		PubDate:     domainArticle.PubDate,
		GUID:        domainArticle.GUID,
		ImageUrl:    domainArticle.ImageUrl,
		Tags:        tags,
		CreatedAt:   domainArticle.CreatedAt,
		UpdatedAt:   domainArticle.UpdatedAt,
	}
}
