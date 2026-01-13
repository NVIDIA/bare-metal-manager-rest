/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	goset "github.com/deckarep/golang-set/v2"
	"github.com/labstack/echo/v4"

	"go.opentelemetry.io/otel/attribute"
	temporalClient "go.temporal.io/sdk/client"
	tp "go.temporal.io/sdk/temporal"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/api/internal/config"
	common "github.com/nvidia/carbide-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
	"github.com/nvidia/carbide-rest/api/pkg/api/pagination"
	sc "github.com/nvidia/carbide-rest/api/pkg/client/site"
	auth "github.com/nvidia/carbide-rest/auth/pkg/authorization"
	cerr "github.com/nvidia/carbide-rest/common/pkg/util"
	sutil "github.com/nvidia/carbide-rest/common/pkg/util"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
	cdbp "github.com/nvidia/carbide-rest/db/pkg/db/paginator"
	swe "github.com/nvidia/carbide-rest/site-workflow/pkg/error"
	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"github.com/nvidia/carbide-rest/workflow/pkg/queue"
)

// ~~~~~ Create Handler ~~~~~ //

// CreateInstanceHandler is the API Handler for creating new Instance
type CreateInstanceHandler struct {
	dbSession  *cdb.Session
	tc         temporalClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewCreateInstanceHandler initializes and returns a new handler for creating Instance
func NewCreateInstanceHandler(dbSession *cdb.Session, tc temporalClient.Client, scp *sc.ClientPool, cfg *config.Config) CreateInstanceHandler {
	return CreateInstanceHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Returns either a default OS or an existing instance OS config.
// apiRequest will be mutated for use in createFromParams.
// osConfig will hold the struct/data for use with Temporal/Carbide calls.
// Errors should be returned in the form of cerr.NewAPIErrorResponse
func (cih CreateInstanceHandler) buildInstanceCreateRequestOsConfig(c echo.Context, logger *zerolog.Logger, apiRequest *model.APIInstanceCreateRequest, siteID uuid.UUID) (*cwssaws.OperatingSystem, *uuid.UUID, *cerr.APIError) {

	ctx := c.Request().Context()

	// If no OS was selected
	if apiRequest.OperatingSystemID == nil || *apiRequest.OperatingSystemID == "" {

		if err := apiRequest.ValidateAndSetOperatingSystemData(cih.cfg, nil); err != nil {
			logger.Error().Err(err).Msg("failed to validate OperatingSystem")
			return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Failed to validate OperatingSystem data", err)
		}

		return &cwssaws.OperatingSystem{
			RunProvisioningInstructionsOnEveryBoot: *apiRequest.AlwaysBootWithCustomIpxe, // Set by the earlier call to ValidateAndSetOperatingSystemData
			PhoneHomeEnabled:                       *apiRequest.PhoneHomeEnabled,         // Set by the earlier call to ValidateAndSetOperatingSystemData
			Variant: &cwssaws.OperatingSystem_Ipxe{
				Ipxe: &cwssaws.IpxeOperatingSystem{
					IpxeScript: *apiRequest.IpxeScript,
					UserData:   apiRequest.UserData,
				},
			},
		}, nil, nil
	}

	// Otherwise, we'll use the OS sent by the caller

	var id uuid.UUID
	var err error

	if id, err = uuid.Parse(*apiRequest.OperatingSystemID); err != nil {
		logger.Error().Err(err).Msg("failed to parse OperatingSystemID")
		return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Unable to parse `operatingSystemId` specified", validation.Errors{
			"operatingSystemId": errors.New(*apiRequest.OperatingSystemID),
		})
	}

	osID := &id

	// Retrieve the details for the OS
	osDAO := cdbm.NewOperatingSystemDAO(cih.dbSession)
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

	// Add the OS ID to the log fields now that we know we have a valid one.
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("OperatingSystem ID", os.ID.String())
	})

	// Confirm ownership between tenant and OS.
	if os.TenantID.String() != apiRequest.TenantID {
		logger.Error().Msg("OperatingSystem in request is not owned by tenant")
		return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "OperatingSystem specified in request is not owned by Tenant", nil)
	}

	// Confirm match between site and OS (only for Image type).
	if os.Type == cdbm.OperatingSystemTypeImage {
		ossaDAO := cdbm.NewOperatingSystemSiteAssociationDAO(cih.dbSession)
		_, ossaCount, err := ossaDAO.GetAll(
			ctx,
			nil,
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

	// Validate any additional properties.
	// `os` could still be nil here if no OS ID was sent
	// in the request.

	err = apiRequest.ValidateAndSetOperatingSystemData(cih.cfg, os)
	if err != nil {
		logger.Error().Msgf("OperatingSystem options validation failed: %s", err)
		return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "OperatingSystem options validation failed", err)
	}

	// Options below should all have been set by the
	// earlier call to ValidateAndSetOperatingSystemData

	if os.Type == cdbm.OperatingSystemTypeIPXE {
		return &cwssaws.OperatingSystem{
			RunProvisioningInstructionsOnEveryBoot: *apiRequest.AlwaysBootWithCustomIpxe,
			PhoneHomeEnabled:                       *apiRequest.PhoneHomeEnabled,
			Variant: &cwssaws.OperatingSystem_Ipxe{
				Ipxe: &cwssaws.IpxeOperatingSystem{
					IpxeScript: *apiRequest.IpxeScript,
					UserData:   apiRequest.UserData,
				},
			},
		}, osID, nil
	} else {
		return &cwssaws.OperatingSystem{
			PhoneHomeEnabled: *apiRequest.PhoneHomeEnabled,
			Variant: &cwssaws.OperatingSystem_OsImageId{
				OsImageId: &cwssaws.UUID{
					Value: os.ID.String(),
				},
			},
			UserData: apiRequest.UserData,
		}, osID, nil
	}
}

