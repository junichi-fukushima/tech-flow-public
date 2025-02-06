package usecase

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"github.com/junichi-fukushima/tech-flow/backend/dto"
)

type ImpUsecase interface {
	CreateRankingEvent(sessionToken string, articles dto.Response, rankingEventID string) (*metaRank.RankingEvent, error)
	GetItemMetadataEventByArticleID(articleID int64) (*metaRank.ItemMetadataEvent, error)
}

type impUsecase struct {
	metaRankRepository metaRank.MetaRankRepository
	authRepository     auth.AuthRepository
}

func NewImpUsecase(metaRankRepository metaRank.MetaRankRepository, authRepository auth.AuthRepository) ImpUsecase {
	return &impUsecase{
		metaRankRepository: metaRankRepository,
		authRepository:     authRepository,
	}
}

func (u *impUsecase) CreateRankingEvent(sessionToken string, response dto.Response, rankingEventID string) (*metaRank.RankingEvent, error) {
	// find user
	user, err := u.authRepository.FindAnonymousUser(sessionToken)
	if err != nil {
		return nil, err
	}

	// find meta user data
	metaUser, err := u.authRepository.FindUserMetaData(user.ID)
	if err != nil {
		return nil, err
	}

	articles, err := createArticleJson(response)
	if err != nil {
		return nil, err
	}

	rankingEvent := metaRank.RankingEvent{
		ID:                  rankingEventID,
		Timestamp:           time.Now(),
		UserMetadataEventID: metaUser.ID,
		Articles:            articles,
	}

	err = u.metaRankRepository.CreateRankingEvent(rankingEvent)
	if err != nil {
		return nil, err
	}
	return &rankingEvent, nil

}

func createArticleJson(response dto.Response) ([]byte, error) {
	var items []map[string]string
	for _, article := range response.Data {
		items = append(items, map[string]string{"id": fmt.Sprintf("%d", article.ID)})
	}

	// sample)  [ {"id": "63ed3885-ed2e-f973-07b1-98f341b759cc"}, ...} ]
	return json.Marshal(items)
}

func (u *impUsecase) GetItemMetadataEventByArticleID(articleID int64) (*metaRank.ItemMetadataEvent, error) {
	return u.metaRankRepository.GetItemMetadataEventByArticleID(articleID)
}
