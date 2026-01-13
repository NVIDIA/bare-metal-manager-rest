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
	"maps"
	"math"
	"net/http"
	"slices"
	"strconv"

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
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/attribute"
	tclient "go.temporal.io/sdk/client"
)

func validateProviderTenantSiteAccess(ctx context.Context, logger zerolog.Logger, dbSession *cdb.Session, infrastructureProvider *cdbm.InfrastructureProvider, tenant *cdbm.Tenant, site *cdbm.Site) (code int, message string) {
	// Validate permissions based on user role
	if infrastructureProvider != nil {
		// Validate that site belongs to the organization's infrastructure provider
		if site.InfrastructureProviderID != infrastructureProvider.ID {
			logger.Warn().Msg("Site is not owned by org's Infrastructure Provider")
			return http.StatusForbidden, "Site is not owned by org's Infrastructure Provider"
		}
	} else if tenant != nil {
		// Check if tenant has an account with the Site's Infrastructure Provider
		taDAO := cdbm.NewTenantAccountDAO(dbSession)
		_, taCount, err := taDAO.GetAll(ctx, nil, cdbm.TenantAccountFilterInput{
			InfrastructureProviderID: &site.InfrastructureProviderID,
			TenantIDs:                []uuid.UUID{tenant.ID},
		}, paginator.PageInput{}, []string{})
		if err != nil {
			logger.Error().Err(err).Msg("error retrieving Tenant Account for Site")
			return http.StatusInternalServerError, "Failed to retrieve Tenant Account with Site's Provider due to DB error"
		}

		if taCount == 0 {
			logger.Error().Msg("Tenant doesn't have an account with Site's Provider")
			return http.StatusForbidden, "Tenant doesn't have an account with Site's Provider"
		}
	}

	// Validate that site is in Registered state
	if site.Status != cdbm.SiteStatusRegistered {
		logger.Warn().Msg("Site is not in Registered state")
		return http.StatusBadRequest, "Site is not in Registered state, cannot perform operation"
	}

	return http.StatusOK, ""
}

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
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate Expected Machine creation data", verr)
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
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate SKU ID in request data due to DB error", nil)
		}
	}

	// Retrieve the Site from the DB
	site, err := common.GetSiteFromIDString(ctx, nil, apiRequest.SiteID, cemh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request data does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request data due to DB error", nil)
	}

	// Validate ProviderTenantSite relationship and site state
	code, message := validateProviderTenantSiteAccess(ctx, logger, cemh.dbSession, infrastructureProvider, tenant, site)
	if code != http.StatusOK {
		return cerr.NewAPIErrorResponse(c, code, message, nil)
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate MAC address uniqueness on Site due to DB error", nil)
	}
	if count > 0 {
		logger.Warn().Str("MacAddress", apiRequest.BmcMacAddress).Msg("Expected Machine with specified MAC address already exists on Site")
		return cerr.NewAPIErrorResponse(c, http.StatusConflict, "Expected Machine with specified MAC address already exists on Site", validation.Errors{
			"id": errors.New(ems[0].ID.String()),
		})
	}

	// Start a db transaction
	tx, err := cdb.BeginTx(ctx, cemh.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Expected Machine due to DB transaction error", nil)
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Expected Machine due to DB error", nil)
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Expected Machine due to DB transaction error", nil)
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
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request data due to DB error", nil)
		}

		// Validate permissions based on user role
		if infrastructureProvider != nil {
			// Validate that site belongs to the organization's infrastructure provider
			if site.InfrastructureProviderID != infrastructureProvider.ID {
				logger.Warn().Msg("Site specified in request data does not belong to current org's Infrastructure Provider")
				return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request data does not belong to current org", nil)
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
				return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant Account with Site's Provider due to DB error", nil)
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
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Sites for org due to DB error", nil)
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Expected Machines due to DB error", nil)
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

	logger = logger.With().Str("ExpectedMachineID", expectedMachineID.String()).Logger()

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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Expected Machine due to DB error", nil)
	}

	// Get Site for the Expected Machine
	siteDAO := cdbm.NewSiteDAO(gemh.dbSession)
	site, err := siteDAO.GetByID(ctx, nil, expectedMachine.SiteID, nil, false)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site details for Expected Machine due to DB error", nil)
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
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant Account with Site's Provider due to DB error", nil)
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
	logger = logger.With().Str("ExpectedMachineID", expectedMachineID.String()).Logger()

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
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate ExpectedMachine update request data", verr)
	}

	// If ID is provided in body, it must match the path ID
	if apiRequest.ID != nil && *apiRequest.ID != expectedMachineID.String() {
		logger.Warn().
			Str("URLID", expectedMachineID.String()).
			Str("RequestDataID", *apiRequest.ID).
			Msg("Mismatched Expected Machine ID between path and body")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "If provided, Expected Machine ID specified in request data must match URL request value", nil)
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
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate SKU ID in request data due to DB error", nil)
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Expected Machine due to DB error", nil)
	}

	// Validate that Site relation exists for the Expected Machine
	site := expectedMachine.Site
	if site == nil {
		logger.Error().Msg("no Site relation found for Expected Machine")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site details for Expected Machine", nil)
	}

	// Validate ProviderTenantSite relationship and site state
	code, message := validateProviderTenantSiteAccess(ctx, logger, uemh.dbSession, infrastructureProvider, tenant, site)
	if code != http.StatusOK {
		return cerr.NewAPIErrorResponse(c, code, message, nil)
	}

	// Start a db tx
	tx, err := cdb.BeginTx(ctx, uemh.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Expected Machine due to DB transaction error", nil)
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Expected Machine due to DB error", nil)
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

// validateSkuIDs queries all provided SKU IDs once and returns the indices (relative
// to the provided skuIDs slice) that are missing in the database.
// Note: we expect empty IDs to ensure that the input array indices match the original Expected Machine array length.
func validateSkuIDs(ctx context.Context, tx *cdb.Tx, skuDAO cdbm.SkuDAO, siteID uuid.UUID, skuIDs []string) ([]int, error) {
	if len(skuIDs) == 0 {
		return nil, nil
	}

	// Collect unique non-empty SKU IDs for DB request
	uniqueSkuIDs := make(map[string]struct{})
	for _, skuID := range skuIDs {
		if skuID != "" {
			uniqueSkuIDs[skuID] = struct{}{}
		}
	}

	skus, _, err := skuDAO.GetAll(ctx, tx, cdbm.SkuFilterInput{
		SkuIDs:  slices.Collect(maps.Keys(uniqueSkuIDs)),
		SiteIDs: []uuid.UUID{siteID},
	}, paginator.PageInput{
		Limit: cdb.GetIntPtr(len(uniqueSkuIDs)),
	})
	if err != nil {
		return nil, err
	}

	// build a set for skuID found in DB
	foundSkuIDs := make(map[string]bool)
	for _, sku := range skus {
		foundSkuIDs[sku.ID] = true
	}

	// identify which input skuID are missing in DB results so that we can report to caller
	missing := make([]int, 0)
	for idx, skuID := range skuIDs {
		if skuID == "" {
			continue
		}
		if !foundSkuIDs[skuID] {
			missing = append(missing, idx)
		}
	}

	return missing, nil
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
	logger = logger.With().Str("ExpectedMachineID", expectedMachineID.String()).Logger()

	demh.tracerSpan.SetAttribute(handlerSpan, attribute.String("expected_machine_id", expectedMachineID.String()), logger)

	// Get ExpectedMachine from DB by ID
	emDAO := cdbm.NewExpectedMachineDAO(demh.dbSession)
	expectedMachine, err := emDAO.Get(ctx, nil, expectedMachineID, []string{cdbm.SiteRelationName}, false)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Could not find Expected Machine with ID: %s", expectedMachineID.String()), nil)
		}
		logger.Error().Err(err).Msg("error retrieving Expected Machine from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Expected Machine due to DB error", nil)
	}

	// Validate that Site relation exists for the Expected Machine
	site := expectedMachine.Site
	if site == nil {
		logger.Error().Msg("no Site relation found for Expected Machine")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site details for Expected Machine", nil)
	}

	// Validate ProviderTenantSite relationship and site state
	code, message := validateProviderTenantSiteAccess(ctx, logger, demh.dbSession, infrastructureProvider, tenant, site)
	if code != http.StatusOK {
		return cerr.NewAPIErrorResponse(c, code, message, nil)
	}

	// Start a db tx
	tx, err := cdb.BeginTx(ctx, demh.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to delete Expected Machine due to DB error", nil)
	}
	// this variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// Delete ExpectedMachine from DB
	err = emDAO.Delete(ctx, tx, expectedMachine.ID)
	if err != nil {
		logger.Error().Err(err).Msg("unable to delete ExpectedMachine record from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to delete Expected Machine due to DB error", nil)
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to delete Expected Machine due to DB transaction error", nil)
	}
	// Set committed so, deferred cleanup functions will do nothing
	txCommitted = true

	logger.Info().Msg("finishing API handler")

	return c.NoContent(http.StatusNoContent)
}

