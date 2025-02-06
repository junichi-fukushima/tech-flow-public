package feed

type FeedRepository interface {
	GetFeedsAll() ([]Feed, error)
	UpsertRss(domainFeed Feed) error
	GetFeedsByIDRange(startID, endID int) ([]Feed, error)
}
