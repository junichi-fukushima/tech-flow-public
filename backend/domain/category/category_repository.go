package category

type CategoryRepository interface {
	GetCategoriesAll() ([]Category, error)
}
