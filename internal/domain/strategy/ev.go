package strategy

import (
	"incangold_strategy/internal/domain"
)

// EVStrategy implements the Expected Value (EV) based decision from the research paper (§3.2):
//
//	E[V] = Upside - Downside
//
// Where:
//
//	Upside   = avgGoodCardValue * goodRate / activePlayers
//	Downside = PocketScore * deathRate
//
// Card classification (from the paper §3.1):
//   - Good cards:    treasure cards + artifact cards remaining in the draw pile
//   - Neutral cards: hazard types with 0 cards revealed so far (drawing one is "safe" this step)
//   - Bad cards:     hazard types with exactly 1 card already revealed (drawing a 2nd causes burst)
//
// The player leaves when E[V] <= 0, meaning staying is no longer worth the risk.
type EVStrategy struct {
	name string
}

func NewEVStrategy(name string) *EVStrategy {
	return &EVStrategy{name: name}
}

func (s *EVStrategy) Decide(game *domain.Game, round *domain.Round, self *domain.Player) domain.TurnChoice {
	remaining := round.Deck.RemainingCards()
	total := len(remaining)
	if total == 0 {
		return domain.TurnChoiceLeave
	}

	activePlayers := len(round.ActivePlayers())
	if activePlayers == 0 {
		return domain.TurnChoiceLeave
	}

	// --- Classify remaining cards ---
	var goodCards []domain.Card // treasure + artifact
	badCount := 0               // hazard types with 1 card already revealed (lethal next draw)

	for _, c := range remaining {
		switch c.Type {
		case domain.CardTypeTreasure, domain.CardTypeArtifact:
			goodCards = append(goodCards, c)
		case domain.CardTypeHazard:
			// A hazard card remaining in the deck is "bad" only if its type has already appeared once
			if round.RevealedHazards[c.HazardType] == 1 {
				badCount++
			}
			// If RevealedHazards == 0, it's a neutral card (not lethal next draw)
			// If RevealedHazards == 2, the round would already be over, so this case won't occur
		}
	}

	goodCount := len(goodCards)
	goodRate := float64(goodCount) / float64(total)
	deathRate := float64(badCount) / float64(total)

	// --- Upside: expected gems gained from the next card draw ---
	upside := 0.0
	if goodCount > 0 {
		// Average value of the good cards (treasure values; artifacts counted at their estimated worth)
		totalGoodValue := 0.0
		for _, c := range goodCards {
			switch c.Type {
			case domain.CardTypeTreasure:
				totalGoodValue += float64(c.Value)
			case domain.CardTypeArtifact:
				// Estimate artifact value based on how many have been collected so far
				totalGoodValue += float64(artifactEstimate(game, 1))
			}
		}
		avgGoodValue := totalGoodValue / float64(goodCount)
		upside = avgGoodValue * goodRate / float64(activePlayers)
	}

	// --- Downside: expected gem loss if we burst next draw ---
	// If we burst, we lose our entire PocketScore for this round.
	downside := float64(self.PocketScore) * deathRate

	// --- Decision ---
	// Leave when the expected gain no longer justifies the risk of losing what we have.
	ev := upside - downside
	if ev <= 0 {
		return domain.TurnChoiceLeave
	}
	return domain.TurnChoiceExplore
}

func (s *EVStrategy) Name() string {
	return s.name
}
