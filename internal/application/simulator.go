package application

import (
	"fmt"
	"sync"
	"time"

	"incangold_strategy/internal/domain"
)

type PlayerFactory func() *domain.Player

type Simulator struct {
	NumGames int
	Factories []PlayerFactory
}

type SimulationResult struct {
	WinCounts  map[string]int
	AvrScores  map[string]float64
	TotalGames int
	Duration   time.Duration
}

func NewSimulator(numGames int, factories []PlayerFactory) *Simulator {
	return &Simulator{
		NumGames:  numGames,
		Factories: factories,
	}
}

func (s *Simulator) Run() SimulationResult {
	start := time.Now()

	winCounts := make(map[string]int)
	totalScores := make(map[string]int)
	var mu sync.Mutex

	var wg sync.WaitGroup
	
	// Workers approach: process games in parallel
	numWorkers := 8
	gamesPerWorker := s.NumGames / numWorkers
	
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		
		go func(workerID int) {
			defer wg.Done()
			
			localWinCounts := make(map[string]int)
			localTotalScores := make(map[string]int)

			gamesToRun := gamesPerWorker
			if workerID == numWorkers-1 {
				// Last worker takes the remainder if not easily divisible
				gamesToRun += s.NumGames % numWorkers
			}

			for i := 0; i < gamesToRun; i++ {
				var players []*domain.Player
				for _, factory := range s.Factories {
					players = append(players, factory())
				}

				game := domain.NewGame(players)
				game.Play()

				winners := game.GetWinner()
				// Track wins
				for _, w := range winners {
					// We use strategy name + id or just strategy name for simple aggregation
					localWinCounts[w.Strategy.Name()]++
				}

				// Track scores
				for _, p := range game.Players {
					localTotalScores[p.Strategy.Name()] += p.TentScore
				}
			}

			// Merge local results into global
			mu.Lock()
			for name, wins := range localWinCounts {
				winCounts[name] += wins
			}
			for name, score := range localTotalScores {
				totalScores[name] += score
			}
			mu.Unlock()

		}(w)
	}

	wg.Wait()

	duration := time.Since(start)

	avrScores := make(map[string]float64)
	for name, count := range totalScores {
		avrScores[name] = float64(count) / float64(s.NumGames)
	}

	return SimulationResult{
		WinCounts:  winCounts,
		AvrScores:  avrScores,
		TotalGames: s.NumGames,
		Duration:   duration,
	}
}

func (r *SimulationResult) Print() {
	fmt.Printf("Simulation completed in %v\n", r.Duration)
	fmt.Printf("Total Games: %d\n", r.TotalGames)
	fmt.Println("-----------------------------------------------------")
	fmt.Printf("%-25s | %-10s | %-10s\n", "Strategy", "Win Rate", "Avg Score")
	fmt.Println("-----------------------------------------------------")
	for name := range r.AvrScores {
		winRate := float64(r.WinCounts[name]) / float64(r.TotalGames) * 100
		fmt.Printf("%-25s | %5.2f%%     | %5.2f\n", name, winRate, r.AvrScores[name])
	}
	fmt.Println("-----------------------------------------------------")
}
