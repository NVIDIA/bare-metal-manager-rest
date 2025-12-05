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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nvidia/carbide-rest/api/internal/config"
	"github.com/nvidia/carbide-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
	"github.com/nvidia/carbide-rest/api/pkg/api/model/util"
	"github.com/nvidia/carbide-rest/api/pkg/api/pagination"
	sc "github.com/nvidia/carbide-rest/api/pkg/client/site"
	cerr "github.com/nvidia/carbide-rest/common/pkg/util"
	sutil "github.com/nvidia/carbide-rest/common/pkg/util"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
	"github.com/nvidia/carbide-rest/db/pkg/db/paginator"
	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"github.com/nvidia/carbide-rest/workflow/pkg/queue"
	"go.opentelemetry.io/otel/attribute"
	tclient "go.temporal.io/sdk/client"
)

// ~~~~~ Create Handler ~~~~~ //

// CreateExpectedMachineHandler is the API Handler for creating new ExpectedMachine
type CreateExpectedMachineHandler struct {
	dbSession  *cdb.Session
	tc         tclient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewCreateExpectedMachineHandler initializes and returns a new handler for creating ExpectedMachine
func NewCreateExpectedMachineHandler(dbSession *cdb.Session, tc tclient.Client, scp *sc.ClientPool, cfg *config.Config) CreateExpectedMachineHandler {
	return CreateExpectedMachineHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Create an ExpectedMachine
// @Description Create an ExpectedMachine
// @Tags ExpectedMachine
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param message body model.APIExpectedMachineCreateRequest true "ExpectedMachine creation request"
// @Success 201 {object} model.APIExpectedMachine
// @Router /v2/org/{org}/carbide/expected-machine [post]
func (cemh CreateExpectedMachineHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Create", "ExpectedMachine", c, cemh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}
	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// ensure our user is a provider or tenant for the org
	infrastructureProvider, tenant, apiError := common.IsProviderOrTenant(ctx, logger, cemh.dbSession, org, dbUser, false, true)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	// Validate request
	// Bind request data to API model
	apiRequest := model.APIExpectedMachineCreateRequest{}
	err := c.Bind(&apiRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data, potentially invalid structure", nil)
	}

	// Validate request attributes
	verr := apiRequest.Validate()
	if verr != nil {
		logger.Warn().Err(verr).Msg("error validating Expected Machine creation request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Error validating Expected Machine creation data", verr)
	}

	// Validate that SKU exists if specified
	if apiRequest.SkuID != nil {
		skuDAO := cdbm.NewSkuDAO(cemh.dbSession)
		_, err = skuDAO.Get(ctx, nil, *apiRequest.SkuID)
		if err != nil {
			if errors.Is(err, cdb.ErrDoesNotExist) {
				logger.Warn().Msg("SKU ID specified in request does not exist")
				return cerr.NewAPIErrorResponse(c, http.StatusUnprocessableEntity, "SKU ID specified in request does not exist", nil)
			}
			logger.Warn().Err(err).Msg("error validating SKU ID in request data")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Error validating SKU ID in request data, DB error", nil)
		}
	}

	// Retrieve the Site from the DB
	site, err := common.GetSiteFromIDString(ctx, nil, apiRequest.SiteID, cemh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request data does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request data, DB error", nil)
	}

	// Validate permissions based on user role
	if infrastructureProvider != nil {
		// Validate that site belongs to the organization's infrastructure provider
		if site.InfrastructureProviderID != infrastructureProvider.ID {
			logger.Warn().Msg("Site specified in request data does not belong to current org's Infrastructure Provider")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request data is does not belong to current org", nil)
		}
	} else if tenant != nil {
		// Check if tenant has an account with the Site's Infrastructure Provider
		taDAO := cdbm.NewTenantAccountDAO(cemh.dbSession)
		_, taCount, err := taDAO.GetAll(ctx, nil, cdbm.TenantAccountFilterInput{
			InfrastructureProviderID: &site.InfrastructureProviderID,
			TenantIDs:                []uuid.UUID{tenant.ID},
		}, paginator.PageInput{}, []string{})
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Tenant Account for Site")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant Account with Site's Provider, DB error", nil)
		}

		if taCount == 0 {
			logger.Error().Msg("Tenant doesn't have an account with Infrastructure Provider")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Tenant doesn't have an account with Provider of Site specified in request", nil)
		}
	}

	// Validate that site is in Registered state
	if site.Status != cdbm.SiteStatusRegistered {
		logger.Warn().Msg("Site specified in request data is not in Registered state")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request data is not in Registered state, cannot create Expected Machine", nil)
	}

	// Check for duplicate MAC address
	// Notes: We do not allow multiple Expected Machines with the same MAC address, but it's not a DB unique constraint so we check here
	emDAO := cdbm.NewExpectedMachineDAO(cemh.dbSession)
	ems, count, err := emDAO.GetAll(ctx, nil, cdbm.ExpectedMachineFilterInput{
		BmcMacAddresses: []string{apiRequest.BmcMacAddress},
		SiteIDs:         []uuid.UUID{site.ID},
	}, paginator.PageInput{
		Limit: cdb.GetIntPtr(1),
	}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error checking for duplicate MAC address on Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate MAC address uniqueness on Site, DB error", nil)
	}
	if count > 0 {
		logger.Warn().Str("mac_address", apiRequest.BmcMacAddress).Msg("Expected Machine with specified MAC address already exists on Site")
		return cerr.NewAPIErrorResponse(c, http.StatusConflict, "Expected Machine with specified MAC address already exists on Site", validation.Errors{
			"id": errors.New(ems[0].ID.String()),
		})
	}

	// Start a db transaction
	tx, err := cdb.BeginTx(ctx, cemh.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Expected Machine, DB transaction error", nil)
	}
	// this variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// Create the ExpectedMachine in DB
	// Note: DefaultBmcUsername and BmcPassword are not stored in DB, only passed to workflow
	expectedMachine, err := emDAO.Create(
		ctx,
		tx,
		cdbm.ExpectedMachineCreateInput{
			ExpectedMachineID:        uuid.New(),
			SiteID:                   site.ID,
			BmcMacAddress:            apiRequest.BmcMacAddress,
			ChassisSerialNumber:      apiRequest.ChassisSerialNumber,
			SkuID:                    apiRequest.SkuID,
			FallbackDpuSerialNumbers: apiRequest.FallbackDPUSerialNumbers,
			Labels:                   apiRequest.Labels,
			CreatedBy:                dbUser.ID,
		},
	)
	if err != nil {
		logger.Error().Err(err).Msg("error creating ExpectedMachine record in DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Expected Machine, DB error", nil)
	}

	// Build the create request for workflow
	// BMC credentials come from API request since they're not stored in DB
	createExpectedMachineRequest := &cwssaws.ExpectedMachine{
		Id:                       &cwssaws.UUID{Value: expectedMachine.ID.String()},
		BmcMacAddress:            expectedMachine.BmcMacAddress,
		ChassisSerialNumber:      expectedMachine.ChassisSerialNumber,
		FallbackDpuSerialNumbers: expectedMachine.FallbackDpuSerialNumbers,
		SkuId:                    expectedMachine.SkuID,
	}

	if apiRequest.DefaultBmcUsername != nil {
		createExpectedMachineRequest.BmcUsername = *apiRequest.DefaultBmcUsername
	}

	if apiRequest.DefaultBmcPassword != nil {
		createExpectedMachineRequest.BmcPassword = *apiRequest.DefaultBmcPassword
	}

	protoLabels := util.ProtobufLabelsFromAPILabels(apiRequest.Labels)
	if protoLabels != nil {
		createExpectedMachineRequest.Metadata = &cwssaws.Metadata{
			Labels: protoLabels,
		}
	}

	logger.Info().Msg("triggering Expected Machine create workflow on Site")

	// Create workflow options
	workflowOptions := tclient.StartWorkflowOptions{
		ID:                       "expected-machine-create-" + expectedMachine.ID.String(),
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	// Get the temporal client for the site we are working with
	stc, err := cemh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Run workflow
	apiErr := common.ExecuteSyncWorkflow(ctx, logger, stc, "CreateExpectedMachine", workflowOptions, createExpectedMachineRequest)
	if apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, apiErr.Data)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing ExpectedMachine transaction to DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Expected Machine, DB transaction error", nil)
	}
	// Set committed so, deferred cleanup functions will do nothing
	txCommitted = true

	// Create response
	apiExpectedMachine := model.NewAPIExpectedMachine(expectedMachine)

	logger.Info().Msg("finishing API handler")
	return c.JSON(http.StatusCreated, apiExpectedMachine)
}

