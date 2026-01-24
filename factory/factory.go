package factory

import (
	"time"

	"github.com/velosypedno/resource-allocation/base"
)

type PlannerStrategy interface {
	Plan([]*base.Job, []*base.Machine, time.Time) (base.Solution, base.MachineTimeSlots)
	Name() string
	Description() string
}

type Factory struct {
	Jobs     []*base.Job
	Machines []*base.Machine
	Planner  PlannerStrategy

	jobCounter     int
	machineCounter int
}

func (f *Factory) AddJob(template base.JobTemplate) {
	f.jobCounter++
	newJob := base.CreateJob(base.JobID(f.jobCounter), template)
	f.Jobs = append(f.Jobs, &newJob)
}

func (f *Factory) AddMachine(machineType base.MachineType) {
	f.machineCounter++
	machine := base.NewMachine(base.MachineID(f.machineCounter), machineType)
	f.Machines = append(f.Machines, &machine)
}

func (f *Factory) SetPlanner(planner PlannerStrategy) {
	f.Planner = planner
}

func (f *Factory) Plan(startTime time.Time) (base.Solution, SchedulingInfo, error) {
	startPlanning := time.Now()
	schedulingMetaInfo := SchedulingMetaInfo{
		StrategyName:        f.Planner.Name(),
		StrategyDescription: f.Planner.Description(),
	}

	if err := f.validate(); err != nil {
		schedulingMetaInfo.SchedulingTime = time.Since(startPlanning)
		return base.Solution{}, SchedulingInfo{}, err
	}

	solution, machineSlotsMap := f.Planner.Plan(f.Jobs, f.Machines, startTime)
	schedulingMetaInfo.SchedulingTime = time.Since(startPlanning)

	workflowPeriod := solution.GetWorkFlowPeriod()
	schedulingInfo := SchedulingInfo{
		SchedulingMetaInfo: schedulingMetaInfo,
		MakeSpan:           workflowPeriod.Duration(),
		UtilizationLevel:   machineSlotsMap.GetUtilizationLevel(workflowPeriod.Duration()),
	}

	return solution, schedulingInfo, nil
}

func (f *Factory) validate() error {
	neededMachineTypesMap := make(map[base.MachineType]struct{})

	for _, job := range f.Jobs {
		for _, machineType := range job.GetNeededMachineTypes() {
			neededMachineTypesMap[machineType] = struct{}{}
		}
	}

	availableMachineTypesMap := make(map[base.MachineType]struct{}, len(f.Machines))
	for _, machine := range f.Machines {
		availableMachineTypesMap[machine.Type] = struct{}{}
	}

	var missing []base.MachineType
	for machineType := range neededMachineTypesMap {
		if _, ok := availableMachineTypesMap[machineType]; !ok {
			missing = append(missing, machineType)
		}
	}

	if len(missing) > 0 {
		return &MissingMachinesError{MissingTypes: missing}
	}

	return nil
}

type SchedulingMetaInfo struct {
	StrategyName        string
	StrategyDescription string
	SchedulingTime      time.Duration
}

type SchedulingInfo struct {
	SchedulingMetaInfo
	MakeSpan         time.Duration
	UtilizationLevel float64
}
