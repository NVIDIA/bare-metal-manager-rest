package managerapi

// InstanceTypeExpansion - InstanceType Expansion
type InstanceTypeExpansion interface{}

// InstanceTypeInterface - Interface for InstanceType
type InstanceTypeInterface interface {
	// List all the APIs for InstanceType here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	GetState() []string
	InstanceTypeExpansion
}
