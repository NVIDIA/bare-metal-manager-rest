package managerapi

// NetworkSecurityGroupExpansion - NetworkSecurityGroup Expansion
type NetworkSecurityGroupExpansion interface{}

// NetworkSecurityGroupInterface - Interface for NetworkSecurityGroup
type NetworkSecurityGroupInterface interface {
	// List all the APIs for NetworkSecurityGroup here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	GetState() []string
	NetworkSecurityGroupExpansion
}
