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

package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	temporalClient "go.temporal.io/sdk/client"
	tp "go.temporal.io/sdk/temporal"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/nvidia/bare-metal-manager-rest/api/internal/config"
	common "github.com/nvidia/bare-metal-manager-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/model"
	sc "github.com/nvidia/bare-metal-manager-rest/api/pkg/client/site"
	auth "github.com/nvidia/bare-metal-manager-rest/auth/pkg/authorization"
	cerr "github.com/nvidia/bare-metal-manager-rest/common/pkg/util"
	sutil "github.com/nvidia/bare-metal-manager-rest/common/pkg/util"
	cdb "github.com/nvidia/bare-metal-manager-rest/db/pkg/db"
	cdbm "github.com/nvidia/bare-metal-manager-rest/db/pkg/db/model"
	cdbp "github.com/nvidia/bare-metal-manager-rest/db/pkg/db/paginator"
	cwssaws "github.com/nvidia/bare-metal-manager-rest/workflow-schema/schema/site-agent/workflows/v1"
	"github.com/nvidia/bare-metal-manager-rest/workflow/pkg/queue"
)

const (
	// batchSuffixLength is the length of the random suffix used for batch instance names
	batchSuffixLength = 6
)

// ~~~~~ Batch Create Handler ~~~~~ //

