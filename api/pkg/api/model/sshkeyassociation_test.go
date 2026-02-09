package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

func TestNewAPISSHKeyAssociation(t *testing.T) {
	ska := cdbm.SSHKeyAssociation{
		ID:            uuid.New(),
		SSHKeyID:      uuid.New(),
		SSHKeyGroupID: uuid.New(),
		Created:       time.Now(),
		Updated:       time.Now(),
	}
	apiska := NewAPISSHKeyAssociation(&ska)
	assert.Equal(t, apiska.ID, ska.ID.String())
	assert.Equal(t, apiska.SSHKeyID, ska.SSHKeyID.String())
	assert.Equal(t, apiska.SSHKeyGroupID, ska.SSHKeyGroupID.String())
}
