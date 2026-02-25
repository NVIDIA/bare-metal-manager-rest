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
// instance create requests, allowing a single ValidateAndBuildOsConfig implementation.
type InstanceOSConfigProvider interface {
	GetOperatingSystemID() *string
	GetTenantID() string
	ValidateAndSetOperatingSystemData(cfg *config.Config, os *cdbm.OperatingSystem) error
	GetAlwaysBootWithCustomIpxe() *bool
	GetPhoneHomeEnabled() *bool
	GetIpxeScript() *string
	GetUserData() *string
}

// InstanceCreateValidator holds infrastructure dependencies for shared instance creation validation.
// Domain data flows explicitly through method parameters and return values.
type InstanceCreateValidator struct {
	DBSession *cdb.Session
	Cfg       *config.Config
	Logger    *zerolog.Logger
}

// NewInstanceCreateValidator creates a validator with infrastructure dependencies.
func NewInstanceCreateValidator(dbSession *cdb.Session, cfg *config.Config, logger *zerolog.Logger) *InstanceCreateValidator {
	return &InstanceCreateValidator{
		DBSession: dbSession,
		Cfg:       cfg,
		Logger:    logger,
	}
}

// InterfaceValidationResult holds the products of ValidateNetworkInterfaces.
type InterfaceValidationResult struct {
	DBInterfaces        []cdbm.Interface
	SubnetIDMap         map[uuid.UUID]*cdbm.Subnet
	VpcPrefixIDMap      map[uuid.UUID]*cdbm.VpcPrefix
	IsDeviceInfoPresent bool
}

// ValidateTenantAndVPC validates tenant ownership, VPC state, and site readiness.
func (icv *InstanceCreateValidator) ValidateTenantAndVPC(ctx context.Context, org, tenantID, vpcID string) (
	*cdbm.Tenant, *cdbm.Vpc, *cdbm.Site, *uuid.UUID, *cerr.APIError,
) {
	logger := icv.Logger

	tenant, err := GetTenantForOrg(ctx, nil, icv.DBSession, org)
	if err != nil {
		if err == ErrOrgTenantNotFound {
			logger.Warn().Err(err).Msg("Org does not have a Tenant associated")
			return nil, nil, nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Org does not have a Tenant associated", nil)
		}
		logger.Error().Err(err).Msg("unable to retrieve tenant for org")
		return nil, nil, nil, nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve tenant for org", nil)
	}

	apiTenant, err := GetTenantFromIDString(ctx, nil, tenantID, icv.DBSession)
	if err != nil {
		logger.Warn().Err(err).Msg("error retrieving tenant from request")
		return nil, nil, nil, nil, cerr.NewAPIError(http.StatusBadRequest, "TenantID in request is not valid", nil)
	}
	if apiTenant.ID != tenant.ID {
		logger.Warn().Msg("tenant id in request does not match tenant in org")
		return nil, nil, nil, nil, cerr.NewAPIError(http.StatusBadRequest, "TenantID in request does not match tenant in org", nil)
	}

	vpc, err := GetVpcFromIDString(ctx, nil, vpcID, []string{cdbm.NVLinkLogicalPartitionRelationName}, icv.DBSession)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return nil, nil, nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Could not find VPC with ID specified in request data", nil)
		}
		logger.Warn().Err(err).Str("vpcId", vpcID).Msg("error retrieving VPC from request")
		return nil, nil, nil, nil, cerr.NewAPIError(http.StatusBadRequest, "VpcID in request is not valid", nil)
	}

	if vpc.TenantID != tenant.ID {
		logger.Warn().Msg("tenant id in request does not match tenant in VPC")
		return nil, nil, nil, nil, cerr.NewAPIError(http.StatusBadRequest, "VPC specified in request is not owned by Tenant", nil)
	}

	if vpc.ControllerVpcID == nil || vpc.Status != cdbm.VpcStatusReady {
		logger.Warn().Msg("VPC specified in request data is not ready")
		return nil, nil, nil, nil, cerr.NewAPIError(http.StatusBadRequest, "VPC specified in request data is not ready", nil)
	}

	var defaultNvllpID *uuid.UUID
	if vpc.NVLinkLogicalPartitionID != nil {
		defaultNvllpID = vpc.NVLinkLogicalPartitionID
	}

	siteDAO := cdbm.NewSiteDAO(icv.DBSession)
	site, err := siteDAO.GetByID(ctx, nil, vpc.SiteID, nil, false)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return nil, nil, nil, nil, cerr.NewAPIError(http.StatusBadRequest, "The Site where this Instance is being created could not be found", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB by ID")
		return nil, nil, nil, nil, cerr.NewAPIError(http.StatusInternalServerError, "The Site where this Instance is being created could not be retrieved", nil)
	}

	if site.Status != cdbm.SiteStatusRegistered {
		logger.Warn().Str("Site ID", site.ID.String()).Str("Site Status", site.Status).
			Msg("The Site where this Instance is being created is not in Registered state")
		return nil, nil, nil, nil, cerr.NewAPIError(http.StatusBadRequest, "The Site where this Instance is being created is not in Registered state", nil)
	}

	return tenant, vpc, site, defaultNvllpID, nil
}

