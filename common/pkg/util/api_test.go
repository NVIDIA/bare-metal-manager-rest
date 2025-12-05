// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIErrorResponse(t *testing.T) {
	type args struct {
		c       echo.Context
		status  int
		message string
		data    error
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"test": true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "initialize and return error response",
			args: args{
				c:       e.NewContext(req, rec),
				status:  400,
				message: "bad request",
				data:    errors.New("bad request"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewAPIErrorResponse(tt.args.c, tt.args.status, tt.args.message, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("NewAPIErrorResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultHTTPErrorHandler(t *testing.T) {
	type args struct {
		err error
	}

	e := echo.New()

	tests := []struct {
		name            string
		args            args
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "test 404 error handler",
			args: args{
				err: echo.ErrNotFound,
			},
			expectedStatus:  http.StatusNotFound,
			expectedMessage: APIErrorNotFound,
		},
		{
			name: "test 500 error handler",
			args: args{
				err: echo.ErrInternalServerError,
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: APIErrorInternalServer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			ec := e.NewContext(req, rec)

			DefaultHTTPErrorHandler(tt.args.err, ec)

			resp := ec.Response()
			assert.Equal(t, tt.expectedStatus, resp.Status)

			rst := &APIError{}
			err := json.Unmarshal(rec.Body.Bytes(), rst)
			assert.NoError(t, err)

			assert.Equal(t, APIErrorSourceForge, rst.Source)
			assert.Equal(t, tt.expectedMessage, rst.Message)
		})
	}
}
