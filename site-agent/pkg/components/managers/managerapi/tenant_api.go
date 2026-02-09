package managerapi

// TenantExpansion - Tenant Expansion
type TenantExpansion interface{}

// TenantInterface - Interface for Tenant
type TenantInterface interface {
	// List all the APIs for Tenant here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	GetState() []string
	TenantExpansion
}
