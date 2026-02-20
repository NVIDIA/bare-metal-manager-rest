/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/nvidia/bare-metal-manager-rest/api/internal/config"
	cam "github.com/nvidia/bare-metal-manager-rest/api/pkg/api/model"
	cerr "github.com/nvidia/bare-metal-manager-rest/common/pkg/util"
	cdb "github.com/nvidia/bare-metal-manager-rest/db/pkg/db"
	cdbm "github.com/nvidia/bare-metal-manager-rest/db/pkg/db/model"
	cdbp "github.com/nvidia/bare-metal-manager-rest/db/pkg/db/paginator"
	cwssaws "github.com/nvidia/bare-metal-manager-rest/workflow-schema/schema/site-agent/workflows/v1"
)

// InstanceOSConfigProvider abstracts OS-related fields shared between single and batch
// instance create requests, allowing a single BuildOsConfig implementation.
type InstanceOSConfigProvider interface {
	GetOperatingSystemID() *string
	GetTenantID() string
	ValidateAndSetOperatingSystemData(cfg *config.Config, os *cdbm.OperatingSystem) error
	GetAlwaysBootWithCustomIpxe() *bool
	GetPhoneHomeEnabled() *bool
	GetIpxeScript() *string
	GetUserData() *string
}

// InstanceCreateContext accumulates validated state during instance creation.
// Both single and batch instance create handlers use this to avoid duplicating validation logic.
type InstanceCreateContext struct {
	DBSession *cdb.Session
	Cfg       *config.Config
	Logger    *zerolog.Logger

	// Validated entities — populated progressively by validation methods.
	Tenant         *cdbm.Tenant
	VPC            *cdbm.Vpc
	Site           *cdbm.Site
	DefaultNvllpID *uuid.UUID

	// Interface validation products
	DBInterfaces        []cdbm.Interface
	SubnetIDMap         map[uuid.UUID]*cdbm.Subnet
	VpcPrefixIDMap      map[uuid.UUID]*cdbm.VpcPrefix
	IsDeviceInfoPresent bool

	// DPU Extension Service validation products
	DESIDMap map[string]*cdbm.DpuExtensionService

	// SSH Key Group validation products
	SSHKeyGroups []cdbm.SSHKeyGroup

	// Machine Capability validation products
	DBIBInterfaces  []cdbm.InfiniBandInterface
	DBNVLInterfaces []cdbm.NVLinkInterface

	// OS Config products
	OsConfig *cwssaws.OperatingSystem
	OsID     *uuid.UUID
}

// NewInstanceCreateContext initializes a new context for shared instance creation logic.
func NewInstanceCreateContext(dbSession *cdb.Session, cfg *config.Config, logger *zerolog.Logger) *InstanceCreateContext {
	return &InstanceCreateContext{
		DBSession: dbSession,
		Cfg:       cfg,
		Logger:    logger,
	}
}

// ValidateTenantAndVPC validates tenant ownership, VPC state, site readiness, and populates
// cc.Tenant, cc.VPC, cc.Site, and cc.DefaultNvllpID.
func (cc *InstanceCreateContext) ValidateTenantAndVPC(ctx context.Context, org, tenantID, vpcID string) *cerr.APIError {
	logger := cc.Logger

	tenant, err := GetTenantForOrg(ctx, nil, cc.DBSession, org)
	if err != nil {
		if err == ErrOrgTenantNotFound {
			logger.Warn().Err(err).Msg("Org does not have a Tenant associated")
			return cerr.NewAPIError(http.StatusBadRequest, "Org does not have a Tenant associated", nil)
		}
		logger.Error().Err(err).Msg("unable to retrieve tenant for org")
		return cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve tenant for org", nil)
	}

	apiTenant, err := GetTenantFromIDString(ctx, nil, tenantID, cc.DBSession)
	if err != nil {
		logger.Warn().Err(err).Msg("error retrieving tenant from request")
		return cerr.NewAPIError(http.StatusBadRequest, "TenantID in request is not valid", nil)
	}
	if apiTenant.ID != tenant.ID {
		logger.Warn().Msg("tenant id in request does not match tenant in org")
		return cerr.NewAPIError(http.StatusBadRequest, "TenantID in request does not match tenant in org", nil)
	}

	vpc, err := GetVpcFromIDString(ctx, nil, vpcID, []string{cdbm.NVLinkLogicalPartitionRelationName}, cc.DBSession)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIError(http.StatusBadRequest, "Could not find VPC with ID specified in request data", nil)
		}
		logger.Warn().Err(err).Str("vpcId", vpcID).Msg("error retrieving VPC from request")
		return cerr.NewAPIError(http.StatusBadRequest, "VpcID in request is not valid", nil)
	}

	if vpc.TenantID != tenant.ID {
		logger.Warn().Msg("tenant id in request does not match tenant in VPC")
		return cerr.NewAPIError(http.StatusBadRequest, "VPC specified in request is not owned by Tenant", nil)
	}

	if vpc.ControllerVpcID == nil || vpc.Status != cdbm.VpcStatusReady {
		logger.Warn().Msg("VPC specified in request data is not ready")
		return cerr.NewAPIError(http.StatusBadRequest, "VPC specified in request data is not ready", nil)
	}

	var defaultNvllpID *uuid.UUID
	if vpc.NVLinkLogicalPartitionID != nil {
		defaultNvllpID = vpc.NVLinkLogicalPartitionID
	}

	siteDAO := cdbm.NewSiteDAO(cc.DBSession)
	site, err := siteDAO.GetByID(ctx, nil, vpc.SiteID, nil, false)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIError(http.StatusBadRequest, "The Site where Instances are being created could not be found", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB by ID")
		return cerr.NewAPIError(http.StatusInternalServerError, "The Site where Instances are being created could not be retrieved", nil)
	}

	if site.Status != cdbm.SiteStatusRegistered {
		logger.Warn().Str("Site ID", site.ID.String()).Str("Site Status", site.Status).
			Msg("The Site where Instances are being created is not in Registered state")
		return cerr.NewAPIError(http.StatusBadRequest, "The Site where Instances are being created is not in Registered state", nil)
	}

	cc.Tenant = tenant
	cc.VPC = vpc
	cc.Site = site
	cc.DefaultNvllpID = defaultNvllpID
	return nil
}