// ~~~~~ CreateExpectedMachines Handler ~~~~~ //

// CreateExpectedMachinesHandler is the API Handler for creating multiple ExpectedMachines
type CreateExpectedMachinesHandler struct {
	dbSession  *cdb.Session
	tc         tclient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewCreateExpectedMachinesHandler initializes and returns a new handler for creating multiple ExpectedMachines
func NewCreateExpectedMachinesHandler(dbSession *cdb.Session, tc tclient.Client, scp *sc.ClientPool, cfg *config.Config) CreateExpectedMachinesHandler {
	return CreateExpectedMachinesHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Create multiple ExpectedMachines
// @Description Create multiple ExpectedMachines in a single request. All machines must belong to the same site.
// @Tags ExpectedMachine
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param message body []model.APIExpectedMachineCreateRequest true "ExpectedMachine batch creation request"
// @Success 201 {object} model.APIExpectedMachineBatchResponse
// @Router /v2/org/{org}/forge/expected-machine/batch [post]
func (cemh CreateExpectedMachinesHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("CreateMultiple", "ExpectedMachine", c, cemh.tracerSpan)
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
	// Bind request data to API model (array payload)
	apiRequests := []model.APIExpectedMachineCreateRequest{}
	err := c.Bind(&apiRequests)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data, potentially invalid structure", nil)
	}

	if len(apiRequests) == 0 {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Request data must contain at least 1 Expected Machine entry", nil)
	}

	if len(apiRequests) > model.ExpectedMachineMaxBatchItems {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("At most %d Expected Machine entries can be created in a batch request", model.ExpectedMachineMaxBatchItems), nil)
	}

	// Validate each item and ensure all site IDs match
	siteID := apiRequests[0].SiteID // we use the first item's Site ID as reference
	bmcMacMap := make(map[string]int)
	serialMap := make(map[string]int)
	for i := range apiRequests {
		if verr := apiRequests[i].Validate(); verr != nil {
			logger.Warn().Err(verr).Int("Index", i).Msg("error validating Expected Machine creation request data")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate Expected Machine creation data", validation.Errors{
				strconv.Itoa(i): verr,
			})
		}

		if apiRequests[i].SiteID != siteID {
			logger.Warn().Str("SiteID", siteID).Str("NewSiteID", apiRequests[i].SiteID).Msg("Mismatch between Expected Machines Site ID in request")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "All Expected Machines in request must belong to the same site", nil)
		}

		if prev, ok := bmcMacMap[apiRequests[i].BmcMacAddress]; ok {
			logger.Warn().Msgf("duplicate BMC MAC address '%s' found at indices %d and %d", apiRequests[i].BmcMacAddress, prev, i)
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Duplicate BMC MAC address '%s' found at indices %d and %d", apiRequests[i].BmcMacAddress, prev, i), nil)
		}
		bmcMacMap[apiRequests[i].BmcMacAddress] = i

		if prev, ok := serialMap[apiRequests[i].ChassisSerialNumber]; ok {
			logger.Warn().Msgf("duplicate chassis serial number '%s' found at indices %d and %d", apiRequests[i].ChassisSerialNumber, prev, i)
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Duplicate chassis serial number '%s' found at indices %d and %d", apiRequests[i].ChassisSerialNumber, prev, i), nil)
		}
		serialMap[apiRequests[i].ChassisSerialNumber] = i
	}

	logger.Info().
		Int("MachineCount", len(apiRequests)).
		Msg("processing CreateExpectedMachines request")

	// Retrieve the Site from the DB
	site, err := common.GetSiteFromIDString(ctx, nil, siteID, cemh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request data does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request data due to DB error", nil)
	}

	// Validate ProviderTenantSite relationship and site state
	code, message := validateProviderTenantSiteAccess(ctx, logger, cemh.dbSession, infrastructureProvider, tenant, site)
	if code != http.StatusOK {
		return cerr.NewAPIErrorResponse(c, code, message, nil)
	}

	// Validate that all specified SKU IDs exist
	skuDAO := cdbm.NewSkuDAO(cemh.dbSession)

	skuIDs := make([]string, 0, len(apiRequests))
	for _, machine := range apiRequests {
		if machine.SkuID != nil {
			skuIDs = append(skuIDs, *machine.SkuID)
		} else {
			skuIDs = append(skuIDs, "")
		}
	}

	missingIndices, verr := validateSkuIDs(ctx, nil, skuDAO, site.ID, skuIDs)
	if verr != nil {
		logger.Warn().Err(verr).Msg("error validating SKU IDs in request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate SKU ID in request data due to DB error", nil)
	}

	if len(missingIndices) > 0 {
		indices, _ := json.Marshal(missingIndices)
		logger.Warn().Msg(fmt.Sprintf("SKU ID specified at indices %s do not exist", indices))
		return cerr.NewAPIErrorResponse(c, http.StatusUnprocessableEntity, fmt.Sprintf("SKU ID specified at indices %s do not exist", indices), nil)
	}

	// Check for duplicate Chassis Serial Number on the Site
	// Note: at this stage Chassis Serial Number is guaranteed to be non-empty due to prior validation.
	emDAO := cdbm.NewExpectedMachineDAO(cemh.dbSession)
	chassisSerialNumbers := make([]string, 0, len(apiRequests))
	for _, machine := range apiRequests {
		chassisSerialNumbers = append(chassisSerialNumbers, machine.ChassisSerialNumber)
	}

	existingMachines, count, err := emDAO.GetAll(ctx, nil, cdbm.ExpectedMachineFilterInput{
		ChassisSerialNumbers: chassisSerialNumbers,
		SiteIDs:              []uuid.UUID{site.ID},
	}, paginator.PageInput{
		Limit: cdb.GetIntPtr(len(chassisSerialNumbers)),
	}, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error checking for duplicate Chassis Serial Number on Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate Chassis Serial Number uniqueness on Site due to DB error", nil)
	}

	if count > 0 {
		// Build list of conflicting Chassis Serial Number
		conflictingSerials := make([]string, 0, len(existingMachines))
		for _, em := range existingMachines {
			conflictingSerials = append(conflictingSerials, em.ChassisSerialNumber)
		}
		logger.Warn().Strs("ChassisSerialNumber", conflictingSerials).Msg("Expected Machines with specified Chassis Serial Number already exist on Site")
		return cerr.NewAPIErrorResponse(c, http.StatusConflict, fmt.Sprintf("Expected Machines with Chassis Serial Number %v already exist on Site", conflictingSerials), nil)
	}

	// Start a db transaction
	tx, err := cdb.BeginTx(ctx, cemh.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Expected Machines due to DB transaction error", nil)
	}
	// this variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	createInputs := make([]cdbm.ExpectedMachineCreateInput, 0, len(apiRequests))
	for _, machineReq := range apiRequests {
		createInputs = append(createInputs, cdbm.ExpectedMachineCreateInput{
			ExpectedMachineID:        uuid.New(),
			SiteID:                   site.ID,
			BmcMacAddress:            machineReq.BmcMacAddress,
			ChassisSerialNumber:      machineReq.ChassisSerialNumber,
			SkuID:                    machineReq.SkuID,
			FallbackDpuSerialNumbers: machineReq.FallbackDPUSerialNumbers,
			Labels:                   machineReq.Labels,
			CreatedBy:                dbUser.ID,
		})
	}

	createdExpectedMachines, err := emDAO.CreateMultiple(ctx, tx, createInputs)
	if err != nil {
		logger.Error().Err(err).Msg("error creating ExpectedMachine records in DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Expected Machine due to DB error", nil)
	}

	workflowMachines := make([]*cwssaws.ExpectedMachine, 0, len(createdExpectedMachines))
	for i, createdMachine := range createdExpectedMachines {
		workflowMachine := &cwssaws.ExpectedMachine{
			Id:                       &cwssaws.UUID{Value: createdMachine.ID.String()},
			BmcMacAddress:            createdMachine.BmcMacAddress,
			ChassisSerialNumber:      createdMachine.ChassisSerialNumber,
			FallbackDpuSerialNumbers: createdMachine.FallbackDpuSerialNumbers,
			SkuId:                    createdMachine.SkuID,
		}

		if apiRequests[i].DefaultBmcUsername != nil {
			workflowMachine.BmcUsername = *apiRequests[i].DefaultBmcUsername
		}

		if apiRequests[i].DefaultBmcPassword != nil {
			workflowMachine.BmcPassword = *apiRequests[i].DefaultBmcPassword
		}

		protoLabels := util.ProtobufLabelsFromAPILabels(apiRequests[i].Labels)
		if protoLabels != nil {
			workflowMachine.Metadata = &cwssaws.Metadata{
				Labels: protoLabels,
			}
		}

		workflowMachines = append(workflowMachines, workflowMachine)
	}

	logger.Info().Int("Count", len(workflowMachines)).Msg("triggering CreateExpectedMachines workflow on Site")

	// Create workflow request
	workflowRequest := &cwssaws.BatchExpectedMachineOperationRequest{
		ExpectedMachines:     &cwssaws.ExpectedMachineList{ExpectedMachines: workflowMachines},
		AcceptPartialResults: false,
	}

	// Create workflow options
	workflowID := fmt.Sprintf("create-expected-machines-%s-%d", site.ID.String(), len(workflowMachines))
	workflowOptions := tclient.StartWorkflowOptions{
		ID:                       workflowID,
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	// Get the temporal client for the site we are working with
	stc, err := cemh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Execute workflow and get results
	workflowRun, err := stc.ExecuteWorkflow(ctx, workflowOptions, "CreateExpectedMachines", workflowRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to schedule CreateExpectedMachines workflow on Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to schedule batch Expected Machine creation workflow on Site: %v", err), nil)
	}

	workflowRunID := workflowRun.GetID()
	logger = logger.With().Str("WorkflowID", workflowRunID).Logger()
	logger.Info().Msg("executing CreateExpectedMachines workflow on Site")

	// Get workflow results
	var workflowResult cwssaws.BatchExpectedMachineOperationResponse

	err = workflowRun.Get(ctx, &workflowResult)
	if err != nil {
		logger.Error().Err(err).Msg("error executing CreateExpectedMachines workflow on Site")
		// Workflow failed entirely - don't commit transaction
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to execute batch Expected Machine creation workflow on Site: %v", err), nil)
	}

	// sanity checks since this is all-or-nothing
	if len(workflowResult.GetResults()) != len(createdExpectedMachines) {
		logger.Warn().Msgf("Workflow returned a different number of Expected Machines (expected %d but got %d)", len(createdExpectedMachines), len(workflowResult.GetResults()))
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing ExpectedMachine CreateExpectedMachines transaction to DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create Expected Machines due to DB transaction error", nil)
	}
	txCommitted = true

	logger.Info().
		Int("SuccessCount", len(createdExpectedMachines)).
		Msg("finishing CreateExpectedMachines API handler")

	// Return only successful machines
	return c.JSON(http.StatusCreated, createdExpectedMachines)
}

