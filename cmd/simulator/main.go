package main

import (
	"flag"
	"fmt"
	"os"

	"incangold_strategy/internal/application"
	"incangold_strategy/internal/domain"
	"incangold_strategy/internal/domain/strategy"
)

// scenarioAll runs all strategies head-to-head.
func scenarioAll() []application.PlayerFactory {
	return []application.PlayerFactory{
		func() *domain.Player {
			return domain.NewPlayer("Random 50%", strategy.NewRandomStrategy(0.5, "Random 50%"))
		},
		func() *domain.Player {
			return domain.NewPlayer("Random 10%", strategy.NewRandomStrategy(0.1, "Random 10%"))
		},
		func() *domain.Player {
			return domain.NewPlayer("Cowardly", strategy.NewThresholdStrategy(1, "Threshold (1 Hazard)"))
		},
		func() *domain.Player {
			return domain.NewPlayer("Normal", strategy.NewThresholdStrategy(2, "Threshold (2 Hazards)"))
		},
		func() *domain.Player {
			// More Risky best params (R=2, S1=4, S2=2) + artifact EV
			return domain.NewPlayer("More Risky (Best+ArtEV)", strategy.NewRiskyStrategy(9, 2, 4, 2, true, "Risky (Best+ArtEV)"))
		},
	}
}

// scenarioRisky runs the best Risky strategy against baseline strategies.
func scenarioRisky() []application.PlayerFactory {
	return []application.PlayerFactory{
		func() *domain.Player {
			return domain.NewPlayer("Normal", strategy.NewThresholdStrategy(2, "Threshold (2 Hazards)"))
		},
		func() *domain.Player {
			// Best from paper, with artifact EV
			return domain.NewPlayer("Risky Best+ArtEV", strategy.NewRiskyStrategy(9, 2, 4, 2, true, "Risky (Best: R=2,S1=4,S2=2+ArtEV)"))
		},
		func() *domain.Player {
			// Best from paper, without artifact EV
			return domain.NewPlayer("Risky Best", strategy.NewRiskyStrategy(9, 2, 4, 2, false, "Risky (Best: R=2,S1=4,S2=2)"))
		},
	}
}

// scenarioRiskyVsRisky compares different parameter configurations of RiskyStrategy
// and also the effect of artifact EV on the best variant.
func scenarioRiskyVsRisky() []application.PlayerFactory {
	return []application.PlayerFactory{
		func() *domain.Player {
			// Paper's best strategy with artifact EV
			return domain.NewPlayer("Risky Best+ArtEV", strategy.NewRiskyStrategy(9, 2, 4, 2, true, "Risky (R=2,S1=4,S2=2+ArtEV)"))
		},
		func() *domain.Player {
			// Paper's best strategy without artifact EV
			return domain.NewPlayer("Risky Best", strategy.NewRiskyStrategy(9, 2, 4, 2, false, "Risky (R=2,S1=4,S2=2)"))
		},
		func() *domain.Player {
			// Aggressive even when trailing: always use S1
			return domain.NewPlayer("Risky Always Aggressive", strategy.NewRiskyStrategy(9, 0, 4, 4, false, "Risky (R=0,S1=4,S2=4)"))
		},
		func() *domain.Player {
			// Catch-up: more aggressive when trailing (S2 is negative = lower bar)
			return domain.NewPlayer("Risky Catch-up", strategy.NewRiskyStrategy(9, 2, 2, -2, false, "Risky (R=2,S1=2,S2=-2)"))
		},
		func() *domain.Player {
			// Flat: no lead adjustment, pure base threshold
			return domain.NewPlayer("Risky Flat", strategy.NewRiskyStrategy(9, 100, 0, 0, false, "Risky (Flat base=9)"))
		},
	}
}

func main() {
	scenarioFlag := flag.String("scenario", "all", "Which strategy scenario to run:\n  all           - All strategies compete head-to-head\n  risky         - Risky strategies vs baselines\n  risky-vs-risky- Compare different Risky parameter sets")
	numGames := flag.Int("games", 100000, "Number of games to simulate")
	flag.Parse()

	var factories []application.PlayerFactory

	switch *scenarioFlag {
	case "all":
		factories = scenarioAll()
	case "risky":
		factories = scenarioRisky()
	case "risky-vs-risky":
		factories = scenarioRiskyVsRisky()
	default:
		fmt.Fprintf(os.Stderr, "Unknown scenario: %q\nUse --help for available options.\n", *scenarioFlag)
		os.Exit(1)
	}

	fmt.Printf("=== Scenario: %s | Games: %d ===\n\n", *scenarioFlag, *numGames)

	sim := application.NewSimulator(*numGames, factories)
	result := sim.Run()
	result.Print()
}
