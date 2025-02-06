package tag

type TagsRepository interface {
	GetTagsAll() ([]*Tag, error)
	GetTagsByCategoryID(categoryID int) ([]*Tag, error)
}
