package usecase

import (
	"github.com/junichi-fukushima/tech-flow/backend/domain/category"
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
)

type CategoryUsecase interface {
	DecideCategory(tags []*tag.Tag) (category.Category, error)
	GetAllCategories() ([]category.Category, error)
}

type categoryUsecase struct {
	categoryRepository category.CategoryRepository
}

func NewCategoryUsecase(categoryRepository category.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{categoryRepository: categoryRepository}
}

// 全カテゴリを取得
func (u *categoryUsecase) GetAllCategories() ([]category.Category, error) {
	categories, err := u.categoryRepository.GetCategoriesAll()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// categoryの判定を実施
func (u *categoryUsecase) DecideCategory(tags []*tag.Tag) (category.Category, error) {
	// 全categories取得
	categories, _ := u.categoryRepository.GetCategoriesAll()

	// カテゴリのMapを作成
	categoryMap := make(map[int]category.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	// Tagと一致するカテゴリのMapを作成
	matchedCategories := make(map[int]category.Category)
	for _, tag := range tags {
		if category, exists := categoryMap[tag.CategoryID]; exists {
			matchedCategories[category.ID] = category
		}
	}

	var matchedCategoryList []category.Category
	for _, category := range matchedCategories {
		matchedCategoryList = append(matchedCategoryList, category)
	}

	// 該当カテゴリが0件→None
	if len(matchedCategories) == 0 {
		categoryNotFound := category.Category{
			ID:   8,
			Name: "NONE",
		}
		return categoryNotFound, nil
	}

	// 該当カテゴリが1件→合致したカテゴリを返す
	if len(matchedCategories) == 1 {
		for _, category := range matchedCategories {
			return category, nil
		}
	}

	var selectedCategory category.Category
	minPriority := int(^uint(0) >> 1) // 最大値
	for _, category := range matchedCategories {
		priority := calculatePriority(&category)
		if priority != -1 && priority < minPriority {
			// 重みづけの値が大きい時だけcategoryを更新する
			selectedCategory = category
			minPriority = priority
		}
	}
	return selectedCategory, nil
}

func calculatePriority(category *category.Category) int {
	switch category.Name {
	case "マネジメント":
		return 7
	case "最新技術":
		return 6
	case "AI":
		return 5
	case "IoT":
		return 4
	case "インフラ":
		return 3
	case "フレームワーク":
		return 2
	case "プログラミング言語":
		return 1
	default:
		return -1
	}
}
