package usecase

import (
	"github.com/google/uuid"
	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
)

type AuthUsecase interface {
	CreateAnonymousUser(token string) error
	FindAnonymousUser(token string) (*auth.User, error)
	FindUserMetaData(userID string) (*auth.UserMetadataEvent, error)
}

type authUsecase struct {
	authRepository auth.AuthRepository
}

func NewAuthUsecase(authRepository auth.AuthRepository) AuthUsecase {
	return &authUsecase{
		authRepository: authRepository,
	}
}

func (u *authUsecase) CreateAnonymousUser(token string) error {
	err := u.authRepository.CreateAnonymousUser(uuid.New().String(), token)
	if err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) FindAnonymousUser(token string) (*auth.User, error) {
	return u.authRepository.FindAnonymousUser(token)
}

func (u *authUsecase) FindUserMetaData(userID string) (*auth.UserMetadataEvent, error) {
	return u.authRepository.FindUserMetaData(userID)
}