// Handle godoc
// @Summary Create an Instance
// @Description Create an Instance for the org.
// @Tags Instance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param message body model.APIInstanceCreateRequest true "Instance create request"
// @Success 201 {object} model.APIInstance
// @Router /v2/org/{org}/carbide/instance [post]
func (cih CreateInstanceHandler) Handle(c echo.Context) error {
	// Execution Steps:
	// 1. Authentication & Authorization
	//    - Extract user from context
	//    - Validate org membership
	//    - Validate Tenant Admin role
	// 2. Request Validation
	//    - Bind and validate request data
	//    - Validate tenant, VPC, site
	//    - Load and validate Interfaces (Subnets, VPC Prefixes)
	//    - Load and validate DPU Extension Service Deployments
	//    - Load and validate Network Security Groups
	//    - Load and validate SSH Key Groups
	//    - Validate OS or iPXE script
	//    - Check instance name uniqueness
	// 3. Database Transaction
	//    - Start transaction
	// 4. Machine Selection
	//    - Path A: Machine ID specified → validate and assign specific machine
	//    - Path B: Instance Type ID specified → acquire advisory lock, verify allocation constraints, find available machine
	// 5. Machine Capability Validation
	//    - Validate InfiniBand interfaces against Instance Type capabilities
	//    - Validate InfiniBand partitions (Site, Tenant, Status)
	//    - Validate DPU interfaces against Instance Type capabilities
	//    - Validate NVLink interfaces against Instance Type capabilities
	//    - Validate NVLink logical partitions (Site, Tenant, Status)
	// 6. Create Instance Records
	//    - Create Instance record
	//    - Update ControllerInstanceID
	//    - Create SSH Key Group associations
	//    - Create Interface records
	//    - Create InfiniBand Interface records
	//    - Create NVLink Interface records
	//    - Create DPU Extension Service Deployment records
	//    - Create status detail record
	// 7. Workflow Trigger
	//    - Build instance allocation request with all configs
	//    - Execute synchronous Temporal workflow (CreateInstanceV2)
	//    - Wait for site-agent to provision the instance
	//    - Handle timeout with workflow termination
	// 8. Commit & Response
	//    - Commit transaction after workflow succeeds
	//    - Return created instance to client

	// ==================== Step 1: Authentication & Authorization ====================

	// Get context
	ctx := c.Request().Context()

	// Get org
	org := c.Param("orgName")

	// Initialize logger
	logger := log.With().Str("Model", "Instance").Str("Handler", "Create").Str("Org", org).Logger()

	logger.Info().Msg("started API handler")

	// Create a child span and set the attributes for current request
	newctx, handlerSpan := cih.tracerSpan.CreateChildInContext(ctx, "CreateInstanceHandler", logger)
	if handlerSpan != nil {
		// Set newly created span context as a current context
		ctx = newctx

		defer handlerSpan.End()

		cih.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, cih.tracerSpan, handlerSpan)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Tenant Admins are allowed to create Instances
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.TenantAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Tenant Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Tenant Admin role with org", nil)
	}

	// ==================== Step 2: Request Validation ====================

	// Validate request
	// Bind request data to API model
	apiRequest := model.APIInstanceCreateRequest{}
	err = c.Bind(&apiRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data, potentially invalid structure", nil)
	}

	// Validate request attributes
	verr := apiRequest.Validate()
	if verr != nil {
		logger.Warn().Err(verr).Msg("error validating Instance creation request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Error validating Instance creation request data", verr)
	}

	// Validate the tenant for which this Instance is being created
	tenant, err := common.GetTenantForOrg(ctx, nil, cih.dbSession, org)
	if err != nil {
		if err == common.ErrOrgTenantNotFound {
			logger.Warn().Err(err).Msg("Org does not have a Tenant associated")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Org does not have a Tenant associated", nil)
		}
		logger.Error().Err(err).Msg("unable to retrieve tenant for org")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve tenant for org", nil)
	}

	// verify tenant-id in request, the api validation ensures non-nil tenantID in request
	apiTenant, err := common.GetTenantFromIDString(ctx, nil, apiRequest.TenantID, cih.dbSession)
	if err != nil {
		logger.Warn().Err(err).Msg("error retrieving tenant from request")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "TenantID in request is not valid", nil)
	}
	if apiTenant.ID != tenant.ID {
		logger.Warn().Msg("tenant id in request does not match tenant in org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "TenantID in request does not match tenant in org", nil)
	}

	// Validate the VPC state
	vpc, err := common.GetVpcFromIDString(ctx, nil, apiRequest.VpcID, []string{cdbm.NVLinkLogicalPartitionRelationName}, cih.dbSession)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not find VPC with ID specified in request data", nil)
		}
		logger.Warn().Err(err).Msg("error retrieving VPC from request")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "VpcID in request is not valid", nil)
	}

	if vpc.TenantID != tenant.ID {
		logger.Warn().Msg("tenant id in request does not match tenant in VPC")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "VPC specified in request is not owned by Tenant", nil)
	}

	if vpc.ControllerVpcID == nil || vpc.Status != cdbm.VpcStatusReady {
		logger.Warn().Msg("VPC specified in request data is not ready")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "VPC specified in request data is not ready", nil)
	}

	var defaultNvllpID *uuid.UUID
	if vpc.NVLinkLogicalPartitionID != nil {
		// NOTE: No validation needed here because the VPC validation ensures the NVLink Logical Partition is valid for this instance
		defaultNvllpID = vpc.NVLinkLogicalPartitionID
	}

	// Verify if site is ready
	stDAO := cdbm.NewSiteDAO(cih.dbSession)
	site, err := stDAO.GetByID(ctx, nil, vpc.SiteID, nil, false)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "The Site where this Instance is being created could not be found", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB by ID")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "The Site where this Instance is being created could not be retrieved", nil)
	}

	if site.Status != cdbm.SiteStatusRegistered {
		logger.Warn().Msg(fmt.Sprintf("The Site: %v where this Instance is being created is not in Registered state", vpc.SiteID.String()))
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "The Site where this Instance is being created is not in Registered state", nil)
	}

	// validate the instance subnet information to create instance subnet records later
	// Verify if subnet is ready
	sbDAO := cdbm.NewSubnetDAO(cih.dbSession)
	vpDAO := cdbm.NewVpcPrefixDAO(cih.dbSession)
	dbifcs := []cdbm.Interface{}

	// We'll need this later for grabbing network segments
	// to send in the carbide request.
	subnets := map[uuid.UUID]*cdbm.Subnet{}
	vpcPrefixes := map[uuid.UUID]*cdbm.VpcPrefix{}
	isDeviceInfoPresent := false

	for _, ifc := range apiRequest.Interfaces {
		if ifc.SubnetID != nil {
			subnetID, err := uuid.Parse(*ifc.SubnetID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing subnet id in instance subnet request")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Subnet ID specified in request data is not valid", nil)
			}

			if subnets[subnetID] == nil {
				subnet, err := sbDAO.GetByID(ctx, nil, subnetID, nil)
				if err != nil {
					if err == cdb.ErrDoesNotExist {
						return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not find Subnet with ID specified in request data", nil)
					}
					logger.Error().Err(err).Msg("error retrieving Subnet from DB by ID")
					return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Subnet with ID specified in request data", nil)
				}

				if subnet.TenantID != tenant.ID {
					logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request is not owned by Tenant", subnetID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request is not owned by Tenant", subnetID), nil)
				}

				if subnet.ControllerNetworkSegmentID == nil || subnet.Status != cdbm.SubnetStatusReady {
					logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request data is not in Ready state", subnetID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request data is not in Ready state", subnetID), nil)
				}

				if subnet.VpcID != vpc.ID {
					logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request does not match with VPC", subnetID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request does not match with VPC", subnetID), nil)
				}

				if vpc.NetworkVirtualizationType != nil && *vpc.NetworkVirtualizationType != cdbm.VpcEthernetVirtualizer {
					logger.Warn().Msg(fmt.Sprintf("VPC: %v specified in request must have Ethernet network virtualization type in order to create Subnet based interfaces", vpc.ID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC: %v specified in request must have Ethernet network virtualization type in order to create Subnet based interfaces", vpc.ID), nil)
				}

				subnets[subnetID] = subnet
			}
			dbifcs = append(dbifcs, cdbm.Interface{SubnetID: &subnetID, IsPhysical: ifc.IsPhysical, Status: cdbm.InterfaceStatusPending})
		}

		if ifc.VpcPrefixID != nil {
			vpcPrefixID, err := uuid.Parse(*ifc.VpcPrefixID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing vpcprefix id in instance vpcprefix request")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "VPC Prefix ID specified in request data is not valid", nil)
			}

			if vpcPrefixes[vpcPrefixID] == nil {
				vpcPrefix, err := vpDAO.GetByID(ctx, nil, vpcPrefixID, nil)
				if err != nil {
					if err == cdb.ErrDoesNotExist {
						return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not find VPC Prefix with ID specified in request data", nil)
					}
					logger.Error().Err(err).Msg("error retrieving vpcprefix from DB by ID")
					return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve VPC Prefix with ID specified in request data", nil)
				}

				if vpcPrefix.TenantID != tenant.ID {
					logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request is not owned by Tenant", vpcPrefixID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request is not owned by Tenant", vpcPrefixID), nil)
				}

				if vpcPrefix.Status != cdbm.VpcPrefixStatusReady {
					logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request data is not in Ready state", vpcPrefixID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request data is not in Ready state", vpcPrefixID), nil)
				}

				if vpcPrefix.VpcID != vpc.ID {
					logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request does not match with VPC", vpcPrefixID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request does not match with VPC", vpcPrefixID), nil)
				}

				if vpc.NetworkVirtualizationType == nil || *vpc.NetworkVirtualizationType != cdbm.VpcFNN {
					logger.Warn().Msg(fmt.Sprintf("VPC: %v specified in request must have FNN network virtualization type in order to create VPC Prefix based interfaces", vpc.ID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC: %v specified in request must have FNN network virtualization type in order to create VPC Prefix based interfaces", vpc.ID), nil)
				}

				vpcPrefixes[vpcPrefixID] = vpcPrefix
			}

			if ifc.Device != nil && ifc.DeviceInstance != nil {
				isDeviceInfoPresent = true
			}

			dbifcs = append(dbifcs, cdbm.Interface{
				VpcPrefixID:       &vpcPrefixID,
				Device:            ifc.Device,
				DeviceInstance:    ifc.DeviceInstance,
				VirtualFunctionID: ifc.VirtualFunctionID,
				IsPhysical:        ifc.IsPhysical,
				Status:            cdbm.InterfaceStatusPending})
		}
	}

	// Validate the DPU Extension Service Deployments
	desDAO := cdbm.NewDpuExtensionServiceDAO(cih.dbSession)
	desIDMap := map[string]*cdbm.DpuExtensionService{}
	for _, adesdr := range apiRequest.DpuExtensionServiceDeployments {
		desID, err := uuid.Parse(adesdr.DpuExtensionServiceID)
		if err != nil {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid DPU Extension Service ID: %s specified in request", adesdr.DpuExtensionServiceID), nil)
		}

		des, err := desDAO.GetByID(ctx, nil, desID, nil)
		if err != nil {
			if err == cdb.ErrDoesNotExist {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Could not find DPU Extension Service with ID: %s", desID), nil)
			}

			logger.Error().Err(err).Str("DPU Extension Service ID", desID.String()).Msg("error retrieving DPU Extension Service from DB by ID")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve DPU Extension Service specified in request, DB error", nil)
		}

		if des.TenantID != tenant.ID {
			logger.Warn().Str("Tenant ID", tenant.ID.String()).Str("DPU Extension Service ID", desID.String()).Msg("DPU Extension Service does not belong to current Tenant")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("DPU Extension Service: %s does not belong to current Tenant", desID.String()), nil)
		}

		if des.SiteID != site.ID {
			logger.Warn().Str("Site ID", site.ID.String()).Str("DPU Extension Service ID", desID.String()).Msg("DPU Extension Service does not belong to Site")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("DPU Extension Service: %s does not belong to Site where Instance is being created", desID.String()), nil)
		}

		versionFound := false
		for _, version := range des.ActiveVersions {
			if version == adesdr.Version {
				versionFound = true
				break
			}
		}
		if !versionFound {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Version: %s was not found for DPU Extension Service: %s", adesdr.Version, desID.String()), nil)
		}

		desIDMap[desID.String()] = des
	}

	// If an NSG was requested, validate it
	if apiRequest.NetworkSecurityGroupID != nil {
		nsgDAO := cdbm.NewNetworkSecurityGroupDAO(cih.dbSession)

		nsg, err := nsgDAO.GetByID(ctx, nil, *apiRequest.NetworkSecurityGroupID, nil)
		if err != nil {
			if err == cdb.ErrDoesNotExist {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Could not find NetworkSecurityGroup with ID: %s specified in request", *apiRequest.NetworkSecurityGroupID), nil)
			}

			logger.Error().Err(err).Msg("error retrieving NetworkSecurityGroup with ID specified in request data")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve NetworkSecurityGroup with ID specified in request data", nil)
		}

		if nsg.SiteID != site.ID {
			logger.Error().Msg("NetworkSecurityGroup in request does not belong to Site")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "NetworkSecurityGroup with ID specified in request data does not belong to Site", nil)
		}

		if nsg.TenantID != tenant.ID {
			logger.Error().Msg("NetworkSecurityGroup in request does not belong to Tenant")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "NetworkSecurityGroup with ID specified in request data does not belong to Tenant", nil)
		}
	}

	// Verify or validate SSH Key Group
	var rdbskg []cdbm.SSHKeyGroup
	skgsaDAO := cdbm.NewSSHKeyGroupSiteAssociationDAO(cih.dbSession)
	for _, skgID := range apiRequest.SSHKeyGroupIDs {
		// Validate the SSH Key for which this SSH Key Group is being associated
		sshkeygroup, serr := common.GetSSHKeyGroupFromIDString(ctx, nil, skgID, cih.dbSession, nil)
		if serr != nil {
			if serr == common.ErrInvalidID {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to create Instance, Invalid SSH Key Group ID: %s", skgID), nil)
			}
			if serr == cdb.ErrDoesNotExist {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to create Instance, Could not find SSH Key Group with ID: %s ", skgID), nil)
			}

			logger.Warn().Err(serr).Str("SSH Key Group ID", skgID).Msg("error retrieving SSH Key Group from DB by ID")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to retrieve SSH Key Group with ID `%s`specified in request, DB error", skgID), nil)
		}

		if sshkeygroup.TenantID != tenant.ID {
			logger.Warn().Str("Tenant ID", tenant.ID.String()).Str("SSH Key Group ID", skgID).Msg("SSH Key Group does not belong to current Tenant")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to create Instance, SSH Key Group with ID: %s does not belong to Tenant", skgID), nil)
		}

		// Verify if SSH Key Group Site Association exists
		_, serr = skgsaDAO.GetBySSHKeyGroupIDAndSiteID(ctx, nil, sshkeygroup.ID, site.ID, nil)
		if serr != nil {
			if serr == cdb.ErrDoesNotExist {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("SSH Key Group with ID: %s is not associated with the Site where Instance is being created", skgID), nil)
			}
			logger.Warn().Err(serr).Str("SSH Key Group ID", skgID).Msg("error retrieving SSH Key Group Site Association from DB by SSH Key Group ID & Site ID")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to determine if SSH Key Group: %s is associated with the Site where Instance is being created, DB error", skgID), nil)
		}

		rdbskg = append(rdbskg, *sshkeygroup)
	}

	// apiRequest will be mutated for use in CreateFromParams.
	// osConfig will hold the struct/data for use with Temporal/Carbide calls.
	// Errors will be returned already in the form of cerr.NewAPIErrorResponse
	osConfig, osID, oserr := cih.buildInstanceCreateRequestOsConfig(c, &logger, &apiRequest, vpc.SiteID)
	if oserr != nil {
		// buildInstanceCreateRequestOsConfig already handles logging,
		// so this is a bit redundant, but this log brings you to the
		// actual call site.  I think buildInstanceCreateRequestOsConfig
		// would ideally return only `error` and let the logging and
		// and cerr.NewAPIErrorResponse(...) happen here, but we
		// have at least one StatusInternalServerError case that would
		// be hidden if we merge it all under StatusBadRequest here.
		logger.Error().Err(errors.New(oserr.Message)).Msg("error building os config for creating Instance")
		return c.JSON(oserr.Code, oserr)
	}

	// ensure we have one and only one of InstanceTypeID or MachineID
	if apiRequest.InstanceTypeID != nil && apiRequest.MachineID != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "request can only have InstanceType ID or Machine ID and not both", nil)
	} else if apiRequest.InstanceTypeID == nil && apiRequest.MachineID == nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "request must have InstanceType ID or Machine ID", nil)
	}

	// common pre-requisites for both InstanceType and Machine ID cases
	var (
		instanceTypeID *uuid.UUID
		machine        *cdbm.Machine
		dbibic         []cdbm.InfiniBandInterface
		dbnvlic        []cdbm.NVLinkInterface
	)
	// allocation constraint of the provided instancetype and tenant site allocation (nil when machine id is provided)
	var currentAllocationConstraint *cdbm.AllocationConstraint

	inDAO := cdbm.NewInstanceDAO(cih.dbSession)

	// Check for name uniqueness for the tenant, ie, tenant cannot have another instance with same name at the site
	// TODO consider doing this with an advisory lock for correctness
	ins, tot, err := inDAO.GetAll(ctx, nil, cdbm.InstanceFilterInput{Names: []string{apiRequest.Name}, TenantIDs: []uuid.UUID{tenant.ID}, SiteIDs: []uuid.UUID{vpc.SiteID}}, cdbp.PageInput{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("db error checking for name uniqueness of tenant instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Instance due to DB error", nil)
	}
	if tot > 0 {
		logger.Warn().Str("tenantId", tenant.ID.String()).Str("name", apiRequest.Name).Msg("instance with same name already exists for tenant")
		return cerr.NewAPIErrorResponse(c, http.StatusConflict, "An Instance with specified name already exists for Tenant", validation.Errors{
			"id": errors.New(ins[0].ID.String()),
		})
	}

	// ==================== Step 3: Database Transaction ====================

	// Start a db tx
	tx, err := cdb.BeginTx(ctx, cih.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Instance", nil)
	}
	// this variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// ==================== Step 4: Machine Selection  ====================

	// if requested a specific machine ID:
	var allowUnhealthyMachine bool
	if apiRequest.MachineID != nil {
		if tenant.Config == nil || !tenant.Config.TargetedInstanceCreation {
			logger.Warn().Msg("tenant does not have capability to create instances from specific machine")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Tenant does not have capability to create instances from specific machine", nil)
		}

		mDAO := cdbm.NewMachineDAO(cih.dbSession)

		// retrieve specified machine
		machine, err = mDAO.GetByID(ctx, nil, *apiRequest.MachineID, nil, false)
		if err != nil {
			if err == cdb.ErrDoesNotExist {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not find Machine with ID specified in request data", nil)
			}
			logger.Error().Err(err).Msg("error retrieving Machine from DB by ID")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Machine with ID specified in request data", nil)
		}
		// validate that the machine is part of the site
		if machine.SiteID != site.ID {
			logger.Warn().Msg("Machine specified in request is not part of the site")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Machine specified in request is not part of the site", nil)
		}
		// validate that the machine is not missing
		if machine.IsMissingOnSite {
			logger.Warn().Msg("Machine specified in request is missing on site")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Machine specified in request is missing on site", nil)
		}
		// validate that machine is not in use
		if machine.IsAssigned {
			logger.Warn().Msg("Machine specified in request is already in use")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Machine specified in request is already in use", nil)
		}

		// ensure specified machine is healthy OR caller allows unhealthy machines
		allowUnhealthyMachine = apiRequest.AllowUnhealthyMachine != nil && *apiRequest.AllowUnhealthyMachine
		if machine.Status == cdbm.MachineStatusError && !allowUnhealthyMachine {
			logger.Warn().Msg("Machine specified in request is not healthy")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Machine specified in request is not healthy", nil)
		}

		// Update the machine status to assigned
		updateInput := cdbm.MachineUpdateInput{
			MachineID:  machine.ID,
			IsAssigned: cdb.GetBoolPtr(true),
		}
		machine, err = mDAO.Update(ctx, tx, updateInput)
		if err != nil {
			if err == cdb.ErrDoesNotExist {
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Could not find Machine with ID specified in request data for update", nil)
			}
			logger.Error().Err(err).Msg("error retrieving Machine from DB by ID")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Machine with ID specified in request data", nil)
		}

		instanceTypeID = machine.InstanceTypeID
	} // if apiRequest.MachineID != nil

	// if we only have an Instance ID then we need to find one machine from that Instance Type pool:
	if apiRequest.InstanceTypeID != nil {
		// Validate the instance type
		apiInstanceTypeID, err := uuid.Parse(*apiRequest.InstanceTypeID)
		if err != nil {
			logger.Warn().Err(err).Msg("error parsing instance type id in request")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Instance Type ID in request is not valid", nil)
		}
		instanceTypeID = &apiInstanceTypeID

		itDAO := cdbm.NewInstanceTypeDAO(cih.dbSession)
		instancetype, err := itDAO.GetByID(ctx, nil, *instanceTypeID, nil)
		if err != nil {
			if err == cdb.ErrDoesNotExist {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not find Instance Type with ID specified in request data", nil)
			}
			logger.Error().Err(err).Msg("error retrieving Instance Type from DB by ID")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance Type with ID specified in request data", nil)
		}

		// acquire an advisory lock on the tenant ID and instancetype ID on which instance is being creating
		// this lock is released when the transaction commits or rolls back
		err = tx.TryAcquireAdvisoryLock(ctx, cdb.GetAdvisoryLockIDFromString(fmt.Sprintf("%s-%s", tenant.ID.String(), instancetype.ID.String())), nil)
		if err != nil {
			// TODO add a retry here
			logger.Error().Err(err).Msg("Failed to acquire advisory lock on Tenant and Instance Type")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Error creating Instance, detected multiple parallel request on Instance Type by Tenant", nil)
		}

		// Ensure that Tenant has an Allocation with specified Tenant InstanceType Site
		aDAO := cdbm.NewAllocationDAO(cih.dbSession)
		allocationFilter := cdbm.AllocationFilterInput{TenantIDs: []uuid.UUID{tenant.ID}, SiteIDs: []uuid.UUID{*instancetype.SiteID}}
		allocationPage := cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}
		tnas, _, serr := aDAO.GetAll(ctx, tx, allocationFilter, allocationPage, nil)
		if serr != nil {
			logger.Error().Err(serr).Msg("error retrieving Allocations from DB for Tenant and Site")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Allocation for Tenant", nil)
		}
		if len(tnas) == 0 {
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden,
				"Tenant does not have any Allocations for Site and Instance Type specified in request data", nil)
		}

		alconstraints, err := common.GetAllocationConstraintsForInstanceType(ctx, tx, cih.dbSession, tenant.ID, instancetype, tnas)
		if err != nil {
			if err == common.ErrAllocationConstraintNotFound {
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "No Allocations for specified Instance Type were found for current Tenant", nil)
			}
			logger.Error().Err(err).Msg("error retrieving Allocation Constraints from DB for InstanceType and Allocation")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Allocations for specified Instance Type, DB error", nil)
		}

		// Getting active instances for the tenant on requested instance type
		var siteIDs []uuid.UUID
		if instancetype.SiteID != nil {
			siteIDs = []uuid.UUID{*instancetype.SiteID}
		}
		instances, insTotal, err := inDAO.GetAll(ctx, tx, cdbm.InstanceFilterInput{TenantIDs: []uuid.UUID{tenant.ID}, SiteIDs: siteIDs, InstanceTypeIDs: []uuid.UUID{instancetype.ID}}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Active Instances from DB for Tenant and InstanceType")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve active instances for Tenant and Instance Type, DB error", nil)
		}

		// Build map allocation constraint ID which has been used by Instance
		usedMapAllocationConstraintIDs := map[uuid.UUID]int{}
		for _, inst := range instances {
			// WARNING
			// TODO: Currently no instances can be created without a constraint ID
			// but that will be changing.  When it does, this will need to be handled differently.
			if inst.AllocationConstraintID == nil {
				logger.Error().Msgf("found Instance missing AllocationConstraintID")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Instance is missing Allocation Constraint ID", nil)
			}
			usedMapAllocationConstraintIDs[*inst.AllocationConstraintID] += 1
		}

		// Calculate total constraint value
		totalConstraintValue := 0
		for _, alcs := range alconstraints {
			totalConstraintValue += alcs.ConstraintValue
		}

		// Allocation constraints
		if len(alconstraints) > 0 && insTotal >= totalConstraintValue {
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden,
				"Tenant has reached the maximum number of Instances for Instance Type specified in request data", nil)
		}

		// Validate the currently active instances of the requested instance type for the tenant with allocation constraints
		for _, alcs := range alconstraints {
			if usedMapAllocationConstraintIDs[alcs.ID] < alcs.ConstraintValue {
				currentAllocationConstraint = &alcs
				break
			}
		}

		// Allocation constraints
		if currentAllocationConstraint == nil {
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
				"Error determining available Allocation for Instance, potential data issue", nil)
		}

		// Select unallocated Machine for the requested instance type
		machine, err = common.GetUnallocatedMachineForInstanceType(ctx, tx, cih.dbSession, instancetype)
		if err != nil {
			if err == common.ErrInstanceTypeMachineNotFound {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest,
					"No Machines are available for specified Instance Type", nil)
			}
			logger.Error().Err(err).Msg("error retrieving Machine from DB for Instance Type")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve available baremetal Machines for specified Instance Type", nil)
		}
	} // if apiRequest.InstanceTypeID != nil

	// ==================== Step 5: Machine Capability Validation  ====================

	mcDAO := cdbm.NewMachineCapabilityDAO(cih.dbSession)

	// Actions needed when an instance type exists either coming directly from request (by instance type id) or because
	// a machine id was provided and the machine already has an instance type
	if instanceTypeID != nil {
		itIbCaps, itIbCapCount, err := mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{*instanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeInfiniBand), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Machine Capabilities from DB for Instance Type")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Machine Capabilities for Instance Type, DB error", nil)
		}

		if itIbCapCount == 0 && len(apiRequest.InfiniBandInterfaces) > 0 {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "InfiniBand Interfaces cannot be specified if Instance Type doesn't have InfiniBand Capability", nil)
		}

		ibpDAO := cdbm.NewInfiniBandPartitionDAO(cih.dbSession)
		for _, ibic := range apiRequest.InfiniBandInterfaces {
			// InfiniBand Partition
			ibpID, err := uuid.Parse(ibic.InfiniBandPartitionID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing infiniband partition id in instance infiniband interface request")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Partition ID: %v specified in request data is not valid", ibic.InfiniBandPartitionID), nil)
			}

			// Validate Instance infiniband interface information to create DB records later
			ibp, err := ibpDAO.GetByID(ctx, nil, ibpID, nil)
			if err != nil {
				if err == cdb.ErrDoesNotExist {
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Could find Partition with ID: %v specified in request data", ibic.InfiniBandPartitionID), nil)
				}
				logger.Error().Err(err).Msg("error retrieving InfiniBand Partition from DB by ID")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Partition with ID specified in request data, DB error", nil)
			}

			if ibp.SiteID != site.ID {
				logger.Warn().Msg(fmt.Sprintf("InfiniBandPartition: %v specified in request does not match with Instance Site", ibpID))
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Partition: %v specified in request does not match with Instance Site", ibpID), nil)
			}

			if ibp.TenantID != tenant.ID {
				logger.Warn().Msg(fmt.Sprintf("InfiniBandPartition: %v specified in request is not owned by Tenant", ibpID))
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Partition: %v specified in request is not owned by Tenant", ibpID), nil)
			}

			if ibp.ControllerIBPartitionID == nil || ibp.Status != cdbm.InfiniBandPartitionStatusReady {
				logger.Warn().Msg(fmt.Sprintf("InfiniBandPartition: %v specified in request data is not in Ready state", ibpID))
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Partition: %v specified in request data is not in Ready state", ibpID), nil)
			}

			dbibic = append(dbibic, cdbm.InfiniBandInterface{InfiniBandPartitionID: ibp.ID, Device: ibic.Device, Vendor: ibic.Vendor, DeviceInstance: ibic.DeviceInstance, IsPhysical: ibic.IsPhysical, VirtualFunctionID: ibic.VirtualFunctionID})
		}

		// Validate InfiniBand Interfaces if Instance Type has InfiniBand Capability
		err = apiRequest.ValidateInfiniBandInterfaces(itIbCaps)
		if err != nil {
			logger.Error().Msgf("InfiniBand interfaces validation failed: %s", err)
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "InfiniBand interfaces validation failed", err)
		}

		// Validate DPU Interfaces if Instance Type has Network Capability with DPU device type
		if isDeviceInfoPresent {
			itDpuCaps, itDpuCapCount, err := mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{*instanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeNetwork), nil, nil, nil, nil, nil, cdb.GetStrPtr(cdbm.MachineCapabilityDeviceTypeDPU), nil, nil, nil, nil, nil)
			if err != nil {
				logger.Error().Err(err).Msg("error retrieving Machine Capabilities from DB for Instance Type")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Machine Capabilities for Instance Type, DB error", nil)
			}

			if itDpuCapCount == 0 {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Device and Device Instance cannot be specified if Instance Type doesn't have Network Capabilities with DPU device type", nil)
			}

			// Validate DPU Interfaces if Instance Type DPU capability is present and matches with the request
			err = apiRequest.ValidateMultiEthernetDeviceInterfaces(itDpuCaps, dbifcs)
			if err != nil {
				logger.Error().Msgf("DPU interfaces validation failed: %s", err)
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "DPU interfaces validation failed", err)
			}
		}
	}

	// Get GPU Capabilities with NVLink device type
	var nvlCaps []cdbm.MachineCapability
	var nvlCapCount int
	if instanceTypeID != nil {
		nvlCaps, nvlCapCount, err = mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{*instanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeGPU), nil, nil, nil, nil, nil, cdb.GetStrPtr(cdbm.MachineCapabilityDeviceTypeNVLink), nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving GPU Machine Capabilities from DB for Instance Type")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve GPU Capabilities for Instance Type, DB error", nil)
		}
	} else {
		nvlCaps, nvlCapCount, err = mcDAO.GetAll(ctx, nil, []string{machine.ID}, nil, cdb.GetStrPtr(cdbm.MachineCapabilityTypeGPU), nil, nil, nil, nil, nil, cdb.GetStrPtr(cdbm.MachineCapabilityDeviceTypeNVLink), nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving GPU Machine Capabilities from DB for Machine")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve GPU Capabilities for Machine, DB error", nil)
		}
	}

	if len(apiRequest.NVLinkInterfaces) > 0 {
		if nvlCapCount == 0 {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "NVLink Interfaces cannot be specified if Instance Type doesn't have GPU Capabilities", nil)
		}

		// Validate NVLink interfaces if Instance Type has GPU Capability
		err = apiRequest.ValidateNVLinkInterfaces(nvlCaps)
		if err != nil {
			logger.Error().Msgf("NVLink interfaces validation failed: %s", err)
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate NVLink interfaces specified in request", err)
		}

		nvllpDAO := cdbm.NewNVLinkLogicalPartitionDAO(cih.dbSession)
		for _, nvlifc := range apiRequest.NVLinkInterfaces {
			// NVLink Logical Partition
			nvllpID, err := uuid.Parse(nvlifc.NVLinkLogicalPartitionID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing NVLink Logical Partition id in instance NVLink Interface request")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition ID: %v specified in request data is not valid", nvlifc.NVLinkLogicalPartitionID), nil)
			}

			// Validate that the NVLink Logical Partition ID matches the default NVLink Logical Partition ID
			if defaultNvllpID != nil {
				if nvllpID != *defaultNvllpID {
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "NVLink Logical Partition specified for NVLink Interface does not match NVLink Logical Partition of VPC", nil)
				}
			} else {
				// Validate NVLink Logical Partition only if it's not the default
				nvllp, err := nvllpDAO.GetByID(ctx, nil, nvllpID, nil)
				if err != nil {
					logger.Error().Err(err).Msg("error retrieving NVLink Logical Partition from DB by ID")
					return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve NVLink Logical Partition with ID specified in request data, DB error", nil)
				}

				if nvllp.SiteID != site.ID {
					logger.Warn().Msg(fmt.Sprintf("NVLink Logical Partition: %v specified in request does not match with Instance Site", nvllpID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition: %v specified in request does not match with Instance Site", nvllpID), nil)
				}

				if nvllp.TenantID != tenant.ID {
					logger.Warn().Msg(fmt.Sprintf("NVLink Logical Partition: %v specified in request data is not owned by Tenant", nvllpID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition: %v specified in request data is not owned by Tenant", nvllpID), nil)
				}

				if nvllp.Status != cdbm.NVLinkLogicalPartitionStatusReady {
					logger.Warn().Msg(fmt.Sprintf("NVLink Logical Partition: %v specified in request data is not in Ready state", nvllpID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition: %v specified in request data is not in Ready state", nvllpID), nil)
				}
			}

			// Validate Instance NVLink Interface information to create DB records later
			dbnvlic = append(dbnvlic, cdbm.NVLinkInterface{NVLinkLogicalPartitionID: nvllpID, DeviceInstance: nvlifc.DeviceInstance})
		}
	} else if defaultNvllpID != nil {
		// Generate Interfaces for the default NVLink Logical Partition
		// For a given Machine, all the GPUs should be connected to the same NVLink Logical Partition
		for _, nvlCap := range nvlCaps {
			if nvlCap.Count != nil {
				for i := 0; i < *nvlCap.Count; i++ {
					dbnvlic = append(dbnvlic, cdbm.NVLinkInterface{NVLinkLogicalPartitionID: *defaultNvllpID, Device: cdb.GetStrPtr(nvlCap.Name), DeviceInstance: i})
				}
			}
		}
	}

	// ==================== Step 6: Create Instance Records  ====================

	instanceCreateInput := cdbm.InstanceCreateInput{
		Name:                     apiRequest.Name,
		Description:              apiRequest.Description,
		TenantID:                 tenant.ID,
		InfrastructureProviderID: machine.InfrastructureProviderID,
		SiteID:                   machine.SiteID,
		VpcID:                    vpc.ID,
		MachineID:                cdb.GetStrPtr(machine.ID),
		OperatingSystemID:        osID,
		IpxeScript:               apiRequest.IpxeScript,
		AlwaysBootWithCustomIpxe: *apiRequest.AlwaysBootWithCustomIpxe,
		PhoneHomeEnabled:         *apiRequest.PhoneHomeEnabled,
		UserData:                 apiRequest.UserData,
		NetworkSecurityGroupID:   apiRequest.NetworkSecurityGroupID,
		Labels:                   apiRequest.Labels,
		IsUpdatePending:          false,
		Status:                   cdbm.InstanceStatusPending,
		PowerStatus:              cdb.GetStrPtr(cdbm.InstancePowerStatusRebooting),
		CreatedBy:                dbUser.ID,
	}

	// NOTE: Set InstanceTypeID only if it is provided in the request. For Instances created with an Instance Type ID, we expect Allocation information
	// to be present. Since Machine ID based Instance creation does not require Allocation information, setting InstanceTypeID will create data integrity issues.
	if apiRequest.InstanceTypeID != nil {
		instanceCreateInput.InstanceTypeID = instanceTypeID
		instanceCreateInput.AllocationID = &currentAllocationConstraint.AllocationID
		instanceCreateInput.AllocationConstraintID = &currentAllocationConstraint.ID
	}

	instance, err := inDAO.Create(ctx, tx, instanceCreateInput)
	if err != nil {
		logger.Error().Err(err).Msg("unable to create Instance record in DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed creating Instance record, DB error", nil)
	}

	// Update the controller ID
	// We need this to match the instance ID.  This was previously handled
	// by the async cloud workflow after successful creation on site.
	instance, err = inDAO.Update(ctx, tx, cdbm.InstanceUpdateInput{InstanceID: instance.ID, ControllerInstanceID: cdb.GetUUIDPtr(instance.ID)})
	if err != nil {
		logger.Error().Err(err).Msg("unable to update Instance record controllerInstanceID in DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed updating new Instance record, DB error", nil)
	}

	// We'll need a list of the IDs as a string slice to
	// send along in the config update request to carbide.
	instanceSshKeyGroupIds := []string{}

	// create the ssh key group instance association in the db
	skgiaDAO := cdbm.NewSSHKeyGroupInstanceAssociationDAO(cih.dbSession)
	for _, skg := range rdbskg {
		_, err := skgiaDAO.CreateFromParams(ctx, tx, skg.ID, site.ID, instance.ID, dbUser.ID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create the SSH Key Group Instance Association record in DB")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to associate one or more SSH Key Group with Instance, DB error", nil)
		}

		instanceSshKeyGroupIds = append(instanceSshKeyGroupIds, skg.ID.String())
	}

	// Prepare interface details to pass to carbide call
	interfaceConfigs := []*cwssaws.InstanceInterfaceConfig{}

	// Create the instance subnet record in the db from info gathered earlier
	// The first Subnet is automatically added to the physical interface
	ifcs := []cdbm.Interface{}
	ifcDAO := cdbm.NewInterfaceDAO(cih.dbSession)
	for _, dbifc := range dbifcs {
		input := cdbm.InterfaceCreateInput{
			InstanceID:        instance.ID,
			SubnetID:          dbifc.SubnetID,
			VpcPrefixID:       dbifc.VpcPrefixID,
			Device:            dbifc.Device,
			DeviceInstance:    dbifc.DeviceInstance,
			VirtualFunctionID: dbifc.VirtualFunctionID,
			IsPhysical:        dbifc.IsPhysical,
			Status:            dbifc.Status,
			CreatedBy:         dbUser.ID,
		}

		retifc, serr := ifcDAO.Create(ctx, tx, input)
		if serr != nil {
			logger.Error().Err(serr).Msg("error creating Instance Subnet DB entry")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Instance Subnet entry for Instance, DB error", nil)
		}

		ifc := *retifc
		ifcs = append(ifcs, ifc)

		interfaceConfig := &cwssaws.InstanceInterfaceConfig{
			FunctionType: cwssaws.InterfaceFunctionType_VIRTUAL_FUNCTION,
		}

		// Assign InstanceInterfaceConfig_SegmentId in case of Subnet
		if dbifc.SubnetID != nil {
			interfaceConfig.NetworkSegmentId = &cwssaws.NetworkSegmentId{
				Value: subnets[*dbifc.SubnetID].ControllerNetworkSegmentID.String(),
			}
			interfaceConfig.NetworkDetails = &cwssaws.InstanceInterfaceConfig_SegmentId{
				SegmentId: &cwssaws.NetworkSegmentId{
					Value: subnets[*dbifc.SubnetID].ControllerNetworkSegmentID.String(),
				},
			}
		}

		// Assign InstanceInterfaceConfig_VpcPrefixId in case of VpcPrefix
		if dbifc.VpcPrefixID != nil {
			interfaceConfig.NetworkDetails = &cwssaws.InstanceInterfaceConfig_VpcPrefixId{
				VpcPrefixId: &cwssaws.VpcPrefixId{Value: dbifc.VpcPrefixID.String()},
			}
		}

		if dbifc.IsPhysical {
			interfaceConfig.FunctionType = cwssaws.InterfaceFunctionType_PHYSICAL_FUNCTION
		}

		// Assign Device and DeviceInstance in case of Multi DPU Interface
		if dbifc.Device != nil && dbifc.DeviceInstance != nil {
			interfaceConfig.Device = dbifc.Device
			interfaceConfig.DeviceInstance = uint32(*dbifc.DeviceInstance)
		}

		if !dbifc.IsPhysical {
			if dbifc.VirtualFunctionID != nil {
				vfID := uint32(*dbifc.VirtualFunctionID)
				interfaceConfig.VirtualFunctionId = &vfID
			}
		}

		interfaceConfigs = append(interfaceConfigs, interfaceConfig)
	}

	//We'll need this later for the carbide call
	ibInterfaceConfigs := []*cwssaws.InstanceIBInterfaceConfig{}

	// Create the instance infiniband interface record in the db from info gathered earlier IF instance type was used
	ibifcs := []cdbm.InfiniBandInterface{}
	ibifcDAO := cdbm.NewInfiniBandInterfaceDAO(cih.dbSession)
	for _, ibifc := range dbibic {
		retibifc, serr := ibifcDAO.Create(
			ctx,
			tx,
			cdbm.InfiniBandInterfaceCreateInput{
				InstanceID:            instance.ID,
				SiteID:                site.ID,
				InfiniBandPartitionID: ibifc.InfiniBandPartitionID,
				Device:                ibifc.Device,
				Vendor:                ibifc.Vendor,
				DeviceInstance:        ibifc.DeviceInstance,
				IsPhysical:            ibifc.IsPhysical,
				VirtualFunctionID:     ibifc.VirtualFunctionID,
				Status:                cdbm.InfiniBandInterfaceStatusPending,
				CreatedBy:             dbUser.ID,
			},
		)
		if serr != nil {
			logger.Error().Err(serr).Msg("error creating Instance InfiniBand Interface DB entry")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Instance InfiniBand Interface entry for Instance, DB error", nil)
		}

		ifc := *retibifc
		ibifcs = append(ibifcs, ifc)

		ibInterfaceConfig := &cwssaws.InstanceIBInterfaceConfig{
			Device:         ifc.Device,
			Vendor:         ifc.Vendor,
			DeviceInstance: uint32(ifc.DeviceInstance),
			FunctionType:   cwssaws.InterfaceFunctionType_PHYSICAL_FUNCTION,
			IbPartitionId:  &cwssaws.IBPartitionId{Value: ifc.InfiniBandPartitionID.String()},
		}
		ibInterfaceConfigs = append(ibInterfaceConfigs, ibInterfaceConfig)

		if !ifc.IsPhysical {
			ibInterfaceConfig.FunctionType = cwssaws.InterfaceFunctionType_VIRTUAL_FUNCTION

			if ifc.VirtualFunctionID != nil {
				vfID := uint32(*ifc.VirtualFunctionID)
				ibInterfaceConfig.VirtualFunctionId = &vfID
			}
		}
	}

	// Create the instance NVLink Interface record in the db from info gathered earlier IF instance type was used
	nvlifcs := []cdbm.NVLinkInterface{}
	nvlifcDAO := cdbm.NewNVLinkInterfaceDAO(cih.dbSession)
	nvlInterfaceConfigs := []*cwssaws.InstanceNVLinkGpuConfig{}
	for _, nvlifc := range dbnvlic {
		retnvlifc, serr := nvlifcDAO.Create(
			ctx,
			tx,
			cdbm.NVLinkInterfaceCreateInput{
				InstanceID:               instance.ID,
				SiteID:                   site.ID,
				NVLinkLogicalPartitionID: nvlifc.NVLinkLogicalPartitionID,
				Device:                   nvlifc.Device,
				DeviceInstance:           nvlifc.DeviceInstance,
				Status:                   cdbm.NVLinkInterfaceStatusPending,
				CreatedBy:                dbUser.ID,
			},
		)
		if serr != nil {
			logger.Error().Err(serr).Msg("error creating Instance NVLink Interface DB entry")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Instance NVLink Interface entry for Instance, DB error", nil)
		}

		nvlfc := *retnvlifc
		nvlifcs = append(nvlifcs, nvlfc)

		nvlInterfaceConfig := &cwssaws.InstanceNVLinkGpuConfig{
			DeviceInstance:     uint32(nvlifc.DeviceInstance),
			LogicalPartitionId: &cwssaws.NVLinkLogicalPartitionId{Value: nvlfc.NVLinkLogicalPartitionID.String()},
		}
		nvlInterfaceConfigs = append(nvlInterfaceConfigs, nvlInterfaceConfig)
	}

	// Create the DpuExtensionServiceDeployment records in DB
	desdConfigs := []*cwssaws.InstanceDpuExtensionServiceConfig{}

	desdDAO := cdbm.NewDpuExtensionServiceDeploymentDAO(cih.dbSession)
	desds := []cdbm.DpuExtensionServiceDeployment{}

	for _, adesdr := range apiRequest.DpuExtensionServiceDeployments {
		desdID, err := uuid.Parse(adesdr.DpuExtensionServiceID)
		if err != nil {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid DPU Extension Service ID: %s specified in request", adesdr.DpuExtensionServiceID), nil)
		}

		desd, err := desdDAO.Create(ctx, tx, cdbm.DpuExtensionServiceDeploymentCreateInput{
			SiteID:                site.ID,
			TenantID:              tenant.ID,
			InstanceID:            instance.ID,
			DpuExtensionServiceID: desdID,
			Version:               adesdr.Version,
			Status:                cdbm.DpuExtensionServiceDeploymentStatusPending,
			CreatedBy:             dbUser.ID,
		})
		if err != nil {
			logger.Error().Err(err).Msg("error creating Instance DpuExtensionServiceDeployment record in DB")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create DPU Extension Service Deployment for Instance, DB error", nil)
		}

		des, _ := desIDMap[desdID.String()]
		desd.DpuExtensionService = des

		desds = append(desds, *desd)

		desdConfigs = append(desdConfigs, &cwssaws.InstanceDpuExtensionServiceConfig{
			ServiceId: desd.DpuExtensionServiceID.String(),
			Version:   desd.Version,
		})
	}

	// Create the status detail record
	sdDAO := cdbm.NewStatusDetailDAO(cih.dbSession)
	ssd, serr := sdDAO.CreateFromParams(ctx, tx, instance.ID.String(), *cdb.GetStrPtr(cdbm.InstanceStatusPending),
		cdb.GetStrPtr("received instance creation request, pending"))
	if serr != nil {
		logger.Error().Err(serr).Msg("error creating Status Detail DB entry")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Status Detail for Instance, DB error", nil)
	}
	if ssd == nil {
		logger.Error().Msg("Status Detail DB entry not returned from CreateFromParams")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get new Status Detail for Instance", nil)
	}

	// ==================== Step 7: Workflow Trigger  ====================

	// Get the temporal client for the site we are working with.
	stc, err := cih.scp.GetClientByID(instance.SiteID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Prepare the labels for the metadata of the carbide call.
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

	// Prepare the create request workflow object
	createInstanceRequest := &cwssaws.InstanceAllocationRequest{
		InstanceId: &cwssaws.InstanceId{Value: common.GetSiteInstanceID(instance).String()},
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
				TenantKeysetIds:      instanceSshKeyGroupIds,
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
	if instanceTypeID != nil {
		createInstanceRequest.InstanceTypeId = cdb.GetStrPtr(instanceTypeID.String())
	}

	workflowOptions := temporalClient.StartWorkflowOptions{
		ID:                       "instance-create-" + instance.ID.String(),
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	logger.Info().Msg("triggering instance update workflow")

	// Add context deadlines
	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	// Trigger Site workflow to update instance
	// TODO: Once Site Agent offers CreateInstanceV2 re-registered as CreateInstance then update workflow name here
	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "CreateInstanceV2", createInstanceRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to synchronously start Temporal workflow to create Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to start sync workflow to create Instance on Site: %s", err), nil)
	}

	wid := we.GetID()
	logger.Info().Str("Workflow ID", wid).Msg("executed synchronous create Instance workflow")

	// Execute the workflow synchronously
	err = we.Get(ctx, nil)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {

			logger.Error().Err(err).Msg("failed to create Instance, timeout occurred executing workflow on Site.")

			// Create a new context deadlines
			newctx, newcancel := context.WithTimeout(context.Background(), common.WorkflowContextNewAfterTimeout)
			defer newcancel()

			// Initiate termination workflow
			serr := stc.TerminateWorkflow(newctx, wid, "", "timeout occurred executing create Instance workflow")
			if serr != nil {
				logger.Error().Err(serr).Msg("failed to terminate Temporal workflow for creating Instance")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to terminate synchronous Instance creation workflow after timeout, Cloud and Site data may be de-synced: %s", serr), nil)
			}

			logger.Info().Str("Workflow ID", wid).Msg("initiated terminate synchronous create Instance workflow successfully")

			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create Instance, timeout occurred executing workflow on Site: %s", err), nil)
		}

		code, err := common.UnwrapWorkflowError(err)
		logger.Error().Err(err).Msg("failed to synchronously execute Temporal workflow to create Instance")
		return cerr.NewAPIErrorResponse(c, code, fmt.Sprintf("Failed to execute sync workflow to create Instance on Site: %s", err), nil)
	}

	logger.Info().Str("Workflow ID", wid).Msg("completed synchronous create Instance workflow")

	// ==================== Step 8: Commit & Response ====================

	// Commit the DB transaction after the synchronous workflow has completed without error
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing instance transaction to DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Instance, DB transaction error", nil)
	}

	// Set committed so, deferred cleanup functions will do nothing
	txCommitted = true

	// Create response
	apiInstance := model.NewAPIInstance(instance, site, ifcs, ibifcs, desds, nvlifcs, rdbskg, []cdbm.StatusDetail{*ssd})

	logger.Info().Msg("finishing API handler")
	return c.JSON(http.StatusCreated, apiInstance)
}

// ~~~~~ Update Handler ~~~~~ //

// UpdateInstanceHandler is the API Handler for updating an Instance
type UpdateInstanceHandler struct {
	dbSession  *cdb.Session
	tc         temporalClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewUpdateInstanceHandler initializes and returns a new handler for updating Instance
func NewUpdateInstanceHandler(dbSession *cdb.Session, tc temporalClient.Client, scp *sc.ClientPool, cfg *config.Config) UpdateInstanceHandler {
	return UpdateInstanceHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

func (uih UpdateInstanceHandler) handleReboot(c echo.Context, logger *zerolog.Logger, apiRequest *model.APIInstanceUpdateRequest, instance *cdbm.Instance) error {
	ctx := c.Request().Context()

	// Get the temporal client for the site we are working with.
	stc, err := uih.scp.GetClientByID(instance.SiteID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Start a database transaction
	tx, err := cdb.BeginTx(ctx, uih.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Error updating Instance", nil)
	}
	// This variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// Prepare DAOs
	sdDAO := cdbm.NewStatusDetailDAO(uih.dbSession)

	// Check for reboot request
	powerStatus := cdb.GetStrPtr(cdbm.InstancePowerStatusRebooting)
	powerStatusMessage := cdb.GetStrPtr("received Instance reboot request, processing")

	// Check if instance request for rebooting with ipxe
	rebootWithCustomIpxe := false
	if apiRequest.RebootWithCustomIpxe != nil && *apiRequest.RebootWithCustomIpxe {
		rebootWithCustomIpxe = true
	}

	// Check if instance request for updating before reboot
	applyUpdatesOnReboot := false
	if apiRequest.ApplyUpdatesOnReboot != nil && *apiRequest.ApplyUpdatesOnReboot {
		applyUpdatesOnReboot = true
		powerStatusMessage = cdb.GetStrPtr("received Instance reboot request with apply updates, processing")
	}

	// Update Instance
	instanceDAO := cdbm.NewInstanceDAO(uih.dbSession)
	ui, err := instanceDAO.Update(ctx, tx,
		cdbm.InstanceUpdateInput{
			InstanceID:  instance.ID,
			Name:        apiRequest.Name,
			Description: apiRequest.Description,
			PowerStatus: powerStatus,
		},
	)
	if err != nil {
		logger.Error().Err(err).Msg("error updating Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Instance", nil)
	}

	_, serr := sdDAO.CreateFromParams(ctx, tx, instance.ID.String(), *cdb.GetStrPtr(cdbm.InstancePowerStatusRebooting), powerStatusMessage)
	if serr != nil {
		logger.Error().Err(serr).Msg("error creating Status Detail DB entry")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Status Detail for Instance reboot", nil)
	}

	// Get the instance subnets record from the db
	ifcDAO := cdbm.NewInterfaceDAO(uih.dbSession)
	retifc, _, err := ifcDAO.GetAll(ctx, tx, cdbm.InterfaceFilterInput{InstanceIDs: []uuid.UUID{instance.ID}}, cdbp.PageInput{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Instance Subnets Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance Subnets for Instance", nil)
	}

	// Get the ssh key group instance associations record from the db
	skgiaDAO := cdbm.NewSSHKeyGroupInstanceAssociationDAO(uih.dbSession)
	var dbskgs []cdbm.SSHKeyGroup
	skgias, _, err := skgiaDAO.GetAll(ctx, nil, nil, []uuid.UUID{instance.Site.ID}, []uuid.UUID{instance.ID}, []string{cdbm.SSHKeyGroupRelationName}, nil, nil, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving ssh key group instance association Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve SSH Key Group Instance Association for Instance", nil)
	}

	for _, skgia := range skgias {
		dbskgs = append(dbskgs, *skgia.SSHKeyGroup)
	}

	// Get status details
	ssds, _, err := sdDAO.GetAllByEntityID(ctx, tx, ui.ID.String(), nil, nil, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Status Details for Instance from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Status Details for Instance", nil)
	}

	// Prepare the config update request workflow object
	rebootInstanceRequest := &cwssaws.InstancePowerRequest{
		MachineId:            &cwssaws.MachineId{Id: *instance.MachineID},
		Operation:            cwssaws.InstancePowerRequest_POWER_RESET,
		BootWithCustomIpxe:   rebootWithCustomIpxe,
		ApplyUpdatesOnReboot: applyUpdatesOnReboot,
	}

	workflowOptions := temporalClient.StartWorkflowOptions{
		ID:                       "instance-reboot-" + instance.ID.String(),
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	logger.Info().Msg("triggering instance reboot workflow")

	// Add context deadlines
	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	// Trigger Site workflow to update instance
	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "RebootInstanceV2", rebootInstanceRequest)

	if err != nil {
		logger.Error().Err(err).Msg("failed to synchronously start Temporal workflow to reboot Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to start sync workflow to reboot Instance on Site: %s", err), nil)
	}

	wid := we.GetID()
	logger.Info().Str("Workflow ID", wid).Msg("executed synchronous reboot Instance workflow")

	// Execute the workflow synchronously
	err = we.Get(ctx, nil)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {

			logger.Error().Err(err).Msg("failed to reboot Instance, timeout occurred executing workflow on Site.")

			// Create a new context deadlines
			newctx, newcancel := context.WithTimeout(context.Background(), common.WorkflowContextNewAfterTimeout)
			defer newcancel()

			// Initiate termination workflow
			serr := stc.TerminateWorkflow(newctx, wid, "", "timeout occurred executing reboot Instance workflow")
			if serr != nil {
				logger.Error().Err(serr).Msg("failed to terminate Temporal workflow for reboot Instance")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to terminate synchronous Instance reboot workflow after timeout, Cloud and Site data may be de-synced: %s", serr), nil)
			}

			logger.Info().Str("Workflow ID", wid).Msg("initiated terminate synchronous reboot Instance workflow successfully")

			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to reboot Instance, timeout occurred executing workflow on Site: %s", err), nil)
		}
		logger.Error().Err(err).Msg("failed to execute Temporal workflow to reboot Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to execute sync workflow to reboot Instance on Site: %s", err), nil)
	}

	logger.Info().Str("Workflow ID", wid).Msg("completed synchronous reboot Instance workflow")

	// Commit the DB transaction after the synchronous workflow has completed without error
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to reboot Instance, DB transaction error", nil)
	}
	txCommitted = true

	// Create response
	apiInstance := model.NewAPIInstance(ui, instance.Site, retifc, nil, nil, nil, dbskgs, ssds)

	logger.Info().Msg("finishing rebootHandler")

	return c.JSON(http.StatusOK, apiInstance)
}

// Returns either an existing instance OS config or an updated OS config based on the incoming request.
// apiRequest will be mutated for use in UpdateFromParams.
// osConfig will hold the struct/data for use with Temporal/Carbide calls.
// Errors should be returned in the form of cerr.NewAPIErrorResponse
func (uih UpdateInstanceHandler) buildInstanceUpdateRequestOsConfig(c echo.Context, logger *zerolog.Logger, apiRequest *model.APIInstanceUpdateRequest, instance *cdbm.Instance, siteID uuid.UUID) (*cwssaws.OperatingSystem, *uuid.UUID, *cerr.APIError) {

	var os *cdbm.OperatingSystem
	var osID *uuid.UUID

	ctx := c.Request().Context()

	// The OS is being cleared.
	if apiRequest.OperatingSystemID != nil && *apiRequest.OperatingSystemID == "" {

		if err := apiRequest.ValidateAndSetOperatingSystemData(uih.cfg, instance, nil); err != nil {
			logger.Error().Err(err).Msg("failed to validate OperatingSystem")
			return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Failed to validate OperatingSystem data", err)
		}

		return &cwssaws.OperatingSystem{
			RunProvisioningInstructionsOnEveryBoot: instance.AlwaysBootWithCustomIpxe,
			PhoneHomeEnabled:                       *apiRequest.PhoneHomeEnabled, // Set by the earlier call to ValidateAndSetOperatingSystemData
			Variant: &cwssaws.OperatingSystem_Ipxe{
				Ipxe: &cwssaws.IpxeOperatingSystem{
					IpxeScript: *apiRequest.IpxeScript,
					UserData:   apiRequest.UserData,
				},
			},
		}, nil, nil
	}

	// If the base OS is either not changing OR the base is changing to another OS and NOT simply being cleared,
	// then we'll need to pull OS data.

	// Default to the OS of the instance.
	osID = instance.OperatingSystemID

	// Use the OS sent by the caller if one was sent in.
	if apiRequest.OperatingSystemID != nil {
		var id uuid.UUID
		var err error

		if id, err = uuid.Parse(*apiRequest.OperatingSystemID); err != nil {
			logger.Error().Err(err).Msg("failed to parse OperatingSystemID")
			return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Unable to parse `operatingSystemId` specified", validation.Errors{
				"operatingSystemId": errors.New(*apiRequest.OperatingSystemID),
			})
		}

		osID = &id
	}

	if osID != nil {
		var serr error

		// Retrieve the details for the OS
		osDAO := cdbm.NewOperatingSystemDAO(uih.dbSession)
		os, serr = osDAO.GetByID(ctx, nil, *osID, nil)
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

		// Add the OS ID to the log fields now that we know we have a valid one.
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("OperatingSystem ID", os.ID.String())
		})

		// Confirm ownership between tenant and OS.
		if os.TenantID.String() != instance.Tenant.ID.String() {
			logger.Error().Msg("OperatingSystem in request is not owned by tenant")
			return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Operating system specified in request is not owned by Tenant", nil)
		}

		// Confirm match between site and OS (only for Image type).
		if os.Type == cdbm.OperatingSystemTypeImage {
			ossaDAO := cdbm.NewOperatingSystemSiteAssociationDAO(uih.dbSession)
			_, ossaCount, err := ossaDAO.GetAll(
				ctx,
				nil,
				cdbm.OperatingSystemSiteAssociationFilterInput{
					OperatingSystemIDs: []uuid.UUID{*osID},
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
	}

	// reject deactivated OS except if OS stays the same:
	if os != nil && !os.IsActive {
		if apiRequest.OperatingSystemID != nil && instance.OperatingSystemID != nil && *apiRequest.OperatingSystemID != instance.OperatingSystemID.String() {
			return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "Operating System specified in request has been deactivated and cannot be used to update an instance", nil)
		}
	}

	// Validate any additional properties.
	// `os` could still be nil here if no OS ID was sent
	// in the request _and_ the instance didn't have an OS ID
	// to begin with or had previously been cleared (nil'ed)
	// by an earlier request.

	err := apiRequest.ValidateAndSetOperatingSystemData(uih.cfg, instance, os)
	if err != nil {
		logger.Error().Msgf("OperatingSystem options validation failed: %s", err)
		return nil, nil, cerr.NewAPIError(http.StatusBadRequest, "OperatingSystem options validation failed", err)
	}

	// Here, we'll default to whatever the instance already had set,
	// but will give precedence to any property being set by the request.
	// Some or all of these might have been set in ValidateAndSetOperatingSystemData
	// to the desired/expected override value(s).

	alwaysBootWithCustomIpxe := instance.AlwaysBootWithCustomIpxe
	if apiRequest.AlwaysBootWithCustomIpxe != nil {
		alwaysBootWithCustomIpxe = *apiRequest.AlwaysBootWithCustomIpxe
	}

	ipxeScript := instance.IpxeScript
	if apiRequest.IpxeScript != nil {
		ipxeScript = apiRequest.IpxeScript
	}

	userData := instance.UserData
	if apiRequest.UserData != nil {
		userData = apiRequest.UserData
	}

	phoneHomeEnabled := instance.PhoneHomeEnabled
	if apiRequest.PhoneHomeEnabled != nil {
		phoneHomeEnabled = *apiRequest.PhoneHomeEnabled
	}

	if os != nil {
		if os.Type == cdbm.OperatingSystemTypeIPXE {
			return &cwssaws.OperatingSystem{
				RunProvisioningInstructionsOnEveryBoot: alwaysBootWithCustomIpxe,
				PhoneHomeEnabled:                       phoneHomeEnabled,
				Variant: &cwssaws.OperatingSystem_Ipxe{
					Ipxe: &cwssaws.IpxeOperatingSystem{
						IpxeScript: *ipxeScript,
						UserData:   userData,
					},
				},
			}, osID, nil
		} else if os.Type == cdbm.OperatingSystemTypeImage {
			return &cwssaws.OperatingSystem{
				PhoneHomeEnabled: phoneHomeEnabled,
				Variant: &cwssaws.OperatingSystem_OsImageId{
					OsImageId: &cwssaws.UUID{
						Value: os.ID.String(),
					},
				},
				UserData: userData,
			}, osID, nil
		}
	}

	return &cwssaws.OperatingSystem{
		RunProvisioningInstructionsOnEveryBoot: alwaysBootWithCustomIpxe,
		PhoneHomeEnabled:                       phoneHomeEnabled,
		Variant: &cwssaws.OperatingSystem_Ipxe{
			Ipxe: &cwssaws.IpxeOperatingSystem{
				IpxeScript: *ipxeScript,
				UserData:   userData,
			},
		},
	}, osID, nil
}

