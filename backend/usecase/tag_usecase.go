package usecase

import (
	"strings"

	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
)

type TagUsecase interface {
	GetAllTags() ([]*tag.Tag, error)
	DecideTags(title string, description *string) ([]*tag.Tag, error)
}

type tagUsecase struct {
	tagsRepository tag.TagsRepository
}

func NewTagUsecase(tagsRepository tag.TagsRepository) TagUsecase {
	return &tagUsecase{tagsRepository: tagsRepository}
}

// 全タグを取得
func (u *tagUsecase) GetAllTags() ([]*tag.Tag, error) {
	tags, err := u.tagsRepository.GetTagsAll()
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// tagの判定をする
func (u *tagUsecase) DecideTags(title string, description *string) ([]*tag.Tag, error) {
	tags, err := u.tagsRepository.GetTagsAll()
	if err != nil {
		return nil, err
	}

	// (第1優先)まずはタイトルからTagの判定を実行する
	var matchedTags []*tag.Tag
	for _, t := range tags {
		if strings.Contains(title, t.Name) {
			matchedTags = append(matchedTags, t)
		}
	}

	// (第2優先)タイトルで判定出来なかった場合のみdescriptionを見るようにする
	if len(matchedTags) == 0 && description != nil {
		for _, t := range tags {
			if strings.Contains(*description, t.Name) {
				matchedTags = append(matchedTags, t)
			}
		}
	}

	// タグが判定できなかった場合、"NONE"タグを返す
	if len(matchedTags) == 0 {
		tagNotFound := []*tag.Tag{
			{
				ID:         108,
				Name:       "NONE",
				CategoryID: 8,
			},
		}
		return tagNotFound, nil
	}

	return matchedTags, nil
}
