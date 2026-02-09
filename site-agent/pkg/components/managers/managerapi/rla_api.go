package managerapi

import (
	"context"

	"github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
)

// RLAExpansion - RLA Expansion
type RLAExpansion interface{}

// RLAInterface - interface to RLA
type RLAInterface interface {
	// List all the apis of RLA here
	Init()
	Start()
	CreateGrpcClient() error
	GetGrpcClient() *client.RlaClient
	UpdateGrpcClientState(err error)
	CreateGrpcClientActivity(ctx context.Context, ResourceID string) (client *client.RlaClient, err error)
	RegisterGrpc()
	GetState() []string
	GetGrpcClientVersion() int64
	RLAExpansion
}
