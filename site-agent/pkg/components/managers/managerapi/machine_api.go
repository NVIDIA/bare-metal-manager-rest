package managerapi

// MachineExpansion - Machine Expansion
type MachineExpansion interface{}

// MachineInterface - interface to Machine
type MachineInterface interface {
	// List all the apis of Machine here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error

	GetState() []string

	MachineExpansion
}
