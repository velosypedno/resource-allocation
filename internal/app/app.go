package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/chart"
	"github.com/velosypedno/resource-allocation/internal/factory"
	"github.com/velosypedno/resource-allocation/internal/parser"
)

type App struct {
	Factory *factory.Factory
}

func New(machinesConfig []parser.MachineConfig, templates []base.JobTemplate, strategies []parser.Strategy) *App {
	f := &factory.Factory{}
	f.Configure(machinesConfig, templates)
	f.SetPlanners(strategies...)

	return &App{
		Factory: f,
	}
}

func (a *App) Run(startTime time.Time, orders []parser.OrderDTO, customName string) error {
	results, err := a.Factory.Plan(orders, startTime)
	if err != nil {
		return fmt.Errorf("during planning: %v", err)

	}

	solutionsChart := chart.GenerateFromSolutions(results, a.Factory.Machines)

	outputDir := "results"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating directory: %v", err)

	}

	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("plan_%s%s.html", timestamp, customName)
	fullPath := filepath.Join(outputDir, fileName)

	outputFile, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("creating output file: %v", err)
	}
	defer outputFile.Close()

	err = solutionsChart.Render(outputFile)
	if err != nil {
		return fmt.Errorf("error rendering chart: %v", err)
	}

	fmt.Printf("Successfully generated %s\n", fullPath)
	return nil
}
