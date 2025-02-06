package auth

// NOTE: 本当はドメインとしてはuser
type AuthRepository interface {
	CreateAnonymousUser(userID string, sessionToken string) error
	FindAnonymousUser(sessionToken string) (*User, error)
	FindUserMetaData(userID string) (*UserMetadataEvent, error)
}
