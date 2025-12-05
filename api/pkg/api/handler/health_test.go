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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
)

func TestHealthCheckHandler_Handle(t *testing.T) {
	type args struct {
		c echo.Context
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	tests := []struct {
		name string
		hch  HealthCheckHandler
		args args
	}{
		{
			name: "test health check API endpoint",
			hch:  HealthCheckHandler{},
			args: args{
				c: e.NewContext(req, rec),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hch := HealthCheckHandler{}
			err := hch.Handle(tt.args.c)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, rec.Code)

			rhc := &model.APIHealthCheck{}

			serr := json.Unmarshal(rec.Body.Bytes(), rhc)
			assert.NoError(t, serr)

			assert.Equal(t, true, rhc.IsHealthy)
		})
	}
}
