package base

import "time"

var BicycleAssemblingOperationTemplate = OperationTemplate{
	Name:           "Bicycle Assembly",
	MachineType:    Assembler,
	ProcessingTime: 60 * time.Minute,
	Children: []OperationTemplate{
		{
			Name:           "Frame Construction",
			MachineType:    Assembler,
			ProcessingTime: 45 * time.Minute,
			Children: []OperationTemplate{
				{Name: "Tube Smelting", MachineType: Smeltery, ProcessingTime: 30 * time.Minute},
				{Name: "Fork Extrusion", MachineType: Extruder, ProcessingTime: 20 * time.Minute},
			},
		},
		{
			Name:           "Wheel Set",
			MachineType:    Extruder,
			ProcessingTime: 25 * time.Minute,
		},
	},
}

var BicycleJobTemplate = JobTemplate{
	Name:       "Bicycle",
	Operations: []OperationTemplate{BicycleAssemblingOperationTemplate},
}

var ScooterAssemblingOperationTemplate = OperationTemplate{
	Name:           "Scooter Assembly",
	MachineType:    Assembler,
	ProcessingTime: 40 * time.Minute,
	Children: []OperationTemplate{
		{
			Name:           "Handlebar Unit",
			MachineType:    Extruder,
			ProcessingTime: 20 * time.Minute,
			Children: []OperationTemplate{
				{Name: "Metal Casting", MachineType: Smeltery, ProcessingTime: 25 * time.Minute},
			},
		},
		{Name: "Deck Sawing", MachineType: Sawmill, ProcessingTime: 15 * time.Minute},
	},
}

var ScooterJobTemplate = JobTemplate{
	Name:       "Scooter",
	Operations: []OperationTemplate{ScooterAssemblingOperationTemplate},
}

var SkateboardAssemblingOperationTemplate = OperationTemplate{
	Name:           "Skateboard Assembly",
	MachineType:    Assembler,
	ProcessingTime: 25 * time.Minute,
	Children: []OperationTemplate{
		{
			Name:           "Deck Lamination",
			MachineType:    Sawmill,
			ProcessingTime: 35 * time.Minute,
			Children: []OperationTemplate{
				{Name: "Wood Cutting", MachineType: Sawmill, ProcessingTime: 15 * time.Minute},
			},
		},
		{Name: "Truck Casting", MachineType: Smeltery, ProcessingTime: 20 * time.Minute},
	},
}

var SkateboardJobTemplate = JobTemplate{
	Name:       "Skateboard",
	Operations: []OperationTemplate{SkateboardAssemblingOperationTemplate},
}
