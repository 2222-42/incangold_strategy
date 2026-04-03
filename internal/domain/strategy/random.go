package strategy

import (
	"math/rand"
	"time"

	"incangold_strategy/internal/domain"
)

type RandomStrategy struct {
	LeaveProb float64
	rnd       *rand.Rand
	name string
}

func NewRandomStrategy(leaveProb float64, name string) *RandomStrategy {
	return &RandomStrategy{
		LeaveProb: leaveProb,
		rnd:       rand.New(rand.NewSource(time.Now().UnixNano())),
		name:      name,
	}
}

func (s *RandomStrategy) Decide(game *domain.Game, round *domain.Round, self *domain.Player) domain.TurnChoice {
	if s.rnd.Float64() < s.LeaveProb {
		return domain.TurnChoiceLeave
	}
	return domain.TurnChoiceExplore
}

func (s *RandomStrategy) Name() string {
	return s.name
}
