package main

import (
	"fmt"
	"os"
	"time"

	"github.com/velosypedno/resource-allocation/chart"
	"github.com/velosypedno/resource-allocation/factory"
	"github.com/velosypedno/resource-allocation/parser"
	"github.com/velosypedno/resource-allocation/strategy/naive"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <factory_config_path> <orders_path>")
		os.Exit(1)
	}

	factoryConfigPath := os.Args[1]
	ordersPath := os.Args[2]

	machinesConfig, templates, err := parser.ParseFactoryConfig(factoryConfigPath)
	if err != nil {
		fmt.Printf("Error parsing factory config: %v\n", err)
		os.Exit(1)
	}

	f := &factory.Factory{}
	f.Configure(machinesConfig, templates)
	f.SetPlanner(&naive.Strategy{})

	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)

	orders, err := parser.ParseOrders(ordersPath)
	if err != nil {
		fmt.Printf("Error parsing orders: %v\n", err)
		os.Exit(1)
	}

	solution, metaInfo, err := f.Plan(orders, startTime)
	if err != nil {
		fmt.Printf("Error during planning: %v\n", err)
		os.Exit(1)
	}

	solutionChart := chart.GenerateFromSolution(solution, f.Machines, metaInfo)

	outputFile, err := os.Create("bar.html")
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	err = solutionChart.Render(outputFile)
	if err != nil {
		fmt.Printf("Error rendering chart: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated bar.html")
}