// ~~~~~ GetAll Handler ~~~~~ //

// GetAllExpectedMachineHandler is the API Handler for getting all ExpectedMachines
type GetAllExpectedMachineHandler struct {
	dbSession  *cdb.Session
	tc         tclient.Client
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewGetAllExpectedMachineHandler initializes and returns a new handler for getting all ExpectedMachines
func NewGetAllExpectedMachineHandler(dbSession *cdb.Session, tc tclient.Client, cfg *config.Config) GetAllExpectedMachineHandler {
	return GetAllExpectedMachineHandler{
		dbSession:  dbSession,
		tc:         tc,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Get all ExpectedMachines
// @Description Get all ExpectedMachines
// @Tags ExpectedMachine
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string false "ID of Site (optional, filters results to specific site)"
// @Param pageNumber query integer false "Page number of results returned"
// @Param includeRelation query string false "Related entities to include in response e.g. 'Site', 'SKU'"
// @Param pageSize query integer false "Number of results per page"
// @Param orderBy query string false "Order by field"
// @Success 200 {object} []model.APIExpectedMachine
// @Router /v2/org/{org}/carbide/expected-machine [get]
func (gaemh GetAllExpectedMachineHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("GetAll", "ExpectedMachine", c, gaemh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}
	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// ensure our user is a provider for the org
	infrastructureProvider, tenant, apiError := common.IsProviderOrTenant(ctx, logger, gaemh.dbSession, org, dbUser, true, true)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	filterInput := cdbm.ExpectedMachineFilterInput{}

	// Get Site ID from query param if specified
	siteIDStr := c.QueryParam("siteId")
	if siteIDStr != "" {
		site, err := common.GetSiteFromIDString(ctx, nil, siteIDStr, gaemh.dbSession)
		if err != nil {
			if errors.Is(err, cdb.ErrDoesNotExist) {
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request data does not exist", nil)
			}
			logger.Error().Err(err).Msg("error retrieving Site from DB")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request data, DB error", nil)
		}

		// Validate permissions based on user role
		if infrastructureProvider != nil {
			// Validate that site belongs to the organization's infrastructure provider
			if site.InfrastructureProviderID != infrastructureProvider.ID {
				logger.Warn().Msg("Site specified in request data does not belong to current org's Infrastructure Provider")
				return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request data is does not belong to current org", nil)
			}
		} else if tenant != nil {
			// Check if tenant has an account with the Site's Infrastructure Provider
			taDAO := cdbm.NewTenantAccountDAO(gaemh.dbSession)
			_, taCount, err := taDAO.GetAll(ctx, nil, cdbm.TenantAccountFilterInput{
				InfrastructureProviderID: &site.InfrastructureProviderID,
				TenantIDs:                []uuid.UUID{tenant.ID},
			}, paginator.PageInput{}, []string{})
			if err != nil {
				logger.Error().Err(err).Msg("error retrieving Tenant Account for Site")
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant Account with Site's Provider, DB error", nil)
			}

			if taCount == 0 {
				logger.Error().Msg("Tenant doesn't have an account with Infrastructure Provider")
				return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Tenant doesn't have an account with Provider of Site specified in request", nil)
			}
		}

		filterInput.SiteIDs = []uuid.UUID{site.ID}
	} else if tenant != nil {
		// Tenants must specify a Site ID
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site ID must be specified in query when retrieving Expected Machines as a Tenant", nil)
	} else {
		// Get all Sites for the org's Infrastructure Provider
		siteDAO := cdbm.NewSiteDAO(gaemh.dbSession)
		sites, _, err := siteDAO.GetAll(ctx, nil,
			cdbm.SiteFilterInput{InfrastructureProviderID: &infrastructureProvider.ID},
			paginator.PageInput{Limit: cdb.GetIntPtr(math.MaxInt)},
			nil,
		)
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Sites from DB")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Sites for org, DB error", nil)
		}

		siteIDs := make([]uuid.UUID, 0, len(sites))
		for _, site := range sites {
			siteIDs = append(siteIDs, site.ID)
		}
		filterInput.SiteIDs = siteIDs
	}

	// Get and validate includeRelation params
	qParams := c.QueryParams()
	qIncludeRelations, errStr := common.GetAndValidateQueryRelations(qParams, cdbm.ExpectedMachineRelatedEntities)
	if errStr != "" {
		logger.Warn().Msg(errStr)
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, errStr, nil)
	}

	// Validate pagination request
	pageRequest := pagination.PageRequest{}
	err := c.Bind(&pageRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding pagination request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request pagination data", nil)
	}

	// Validate pagination attributes
	err = pageRequest.Validate(cdbm.ExpectedMachineOrderByFields)
	if err != nil {
		logger.Warn().Err(err).Msg("error validating pagination request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate pagination request data", err)
	}

	// Get Expected Machines from DB
	emDAO := cdbm.NewExpectedMachineDAO(gaemh.dbSession)
	expectedMachines, total, err := emDAO.GetAll(
		ctx,
		nil,
		filterInput,
		paginator.PageInput{
			Offset:  pageRequest.Offset,
			Limit:   pageRequest.Limit,
			OrderBy: pageRequest.OrderBy,
		}, qIncludeRelations,
	)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Expected Machines from db")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Expected Machines, DB error", nil)
	}

	// Create response
	apiExpectedMachines := []*model.APIExpectedMachine{}
	for _, em := range expectedMachines {
		apiExpectedMachine := model.NewAPIExpectedMachine(&em)
		apiExpectedMachines = append(apiExpectedMachines, apiExpectedMachine)
	}

	// Create pagination response header
	pageResponse := pagination.NewPageResponse(*pageRequest.PageNumber, *pageRequest.PageSize, total, pageRequest.OrderByStr)
	pageHeader, err := json.Marshal(pageResponse)
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling pagination response")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to generate pagination response header", nil)
	}

	c.Response().Header().Set(pagination.ResponseHeaderName, string(pageHeader))

	logger.Info().Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiExpectedMachines)
}

