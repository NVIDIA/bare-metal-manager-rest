package managerapi

// SKUExpansion - SKU Expansion
type SKUExpansion interface{}

// SKUInterface - interface to SKU
type SKUInterface interface {
	// List all the apis of SKU here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	GetState() []string
	SKUExpansion
}
