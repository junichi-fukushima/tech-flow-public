package repository

import (
	"time"

	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/gorm"
)

type User struct {
	ID                    string    `gorm:"column:id;primaryKey;autoIncrement"`
	SessionToken          string    `gorm:"column:session_token;type:varchar(255);not null"`
	HasFavoriteCategories bool      `gorm:"column:has_favorite_categories;not null"`
	CreatedAt             time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

type UserCategory struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	User       User      `gorm:"foreignKey:UserID;references:ID"`
	UserID     string    `gorm:"column:user_id;not null;index"`
	Category   Category  `gorm:"foreignKey:CategoryID;references:ID"`
	CategoryID int       `gorm:"column:category_id;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

type GormUserRepository struct {
	db *gorm.SQLHandler
}

func NewUserRepository() (auth.UserRepository, error) {
	sqlHandler, err := gorm.GetSQLHandler()
	if err != nil {
		return nil, err
	}
	// repository
	return &GormUserRepository{db: sqlHandler}, nil
}

func (repo *GormUserRepository) BulkCreateFavoriteCategories(userID string, categoryIDs []int) error {
	userCategories := make([]UserCategory, len(categoryIDs))

	for i, categoryID := range categoryIDs {
		userCategories[i] = UserCategory{
			UserID:     userID,
			CategoryID: categoryID,
		}
	}

	if err := repo.db.DB.Create(&userCategories).Error; err != nil {
		return err
	}

	// ユーザーの HasFavoriteCategories を true に更新
	if err := repo.db.DB.Model(&User{}).Where("id = ?", userID).Update("has_favorite_categories", true).Error; err != nil {
		return err
	}

	return nil
}
