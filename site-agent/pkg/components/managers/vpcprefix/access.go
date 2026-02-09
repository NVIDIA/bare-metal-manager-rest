package vpcprefix

import (
	Manager "github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/managerapi"
	"github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/elektratypes"
)

// ManagerAccess - access to all managers
var ManagerAccess *Manager.ManagerAccess

// API - all API interface
type API struct{}

// NewVpcPrefixManager - returns a new instance of helm manager
func NewVpcPrefixManager(superForge *elektratypes.Elektra, superAPI *Manager.ManagerAPI, superConf *Manager.ManagerConf) *API {
	ManagerAccess = &Manager.ManagerAccess{
		Data: &Manager.ManagerData{
			EB: superForge,
		},
		API:  superAPI,
		Conf: superConf,
	}
	return &API{}
}
