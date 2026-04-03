package main

import (
	"incangold_strategy/internal/application"
	"incangold_strategy/internal/domain"
	"incangold_strategy/internal/domain/strategy"
)

func main() {
	// Configure the players that will participate
	factories := []application.PlayerFactory{
		func() *domain.Player {
			return domain.NewPlayer("Random 50%", strategy.NewRandomStrategy(0.5, "Random 50%"))
		},
		func() *domain.Player {
			return domain.NewPlayer("Random 10%", strategy.NewRandomStrategy(0.1, "Random 10%"))
		},
		func() *domain.Player {
			return domain.NewPlayer("Greedy", strategy.NewGreedyStrategy(0, "Greedy (Never Leave)"))
		},
		func() *domain.Player {
			return domain.NewPlayer("Cowardly (1 Hazard)", strategy.NewThresholdStrategy(1, "Threshold (1 Hazard)"))
		},
		func() *domain.Player {
			return domain.NewPlayer("Normal (2 Hazards)", strategy.NewThresholdStrategy(2, "Threshold (2 Hazards)"))
		},
		func() *domain.Player {
			return domain.NewPlayer("Target Pokcet 10", strategy.NewGreedyStrategy(10, "Greedy (Target 10)"))
		},
	}

	sim := application.NewSimulator(100000, factories)
	result := sim.Run()
	result.Print()
}