// ValidateNetworkInterfaces validates subnet and VPC prefix interfaces.
func (icv *InstanceCreateValidator) ValidateNetworkInterfaces(ctx context.Context, tenant *cdbm.Tenant, vpc *cdbm.Vpc, interfaces []cam.APIInterfaceCreateOrUpdateRequest) (*InterfaceValidationResult, *cerr.APIError) {
	logger := icv.Logger

	subnetDAO := cdbm.NewSubnetDAO(icv.DBSession)
	vpDAO := cdbm.NewVpcPrefixDAO(icv.DBSession)

	subnetIDs := []uuid.UUID{}
	vpcPrefixIDs := []uuid.UUID{}

	for _, ifc := range interfaces {
		if ifc.SubnetID != nil {
			subnetID, err := uuid.Parse(*ifc.SubnetID)
			if err != nil {
				logger.Error().Err(err).Str("subnetID", *ifc.SubnetID).Msg("error parsing subnet id")
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Subnet ID: %s specified in interfaces data in request is not valid", *ifc.SubnetID), nil)
			}
			subnetIDs = append(subnetIDs, subnetID)
		}
		if ifc.VpcPrefixID != nil {
			vpcPrefixID, err := uuid.Parse(*ifc.VpcPrefixID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing vpcprefix id in instance vpcprefix request")
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC Prefix ID: %s specified in interfaces data in request is not valid", *ifc.VpcPrefixID), nil)
			}
			vpcPrefixIDs = append(vpcPrefixIDs, vpcPrefixID)
		}
	}

	subnetIDMap := make(map[uuid.UUID]*cdbm.Subnet)
	if len(subnetIDs) > 0 {
		subnets, _, err := subnetDAO.GetAll(ctx, nil, cdbm.SubnetFilterInput{SubnetIDs: subnetIDs}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Subnets from DB by IDs")
			return nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve Subnets from DB by IDs", nil)
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
			return nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve VPC Prefixes from DB by IDs", nil)
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
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Could not find Subnet with ID: %v specified in request data", subnetID), nil)
			}

			if subnet.TenantID != tenant.ID {
				logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request is not owned by Tenant", subnetID))
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request is not owned by Tenant", subnetID), nil)
			}

			if subnet.ControllerNetworkSegmentID == nil || subnet.Status != cdbm.SubnetStatusReady {
				logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request data is not in Ready state", subnetID))
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request data is not in Ready state", subnetID), nil)
			}

			if subnet.VpcID != vpc.ID {
				logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request does not match with VPC", subnetID))
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request does not match with VPC", subnetID), nil)
			}

			if vpc.NetworkVirtualizationType != nil && *vpc.NetworkVirtualizationType != cdbm.VpcEthernetVirtualizer {
				logger.Warn().Msg(fmt.Sprintf("VPC: %v specified in request must have Ethernet network virtualization type in order to create Subnet based interfaces", vpc.ID))
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC: %v specified in request must have Ethernet network virtualization type in order to create Subnet based interfaces", vpc.ID), nil)
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
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Could not find VPC Prefix with ID: %v specified in request data", vpcPrefixUUID), nil)
			}

			if vpcPrefix.TenantID != tenant.ID {
				logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request is not owned by Tenant", vpcPrefixUUID))
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request is not owned by Tenant", vpcPrefixUUID), nil)
			}

			if vpcPrefix.Status != cdbm.VpcPrefixStatusReady {
				logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request data is not in Ready state", vpcPrefixUUID))
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request data is not in Ready state", vpcPrefixUUID), nil)
			}

			if vpcPrefix.VpcID != vpc.ID {
				logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request does not match with VPC", vpcPrefixUUID))
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request does not match with VPC", vpcPrefixUUID), nil)
			}

			if vpc.NetworkVirtualizationType == nil || *vpc.NetworkVirtualizationType != cdbm.VpcFNN {
				logger.Warn().Msg(fmt.Sprintf("VPC: %v specified in request must have FNN network virtualization type in order to create VPC Prefix based interfaces", vpc.ID))
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("VPC: %v specified in request must have FNN network virtualization type in order to create VPC Prefix based interfaces", vpc.ID), nil)
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

	return &InterfaceValidationResult{
		DBInterfaces:        dbInterfaces,
		SubnetIDMap:         subnetIDMap,
		VpcPrefixIDMap:      vpcPrefixIDMap,
		IsDeviceInfoPresent: isDeviceInfoPresent,
	}, nil
}

// ValidateDPUExtensionServices validates DPU extension service deployments.
func (icv *InstanceCreateValidator) ValidateDPUExtensionServices(ctx context.Context, tenant *cdbm.Tenant, site *cdbm.Site, deployments []cam.APIDpuExtensionServiceDeploymentRequest) (map[string]*cdbm.DpuExtensionService, *cerr.APIError) {
	logger := icv.Logger

	desIDMap := map[string]*cdbm.DpuExtensionService{}

	if len(deployments) == 0 {
		return desIDMap, nil
	}

	desIDs := make([]uuid.UUID, 0, len(deployments))
	uniqueDesIDs := make([]uuid.UUID, 0, len(deployments))
	seenDesIDs := make(map[uuid.UUID]bool, len(deployments))
	for _, adesdr := range deployments {
		desID, err := uuid.Parse(adesdr.DpuExtensionServiceID)
		if err != nil {
			logger.Warn().Err(err).Str("serviceID", adesdr.DpuExtensionServiceID).
				Msg("error parsing DPU Extension Service ID")
			return nil, cerr.NewAPIError(http.StatusBadRequest,
				fmt.Sprintf("Invalid DPU Extension Service ID: %s specified in request", adesdr.DpuExtensionServiceID), nil)
		}
		desIDs = append(desIDs, desID)
		if !seenDesIDs[desID] {
			seenDesIDs[desID] = true
			uniqueDesIDs = append(uniqueDesIDs, desID)
		}
	}

	desDAO := cdbm.NewDpuExtensionServiceDAO(icv.DBSession)
	desList, _, err := desDAO.GetAll(ctx, nil, cdbm.DpuExtensionServiceFilterInput{
		DpuExtensionServiceIDs: uniqueDesIDs,
	}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving DPU Extension Services from DB")
		return nil, cerr.NewAPIError(http.StatusInternalServerError,
			"Failed to retrieve DPU Extension Services from DB by IDs", nil)
	}

	desMap := make(map[uuid.UUID]*cdbm.DpuExtensionService, len(desList))
	for i := range desList {
		desMap[desList[i].ID] = &desList[i]
	}

	for i, desID := range desIDs {
		des, exists := desMap[desID]
		if !exists {
			return nil, cerr.NewAPIError(http.StatusBadRequest,
				fmt.Sprintf("Could not find DPU Extension Service with ID: %v specified in request data", desID), nil)
		}

		if des.TenantID != tenant.ID {
			logger.Warn().Str("tenantID", tenant.ID.String()).Str("serviceID", desID.String()).
				Msg("DPU Extension Service does not belong to current Tenant")
			return nil, cerr.NewAPIError(http.StatusForbidden,
				fmt.Sprintf("DPU Extension Service: %s does not belong to current Tenant", desID.String()), nil)
		}

		if des.SiteID != site.ID {
			logger.Warn().Str("siteID", site.ID.String()).Str("serviceID", desID.String()).
				Msg("DPU Extension Service does not belong to Site")
			return nil, cerr.NewAPIError(http.StatusForbidden,
				fmt.Sprintf("DPU Extension Service: %s does not belong to Site where Instance is being created", desID.String()), nil)
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
			return nil, cerr.NewAPIError(http.StatusBadRequest,
				fmt.Sprintf("Version: %s was not found for DPU Extension Service: %s", requestedVersion, desID.String()), nil)
		}

		desIDMap[desID.String()] = des
	}

	return desIDMap, nil
}

// ValidateNetworkSecurityGroup validates the network security group, if specified.
func (icv *InstanceCreateValidator) ValidateNetworkSecurityGroup(ctx context.Context, tenant *cdbm.Tenant, site *cdbm.Site, nsgID *string) *cerr.APIError {
	if nsgID == nil {
		return nil
	}

	logger := icv.Logger

	nsgDAO := cdbm.NewNetworkSecurityGroupDAO(icv.DBSession)
	nsg, err := nsgDAO.GetByID(ctx, nil, *nsgID, nil)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIError(http.StatusBadRequest,
				fmt.Sprintf("Could not find NetworkSecurityGroup with ID: %s specified in request", *nsgID), nil)
		}
		logger.Error().Err(err).Str("nsgID", *nsgID).Msg("error retrieving NetworkSecurityGroup with ID specified in request data")
		return cerr.NewAPIError(http.StatusInternalServerError,
			"Failed to retrieve NetworkSecurityGroup with ID specified in request data", nil)
	}

	if nsg.SiteID != site.ID {
		logger.Error().Str("siteID", site.ID.String()).Str("nsgID", *nsgID).
			Msg("NetworkSecurityGroup in request does not belong to Site")
		return cerr.NewAPIError(http.StatusForbidden,
			"NetworkSecurityGroup with ID specified in request data does not belong to Site", nil)
	}

	if nsg.TenantID != tenant.ID {
		logger.Error().Str("tenantID", tenant.ID.String()).Str("nsgID", *nsgID).
			Msg("NetworkSecurityGroup in request does not belong to Tenant")
		return cerr.NewAPIError(http.StatusForbidden,
			"NetworkSecurityGroup with ID specified in request data does not belong to Tenant", nil)
	}

	return nil
}