// Handle godoc
// @Summary Update an existing Instance
// @Description Update an existing Instance for the org
// @Tags Instance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Instance"
// @Param message body model.APIInstanceUpdateRequest true "Instance update request"
// @Success 200 {object} model.APIInstance
// @Router /v2/org/{org}/carbide/instance/{id} [patch]
func (uih UpdateInstanceHandler) Handle(c echo.Context) error {
	// Get context
	ctx := c.Request().Context()

	// Get org
	org := c.Param("orgName")

	// Initialize logger
	logger := log.With().Str("Model", "Instance").Str("Handler", "Update").Str("Org", org).Logger()

	logger.Info().Msg("started API handler")

	// Create a child span and set the attributes for current request
	newctx, handlerSpan := uih.tracerSpan.CreateChildInContext(ctx, "UpdateInstanceHandler", logger)
	if handlerSpan != nil {
		// Set newly created span context as a current context
		ctx = newctx

		defer handlerSpan.End()

		uih.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, uih.tracerSpan, handlerSpan)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Tenant Admins are allowed to update Instances
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.TenantAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Tenant Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Tenant Admin role with org", nil)
	}

	// Get instance ID from URL param
	instanceStrID := c.Param("id")
	instanceID, err := uuid.Parse(instanceStrID)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Instance ID in URL", nil)
	}

	uih.tracerSpan.SetAttribute(handlerSpan, attribute.String("instance_id", instanceStrID), logger)

	// Add the instance ID to the log fields now that we know we have a valid one.
	logger = logger.With().Str("Instance ID", instanceID.String()).Logger()

	// Validate request
	// Bind request data to API model
	apiRequest := model.APIInstanceUpdateRequest{}
	err = c.Bind(&apiRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data, potentially invalid structure", nil)
	}

	// Validate request attributes
	verr := apiRequest.Validate()
	if verr != nil {
		logger.Warn().Err(verr).Msg("error validating Instance update request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate Instance update request data", verr)
	}

	instanceDAO := cdbm.NewInstanceDAO(uih.dbSession)

	// Check that Instance exists
	instance, err := instanceDAO.GetByID(ctx, nil, instanceID, []string{cdbm.SiteRelationName, cdbm.TenantRelationName, cdbm.VpcRelationName})
	if err != nil {
		logger.Warn().Err(err).Msg("error retrieving Instance DB entity")
		return cerr.NewAPIErrorResponse(c, http.StatusNotFound, "Could not retrieve Instance to update", nil)
	}

	if instance.Site == nil {
		logger.Error().Msg("error retrieving Site as included relation for Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site details for Instance", nil)
	}

	if instance.Tenant == nil {
		logger.Error().Msg("error retrieving Tenant as included relation for Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant details for Instance", nil)
	}

	if instance.Vpc == nil {
		logger.Error().Msg("error retrieving VPC as included relation for Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve VPC details for Instance", nil)
	}

	tenant := instance.Tenant
	site := instance.Site
	vpc := instance.Vpc

	// Confirm that the Instance's org matches the org sent in the request
	if tenant.Org != org {
		logger.Error().Err(err).Msg("org specified in request does not match org of Tenant associated with Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Org specified in request does not match Org of Tenant associated with Instance", nil)
	}

	// Add the tenant to the log fields
	logger = logger.With().Str("Tenant ID", tenant.ID.String()).Logger()

	if site.Status != cdbm.SiteStatusRegistered {
		logger.Error().Str("Site ID", site.ID.String()).Str("Site Status", site.Status).Msg("Unable to update Instance, Site is not in Registered state")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site is not in Registered state - cannot update Instance", nil)
	}

	// If the instance is in some stage of deprovisioning, there's nothing to update.
	// We could move this up even higher, but we might not want to reveal status at all until
	// we know the caller has access to this instance.
	if instance.Status == cdbm.InstanceStatusTerminating || instance.Status == cdbm.InstanceStatusTerminated {
		return cerr.NewAPIErrorResponse(c, http.StatusConflict, "Instance is terminating and cannot be updated", nil)
	}

	// check for name uniqueness for the tenant, ie, tenant cannot have another instance with same name at the site
	if apiRequest.Name != nil && *apiRequest.Name != instance.Name {
		ins, tot, serr := instanceDAO.GetAll(ctx, nil,
			cdbm.InstanceFilterInput{
				Names:     []string{*apiRequest.Name},
				TenantIDs: []uuid.UUID{tenant.ID},
				SiteIDs:   []uuid.UUID{site.ID},
			},
			cdbp.PageInput{},
			nil,
		)
		if serr != nil {
			logger.Error().Err(serr).Msg("db error checking for name uniqueness of tenant instance")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Instance due to DB error", nil)
		}
		if tot > 0 {
			return cerr.NewAPIErrorResponse(c, http.StatusConflict, "Another Instance with specified name already exists for Tenant", validation.Errors{
				"id": errors.New(ins[0].ID.String()),
			})
		}
	}

	// Check for reboot request
	instanceRebootRequest := false
	if apiRequest.TriggerReboot != nil && *apiRequest.TriggerReboot {
		instanceRebootRequest = true
	}

	// If this was only a reboot request, handle it and return.
	if instanceRebootRequest {
		return uih.handleReboot(c, &logger, &apiRequest, instance)
	}

	// Otherwise, this is a real Instance config update.
	var instanceStatusConfiguring *string
	if apiRequest.IsInterfaceUpdateRequest() {
		instanceStatusConfiguring = cdb.GetStrPtr(cdbm.InstanceStatusConfiguring)
	}

	// If an NSG was requested, validate it.
	// A blank NSG ID means the user is updating to clear the field.
	var nsgID *string
	if apiRequest.NetworkSecurityGroupID != nil && *apiRequest.NetworkSecurityGroupID != "" {
		nsgDAO := cdbm.NewNetworkSecurityGroupDAO(uih.dbSession)

		nsg, err := nsgDAO.GetByID(ctx, nil, *apiRequest.NetworkSecurityGroupID, nil)
		if err != nil {
			if err == cdb.ErrDoesNotExist {
				logger.Error().Err(err).Msg("could not find NetworkSecurityGroup with ID specified in request data")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not find NetworkSecurityGroup with ID specified in request data", nil)
			}

			logger.Error().Err(err).Msg("error retrieving NetworkSecurityGroup with ID specified in request data")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve NetworkSecurityGroup with ID specified in request data", nil)
		}

		if nsg.SiteID != site.ID {
			logger.Error().Msg("NetworkSecurityGroup in request does not belong to Site")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "NetworkSecurityGroup with ID specified in request data does not belong to Site", nil)
		}

		if nsg.TenantID != tenant.ID {
			logger.Error().Msg("NetworkSecurityGroup in request does not belong to Tenant")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "NetworkSecurityGroup with ID specified in request data does not belong to Tenant", nil)
		}

		nsgID = cdb.GetStrPtr(nsg.ID)
	}

	// Validate Interfaces if present
	sbDAO := cdbm.NewSubnetDAO(uih.dbSession)
	vpDAO := cdbm.NewVpcPrefixDAO(uih.dbSession)
	dbifcs := []cdbm.Interface{}

	// We'll need this later for grabbing network segments
	// to send in the carbide request.
	subnets := map[uuid.UUID]*cdbm.Subnet{}
	vpcPrefixes := map[uuid.UUID]*cdbm.VpcPrefix{}
	isDeviceInfoPresent := false

	for _, ifc := range apiRequest.Interfaces {
		if ifc.SubnetID != nil {
			subnetID, err := uuid.Parse(*ifc.SubnetID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing subnet id in instance subnet request")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Subnet ID specified in request data is not valid", nil)
			}

			if subnets[subnetID] == nil {
				subnet, err := sbDAO.GetByID(ctx, nil, subnetID, nil)
				if err != nil {
					if err == cdb.ErrDoesNotExist {
						return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not find Subnet with ID specified in request data", nil)
					}
					logger.Error().Err(err).Msg("error retrieving Subnet from DB by ID")
					return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Subnet with ID specified in request data", nil)
				}

				if subnet.TenantID != tenant.ID {
					logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request is not owned by Tenant", subnetID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request is not owned by Tenant", subnetID), nil)
				}

				if subnet.ControllerNetworkSegmentID == nil || subnet.Status != cdbm.SubnetStatusReady {
					logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request data is not in Ready state", subnetID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request data is not in Ready state", subnetID), nil)
				}

				if subnet.VpcID != vpc.ID {
					logger.Warn().Msg(fmt.Sprintf("Subnet: %v specified in request does not match with VPC", subnetID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Subnet: %v specified in request does not match with VPC", subnetID), nil)
				}

				if vpc.NetworkVirtualizationType != nil && *vpc.NetworkVirtualizationType != cdbm.VpcEthernetVirtualizer {
					logger.Warn().Msg(fmt.Sprintf("VPC: %v specified in request must have Ethernet network virtualization type in order to create Subnet based interfaces", instance.VpcID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC: %v specified in request must have Ethernet network virtualization type in order to create Subnet based interfaces", instance.VpcID), nil)
				}

				subnets[subnetID] = subnet
			}
			dbifcs = append(dbifcs, cdbm.Interface{SubnetID: &subnetID, IsPhysical: ifc.IsPhysical, Status: cdbm.InterfaceStatusPending})
		}

		if ifc.VpcPrefixID != nil {
			vpcPrefixID, err := uuid.Parse(*ifc.VpcPrefixID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing vpcprefix id in instance vpcprefix request")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "VPC Prefix ID specified in request data is not valid", nil)
			}

			if vpcPrefixes[vpcPrefixID] == nil {
				vpcPrefix, err := vpDAO.GetByID(ctx, nil, vpcPrefixID, nil)
				if err != nil {
					if err == cdb.ErrDoesNotExist {
						return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not find VPC Prefix with ID specified in request data", nil)
					}
					logger.Error().Err(err).Msg("error retrieving vpcprefix from DB by ID")
					return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve VPC Prefix with ID specified in request data", nil)
				}

				if vpcPrefix.TenantID != tenant.ID {
					logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request is not owned by Tenant", vpcPrefixID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request is not owned by Tenant", vpcPrefixID), nil)
				}

				if vpcPrefix.Status != cdbm.VpcPrefixStatusReady {
					logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request data is not in Ready state", vpcPrefixID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request data is not in Ready state", vpcPrefixID), nil)
				}

				if vpcPrefix.VpcID != vpc.ID {
					logger.Warn().Msg(fmt.Sprintf("VPC Prefix: %v specified in request does not match with VPC", vpcPrefixID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC Prefix: %v specified in request does not match with VPC", vpcPrefixID), nil)
				}

				if vpc.NetworkVirtualizationType == nil || *vpc.NetworkVirtualizationType != cdbm.VpcFNN {
					logger.Warn().Msg(fmt.Sprintf("VPC: %v specified in request must have FNN network virtualization type in order to create VPC Prefix based interfaces", instance.VpcID))
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("VPC: %v specified in request must have FNN network virtualization type in order to create VPC Prefix based interfaces", instance.VpcID), nil)
				}

				vpcPrefixes[vpcPrefixID] = vpcPrefix
			}

			if ifc.Device != nil && ifc.DeviceInstance != nil {
				isDeviceInfoPresent = true
			}

			dbifcs = append(dbifcs, cdbm.Interface{
				VpcPrefixID:       &vpcPrefixID,
				Device:            ifc.Device,
				DeviceInstance:    ifc.DeviceInstance,
				VirtualFunctionID: ifc.VirtualFunctionID,
				IsPhysical:        ifc.IsPhysical,
				Status:            cdbm.InterfaceStatusPending})
		}
	}

	mcDAO := cdbm.NewMachineCapabilityDAO(uih.dbSession)

	// Validate DPU Interfaces if Instance Type has Network Capability with DPU device type
	if isDeviceInfoPresent {
		// Get Network Capabilities with DPU device type
		var itDpuCaps []cdbm.MachineCapability
		var itDpuCapCount int
		if instance.InstanceTypeID != nil {
			itDpuCaps, itDpuCapCount, err = mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{*instance.InstanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeNetwork), nil, nil, nil, nil, nil, cdb.GetStrPtr(cdbm.MachineCapabilityDeviceTypeDPU), nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
			if err != nil {
				logger.Error().Err(err).Msg("error retrieving Machine Capabilities from DB for Instance Type")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Machine Capabilities for Instance Type", nil)
			}
		} else {
			itDpuCaps, itDpuCapCount, err = mcDAO.GetAll(ctx, nil, []string{*instance.MachineID}, nil, cdb.GetStrPtr(cdbm.MachineCapabilityTypeNetwork), nil, nil, nil, nil, nil, cdb.GetStrPtr(cdbm.MachineCapabilityDeviceTypeDPU), nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
			if err != nil {
				logger.Error().Err(err).Msg("error retrieving Machine Capabilities from DB for Instance Type")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Machine Capabilities for Instance Type", nil)
			}
		}

		if itDpuCapCount == 0 {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Device and Device Instance cannot be specified if Instance Type doesn't have Network Capabilities with DPU device type", nil)
		}

		// Validate DPU Interfaces if Instance Type DPU capability is present and matches with the request
		err = apiRequest.ValidateMultiEthernetDeviceInterfaces(itDpuCaps, dbifcs)
		if err != nil {
			logger.Error().Msgf("Failed to validate configuration for one or more multi-Ethernet device Interfaces: %s", err)
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate configuration for one or more multi-Ethernet device Interfaces", err)
		}
	}

	ibpDAO := cdbm.NewInfiniBandPartitionDAO(uih.dbSession)
	for _, ibic := range apiRequest.InfiniBandInterfaces {
		// InfiniBand Partition
		ibpID, err := uuid.Parse(ibic.InfiniBandPartitionID)
		if err != nil {
			logger.Warn().Err(err).Msg("error parsing infiniband partition id in instance infiniband interface request")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Partition ID: %v specified in request data is not valid", ibic.InfiniBandPartitionID), nil)
		}

		// Validate Instance infiniband interface information to create DB records later
		ibp, err := ibpDAO.GetByID(ctx, nil, ibpID, nil)
		if err != nil {
			if err == cdb.ErrDoesNotExist {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Could find Partition with ID: %v specified in request data", ibic.InfiniBandPartitionID), nil)
			}
			logger.Error().Err(err).Msg("error retrieving InfiniBand Partition from DB by ID")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Partition with ID specified in request data, DB error", nil)
		}

		if ibp.SiteID != site.ID {
			logger.Warn().Msg(fmt.Sprintf("InfiniBandPartition: %v specified in request does not match with Instance Site", ibpID))
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Partition: %v specified in request does not match with Instance Site", ibpID), nil)
		}

		if ibp.TenantID != tenant.ID {
			logger.Warn().Msg(fmt.Sprintf("InfiniBandPartition: %v specified in request is not owned by Tenant", ibpID))
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Partition: %v specified in request is not owned by Tenant", ibpID), nil)
		}

		if ibp.ControllerIBPartitionID == nil || ibp.Status != cdbm.InfiniBandPartitionStatusReady {
			logger.Warn().Msg(fmt.Sprintf("InfiniBandPartition: %v specified in request data is not in Ready state", ibpID))
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Partition: %v specified in request data is not in Ready state", ibpID), nil)
		}
	}

	if len(apiRequest.NVLinkInterfaces) > 0 {
		nvlIfcDAO := cdbm.NewNVLinkInterfaceDAO(uih.dbSession)
		nvlIfcs, _, err := nvlIfcDAO.GetAll(ctx, nil, cdbm.NVLinkInterfaceFilterInput{InstanceIDs: []uuid.UUID{instance.ID}}, cdbp.PageInput{}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving NVLink Interfaces from DB for Instance")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve NVLink Interfaces for Instance", nil)
		}

		// Discard if VPC has default NVLink Logical Partition specified and NVLink Interfaces are exists
		if len(nvlIfcs) > 0 && vpc.NVLinkLogicalPartitionID != nil {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Cannot update NVLink Interfaces if VPC has default NVLink Logical Partition and NVLink Interfaces already exist for the Instance", nil)
		}
	}

	nvllpDAO := cdbm.NewNVLinkLogicalPartitionDAO(uih.dbSession)
	dbnvlic := []cdbm.NVLinkInterface{}
	for _, nvlifc := range apiRequest.NVLinkInterfaces {
		// NVLink Logical Partition
		nvllpID, err := uuid.Parse(nvlifc.NVLinkLogicalPartitionID)
		if err != nil {
			logger.Warn().Err(err).Msg("error parsing NVLink Logical Partition id in instance NVLink Interface request")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition ID: %v specified in request data is not valid", nvlifc.NVLinkLogicalPartitionID), nil)
		}

		// Validate NVLink Logical Partition
		nvllp, err := nvllpDAO.GetByID(ctx, nil, nvllpID, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving NVLink Logical Partition from DB by ID")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve NVLink Logical Partition with ID specified in request data, DB error", nil)
		}

		if nvllp.SiteID != instance.SiteID {
			logger.Warn().Msg(fmt.Sprintf("NVLink Logical Partition: %v specified in request does not match with Instance Site", nvllpID))
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition: %v specified in request does not match with Instance Site", nvllpID), nil)
		}

		if nvllp.TenantID != instance.TenantID {
			logger.Warn().Msg(fmt.Sprintf("NVLink Logical Partition: %v specified in request data is not owned by Tenant", nvllpID))
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition: %v specified in request data is not owned by Tenant", nvllpID), nil)
		}

		if nvllp.Status != cdbm.NVLinkLogicalPartitionStatusReady {
			logger.Warn().Msg(fmt.Sprintf("NVLink Logical Partition: %v specified in request data is not in Ready state", nvllpID))
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition: %v specified in request data is not in Ready state", nvllpID), nil)
		}

		dbnvlic = append(dbnvlic, cdbm.NVLinkInterface{NVLinkLogicalPartitionID: nvllp.ID, DeviceInstance: nvlifc.DeviceInstance})
	}

	// Get InfiniBand Capabilities
	var itIbCaps []cdbm.MachineCapability

	if instance.InstanceTypeID != nil {
		itIbCaps, _, err = mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{*instance.InstanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeInfiniBand), nil, nil, nil, nil, nil, nil, nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Machine Capabilities from DB for Instance Type")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Machine Capabilities for Instance Type", nil)
		}
	} else {
		itIbCaps, _, err = mcDAO.GetAll(ctx, nil, []string{*instance.MachineID}, nil, cdb.GetStrPtr(cdbm.MachineCapabilityTypeInfiniBand), nil, nil, nil, nil, nil, nil, nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Machine Capabilities from DB for Instance Type")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Machine Capabilities for Instance Type", nil)
		}
	}

	// Validate InfiniBand Interfaces if Instance Type has InfiniBand Capability
	err = apiRequest.ValidateInfiniBandInterfaces(itIbCaps)
	if err != nil {
		logger.Error().Msgf("InfiniBand interfaces validation failed: %s", err)
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate InfiniBand Interfaces specified in request", err)
	}

	desDAO := cdbm.NewDpuExtensionServiceDAO(uih.dbSession)
	desIDMap := map[string]*cdbm.DpuExtensionService{}

	for _, adesdr := range apiRequest.DpuExtensionServiceDeployments {
		desID, err := uuid.Parse(adesdr.DpuExtensionServiceID)
		if err != nil {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid DPU Extension Service ID: %s specified in request", adesdr.DpuExtensionServiceID), nil)
		}

		des, err := desDAO.GetByID(ctx, nil, desID, nil)
		if err != nil {
			if err == cdb.ErrDoesNotExist {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Could not find DPU Extension Service with ID: %s", desID), nil)
			}

			logger.Error().Err(err).Str("DPU Extension Service ID", desID.String()).Msg("error retrieving DPU Extension Service from DB by ID")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve DPU Extension Service specified in request, DB error", nil)
		}

		if des.TenantID != tenant.ID {
			logger.Warn().Str("Tenant ID", tenant.ID.String()).Str("DPU Extension Service ID", desID.String()).Msg("DPU Extension Service does not belong to current Tenant")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("DPU Extension Service: %s does not belong to current Tenant", desID.String()), nil)
		}

		if des.SiteID != site.ID {
			logger.Warn().Str("Site ID", site.ID.String()).Str("DPU Extension Service ID", desID.String()).Msg("DPU Extension Service does not belong to Site")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("DPU Extension Service: %s does not belong to Site where Instance is being created", desID.String()), nil)
		}

		versionFound := false
		for _, version := range des.ActiveVersions {
			if version == adesdr.Version {
				versionFound = true
				break
			}
		}
		if !versionFound {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Version: %s was not found for DPU Extension Service: %s", adesdr.Version, desID.String()), nil)
		}

		desIDMap[desID.String()] = des
	}

	// Validate NVLink interfaces if Instance Type has GPU Capability
	if len(dbnvlic) > 0 {
		var nvlCapCount int
		var nvlCaps []cdbm.MachineCapability
		if instance.InstanceTypeID != nil {
			nvlCaps, nvlCapCount, err = mcDAO.GetAll(ctx, nil, nil, []uuid.UUID{*instance.InstanceTypeID}, cdb.GetStrPtr(cdbm.MachineCapabilityTypeGPU), nil, nil, nil, nil, nil, nil, nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
			if err != nil {
				logger.Error().Err(err).Msg("error retrieving Machine Capabilities from DB for Instance Type")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Machine Capabilities for Instance Type", nil)
			}
		} else {
			nvlCaps, nvlCapCount, err = mcDAO.GetAll(ctx, nil, []string{*instance.MachineID}, nil, cdb.GetStrPtr(cdbm.MachineCapabilityTypeGPU), nil, nil, nil, nil, nil, nil, nil, nil, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
			if err != nil {
				logger.Error().Err(err).Msg("error retrieving Machine Capabilities from DB for Instance's Machine")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Machine Capabilities for Instance's Machine", nil)
			}
		}

		if nvlCapCount == 0 {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "NVLink interfaces cannot be specified if Instance Type doesn't have GPU Capabilities", nil)
		}

		// Validate NVLink interfaces if Instance Type has GPU Capability
		err = apiRequest.ValidateNVLinkInterfaces(nvlCaps)
		if err != nil {
			logger.Error().Msgf("NVLink interfaces validation failed: %s", err)
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate NVLink interfaces specified in request", err)
		}
	}

	// Start a database transaction
	tx, err := cdb.BeginTx(ctx, uih.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Error updating Instance", nil)
	}

	// This variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// Prepare DAOs
	sdDAO := cdbm.NewStatusDetailDAO(uih.dbSession)

	// apiRequest will be mutated for use in UpdateFromParams.
	// osConfig will hold the struct/data for use with Temporal/Carbide calls.
	// Errors will be returned already in the form of cerr.NewAPIError
	osConfig, osID, oserr := uih.buildInstanceUpdateRequestOsConfig(c, &logger, &apiRequest, instance, site.ID)
	if oserr != nil {
		// buildInstanceUpdateRequestOsConfig already handles logging,
		// so this is a bit redundant, but this log brings you to the
		// actual call site.  I think buildInstanceUpdateRequestOsConfig
		// would ideally return only `error` and let the logging and
		// and cerr.NewAPIErrorResponse(...) happen here, but we
		// have at least one StatusInternalServerError case that would
		// be hidden if we merge it all under StatusBadRequest here.
		logger.Error().Err(err).Msg("error building os config for updating Instance")
		return c.JSON(oserr.Code, oserr)
	}

	// Update Instance
	// Once details are fully built, we can just fill out all the columns we have.
	// Postgres either updates a row or it does not.
	// HOT update should not apply here.
	ui, err := instanceDAO.Update(ctx, tx,
		cdbm.InstanceUpdateInput{
			InstanceID:               instanceID,
			Name:                     apiRequest.Name,
			Description:              apiRequest.Description,
			OperatingSystemID:        osID,
			IpxeScript:               apiRequest.IpxeScript,
			AlwaysBootWithCustomIpxe: apiRequest.AlwaysBootWithCustomIpxe,
			NetworkSecurityGroupID:   nsgID,
			PhoneHomeEnabled:         apiRequest.PhoneHomeEnabled,
			Status:                   instanceStatusConfiguring,
			UserData:                 apiRequest.UserData,
			Labels:                   apiRequest.Labels,
		},
	)
	if err != nil {
		logger.Error().Err(err).Msg("error updating Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Instance", nil)
	}

	clearInput := cdbm.InstanceClearInput{InstanceID: instanceID}
	shouldClear := false
	// If this request is attempting to clear the OS for the instance, set it.
	if apiRequest.OperatingSystemID != nil && *apiRequest.OperatingSystemID == "" {
		clearInput.OperatingSystemID = true
		shouldClear = true
	}

	// If this request is attempting to clear the NSG for the instance, set it.
	if apiRequest.NetworkSecurityGroupID != nil {
		if *apiRequest.NetworkSecurityGroupID == "" {
			clearInput.NetworkSecurityGroupID = true
		}

		// We should always clear details for any NSG change so that users don't see stale
		// status.
		clearInput.NetworkSecurityGroupPropagationDetails = true
		shouldClear = true
	}

	// Clear it in the db if something should be cleared.
	if shouldClear {
		ui, err = instanceDAO.Clear(ctx, tx, clearInput)
		if err != nil {
			logger.Error().Err(err).Msg("error clearing requested Instance properties")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to clear requested Instance properties", nil)
		}
	}

	// Save update status in DB
	// Create status detail for instance based on updates requested
	statusMessage := cdb.GetStrPtr("received Instance config update request, processing")

	_, serr := sdDAO.CreateFromParams(ctx, tx, ui.ID.String(), *cdb.GetStrPtr(cdbm.InstanceStatusConfiguring), statusMessage)
	if serr != nil {
		logger.Error().Err(serr).Msg("error creating Status Detail DB entry")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create status detail for Instance update", nil)
	}

	// Get the existing ssh key group instance associations records from the db
	skgiaDAO := cdbm.NewSSHKeyGroupInstanceAssociationDAO(uih.dbSession)
	skgias, _, err := skgiaDAO.GetAll(ctx, nil, nil, []uuid.UUID{site.ID}, []uuid.UUID{instanceID}, []string{cdbm.SSHKeyGroupRelationName}, nil, nil, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving ssh key group instance association Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve SSH Key Group Instance Association for Instance", nil)
	}

	var dbskgs []cdbm.SSHKeyGroup

	// We'll need a list of the IDs as a string slice to
	// send along in the config update request to carbide.
	var instanceSshKeyGroupIds []string

	if apiRequest.SSHKeyGroupIDs == nil {
		// If no change in keygroups, we just need to build the keygroup and keygroup ID
		// lists from the existing instance data so we can send to carbide
		// and return it to the client.
		for _, skgia := range skgias {
			dbskgs = append(dbskgs, *skgia.SSHKeyGroup)
			instanceSshKeyGroupIds = append(instanceSshKeyGroupIds, skgia.SSHKeyGroupID.String())
		}
	} else {
		existingSkgiasBySkg := map[uuid.UUID]*cdbm.SSHKeyGroupInstanceAssociation{}
		for _, skgia := range skgias {
			existingSkgiasBySkg[skgia.SSHKeyGroupID] = &skgia
		}

		skgDAO := cdbm.NewSSHKeyGroupDAO(uih.dbSession)
		skgsaDAO := cdbm.NewSSHKeyGroupSiteAssociationDAO(uih.dbSession)

		incomingSkgMap := map[uuid.UUID]bool{}

		// Determine which SSH Key Group Associations to add.
		for _, skgIDStr := range apiRequest.SSHKeyGroupIDs {
			skgID, err := uuid.Parse(skgIDStr)
			if err != nil {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to update Instance, Invalid SSH Key Group ID: %s", skgIDStr), nil)
			}

			// If the user request has a duplicate, we can skip it.
			if incomingSkgMap[skgID] {
				continue
			}

			incomingSkgMap[skgID] = true

			skgia, found := existingSkgiasBySkg[skgID]
			// If the SKG is already associated with the Instance, we can
			// skip any DB work and just add the keygroup to the lists we'll
			// send to carbide and back to the client.
			if found {
				dbskgs = append(dbskgs, *skgia.SSHKeyGroup)
				instanceSshKeyGroupIds = append(instanceSshKeyGroupIds, skgID.String())
				continue
			}

			// If the SKG is new and not already associated with the Instance
			// we need to associate the SSH Key Group to the Instance.

			// Validate the SSH Key for which this SSH Key Group is being associated.
			sshkeygroup, serr := skgDAO.GetByID(ctx, tx, skgID, nil)
			if serr != nil {
				if serr == common.ErrInvalidID {
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to update Instance, Invalid SSH Key Group ID: %s", skgID), nil)
				}
				if serr == cdb.ErrDoesNotExist {
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to update Instance, Could not find SSH Key Group with ID: %s ", skgID), nil)
				}

				logger.Warn().Err(serr).Str("SSH Key Group ID", skgID.String()).Msg("error retrieving SSH Key Group from DB by ID")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to retrieve SSH Key Group with ID `%s`specified in request, DB error", skgID), nil)
			}

			if sshkeygroup.TenantID != ui.TenantID {
				logger.Warn().Str("Tenant ID", ui.TenantID.String()).Str("SSH Key Group ID", skgID.String()).Msg("SSH Key Group does not belong to current Tenant")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to update Instance, SSH Key Group with ID: %s does not belong to Tenant", skgID), nil)
			}

			// Verify if the SSHKeyGroupSiteAssociation exists
			_, serr = skgsaDAO.GetBySSHKeyGroupIDAndSiteID(ctx, nil, sshkeygroup.ID, site.ID, nil)
			if serr != nil {
				if serr == cdb.ErrDoesNotExist {
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("SSH Key Group with ID: %s is not associated with the Site where Instance is being updated", skgID), nil)
				}
				logger.Warn().Err(serr).Str("SSH Key Group ID", skgID.String()).Msg("error retrieving SSH Key Group Site Association from DB by SSH Key Group ID & Site ID")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to determine if SSH Key Group: %s is associated with the Site where Instance is being updated, DB error", skgID), nil)
			}

			_, err = skgiaDAO.CreateFromParams(ctx, tx, skgID, site.ID, instance.ID, dbUser.ID)
			if err != nil {
				logger.Error().Err(serr).Msg("failed to create the SSH Key Group Instance Association record in DB")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to associate one or more SSH Key Group with Instance, DB error", nil)
			}

			dbskgs = append(dbskgs, *sshkeygroup)
			instanceSshKeyGroupIds = append(instanceSshKeyGroupIds, skgID.String())
		}

		// Determine which SSH Key Group Associations to remove
		for skgID := range existingSkgiasBySkg {
			// Ignore anything see in the users's request.
			// We want to keep those.
			if incomingSkgMap[skgID] {
				continue
			}

			// If not found, we need to disassociate the SSH Key Group from the Instance.
			skgia := existingSkgiasBySkg[skgID]
			err := skgiaDAO.DeleteByID(ctx, tx, skgia.ID)
			if err != nil {
				logger.Error().Err(serr).Str("SSHKeyGroupInstanceAssociation", skgia.ID.String()).Msg("error removing SSH Key Group Instance Association from DB by SSH Key Group Instance Association ID")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to update Instance: %s is associated with the Site where Instance is being updated, DB error", skgia.ID), nil)
			}
		}
	}

	// Create new Interface records in the DB if specified in request

	ifcDAO := cdbm.NewInterfaceDAO(uih.dbSession)

	// OrderAscending is our best-effort to make sure we send
	// Carbide the interfaces in the order it originally received them
	// so the config doesn't get rejected.
	existingIfcs, _, err := ifcDAO.GetAll(ctx, tx, cdbm.InterfaceFilterInput{InstanceIDs: []uuid.UUID{instance.ID}}, cdbp.PageInput{OrderBy: &cdbp.OrderBy{Field: cdbm.InterfaceOrderByCreated, Order: cdbp.OrderAscending}}, []string{cdbm.SubnetRelationName, cdbm.VpcPrefixRelationName})
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve current Ethernet Interfaces details for Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current Ethernet Interfaces for Instance, DB error", nil)
	}

	// Create new Interface records in the DB if specified in request
	var newdbIfcs []cdbm.Interface
	if len(apiRequest.Interfaces) > 0 {
		for _, dbifc := range dbifcs {
			input := cdbm.InterfaceCreateInput{
				InstanceID:        instance.ID,
				SubnetID:          dbifc.SubnetID,
				VpcPrefixID:       dbifc.VpcPrefixID,
				Device:            dbifc.Device,
				DeviceInstance:    dbifc.DeviceInstance,
				VirtualFunctionID: dbifc.VirtualFunctionID,
				IsPhysical:        dbifc.IsPhysical,
				Status:            dbifc.Status,
				CreatedBy:         dbUser.ID,
			}

			dbifc, serr := ifcDAO.Create(ctx, tx, input)
			if serr != nil {
				logger.Error().Err(serr).Msg("error creating Instance Interface DB entry")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Instance Interface entry for Instance, DB error", nil)
			}

			ifc := *dbifc
			// Add the new Interface to the list of new Interfaces
			newdbIfcs = append(newdbIfcs, ifc)
		}

		// Update status of existing Interfaces to Deleting
		for i := range existingIfcs {
			existingIfcs[i].Status = cdbm.InterfaceStatusDeleting
			_, err := ifcDAO.Update(ctx, tx, cdbm.InterfaceUpdateInput{InterfaceID: existingIfcs[i].ID, Status: cdb.GetStrPtr(cdbm.InterfaceStatusDeleting)})
			if err != nil {
				logger.Error().Err(err).Msg("failed to update Interface record in DB")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Interface for Instance, DB error", nil)
			}
		}
	} else {
		newdbIfcs = existingIfcs
	}

	// Create new InfiniBand Interface records in the DB if specified in request
	var newIbIfcs []cdbm.InfiniBandInterface

	ibiDAO := cdbm.NewInfiniBandInterfaceDAO(uih.dbSession)

	// OrderAscending is our best-effort to make sure we send Carbide the interfaces in the order it originally received them. so the config doesn't get rejected
	existingIbIfcs, _, err := ibiDAO.GetAll(ctx, tx, cdbm.InfiniBandInterfaceFilterInput{InstanceIDs: []uuid.UUID{instanceID}}, cdbp.PageInput{OrderBy: &cdbp.OrderBy{Field: cdbm.InfiniBandInterfaceOrderByCreated, Order: cdbp.OrderAscending}}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve InfinibandInterface details for Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Infiniband Interfaces for Instance, DB error", nil)
	}

	if apiRequest.InfiniBandInterfaces != nil {
		for _, apiibifc := range apiRequest.InfiniBandInterfaces {
			// NOTE: This is redundant due to earlier validation, but we handle it anyway
			ibpID, err := uuid.Parse(apiibifc.InfiniBandPartitionID)
			if err != nil {
				logger.Error().Err(err).Msg("failed to parse InfinibandPartitionID")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to parse InfiniBand Partition ID specified in request: %s", apiibifc.InfiniBandPartitionID), nil)
			}

			dbibifc, err := ibiDAO.Create(ctx, tx, cdbm.InfiniBandInterfaceCreateInput{
				InstanceID:            instanceID,
				SiteID:                site.ID,
				InfiniBandPartitionID: ibpID,
				Device:                apiibifc.Device,
				Vendor:                apiibifc.Vendor,
				DeviceInstance:        apiibifc.DeviceInstance,
				IsPhysical:            apiibifc.IsPhysical,
				VirtualFunctionID:     apiibifc.VirtualFunctionID,
				Status:                cdbm.InfiniBandInterfaceStatusPending,
				CreatedBy:             dbUser.ID,
			})

			if err != nil {
				logger.Error().Err(err).Msg("failed to create Infiniband Interface record in DB")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Infiniband Interface for Instance, DB error", nil)
			}

			newIbIfcs = append(newIbIfcs, *dbibifc)
		}

		// Update status of existing InfiniBand Interfaces to Deleting
		for i := range existingIbIfcs {
			existingIbIfcs[i].Status = cdbm.InfiniBandInterfaceStatusDeleting
			_, err = ibiDAO.Update(ctx, tx, cdbm.InfiniBandInterfaceUpdateInput{
				InfiniBandInterfaceID: existingIbIfcs[i].ID,
				Status:                cdb.GetStrPtr(cdbm.InfiniBandInterfaceStatusDeleting),
			})
			if err != nil {
				logger.Error().Err(err).Msg("failed to update Infiniband Interface record in DB")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Infiniband Interface for Instance, DB error", nil)
			}
		}
	} else {
		newIbIfcs = existingIbIfcs
	}

	// Fetch existing DPU Extension Service Deployments for the Instance
	desdDAO := cdbm.NewDpuExtensionServiceDeploymentDAO(uih.dbSession)
	existingDesds, _, err := desdDAO.GetAll(ctx, tx, cdbm.DpuExtensionServiceDeploymentFilterInput{
		InstanceIDs: []uuid.UUID{instance.ID},
	}, cdbp.PageInput{
		OrderBy: &cdbp.OrderBy{Field: cdbm.DpuExtensionServiceDeploymentOrderByDefault, Order: cdbp.OrderAscending},
		Limit:   cdb.GetIntPtr(cdbp.TotalLimit),
	}, []string{cdbm.DpuExtensionServiceRelationName})
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve DpuExtensionServiceDeployment details for Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve existing DPU Extension Service Deployments for Instance, DB error", nil)
	}

	// Check if any DPU Extension Service Deployments are being requested to be created or removed
	existingDesdMap := map[string]*cdbm.DpuExtensionServiceDeployment{}
	for _, desd := range existingDesds {
		existingDesdMap[fmt.Sprintf("%s:%s", desd.DpuExtensionServiceID.String(), desd.Version)] = &desd
	}

	updateDesds := []cdbm.DpuExtensionServiceDeployment{}
	updatedDesdMap := map[string]*cdbm.DpuExtensionServiceDeployment{}

	if len(apiRequest.DpuExtensionServiceDeployments) > 0 {
		for _, adesdr := range apiRequest.DpuExtensionServiceDeployments {
			desdID, err := uuid.Parse(adesdr.DpuExtensionServiceID)
			if err != nil {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid DPU Extension Service ID: %s specified in request", adesdr.DpuExtensionServiceID), nil)
			}

			desvID := fmt.Sprintf("%s:%s", desdID.String(), adesdr.Version)

			existingDesd, exists := existingDesdMap[desvID]
			if exists {
				updateDesds = append(updateDesds, *existingDesd)
				updatedDesdMap[desvID] = existingDesd
			} else {
				newDesd, serr := desdDAO.Create(ctx, tx, cdbm.DpuExtensionServiceDeploymentCreateInput{
					SiteID:                site.ID,
					TenantID:              tenant.ID,
					InstanceID:            instance.ID,
					DpuExtensionServiceID: desdID,
					Version:               adesdr.Version,
					Status:                cdbm.DpuExtensionServiceDeploymentStatusPending,
					CreatedBy:             dbUser.ID,
				})
				if serr != nil {
					logger.Error().Err(serr).Msg("error creating Instance DpuExtensionServiceDeployment record in DB")
					return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create DPU Extension Service Deployment for Instance, DB error", nil)
				}
				des, _ := desIDMap[desdID.String()]
				newDesd.DpuExtensionService = des
				updateDesds = append(updateDesds, *newDesd)
				updatedDesdMap[desvID] = newDesd
			}
		}

		for _, existingDesd := range existingDesds {
			desvID := fmt.Sprintf("%s:%s", existingDesd.DpuExtensionServiceID.String(), existingDesd.Version)
			_, exists := updatedDesdMap[desvID]
			if !exists && existingDesd.Status != cdbm.DpuExtensionServiceDeploymentStatusTerminating {
				// TH deployment is not present in request sent by user, update status to Terminating if not already in that state
				_, err = desdDAO.Update(ctx, tx, cdbm.DpuExtensionServiceDeploymentUpdateInput{
					DpuExtensionServiceDeploymentID: existingDesd.ID,
					Status:                          cdb.GetStrPtr(cdbm.DpuExtensionServiceDeploymentStatusTerminating)})
				if err != nil {
					logger.Error().Err(err).Msg("failed to update DpuExtensionServiceDeployment record in DB")
					return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update DPU Extension Service Deployment for Instance, DB error", nil)
				}
			}
		}
	}

	// Create new NVLink Interface records in the DB if specified in request
	var newNvlIfcs []cdbm.NVLinkInterface
	nvlIfcDAO := cdbm.NewNVLinkInterfaceDAO(uih.dbSession)

	// OrderAscending is our best-effort to make sure we send Carbide the interfaces in the order it originally received them. so the config doesn't get rejected
	existingNvlIfcs, _, err := nvlIfcDAO.GetAll(ctx, tx, cdbm.NVLinkInterfaceFilterInput{InstanceIDs: []uuid.UUID{instanceID}}, cdbp.PageInput{OrderBy: &cdbp.OrderBy{Field: cdbm.NVLinkInterfaceOrderByCreated, Order: cdbp.OrderAscending}}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve NVLink Interfaces details for Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve NVLink interfaces for Instance, DB error", nil)
	}

	if apiRequest.NVLinkInterfaces != nil {
		for _, apiNvlIfc := range apiRequest.NVLinkInterfaces {
			// NVLink Logical Partition
			nvllPartitionID, err := uuid.Parse(apiNvlIfc.NVLinkLogicalPartitionID)
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing NVLink Logical Partition id in instance NVLink Interface request")
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("NVLink Logical Partition ID: %v specified in request data is not valid", apiNvlIfc.NVLinkLogicalPartitionID), nil)
			}

			// Validate NVLink Logical Partition
			nvllPartition, err := nvllpDAO.GetByID(ctx, nil, nvllPartitionID, nil)
			if err != nil {
				logger.Error().Err(err).Msg("error retrieving NVLink Logical Partition from DB by ID")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve NVLink Logical Partition with ID specified in request data, DB error", nil)
			}

			newNvlIfc, err := nvlIfcDAO.Create(ctx, tx, cdbm.NVLinkInterfaceCreateInput{
				InstanceID:               instanceID,
				SiteID:                   site.ID,
				NVLinkLogicalPartitionID: nvllPartition.ID,
				DeviceInstance:           apiNvlIfc.DeviceInstance,
				Status:                   cdbm.NVLinkInterfaceStatusPending,
				CreatedBy:                dbUser.ID,
			})

			if err != nil {
				logger.Error().Err(err).Msg("failed to create NVLink Interface record in DB")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create NVLink Interface for Instance, DB error", nil)
			}
			newNvlIfcs = append(newNvlIfcs, *newNvlIfc)
		}

		// Update status of existing NVLink interfaces to Deleting
		for i := range existingNvlIfcs {
			existingNvlIfcs[i].Status = cdbm.NVLinkInterfaceStatusDeleting
			_, err := nvlIfcDAO.Update(ctx, tx, cdbm.NVLinkInterfaceUpdateInput{
				NVLinkInterfaceID: existingNvlIfcs[i].ID,
				Status:            cdb.GetStrPtr(cdbm.NVLinkInterfaceStatusDeleting),
			})
			if err != nil {
				logger.Error().Err(err).Msg("failed to update NVLink Interface record in DB")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update NVLink Interface for Instance, DB error", nil)
			}
		}

	} else {
		newNvlIfcs = existingNvlIfcs
	}

	// Get Status Details
	ssds, _, err := sdDAO.GetAllByEntityID(ctx, tx, ui.ID.String(), nil, nil, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Status Details for Instance from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Status Details for Instance", nil)
	}

	// Get the temporal client for the site we are working with.
	stc, err := uih.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Prepare the labels for the metadata of the carbide call.
	labels := []*cwssaws.Label{}
	for k, v := range ui.Labels {
		labels = append(labels, &cwssaws.Label{
			Key:   k,
			Value: &v,
		})
	}

	description := ""
	if ui.Description != nil {
		description = *ui.Description
	}

	interfaceConfigs := make([]*cwssaws.InstanceInterfaceConfig, len(newdbIfcs))
	for i, ifc := range newdbIfcs {
		if ifc.Status == cdbm.InterfaceStatusDeleting {
			// NOTE: Don't send any Interfaces that are being deleted
			continue
		}

		interfaceConfig := &cwssaws.InstanceInterfaceConfig{
			FunctionType: cwssaws.InterfaceFunctionType_VIRTUAL_FUNCTION,
		}

		if ifc.SubnetID != nil {
			interfaceConfig.NetworkSegmentId = &cwssaws.NetworkSegmentId{Value: ifc.SubnetID.String()}
			interfaceConfig.NetworkDetails = &cwssaws.InstanceInterfaceConfig_SegmentId{
				SegmentId: &cwssaws.NetworkSegmentId{Value: ifc.SubnetID.String()},
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

		// Assign Device and DeviceInstance in case of Multi DPU Interface
		if ifc.Device != nil && ifc.DeviceInstance != nil {
			interfaceConfig.Device = ifc.Device
			interfaceConfig.DeviceInstance = uint32(*ifc.DeviceInstance)
		}

		if !ifc.IsPhysical {
			if ifc.VirtualFunctionID != nil {
				vfID := uint32(*ifc.VirtualFunctionID)
				interfaceConfig.VirtualFunctionId = &vfID
			}
		}

		interfaceConfigs[i] = interfaceConfig
	}

	// Populate InfiniBand Interface details for Site Controller request
	ibInterfaceConfigs := []*cwssaws.InstanceIBInterfaceConfig{}

	for _, newIbIfc := range newIbIfcs {
		if newIbIfc.Status == cdbm.InfiniBandInterfaceStatusDeleting {
			// NOTE: Don't send any InfiniBand Interfaces that are being deleted
			continue
		}

		ibInterfaceConfig := &cwssaws.InstanceIBInterfaceConfig{
			Device:         newIbIfc.Device,
			Vendor:         newIbIfc.Vendor,
			DeviceInstance: uint32(newIbIfc.DeviceInstance),
			FunctionType:   cwssaws.InterfaceFunctionType_PHYSICAL_FUNCTION,
			IbPartitionId:  &cwssaws.IBPartitionId{Value: newIbIfc.InfiniBandPartitionID.String()},
		}

		// NOTE: Not supported yet, but ensures future compatibility
		if !newIbIfc.IsPhysical {
			ibInterfaceConfig.FunctionType = cwssaws.InterfaceFunctionType_VIRTUAL_FUNCTION

			if newIbIfc.VirtualFunctionID != nil {
				vfID := uint32(*newIbIfc.VirtualFunctionID)
				ibInterfaceConfig.VirtualFunctionId = &vfID
			}
		}

		ibInterfaceConfigs = append(ibInterfaceConfigs, ibInterfaceConfig)
	}

	// Populate DPU Extension Service Deployment details for Site Controller request
	desdConfigs := []*cwssaws.InstanceDpuExtensionServiceConfig{}
	for _, desd := range updateDesds {
		desdConfigs = append(desdConfigs, &cwssaws.InstanceDpuExtensionServiceConfig{
			ServiceId: desd.DpuExtensionServiceID.String(),
			Version:   desd.Version,
		})
	}

	// Populate NVLink Interface details for Site Controller request
	nvlInterfaceConfigs := []*cwssaws.InstanceNVLinkGpuConfig{}
	for _, newNvlIfc := range newNvlIfcs {
		if newNvlIfc.Status == cdbm.NVLinkInterfaceStatusDeleting {
			// NOTE: Don't send any NVLink interfaces that are being deleted
			continue
		}
		nvlInterfaceConfig := &cwssaws.InstanceNVLinkGpuConfig{
			DeviceInstance:     uint32(newNvlIfc.DeviceInstance),
			LogicalPartitionId: &cwssaws.NVLinkLogicalPartitionId{Value: newNvlIfc.NVLinkLogicalPartitionID.String()},
		}
		nvlInterfaceConfigs = append(nvlInterfaceConfigs, nvlInterfaceConfig)
	}

	// Prepare the config update request workflow object
	updateInstanceRequest := &cwssaws.InstanceConfigUpdateRequest{
		InstanceId: &cwssaws.InstanceId{Value: common.GetSiteInstanceID(instance).String()},
		Metadata: &cwssaws.Metadata{
			Name:        ui.Name,
			Description: description,
			Labels:      labels,
		},
		Config: &cwssaws.InstanceConfig{
			NetworkSecurityGroupId: ui.NetworkSecurityGroupID,
			Tenant: &cwssaws.TenantConfig{
				TenantOrganizationId: tenant.Org,
				TenantKeysetIds:      instanceSshKeyGroupIds,
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
	}

	workflowOptions := temporalClient.StartWorkflowOptions{
		ID:                       "instance-update-" + instance.ID.String(),
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	logger.Info().Msg("triggering instance update workflow")

	// Add context deadlines
	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	// Trigger Site workflow to update instance
	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "UpdateInstance", updateInstanceRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to synchronously start Temporal workflow to update Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed start sync workflow to update Instance on Site: %s", err), nil)
	}

	wid := we.GetID()
	logger.Info().Str("Workflow ID", wid).Msg("executed synchronous update Instance workflow")

	// Execute the workflow synchronously
	err = we.Get(ctx, nil)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {

			logger.Error().Err(err).Msg("failed to update Instance, timeout occurred executing workflow on Site.")

			// Create a new context deadlines
			newctx, newcancel := context.WithTimeout(context.Background(), common.WorkflowContextNewAfterTimeout)
			defer newcancel()

			// Initiate termination workflow
			serr := stc.TerminateWorkflow(newctx, wid, "", "timeout occurred executing update Instance workflow")
			if serr != nil {
				logger.Error().Err(serr).Msg("failed to terminate Temporal workflow for update Instance")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to terminate synchronous Instance update workflow after timeout, Cloud and Site data may be de-synced: %s", serr), nil)
			}

			logger.Info().Str("Workflow ID", wid).Msg("initiated terminate synchronous update Instance workflow successfully")

			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update Instance, timeout occurred executing workflow on Site: %s", err), nil)
		}

		code, err := common.UnwrapWorkflowError(err)
		logger.Error().Err(err).Msg("failed to synchronously execute Temporal workflow to update Instance")
		return cerr.NewAPIErrorResponse(c, code, fmt.Sprintf("Failed to execute sync workflow to update Instance on Site: %s", err), nil)
	}

	logger.Info().Str("Workflow ID", wid).Msg("completed synchronous update Instance workflow")

	// Commit the DB transaction after the synchronous workflow has completed without error
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Instance, DB transaction error", nil)
	}
	txCommitted = true

	// If existing Interfaces were updated, add them to the response
	if existingIfcs != nil && len(existingIfcs) > 0 {
		// Add the existing Interfaces to the response
		newdbIfcs = append(newdbIfcs, existingIfcs...)
	}

	// If existing InfiniBand Interfaces were updated, add them to the response
	if existingIbIfcs != nil && len(existingIbIfcs) > 0 {
		// Add the existing InfiniBand Interfaces to the response
		newIbIfcs = append(newIbIfcs, existingIbIfcs...)
	}

	// If existing NVLink Interfaces were updated, add them to the response
	if existingNvlIfcs != nil && len(existingNvlIfcs) > 0 {
		// Add the existing NVLink Interfaces to the response
		newNvlIfcs = append(newNvlIfcs, existingNvlIfcs...)
	}

	// Create response
	apiInstance := model.NewAPIInstance(ui, site, newdbIfcs, newIbIfcs, updateDesds, newNvlIfcs, dbskgs, ssds)

	logger.Info().Msg("finishing API handler")
	return c.JSON(http.StatusOK, apiInstance)
}

// ~~~~~ Get Handler ~~~~~ //

// GetInstanceHandler is the API Handler for getting an Instance
type GetInstanceHandler struct {
	dbSession  *cdb.Session
	tc         temporalClient.Client
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewGetInstanceHandler initializes and returns a new handler for getting Instance
func NewGetInstanceHandler(dbSession *cdb.Session, tc temporalClient.Client, cfg *config.Config) GetInstanceHandler {
	return GetInstanceHandler{
		dbSession:  dbSession,
		tc:         tc,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Get an Instance
// @Description Get an Instance for the org
// @Tags Instance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Instance"
// @Param includeRelation query string false "Related entities to include in response e.g. 'InfrastructureProvider', 'Tenant', 'Site'"
// @Success 200 {object} model.APIInstance
// @Router /v2/org/{org}/carbide/instance/{id} [get]
func (gih GetInstanceHandler) Handle(c echo.Context) error {
	// Get context
	ctx := c.Request().Context()

	// Get org
	org := c.Param("orgName")

	// Initialize logger
	logger := log.With().Str("Model", "Instance").Str("Handler", "Get").Str("Org", org).Logger()

	logger.Info().Msg("started API handler")

	// Create a child span and set the attributes for current request
	newctx, handlerSpan := gih.tracerSpan.CreateChildInContext(ctx, "GetInstanceHandler", logger)
	if handlerSpan != nil {
		// Set newly created span context as a current context
		ctx = newctx

		defer handlerSpan.End()

		gih.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, gih.tracerSpan, handlerSpan)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Tenant Admins are allowed to retrieve Instances
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.TenantAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Tenant Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Tenant Admin role with org", nil)
	}

	// Get and validate includeRelation params
	qParams := c.QueryParams()
	qIncludeRelations, errMsg := common.GetAndValidateQueryRelations(qParams, cdbm.InstanceRelatedEntities)
	if errMsg != "" {
		logger.Warn().Msg(errMsg)
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, errMsg, nil)
	}

	// Get Instance ID from URL param
	instanceStrID := c.Param("id")
	instanceID, err := uuid.Parse(instanceStrID)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Instance ID in URL", nil)
	}

	gih.tracerSpan.SetAttribute(handlerSpan, attribute.String("instance_id", instanceStrID), logger)

	// Get Instance
	instanceDAO := cdbm.NewInstanceDAO(gih.dbSession)

	instance, err := instanceDAO.GetByID(ctx, nil, instanceID, qIncludeRelations)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIErrorResponse(c, http.StatusNotFound, "Could not find Instance with specified ID", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Instance from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance", nil)
	}

	// Get Tenant for this org
	tnDAO := cdbm.NewTenantDAO(gih.dbSession)

	tenants, err := tnDAO.GetAllByOrg(ctx, nil, org, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Tenant for this org")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant", nil)
	}

	if len(tenants) == 0 {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Org does not have a Tenant associated", nil)
	}
	tenant := tenants[0]

	// Check if Instance belongs to Tenant
	if instance.TenantID != tenant.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Instance does not belong to current Tenant", nil)
	}

	// Get Site for this Instance
	siteDAO := cdbm.NewSiteDAO(gih.dbSession)
	site, err := siteDAO.GetByID(ctx, nil, instance.SiteID, nil, false)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Site DB entity")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site for Instance", nil)
	}

	// Get the instance subnets record from the db
	ifcDAO := cdbm.NewInterfaceDAO(gih.dbSession)
	ifcs, _, err := ifcDAO.GetAll(ctx, nil, cdbm.InterfaceFilterInput{InstanceIDs: []uuid.UUID{instance.ID}}, cdbp.PageInput{}, []string{cdbm.SubnetRelationName, cdbm.VpcPrefixRelationName})
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving instance Subnet Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance Subnets for Instance", nil)
	}

	// Get the instance infiniband interface record from the db
	ibifcDAO := cdbm.NewInfiniBandInterfaceDAO(gih.dbSession)
	ibIfcs, _, err := ibifcDAO.GetAll(
		ctx,
		nil,
		cdbm.InfiniBandInterfaceFilterInput{
			InstanceIDs: []uuid.UUID{instanceID},
		},
		cdbp.PageInput{},
		[]string{cdbm.InfiniBandPartitionRelationName},
	)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving instance InfiniBand Interfaces Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance InfiniBand Interfaces for Instance", nil)
	}

	// Get DPU Extension Service Deployments for the instance
	desdDAO := cdbm.NewDpuExtensionServiceDeploymentDAO(gih.dbSession)
	desds, _, err := desdDAO.GetAll(
		ctx,
		nil,
		cdbm.DpuExtensionServiceDeploymentFilterInput{
			InstanceIDs: []uuid.UUID{instanceID},
		},
		cdbp.PageInput{
			OrderBy: &cdbp.OrderBy{Field: cdbm.DpuExtensionServiceDeploymentOrderByDefault, Order: cdbp.OrderAscending},
			Limit:   cdb.GetIntPtr(cdbp.TotalLimit),
		},
		[]string{cdbm.DpuExtensionServiceRelationName},
	)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving DPU Extension Service Deployments for instance from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve DPU Extension Service Deployments for instance", nil)
	}

	// Get the instance NVLink Interface record from the db
	nvlDAO := cdbm.NewNVLinkInterfaceDAO(gih.dbSession)
	nvlIfcs, _, err := nvlDAO.GetAll(ctx, nil, cdbm.NVLinkInterfaceFilterInput{InstanceIDs: []uuid.UUID{instanceID}}, cdbp.PageInput{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving instance NVLink interfaces Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance NVLink interfaces for Instance", nil)
	}

	// Get the ssh key group instance associations record from the db
	skgiaDAO := cdbm.NewSSHKeyGroupInstanceAssociationDAO(gih.dbSession)
	var dbskgs []cdbm.SSHKeyGroup
	skgias, _, err := skgiaDAO.GetAll(ctx, nil, nil, []uuid.UUID{site.ID}, []uuid.UUID{instanceID}, []string{cdbm.SSHKeyGroupRelationName}, nil, nil, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving ssh key group instance association Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve SSH Key Group Instance Association for Instance", nil)
	}

	for _, skgia := range skgias {
		dbskgs = append(dbskgs, *skgia.SSHKeyGroup)
	}

	// Get status details
	sdDAO := cdbm.NewStatusDetailDAO(gih.dbSession)
	ssds, err := sdDAO.GetRecentByEntityIDs(ctx, nil, []string{instanceID.String()}, common.RECENT_STATUS_DETAIL_COUNT)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Status Details for instance from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Status Details for instance", nil)
	}

	// Create response
	ins := model.NewAPIInstance(instance, site, ifcs, ibIfcs, desds, nvlIfcs, dbskgs, ssds)

	// If the instance has no NSG ID, then we need to check if its parent VPC does.
	// We'll need to pull that separately because the user might not have asked for
	// the VPC relation, so we can't assume that it's there.

	if instance.NetworkSecurityGroupID == nil {
		vpcDAO := cdbm.NewVpcDAO(gih.dbSession)

		vpc, err := vpcDAO.GetByID(ctx, nil, instance.VpcID, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving VPC DB entity")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve VPC for Instance", nil)
		}

		// We've only inherited if our NSG ID is null _and_ the parent
		// NSG ID is not null.
		ins.NetworkSecurityGroupInherited = vpc.NetworkSecurityGroupID != nil

		// If we inherited our NSG, then see if we're propagated.
		if ins.NetworkSecurityGroupInherited {

			// We can default to non/configuring and then switch
			// if we're actually propagated.
			ins.NetworkSecurityGroupPropagationDetails = &model.APINetworkSecurityGroupPropagationDetails{
				ObjectID:       instanceID.String(),
				DetailedStatus: model.APINetworkSecurityGroupPropagationDetailedStatusNone,
				Status:         model.APINetworkSecurityGroupPropagationStatusSynchronizing,
			}

			if vpc.NetworkSecurityGroupPropagationDetails != nil {
				// If the instance wasn't found in the list of unpropagated instances, then we're propagated.
				if !slices.Contains(vpc.NetworkSecurityGroupPropagationDetails.GetUnpropagatedInstanceIds(), instanceID.String()) {
					ins.NetworkSecurityGroupPropagationDetails.DetailedStatus = model.APINetworkSecurityGroupPropagationDetailedStatusFull
					ins.NetworkSecurityGroupPropagationDetails.Status = model.APINetworkSecurityGroupPropagationStatusSynchronized
				}
			}
		}
	}

	logger.Info().Msg("finishing API handler")

	return c.JSON(http.StatusOK, ins)
}

