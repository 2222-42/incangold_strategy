package domain

import (
	"github.com/google/uuid"
)

type ArtifactRule struct {
	Enabled               bool
	FirstThreeValue       int
	FourthAndBeyondValue  int
}

func DefaultArtifactRule() ArtifactRule {
	return ArtifactRule{
		Enabled:              true,
		FirstThreeValue:      5,
		FourthAndBeyondValue: 10,
	}
}

type Game struct {
	ID                 uuid.UUID
	Players            []*Player
	Deck               *Deck
	Rounds             []*Round
	ArtifactsCollected int
	ArtifactRule       ArtifactRule
}

func NewGame(players []*Player) *Game {
	return &Game{
		ID:                 uuid.New(),
		Players:            players,
		Deck:               NewDeck(),
		Rounds:             []*Round{},
		ArtifactsCollected: 0,
		ArtifactRule:       DefaultArtifactRule(),
	}
}

// SetArtifactRule allows customizing the artifact scoring logic
func (g *Game) SetArtifactRule(rule ArtifactRule) {
	g.ArtifactRule = rule
}

// ClaimArtifacts calculates the score for a given number of artifacts and updates the collected count
func (g *Game) ClaimArtifacts(count int) int {
	if !g.ArtifactRule.Enabled || count == 0 {
		return 0
	}
	
	totalScore := 0
	for i := 0; i < count; i++ {
		g.ArtifactsCollected++
		if g.ArtifactsCollected <= 3 {
			totalScore += g.ArtifactRule.FirstThreeValue
		} else {
			totalScore += g.ArtifactRule.FourthAndBeyondValue
		}
	}
	return totalScore
}

func (g *Game) Play() {
	artifactIDBase := 100
	for roundNum := 1; roundNum <= 5; roundNum++ {
		// Add 1 new artifact to the available deck pool for the round
		g.Deck.AddAvailableCard(Card{
			Type:       CardTypeArtifact,
			HazardType: HazardTypeNone,
			Value:      0, // Dynamic value evaluated internally
			ID:         artifactIDBase + roundNum,
		})

		round := NewRound(roundNum, g.Deck, g.Players, g)
		g.Rounds = append(g.Rounds, round)
		round.Play()
	}
}

// GetWinner returns the player(s) with the highest tent score
func (g *Game) GetWinner() []*Player {
	var winners []*Player
	maxScore := -1

	for _, p := range g.Players {
		if p.TentScore > maxScore {
			maxScore = p.TentScore
			winners = []*Player{p}
		} else if p.TentScore == maxScore {
			winners = append(winners, p)
		}
	}

	return winners
}