// BatchCreateInstanceHandler is the API Handler for creating multiple instances with topology-optimized allocation
type BatchCreateInstanceHandler struct {
	dbSession  *cdb.Session
	tc         temporalClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewBatchCreateInstanceHandler initializes and returns a new handler for batch creating Instances
func NewBatchCreateInstanceHandler(dbSession *cdb.Session, tc temporalClient.Client, scp *sc.ClientPool, cfg *config.Config) BatchCreateInstanceHandler {
	return BatchCreateInstanceHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Batch create multiple Instances with topology-optimized allocation
// @Description Create multiple Instances with topology-optimized machine allocation. If topologyOptimized is true, all instances must be allocated on the same rack/topology domain (e.g., for NVLink). If false, instances can be spread across racks.
// @Tags Instance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param message body model.APIBatchInstanceCreateRequest true "Batch instance creation request"
// @Success 201 {object} []model.APIInstance
// @Router /v2/org/{org}/carbide/instance/batch [post]
func (bcih BatchCreateInstanceHandler) Handle(c echo.Context) error {
	// Execution Steps:
	// 1. Authentication & Authorization
	//    - Extract user from context
	//    - Validate org membership
	//    - Validate Tenant Admin role
	// 2. Request Validation
	//    - Bind and validate batch request data (count, namePrefix, topology flag)
	//    - Validate tenant, instance type, VPC, site
	//    - Load and validate Interfaces (Subnets, VPC Prefixes) - shared across all instances
	//    - Load and validate DPU Extension Service Deployments - shared across all instances
	//    - Load and validate Network Security Groups - shared across all instances
	//    - Load and validate SSH Key Groups - shared across all instances
	//    - Validate OS or iPXE script
	//    - Generate unique instance names and check for conflicts
	// 3. Database Transaction
	//    - Start transaction
	//    - Acquire advisory lock on Tenant + Instance Type
	// 4. Machine Selection
	//    - Verify allocation constraints have sufficient quota for batch
	//    - Build list of available allocation constraints with capacity
	//    - Allocate multiple machines (topology-optimized or cross-rack)
	//    - Mark all allocated machines as assigned with per-machine advisory locks
	// 5. Machine Capability Validation
	//    - Validate InfiniBand interfaces against Instance Type capabilities - shared across all instances
	//    - Validate InfiniBand partitions (Site, Tenant, Status)
	//    - Validate DPU interfaces against Instance Type capabilities - shared across all instances
	//    - Validate NVLink interfaces against Instance Type capabilities - shared across all instances
	//    - Validate NVLink logical partitions (Site, Tenant, Status)
	// 6. Create Instance Records (loop for each instance)
	//    - Create Instance record with allocated machine
	//    - Update ControllerInstanceID
	//    - Create SSH Key Group associations
	//    - Create Interface records
	//    - Create InfiniBand Interface records
	//    - Create NVLink Interface records
	//    - Create DPU Extension Service Deployment records
	//    - Create status detail record
	//    - Switch to next allocation constraint when current reaches capacity
	// 7. Workflow Trigger
	//    - Build batch instance allocation request with all configs
	//    - Execute synchronous Temporal workflow (CreateInstances)
	//    - Wait for site-agent to provision all instances
	//    - Handle timeout with workflow termination
	// 8. Commit & Response
	//    - Commit transaction after workflow succeeds
	//    - Return array of created instances to client
	//
	// Key Differences from Single Instance API:
	// - Creates multiple instances in single transaction (all-or-nothing)
	// - Topology-optimized machine allocation (same rack or cross-rack based on flag)
	// - Automatically distributes instances across multiple allocation constraints
	// - Shared interface configurations across all instances (Interfaces, InfiniBandInterfaces, NVLinkInterfaces)
	// - Batch workflow for atomic provisioning of all instances

	// ==================== Step 1: Authentication & Authorization ====================

	// Get context
	ctx := c.Request().Context()

	// Get org
	org := c.Param("orgName")

	// Initialize logger
	logger := log.With().Str("Model", "Instance").Str("Handler", "BatchCreate").Str("Org", org).Logger()

	logger.Info().Msg("started API handler for batch instance creation")

	// Create a child span and set the attributes for current request
	newctx, handlerSpan := bcih.tracerSpan.CreateChildInContext(ctx, "BatchCreateInstanceHandler", logger)
	if handlerSpan != nil {
		// Set newly created span context as a current context
		ctx = newctx

		defer handlerSpan.End()

		bcih.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, bcih.tracerSpan, handlerSpan)
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
	apiRequest := model.APIBatchInstanceCreateRequest{}
	err = c.Bind(&apiRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data, potentially invalid structure", nil)
	}

	// Validate request attributes
	verr := apiRequest.Validate()
	if verr != nil {
		logger.Warn().Err(verr).Msg("error validating batch instance creation request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Error validating batch instance creation request data", verr)
	}

	// Set default for TopologyOptimized if not provided
	// Default to true for better performance and locality
	topologyOptimized := true
	if apiRequest.TopologyOptimized != nil {
		topologyOptimized = *apiRequest.TopologyOptimized
	}

	logger.Info().Int("Count", apiRequest.Count).Bool("TopologyOptimized", topologyOptimized).Msg("Input validation completed for batch Instance creation request")

	icv := common.NewInstanceCreateValidator(bcih.dbSession, bcih.cfg, &logger)

	tenant, vpc, site, defaultNvllpID, apiErr := icv.ValidateTenantAndVPC(ctx, org, apiRequest.TenantID, apiRequest.VpcID)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}

	// Validate the instance type (batch-specific: required in request)
	apiInstanceTypeID, err := uuid.Parse(apiRequest.InstanceTypeID)
	if err != nil {
		logger.Warn().Err(err).Msg("error parsing instance type id in request")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Instance Type ID in request is not valid", nil)
	}

	itDAO := cdbm.NewInstanceTypeDAO(bcih.dbSession)
	instancetype, err := itDAO.GetByID(ctx, nil, apiInstanceTypeID, nil)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Could not find Instance Type with ID specified in request data", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Instance Type from DB by ID")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance Type with ID specified in request data", nil)
	}

	// Verify VPC and InstanceType are on the same Site
	if vpc.SiteID != *instancetype.SiteID {
		logger.Warn().
			Str("Site ID for VPC", vpc.SiteID.String()).
			Str("Site ID for Instance Type", instancetype.SiteID.String()).
			Msg("VPC and InstanceType are not on the same Site")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "VPC and Instance Type specified in request data do not belong to the same Site", nil)
	}

	ifcResult, apiErr := icv.ValidateInterfaces(ctx, tenant, vpc, apiRequest.Interfaces)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}

	logger.Info().Int("uniqueSubnetCount", len(ifcResult.SubnetIDMap)).Int("uniqueVpcPrefixCount", len(ifcResult.VpcPrefixIDMap)).
		Msg("validated all Subnets and VPC Prefixes (shared across all instances)")

	_, apiErr = icv.ValidateDPUExtensionServices(ctx, tenant, site, apiRequest.DpuExtensionServiceDeployments)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}

	logger.Info().Int("dpuExtensionServiceCount", len(apiRequest.DpuExtensionServiceDeployments)).
		Msg("validated DPU Extension Service Deployments")

	if apiErr := icv.ValidateNSG(ctx, tenant, site, apiRequest.NetworkSecurityGroupID); apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}

	sshKeyGroups, apiErr := icv.ValidateSSHKeyGroups(ctx, tenant, site, apiRequest.SSHKeyGroupIDs)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}

	logger.Info().Int("sshKeyGroupCount", len(sshKeyGroups)).
		Msg("validated SSH Key Groups")

	osConfig, osID, apiErr := icv.BuildOsConfig(ctx, &apiRequest, site.ID)
	if apiErr != nil {
		logger.Error().Err(errors.New(apiErr.Message)).Msg("error building os config for creating Instances")
		return c.JSON(apiErr.Code, apiErr)
	}

	// Generate instance names with random suffix to avoid name conflicts
	// Format: namePrefix-randomSuffix-index (e.g., myapp-a1b2c3-1)
	batchSuffix := uuid.New().String()[:batchSuffixLength]
	generateInstanceName := func(index int) string {
		return fmt.Sprintf("%s-%s-%d", apiRequest.NamePrefix, batchSuffix, index+1)
	}

	// Check for name conflicts before allocating any resources (safety check)
	// Build all instance names first, then check in a single batch query
	inDAO := cdbm.NewInstanceDAO(bcih.dbSession)
	allInstanceNames := make([]string, apiRequest.Count)
	for i := 0; i < apiRequest.Count; i++ {
		allInstanceNames[i] = generateInstanceName(i)
	}

	// Single batch query to check all names at once
	existing, tot, err := inDAO.GetAll(ctx, nil,
		cdbm.InstanceFilterInput{
			Names:     allInstanceNames,
			TenantIDs: []uuid.UUID{tenant.ID},
			SiteIDs:   []uuid.UUID{site.ID},
		},
		cdbp.PageInput{}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error checking for name uniqueness")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to check instance name uniqueness, DB error", nil)
	}
	if tot > 0 {
		logger.Warn().Str("existingInstanceName", existing[0].Name).Str("existingInstanceID", existing[0].ID.String()).
			Msg("instance with same name already exists for tenant at site")
		return cerr.NewAPIErrorResponse(c, http.StatusConflict,
			fmt.Sprintf("Instance with name '%s' already exists for Tenant at this Site. Please choose a different name prefix.", existing[0].Name), nil)
	}

	logger.Info().Int("instanceCount", apiRequest.Count).Str("batchSuffix", batchSuffix).
		Msg("validated instance names - all pre-transaction validations completed successfully")

	// ==================== Step 3: Database Transaction ====================

	// Start a db tx
	tx, err := cdb.BeginTx(ctx, bcih.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Instances", nil)
	}
	// this variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// ==================== Step 4: Machine Selection ====================

	// Acquire an advisory lock on the tenant ID and instancetype ID
	// This prevents concurrent instance creation (single or batch) from the same tenant on the same instance type
	err = tx.TryAcquireAdvisoryLock(ctx, cdb.GetAdvisoryLockIDFromString(fmt.Sprintf("%s-%s", tenant.ID.String(), instancetype.ID.String())), nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to acquire advisory lock on Tenant and Instance Type")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Error creating Instances, detected multiple parallel requests on Instance Type by Tenant", nil)
	}

	// Ensure that Tenant has an Allocation with specified Tenant InstanceType Site
	aDAO := cdbm.NewAllocationDAO(bcih.dbSession)
	allocationFilter := cdbm.AllocationFilterInput{TenantIDs: []uuid.UUID{tenant.ID}, SiteIDs: []uuid.UUID{*instancetype.SiteID}}
	allocationPage := cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}
	tnas, _, serr := aDAO.GetAll(ctx, tx, allocationFilter, allocationPage, nil)
	if serr != nil {
		logger.Error().Err(serr).Msg("error retrieving allocations for tenant")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Allocation for Tenant", nil)
	}
	if len(tnas) == 0 {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden,
			"Tenant does not have any Allocations for Site and Instance Type specified in request data", nil)
	}

	alconstraints, err := common.GetAllocationConstraintsForInstanceType(ctx, tx, bcih.dbSession, tenant.ID, instancetype, tnas)
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

	// Check if we have enough allocation for all requested instances
	if insTotal+apiRequest.Count > totalConstraintValue {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden,
			fmt.Sprintf("Tenant has reached the maximum number of Instances for Instance Type. Current: %d, Requested: %d, Max: %d", insTotal, apiRequest.Count, totalConstraintValue), nil)
	}

	// Build list of allocation constraints with available capacity
	// Instances will be distributed across multiple constraints as needed
	availableConstraints := []struct {
		constraint *cdbm.AllocationConstraint
		available  int
	}{}

	for _, alc := range alconstraints {
		used := usedMapAllocationConstraintIDs[alc.ID]
		available := alc.ConstraintValue - used
		if available > 0 {
			availableConstraints = append(availableConstraints, struct {
				constraint *cdbm.AllocationConstraint
				available  int
			}{
				constraint: &alc,
				available:  available,
			})
		}
	}

	// Allocate machines with topology optimization
	machines, apiErr := allocateMachinesForBatch(ctx, tx, bcih.dbSession, instancetype, apiRequest.Count, topologyOptimized, logger)
	if apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, apiErr.Data)
	}

	// ==================== Step 5: Machine Capability Validation ====================

	dbibic, apiErr := icv.ValidateIBInterfaces(ctx, tenant, site, apiRequest.InfiniBandInterfaces, apiInstanceTypeID)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}

	if apiErr := icv.ValidateDPUInterfaces(ctx, ifcResult.DBInterfaces, ifcResult.IsDeviceInfoPresent, apiInstanceTypeID); apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}

	dbnvlic, apiErr := icv.ValidateNVLinkInterfaces(ctx, tenant, site, defaultNvllpID, apiRequest.NVLinkInterfaces, apiInstanceTypeID)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}

	logger.Info().Msg("completed machine capability validation (Step 5)")

	// ==================== Step 6: Batch Instance Creation (Optimized with Batch DB Operations) ====================

	// Build SSH key group IDs for temporal workflow (shared by all instances)
	instanceSshKeyGroupIds := make([]string, 0, len(sshKeyGroups))
	for _, skg := range sshKeyGroups {
		instanceSshKeyGroupIds = append(instanceSshKeyGroupIds, skg.ID.String())
	}

	// Pre-parse DPU Extension Service IDs (shared validation, done once)
	dpuServiceIDs := make([]uuid.UUID, 0, len(apiRequest.DpuExtensionServiceDeployments))
	for _, adesdr := range apiRequest.DpuExtensionServiceDeployments {
		desdID, err := uuid.Parse(adesdr.DpuExtensionServiceID)
		if err != nil {
			logger.Warn().Err(err).Str("serviceID", adesdr.DpuExtensionServiceID).
				Msg("error parsing DPU Extension Service ID")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest,
				fmt.Sprintf("Invalid DPU Extension Service ID: %s", adesdr.DpuExtensionServiceID), nil)
		}
		dpuServiceIDs = append(dpuServiceIDs, desdID)
	}

	// --- Build all InstanceCreateInputs ---
	instanceCreateInputs := make([]cdbm.InstanceCreateInput, 0, len(machines))
	constraintIdx := 0
	constraintUsedCount := 0

	for i, machine := range machines {
		// Check if we need to switch to the next allocation constraint
		if constraintUsedCount >= availableConstraints[constraintIdx].available {
			constraintIdx++
			constraintUsedCount = 0
			if constraintIdx >= len(availableConstraints) {
				logger.Error().Msg("ran out of allocation constraints (should not happen)")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
					"Failed to allocate instances: insufficient allocation constraints", nil)
			}
		}
		currentConstraint := availableConstraints[constraintIdx].constraint
		constraintUsedCount++

		instanceCreateInputs = append(instanceCreateInputs, cdbm.InstanceCreateInput{
			Name:                     generateInstanceName(i),
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
			InstanceTypeID:           &apiInstanceTypeID,
			AllocationID:             &currentConstraint.AllocationID,
			AllocationConstraintID:   &currentConstraint.ID,
			IsUpdatePending:          false,
			Status:                   cdbm.InstanceStatusPending,
			PowerStatus:              cdb.GetStrPtr(cdbm.InstancePowerStatusRebooting),
			CreatedBy:                dbUser.ID,
		})
	}

	// --- Batch create all instances ---
	createdInstances, err := inDAO.CreateMultiple(ctx, tx, instanceCreateInputs)
	if err != nil {
		logger.Error().Err(err).Msg("failed to batch create instance records")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
			fmt.Sprintf("Failed to batch create instances: %v", err), nil)
	}
	logger.Info().Int("count", len(createdInstances)).Msg("batch created all instance records")

	// --- Build and batch update ControllerInstanceIDs ---
	instanceUpdateInputs := make([]cdbm.InstanceUpdateInput, 0, len(createdInstances))
	for _, inst := range createdInstances {
		instanceUpdateInputs = append(instanceUpdateInputs, cdbm.InstanceUpdateInput{
			InstanceID:           inst.ID,
			ControllerInstanceID: cdb.GetUUIDPtr(inst.ID),
		})
	}

	updatedInstances, err := inDAO.UpdateMultiple(ctx, tx, instanceUpdateInputs)
	if err != nil {
		logger.Error().Err(err).Msg("failed to batch update controller instance IDs")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
			fmt.Sprintf("Failed to batch update instances: %v", err), nil)
	}
	logger.Info().Int("count", len(updatedInstances)).Msg("batch updated all controller instance IDs")

	// Build instance map for easy lookup
	instanceMap := make(map[uuid.UUID]*cdbm.Instance, len(updatedInstances))
	for i := range updatedInstances {
		instanceMap[updatedInstances[i].ID] = &updatedInstances[i]
	}

	// --- Build and batch create SSH Key Group Instance Associations ---
	skgiaInputs := make([]cdbm.SSHKeyGroupInstanceAssociationCreateInput, 0, len(updatedInstances)*len(sshKeyGroups))
	for _, inst := range updatedInstances {
		for _, skg := range sshKeyGroups {
			skgiaInputs = append(skgiaInputs, cdbm.SSHKeyGroupInstanceAssociationCreateInput{
				SSHKeyGroupID: skg.ID,
				SiteID:        site.ID,
				InstanceID:    inst.ID,
				CreatedBy:     dbUser.ID,
			})
		}
	}

	if len(skgiaInputs) > 0 {
		skgiaDAO := cdbm.NewSSHKeyGroupInstanceAssociationDAO(bcih.dbSession)
		_, err = skgiaDAO.CreateMultiple(ctx, tx, skgiaInputs)
		if err != nil {
			logger.Error().Err(err).Msg("failed to batch create SSH key group associations")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
				fmt.Sprintf("Failed to batch create SSH key group associations: %v", err), nil)
		}
		logger.Info().Int("count", len(skgiaInputs)).Msg("batch created all SSH key group associations")
	}

	// --- Build and batch create Interfaces ---
	ifcInputs := make([]cdbm.InterfaceCreateInput, 0, len(updatedInstances)*len(ifcResult.DBInterfaces))
	for _, inst := range updatedInstances {
		for _, dbifc := range ifcResult.DBInterfaces {
			ifcInputs = append(ifcInputs, cdbm.InterfaceCreateInput{
				InstanceID:        inst.ID,
				SubnetID:          dbifc.SubnetID,
				VpcPrefixID:       dbifc.VpcPrefixID,
				Device:            dbifc.Device,
				DeviceInstance:    dbifc.DeviceInstance,
				VirtualFunctionID: dbifc.VirtualFunctionID,
				IsPhysical:        dbifc.IsPhysical,
				Status:            cdbm.InterfaceStatusPending,
				CreatedBy:         dbUser.ID,
			})
		}
	}

	var createdIfcsAll []cdbm.Interface
	if len(ifcInputs) > 0 {
		ifcDAO := cdbm.NewInterfaceDAO(bcih.dbSession)
		createdIfcsAll, err = ifcDAO.CreateMultiple(ctx, tx, ifcInputs)
		if err != nil {
			logger.Error().Err(err).Msg("failed to batch create interfaces")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
				fmt.Sprintf("Failed to batch create interfaces: %v", err), nil)
		}
		logger.Info().Int("count", len(createdIfcsAll)).Msg("batch created all interfaces")
	}

	// --- Build and batch create InfiniBand Interfaces ---
	var createdIbIfcsAll []cdbm.InfiniBandInterface
	if len(dbibic) > 0 {
		ibifcInputs := make([]cdbm.InfiniBandInterfaceCreateInput, 0, len(updatedInstances)*len(dbibic))
		for _, inst := range updatedInstances {
			for _, ibic := range dbibic {
				ibifcInputs = append(ibifcInputs, cdbm.InfiniBandInterfaceCreateInput{
					InstanceID:            inst.ID,
					SiteID:                site.ID,
					InfiniBandPartitionID: ibic.InfiniBandPartitionID,
					Device:                ibic.Device,
					DeviceInstance:        ibic.DeviceInstance,
					VirtualFunctionID:     ibic.VirtualFunctionID,
					IsPhysical:            ibic.IsPhysical,
					Vendor:                ibic.Vendor,
					Status:                cdbm.InfiniBandInterfaceStatusPending,
					CreatedBy:             dbUser.ID,
				})
			}
		}

		ibifcDAO := cdbm.NewInfiniBandInterfaceDAO(bcih.dbSession)
		createdIbIfcsAll, err = ibifcDAO.CreateMultiple(ctx, tx, ibifcInputs)
		if err != nil {
			logger.Error().Err(err).Msg("failed to batch create InfiniBand interfaces")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
				fmt.Sprintf("Failed to batch create InfiniBand interfaces: %v", err), nil)
		}
		logger.Info().Int("count", len(createdIbIfcsAll)).Msg("batch created all InfiniBand interfaces")
	}

	// --- Build and batch create NVLink Interfaces ---
	var createdNvlIfcsAll []cdbm.NVLinkInterface
	if len(dbnvlic) > 0 {
		nvlifcInputs := make([]cdbm.NVLinkInterfaceCreateInput, 0, len(updatedInstances)*len(dbnvlic))
		for _, inst := range updatedInstances {
			for _, nvlic := range dbnvlic {
				nvlifcInputs = append(nvlifcInputs, cdbm.NVLinkInterfaceCreateInput{
					InstanceID:               inst.ID,
					SiteID:                   site.ID,
					NVLinkLogicalPartitionID: nvlic.NVLinkLogicalPartitionID,
					Device:                   nvlic.Device,
					DeviceInstance:           nvlic.DeviceInstance,
					Status:                   cdbm.NVLinkInterfaceStatusPending,
					CreatedBy:                dbUser.ID,
				})
			}
		}

		nvlifcDAO := cdbm.NewNVLinkInterfaceDAO(bcih.dbSession)
		createdNvlIfcsAll, err = nvlifcDAO.CreateMultiple(ctx, tx, nvlifcInputs)
		if err != nil {
			logger.Error().Err(err).Msg("failed to batch create NVLink interfaces")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
				fmt.Sprintf("Failed to batch create NVLink interfaces: %v", err), nil)
		}
		logger.Info().Int("count", len(createdNvlIfcsAll)).Msg("batch created all NVLink interfaces")
	}

	// --- Build and batch create DPU Extension Service Deployments ---
	var createdDesdsAll []cdbm.DpuExtensionServiceDeployment
	if len(apiRequest.DpuExtensionServiceDeployments) > 0 {
		desdInputs := make([]cdbm.DpuExtensionServiceDeploymentCreateInput, 0, len(updatedInstances)*len(dpuServiceIDs))
		for _, inst := range updatedInstances {
			for j, desdID := range dpuServiceIDs {
				desdInputs = append(desdInputs, cdbm.DpuExtensionServiceDeploymentCreateInput{
					SiteID:                site.ID,
					TenantID:              tenant.ID,
					InstanceID:            inst.ID,
					DpuExtensionServiceID: desdID,
					Version:               apiRequest.DpuExtensionServiceDeployments[j].Version,
					Status:                cdbm.DpuExtensionServiceDeploymentStatusPending,
					CreatedBy:             dbUser.ID,
				})
			}
		}

		desdDAO := cdbm.NewDpuExtensionServiceDeploymentDAO(bcih.dbSession)
		createdDesdsAll, err = desdDAO.CreateMultiple(ctx, tx, desdInputs)
		if err != nil {
			logger.Error().Err(err).Msg("failed to batch create DPU Extension Service Deployments")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
				fmt.Sprintf("Failed to batch create DPU Extension Service Deployments: %v", err), nil)
		}
		logger.Info().Int("count", len(createdDesdsAll)).Msg("batch created all DPU Extension Service Deployments")
	}

	// --- Build and batch create Status Details ---
	sdInputs := make([]cdbm.StatusDetailCreateInput, 0, len(updatedInstances))
	for _, inst := range updatedInstances {
		sdInputs = append(sdInputs, cdbm.StatusDetailCreateInput{
			EntityID: inst.ID.String(),
			Status:   cdbm.InstanceStatusPending,
			Message:  cdb.GetStrPtr("received instance creation request, pending"),
		})
	}

	sdDAO := cdbm.NewStatusDetailDAO(bcih.dbSession)
	createdSdsAll, err := sdDAO.CreateMultiple(ctx, tx, sdInputs)
	if err != nil {
		logger.Error().Err(err).Msg("failed to batch create status details")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
			fmt.Sprintf("Failed to batch create status details: %v", err), nil)
	}
	logger.Info().Int("count", len(createdSdsAll)).Msg("batch created all status details")

	// --- Organize created records by instance and build workflow configs ---
	// Track all created data for building API responses and temporal workflow
	type instanceData struct {
		instance *cdbm.Instance
		ifcs     []cdbm.Interface
		ibifcs   []cdbm.InfiniBandInterface
		nvlifcs  []cdbm.NVLinkInterface
		desds    []cdbm.DpuExtensionServiceDeployment
		ssd      *cdbm.StatusDetail
		// Temporal workflow configs
		interfaceConfigs    []*cwssaws.InstanceInterfaceConfig
		ibInterfaceConfigs  []*cwssaws.InstanceIBInterfaceConfig
		nvlInterfaceConfigs []*cwssaws.InstanceNVLinkGpuConfig
		desdConfigs         []*cwssaws.InstanceDpuExtensionServiceConfig
	}

	createdInstancesData := make([]instanceData, len(updatedInstances))

	// Initialize data structures for each instance
	for i, inst := range updatedInstances {
		instCopy := inst // Make a copy to avoid loop variable capture
		createdInstancesData[i] = instanceData{
			instance:            &instCopy,
			ifcs:                make([]cdbm.Interface, 0, len(ifcResult.DBInterfaces)),
			ibifcs:              make([]cdbm.InfiniBandInterface, 0, len(dbibic)),
			nvlifcs:             make([]cdbm.NVLinkInterface, 0, len(dbnvlic)),
			desds:               make([]cdbm.DpuExtensionServiceDeployment, 0, len(dpuServiceIDs)),
			interfaceConfigs:    make([]*cwssaws.InstanceInterfaceConfig, 0, len(ifcResult.DBInterfaces)),
			ibInterfaceConfigs:  make([]*cwssaws.InstanceIBInterfaceConfig, 0, len(dbibic)),
			nvlInterfaceConfigs: make([]*cwssaws.InstanceNVLinkGpuConfig, 0, len(dbnvlic)),
			desdConfigs:         make([]*cwssaws.InstanceDpuExtensionServiceConfig, 0, len(dpuServiceIDs)),
		}
	}

	// Build instance ID to index map for efficient lookup
	instanceIDToIdx := make(map[uuid.UUID]int, len(updatedInstances))
	for i, inst := range updatedInstances {
		instanceIDToIdx[inst.ID] = i
	}

	// Distribute Interfaces and build workflow configs
	for _, ifc := range createdIfcsAll {
		idx := instanceIDToIdx[ifc.InstanceID]
		createdInstancesData[idx].ifcs = append(createdInstancesData[idx].ifcs, ifc)

		// Build temporal workflow config
		interfaceConfig := &cwssaws.InstanceInterfaceConfig{
			FunctionType: cwssaws.InterfaceFunctionType_VIRTUAL_FUNCTION,
		}
		if ifc.SubnetID != nil {
			interfaceConfig.NetworkSegmentId = &cwssaws.NetworkSegmentId{
				Value: ifcResult.SubnetIDMap[*ifc.SubnetID].ControllerNetworkSegmentID.String(),
			}
			interfaceConfig.NetworkDetails = &cwssaws.InstanceInterfaceConfig_SegmentId{
				SegmentId: &cwssaws.NetworkSegmentId{
					Value: ifcResult.SubnetIDMap[*ifc.SubnetID].ControllerNetworkSegmentID.String(),
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
		createdInstancesData[idx].interfaceConfigs = append(createdInstancesData[idx].interfaceConfigs, interfaceConfig)
	}

	// Distribute InfiniBand Interfaces and build workflow configs
	for _, ibifc := range createdIbIfcsAll {
		idx := instanceIDToIdx[ibifc.InstanceID]
		createdInstancesData[idx].ibifcs = append(createdInstancesData[idx].ibifcs, ibifc)

		// Build temporal workflow config
		ibInterfaceConfig := &cwssaws.InstanceIBInterfaceConfig{
			Device:         ibifc.Device,
			Vendor:         ibifc.Vendor,
			DeviceInstance: uint32(ibifc.DeviceInstance),
			FunctionType:   cwssaws.InterfaceFunctionType_PHYSICAL_FUNCTION,
			IbPartitionId:  &cwssaws.IBPartitionId{Value: ibifc.InfiniBandPartitionID.String()},
		}
		if !ibifc.IsPhysical {
			ibInterfaceConfig.FunctionType = cwssaws.InterfaceFunctionType_VIRTUAL_FUNCTION
			if ibifc.VirtualFunctionID != nil {
				vfID := uint32(*ibifc.VirtualFunctionID)
				ibInterfaceConfig.VirtualFunctionId = &vfID
			}
		}
		createdInstancesData[idx].ibInterfaceConfigs = append(createdInstancesData[idx].ibInterfaceConfigs, ibInterfaceConfig)
	}

	// Distribute NVLink Interfaces and build workflow configs
	for _, nvlifc := range createdNvlIfcsAll {
		idx := instanceIDToIdx[nvlifc.InstanceID]
		createdInstancesData[idx].nvlifcs = append(createdInstancesData[idx].nvlifcs, nvlifc)

		// Build temporal workflow config
		nvlInterfaceConfig := &cwssaws.InstanceNVLinkGpuConfig{
			DeviceInstance:     uint32(nvlifc.DeviceInstance),
			LogicalPartitionId: &cwssaws.NVLinkLogicalPartitionId{Value: nvlifc.NVLinkLogicalPartitionID.String()},
		}
		createdInstancesData[idx].nvlInterfaceConfigs = append(createdInstancesData[idx].nvlInterfaceConfigs, nvlInterfaceConfig)
	}

	// Distribute DPU Extension Service Deployments and build workflow configs
	for _, desd := range createdDesdsAll {
		idx := instanceIDToIdx[desd.InstanceID]
		createdInstancesData[idx].desds = append(createdInstancesData[idx].desds, desd)

		// Build temporal workflow config
		desdConfig := &cwssaws.InstanceDpuExtensionServiceConfig{
			ServiceId: desd.DpuExtensionServiceID.String(),
			Version:   desd.Version,
		}
		createdInstancesData[idx].desdConfigs = append(createdInstancesData[idx].desdConfigs, desdConfig)
	}

	// Distribute Status Details
	for i := range createdSdsAll {
		// Status details are created in the same order as instances
		createdInstancesData[i].ssd = &createdSdsAll[i]
	}

	// Log allocation constraint usage summary
	constraintsUsed := constraintIdx + 1
	logger.Info().
		Int("instanceCount", len(createdInstancesData)).
		Int("constraintsUsed", constraintsUsed).
		Int("totalConstraintsAvailable", len(availableConstraints)).
		Msg("all instance records created using batch operations, now triggering batch Temporal workflow before commit")

	// ==================== Step 7: Workflow Trigger ====================

	// Get Temporal client for the site
	stc, err := bcih.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build batch workflow request using pre-built configs (no DB queries)
	batchRequest := &cwssaws.BatchInstanceAllocationRequest{
		InstanceRequests: make([]*cwssaws.InstanceAllocationRequest, 0, len(createdInstancesData)),
	}

	for _, data := range createdInstancesData {
		instanceRequest := common.BuildInstanceAllocationRequest(
			data.instance, tenant, osConfig, instanceSshKeyGroupIds,
			data.interfaceConfigs, data.ibInterfaceConfigs,
			data.nvlInterfaceConfigs, data.desdConfigs, false,
		)
		batchRequest.InstanceRequests = append(batchRequest.InstanceRequests, instanceRequest)
	}

	logger.Info().Int("instanceCount", len(createdInstances)).
		Msg("triggering batch create Instances workflow")

	// Trigger batch workflow (use batchSuffix for consistency with instance names)
	workflowID := "instance-batch-create-" + batchSuffix
	workflowOptions := temporalClient.StartWorkflowOptions{
		ID: workflowID,
		// TODO: temporary config, to be tuned
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	// Add context timeout
	workflowCtx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	// Execute batch workflow with full request
	we, err := stc.ExecuteWorkflow(workflowCtx, workflowOptions, "CreateInstances", batchRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to synchronously start Temporal workflow to create batch Instances")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to start sync workflow to create batch Instances on Site: %s", err), nil)
	}

	wid := we.GetID()
	logger.Info().Str("Workflow ID", wid).Msg("executed synchronous batch create Instances workflow")

	// Wait for workflow to complete (synchronous, matching single API error handling)
	err = we.Get(workflowCtx, nil)
	if err != nil {
		// Handle timeout errors specially (matching single API pattern)
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || workflowCtx.Err() != nil {

			logger.Error().Err(err).Str("Workflow ID", wid).
				Msg("failed to create batch Instances, timeout occurred executing workflow on Site.")

			// Create a new context for termination
			newctx, newcancel := context.WithTimeout(context.Background(), common.WorkflowContextNewAfterTimeout)
			defer newcancel()

			// Initiate termination workflow
			serr := stc.TerminateWorkflow(newctx, wid, "", "timeout occurred executing batch create Instances workflow")
			if serr != nil {
				logger.Error().Err(serr).Str("Workflow ID", wid).Msg("failed to terminate Temporal workflow for creating batch Instances")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to terminate synchronous batch Instances creation workflow after timeout, Cloud and Site data may be de-synced: %s", serr), nil)
			}

			logger.Info().Str("Workflow ID", wid).Msg("initiated terminate synchronous batch create Instances workflow successfully")

			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError,
				fmt.Sprintf("Failed to create batch Instances, timeout occurred executing workflow on Site: %s", err), nil)
		}

		// Handle other workflow errors (matching single API pattern)
		code, err := common.UnwrapWorkflowError(err)
		logger.Error().Err(err).Str("Workflow ID", wid).Msg("failed to synchronously execute Temporal workflow to create batch Instances")
		return cerr.NewAPIErrorResponse(c, code,
			fmt.Sprintf("Failed to execute sync workflow to create batch Instances on Site: %s", err), nil)
	}

	logger.Info().Str("Workflow ID", wid).Int("instanceCount", len(createdInstances)).
		Msg("completed synchronous batch create Instances workflow")

	// ==================== Step 8: Commit & Response ====================

	// Commit transaction only after batch workflow succeeds
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing batch instance transaction to DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create batch Instances, DB transaction error", nil)
	}
	// Set committed so, deferred cleanup functions will do nothing
	txCommitted = true

	// Build complete API response with relations
	// All data was collected during the creation loop above
	apiInstances := []model.APIInstance{}
	for _, data := range createdInstancesData {
		// Build complete API instance using the same method as Single API
		sds := []cdbm.StatusDetail{}
		if data.ssd != nil {
			sds = append(sds, *data.ssd)
		}
		apiInstance := model.NewAPIInstance(data.instance, site, data.ifcs, data.ibifcs, data.desds, data.nvlifcs, sshKeyGroups, sds)
		apiInstances = append(apiInstances, *apiInstance)
	}

	logger.Info().Int("instancesCreated", len(createdInstancesData)).Msg("finishing API handler for batch instance creation")
	return c.JSON(http.StatusCreated, apiInstances)
}

