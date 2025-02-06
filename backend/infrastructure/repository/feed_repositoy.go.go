package repository

import (
	"time"

	"github.com/junichi-fukushima/tech-flow/backend/domain/feed"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/gorm"
	gormclause "gorm.io/gorm/clause"
)

type Feed struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Title         string     `gorm:"column:title;type:varchar(255);not null"`
	Link          string     `gorm:"column:link;type:text;not null"`
	Description   *string    `gorm:"column:description;type:text"`
	Category      *string    `gorm:"column:category;type:text"`
	Image         *string    `gorm:"column:image;type:text;"`
	Language      *string    `gorm:"column:language;type:varchar(10);"`
	LastBuildDate *time.Time `gorm:"column:last_build_date;type:datetime;"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

type GormFeedRepository struct {
	db *gorm.SQLHandler
}

func NewFeedRepository() (feed.FeedRepository, error) {
	sqlHandler, err := gorm.GetSQLHandler()
	if err != nil {
		return nil, err
	}

	return &GormFeedRepository{db: sqlHandler}, nil
}

// feedを全て取得する
func (repo *GormFeedRepository) GetFeedsAll() ([]feed.Feed, error) {
	var feeds []Feed
	err := repo.db.DB.Find(&feeds).Error
	if err != nil {
		return nil, err
	}

	var domainFeeds []feed.Feed
	for _, feed := range feeds {
		domainFeed, err := feed.toDomain()
		if err != nil {
			return nil, err
		}
		domainFeeds = append(domainFeeds, domainFeed)
	}
	return domainFeeds, nil
}

// idの範囲を指定取得
// FIXME: あんま綺麗ではない。妥協案
func (repo *GormFeedRepository) GetFeedsByIDRange(startID, endID int) ([]feed.Feed, error) {
	var feeds []Feed
	err := repo.db.DB.Where("id >= ? AND id < ?", startID, endID).Find(&feeds).Error
	if err != nil {
		return nil, err
	}

	var domainFeeds []feed.Feed
	for _, feed := range feeds {
		domainFeed, err := feed.toDomain()
		if err != nil {
			return nil, err
		}
		domainFeeds = append(domainFeeds, domainFeed)
	}
	return domainFeeds, nil
}

func (repo *GormFeedRepository) UpsertRss(domainFeed feed.Feed) error {
	// ドメインからDB用の構造体に変換
	feedEntity, err := FromDomain(domainFeed)
	if err != nil {
		return err
	}

	err = repo.db.DB.Clauses(
		gormclause.OnConflict{
			Columns: []gormclause.Column{{Name: "id"}},
			DoUpdates: gormclause.AssignmentColumns([]string{
				"title",
				"link",
				"description",
				"category",
				"image",
				"language",
				"last_build_date",
				"updated_at",
			}), // 更新するカラムを指定
		},
	).Create(&feedEntity).Error

	return err
}

func (f *Feed) toDomain() (feed.Feed, error) {
	return feed.Feed{
		ID:             f.ID,
		Title:          f.Title,
		Link:           f.Link,
		Description:    f.Description,
		CategoryOfFeed: f.Category, // 記事で管理するカテゴリとは異なる！
		Image:          f.Image,
		Language:       f.Language,
		LastBuildDate:  f.LastBuildDate,
		CreatedAt:      f.CreatedAt,
		UpdatedAt:      f.UpdatedAt,
	}, nil
}

func FromDomain(domainFeed feed.Feed) (Feed, error) {
	return Feed{
		ID:            domainFeed.ID,
		Title:         domainFeed.Title,
		Link:          domainFeed.Link,
		Description:   domainFeed.Description,
		Category:      domainFeed.CategoryOfFeed, // 記事で管理するカテゴリとは異なる！
		Image:         domainFeed.Image,
		Language:      domainFeed.Language,
		LastBuildDate: domainFeed.LastBuildDate,
		CreatedAt:     domainFeed.CreatedAt,
		UpdatedAt:     domainFeed.UpdatedAt,
	}, nil
}
