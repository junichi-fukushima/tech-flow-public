package feed

import "time"

type Feed struct {
	ID             int64
	Title          string
	Link           string
	Description    *string
	CategoryOfFeed *string
	Image          *string
	Language       *string
	LastBuildDate  *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewFeed(
	id int64,
	title string,
	link string,
	description *string,
	categoryOfFeed *string,
	image *string,
	language *string,
	lastBuildDate *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) Feed {
	return Feed{
		ID:             id,
		Title:          title,
		Link:           link,
		Description:    description,
		CategoryOfFeed: categoryOfFeed,
		Image:          image,
		Language:       language,
		LastBuildDate:  lastBuildDate,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
