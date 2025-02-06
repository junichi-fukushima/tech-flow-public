package auth

// NOTE: 本当はドメインとしてはuserなので、AuthRepositoryのものをこちらに統合したい
type UserRepository interface {
	BulkCreateFavoriteCategories(userID string, categoryIDs []int) error
}
