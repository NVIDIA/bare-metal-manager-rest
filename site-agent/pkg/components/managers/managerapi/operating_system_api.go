package managerapi

// OperatingSystemExpansion - Operating System Expansion
type OperatingSystemExpansion interface{}

// OperatingSystemInterface - Interface for Operating System
type OperatingSystemInterface interface {
	// List all the APIs for Operating System here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	GetState() []string
	OperatingSystemExpansion
}