// ValidateInterfaces validates subnet and VPC prefix interfaces, populating
// cc.DBInterfaces, cc.SubnetIDMap, cc.VpcPrefixIDMap, and cc.IsDeviceInfoPresent.
// Requires cc.Tenant, cc.VPC, cc.Site to be populated.
func (cc *InstanceCreateContext) ValidateInterfaces(ctx context.Context, interfaces []cam.APIInterfaceCreateOrUpdateRequest) *cerr.APIError {
	logger := cc.Logger
	tenant := cc.Tenant
	vpc := cc.VPC

	subnetDAO := cdbm.NewSubnetDAO(cc.DBSession)
	vpDAO := cdbm.NewVpcPrefixDAO(cc.DBSession)

	subnetIDs := []uuid.UUID{}
	vpcPrefixIDs := []uuid.UUID{}

	for _, ifc := range interfaces {
		if ifc.SubnetID != nil {
			subnetID, err := uuid.Parse(*ifc.SubnetID)
			if err != nil {
				logger.Error().Err(err).Str("subnetID", *ifc.SubnetID).Msg("error parsing subnet id")
				return cerr.NewAPIError(http.StatusBadRequest, "Invalid Subnet ID format", nil)
			}
			subnetIDs = append(subnetIDs, subnetID)
		}
		if ifc.VpcPrefixID != nil {
			vpcPrefixID, err := uuid.Parse(*ifc.VpcPrefixID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing vpcprefix id in instance vpcprefix request")
				return cerr.NewAPIError(http.StatusBadRequest, "VPC Prefix ID specified in request data is not valid", nil)
			}
			vpcPrefixIDs = append(vpcPrefixIDs, vpcPrefixID)
		}
	}

	subnetIDMap := make(map[uuid.UUID]*cdbm.Subnet)
	if len(subnetIDs) > 0 {
		subnets, _, err := subnetDAO.GetAll(ctx, nil, cdbm.SubnetFilterInput{SubnetIDs: subnetIDs}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Subnets from DB by IDs")
			return cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve Subnets from DB by IDs", nil)
		}
		for i := range subnets {
			subnetIDMap[subnets[i].ID] = &subnets[i]
		}
	}

	vpcPrefixIDMap := make(map[uuid.UUID]*cdbm.VpcPrefix)
	if len(vpcPrefixIDs) > 0 {
		vpcPrefixes, _, err := vpDAO.GetAll(ctx, nil, cdbm.VpcPrefixFilterInput{VpcPrefixIDs: vpcPrefixIDs}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving VPC Prefixes from DB by IDs")
			return cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve VPC Prefixes from DB by IDs", nil)
		}
		for i := range vpcPrefixes {
			vpcPrefixIDMap[vpcPrefixes[i].ID] = &vpcPrefixes[i]
		}
	}

	dbInterfaces := []cdbm.Interface{}
	isDeviceInfoPresent := false

	for _, ifc := range interfaces {
		if ifc.SubnetID != nil {
			subnetID := uuid.MustParse(*ifc.SubnetID)

			subnet, ok := subnetIDMap[subnetID]
			if !ok {
				return cerr.NewAPIError(http.StatusBadRequest, "Could not find Subnet with ID specified in request data", nil)
			}

			if subnet.TenantID != tenant.ID {
				logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request is not owned by Tenant", subnetID))
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request is not owned by Tenant", subnetID), nil)
			}

			if subnet.ControllerNetworkSegmentID == nil || subnet.Status != cdbm.SubnetStatusReady {
				logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request data is not in Ready state", subnetID))
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request data is not in Ready state", subnetID), nil)
			}

			if subnet.VpcID != vpc.ID {
				logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request does not match with VPC", subnetID))
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request does not match with VPC", subnetID), nil)
			}

			if vpc.NetworkVirtualizationType != nil && *vpc.NetworkVirtualizationType != cdbm.VpcEthernetVirtualizer {
				logger.Warn().Msg(fmt.Sprintf("VPC: %v specified in request must have Ethernet network virtualization type in order to create Subnet based interfaces", vpc.ID))
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC: %v specified in request must have Ethernet network virtualization type in order to create Subnet based interfaces", vpc.ID), nil)
			}

			dbInterfaces = append(dbInterfaces, cdbm.Interface{
				SubnetID:   &subnetID,
				IsPhysical: ifc.IsPhysical,
				Status:     cdbm.InterfaceStatusPending,
			})
		}

		if ifc.VpcPrefixID != nil {
			vpcPrefixUUID := uuid.MustParse(*ifc.VpcPrefixID)

			vpcPrefix, ok := vpcPrefixIDMap[vpcPrefixUUID]
			if !ok {
				return cerr.NewAPIError(http.StatusBadRequest, "Could not find VPC Prefix with ID specified in request data", nil)
			}

			if vpcPrefix.TenantID != tenant.ID {
				logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request is not owned by Tenant", vpcPrefixUUID))
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request is not owned by Tenant", vpcPrefixUUID), nil)
			}

			if vpcPrefix.Status != cdbm.VpcPrefixStatusReady {
				logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request data is not in Ready state", vpcPrefixUUID))
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request data is not in Ready state", vpcPrefixUUID), nil)
			}

			if vpcPrefix.VpcID != vpc.ID {
				logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request does not match with VPC", vpcPrefixUUID))
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request does not match with VPC", vpcPrefixUUID), nil)
			}

			if vpc.NetworkVirtualizationType == nil || *vpc.NetworkVirtualizationType != cdbm.VpcFNN {
				logger.Warn().Msg(fmt.Sprintf("VPC: %v specified in request must have FNN network virtualization type in order to create VPC Prefix based interfaces", vpc.ID))
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC: %v specified in request must have FNN network virtualization type in order to create VPC Prefix based interfaces", vpc.ID), nil)
			}

			if ifc.Device != nil && ifc.DeviceInstance != nil {
				isDeviceInfoPresent = true
			}

			dbInterfaces = append(dbInterfaces, cdbm.Interface{
				VpcPrefixID:       &vpcPrefixUUID,
				Device:            ifc.Device,
				DeviceInstance:    ifc.DeviceInstance,
				VirtualFunctionID: ifc.VirtualFunctionID,
				IsPhysical:        ifc.IsPhysical,
				Status:            cdbm.InterfaceStatusPending,
			})
		}
	}

	cc.DBInterfaces = dbInterfaces
	cc.SubnetIDMap = subnetIDMap
	cc.VpcPrefixIDMap = vpcPrefixIDMap
	cc.IsDeviceInfoPresent = isDeviceInfoPresent
	return nil
}

