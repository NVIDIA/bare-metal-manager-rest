package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestNewAPISSHKeyGroupSiteAssociation(t *testing.T) {
	skgsa := cdbm.SSHKeyGroupSiteAssociation{
		ID:            uuid.New(),
		SSHKeyGroupID: uuid.New(),
		SiteID:        uuid.New(),
		Version:       db.GetStrPtr("1234"),
		Status:        cdbm.SSHKeyGroupSiteAssociationStatusSyncing,
		Created:       time.Now(),
		Updated:       time.Now(),
	}
	apiskgsa := NewAPISSHKeyGroupSiteAssociation(&skgsa, nil)
	assert.Equal(t, apiskgsa.ControllerKeySetVersion, skgsa.Version)
	assert.Equal(t, apiskgsa.Status, skgsa.Status)
}
