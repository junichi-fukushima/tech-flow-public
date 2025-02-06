package repository

import (
	"encoding/json"
	"fmt"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/metarank"
	"time"

	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/gorm"
)

type RankingEvent struct {
	ID        int64
	Timestamp time.Time
	Fields    *json.RawMessage `gorm:"type:json"`
	UserID    string
	Articles  json.RawMessage `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GormMetaRankRepository struct {
	db     *gorm.SQLHandler
	client *metarank.Client
}

func NewMetaRankRepository() (metaRank.MetaRankRepository, error) {
	sqlHandler, err := gorm.GetSQLHandler()
	if err != nil {
		return nil, err
	}
	client := metarank.NewClient()
	return &GormMetaRankRepository{db: sqlHandler, client: client}, nil
}

// InteractionEvent----------------
func (repo *GormMetaRankRepository) CreateClickEvent(clickEvent metaRank.InteractionEvent) error {
	return repo.db.DB.Create(&clickEvent).Error
}

// RankingEvent----------------
func (repo *GormMetaRankRepository) CreateRankingEvent(rankingEvent metaRank.RankingEvent) error {
	return repo.db.DB.Create(&rankingEvent).Error
}

// ItemMetadataEvents----------------
func (repo *GormMetaRankRepository) GetItemMetadataEventsByArticleIDs(articleIDs []int64) ([]metaRank.ItemMetadataEvent, error) {
	if len(articleIDs) == 0 {
		return nil, fmt.Errorf("articleIDs slice is empty")
	}

	var events []metaRank.ItemMetadataEvent

	err := repo.db.DB.
		Model(&metaRank.ItemMetadataEvent{}).
		Where("article_id IN ?", articleIDs).
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch item metadata events: %w", err)
	}

	return events, nil
}

func (repo *GormMetaRankRepository) CreateItemMetadataEvents(events []metaRank.ItemMetadataEvent) error {
	if len(events) == 0 {
		return nil // 処理するイベントがない場合は成功として返す
	}

	err := repo.db.DB.Create(&events).Error
	if err != nil {
		return fmt.Errorf("failed to create item metadata events: %w", err)
	}

	return nil
}

func (repo *GormMetaRankRepository) GetItemMetadataEventByIDs(IDs []string) ([]metaRank.ItemMetadataEvent, error) {
	events := make([]metaRank.ItemMetadataEvent, 0, len(IDs))
	err := repo.db.DB.
		Preload("Article").
		Preload("Article.Category").
		Preload("Article.Tags").
		Find(&events, IDs).
		Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (repo *GormMetaRankRepository) GetUserMetadataEventByUserID(userID string) (*metaRank.UserMetadataEvent, error) {
	var event metaRank.UserMetadataEvent

	err := repo.db.DB.Where("user_id = ?", userID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (repo *GormMetaRankRepository) GetItemMetadataEventByArticleID(articleID int64) (*metaRank.ItemMetadataEvent, error) {
	var event metaRank.ItemMetadataEvent

	err := repo.db.DB.Where("article_id = ?", articleID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// SendFeedback sends user feedback to metarank
func (repo *GormMetaRankRepository) SendFeedback(event any) error {
	return repo.client.SendFeedback(event)
}

// GetTrending get trending articles from metarank
func (repo *GormMetaRankRepository) GetTrending(event any) (*metaRank.Trending, error) {
	return repo.client.GetTrending(event)
}

// GetRecommendation get recommended articles from metarank
func (repo *GormMetaRankRepository) GetRecommendation(event any) (*metaRank.Recommend, error) {
	return repo.client.GetRecommendation(event)
}
