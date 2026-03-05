package clean

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/velosypedno/resource-allocation/base"
	"github.com/velosypedno/resource-allocation/strategy"
)

type Simulator interface {
	Simulate(weights []float64) *strategy.SimulationResult
	TotalOperations() int
}

type AnnealingStrategy struct {
	InitialTemp      float64
	MinTemp          float64
	Alpha            float64
	Iterations       int
	SwapsPerMutation int
}

func NewAnnealingStrategy(initialTemp, minTemp, alpha float64, iterations int, swaps int) *AnnealingStrategy {
	return &AnnealingStrategy{
		InitialTemp:      initialTemp,
		MinTemp:          minTemp,
		Alpha:            alpha,
		Iterations:       iterations,
		SwapsPerMutation: swaps,
	}
}

func (s *AnnealingStrategy) Name() string {
	return "Simulated Annealing"
}

func (s *AnnealingStrategy) Description() string {
	return fmt.Sprintf(
		"Optimization using thermodynamic annealing. It explores the solution space by occasionally accepting worse moves.\n\n"+
			"| Parameter         | Value      |\n"+
			"|-------------------|------------|\n"+
			"| Initial Temp      | %.2f       |\n"+
			"| Min Temp          | %.4f       |\n"+
			"| Alpha (Cooling)   | %.4f       |\n"+
			"| Iterations/T      | %d         |\n"+
			"| Swaps Per Mutate  | %d         |",
		s.InitialTemp, s.MinTemp, s.Alpha, s.Iterations, s.SwapsPerMutation,
	)
}

func (s *AnnealingStrategy) Plan(
	jobs []*base.Job,
	machines []*base.Machine,
	startTime time.Time,
) (*base.Solution, base.MachineTimeSlots) {

	sim := strategy.NewFactorySimulator(jobs, machines, startTime)

	n := sim.TotalOperations()
	currentWeights := make([]float64, n)
	for i := range currentWeights {
		currentWeights[i] = rand.Float64()
	}

	currentRes := sim.Simulate(currentWeights)

	bestRes := currentRes
	temp := s.InitialTemp
	for temp > s.MinTemp {
		for i := 0; i < s.Iterations; i++ {
			nextWeights := s.mutate(currentWeights)

			nextRes := sim.Simulate(nextWeights)

			if s.shouldAccept(currentRes.Cost, nextRes.Cost, temp) {
				currentRes = nextRes

				if currentRes.Cost < bestRes.Cost {
					bestRes = currentRes
					fmt.Printf("[SA] New Best: %.2fs at Temp: %.2f\n", bestRes.Cost, temp)
				}
			}
		}
		temp *= s.Alpha
	}

	return bestRes.Solution, bestRes.MachineSlots
}

func (s *AnnealingStrategy) mutate(weights []float64) []float64 {
	next := make([]float64, len(weights))
	copy(next, weights)

	for i := 0; i < s.SwapsPerMutation; i++ {
		idx1 := rand.Intn(len(next))
		idx2 := rand.Intn(len(next))
		next[idx1], next[idx2] = next[idx2], next[idx1]
	}

	return next
}

func (s *AnnealingStrategy) shouldAccept(current, next, temp float64) bool {
	if next < current {
		return true
	}
	probability := math.Exp((current - next) / temp)
	return rand.Float64() < probability
}
