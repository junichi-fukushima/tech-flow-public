package repository

import (
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/gorm"
)

type Tag struct {
	ID         int      `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string   `gorm:"column:name;type:varchar(255);not null;unique"`
	CategoryID int      `gorm:"column:category_id;not null"`
	Category   Category `gorm:"foreignKey:CategoryID;references:ID"`
}

type GormTagRepository struct {
	db *gorm.SQLHandler
}

func NewTagsRepository() (tag.TagsRepository, error) {
	sqlHandler, err := gorm.GetSQLHandler()
	if err != nil {
		return nil, err
	}
	return &GormTagRepository{db: sqlHandler}, nil
}

func (repo *GormTagRepository) GetTagsAll() ([]*tag.Tag, error) {
	var tags []Tag
	err := repo.db.Select(&tags, "SELECT * FROM tags")
	if err != nil {
		return nil, err
	}

	var result []*tag.Tag
	for _, tag := range tags {
		result = append(result, tag.toDomain())
	}
	return result, nil
}

func (repo *GormTagRepository) GetTagsByCategoryID(categoryID int) ([]*tag.Tag, error) {
	var tags []Tag
	err := repo.db.DB.Where("category_id = ?", categoryID).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	// ドメインモデルに変換
	var result []*tag.Tag
	for _, tag := range tags {
		result = append(result, tag.toDomain())
	}
	return result, nil
}

func (t *Tag) toDomain() *tag.Tag {
	return &tag.Tag{
		ID:         t.ID,
		Name:       t.Name,
		CategoryID: t.CategoryID,
	}
}

func (t *Tag) fromDomain(domainTag *tag.Tag) Tag {
	return Tag{
		ID:         domainTag.ID,
		Name:       domainTag.Name,
		CategoryID: domainTag.CategoryID,
	}
}
