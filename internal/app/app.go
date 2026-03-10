package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/chart"
	"github.com/velosypedno/resource-allocation/internal/parser"
	"github.com/velosypedno/resource-allocation/internal/scheduler"
)

type App struct {
	Scheduler *scheduler.Scheduler
}

func New(machinesConfig []parser.MachineConfig, templates []base.JobTemplate, strategies []parser.Strategy) *App {
	s := &scheduler.Scheduler{}
	s.Configure(machinesConfig, templates)
	s.SetPlanners(strategies...)

	return &App{
		Scheduler: s,
	}
}

func (a *App) Run(startTime time.Time, orders []parser.OrderDTO, customName string) error {
	results, err := a.Scheduler.Plan(orders, startTime)
	if err != nil {
		return fmt.Errorf("during planning: %v", err)

	}

	solutionsChart := chart.GenerateFromSolutions(results, a.Scheduler.Machines)

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
