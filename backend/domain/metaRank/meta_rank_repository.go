package metaRank

type MetaRankRepository interface {
	// for imp
	CreateRankingEvent(rankingEvent RankingEvent) error
	// for click
	CreateClickEvent(InteractionEvent InteractionEvent) error
	// for create itemMetadataEvents
	GetItemMetadataEventsByArticleIDs(articleIDs []int64) ([]ItemMetadataEvent, error)
	CreateItemMetadataEvents(events []ItemMetadataEvent) error
	GetItemMetadataEventByArticleID(articleID int64) (*ItemMetadataEvent, error)
	GetItemMetadataEventByIDs(IDs []string) ([]ItemMetadataEvent, error)
	GetUserMetadataEventByUserID(userID string) (*UserMetadataEvent, error)

	// SendFeedback sends user feedback to the metarank
	SendFeedback(event any) error
	// GetTrending get trending articles from metarank
	GetTrending(event any) (*Trending, error)
	// GetRecommendation get recommended articles from metarank
	GetRecommendation(event any) (*Recommend, error)
}
