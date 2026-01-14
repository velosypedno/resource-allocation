package naive

import (
	"errors"
	"sort"
	"time"

	"github.com/velosypedno/resource-allocation/base"
)

const name = "Greedy"
const description = `Greedy Earliest Completion Time scheduling. Each operation is assigned to the machine that 
provides the earliest completion time, taking into account the technological sequence 
(dependence on child operations) and already occupied time slots.`

type Strategy struct{}

func (Strategy) Name() string {
	return name
}

func (Strategy) Description() string {
	return description
}

func (s *Strategy) Plan(
	jobs []*base.Job,
	machines []*base.Machine,
	startTime time.Time,
) (base.Solution, base.MachineTimeSlots) {
	occupiedMap := InitTimeSlotsMap(machines)
	machineTypeIndex := InitMachineTypeIndex(machines)

	solution := base.Solution{}

	for _, job := range jobs {
		jobSolution := PlanJob(job, machines, startTime, occupiedMap, machineTypeIndex)
		solution.Jobs = append(solution.Jobs, jobSolution)
	}

	return solution, occupiedMap
}

func InitTimeSlotsMap(machines []*base.Machine) base.MachineTimeSlots {
	timeSlotsMap := make(map[base.MachineID][]base.Period)
	for _, machine := range machines {
		timeSlotsMap[machine.ID] = []base.Period{}
	}
	return timeSlotsMap
}

func InitMachineTypeIndex(machines []*base.Machine) base.MachineTypeIndex {
	machineTypeIndex := make(map[base.MachineType][]base.MachineID)
	for _, machine := range machines {
		machineTypeIndex[machine.Type] = append(machineTypeIndex[machine.Type], machine.ID)
	}
	return machineTypeIndex
}

func PlanJob(
	job *base.Job,
	machines []*base.Machine,
	startTime time.Time,
	occupiedMap base.MachineTimeSlots,
	machineTypeIndex base.MachineTypeIndex,
) *base.JobSolution {
	jobSolution := &base.JobSolution{
		Job:                job,
		OperationSolutions: []*base.OperationSolution{},
	}
	for _, operation := range job.Operations {
		operationSolution := PlanOperation(operation, machines, startTime, occupiedMap, machineTypeIndex)
		jobSolution.OperationSolutions = append(jobSolution.OperationSolutions, operationSolution)
	}

	return jobSolution
}

func PlanOperation(
	operation *base.Operation,
	machines []*base.Machine,
	startTime time.Time,
	occupiedMap base.MachineTimeSlots,
	machineTypeIndex base.MachineTypeIndex,
) *base.OperationSolution {
	operationSolution := &base.OperationSolution{
		Operation:      operation,
		ChildSolutions: []*base.OperationSolution{},
	}

	for _, child := range operation.ChildOperations {
		operationSolution.ChildSolutions = append(
			operationSolution.ChildSolutions,
			PlanOperation(child, machines, startTime, occupiedMap, machineTypeIndex))
	}
	lastChildEndTime, err := operationSolution.GetLastChildCompletionTime()
	if errors.Is(err, base.ErrNoChildrenFound) {
		targetMachineID, targetPeriod := FindBestSlot(
			startTime,
			operation.Duration,
			operation.MachineType,
			occupiedMap,
			machineTypeIndex,
		)
		operationSolution.MachineID = targetMachineID
		operationSolution.Period = targetPeriod
		occupiedMap[targetMachineID] = append(occupiedMap[targetMachineID], targetPeriod)
		return operationSolution
	}
	if err != nil {
		panic(err)
	}

	targetMachineID, targetPeriod := FindBestSlot(
		lastChildEndTime,
		operation.Duration,
		operation.MachineType,
		occupiedMap,
		machineTypeIndex,
	)
	operationSolution.MachineID = targetMachineID
	operationSolution.Period = targetPeriod
	occupiedMap[targetMachineID] = append(occupiedMap[targetMachineID], targetPeriod)

	return operationSolution

}

func FindBestSlot(
	startTime time.Time,
	duration time.Duration,
	machineType base.MachineType,
	occupiedMap base.MachineTimeSlots,
	machineTypeIndex base.MachineTypeIndex,
) (base.MachineID, base.Period) {
	targetMachineIDs := machineTypeIndex[machineType]

	var bestMachineID base.MachineID
	var bestPeriod base.Period
	firstFound := false

	for _, mID := range targetMachineIDs {
		currentPeriod := findEarliestGap(startTime, duration, occupiedMap[mID])

		if !firstFound || currentPeriod.End.Before(bestPeriod.End) {
			bestPeriod = currentPeriod
			bestMachineID = mID
			firstFound = true
		}
	}

	return bestMachineID, bestPeriod
}

func findEarliestGap(startTime time.Time, duration time.Duration, occupied []base.Period) base.Period {
	sort.Slice(occupied, func(i, j int) bool {
		return occupied[i].Start.Before(occupied[j].Start)
	})

	candidateStart := startTime

	for _, slot := range occupied {
		if slot.End.Before(candidateStart) {
			continue
		}
		if slot.Start.Sub(candidateStart) >= duration {
			return base.Period{
				Start: candidateStart,
				End:   candidateStart.Add(duration),
			}
		}

		if slot.End.After(candidateStart) {
			candidateStart = slot.End
		}
	}

	return base.Period{
		Start: candidateStart,
		End:   candidateStart.Add(duration),
	}
}