// ValidateSSHKeyGroups validates SSH key groups and their site associations.
func (icv *InstanceCreateValidator) ValidateSSHKeyGroups(ctx context.Context, tenant *cdbm.Tenant, site *cdbm.Site, sshKeyGroupIDs []string) ([]cdbm.SSHKeyGroup, *cerr.APIError) {
	logger := icv.Logger

	sshKeyGroups := []cdbm.SSHKeyGroup{}

	skgsaDAO := cdbm.NewSSHKeyGroupSiteAssociationDAO(icv.DBSession)

	for _, skgIDStr := range sshKeyGroupIDs {
		sshkeygroup, serr := GetSSHKeyGroupFromIDString(ctx, nil, skgIDStr, icv.DBSession, nil)
		if serr != nil {
			if serr == ErrInvalidID {
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Invalid SSH Key Group ID: %s specified in request", skgIDStr), nil)
			}
			if serr == cdb.ErrDoesNotExist {
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Could not find SSH Key Group with ID: %s specified in request data", skgIDStr), nil)
			}
			logger.Warn().Err(serr).Str("SSH Key Group ID", skgIDStr).Msg("error retrieving SSH Key Group from DB by ID")
			return nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve SSH Key Groups from DB by IDs", nil)
		}

		if sshkeygroup.TenantID != tenant.ID {
			logger.Warn().Str("Tenant ID", tenant.ID.String()).Str("SSH Key Group ID", skgIDStr).Msg("SSH Key Group does not belong to current Tenant")
			return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Failed to create Instance, SSH Key Group with ID: %s does not belong to Tenant", skgIDStr), nil)
		}

		_, serr = skgsaDAO.GetBySSHKeyGroupIDAndSiteID(ctx, nil, sshkeygroup.ID, site.ID, nil)
		if serr != nil {
			if serr == cdb.ErrDoesNotExist {
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("SSH Key Group: %s specified in request data is not associated with the Site where Instance is being created", skgIDStr), nil)
			}
			logger.Warn().Err(serr).Str("SSH Key Group ID", skgIDStr).Msg("error retrieving SSH Key Group Site Association from DB by SSH Key Group ID & Site ID")
			return nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve SSH Key Group Site Associations from DB", nil)
		}

		sshKeyGroups = append(sshKeyGroups, *sshkeygroup)
	}

	return sshKeyGroups, nil
}

