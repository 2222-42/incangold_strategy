package domain

import (
	"fmt"
)

type Round struct {
	RoundNum          int
	Deck              *Deck
	Players           []*Player
	BoardGems         int
	RevealedHazards   map[HazardType]int
	RevealedArtifacts int
	PathCards         []Card
	IsOver            bool
	Game              *Game // Reference to the parent game to pass to strategy
}

func NewRound(roundNum int, deck *Deck, players []*Player, game *Game) *Round {
	return &Round{
		RoundNum:          roundNum,
		Deck:              deck,
		Players:           players,
		BoardGems:         0,
		RevealedHazards:   make(map[HazardType]int),
		RevealedArtifacts: 0,
		PathCards:         []Card{},
		IsOver:            false,
		Game:              game,
	}
}

func (r *Round) ActivePlayers() []*Player {
	var active []*Player
	for _, p := range r.Players {
		if p.Status == PlayerStatusExploring {
			active = append(active, p)
		}
	}
	return active
}

func (r *Round) Play() {
	r.Deck.PrepareRound()

	for _, p := range r.Players {
		p.ResetForRound()
	}

	for !r.IsOver {
		r.Step()
	}
}

// Step performs a single turn: decision, process leavers, then draw a card if anyone is still exploring
func (r *Round) Step() {
	active := r.ActivePlayers()

	if len(active) == 0 {
		r.IsOver = true
		return
	}

	// 1. All active players make a decision simultaneously
	choices := make(map[*Player]TurnChoice)
	for _, p := range active {
		choices[p] = p.Strategy.Decide(r.Game, r, p)
	}

	// 2. Process Leaves
	var leavingPlayers []*Player
	for p, choice := range choices {
		if choice == TurnChoiceLeave {
			leavingPlayers = append(leavingPlayers, p)
		}
	}

	if len(leavingPlayers) > 0 {
		// Split board gems
		if r.BoardGems > 0 {
			share := r.BoardGems / len(leavingPlayers)
			for _, p := range leavingPlayers {
				p.Collect(share)
			}
			r.BoardGems = r.BoardGems % len(leavingPlayers)
		}

		// Claims artifacts if only ONE player leaves
		if len(leavingPlayers) == 1 && r.RevealedArtifacts > 0 {
			score := r.Game.ClaimArtifacts(r.RevealedArtifacts)
			leavingPlayers[0].AddArtifactScore(score)
			
			// Remove collected artifacts from the deck completely
			for _, pc := range r.PathCards {
				if pc.Type == CardTypeArtifact {
					r.Deck.RemoveAvailableCard(pc)
				}
			}
			r.RevealedArtifacts = 0
			// Remove from PathCards visually if needed, but not strictly necessary for stats
		}

		for _, p := range leavingPlayers {
			p.Leave() // Moves pocket score to tent
		}
	}

	// Check if anyone is left to explore
	active = r.ActivePlayers()
	if len(active) == 0 {
		r.IsOver = true
		// Everyone left, the round ends safely
		return
	}

	// 3. Draw a Card
	drawn := r.Deck.Draw()
	if drawn == nil {
		fmt.Println("Warning: Deck is empty")
		r.IsOver = true
		return
	}

	r.PathCards = append(r.PathCards, *drawn)

	switch drawn.Type {
	case CardTypeTreasure:
		share := drawn.Value / len(active)
		for _, p := range active {
			p.Collect(share)
		}
		r.BoardGems += drawn.Value % len(active)

	case CardTypeHazard:
		r.RevealedHazards[drawn.HazardType]++
		if r.RevealedHazards[drawn.HazardType] == 2 {
			// Burst!
			r.IsOver = true
			for _, p := range active {
				p.Burst()
			}
			// The causing hazard card is removed from future rounds
			r.Deck.RemoveAvailableCard(*drawn)
		}

	case CardTypeArtifact:
		r.RevealedArtifacts++
	}
}


