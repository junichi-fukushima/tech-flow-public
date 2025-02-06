package usecase

import (
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
)

type ClickUsecase interface {
	CreateInteractionEvent(clickEvent metaRank.InteractionEvent) error
}

type clickUsecase struct {
	metaRankRepository metaRank.MetaRankRepository
}

func NewClickUsecase(metaRankRepository metaRank.MetaRankRepository) ClickUsecase {
	return &clickUsecase{
		metaRankRepository: metaRankRepository,
	}
}

func (u *clickUsecase) CreateInteractionEvent(clickEvent metaRank.InteractionEvent) error {
	return u.metaRankRepository.CreateClickEvent(clickEvent)
}