// allocateMachinesForBatch allocates machines for batch instance creation with optional topology optimization.
// If topologyOptimized is true, all machines must be allocated on the same NVLink domain.
// If false, machines can be allocated across different NVLink domains without topology consideration.
//
// Returns:
//   - machines: the allocated machines
//   - error: API error if allocation fails
func allocateMachinesForBatch(
	ctx context.Context,
	tx *cdb.Tx,
	dbSession *cdb.Session,
	instancetype *cdbm.InstanceType,
	count int,
	topologyOptimized bool,
	logger zerolog.Logger,
) ([]cdbm.Machine, *cerr.APIError) {
	if instancetype == nil || count <= 0 {
		return nil, cerr.NewAPIError(http.StatusBadRequest, "Invalid parameters for machine allocation", nil)
	}

	if tx == nil {
		return nil, cerr.NewAPIError(http.StatusInternalServerError, "Transaction required for machine allocation", nil)
	}

	mcDAO := cdbm.NewMachineDAO(dbSession)

	// Get all available Machines for the Instance Type
	filterInput := cdbm.MachineFilterInput{
		SiteID:          instancetype.SiteID,
		InstanceTypeIDs: []uuid.UUID{instancetype.ID},
		IsAssigned:      cdb.GetBoolPtr(false),
		Statuses:        []string{cdbm.MachineStatusReady},
	}
	machines, _, err := mcDAO.GetAll(ctx, tx, filterInput, cdbp.PageInput{Limit: cdb.GetIntPtr(cdbp.TotalLimit)}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve available machines")
		return nil, cerr.NewAPIError(http.StatusInternalServerError, "Failed to retrieve available machines", nil)
	}

	if len(machines) < count {
		logger.Warn().Int("available", len(machines)).Int("requested", count).
			Msg("insufficient machines available")
		return nil, cerr.NewAPIError(http.StatusConflict,
			fmt.Sprintf("Insufficient machines available: requested %d, available %d", count, len(machines)), nil)
	}

	var candidateMachines []*cdbm.Machine

	// If topologyOptimized is false, no need to consider NVLink domain - just use all available machines
	if !topologyOptimized {
		logger.Info().Int("available", len(machines)).Int("requested", count).
			Msg("topology optimization disabled - allocating machines without NVLink domain constraint")
		candidateMachines = make([]*cdbm.Machine, len(machines))
		for i := range machines {
			candidateMachines[i] = &machines[i]
		}
	} else {
		// topologyOptimized is true - must allocate all machines from the same NVLink domain
		logger.Info().Int("available", len(machines)).Int("requested", count).
			Msg("topology optimization enabled - ensuring all machines from same NVLink domain")

		// Group machines by NVLink Domain ID from Metadata
		nvlinkDomainMap := make(map[string][]*cdbm.Machine)
		noNvlinkDomainMachines := []*cdbm.Machine{}

		for idx := range machines {
			machine := &machines[idx]
			domainID := getNVLinkDomainID(machine)
			if domainID != "" {
				nvlinkDomainMap[domainID] = append(nvlinkDomainMap[domainID], machine)
			} else {
				noNvlinkDomainMachines = append(noNvlinkDomainMachines, machine)
			}
		}

		// Find the NVLink domain with the most available machines
		var bestDomainID string
		for domainID, domainMachines := range nvlinkDomainMap {
			if len(domainMachines) > len(nvlinkDomainMap[bestDomainID]) {
				bestDomainID = domainID
			}
		}

		// Check if we have enough machines in a single NVLink domain
		if len(nvlinkDomainMap[bestDomainID]) < count {
			logger.Warn().Str("bestDomainID", bestDomainID).Int("bestDomainCount", len(nvlinkDomainMap[bestDomainID])).Int("requested", count).
				Msg("topology optimization requires same NVLink domain but insufficient machines in any single domain")
			return nil, cerr.NewAPIError(http.StatusConflict,
				fmt.Sprintf("Topology optimization requires all %d machines on same NVLink domain, but best domain only has %d available", count, len(nvlinkDomainMap[bestDomainID])), nil)
		}

		// Use machines from the best NVLink domain
		candidateMachines = nvlinkDomainMap[bestDomainID]
		logger.Info().Str("nvlinkDomainID", bestDomainID).Int("available", len(nvlinkDomainMap[bestDomainID])).Int("requested", count).
			Msg("allocating all machines from single NVLink domain")
	}

	// Randomize the list of machines to help distribute load and avoid bad machines
	rand.Shuffle(
		len(candidateMachines),
		func(i, j int) {
			candidateMachines[i], candidateMachines[j] = candidateMachines[j], candidateMachines[i]
		},
	)

	// Phase 1: Acquire locks and verify machines, collect update inputs
	updateInputs := make([]cdbm.MachineUpdateInput, 0, count)
	verifiedMachines := make([]*cdbm.Machine, 0, count)

	for _, mc := range candidateMachines {
		if len(verifiedMachines) >= count {
			break
		}

		// Acquire an advisory lock on the MachineID
		err = tx.TryAcquireAdvisoryLock(ctx, cdb.GetAdvisoryLockIDFromString(mc.ID), nil)
		if err != nil {
			continue
		}

		// Re-obtain the Machine record to ensure it is still available
		umc, err := mcDAO.GetByID(ctx, tx, mc.ID, nil, false)
		if err != nil {
			continue
		}

		if umc.Status != cdbm.MachineStatusReady {
			continue
		}

		if umc.IsAssigned {
			continue
		}

		// Machine is verified, add to batch update list
		updateInputs = append(updateInputs, cdbm.MachineUpdateInput{
			MachineID:  mc.ID,
			IsAssigned: cdb.GetBoolPtr(true),
		})
		verifiedMachines = append(verifiedMachines, umc)
	}

	if len(verifiedMachines) < count {
		logger.Error().Int("verified", len(verifiedMachines)).Int("requested", count).
			Msg("could not verify sufficient machines for allocation")
		return nil, cerr.NewAPIError(http.StatusConflict,
			fmt.Sprintf("Could not allocate sufficient machines: requested %d, verified %d", count, len(verifiedMachines)), nil)
	}

	// Phase 2: Batch update all verified machines to assigned
	allocatedMachines, err := mcDAO.UpdateMultiple(ctx, tx, updateInputs)
	if err != nil {
		logger.Error().Err(err).Msg("failed to batch update machines to assigned")
		return nil, cerr.NewAPIError(http.StatusInternalServerError,
			fmt.Sprintf("Failed to batch update machines: %v", err), nil)
	}

	// Log NVLink domain distribution for observability
	nvlinkDomainDistribution := make(map[string]int)
	for _, machine := range allocatedMachines {
		domainID := getNVLinkDomainID(&machine)
		nvlinkDomainDistribution[domainID]++
	}
	logger.Info().Interface("nvlinkDomainDistribution", nvlinkDomainDistribution).
		Bool("topologyOptimized", topologyOptimized).
		Int("nvlinkDomainCount", len(nvlinkDomainDistribution)).
		Int("machinesAllocated", len(allocatedMachines)).
		Msg("successfully allocated machines for batch creation")

	return allocatedMachines, nil
}

// getNVLinkDomainID extracts the NVLink domain ID from a machine's metadata.
// Returns empty string if the machine has no NVLink domain information.
func getNVLinkDomainID(machine *cdbm.Machine) string {
	if machine.Metadata != nil {
		if nvlinkInfo := machine.Metadata.GetNvlinkInfo(); nvlinkInfo != nil {
			if domainUuid := nvlinkInfo.GetDomainUuid(); domainUuid != nil {
				return domainUuid.GetValue()
			}
		}
	}
	return ""
}