// ~~~~~ Get Handler ~~~~~ //

// GetExpectedMachineHandler is the API Handler for retrieving ExpectedMachine
type GetExpectedMachineHandler struct {
	dbSession  *cdb.Session
	tc         tclient.Client
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewGetExpectedMachineHandler initializes and returns a new handler to retrieve ExpectedMachine
func NewGetExpectedMachineHandler(dbSession *cdb.Session, tc tclient.Client, cfg *config.Config) GetExpectedMachineHandler {
	return GetExpectedMachineHandler{
		dbSession:  dbSession,
		tc:         tc,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Retrieve the ExpectedMachine
// @Description Retrieve the ExpectedMachine by ID
// @Tags ExpectedMachine
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Expected Machine"
// @Param includeRelation query string false "Related entities to include in response e.g. 'Site', 'SKU'"
// @Success 200 {object} model.APIExpectedMachine
// @Router /v2/org/{org}/carbide/expected-machine/{id} [get]
func (gemh GetExpectedMachineHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Get", "ExpectedMachine", c, gemh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}
	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// ensure our user is a provider for the org
	infrastructureProvider, tenant, apiError := common.IsProviderOrTenant(ctx, logger, gemh.dbSession, org, dbUser, true, true)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	// Get Expected Machine ID from URL param
	expectedMachineIDStr := c.Param("id")
	expectedMachineID, err := uuid.Parse(expectedMachineIDStr)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Expected Machine ID in URL", nil)
	}

	logger = logger.With().Str("Expected Machine ID", expectedMachineID.String()).Logger()

	gemh.tracerSpan.SetAttribute(handlerSpan, attribute.String("expected_machine_id", expectedMachineID.String()), logger)

	// Get and validate includeRelation params
	qParams := c.QueryParams()
	qIncludeRelations, errStr := common.GetAndValidateQueryRelations(qParams, cdbm.ExpectedMachineRelatedEntities)
	if errStr != "" {
		logger.Warn().Msg(errStr)
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, errStr, nil)
	}

	// Get ExpectedMachine from DB by ID
	emDAO := cdbm.NewExpectedMachineDAO(gemh.dbSession)
	expectedMachine, err := emDAO.Get(ctx, nil, expectedMachineID, qIncludeRelations, false)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Could not find Expected Machine with ID: %s", expectedMachineID.String()), nil)
		}
		logger.Error().Err(err).Msg("error retrieving Expected Machine from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Expected Machine, DB	error", nil)
	}

	// Get Site for the Expected Machine
	siteDAO := cdbm.NewSiteDAO(gemh.dbSession)
	site, err := siteDAO.GetByID(ctx, nil, expectedMachine.SiteID, nil, false)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site details for Expected Machine, DB error", nil)
	}

	// Validate permissions based on user role
	if infrastructureProvider != nil {
		// Validate that site belongs to the organization's infrastructure provider
		if site.InfrastructureProviderID != infrastructureProvider.ID {
			logger.Warn().Msg("Expected Machine does not belong to a Site owned by org's Infrastructure Provider")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Expected Machine does not belong to a Site owned by current org", nil)
		}
	} else if tenant != nil {
		// Check if tenant has an account with the Site's Infrastructure Provider
		taDAO := cdbm.NewTenantAccountDAO(gemh.dbSession)
		_, taCount, err := taDAO.GetAll(ctx, nil, cdbm.TenantAccountFilterInput{
			InfrastructureProviderID: &site.InfrastructureProviderID,
			TenantIDs:                []uuid.UUID{tenant.ID},
		}, paginator.PageInput{}, []string{})
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Tenant Account with Site's Infrastructure Provider")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant Account with Site's Provider, DB error", nil)
		}

		if taCount == 0 {
			logger.Error().Msg("Tenant doesn't have an account with Site's Provider")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Tenant doesn't have an account with Provider of Expected Machine's Site", nil)
		}
	}

	// Create response
	apiExpectedMachine := model.NewAPIExpectedMachine(expectedMachine)

	logger.Info().Msg("finishing API handler")
	return c.JSON(http.StatusOK, apiExpectedMachine)
}

