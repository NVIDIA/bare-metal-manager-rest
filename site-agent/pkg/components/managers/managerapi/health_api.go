package managerapi

// HealthExpansion - Health Expansion
type HealthExpansion interface{}

// HealthInterface - interface to Health
type HealthInterface interface {
	// List all the apis of Health here
	Init()
	// RegisterSubscriber() error
	// RegisterPublisher() error
	HealthExpansion
}