// ValidateDPUExtensionServices validates DPU extension service deployments, populating cc.DESIDMap.
// Requires cc.Tenant and cc.Site to be populated.
func (cc *InstanceCreateContext) ValidateDPUExtensionServices(ctx context.Context, deployments []cam.APIDpuExtensionServiceDeploymentRequest) *cerr.APIError {
	logger := cc.Logger
	tenant := cc.Tenant
	site := cc.Site

	cc.DESIDMap = map[string]*cdbm.DpuExtensionService{}

	if len(deployments) == 0 {
		return nil
	}

	desIDs := make([]uuid.UUID, 0, len(deployments))
	uniqueDesIDs := make([]uuid.UUID, 0, len(deployments))
	seenDesIDs := make(map[uuid.UUID]bool, len(deployments))
	for _, adesdr := range deployments {
		desID, err := uuid.Parse(adesdr.DpuExtensionServiceID)
		if err != nil {
			logger.Warn().Err(err).Str("serviceID", adesdr.DpuExtensionServiceID).
				Msg("error parsing DPU Extension Service ID")
			return cerr.NewAPIError(http.StatusBadRequest,
				fmt.Sprintf("Invalid DPU Extension Service ID: %s", adesdr.DpuExtensionServiceID), nil)
		}
		desIDs = append(desIDs, desID)
		if !seenDesIDs[desID] {
			seenDesIDs[desID] = true
			uniqueDesIDs = append(uniqueDesIDs, desID)
		}
	}

	desDAO := cdbm.NewDpuExtensionServiceDAO(cc.DBSession)
	desList, _, err := desDAO.GetAll(ctx, nil, cdbm.DpuExtensionServiceFilterInput{
		DpuExtensionServiceIDs: uniqueDesIDs,
	}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving DPU Extension Services from DB")
		return cerr.NewAPIError(http.StatusInternalServerError,
			"Failed to retrieve DPU Extension Services specified in request, DB error", nil)
	}

	desMap := make(map[uuid.UUID]*cdbm.DpuExtensionService, len(desList))
	for i := range desList {
		desMap[desList[i].ID] = &desList[i]
	}

	for i, desID := range desIDs {
		des, exists := desMap[desID]
		if !exists {
			return cerr.NewAPIError(http.StatusBadRequest,
				fmt.Sprintf("Could not find DPU Extension Service with ID: %s", desID), nil)
		}

		if des.TenantID != tenant.ID {
			logger.Warn().Str("tenantID", tenant.ID.String()).Str("serviceID", desID.String()).
				Msg("DPU Extension Service does not belong to current Tenant")
			return cerr.NewAPIError(http.StatusForbidden,
				fmt.Sprintf("DPU Extension Service: %s does not belong to current Tenant", desID.String()), nil)
		}

		if des.SiteID != site.ID {
			logger.Warn().Str("siteID", site.ID.String()).Str("serviceID", desID.String()).
				Msg("DPU Extension Service does not belong to Site")
			return cerr.NewAPIError(http.StatusForbidden,
				fmt.Sprintf("DPU Extension Service: %s does not belong to Site where Instances are being created", desID.String()), nil)
		}

		versionFound := false
		requestedVersion := deployments[i].Version
		for _, version := range des.ActiveVersions {
			if version == requestedVersion {
				versionFound = true
				break
			}
		}
		if !versionFound {
			return cerr.NewAPIError(http.StatusBadRequest,
				fmt.Sprintf("Version: %s was not found for DPU Extension Service: %s", requestedVersion, desID.String()), nil)
		}

		cc.DESIDMap[desID.String()] = des
	}

	return nil
}