// ~~~~~ Update Handler ~~~~~ //

// UpdateExpectedMachineHandler is the API Handler for updating a ExpectedMachine
type UpdateExpectedMachineHandler struct {
	dbSession  *cdb.Session
	tc         tclient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewUpdateExpectedMachineHandler initializes and returns a new handler for updating ExpectedMachine
func NewUpdateExpectedMachineHandler(dbSession *cdb.Session, tc tclient.Client, scp *sc.ClientPool, cfg *config.Config) UpdateExpectedMachineHandler {
	return UpdateExpectedMachineHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Update an existing ExpectedMachine
// @Description Update an existing ExpectedMachine by ID
// @Tags ExpectedMachine
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Expected Machine"
// @Param message body model.APIExpectedMachineUpdateRequest true "ExpectedMachine update request"
// @Success 200 {object} model.APIExpectedMachine
// @Router /v2/org/{org}/carbide/expected-machine/{id} [patch]
func (uemh UpdateExpectedMachineHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Update", "ExpectedMachine", c, uemh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Ensure our user is a provider for the org
	infrastructureProvider, tenant, apiError := common.IsProviderOrTenant(ctx, logger, uemh.dbSession, org, dbUser, false, true)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	// Get Expected Machine ID from URL param
	expectedMachineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Expected Machine ID in URL", nil)
	}
	logger = logger.With().Str("Expected Machine ID", expectedMachineID.String()).Logger()

	uemh.tracerSpan.SetAttribute(handlerSpan, attribute.String("expected_machine_id", expectedMachineID.String()), logger)

	// Validate request
	// Bind request data to API model
	apiRequest := model.APIExpectedMachineUpdateRequest{}
	err = c.Bind(&apiRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data, potentially invalid structure", nil)
	}
	// Validate request attributes
	verr := apiRequest.Validate()
	if verr != nil {
		logger.Warn().Err(verr).Msg("error validating ExpectedMachine update request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Error validating ExpectedMachine update request data", verr)
	}

	// Validate that SKU exists if specified
	if apiRequest.SkuID != nil {
		skuDAO := cdbm.NewSkuDAO(uemh.dbSession)
		_, err = skuDAO.Get(ctx, nil, *apiRequest.SkuID)
		if err != nil {
			if errors.Is(err, cdb.ErrDoesNotExist) {
				logger.Warn().Msg("SKU ID specified in request does not exist")
				return cerr.NewAPIErrorResponse(c, http.StatusUnprocessableEntity, "SKU ID specified in request does not exist", nil)
			}
			logger.Warn().Err(err).Msg("error validating SKU ID in request data")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Error validating SKU ID in request data, DB error", nil)
		}
	}

	// Get ExpectedMachine from DB by ID
	emDAO := cdbm.NewExpectedMachineDAO(uemh.dbSession)
	expectedMachine, err := emDAO.Get(ctx, nil, expectedMachineID, []string{cdbm.SiteRelationName}, false)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Could not find Expected Machine with ID: %s", expectedMachineID.String()), nil)
		}
		logger.Error().Err(err).Msg("error retrieving Expected Machine from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Expected Machine, DB	error", nil)
	}

	// Validate that Site relation exists for the Expected Machine
	site := expectedMachine.Site
	if site == nil {
		logger.Error().Msg("no Site relation found for Expected Machine")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site details for Expected Machine", nil)
	}

	// Validate permissions based on user role
	if infrastructureProvider != nil {
		// Validate that site belongs to the organization's infrastructure provider
		if site.InfrastructureProviderID != infrastructureProvider.ID {
			logger.Warn().Msg("Expected Machine does not belong to a Site owned by org's Infrastructure Provider")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Expected Machine does not belong to a Site owned by current org", nil)
		}
	} else if tenant != nil {
		// Check if tenant has an account with the Site's Infrastructure Provider
		taDAO := cdbm.NewTenantAccountDAO(uemh.dbSession)
		_, taCount, err := taDAO.GetAll(ctx, nil, cdbm.TenantAccountFilterInput{
			InfrastructureProviderID: &site.InfrastructureProviderID,
			TenantIDs:                []uuid.UUID{tenant.ID},
		}, paginator.PageInput{}, []string{})
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Tenant Account for Site")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant Account with Site's Provider, DB error", nil)
		}

		if taCount == 0 {
			logger.Error().Msg("Tenant doesn't have an account with Site's Provider")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Tenant doesn't have an account with Provider of Expected Machine's Site", nil)
		}
	}

	// Validate that site is in Registered state
	if site.Status != cdbm.SiteStatusRegistered {
		logger.Warn().Msg("Expected Machine's Site is not in Registered state")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Expected Machine's Site is not in Registered state, cannot execute update", nil)
	}

	// Start a db tx
	tx, err := cdb.BeginTx(ctx, uemh.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Expected Machine, DB transaction error", nil)
	}
	// this variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// Update ExpectedMachine in DB
	// Note: DefaultBmcUsername and BmcPassword are not stored in DB, only passed to workflow

	updatedExpectedMachine, err := emDAO.Update(
		ctx,
		tx,
		cdbm.ExpectedMachineUpdateInput{
			ExpectedMachineID:        expectedMachine.ID,
			BmcMacAddress:            apiRequest.BmcMacAddress,
			ChassisSerialNumber:      apiRequest.ChassisSerialNumber,
			SkuID:                    apiRequest.SkuID,
			FallbackDpuSerialNumbers: apiRequest.FallbackDPUSerialNumbers,
			Labels:                   apiRequest.Labels,
		},
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update ExpectedMachine record in DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Expected Machine, DB error", nil)
	}

	// Build the update request for workflow
	// BMC credentials come from API request since they're not stored in DB
	updateExpectedMachineRequest := &cwssaws.ExpectedMachine{
		Id:                       &cwssaws.UUID{Value: expectedMachine.ID.String()},
		BmcMacAddress:            updatedExpectedMachine.BmcMacAddress,
		ChassisSerialNumber:      updatedExpectedMachine.ChassisSerialNumber,
		FallbackDpuSerialNumbers: updatedExpectedMachine.FallbackDpuSerialNumbers,
		SkuId:                    updatedExpectedMachine.SkuID,
	}

	if apiRequest.DefaultBmcUsername != nil {
		updateExpectedMachineRequest.BmcUsername = *apiRequest.DefaultBmcUsername
	}

	if apiRequest.DefaultBmcPassword != nil {
		updateExpectedMachineRequest.BmcPassword = *apiRequest.DefaultBmcPassword
	}

	protoLabels := util.ProtobufLabelsFromAPILabels(apiRequest.Labels)
	if protoLabels != nil {
		updateExpectedMachineRequest.Metadata = &cwssaws.Metadata{
			Labels: protoLabels,
		}
	}

	logger.Info().Msg("triggering ExpectedMachine update workflow")

	// Create workflow options
	workflowOptions := tclient.StartWorkflowOptions{
		ID:                       "expected-machine-update-" + expectedMachine.ID.String(),
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	// Get the Temporal client for the site we are working with
	stc, err := uemh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Run workflow
	apiErr := common.ExecuteSyncWorkflow(ctx, logger, stc, "UpdateExpectedMachine", workflowOptions, updateExpectedMachineRequest)
	if apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, apiErr.Data)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing ExpectedMachine update transaction to DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update ExpectedMachine", nil)
	}
	// Set committed so, deferred cleanup functions will do nothing
	txCommitted = true

	// Create response
	apiExpectedMachine := model.NewAPIExpectedMachine(updatedExpectedMachine)

	logger.Info().Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiExpectedMachine)
}

