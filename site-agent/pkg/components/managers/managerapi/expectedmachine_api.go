package managerapi

// ExpectedMachineExpansion - ExpectedMachine Expansion
type ExpectedMachineExpansion interface{}

// ExpectedMachineInterface - interface to ExpectedMachine
type ExpectedMachineInterface interface {
	// List all the apis of ExpectedMachine here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	GetState() []string
	ExpectedMachineExpansion
}
