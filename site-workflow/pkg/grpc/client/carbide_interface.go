package client

// CarbideInterface is the interface for the Carbide client
type CarbideInterface interface {
	ComputeGetter
	NetworkGetter
	StorageGetter
}
