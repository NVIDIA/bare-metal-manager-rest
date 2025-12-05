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


package middleware

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/nvidia/carbide-rest/api/pkg/metadata"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

var auditSkipMethods = map[string]bool{
	echo.GET:     true,
	echo.HEAD:    true,
	echo.CONNECT: true,
	echo.TRACE:   true,
	echo.OPTIONS: true,
}

func AuditLog(dbSession *cdb.Session) echo.MiddlewareFunc {
	aeDAO := cdbm.NewAuditEntryDAO(dbSession)
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			// check if method should be skipped
			return auditSkipMethods[c.Request().Method]
		},
		LogURIPath:  true,
		LogMethod:   true,
		LogRemoteIP: true,
		LogStatus:   true,
		LogError:    true,
		LogLatency:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			createInput := cdbm.AuditEntryCreateInput{
				Endpoint:    v.URIPath,
				QueryParams: c.QueryParams(),
				Method:      v.Method,
				StatusCode:  v.Status,
				ClientIP:    v.RemoteIP,
				OrgName:     c.Param("orgName"),
				Timestamp:   v.StartTime,
				Duration:    v.Latency,
				APIVersion:  metadata.Version,
			}
			// get user
			if dbUser, ok := c.Get("user").(*cdbm.User); ok && dbUser != nil {
				createInput.UserID = &dbUser.ID
			}
			if auditEntry, err := aeDAO.Create(c.Request().Context(), nil, createInput); err != nil {
				log.Error().Err(err).Msg("Failed to create audit entry")
			} else if auditEntry != nil {
				c.Set("auditEntryID", auditEntry.ID)
			}
			return nil
		},
	})
}

// fields that should be obfuscated when recording body of the request
var obfuscateFields = []string{
	"ipxeScript",
	"userData",
	"publicKey",
	"defaultBmcUsername",
	"defaultBmcPassword",
}

type ResponseError struct {
	Source  string `json:"source"`
	Message string `json:"message"`
}

func AuditBody(dbSession *cdb.Session) echo.MiddlewareFunc {
	aeDAO := cdbm.NewAuditEntryDAO(dbSession)
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: func(c echo.Context) bool {
			// check if method should be skipped
			return auditSkipMethods[c.Request().Method]
		},
		Handler: func(c echo.Context, reqBody []byte, resBody []byte) {
			auditEntryID, ok := c.Get("auditEntryID").(uuid.UUID)
			if !ok {
				return
			}
			updateInput := cdbm.AuditEntryUpdateInput{
				ID: auditEntryID,
			}
			// save status message
			if c.Response().Status >= 400 {
				responseError := ResponseError{}
				if err := json.Unmarshal(resBody, &responseError); err != nil {
					log.Error().Err(err).Msgf("failed to unmarshall error response %s for audit entry %s", string(resBody), auditEntryID)
				} else {
					updateInput.StatusMessage = cdb.GetStrPtr(responseError.Message)
				}
			}
			// save request body
			if len(reqBody) > 0 {
				var bodyMap map[string]interface{}
				if err := json.Unmarshal(reqBody, &bodyMap); err != nil {
					log.Error().Err(err).Msgf("failed to unmarshall body for audit entry %s", auditEntryID)
				}
				for _, field := range obfuscateFields {
					if _, ok := bodyMap[field]; ok {
						bodyMap[field] = "*******************"
					}
				}
				updateInput.Body = bodyMap
			}
			// update
			if _, err := aeDAO.Update(c.Request().Context(), nil, updateInput); err != nil {
				log.Error().Err(err).Msgf("failed to update audit entry %s", auditEntryID)
			}
		},
	})
}
