package strategy

import (
	"incangold_strategy/internal/domain"
)

// RiskyStrategy implements the "More Risky" strategy from the research paper.
// It adjusts the gem threshold based on the player's lead over opponents.
//
// Parameters (from the paper's best result: base=9, R=2, S1=4, S2=2):
//   - BaseThreshold: Base number of gems in pocket before considering leaving (e.g. 9)
//   - R:             Lead threshold — if (myScore - bestOpponent) >= R, we're "leading"
//   - S1:            Extra gems to wait if leading (leading threshold = Base + S1, e.g. 13)
//   - S2:            Extra gems to wait if trailing (trailing threshold = Base + S2, e.g. 11)
//
// Decision logic:
//  1. Calculate effective pocket value = PocketScore + estimated artifact bonus
//  2. Compare against dynamic threshold to decide Stay/Leave
type RiskyStrategy struct {
	BaseThreshold    int
	R                int
	S1               int
	S2               int
	WithArtifactEV   bool // Whether to factor in artifact expected value
	name             string
}

// NewRiskyStrategy creates a new RiskyStrategy.
// WithArtifactEV=true will add the potential artifact value to the effective score,
// making the player more eager to leave alone when artifacts are on the path.
func NewRiskyStrategy(baseThreshold, r, s1, s2 int, withArtifactEV bool, name string) *RiskyStrategy {
	return &RiskyStrategy{
		BaseThreshold:  baseThreshold,
		R:              r,
		S1:             s1,
		S2:             s2,
		WithArtifactEV: withArtifactEV,
		name:           name,
	}
}

func (s *RiskyStrategy) Decide(game *domain.Game, round *domain.Round, self *domain.Player) domain.TurnChoice {
	// --- Step 1: Calculate dynamic threshold based on lead ---
	myTotalEstimate := self.TentScore + self.PocketScore

	maxOtherScore := 0
	for _, p := range game.Players {
		if p.ID == self.ID {
			continue
		}
		otherEstimate := p.TentScore + p.PocketScore
		if otherEstimate > maxOtherScore {
			maxOtherScore = otherEstimate
		}
	}

	lead := myTotalEstimate - maxOtherScore

	targetThreshold := s.BaseThreshold
	if lead >= s.R {
		targetThreshold += s.S1
	} else {
		targetThreshold += s.S2
	}

	// --- Step 2: Calculate effective pocket value ---
	// Base: gems already in pocket this round
	effectivePocket := self.PocketScore

	// Artifact EV: If there are visible artifacts and we might leave alone,
	// factor the artifact value into our effective score.
	// This lowers the "remaining gap" to the threshold, nudging toward leaving.
	if s.WithArtifactEV && round.RevealedArtifacts > 0 {
		activeCount := len(round.ActivePlayers())
		// Simple estimate: artifact value / number of currently active players
		// (reflects the probability of being the sole leaver)
		estimatedArtifactValue := artifactEstimate(game, round.RevealedArtifacts)
		if activeCount > 0 {
			effectivePocket += estimatedArtifactValue / activeCount
		}
	}

	// --- Step 3: Decide ---
	if effectivePocket >= targetThreshold {
		return domain.TurnChoiceLeave
	}
	return domain.TurnChoiceExplore
}

// artifactEstimate estimates the value of claiming the given number of artifacts right now.
// Returns 0 if artifacts are disabled or count is 0, consistent with Game.ClaimArtifacts.
func artifactEstimate(game *domain.Game, count int) int {
	if !game.ArtifactRule.Enabled || count == 0 {
		return 0
	}
	total := 0
	collected := game.ArtifactsCollected
	for i := 0; i < count; i++ {
		collected++
		if collected <= 3 {
			total += game.ArtifactRule.FirstThreeValue
		} else {
			total += game.ArtifactRule.FourthAndBeyondValue
		}
	}
	return total
}

func (s *RiskyStrategy) Name() string {
	return s.name
}

