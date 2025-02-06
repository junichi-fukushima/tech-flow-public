package repository

import (
	"github.com/google/uuid"
	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/gorm"
	g "gorm.io/gorm"
	"time"
)

type GormAuthRepository struct {
	db *gorm.SQLHandler
}

func NewAuthRepository() (auth.AuthRepository, error) {
	sqlHandler, err := gorm.GetSQLHandler()
	if err != nil {
		return nil, err
	}

	return &GormAuthRepository{db: sqlHandler}, nil
}

func (repo GormAuthRepository) CreateAnonymousUser(userID string, sessionToken string) error {
	user := auth.User{
		ID:           userID,
		SessionToken: sessionToken,
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	userEvent := auth.UserMetadataEvent{
		ID:        id.String(),
		Timestamp: time.Now(),
		UserID:    userID,
	}

	err = repo.db.DB.Transaction(func(tx *g.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&userEvent).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (repo GormAuthRepository) FindAnonymousUser(sessionToken string) (*auth.User, error) {
	var user auth.User

	err := repo.db.DB.Where("session_token = ?", sessionToken).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo GormAuthRepository) FindUserMetaData(userID string) (*auth.UserMetadataEvent, error) {
	var user auth.UserMetadataEvent

	err := repo.db.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