// ValidateNSG validates the network security group, if specified.
// Requires cc.Tenant and cc.Site to be populated.
func (cc *InstanceCreateContext) ValidateNSG(ctx context.Context, nsgID *string) *cerr.APIError {
	if nsgID == nil {
		return nil
	}

	logger := cc.Logger

	nsgDAO := cdbm.NewNetworkSecurityGroupDAO(cc.DBSession)
	nsg, err := nsgDAO.GetByID(ctx, nil, *nsgID, nil)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIError(http.StatusBadRequest,
				fmt.Sprintf("Could not find Network Security Group with ID: %s", *nsgID), nil)
		}
		logger.Error().Err(err).Str("nsgID", *nsgID).Msg("error retrieving Network Security Group from DB")
		return cerr.NewAPIError(http.StatusInternalServerError,
			"Failed to retrieve Network Security Group specified in request, DB error", nil)
	}

	if nsg.SiteID != cc.Site.ID {
		logger.Error().Str("siteID", cc.Site.ID.String()).Str("nsgID", *nsgID).
			Msg("Network Security Group does not belong to Site")
		return cerr.NewAPIError(http.StatusForbidden,
			"Network Security Group specified in request does not belong to Site", nil)
	}

	if nsg.TenantID != cc.Tenant.ID {
		logger.Error().Str("tenantID", cc.Tenant.ID.String()).Str("nsgID", *nsgID).
			Msg("Network Security Group does not belong to Tenant")
		return cerr.NewAPIError(http.StatusForbidden,
			"Network Security Group specified in request does not belong to Tenant", nil)
	}

	return nil
}

// ValidateSSHKeyGroups validates SSH key groups and their site associations, populating cc.SSHKeyGroups.
// Requires cc.Tenant and cc.Site to be populated.
func (cc *InstanceCreateContext) ValidateSSHKeyGroups(ctx context.Context, sshKeyGroupIDs []string) *cerr.APIError {
	logger := cc.Logger
	tenant := cc.Tenant
	site := cc.Site

	cc.SSHKeyGroups = []cdbm.SSHKeyGroup{}

	skgsaDAO := cdbm.NewSSHKeyGroupSiteAssociationDAO(cc.DBSession)

	for _, skgIDStr := range sshKeyGroupIDs {
		sshkeygroup, serr := GetSSHKeyGroupFromIDString(ctx, nil, skgIDStr, cc.DBSession, nil)
		if serr != nil {
			if serr == ErrInvalidID {
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Failed to create Instances, Invalid SSH Key Group ID: %s", skgIDStr), nil)
			}
			if serr == cdb.ErrDoesNotExist {
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Failed to create Instances, Could not find SSH Key Group with ID: %s", skgIDStr), nil)
			}
			logger.Warn().Err(serr).Str("SSH Key Group ID", skgIDStr).Msg("error retrieving SSH Key Group from DB by ID")
			return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Failed to retrieve SSH Key Group with ID `%s` specified in request, DB error", skgIDStr), nil)
		}

		if sshkeygroup.TenantID != tenant.ID {
			logger.Warn().Str("Tenant ID", tenant.ID.String()).Str("SSH Key Group ID", skgIDStr).Msg("SSH Key Group does not belong to current Tenant")
			return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Failed to create Instances, SSH Key Group with ID: %s does not belong to Tenant", skgIDStr), nil)
		}

		_, serr = skgsaDAO.GetBySSHKeyGroupIDAndSiteID(ctx, nil, sshkeygroup.ID, site.ID, nil)
		if serr != nil {
			if serr == cdb.ErrDoesNotExist {
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("SSH Key Group with ID: %s is not associated with the Site where Instances are being created", skgIDStr), nil)
			}
			logger.Warn().Err(serr).Str("SSH Key Group ID", skgIDStr).Msg("error retrieving SSH Key Group Site Association from DB by SSH Key Group ID & Site ID")
			return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Failed to determine if SSH Key Group: %s is associated with the Site where Instances are being created, DB error", skgIDStr), nil)
		}

		cc.SSHKeyGroups = append(cc.SSHKeyGroups, *sshkeygroup)
	}

	return nil
}

