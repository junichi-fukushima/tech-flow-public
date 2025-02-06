package usecase

import (
	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
)

type UserUsecase interface {
	UpsertFavCategories(user auth.User, categoryIDs []int) error
}

type userUsecase struct {
	userRepository auth.UserRepository
}

func NewUserUsecase(userRepository auth.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (u *userUsecase) UpsertFavCategories(user auth.User, categoryIDs []int) error {
	// user情報からお気に入り登録があるかどうかをチェック
	if user.HasFavoriteCategories {
		// ない場合処理終了
		return nil
	}

	// 好みのカテゴリ情報を登録し、ユーザーの好み登録状態をtrueにする
	err := u.userRepository.BulkCreateFavoriteCategories(user.ID, categoryIDs)
	if err != nil {
		return err
	}

	return nil
}
