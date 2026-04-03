package domain

type TurnChoice string

const (
	TurnChoiceExplore TurnChoice = "explore"
	TurnChoiceLeave   TurnChoice = "leave"
)

type Strategy interface {
	Decide(game *Game, round *Round, self *Player) TurnChoice
	Name() string
}
