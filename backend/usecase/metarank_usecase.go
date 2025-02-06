package usecase

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"github.com/junichi-fukushima/tech-flow/backend/dto"
)

type MetarankUsecase interface {
	SendUserFeedback(event *auth.UserMetadataEvent) error
	SendInteractionFeedback(event *metaRank.InteractionEvent, userID string, articleID string) error
	SendRankingFeedback(event *metaRank.RankingEvent, userID string) error
	GetTrendingArticles(limit int) (*metaRank.Trending, error)
	GetRecommendation(userID string, items []metaRank.ItemDetail) (*metaRank.Recommend, error)
	GetUserMetadataEventByUserID(userID string) (*metaRank.UserMetadataEvent, error)
}

type metarankUsecase struct {
	metaRankRepository metaRank.MetaRankRepository
}

func NewMetarankUsecase(metaRankRepository metaRank.MetaRankRepository) MetarankUsecase {
	return &metarankUsecase{
		metaRankRepository: metaRankRepository,
	}
}

func (m *metarankUsecase) SendUserFeedback(event *auth.UserMetadataEvent) error {
	req := (&dto.UserMetadataEventRequest{}).FromUserMetadataEvent(event)
	return m.metaRankRepository.SendFeedback(req)
}

func (m *metarankUsecase) SendInteractionFeedback(event *metaRank.InteractionEvent, userID string, articleID string) error {
	req := (&dto.ClickMetadataEventRequest{}).FromClickMetadataEvent(event, userID, articleID)
	return m.metaRankRepository.SendFeedback(req)
}

func (m *metarankUsecase) SendRankingFeedback(event *metaRank.RankingEvent, userID string) error {
	req := (&dto.ImpMetadataEventRequest{}).FromImpMetadataEvent(event, userID)
	return m.metaRankRepository.SendFeedback(req)
}

func (m *metarankUsecase) GetTrendingArticles(limit int) (*metaRank.Trending, error) {
	req := (&dto.TrendingMetadataEventRequest{}).FromTrendingMetadataEvent(limit)
	return m.metaRankRepository.GetTrending(req)
}

func (m *metarankUsecase) GetRecommendation(userID string, items []metaRank.ItemDetail) (*metaRank.Recommend, error) {
	uuID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	itemJson, err := createItemJson(items)
	if err != nil {
		return nil, err
	}

	req := (&dto.RecommendationMetadataEventRequest{}).FromRecommendationMetadataEvent(uuID.String(), userID, itemJson)
	return m.metaRankRepository.GetRecommendation(req)
}

func (m *metarankUsecase) GetUserMetadataEventByUserID(userID string) (*metaRank.UserMetadataEvent, error) {
	return m.metaRankRepository.GetUserMetadataEventByUserID(userID)
}

func createItemJson(itemDetails []metaRank.ItemDetail) ([]byte, error) {
	var items []map[string]string
	for _, article := range itemDetails {
		items = append(items, map[string]string{"id": article.Item})
	}

	// sample) "items": [ {"id": "63ed3885-ed2e-f973-07b1-98f341b759cc"}, ...} ]
	return json.Marshal(items)
}
