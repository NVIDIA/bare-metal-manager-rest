package managerapi

// DpuExtensionServiceExpansion - DPU Extension Service Expansion
type DpuExtensionServiceExpansion interface{}

// DpuExtensionServiceInterface - Interface for DPU Extension Service
type DpuExtensionServiceInterface interface {
	// List all the APIs for DPU Extension Service here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	GetState() []string
	DpuExtensionServiceExpansion
}