// BuildOsConfig validates and builds the OS configuration for the Temporal workflow.
// Populates cc.OsConfig and cc.OsID.
// Requires cc.Site to be populated.
func (cc *InstanceCreateContext) BuildOsConfig(ctx context.Context, req InstanceOSConfigProvider, siteID uuid.UUID) *cerr.APIError {
	logger := cc.Logger

	if req.GetOperatingSystemID() == nil || *req.GetOperatingSystemID() == "" {
		if err := req.ValidateAndSetOperatingSystemData(cc.Cfg, nil); err != nil {
			logger.Error().Err(err).Msg("failed to validate OperatingSystem")
			return cerr.NewAPIError(http.StatusBadRequest, "Failed to validate OperatingSystem data", err)
		}

		cc.OsConfig = &cwssaws.OperatingSystem{
			RunProvisioningInstructionsOnEveryBoot: *req.GetAlwaysBootWithCustomIpxe(),
			PhoneHomeEnabled:                       *req.GetPhoneHomeEnabled(),
			Variant: &cwssaws.OperatingSystem_Ipxe{
				Ipxe: &cwssaws.IpxeOperatingSystem{
					IpxeScript: *req.GetIpxeScript(),
				},
			},
			UserData: req.GetUserData(),
		}
		cc.OsID = nil
		return nil
	}

	id, err := uuid.Parse(*req.GetOperatingSystemID())
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse OperatingSystemID")
		return cerr.NewAPIError(http.StatusBadRequest, "Unable to parse `operatingSystemId` specified", validation.Errors{
			"operatingSystemId": errors.New(*req.GetOperatingSystemID()),
		})
	}

	osID := &id

	osDAO := cdbm.NewOperatingSystemDAO(cc.DBSession)
	os, serr := osDAO.GetByID(ctx, nil, *osID, nil)
	if serr != nil {
		if serr == cdb.ErrDoesNotExist {
			return cerr.NewAPIError(http.StatusBadRequest, "Could not find OperatingSystem with ID specified in request data", validation.Errors{
				"id": errors.New(osID.String()),
			})
		}
		logger.Error().Err(serr).Msg("error retrieving OperatingSystem from DB by ID")
		return cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve OperatingSystem with ID specified in request data, DB error", validation.Errors{
			"id": errors.New(osID.String()),
		})
	}

	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("OperatingSystem ID", os.ID.String())
	})

	if os.TenantID.String() != req.GetTenantID() {
		logger.Error().Msg("OperatingSystem in request is not owned by tenant")
		return cerr.NewAPIError(http.StatusBadRequest, "OperatingSystem specified in request is not owned by Tenant", nil)
	}

	if os.Type == cdbm.OperatingSystemTypeImage {
		ossaDAO := cdbm.NewOperatingSystemSiteAssociationDAO(cc.DBSession)
		_, ossaCount, err := ossaDAO.GetAll(
			ctx, nil,
			cdbm.OperatingSystemSiteAssociationFilterInput{
				OperatingSystemIDs: []uuid.UUID{id},
				SiteIDs:            []uuid.UUID{siteID},
			},
			cdbp.PageInput{Limit: cdb.GetIntPtr(1)},
			nil,
		)
		if err != nil {
			logger.Error().Msgf("Error retrieving OperatingSystemAssociations for OS: %s", err)
			return cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve OperatingSystemAssociations for OS with ID specified in request data, DB error", validation.Errors{
				"id": errors.New(osID.String()),
			})
		}
		if ossaCount == 0 {
			logger.Error().Msg("OperatingSystem does not belong to VPC site")
			return cerr.NewAPIError(http.StatusBadRequest, "OperatingSystem specified in request is not in VPC site", nil)
		}
	}

	err = req.ValidateAndSetOperatingSystemData(cc.Cfg, os)
	if err != nil {
		logger.Error().Msgf("OperatingSystem options validation failed: %s", err)
		return cerr.NewAPIError(http.StatusBadRequest, "OperatingSystem options validation failed", err)
	}

	if os.Type == cdbm.OperatingSystemTypeIPXE {
		cc.OsConfig = &cwssaws.OperatingSystem{
			RunProvisioningInstructionsOnEveryBoot: *req.GetAlwaysBootWithCustomIpxe(),
			PhoneHomeEnabled:                       *req.GetPhoneHomeEnabled(),
			Variant: &cwssaws.OperatingSystem_Ipxe{
				Ipxe: &cwssaws.IpxeOperatingSystem{
					IpxeScript: *req.GetIpxeScript(),
				},
			},
			UserData: req.GetUserData(),
		}
	} else {
		cc.OsConfig = &cwssaws.OperatingSystem{
			PhoneHomeEnabled: *req.GetPhoneHomeEnabled(),
			Variant: &cwssaws.OperatingSystem_OsImageId{
				OsImageId: &cwssaws.UUID{
					Value: os.ID.String(),
				},
			},
			UserData: req.GetUserData(),
		}
	}

	cc.OsID = osID
	return nil
}

// ValidateIBInterfaces validates InfiniBand interfaces against instance type capabilities,
// validates partition ownership and state, and populates cc.DBIBInterfaces.
// Requires cc.Tenant and cc.Site to be populated.
func (cc *InstanceCreateContext) ValidateIBInterfaces(ctx context.Context, ibIfcs []cam.APIInfiniBandInterfaceCreateOrUpdateRequest, instanceTypeID uuid.UUID) *cerr.APIError {
	logger := cc.Logger
	tenant := cc.Tenant
	site := cc.Site

	cc.DBIBInterfaces = nil

	if len(ibIfcs) == 0 {
		return nil
	}

	mcDAO := cdbm.NewMachineCapabilityDAO(cc.DBSession)
	itIbCaps, itIbCapCount, err := mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{instanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeInfiniBand), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving InfiniBand Machine Capabilities from DB for Instance Type")
		return cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve InfiniBand Capabilities for Instance Type, DB error", nil)
	}
	if itIbCapCount == 0 {
		logger.Warn().Msg("InfiniBand interfaces specified but Instance Type doesn't have InfiniBand Capability")
		return cerr.NewAPIError(http.StatusBadRequest, "InfiniBand Interfaces cannot be specified if Instance Type doesn't have InfiniBand Capability", nil)
	}

	ibpIDs := make([]uuid.UUID, 0, len(ibIfcs))
	for _, ibic := range ibIfcs {
		ibpID, err := uuid.Parse(ibic.InfiniBandPartitionID)
		if err != nil {
			logger.Warn().Err(err).Msg("error parsing infiniband partition id")
			return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Partition ID %v is not valid", ibic.InfiniBandPartitionID), nil)
		}
		ibpIDs = append(ibpIDs, ibpID)
	}

	ibpDAO := cdbm.NewInfiniBandPartitionDAO(cc.DBSession)
	ibpList, _, err := ibpDAO.GetAll(ctx, nil, cdbm.InfiniBandPartitionFilterInput{
		InfiniBandPartitionIDs: ibpIDs,
	}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving InfiniBand Partitions from DB")
		return cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve Partitions specified in request data, DB error", nil)
	}

	ibpMap := make(map[uuid.UUID]*cdbm.InfiniBandPartition, len(ibpList))
	for i := range ibpList {
		ibpMap[ibpList[i].ID] = &ibpList[i]
	}

	dbibic := []cdbm.InfiniBandInterface{}
	for i, ibic := range ibIfcs {
		ibpID := ibpIDs[i]
		ibp, exists := ibpMap[ibpID]
		if !exists {
			return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Could not find Partition with ID %v", ibic.InfiniBandPartitionID), nil)
		}

		if ibp.SiteID != site.ID {
			logger.Warn().Msgf("InfiniBandPartition: %v does not match with Instance Site", ibpID)
			return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Partition %v does not match with Instance Site", ibpID), nil)
		}

		if ibp.TenantID != tenant.ID {
			logger.Warn().Msgf("InfiniBandPartition: %v is not owned by Tenant", ibpID)
			return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Partition %v is not owned by Tenant", ibpID), nil)
		}

		if ibp.ControllerIBPartitionID == nil || ibp.Status != cdbm.InfiniBandPartitionStatusReady {
			logger.Warn().Msgf("InfiniBandPartition: %v is not in Ready state", ibpID)
			return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Partition %v is not in Ready state", ibpID), nil)
		}

		dbibic = append(dbibic, cdbm.InfiniBandInterface{
			InfiniBandPartitionID: ibp.ID,
			Device:                ibic.Device,
			Vendor:                ibic.Vendor,
			DeviceInstance:        ibic.DeviceInstance,
			IsPhysical:            ibic.IsPhysical,
			VirtualFunctionID:     ibic.VirtualFunctionID,
		})
	}

	err = cam.ValidateInfiniBandInterfaces(itIbCaps, ibIfcs)
	if err != nil {
		logger.Error().Err(err).Msg("InfiniBand interfaces validation failed")
		return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("InfiniBand interfaces validation failed: %v", err), err)
	}

	cc.DBIBInterfaces = dbibic
	return nil
}