// ValidateAndBuildOsConfig validates and builds the OS configuration for the Temporal workflow.
func (icv *InstanceCreateValidator) ValidateAndBuildOsConfig(ctx context.Context, req InstanceOSConfigProvider, siteID uuid.UUID) (*cwssaws.OperatingSystem, *uuid.UUID, *cerr.APIError) {
	logger := icv.Logger

	if req.GetOperatingSystemID() == nil || *req.GetOperatingSystemID() == "" {
		if err := req.ValidateAndSetOperatingSystemData(icv.Cfg, nil); err != nil {
			logger.Error().Err(err).Msg("failed to validate OperatingSystem")
			return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Failed to validate OperatingSystem data", err)
		}

		osConfig := &cwssaws.OperatingSystem{
			RunProvisioningInstructionsOnEveryBoot: *req.GetAlwaysBootWithCustomIpxe(),
			PhoneHomeEnabled:                       *req.GetPhoneHomeEnabled(),
			Variant: &cwssaws.OperatingSystem_Ipxe{
				Ipxe: &cwssaws.IpxeOperatingSystem{
					IpxeScript: *req.GetIpxeScript(),
				},
			},
			UserData: req.GetUserData(),
		}
		return osConfig, nil, nil
	}

	id, err := uuid.Parse(*req.GetOperatingSystemID())
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse OperatingSystemID")
		return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Unable to parse `operatingSystemId` specified", validation.Errors{
			"operatingSystemId": errors.New(*req.GetOperatingSystemID()),
		})
	}

	osID := &id

	osDAO := cdbm.NewOperatingSystemDAO(icv.DBSession)
	os, serr := osDAO.GetByID(ctx, nil, *osID, nil)
	if serr != nil {
		if serr == cdb.ErrDoesNotExist {
			return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Could not find OperatingSystem with ID specified in request data", validation.Errors{
				"id": errors.New(osID.String()),
			})
		}
		logger.Error().Err(serr).Msg("error retrieving OperatingSystem from DB by ID")
		return nil, nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve OperatingSystem with ID specified in request data, DB error", validation.Errors{
			"id": errors.New(osID.String()),
		})
	}

	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("OperatingSystem ID", os.ID.String())
	})

	if os.TenantID.String() != req.GetTenantID() {
		logger.Error().Msg("OperatingSystem in request is not owned by tenant")
		return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "OperatingSystem specified in request is not owned by Tenant", nil)
	}

	if os.Type == cdbm.OperatingSystemTypeImage {
		ossaDAO := cdbm.NewOperatingSystemSiteAssociationDAO(icv.DBSession)
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
			return nil, nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve OperatingSystemAssociations for OS with ID specified in request data, DB error", validation.Errors{
				"id": errors.New(osID.String()),
			})
		}
		if ossaCount == 0 {
			logger.Error().Msg("OperatingSystem does not belong to VPC site")
			return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "OperatingSystem specified in request is not in VPC site", nil)
		}
	}

	err = req.ValidateAndSetOperatingSystemData(icv.Cfg, os)
	if err != nil {
		logger.Error().Msgf("OperatingSystem options validation failed: %s", err)
		return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "OperatingSystem options validation failed", err)
	}

	var osConfig *cwssaws.OperatingSystem
	if os.Type == cdbm.OperatingSystemTypeIPXE {
		osConfig = &cwssaws.OperatingSystem{
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
		osConfig = &cwssaws.OperatingSystem{
			PhoneHomeEnabled: *req.GetPhoneHomeEnabled(),
			Variant: &cwssaws.OperatingSystem_OsImageId{
				OsImageId: &cwssaws.UUID{
					Value: os.ID.String(),
				},
			},
			UserData: req.GetUserData(),
		}
	}

	return osConfig, osID, nil
}

