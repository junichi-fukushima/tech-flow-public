package tag

type Tag struct {
	ID         int
	Name       string
	CategoryID int
}

func NewTag(
	id int,
	name string,
	categoryID int,
) Tag {
	return Tag{
		ID:         id,
		Name:       name,
		CategoryID: categoryID,
	}
}
