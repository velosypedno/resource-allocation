package base

import "fmt"

type (
	MachineType uint
	MachineID   int
)

type MachineTimeSlots map[MachineID][]Period
type MachineTypeIndex map[MachineType][]MachineID

const (
	None MachineType = iota
	Assembler
	Smeltery
	Extruder
	Sawmill
)

func (m MachineType) String() string {
	switch m {
	case None:
		return "None"
	case Assembler:
		return "Assembler"
	case Smeltery:
		return "Smeltery"
	case Extruder:
		return "Extruder"
	case Sawmill:
		return "Sawmill"
	default:
		return "Unknown"
	}
}

type Machine struct {
	ID   MachineID
	Type MachineType
}

func (m Machine) String() string {
	return fmt.Sprintf("ID: %d, Type: %s", m.ID, m.Type)
}

func NewMachine(id MachineID, machineType MachineType) Machine {
	return Machine{
		ID:   id,
		Type: machineType,
	}
}