// ValidateInfiniBandInterfaces validates InfiniBand interfaces against instance type capabilities,
// validates partition ownership and state.
func (icv *InstanceCreateValidator) ValidateInfiniBandInterfaces(ctx context.Context, tenant *cdbm.Tenant, site *cdbm.Site, ibIfcs []cam.APIInfiniBandInterfaceCreateOrUpdateRequest, instanceTypeID uuid.UUID) ([]cdbm.InfiniBandInterface, *cerr.APIError) {
	logger := icv.Logger

	if len(ibIfcs) == 0 {
		return nil, nil
	}

	mcDAO := cdbm.NewMachineCapabilityDAO(icv.DBSession)
	itIbCaps, itIbCapCount, err := mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{instanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeInfiniBand), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving InfiniBand Machine Capabilities from DB for Instance Type")
		return nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve InfiniBand Capabilities for Instance Type, DB error", nil)
	}
	if itIbCapCount == 0 {
		logger.Warn().Msg("InfiniBand interfaces specified but Instance Type doesn't have InfiniBand Capability")
		return nil, cerr.NewAPIError(http.StatusBadRequest, "InfiniBand Interfaces cannot be specified if Instance Type doesn't have InfiniBand Capability", nil)
	}

	ibpIDs := make([]uuid.UUID, 0, len(ibIfcs))
	for _, ibic := range ibIfcs {
		ibpID, err := uuid.Parse(ibic.InfiniBandPartitionID)
		if err != nil {
			logger.Warn().Err(err).Msg("error parsing infiniband partition id")
			return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Partition ID %v is not valid", ibic.InfiniBandPartitionID), nil)
		}
		ibpIDs = append(ibpIDs, ibpID)
	}

	ibpDAO := cdbm.NewInfiniBandPartitionDAO(icv.DBSession)
	ibpList, _, err := ibpDAO.GetAll(ctx, nil, cdbm.InfiniBandPartitionFilterInput{
		InfiniBandPartitionIDs: ibpIDs,
	}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving InfiniBand Partitions from DB")
		return nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve Partitions specified in request data, DB error", nil)
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
			return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Could not find Partition with ID %v", ibic.InfiniBandPartitionID), nil)
		}

		if ibp.SiteID != site.ID {
			logger.Warn().Msgf("InfiniBandPartition: %v does not match with Instance Site", ibpID)
			return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Partition %v does not match with Instance Site", ibpID), nil)
		}

		if ibp.TenantID != tenant.ID {
			logger.Warn().Msgf("InfiniBandPartition: %v is not owned by Tenant", ibpID)
			return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Partition %v is not owned by Tenant", ibpID), nil)
		}

		if ibp.ControllerIBPartitionID == nil || ibp.Status != cdbm.InfiniBandPartitionStatusReady {
			logger.Warn().Msgf("InfiniBandPartition: %v is not in Ready state", ibpID)
			return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Partition %v is not in Ready state", ibpID), nil)
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
		return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("InfiniBand interfaces validation failed: %v", err), err)
	}

	return dbibic, nil
}