// ~~~~~ Delete Handler ~~~~~ //

// DeleteExpectedMachineHandler is the API Handler for deleting a ExpectedMachine
type DeleteExpectedMachineHandler struct {
	dbSession  *cdb.Session
	tc         tclient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewDeleteExpectedMachineHandler initializes and returns a new handler for deleting ExpectedMachine
func NewDeleteExpectedMachineHandler(dbSession *cdb.Session, tc tclient.Client, scp *sc.ClientPool, cfg *config.Config) DeleteExpectedMachineHandler {
	return DeleteExpectedMachineHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Delete an existing ExpectedMachine
// @Description Delete an existing ExpectedMachine by ID
// @Tags ExpectedMachine
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Expected Machine"
// @Success 204
// @Router /v2/org/{org}/carbide/expected-machine/{id} [delete]
func (demh DeleteExpectedMachineHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Delete", "ExpectedMachine", c, demh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}
	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Ensure our user is a provider for the org
	infrastructureProvider, tenant, apiError := common.IsProviderOrTenant(ctx, logger, demh.dbSession, org, dbUser, false, true)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	// Get Expected Machine ID from URL param
	expectedMachineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Expected Machine ID in URL", nil)
	}
	logger = logger.With().Str("Expected Machine ID", expectedMachineID.String()).Logger()

	demh.tracerSpan.SetAttribute(handlerSpan, attribute.String("expected_machine_id", expectedMachineID.String()), logger)

	// Get ExpectedMachine from DB by ID
	emDAO := cdbm.NewExpectedMachineDAO(demh.dbSession)
	expectedMachine, err := emDAO.Get(ctx, nil, expectedMachineID, []string{cdbm.SiteRelationName}, false)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Could not find Expected Machine with ID: %s", expectedMachineID.String()), nil)
		}
		logger.Error().Err(err).Msg("error retrieving Expected Machine from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Expected Machine, DB	error", nil)
	}

	// Validate that Site relation exists for the Expected Machine
	site := expectedMachine.Site
	if site == nil {
		logger.Error().Msg("no Site relation found for Expected Machine")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site details for Expected Machine", nil)
	}

	// Validate permissions based on user role
	if infrastructureProvider != nil {
		// Validate that site belongs to the organization's infrastructure provider
		if site.InfrastructureProviderID != infrastructureProvider.ID {
			logger.Warn().Msg("Expected Machine does not belong to a Site owned by org's Infrastructure Provider")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Expected Machine does not belong to a Site owned by current org", nil)
		}
	} else if tenant != nil {
		// Check if tenant has an account with the Site's Infrastructure Provider
		taDAO := cdbm.NewTenantAccountDAO(demh.dbSession)
		_, taCount, err := taDAO.GetAll(ctx, nil, cdbm.TenantAccountFilterInput{
			InfrastructureProviderID: &site.InfrastructureProviderID,
			TenantIDs:                []uuid.UUID{tenant.ID},
		}, paginator.PageInput{}, []string{})
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Tenant Account for Site")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant Account with Site's Provider, DB error", nil)
		}

		if taCount == 0 {
			logger.Error().Msg("Tenant doesn't have an account with Site's Provider")
			return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Tenant doesn't have an account with Provider of Expected Machine's Site", nil)
		}
	}

	// Validate that site is in Registered state
	if site.Status != cdbm.SiteStatusRegistered {
		logger.Warn().Msg("Expected Machine's Site is not in Registered state")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Expected Machine's Site is not in Registered state, cannot execute delete", nil)
	}

	// Start a db tx
	tx, err := cdb.BeginTx(ctx, demh.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to delete Expected Machine, DB error", nil)
	}
	// this variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// Delete ExpectedMachine from DB
	err = emDAO.Delete(ctx, tx, expectedMachine.ID)
	if err != nil {
		logger.Error().Err(err).Msg("unable to delete ExpectedMachine record from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to delete Expected Machine, DB error", nil)
	}

	// Build the delete request for workflow
	deleteExpectedMachineRequest := &cwssaws.ExpectedMachineRequest{
		Id: &cwssaws.UUID{Value: expectedMachine.ID.String()},
	}

	logger.Info().Msg("triggering ExpectedMachine delete workflow")

	// Create workflow options
	workflowOptions := tclient.StartWorkflowOptions{
		ID:                       "expected-machine-delete-" + expectedMachine.ID.String(),
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	// Get the temporal client for the site we are working with
	stc, err := demh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Run workflow
	apiErr := common.ExecuteSyncWorkflow(ctx, logger, stc, "DeleteExpectedMachine", workflowOptions, deleteExpectedMachineRequest)
	if apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, apiErr.Data)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing ExpectedMachine delete transaction to DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to delete Expected Machine, DB transaction error", nil)
	}
	// Set committed so, deferred cleanup functions will do nothing
	txCommitted = true

	logger.Info().Msg("finishing API handler")

	return c.NoContent(http.StatusNoContent)
}