// ValidateNVLinkInterfaces validates NVLink interfaces against instance type capabilities,
// validates logical partition ownership and state, and populates cc.DBNVLInterfaces.
// If no NVLink interfaces are specified but a default NVLink logical partition exists on the VPC,
// it generates default interfaces from the GPU capabilities.
// Requires cc.VPC, cc.Site, cc.Tenant, and cc.DefaultNvllpID to be populated.
func (cc *InstanceCreateContext) ValidateNVLinkInterfaces(ctx context.Context, nvlIfcs []cam.APINVLinkInterfaceCreateOrUpdateRequest, instanceTypeID uuid.UUID) *cerr.APIError {
	logger := cc.Logger
	tenant := cc.Tenant
	site := cc.Site

	mcDAO := cdbm.NewMachineCapabilityDAO(cc.DBSession)
	itNvlCaps, itNvlCapCount, err := mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{instanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeGPU), nil, nil, nil, nil, nil, cdb.GetStrPtr(cdbm.MachineCapabilityDeviceTypeNVLink), nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving GPU (NVLink) Machine Capabilities from DB for Instance Type")
		return cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve GPU Capabilities for Instance Type, DB error", nil)
	}

	cc.DBNVLInterfaces = nil

	if len(nvlIfcs) > 0 {
		if itNvlCapCount == 0 {
			logger.Warn().Msg("NVLink interfaces specified but Instance Type doesn't have GPU (NVLink) Capability")
			return cerr.NewAPIError(http.StatusBadRequest, "NVLink Interfaces cannot be specified if Instance Type doesn't have GPU Capabilities", nil)
		}

		err = cam.ValidateNVLinkInterfaces(itNvlCaps, nvlIfcs)
		if err != nil {
			logger.Error().Err(err).Msg("NVLink interfaces validation failed")
			return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink interfaces validation failed: %v", err), err)
		}

		nvllpIDs := make([]uuid.UUID, 0, len(nvlIfcs))
		for _, nvlifc := range nvlIfcs {
			nvllpID, err := uuid.Parse(nvlifc.NVLinkLogicalPartitionID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing NVLink Logical Partition id")
				return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition ID %v is not valid", nvlifc.NVLinkLogicalPartitionID), nil)
			}
			nvllpIDs = append(nvllpIDs, nvllpID)
		}

		var nvllpMap map[uuid.UUID]*cdbm.NVLinkLogicalPartition
		if cc.DefaultNvllpID == nil {
			nvllpDAO := cdbm.NewNVLinkLogicalPartitionDAO(cc.DBSession)
			uniqueNvllpIDs := make([]uuid.UUID, 0, len(nvllpIDs))
			seenIDs := make(map[uuid.UUID]bool)
			for _, id := range nvllpIDs {
				if !seenIDs[id] {
					seenIDs[id] = true
					uniqueNvllpIDs = append(uniqueNvllpIDs, id)
				}
			}

			nvllpList, _, err := nvllpDAO.GetAll(ctx, nil, cdbm.NVLinkLogicalPartitionFilterInput{
				NVLinkLogicalPartitionIDs: uniqueNvllpIDs,
			}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
			if err != nil {
				logger.Error().Err(err).Msg("error retrieving NVLink Logical Partitions from DB")
				return cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve NVLink Logical Partitions specified in request data, DB error", nil)
			}

			nvllpMap = make(map[uuid.UUID]*cdbm.NVLinkLogicalPartition, len(nvllpList))
			for i := range nvllpList {
				nvllpMap[nvllpList[i].ID] = &nvllpList[i]
			}
		}

		dbnvlic := []cdbm.NVLinkInterface{}
		for i, nvlifc := range nvlIfcs {
			nvllpID := nvllpIDs[i]

			if cc.DefaultNvllpID != nil {
				if nvllpID != *cc.DefaultNvllpID {
					return cerr.NewAPIError(http.StatusBadRequest, "NVLink Logical Partition specified for NVLink Interface does not match NVLink Logical Partition of VPC", nil)
				}
			} else {
				nvllp, exists := nvllpMap[nvllpID]
				if !exists {
					return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Could not find NVLink Logical Partition with ID %v", nvllpID), nil)
				}

				if nvllp.SiteID != site.ID {
					logger.Warn().Msgf("NVLink Logical Partition: %v does not match with Instance Site", nvllpID)
					return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition %v does not match with Instance Site", nvllpID), nil)
				}

				if nvllp.TenantID != tenant.ID {
					logger.Warn().Msgf("NVLink Logical Partition: %v is not owned by Tenant", nvllpID)
					return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition %v is not owned by Tenant", nvllpID), nil)
				}

				if nvllp.Status != cdbm.NVLinkLogicalPartitionStatusReady {
					logger.Warn().Msgf("NVLink Logical Partition: %v is not in Ready state", nvllpID)
					return cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition %v is not in Ready state", nvllpID), nil)
				}
			}

			dbnvlic = append(dbnvlic, cdbm.NVLinkInterface{
				NVLinkLogicalPartitionID: nvllpID,
				DeviceInstance:           nvlifc.DeviceInstance,
			})
		}
		cc.DBNVLInterfaces = dbnvlic
	} else if cc.DefaultNvllpID != nil {
		dbnvlic := []cdbm.NVLinkInterface{}
		for _, nvlCap := range itNvlCaps {
			if nvlCap.Count != nil {
				for i := 0; i < *nvlCap.Count; i++ {
					dbnvlic = append(dbnvlic, cdbm.NVLinkInterface{
						NVLinkLogicalPartitionID: *cc.DefaultNvllpID,
						Device:                   cdb.GetStrPtr(nvlCap.Name),
						DeviceInstance:           i,
					})
				}
			}
		}
		cc.DBNVLInterfaces = dbnvlic
	}

	return nil
}

