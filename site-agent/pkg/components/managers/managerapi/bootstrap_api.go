package managerapi

// BootstrapExpansion - Bootstrap Expansion
type BootstrapExpansion interface{}

// BootstrapInterface - interface to Bootstrap
type BootstrapInterface interface {
	// List all the apis of Bootstrap here
	Init()
	Start()
	DownloadAndStoreCreds(otpOverride []byte) error
	GetState() []string

	BootstrapExpansion

	// Temporal Workflows - Subscriber
	RegisterSubscriber() error
}
