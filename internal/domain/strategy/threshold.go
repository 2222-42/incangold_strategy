package strategy

import (
	"incangold_strategy/internal/domain"
)

type ThresholdStrategy struct {
	MaxHazards int // Leave when this many unique hazards are revealed (usually 0 to 5)
	name string
}

func NewThresholdStrategy(maxHazards int, name string) *ThresholdStrategy {
	return &ThresholdStrategy{MaxHazards: maxHazards, name: name}
}

func (s *ThresholdStrategy) Decide(game *domain.Game, round *domain.Round, self *domain.Player) domain.TurnChoice {
	// Count how many hazards have been revealed this round
	hazardsRevealed := 0
	for count := range round.RevealedHazards {
		// All revealed hazards are only 1 (if 2, it busts and round ends immediately)
		if count != domain.HazardTypeNone {
			hazardsRevealed++
		}
	}

	if hazardsRevealed >= s.MaxHazards {
		return domain.TurnChoiceLeave
	}

	return domain.TurnChoiceExplore
}

func (s *ThresholdStrategy) Name() string {
	return s.name
}