// ValidateDPUInterfaces validates DPU interface capabilities against instance type machine capabilities.
// Requires cc.DBInterfaces and cc.IsDeviceInfoPresent to be populated.
func (cc *InstanceCreateContext) ValidateDPUInterfaces(ctx context.Context, instanceTypeID uuid.UUID) *cerr.APIError {
	if !cc.IsDeviceInfoPresent {
		return nil
	}

	logger := cc.Logger

	mcDAO := cdbm.NewMachineCapabilityDAO(cc.DBSession)
	itDpuCaps, itDpuCapCount, err := mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{instanceTypeID},
		cdb.GetStrPtr(cdbm.MachineCapabilityTypeNetwork), nil, nil, nil, nil, nil,
		cdb.GetStrPtr(cdbm.MachineCapabilityDeviceTypeDPU), nil, nil, nil, nil, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving DPU Machine Capabilities")
		return cerr.NewAPIError(http.StatusInternalServerError,
			"Failed to retrieve DPU Capabilities for Instance Type, DB error", nil)
	}

	if itDpuCapCount == 0 {
		logger.Warn().Msg("Device/DeviceInstance specified but Instance Type doesn't have DPU Capability")
		return cerr.NewAPIError(http.StatusBadRequest,
			"Device and DeviceInstance cannot be specified if Instance Type doesn't have Network Capabilities with DPU device type", nil)
	}

	err = cam.ValidateMultiEthernetDeviceInterfaces(itDpuCaps, cc.DBInterfaces)
	if err != nil {
		logger.Error().Err(err).Msg("DPU interfaces validation failed")
		return cerr.NewAPIError(http.StatusBadRequest,
			fmt.Sprintf("DPU interfaces validation failed: %v", err), err)
	}

	return nil
}

// BuildInterfaceConfigs builds Temporal workflow InstanceInterfaceConfig from created interfaces.
func BuildInterfaceConfigs(ifcs []cdbm.Interface, subnetIDMap map[uuid.UUID]*cdbm.Subnet) []*cwssaws.InstanceInterfaceConfig {
	configs := make([]*cwssaws.InstanceInterfaceConfig, 0, len(ifcs))

	for _, ifc := range ifcs {
		interfaceConfig := &cwssaws.InstanceInterfaceConfig{
			FunctionType: cwssaws.InterfaceFunctionType_VIRTUAL_FUNCTION,
		}
		if ifc.SubnetID != nil {
			interfaceConfig.NetworkSegmentId = &cwssaws.NetworkSegmentId{
				Value: subnetIDMap[*ifc.SubnetID].ControllerNetworkSegmentID.String(),
			}
			interfaceConfig.NetworkDetails = &cwssaws.InstanceInterfaceConfig_SegmentId{
				SegmentId: &cwssaws.NetworkSegmentId{
					Value: subnetIDMap[*ifc.SubnetID].ControllerNetworkSegmentID.String(),
				},
			}
		}
		if ifc.VpcPrefixID != nil {
			interfaceConfig.NetworkDetails = &cwssaws.InstanceInterfaceConfig_VpcPrefixId{
				VpcPrefixId: &cwssaws.VpcPrefixId{Value: ifc.VpcPrefixID.String()},
			}
		}
		if ifc.IsPhysical {
			interfaceConfig.FunctionType = cwssaws.InterfaceFunctionType_PHYSICAL_FUNCTION
		}
		if ifc.Device != nil && ifc.DeviceInstance != nil {
			interfaceConfig.Device = ifc.Device
			interfaceConfig.DeviceInstance = uint32(*ifc.DeviceInstance)
		}
		if !ifc.IsPhysical && ifc.VirtualFunctionID != nil {
			vfID := uint32(*ifc.VirtualFunctionID)
			interfaceConfig.VirtualFunctionId = &vfID
		}
		configs = append(configs, interfaceConfig)
	}

	return configs
}