// ~~~~~ Batch Update Handler ~~~~~ //

// UpdateExpectedMachinesHandler is the API Handler for batch updating ExpectedMachines
type UpdateExpectedMachinesHandler struct {
	dbSession  *cdb.Session
	tc         tclient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewUpdateExpectedMachinesHandler initializes and returns a new handler for batch updating ExpectedMachines
func NewUpdateExpectedMachinesHandler(dbSession *cdb.Session, tc tclient.Client, scp *sc.ClientPool, cfg *config.Config) UpdateExpectedMachinesHandler {
	return UpdateExpectedMachinesHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Batch update ExpectedMachines
// @Description Update multiple ExpectedMachines in a single request. All machines must belong to the same site.
// @Tags ExpectedMachine
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param message body []model.APIExpectedMachineUpdateRequest true "ExpectedMachine UpdateExpectedMachines request"
// @Success 200 {object} model.APIExpectedMachineBatchResponse
// @Router /v2/org/{org}/forge/expected-machine/batch [patch]
func (uemh UpdateExpectedMachinesHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("UpdateMultiple", "ExpectedMachine", c, uemh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Ensure our user is a provider or tenant for the org
	infrastructureProvider, tenant, apiError := common.IsProviderOrTenant(ctx, logger, uemh.dbSession, org, dbUser, false, true)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	// Validate request
	// Bind request data to API model (array payload)
	apiRequests := []model.APIExpectedMachineUpdateRequest{}
	err := c.Bind(&apiRequests)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data, potentially invalid structure", nil)
	}

	if len(apiRequests) == 0 {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Request data must contain at least 1 Expected Machine entry", nil)
	}

	if len(apiRequests) > model.ExpectedMachineMaxBatchItems {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("At most %d Expected Machine entries can be created in a batch request", model.ExpectedMachineMaxBatchItems), nil)
	}

	// Validate each item:
	// - ID is required and must be unique
	// - BMC address is optional but must be unique
	// - Serial Number s optional but must be unique
	// Note: this is early partial validation before we try to call the DB. Full unicity validation including DB data will
	// be done later within a transaction.
	idMap := make(map[uuid.UUID]int)
	bmcMacMap := make(map[string]int)
	serialMap := make(map[string]int)
	for i, req := range apiRequests {
		if verr := req.Validate(); verr != nil {
			logger.Warn().Err(verr).Int("Index", i).Msg("error validating Expected Machine update request data")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate Expected Machine update data", validation.Errors{
				strconv.Itoa(i): verr,
			})
		}
		// validation must accept nil ID for single update use case so we need to check for nil ID here
		if req.ID == nil {
			logger.Warn().Int("Index", i).Msg("missing required Expected Machine ID in update request")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Missing required Expected Machine ID at index %d", i), nil)
		}

		// extract already validated UUID
		mid, _ := uuid.Parse(*req.ID)

		if prev, ok := idMap[mid]; ok {
			logger.Warn().Msgf("duplicate Expected Machine ID '%s' found at indices %d and %d", *req.ID, prev, i)
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Duplicate Expected Machine ID '%s' found at indices %d and %d", *req.ID, prev, i), nil)
		}
		idMap[mid] = i

		if req.BmcMacAddress != nil {
			if prev, ok := bmcMacMap[*req.BmcMacAddress]; ok {
				logger.Warn().Msgf("duplicate BMC MAC address '%s' found at indices %d and %d", *req.BmcMacAddress, prev, i)
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Duplicate BMC MAC address '%s' found at indices %d and %d", *req.BmcMacAddress, prev, i), nil)
			}
			bmcMacMap[*req.BmcMacAddress] = i
		}

		if req.ChassisSerialNumber != nil {
			if prev, ok := serialMap[*req.ChassisSerialNumber]; ok {
				logger.Warn().Msgf("duplicate chassis serial number '%s' found at indices %d and %d", *req.ChassisSerialNumber, prev, i)
				return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Duplicate chassis serial number '%s' found at indices %d and %d", *req.ChassisSerialNumber, prev, i), nil)
			}
			serialMap[*req.ChassisSerialNumber] = i
		}
	}

	logger.Info().
		Int("MachineCount", len(apiRequests)).
		Msg("processing UpdateExpectedMachines request")

	// Since we only have a list of Expected Machine ID as input we can only learn the SiteIDs involved by querying the DB
	// but we also want to retrieve full Expected Machine to check for Serial uniqueness.
	// We will split into three queries:
	// 1. Only retrieve SiteIDs involved and check for SiteID unicity
	// 2. Load site record
	// 3. Retrieve Expected Machines for Site to check for serial uniqueness.
	// Pros:
	// - no need to load associated sites for every ExpectedMachine on Site
	// - we can do Provider/Tenant/Site validation before starting transaction and doing any heavy querying/locking which
	//   match our regular pattern
	// Cons:
	// - more queries (3 separate queries instead of 1)
	// Note: we do the Serial uniqueness as best-effort and not under lock or transaction. If we want stricter constraints
	//       we should look at below suggestion.
	// TODO: now that we have a unique index on (mac,siteID) we should reconsider adding unique indices on (serial,siteID)
	//       since it would remove a lot of code for unicity checks. At this time it is expected that existing serial data
	//       may not be unique so we cannot add such an index without cleaning existing data first.

	// Query database directly for unique Site IDs matching our Expected Machine IDs
	var uniqueSiteIDs []uuid.UUID
	err = uemh.dbSession.DB.NewSelect().
		Model((*cdbm.ExpectedMachine)(nil)).
		Column("site_id").
		Where("id IN (?)", bun.In(slices.Collect(maps.Keys(idMap)))).
		Distinct().
		Scan(ctx, &uniqueSiteIDs)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving unique Site IDs for Expected Machines")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site ID for Expected Machines due to DB error", nil)
	}
	if len(uniqueSiteIDs) == 0 {
		logger.Warn().Msg("No Expected Machines found for provided IDs")
		return cerr.NewAPIErrorResponse(c, http.StatusNotFound, "No Expected Machines found for provided IDs", nil)
	}
	if len(uniqueSiteIDs) > 1 {
		logger.Warn().Int("SiteIDCount", len(uniqueSiteIDs)).Msg("all Expected Machines must belong to the same site")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "All Expected Machines in batch must belong to the same site", nil)
	}
	// Get our unique Site ID
	siteID := uniqueSiteIDs[0]

	// Retrieve the Site from the DB
	site, err := cdbm.NewSiteDAO(uemh.dbSession).GetByID(ctx, nil, siteID, nil, false)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request data does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request data due to DB error", nil)
	}

	// Validate ProviderTenantSite relationship and state
	code, message := validateProviderTenantSiteAccess(ctx, logger, uemh.dbSession, infrastructureProvider, tenant, site)
	if code != http.StatusOK {
		return cerr.NewAPIErrorResponse(c, code, message, nil)
	}

	// Retrieve ExpectedMachines to update from DB to allow unicity checks
	// Note: we MUST retrieve all records from the site to check for unicity for the full site.
	emDAO := cdbm.NewExpectedMachineDAO(uemh.dbSession)
	expectedMachines, _, err := emDAO.GetAll(ctx, nil, cdbm.ExpectedMachineFilterInput{
		SiteIDs: []uuid.UUID{siteID},
	}, paginator.PageInput{
		Limit: cdb.GetIntPtr(paginator.TotalLimit), // we want ALL records on site
	}, []string{})
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Expected Machines from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Expected Machines due to DB error", nil)
	}

	// Verify unicity of Serial Numbers with existing records on Site

	// Create maps for easier lookup
	// Note: we do a sanity check to ensure retrieved records are unique since we don't want to report an error later to the caller not related to their input.
	expectedMachineSerialNumberMap := make(map[uuid.UUID]string)
	uniqueSerialNumbers := make(map[string]bool)
	for i := range expectedMachines { // iterate on ALL Expected Machine on Site
		em := &expectedMachines[i]
		if em.ChassisSerialNumber == "" {
			continue
		}
		// loaded serial numbers should be unique
		if uniqueSerialNumbers[em.ChassisSerialNumber] {
			logger.Error().Str("ChassisSerialNumber", em.ChassisSerialNumber).Msg("duplicate chassis serial number found in DB for Expected Machines being updated")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate chassis serial number uniqueness for Expected Machines due to DB data error", nil)
		}
		uniqueSerialNumbers[em.ChassisSerialNumber] = true
		expectedMachineSerialNumberMap[em.ID] = em.ChassisSerialNumber
	}
	// Apply changes to Serial unicity maps and check for conflicts.
	// We can only check once all changes are applied (there could be some swap).
	for _, req := range apiRequests {
		mid, _ := uuid.Parse(*req.ID)
		if req.ChassisSerialNumber != nil {
			expectedMachineSerialNumberMap[mid] = *req.ChassisSerialNumber
		}
	}

	// Re-validate unicity now including fully applied changes:
	uniqueSerialNumbers = make(map[string]bool)
	for i := range expectedMachines { // iterate on ALL Expected Machine on Site
		em := &expectedMachines[i]
		if em.ChassisSerialNumber == "" {
			continue
		}
		// loaded and updated serial numbers should be unique
		if uniqueSerialNumbers[em.ChassisSerialNumber] {
			logger.Error().Str("ChassisSerialNumber", em.ChassisSerialNumber).Msg("duplicate chassis serial number found in DB for Expected Machines if update was applied")
			return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate chassis serial number uniqueness for Expected Machines update", nil)
		}
		uniqueSerialNumbers[em.ChassisSerialNumber] = true
	}

	// Validate that all specified SKU IDs exist
	skuDAO := cdbm.NewSkuDAO(uemh.dbSession)

	skuIDs := make([]string, 0, len(apiRequests))
	for _, machine := range apiRequests {
		if machine.SkuID != nil {
			skuIDs = append(skuIDs, *machine.SkuID)
		} else {
			skuIDs = append(skuIDs, "")
		}
	}

	missingIndices, verr := validateSkuIDs(ctx, nil, skuDAO, site.ID, skuIDs)
	if verr != nil {
		logger.Warn().Err(verr).Msg("error validating SKU IDs in request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate SKU ID in request data due to DB error", nil)
	}

	if len(missingIndices) > 0 {
		indices, _ := json.Marshal(missingIndices)
		logger.Warn().Msg(fmt.Sprintf("SKU ID specified at indices %s do not exist", indices))
		return cerr.NewAPIErrorResponse(c, http.StatusUnprocessableEntity, fmt.Sprintf("SKU ID specified at indices %s do not exist", indices), nil)
	}

	// Prepare ExpectedMachines input for DB
	updateInputs := make([]cdbm.ExpectedMachineUpdateInput, 0, len(apiRequests))
	for _, machineReq := range apiRequests {
		// APIExpectedMachineUpdateRequest must allow nil ID for single update use case. If present here, it has already been validated.
		if machineReq.ID == nil {
			logger.Error().Msg("Expected Machine ID cannot be nil")
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Expected Machine ID cannot be nil", nil)
		}

		emID, _ := uuid.Parse(*machineReq.ID)
		updateInputs = append(updateInputs, cdbm.ExpectedMachineUpdateInput{
			ExpectedMachineID:        emID,
			BmcMacAddress:            machineReq.BmcMacAddress,
			ChassisSerialNumber:      machineReq.ChassisSerialNumber,
			SkuID:                    machineReq.SkuID,
			FallbackDpuSerialNumbers: machineReq.FallbackDPUSerialNumbers,
			Labels:                   machineReq.Labels,
		})
	}

	// Start a db transaction
	tx, err := cdb.BeginTx(ctx, uemh.dbSession, &sql.TxOptions{})
	if err != nil {
		logger.Error().Err(err).Msg("unable to start transaction")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Expected Machines due to DB transaction error", nil)
	}
	// this variable is used in cleanup actions to indicate if this transaction committed
	txCommitted := false
	defer common.RollbackTx(ctx, tx, &txCommitted)

	// Update provided ExpectedMachines in DB
	updatedExpectedMachines, err := emDAO.UpdateMultiple(ctx, tx, updateInputs)
	if err != nil {
		logger.Error().Err(err).Msg("error updating ExpectedMachine records in DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Expected Machine due to DB error", nil)
	}

	workflowMachines := make([]*cwssaws.ExpectedMachine, 0, len(updatedExpectedMachines))
	for i, updatedMachine := range updatedExpectedMachines {
		workflowMachine := &cwssaws.ExpectedMachine{
			Id:                       &cwssaws.UUID{Value: updatedMachine.ID.String()},
			BmcMacAddress:            updatedMachine.BmcMacAddress,
			ChassisSerialNumber:      updatedMachine.ChassisSerialNumber,
			FallbackDpuSerialNumbers: updatedMachine.FallbackDpuSerialNumbers,
			SkuId:                    updatedMachine.SkuID,
		}

		if apiRequests[i].DefaultBmcUsername != nil {
			workflowMachine.BmcUsername = *apiRequests[i].DefaultBmcUsername
		}

		if apiRequests[i].DefaultBmcPassword != nil {
			workflowMachine.BmcPassword = *apiRequests[i].DefaultBmcPassword
		}

		protoLabels := util.ProtobufLabelsFromAPILabels(updatedMachine.Labels)
		if protoLabels != nil {
			workflowMachine.Metadata = &cwssaws.Metadata{
				Labels: protoLabels,
			}
		}

		workflowMachines = append(workflowMachines, workflowMachine)
	}

	logger.Info().Int("Count", len(workflowMachines)).Msg("triggering Expected Machine update workflow")

	// Create workflow request
	workflowRequest := &cwssaws.BatchExpectedMachineOperationRequest{
		ExpectedMachines:     &cwssaws.ExpectedMachineList{ExpectedMachines: workflowMachines},
		AcceptPartialResults: false,
	}

	// Create workflow options
	workflowID := fmt.Sprintf("expected-machines-update-batch-%s-%d", site.ID.String(), len(workflowMachines))
	workflowOptions := tclient.StartWorkflowOptions{
		ID:                       workflowID,
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	// Get the Temporal client for the site we are working with
	stc, err := uemh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Execute workflow and get results
	workflowRun, err := stc.ExecuteWorkflow(ctx, workflowOptions, "UpdateExpectedMachines", workflowRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to schedule batch Expected Machine update workflow on Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to schedule batch Expected Machine update workflow on Site: %v", err), nil)
	}

	workflowRunID := workflowRun.GetID()
	logger = logger.With().Str("WorkflowID", workflowRunID).Logger()
	logger.Info().Msg("executing Expected Machine update workflow on Site")

	// Get workflow results
	var workflowResult cwssaws.BatchExpectedMachineOperationResponse

	err = workflowRun.Get(ctx, &workflowResult)
	if err != nil {
		logger.Error().Err(err).Msg("error executing batch Expected Machine update workflow on Site")
		// Workflow failed entirely - don't commit transaction, changes will be rolled back
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to execute batch Expected Machine update workflow on Site: %v", err), nil)
	}

	// sanity checks since this is all-or-nothing
	if len(workflowResult.GetResults()) != len(updatedExpectedMachines) {
		logger.Warn().Msgf("workflow returned a different number of Expected Machines (expected %d but got %d)", len(updatedExpectedMachines), len(workflowResult.GetResults()))
	}

	// Commit transaction - only successful updates remain in DB
	err = tx.Commit()
	if err != nil {
		logger.Error().Err(err).Msg("error committing ExpectedMachine UpdateExpectedMachines transaction to DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to update Expected Machines due to DB transaction error", nil)
	}
	txCommitted = true

	logger.Info().
		Int("SuccessCount", len(updatedExpectedMachines)).
		Msg("finishing UpdateExpectedMachines API handler")

	// Return only successful machines
	return c.JSON(http.StatusOK, updatedExpectedMachines)
}
