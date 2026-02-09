package managerapi

// MachineValidationExpansion - MachineValidation Expansion
type MachineValidationExpansion interface{}

// MachineValidationInterface - Interface for MachineValidation
type MachineValidationInterface interface {
	// List all the APIs for MachineValidation here
	Init()
	RegisterSubscriber() error
	GetState() []string
	MachineValidationExpansion
}