// ValidateNVLinkInterfaces validates NVLink interfaces against instance type capabilities,
// validates logical partition ownership and state.
// If no NVLink interfaces are specified but a default NVLink logical partition exists on the VPC,
// it generates default interfaces from the GPU capabilities.
func (icv *InstanceCreateValidator) ValidateNVLinkInterfaces(ctx context.Context, tenant *cdbm.Tenant, site *cdbm.Site, defaultNvllpID *uuid.UUID, nvlIfcs []cam.APINVLinkInterfaceCreateOrUpdateRequest, instanceTypeID uuid.UUID) ([]cdbm.NVLinkInterface, *cerr.APIError) {
	logger := icv.Logger

	mcDAO := cdbm.NewMachineCapabilityDAO(icv.DBSession)
	itNvlCaps, itNvlCapCount, err := mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{instanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeGPU), nil, nil, nil, nil, nil, cdb.GetStrPtr(cdbm.MachineCapabilityDeviceTypeNVLink), nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving GPU (NVLink) Machine Capabilities from DB for Instance Type")
		return nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve GPU Capabilities for Instance Type, DB error", nil)
	}

	if len(nvlIfcs) > 0 {
		if itNvlCapCount == 0 {
			logger.Warn().Msg("NVLink interfaces specified but Instance Type doesn't have GPU (NVLink) Capability")
			return nil, cerr.NewAPIError(http.StatusBadRequest, "NVLink Interfaces cannot be specified if Instance Type doesn't have GPU Capabilities", nil)
		}

		err = cam.ValidateNVLinkInterfaces(itNvlCaps, nvlIfcs)
		if err != nil {
			logger.Error().Err(err).Msg("NVLink interfaces validation failed")
			return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink interfaces validation failed: %v", err), err)
		}

		nvllpIDs := make([]uuid.UUID, 0, len(nvlIfcs))
		for _, nvlifc := range nvlIfcs {
			nvllpID, err := uuid.Parse(nvlifc.NVLinkLogicalPartitionID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing NVLink Logical Partition id")
				return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition ID %v is not valid", nvlifc.NVLinkLogicalPartitionID), nil)
			}
			nvllpIDs = append(nvllpIDs, nvllpID)
		}

		var nvllpMap map[uuid.UUID]*cdbm.NVLinkLogicalPartition
		if defaultNvllpID == nil {
			nvllpDAO := cdbm.NewNVLinkLogicalPartitionDAO(icv.DBSession)
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
				return nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve NVLink Logical Partitions specified in request data, DB error", nil)
			}

			nvllpMap = make(map[uuid.UUID]*cdbm.NVLinkLogicalPartition, len(nvllpList))
			for i := range nvllpList {
				nvllpMap[nvllpList[i].ID] = &nvllpList[i]
			}
		}

		dbnvlic := []cdbm.NVLinkInterface{}
		for i, nvlifc := range nvlIfcs {
			nvllpID := nvllpIDs[i]

			if defaultNvllpID != nil {
				if nvllpID != *defaultNvllpID {
					return nil, cerr.NewAPIError(http.StatusBadRequest, "NVLink Logical Partition specified for NVLink Interface does not match NVLink Logical Partition of VPC", nil)
				}
			} else {
				nvllp, exists := nvllpMap[nvllpID]
				if !exists {
					return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("Could not find NVLink Logical Partition with ID %v", nvllpID), nil)
				}

				if nvllp.SiteID != site.ID {
					logger.Warn().Msgf("NVLink Logical Partition: %v does not match with Instance Site", nvllpID)
					return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition %v does not match with Instance Site", nvllpID), nil)
				}

				if nvllp.TenantID != tenant.ID {
					logger.Warn().Msgf("NVLink Logical Partition: %v is not owned by Tenant", nvllpID)
					return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition %v is not owned by Tenant", nvllpID), nil)
				}

				if nvllp.Status != cdbm.NVLinkLogicalPartitionStatusReady {
					logger.Warn().Msgf("NVLink Logical Partition: %v is not in Ready state", nvllpID)
					return nil, cerr.NewAPIError(http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition %v is not in Ready state", nvllpID), nil)
				}
			}

			dbnvlic = append(dbnvlic, cdbm.NVLinkInterface{
				NVLinkLogicalPartitionID: nvllpID,
				DeviceInstance:           nvlifc.DeviceInstance,
			})
		}
		return dbnvlic, nil
	} else if defaultNvllpID != nil {
		dbnvlic := []cdbm.NVLinkInterface{}
		for _, nvlCap := range itNvlCaps {
			if nvlCap.Count != nil {
				for i := 0; i < *nvlCap.Count; i++ {
					dbnvlic = append(dbnvlic, cdbm.NVLinkInterface{
						NVLinkLogicalPartitionID: *defaultNvllpID,
						Device:                   cdb.GetStrPtr(nvlCap.Name),
						DeviceInstance:           i,
					})
				}
			}
		}
		return dbnvlic, nil
	}

	return nil, nil
}

// ValidateDPUCapabilities validates DPU interface capabilities against instance type machine capabilities.
func (icv *InstanceCreateValidator) ValidateDPUCapabilities(ctx context.Context, dbInterfaces []cdbm.Interface, isDeviceInfoPresent bool, instanceTypeID uuid.UUID) *cerr.APIError {
	if !isDeviceInfoPresent {
		return nil
	}

	logger := icv.Logger

	mcDAO := cdbm.NewMachineCapabilityDAO(icv.DBSession)
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

	err = cam.ValidateMultiEthernetDeviceInterfaces(itDpuCaps, dbInterfaces)
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
