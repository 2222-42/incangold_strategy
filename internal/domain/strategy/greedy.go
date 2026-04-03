package strategy

import (
	"incangold_strategy/internal/domain"
)

type GreedyStrategy struct {
	TargetPocket int // e.g. 0 means never leave, else leave if PocketScore >= TargetPocket
	name string
}

func NewGreedyStrategy(target int, name string) *GreedyStrategy {
	return &GreedyStrategy{TargetPocket: target, name: name}
}

func (s *GreedyStrategy) Decide(game *domain.Game, round *domain.Round, self *domain.Player) domain.TurnChoice {
	if s.TargetPocket > 0 && self.PocketScore >= s.TargetPocket {
		return domain.TurnChoiceLeave
	}
	// Default greedy: always explore
	return domain.TurnChoiceExplore
}

func (s *GreedyStrategy) Name() string {
	return s.name
}
