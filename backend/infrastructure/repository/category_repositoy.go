package repository

import (
	"github.com/junichi-fukushima/tech-flow/backend/domain/category"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/gorm"
)

type Category struct {
	ID   int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name;type:varchar(255);not null"`
}

type GormCategoryRepository struct {
	db *gorm.SQLHandler
}

func NewCategoryRepository() (category.CategoryRepository, error) {
	sqlHandler, err := gorm.GetSQLHandler()
	if err != nil {
		return nil, err
	}
	return &GormCategoryRepository{db: sqlHandler}, nil
}

func (repo *GormCategoryRepository) GetCategoriesAll() ([]category.Category, error) {
	var categories []Category
	err := repo.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}

	var result []category.Category
	for _, category := range categories {
		result = append(result, category.toDomain())
	}
	return result, nil
}

func (c *Category) toDomain() category.Category {
	return category.Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (c *Category) fromDomain(domainCategory category.Category) Category {
	return Category{
		ID:   domainCategory.ID,
		Name: domainCategory.Name,
	}
}