// BuildIBInterfaceConfigs builds Temporal workflow InstanceIBInterfaceConfig from created InfiniBand interfaces.
func BuildIBInterfaceConfigs(ibifcs []cdbm.InfiniBandInterface) []*cwssaws.InstanceIBInterfaceConfig {
	configs := make([]*cwssaws.InstanceIBInterfaceConfig, 0, len(ibifcs))

	for _, ifc := range ibifcs {
		ibInterfaceConfig := &cwssaws.InstanceIBInterfaceConfig{
			Device:         ifc.Device,
			Vendor:         ifc.Vendor,
			DeviceInstance: uint32(ifc.DeviceInstance),
			FunctionType:   cwssaws.InterfaceFunctionType_PHYSICAL_FUNCTION,
			IbPartitionId:  &cwssaws.IBPartitionId{Value: ifc.InfiniBandPartitionID.String()},
		}
		if !ifc.IsPhysical {
			ibInterfaceConfig.FunctionType = cwssaws.InterfaceFunctionType_VIRTUAL_FUNCTION
			if ifc.VirtualFunctionID != nil {
				vfID := uint32(*ifc.VirtualFunctionID)
				ibInterfaceConfig.VirtualFunctionId = &vfID
			}
		}
		configs = append(configs, ibInterfaceConfig)
	}

	return configs
}

// BuildNVLinkInterfaceConfigs builds Temporal workflow InstanceNVLinkGpuConfig from created NVLink interfaces.
func BuildNVLinkInterfaceConfigs(nvlifcs []cdbm.NVLinkInterface) []*cwssaws.InstanceNVLinkGpuConfig {
	configs := make([]*cwssaws.InstanceNVLinkGpuConfig, 0, len(nvlifcs))

	for _, nvlifc := range nvlifcs {
		nvlInterfaceConfig := &cwssaws.InstanceNVLinkGpuConfig{
			DeviceInstance:     uint32(nvlifc.DeviceInstance),
			LogicalPartitionId: &cwssaws.NVLinkLogicalPartitionId{Value: nvlifc.NVLinkLogicalPartitionID.String()},
		}
		configs = append(configs, nvlInterfaceConfig)
	}

	return configs
}

// BuildDPUServiceConfigs builds Temporal workflow InstanceDpuExtensionServiceConfig from created deployments.
func BuildDPUServiceConfigs(desds []cdbm.DpuExtensionServiceDeployment) []*cwssaws.InstanceDpuExtensionServiceConfig {
	configs := make([]*cwssaws.InstanceDpuExtensionServiceConfig, 0, len(desds))

	for _, desd := range desds {
		configs = append(configs, &cwssaws.InstanceDpuExtensionServiceConfig{
			ServiceId: desd.DpuExtensionServiceID.String(),
			Version:   desd.Version,
		})
	}

	return configs
}

// BuildInstanceAllocationRequest builds a single InstanceAllocationRequest for the Temporal workflow.
func BuildInstanceAllocationRequest(
	instance *cdbm.Instance,
	tenant *cdbm.Tenant,
	osConfig *cwssaws.OperatingSystem,
	sshKeyGroupIDs []string,
	interfaceConfigs []*cwssaws.InstanceInterfaceConfig,
	ibInterfaceConfigs []*cwssaws.InstanceIBInterfaceConfig,
	nvlInterfaceConfigs []*cwssaws.InstanceNVLinkGpuConfig,
	desdConfigs []*cwssaws.InstanceDpuExtensionServiceConfig,
	allowUnhealthyMachine bool,
) *cwssaws.InstanceAllocationRequest {
	createLabels := []*cwssaws.Label{}
	for k, v := range instance.Labels {
		createLabels = append(createLabels, &cwssaws.Label{
			Key:   k,
			Value: &v,
		})
	}

	description := ""
	if instance.Description != nil {
		description = *instance.Description
	}

	req := &cwssaws.InstanceAllocationRequest{
		InstanceId: &cwssaws.InstanceId{Value: GetSiteInstanceID(instance).String()},
		MachineId:  &cwssaws.MachineId{Id: *instance.MachineID},
		Metadata: &cwssaws.Metadata{
			Name:        instance.Name,
			Description: description,
			Labels:      createLabels,
		},
		Config: &cwssaws.InstanceConfig{
			NetworkSecurityGroupId: instance.NetworkSecurityGroupID,
			Tenant: &cwssaws.TenantConfig{
				TenantOrganizationId: tenant.Org,
				TenantKeysetIds:      sshKeyGroupIDs,
			},
			Os: osConfig,
			Network: &cwssaws.InstanceNetworkConfig{
				Interfaces: interfaceConfigs,
			},
			Infiniband: &cwssaws.InstanceInfinibandConfig{
				IbInterfaces: ibInterfaceConfigs,
			},
			DpuExtensionServices: &cwssaws.InstanceDpuExtensionServicesConfig{
				ServiceConfigs: desdConfigs,
			},
			Nvlink: &cwssaws.InstanceNVLinkConfig{
				GpuConfigs: nvlInterfaceConfigs,
			},
		},
		AllowUnhealthyMachine: allowUnhealthyMachine,
	}

	if instance.InstanceTypeID != nil {
		req.InstanceTypeId = cdb.GetStrPtr(instance.InstanceTypeID.String())
	}

	return req
}
