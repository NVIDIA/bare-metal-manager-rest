package carbide

import (
	Manager "github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/managerapi"
	"github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/elektratypes"
)

// ManagerAccess - access to all managers
var ManagerAccess *Manager.ManagerAccess

// API - all API interface
//
//nolint:all
type API struct{} //nolint:all

// NewCarbideManager - returns a new instance of helm manager
func NewCarbideManager(superForge *elektratypes.Elektra, superAPI *Manager.ManagerAPI, superConf *Manager.ManagerConf) *API {
	ManagerAccess = &Manager.ManagerAccess{
		Data: &Manager.ManagerData{
			EB: superForge,
		},
		API:  superAPI,
		Conf: superConf,
	}
	return &API{}
}
