package domain

import "github.com/google/uuid"

type PlayerStatus string

const (
	PlayerStatusExploring PlayerStatus = "exploring"
	PlayerStatusLeft      PlayerStatus = "left"
)

type Player struct {
	ID          uuid.UUID
	Name        string
	TentScore   int    // Confirmed score stored in the tent
	PocketScore int    // Temporary score for the current round
	Status      PlayerStatus
	Strategy    Strategy // The logic this player uses
}

func NewPlayer(name string, strategy Strategy) *Player {
	return &Player{
		ID:          uuid.New(),
		Name:        name,
		TentScore:   0,
		PocketScore: 0,
		Status:      PlayerStatusExploring,
		Strategy:    strategy,
	}
}

// ResetForRound resets the pocket and exploring status at the start of a round
func (p *Player) ResetForRound() {
	p.PocketScore = 0
	p.Status = PlayerStatusExploring
}

// Leave lets the player return to tent safely
func (p *Player) Leave() {
	p.Status = PlayerStatusLeft
	p.TentScore += p.PocketScore
	p.PocketScore = 0
}

// Burst handles the case when 2 identical hazards appear: player loses pocket
func (p *Player) Burst() {
	p.Status = PlayerStatusLeft
	p.PocketScore = 0
}

// Collect returns some gems into the pocket
func (p *Player) Collect(amount int) {
	p.PocketScore += amount
}

// AddArtifactScore adds artifact score to the tent
func (p *Player) AddArtifactScore(amount int) {
	// Artifacts go straight to the tent since they are only collected when leaving alone
	p.TentScore += amount
}