// ~~~~~ GetAll Handler ~~~~~ //

// GetAllInstanceHandler is the API Handler for retrieving all Instances
type GetAllInstanceHandler struct {
	dbSession  *cdb.Session
	tc         temporalClient.Client
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewGetAllInstanceHandler initializes and returns a new handler for retreiving all Instances
func NewGetAllInstanceHandler(dbSession *cdb.Session, tc temporalClient.Client, cfg *config.Config) GetAllInstanceHandler {
	return GetAllInstanceHandler{
		dbSession:  dbSession,
		tc:         tc,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Get all Instances
// @Description Get all Instances for the org
// @Tags Instance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param infrastructureProviderId query string true "Infrastructure Provider ID"
// @Param siteId query string true "ID of Site"
// @Param vpcId query string true "ID of Vpc"
// @Param instanceTypeId query string false "ID of Instance Type"
// @Param operatingSystemId query string false "ID of Operating System"
// @Param name query string false "Filter by Instance name"
// @Param status query string false "Filter by status" e.g. 'Pending', 'Error'"
// @Param ipAddress query string false "Filter by IP address. Can be specified multiple times to filter on more than one IP address."
// @Param query query string false "Query input for full text search"
// @Param includeRelation query string false "Related entities to include in response e.g. 'InfrastructureProvider', 'Tenant', 'Site'"
// @Param pageNumber query integer false "Page number of results returned"
// @Param pageSize query integer false "Number of results per page"
// @Param orderBy query string false "Order by field"
// @Success 200 {array} []model.APIInstance
// @Router /v2/org/{org}/carbide/instance [get]
func (gaih GetAllInstanceHandler) Handle(c echo.Context) error {
	// Get context
	ctx := c.Request().Context()

	// Get org
	org := c.Param("orgName")

	// Initialize logger
	logger := log.With().Str("Model", "Instance").Str("Handler", "GetAll").Str("Org", org).Logger()

	logger.Info().Msg("started API handler")

	// Create a child span and set the attributes for current request
	newctx, handlerSpan := gaih.tracerSpan.CreateChildInContext(ctx, "GetAllInstanceHandler", logger)
	if handlerSpan != nil {
		// Set newly created span context as a current context
		ctx = newctx

		defer handlerSpan.End()

		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, gaih.tracerSpan, handlerSpan)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Tenant Admins are allowed to retrieve Instances
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.TenantAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Tenant Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Tenant Admin role with org", nil)
	}

	// Validate pagination request
	pageRequest := pagination.PageRequest{}
	err = c.Bind(&pageRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding pagination request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request pagination data", nil)
	}

	// Validate pagination request attributes
	err = pageRequest.Validate(cdbm.InstanceOrderByFields)
	if err != nil {
		logger.Warn().Err(err).Msg("error validating pagination request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest,
			"Failed to validate pagination request data", err)
	}

	// Get and validate includeRelation params
	qParams := c.QueryParams()
	qIncludeRelations, errMsg := common.GetAndValidateQueryRelations(qParams, cdbm.InstanceRelatedEntities)
	if errMsg != "" {
		logger.Warn().Msg(errMsg)
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, errMsg, nil)
	}

	filter := cdbm.InstanceFilterInput{}
	// Get infrastructure ID from query param if specified
	infrastructureProviderIDStr := c.QueryParam("infrastructureProviderId")
	if infrastructureProviderIDStr != "" {
		parsedID, serr := uuid.Parse(infrastructureProviderIDStr)
		if serr != nil {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Infrastructure Provider ID in query", nil)
		}

		// Check for Provider existence
		ifpDAO := cdbm.NewInfrastructureProviderDAO(gaih.dbSession)
		_, verr := ifpDAO.GetByID(ctx, nil, parsedID, nil)
		if verr != nil {
			logger.Warn().Err(verr).Msg("error retrieving Infrastructure Provider from DB by ID")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not retrieve Infrastructure Provider with ID specified in query", nil)
		}
		filter.InfrastructureProviderIDs = []uuid.UUID{parsedID}
	}

	// Get Tenant for this org
	tnDAO := cdbm.NewTenantDAO(gaih.dbSession)

	tenants, err := tnDAO.GetAllByOrg(ctx, nil, org, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Tenant for this org")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant for org", nil)
	}

	if len(tenants) == 0 {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Org does not have a Tenant associated", nil)
	}
	tenant := tenants[0]
	filter.TenantIDs = append(filter.TenantIDs, tenant.ID)

	// Get site IDs from query param and validate
	stDAO := cdbm.NewSiteDAO(gaih.dbSession)

	var siteIDs []uuid.UUID
	sitesByID := map[uuid.UUID]*cdbm.Site{}
	siteIDStrs := qParams["siteId"]

	for _, siteIDStr := range siteIDStrs {
		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.StringSlice("siteId", siteIDStrs), logger)
		parsedID, err := uuid.Parse(siteIDStr)
		if err != nil {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid site ID %v in query", siteIDStr), nil)
		}

		siteIDs = append(siteIDs, parsedID)
	}

	if siteIDs != nil {
		siteIDs = goset.NewSet(siteIDs...).ToSlice()

		sites, _, err := stDAO.GetAll(ctx, nil, cdbm.SiteFilterInput{SiteIDs: siteIDs}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Sites from DB")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Sites s", nil)
		}
		for _, site := range sites {
			sitesByID[site.ID] = &site
		}

		if len(sites) != len(siteIDs) {
			return cerr.NewAPIErrorResponse(c, http.StatusNotFound, "Could not find one or more Sites specified in query", nil)
		}
	} else {
		tsDAO := cdbm.NewTenantSiteDAO(gaih.dbSession)
		tss, _, err := tsDAO.GetAll(ctx, nil, cdbm.TenantSiteFilterInput{TenantIDs: []uuid.UUID{tenant.ID}}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, []string{cdbm.SiteRelationName})
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Sites for Tenant from DB")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Sites for Tenant, DB error", nil)
		}
		for _, ts := range tss {
			// Check if Site relation was loaded successfully
			if ts.Site != nil {
				sitesByID[ts.Site.ID] = ts.Site
			}
		}
	}

	// Check TenantSite entry
	if len(siteIDs) > 0 {
		tsDAO := cdbm.NewTenantSiteDAO(gaih.dbSession)
		_, count, err := tsDAO.GetAll(
			ctx,
			nil,
			cdbm.TenantSiteFilterInput{
				TenantIDs: []uuid.UUID{tenant.ID},
				SiteIDs:   siteIDs,
			},
			cdbp.PageInput{},
			nil,
		)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving TenantSite entry")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to determine Tenant's association with Site", nil)
		}

		// We've ensured that the set of siteIDs is unique earlier,
		// so if the counts don't match, then something wasn't found.
		if count != len(siteIDs) {
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden,
				"Tenant is not associated with one or more of the Sites specified in query", nil)
		}
	}

	filter.SiteIDs = siteIDs

	// Get query text for full text search from query param
	if searchQueryStr := c.QueryParam("query"); searchQueryStr != "" {
		filter.SearchQuery = &searchQueryStr
		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.String("query", searchQueryStr), logger)
	}

	// Get status from query param
	if statusStrings := qParams["status"]; len(statusStrings) != 0 {
		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.StringSlice("status", statusStrings), logger)
		for _, status := range statusStrings {
			gaih.tracerSpan.SetAttribute(handlerSpan, attribute.String("status", status), logger)
			_, ok := cdbm.InstanceStatusMap[status]
			if !ok {
				logger.Warn().Msg(fmt.Sprintf("invalid value in status query: %v", status))
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Status value in query", nil)
			}
			filter.Statuses = append(filter.Statuses, status)
		}
	}

	// Get VPC IDs from query param
	if vpcIDStrs := qParams["vpcId"]; len(vpcIDStrs) != 0 {
		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.StringSlice("vpcId", vpcIDStrs), logger)
		for _, vpcIDStr := range vpcIDStrs {
			// Check for Vpc existence
			vpc, verr := common.GetVpcFromIDString(ctx, nil, vpcIDStr, nil, gaih.dbSession)
			if verr != nil {
				if verr == common.ErrInvalidID {
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid VPC ID %v in query", vpcIDStr), nil)
				}
				if verr == cdb.ErrDoesNotExist {
					return cerr.NewAPIErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Could not find VPC with ID %v specified in query", vpcIDStr), nil)
				}
				logger.Error().Err(verr).Msg("error retrieving Vpc from DB")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve VPC with ID %v specified in query", vpcIDStr), nil)
			}
			filter.VpcIDs = append(filter.VpcIDs, vpc.ID)
		}
	}

	// Get instance type IDs from query param
	if instanceTypeIDStrs := qParams["instanceTypeId"]; len(instanceTypeIDStrs) != 0 {
		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.StringSlice("instanceTypeId", instanceTypeIDStrs), logger)
		for _, instanceTypeStr := range instanceTypeIDStrs {
			// Check for instance type existence
			instanceType, verr := common.GetInstanceTypeFromIDString(ctx, nil, instanceTypeStr, gaih.dbSession)
			if verr != nil {
				if verr == common.ErrInvalidID {
					return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid instance type ID %v in query", instanceTypeStr), nil)
				}
				if verr == cdb.ErrDoesNotExist {
					return cerr.NewAPIErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Could not find instance type with ID %v specified in query", instanceTypeStr), nil)
				}
				logger.Error().Err(verr).Msg("error retrieving instance type from DB")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve instance type with ID %v specified in query", instanceTypeStr), nil)
			}
			filter.InstanceTypeIDs = append(filter.InstanceTypeIDs, instanceType.ID)
		}
	}

	// Get network security group IDs from query param
	if len(qParams["networkSecurityGroupId"]) > 0 {
		filter.NetworkSecurityGroupIDs = qParams["networkSecurityGroupId"]
	}

	// Get operating system IDs from query param
	if operatingSystemIDStrs := qParams["operatingSystemId"]; len(operatingSystemIDStrs) != 0 {
		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.StringSlice("operatingSystemId", operatingSystemIDStrs), logger)
		operatingSystemDAO := cdbm.NewOperatingSystemDAO(gaih.dbSession)
		for _, operatingSystemStr := range operatingSystemIDStrs {
			parsedID, err := uuid.Parse(operatingSystemStr)
			if err != nil {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid operating system ID %v in query", operatingSystemStr), nil)
			}

			// Check for operating system existence
			_, verr := operatingSystemDAO.GetByID(ctx, nil, parsedID, nil)
			if verr != nil {
				if verr == cdb.ErrDoesNotExist {
					return cerr.NewAPIErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Could not find operating system with ID %v specified in query", operatingSystemStr), nil)
				}
				logger.Error().Err(err).Msg("error retrieving operating system from DB")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve operating system with ID %v specified in query", operatingSystemStr), nil)
			}
			filter.OperatingSystemIDs = append(filter.OperatingSystemIDs, parsedID)
		}
	}

	// Get machine IDs from query param
	if machineIDs := qParams["machineId"]; len(machineIDs) != 0 {
		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.StringSlice("machineId", machineIDs), logger)
		machineDAO := cdbm.NewMachineDAO(gaih.dbSession)
		machines, _, err := machineDAO.GetAll(ctx, nil, cdbm.MachineFilterInput{MachineIDs: machineIDs}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving machines from DB")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve machines with IDs %v specified in query", strings.Join(machineIDs, ", ")), nil)
		}
		machineMap := map[string]bool{}
		for _, machine := range machines {
			machineMap[machine.ID] = true
		}
		hasValidMachineID := false
		for _, machineID := range machineIDs {
			_, ok := machineMap[machineID]
			if ok {
				filter.MachineIDs = append(filter.MachineIDs, machineID)
				hasValidMachineID = true
			}
		}
		if !hasValidMachineID {
			// Create pagination response header
			pageReponse := pagination.NewPageResponse(*pageRequest.PageNumber, *pageRequest.PageSize, 0, pageRequest.OrderByStr)
			pageHeader, err := json.Marshal(pageReponse)
			if err != nil {
				logger.Error().Err(err).Msg("error marshaling pagination response")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to generate pagination response header", nil)
			}
			c.Response().Header().Set(pagination.ResponseHeaderName, string(pageHeader))
			return c.JSON(http.StatusOK, []model.APIInstance{})
		}
	}

	// Get instance name from query param
	if name := c.QueryParam("name"); name != "" {
		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.String("name", name), logger)
		filter.Names = []string{name}
	}

	// Get IP addresses from query param and filter by interface IPs
	if ipAddresses := qParams["ipAddress"]; len(ipAddresses) != 0 {
		gaih.tracerSpan.SetAttribute(handlerSpan, attribute.StringSlice("ipAddress", ipAddresses), logger)

		// GetAll interfaces matching specified IP addresses
		ifcDAO := cdbm.NewInterfaceDAO(gaih.dbSession)
		matchingIfcs, _, err := ifcDAO.GetAll(ctx, nil, cdbm.InterfaceFilterInput{IPAddresses: ipAddresses}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Interfaces for IP filtering")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Interfaces for IP filtering", nil)
		}

		// Collect Instance ID attribute of matching Interfaces
		instanceIDsWithMatchingIPs := goset.NewSet[uuid.UUID]()
		for _, ifc := range matchingIfcs {
			instanceIDsWithMatchingIPs.Add(ifc.InstanceID)
		}

		// Add InstanceIDs to filter object
		if instanceIDsWithMatchingIPs.Cardinality() > 0 {
			filter.InstanceIDs = instanceIDsWithMatchingIPs.ToSlice()
		} else {
			// No instances match the IP filter, set empty list to get no results
			filter.InstanceIDs = []uuid.UUID{}
		}
	}

	// Get all Instances by Tenant, and Site, if specified
	instanceDAO := cdbm.NewInstanceDAO(gaih.dbSession)

	dbInstances, total, serr := instanceDAO.GetAll(ctx, nil,
		filter,
		cdbp.PageInput{
			Limit:   pageRequest.Limit,
			Offset:  pageRequest.Offset,
			OrderBy: pageRequest.OrderBy,
		},
		qIncludeRelations,
	)
	if serr != nil {
		logger.Error().Err(serr).Msg("error retrieving Instances for this Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instances for Site", nil)
	}

	// Get status details
	sdDAO := cdbm.NewStatusDetailDAO(gaih.dbSession)

	sdEntityIDs := []string{}
	insIDs := []uuid.UUID{}
	for _, ins := range dbInstances {
		sdEntityIDs = append(sdEntityIDs, ins.ID.String())
		insIDs = append(insIDs, ins.ID)
	}

	ssds, serr := sdDAO.GetRecentByEntityIDs(ctx, nil, sdEntityIDs, common.RECENT_STATUS_DETAIL_COUNT)
	if serr != nil {
		logger.Warn().Err(serr).Msg("error retrieving Status Details for Instances from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to populate status history for Instances", nil)
	}
	ssdMap := map[string][]cdbm.StatusDetail{}
	for _, ssd := range ssds {
		cssd := ssd
		ssdMap[ssd.EntityID] = append(ssdMap[ssd.EntityID], cssd)
	}

	// Create response
	ifcDAO := cdbm.NewInterfaceDAO(gaih.dbSession)
	ibifcDAO := cdbm.NewInfiniBandInterfaceDAO(gaih.dbSession)
	nvlDAO := cdbm.NewNVLinkInterfaceDAO(gaih.dbSession)
	skgiaDAO := cdbm.NewSSHKeyGroupInstanceAssociationDAO(gaih.dbSession)

	ifcs, _, serr := ifcDAO.GetAll(ctx, nil, cdbm.InterfaceFilterInput{InstanceIDs: insIDs}, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, []string{cdbm.SubnetRelationName, cdbm.VpcPrefixRelationName})
	if serr != nil {
		// Log error and continue
		logger.Error().Err(serr).Msg("error retrieving Instance Subnets for Instance from DB")
	}
	ifcMap := map[uuid.UUID][]cdbm.Interface{}
	for _, ifc := range ifcs {
		cifc := ifc
		ifcMap[ifc.InstanceID] = append(ifcMap[ifc.InstanceID], cifc)
	}

	// Get the instance infiniband interface record from the db
	ibifcs, _, serr := ibifcDAO.GetAll(
		ctx,
		nil,
		cdbm.InfiniBandInterfaceFilterInput{
			InstanceIDs: insIDs,
			SiteIDs:     siteIDs,
		},
		cdbp.PageInput{
			Limit: cdb.GetIntPtr(cdbp.TotalLimit),
		},
		[]string{cdbm.InfiniBandPartitionRelationName},
	)
	if serr != nil {
		logger.Error().Err(serr).Msg("error retrieving instance InfiniBand Interfaces Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance InfiniBand Interfaces for Instance", nil)
	}
	ibifcMap := map[uuid.UUID][]cdbm.InfiniBandInterface{}
	for _, ibifc := range ibifcs {
		cibifc := ibifc
		ibifcMap[ibifc.InstanceID] = append(ibifcMap[ibifc.InstanceID], cibifc)
	}

	// Get the instance NVLink Interface record from the db
	retnvlifc, _, serr := nvlDAO.GetAll(ctx, nil, cdbm.NVLinkInterfaceFilterInput{InstanceIDs: insIDs}, cdbp.PageInput{}, nil)
	if serr != nil {
		logger.Error().Err(serr).Msg("error retrieving instance NVLink interfaces Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance NVLink interfaces for Instance", nil)
	}

	// Get the instance NVLink Interface record from the db
	nvlifcMap := map[uuid.UUID][]cdbm.NVLinkInterface{}
	for _, nvlifc := range retnvlifc {
		cnvlifc := nvlifc
		nvlifcMap[nvlifc.InstanceID] = append(nvlifcMap[nvlifc.InstanceID], cnvlifc)
	}

	// Get DPU Extension Service Deployments for all instances
	desdDAO := cdbm.NewDpuExtensionServiceDeploymentDAO(gaih.dbSession)
	desds, _, err := desdDAO.GetAll(
		ctx,
		nil,
		cdbm.DpuExtensionServiceDeploymentFilterInput{
			InstanceIDs: insIDs,
		},
		cdbp.PageInput{
			Limit: cdb.GetIntPtr(cdbp.TotalLimit),
		},
		[]string{cdbm.DpuExtensionServiceRelationName},
	)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving DPU Extension Service Deployments for instances from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve DPU Extension Service Deployments for Instances", nil)
	}
	desdsMap := map[uuid.UUID][]cdbm.DpuExtensionServiceDeployment{}
	for _, desd := range desds {
		cdesd := desd
		desdsMap[desd.InstanceID] = append(desdsMap[desd.InstanceID], cdesd)
	}

	// Get SSH Key Group Instance Associations for all Instances
	skgias, _, err := skgiaDAO.GetAll(ctx, nil, nil, siteIDs, insIDs, []string{cdbm.SSHKeyGroupRelationName}, nil, cdb.GetIntPtr(cdbp.TotalLimit), nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving ssh key group instance association Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve SSH Key Group Instance Association for Instance", nil)
	}
	skgiasMap := map[uuid.UUID][]cdbm.SSHKeyGroup{}
	for _, skgia := range skgias {
		cskgia := skgia
		skgiasMap[skgia.InstanceID] = append(skgiasMap[skgia.InstanceID], *cskgia.SSHKeyGroup)
	}

	// We'll need to pull the VPC details for any instances that aren't setting NSG ID
	// to decide if they're inheriting one from their parent VPC, and then figure out
	// their propagation status.  This needs to done separately because we can't assume
	// that the user requested the 	VPC relation.

	inheritVpcIDs := goset.NewSet[uuid.UUID]()

	for _, ins := range dbInstances {
		// Only instances with no NSG attached directly
		// could possible be inheriting from their VPC.
		if ins.NetworkSecurityGroupID == nil {
			inheritVpcIDs.Add(ins.VpcID)
		}
	}

	vpcs := map[uuid.UUID]*cdbm.Vpc{}

	// Only if there's at least one possible case
	// of inheritence
	if !inheritVpcIDs.IsEmpty() {

		vpcDAO := cdbm.NewVpcDAO(gaih.dbSession)

		vpcFilter := cdbm.VpcFilterInput{
			VpcIDs: inheritVpcIDs.ToSlice(),
		}

		dbVpcs, _, err := vpcDAO.GetAll(ctx, nil, vpcFilter, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving VPCs DB entities")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve VPCs for Instances", nil)
		}

		for _, vpc := range dbVpcs {
			vpcs[vpc.ID] = &vpc
		}
	}

	apiInstances := []model.APIInstance{}
	for _, ins := range dbInstances {
		// Create response
		dbInstance := ins
		apiInstance := model.NewAPIInstance(&dbInstance, sitesByID[dbInstance.SiteID], ifcMap[dbInstance.ID], ibifcMap[dbInstance.ID], desdsMap[dbInstance.ID], nvlifcMap[dbInstance.ID], skgiasMap[dbInstance.ID], ssdMap[ins.ID.String()])

		if ins.NetworkSecurityGroupID == nil {

			vpc, exists := vpcs[ins.VpcID]

			// We've only inherited if our NSG ID is null _and_ the parent
			// NSG ID is not null.
			if exists {
				apiInstance.NetworkSecurityGroupInherited = vpc.NetworkSecurityGroupID != nil

				// If we inherited our NSG, then see if we're propagated.
				if apiInstance.NetworkSecurityGroupInherited {

					// We can default to non/configuring and then switch
					// if we're actually propagated.
					apiInstance.NetworkSecurityGroupPropagationDetails = &model.APINetworkSecurityGroupPropagationDetails{
						ObjectID:       apiInstance.ID,
						DetailedStatus: model.APINetworkSecurityGroupPropagationDetailedStatusNone,
						Status:         model.APINetworkSecurityGroupPropagationStatusSynchronizing,
					}

					if vpc.NetworkSecurityGroupPropagationDetails != nil {
						// If the instance wasn't found in the list of unpropagated instances, then we're propagated.
						if !slices.Contains(vpc.NetworkSecurityGroupPropagationDetails.GetUnpropagatedInstanceIds(), apiInstance.ID) {
							apiInstance.NetworkSecurityGroupPropagationDetails.DetailedStatus = model.APINetworkSecurityGroupPropagationDetailedStatusFull
							apiInstance.NetworkSecurityGroupPropagationDetails.Status = model.APINetworkSecurityGroupPropagationStatusSynchronized
						}
					}
				}
			}

		}

		apiInstances = append(apiInstances, *apiInstance)
	}

	// Create pagination response header
	pageReponse := pagination.NewPageResponse(*pageRequest.PageNumber, *pageRequest.PageSize, total, pageRequest.OrderByStr)
	pageHeader, err := json.Marshal(pageReponse)
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling pagination response")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to generate pagination response header", nil)
	}
	c.Response().Header().Set(pagination.ResponseHeaderName, string(pageHeader))

	logger.Info().Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiInstances)
}

