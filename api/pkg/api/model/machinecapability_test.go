package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestMachineCapability_NewAPIMachineCapability(t *testing.T) {
	dbmc := &cdbm.MachineCapability{
		Type:      cdbm.MachineCapabilityTypeCPU,
		Name:      "AMD Opteron Series x10",
		Frequency: cdb.GetStrPtr("3.0GHz"),
		Capacity:  cdb.GetStrPtr("3.0GHz"),
		Vendor:    cdb.GetStrPtr("AMD"),
		Count:     cdb.GetIntPtr(2),
	}

	apimc := NewAPIMachineCapability(dbmc)
	assert.Equal(t, dbmc.Type, apimc.Type)
	assert.Equal(t, dbmc.Name, apimc.Name)
	assert.Equal(t, *dbmc.Frequency, *apimc.Frequency)
	assert.Equal(t, *dbmc.Capacity, *apimc.Capacity)
	assert.Equal(t, *dbmc.Vendor, *apimc.Vendor)
	assert.Equal(t, *dbmc.Count, *apimc.Count)
}
