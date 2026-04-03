package domain

import (
	"math/rand"
	"time"
)

type CardType string

const (
	CardTypeTreasure CardType = "treasure"
	CardTypeHazard   CardType = "hazard"
	CardTypeArtifact CardType = "artifact"
)

type HazardType string

const (
	HazardTypeSpider  HazardType = "spider"
	HazardTypeSnake   HazardType = "snake"
	HazardTypeLava    HazardType = "lava"
	HazardTypeBoulder HazardType = "boulder"
	HazardTypeMummy   HazardType = "mummy"
	HazardTypeNone    HazardType = "none" // For non-hazard cards
)

type Card struct {
	Type       CardType
	HazardType HazardType
	Value      int
	ID         int // to distinguish cards if needed
}

// Default Deck Configuration
var TreasureValues = []int{1, 2, 3, 4, 5, 5, 7, 7, 9, 11, 11, 13, 14, 15}

var HazardTypes = []HazardType{
	HazardTypeSpider,
	HazardTypeSnake,
	HazardTypeLava,
	HazardTypeBoulder,
	HazardTypeMummy,
}

type Deck struct {
	availableCards []Card
	drawPile       []Card
}

func NewDeck() *Deck {
	d := &Deck{
		availableCards: []Card{},
		drawPile:       []Card{},
	}
	idCounter := 0

	// Add initial Treasures
	for _, val := range TreasureValues {
		idCounter++
		d.availableCards = append(d.availableCards, Card{
			Type:       CardTypeTreasure,
			HazardType: HazardTypeNone,
			Value:      val,
			ID:         idCounter,
		})
	}

	// Add initial Hazards
	for _, ht := range HazardTypes {
		for i := 0; i < 3; i++ {
			idCounter++
			d.availableCards = append(d.availableCards, Card{
				Type:       CardTypeHazard,
				HazardType: ht,
				Value:      0,
				ID:         idCounter,
			})
		}
	}
	return d
}

// AddAvailableCard adds a card (like a new Artifact) to the persistent available pool
func (d *Deck) AddAvailableCard(c Card) {
	d.availableCards = append(d.availableCards, c)
}

// RemoveAvailableCard removes a card permanently from the game (e.g. collected artifact, burst hazard)
func (d *Deck) RemoveAvailableCard(c Card) {
	for i, ac := range d.availableCards {
		if ac.ID == c.ID {
			d.availableCards = append(d.availableCards[:i], d.availableCards[i+1:]...)
			return
		}
	}
}

// PrepareRound prepares the draw pile for a new round
func (d *Deck) PrepareRound() {
	d.drawPile = make([]Card, len(d.availableCards))
	copy(d.drawPile, d.availableCards)
	d.Shuffle()
}

func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(d.drawPile) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		d.drawPile[i], d.drawPile[j] = d.drawPile[j], d.drawPile[i]
	}
}

func (d *Deck) Draw() *Card {
	if len(d.drawPile) == 0 {
		return nil
	}
	card := d.drawPile[0]
	d.drawPile = d.drawPile[1:]
	return &card
}