// ~~~~~ Delete Handler ~~~~~ //

// DeleteInstanceHandler is the API Handler for deleting an Instance
type DeleteInstanceHandler struct {
	dbSession  *cdb.Session
	tc         temporalClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewDeleteInstanceHandler initializes and r`eturns a new handler for deleting an Instance
func NewDeleteInstanceHandler(dbSession *cdb.Session, tc temporalClient.Client, scp *sc.ClientPool, cfg *config.Config) DeleteInstanceHandler {
	return DeleteInstanceHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Delete an Instance
// @Description Delete an Instance fro the org
// @Tags Instance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Instance"
// @Success 202
// @Router /v2/org/{org}/carbide/instance/{id} [delete]
func (dih DeleteInstanceHandler) Handle(c echo.Context) error {
	// Get context
	ctx := c.Request().Context()

	// Get org
	org := c.Param("orgName")

	// Initialize logger
	logger := log.With().Str("Model", "Instance").Str("Handler", "Delete").Str("Org", org).Logger()

	logger.Info().Msg("started API handler")

	// Create a child span and set the attributes for current request
	newctx, handlerSpan := dih.tracerSpan.CreateChildInContext(ctx, "DeleteInstanceHandler", logger)
	if handlerSpan != nil {
		// Set newly created span context as a current context
		ctx = newctx

		defer handlerSpan.End()

		dih.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, dih.tracerSpan, handlerSpan)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Tenant Admins are allowed to delete Instances
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.TenantAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Tenant Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Tenant Admin role with org", nil)
	}

	// Get Instance ID from URL param
	instanceStrID := c.Param("id")
	instanceID, err := uuid.Parse(instanceStrID)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Instance ID in URL", nil)
	}

	dih.tracerSpan.SetAttribute(handlerSpan, attribute.String("instance_id", instanceStrID), logger)

	// Get Instance
	instanceDAO := cdbm.NewInstanceDAO(dih.dbSession)

	instance, err := instanceDAO.GetByID(ctx, nil, instanceID, []string{cdbm.SiteRelationName, cdbm.TenantRelationName})
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIErrorResponse(c, http.StatusNotFound, "Could not find Instance with specified ID", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Instance from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance", nil)
	}

	if instance.Tenant == nil {
		logger.Error().Err(err).Msg("error retrieving Tenant")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant details", nil)
	}

	// Confirm that the Instance org (via the Tenant org)
	// matches the org sent in the request.
	if instance.Tenant.Org != org {
		logger.Error().Msg("org specified in request does not match org of Tenant associated with Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Org specified in request does not match Org of Tenant associated with Instance", nil)
	}

	// Verify that the instance is associated with a site and then that the site is
	// in a valid state.
	if instance.Site == nil {
		logger.Error().Msg("failed to pull site data for Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site details for Instance", nil)
	}

	if instance.Site.Status != cdbm.SiteStatusRegistered {
		logger.Error().Msg("site not in registered state - cannot delete Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Site is not in Registered state - cannot delete Instance", nil)
	}

	// Bind request data to API model
	apiRequest := model.APIInstanceDeleteRequest{}
	if c.Request().Body != http.NoBody {
		if err := c.Bind(&apiRequest); err != nil {
			logger.Warn().Err(err).Msg("error binding request data into API model")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data, potentially invalid structure", nil)
		}
	}

	// Validate request attributes
	if verr := apiRequest.Validate(); verr != nil {
		logger.Warn().Err(verr).Msg("error validating Instance deletion request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Error validating Instance deletion request data", verr)
	}

	// Start a DB transaction
	tx, err := cdb.BeginTx(ctx, dih.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to delete Instance", nil)
	}

	// This variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// Update Instance to set status to Deleting
	_, err = instanceDAO.Update(ctx, tx, cdbm.InstanceUpdateInput{InstanceID: instance.ID, Status: cdb.GetStrPtr(cdbm.InstanceStatusTerminating)})
	if err != nil {
		logger.Error().Err(err).Msg("error updating Instance in DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to delete Instance", nil)
	}

	// Create status detail
	sdDAO := cdbm.NewStatusDetailDAO(dih.dbSession)
	_, err = sdDAO.CreateFromParams(ctx, tx, instance.ID.String(), *cdb.GetStrPtr(cdbm.InstanceStatusTerminating),
		cdb.GetStrPtr("Instance deletion successfully initiated on Site"))
	if err != nil {
		logger.Error().Err(err).Msg("error creating Status Detail DB entry")
	}

	// Get the temporal client for the site we are working with.
	stc, err := dih.scp.GetClientByID(instance.SiteID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Prepare the delete/release request workflow object
	releaseInstanceRequest := &cwssaws.InstanceReleaseRequest{
		Id: &cwssaws.InstanceId{Value: common.GetSiteInstanceID(instance).String()},
	}

	// This is for enhanced break-fix flow:
	if apiRequest.MachineHealthIssue != nil {
		releaseInstanceRequest.Issue = &cwssaws.Issue{
			Category: cwssaws.IssueCategory(model.MachineIssueCategoriesFromAPIToProtobuf[apiRequest.MachineHealthIssue.Category]),
		}
		if apiRequest.MachineHealthIssue.Summary != nil {
			releaseInstanceRequest.Issue.Summary = *apiRequest.MachineHealthIssue.Summary
		}
		if apiRequest.MachineHealthIssue.Details != nil {
			releaseInstanceRequest.Issue.Details = *apiRequest.MachineHealthIssue.Details
		}
	}
	// if caller attempt to set IsRepairTenant then it must be a tenant with targetedInstanceCreation capability
	if apiRequest.IsRepairTenant != nil && *apiRequest.IsRepairTenant {
		if instance.Tenant.Config == nil || !instance.Tenant.Config.TargetedInstanceCreation {
			logger.Warn().Msg("tenant does not have capability to set IsRepairTenant")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Tenant does not have capability to set IsRepairTenant", nil)
		}
		releaseInstanceRequest.IsRepairTenant = apiRequest.IsRepairTenant
	}

	workflowOptions := temporalClient.StartWorkflowOptions{
		ID:                       "instance-delete-" + instance.ID.String(),
		TaskQueue:                queue.SiteTaskQueue,
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
	}

	logger.Info().Msg("triggering instance delete workflow")

	// Add context deadline.
	// The client (NGC or its downstream client) could cancel the parent deadline
	// at any time by closing the connection, HTTP2 reset, etc.  So, the real
	// deadline could be shorter than WorkflowContextTimeout.  We're only
	// enforcing an upper limit here.
	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	// Trigger Site workflow to update instance
	// TODO: Once Site Agent offers DeleteInstanceV2 re-registered as DeleteInstance then update workflow name here
	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "DeleteInstanceV2", releaseInstanceRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to synchronously start Temporal workflow to delete Instance")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to start sync workflow to delete Instance on Site: %s", err), nil)
	}

	wid := we.GetID()
	logger.Info().Str("Workflow ID", wid).Msg("executed synchronous delete Instance workflow")

	// Execute the workflow synchronously
	err = we.Get(ctx, nil)

	// Handle skippable errors
	if err != nil {
		// If this was a 404 back from Carbide, we can treat the object as already having been deleted and allow things to proceed.
		var applicationErr *tp.ApplicationError
		if errors.As(err, &applicationErr) && applicationErr.Type() == swe.ErrTypeCarbideObjectNotFound {
			logger.Warn().Msg(swe.ErrTypeCarbideObjectNotFound + " received from Site")
			// Reset error to nil
			err = nil
		}
	}

	// Check if err is still nil now that we've handled any skippable errors.
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || ctx.Err() != nil {

			logger.Error().Err(err).Msg("failed to delete Instance, timeout occurred executing workflow on Site.")

			// Create a new context deadlines
			newctx, newcancel := context.WithTimeout(context.Background(), common.WorkflowContextNewAfterTimeout)
			defer newcancel()

			// Initiate termination workflow
			serr := stc.TerminateWorkflow(newctx, wid, "", "timeout occurred executing delete Instance workflow")
			if serr != nil {
				logger.Error().Err(serr).Msg("failed to terminate Temporal workflow for deleting Instance")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to terminate synchronous Instance deletion workflow after timeout, Cloud and Site data may be de-synced: %s", serr), nil)
			}

			logger.Info().Str("Workflow ID", wid).Msg("initiated terminate synchronous delete Instance workflow successfully")

			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to delete Instance, timeout occurred executing workflow on Site: %s", err), nil)
		}

		code, err := common.UnwrapWorkflowError(err)
		logger.Error().Err(err).Msg("failed to synchronously execute Temporal workflow to delete Instance")
		return cerr.NewAPIErrorResponse(c, code, fmt.Sprintf("Failed to execute sync workflow to delete Instance on Site: %s", err), nil)
	}

	logger.Info().Str("Workflow ID", wid).Msg("completed synchronous delete Instance workflow")

	// Commit the DB transaction after the synchronous workflow has completed without error
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing instance transaction to DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to delete Instance, DB transaction error", nil)
	}

	// Set committed so, deferred cleanup functions will do nothing
	txCommitted = true

	// Return response
	logger.Info().Msg("finishing API handler")

	return c.String(http.StatusAccepted, "Deletion request was accepted")
}

// GetInstanceStatusDetailsHandler is the API Handler for getting Instance StatusDetail records
type GetInstanceStatusDetailsHandler struct {
	dbSession  *cdb.Session
	tracerSpan *sutil.TracerSpan
}

// NewGetInstanceStatusDetailsHandler initializes and returns a new handler to retrieve Instance StatusDetail records
func NewGetInstanceStatusDetailsHandler(dbSession *cdb.Session) GetInstanceStatusDetailsHandler {
	return GetInstanceStatusDetailsHandler{
		dbSession:  dbSession,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Get Instance StatusDetails
// @Description Get all StatusDetails for Instance
// @Tags instance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Instance"
// @Success 200 {object} []model.APIStatusDetail
// @Router /v2/org/{org}/carbide/instance/{id}/status-history [get]
func (gisdh GetInstanceStatusDetailsHandler) Handle(c echo.Context) error {
	// Get context
	ctx := c.Request().Context()

	// Get org
	org := c.Param("orgName")

	// Initialize logger
	logger := log.With().Str("Model", "Instance").Str("Handler", "Get").Str("Org", org).Logger()

	logger.Info().Msg("started API handler")

	// Create a child span and set the attributes for current request
	newctx, handlerSpan := gisdh.tracerSpan.CreateChildInContext(ctx, "GetInstanceStatusDetailsHandler", logger)
	if handlerSpan != nil {
		// Set newly created span context as a current context
		ctx = newctx
		defer handlerSpan.End()
		gisdh.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, gisdh.tracerSpan, handlerSpan)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Tenant Admins are allowed to retrieve Instances
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.TenantAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Tenant Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Tenant Admin role with org", nil)
	}

	// Get and validate includeRelation params
	qParams := c.QueryParams()
	qIncludeRelations, errMsg := common.GetAndValidateQueryRelations(qParams, cdbm.InstanceRelatedEntities)
	if errMsg != "" {
		logger.Warn().Msg(errMsg)
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, errMsg, nil)
	}

	// Get Instance ID from URL param
	instanceStrID := c.Param("id")
	instanceID, err := uuid.Parse(instanceStrID)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Instance ID in URL", nil)
	}

	gisdh.tracerSpan.SetAttribute(handlerSpan, attribute.String("instance_id", instanceStrID), logger)

	// Get Instance
	instanceDAO := cdbm.NewInstanceDAO(gisdh.dbSession)
	instance, err := instanceDAO.GetByID(ctx, nil, instanceID, qIncludeRelations)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIErrorResponse(c, http.StatusNotFound, "Could not find Instance with specified ID", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Instance from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance", nil)
	}

	// Get Tenant for this org
	tnDAO := cdbm.NewTenantDAO(gisdh.dbSession)
	tenants, err := tnDAO.GetAllByOrg(ctx, nil, org, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Tenant for this org")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant", nil)
	}

	if len(tenants) == 0 {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Org does not have a Tenant associated", nil)
	}
	tenant := tenants[0]

	// Check if Instance belongs to Tenant
	if instance.TenantID != tenant.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Instance does not belong to current Tenant", nil)
	}

	// handle retrieving and building status details response
	apiSds, err := handleEntityStatusDetails(ctx, c, gisdh.dbSession, instanceID.String(), logger)
	if err != nil {
		return err
	}

	logger.Info().Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiSds)
}
