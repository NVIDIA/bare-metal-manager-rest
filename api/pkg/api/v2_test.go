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

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/nvidia/carbide-rest/api/internal/config"
	"github.com/nvidia/carbide-rest/api/pkg/api/handler/util/common"
	sc "github.com/nvidia/carbide-rest/api/pkg/client/site"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"

	temporalClient "go.temporal.io/sdk/client"
	tmocks "go.temporal.io/sdk/mocks"
)

func TestNewV2APIRoutes(t *testing.T) {
	type args struct {
		dbSession *cdb.Session
		tc        temporalClient.Client
		tnc       temporalClient.NamespaceClient
		scp       *sc.ClientPool
		cfg       *config.Config
	}

	tc := &tmocks.Client{}
	tnc := &tmocks.NamespaceClient{}

	cfg := common.GetTestConfig()
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	routeCount := map[string]int{
		"metadata":                1,
		"service-account":         1,
		"infrastructure-provider": 4,
		"tenant":                  4,
		"tenant-account":          5,
		"site":                    6,
		"vpc":                     6,
		"vpcprefix":               5,
		"ip-block":                6,
		"instance":                7,
		"interface":               1,
		"infiniband-partition":    5,
		"nvlinklogical-partition": 5,
		"expected-machine":        5,
		"instance-type":           5,
		"machine":                 5,
		"allocation":              6,
		"subnet":                  5,
		"machine-instance-type":   3,
		"user":                    1,
		"operating-system":        5,
		"sshkey":                  5,
		"sshkeygroup":             5,
		"machine-capability":      1,
		"audit":                   2,
		"network-security-group":  5,
		"machine-validation":      11,
		"dpu-extension-service":   7,
		"sku":                     2,
	}

	totalRouteCount := 0
	for _, v := range routeCount {
		totalRouteCount += v
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test initializing v2 API routes",
			args: args{
				dbSession: &cdb.Session{},
				tc:        tc,
				tnc:       tnc,
				scp:       scp,
				cfg:       cfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewV2APIRoutes(tt.args.dbSession, tt.args.tc, tt.args.tnc, tt.args.scp, tt.args.cfg)

			assert.Equal(t, totalRouteCount, len(got))
		})
	}
}
