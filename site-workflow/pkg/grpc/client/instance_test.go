package client

import (
	"context"
	"testing"

	"github.com/gogo/status"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

func TestInstance_DeleteInstance(t *testing.T) {
	mockCarbide := NewMockCarbideClient()

	type fields struct {
		CarbideClient *CarbideClient
	}
	type args struct {
		ctx     context.Context
		request *wflows.DeleteInstanceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test delete instance success",
			fields: fields{
				CarbideClient: mockCarbide,
			},
			args: args{
				ctx: context.Background(),
				request: &wflows.DeleteInstanceRequest{
					InstanceId: &wflows.UUID{Value: uuid.New().String()},
				},
			},
			wantErr: false,
		},
		{
			name: "test delete instance failed, NotFound",
			fields: fields{
				CarbideClient: mockCarbide,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "wantError", status.Error(codes.NotFound, "instance not found: ")),
				request: &wflows.DeleteInstanceRequest{
					InstanceId: &wflows.UUID{Value: uuid.New().String()},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := &compute{
				carbide: tt.fields.CarbideClient.carbide,
			}
			_, err := cc.DeleteInstance(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInstance_CreateInstance(t *testing.T) {
	mockCarbide := NewMockCarbideClient()

	type fields struct {
		CarbideClient *CarbideClient
	}
	type args struct {
		ctx     context.Context
		request *wflows.CreateInstanceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test create instance success",
			fields: fields{
				CarbideClient: mockCarbide,
			},
			args: args{
				ctx: context.Background(),
				request: &wflows.CreateInstanceRequest{
					InstanceId:       &wflows.UUID{Value: uuid.New().String()},
					MachineId:        &wflows.MachineId{Id: uuid.New().String()},
					TenantOrg:        "testOrg",
					PhoneHomeEnabled: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := &compute{
				carbide: tt.fields.CarbideClient.carbide,
			}
			_, err := cc.CreateInstance(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
