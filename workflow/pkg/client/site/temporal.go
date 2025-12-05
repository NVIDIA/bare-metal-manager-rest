// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package site

import (
	"os"

	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	zlogadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"

	tsdkClient "go.temporal.io/sdk/client"
	tsdkConverter "go.temporal.io/sdk/converter"

	"github.com/nvidia/carbide-rest/workflow/internal/config"
)

// ClientPool contains Temporal clients for different site agents
type ClientPool struct {
	tcfg        *config.TemporalConfig
	IDClientMap map[string]tsdkClient.Client
	mutex       sync.RWMutex
}

// GetClientByID returns a Temporal client for given cluster ID
func (cp *ClientPool) GetClientByID(siteID uuid.UUID) (tsdkClient.Client, error) {
	cp.mutex.RLock()

	client, found := cp.IDClientMap[siteID.String()]
	if found {
		cp.mutex.RUnlock()
		return client, nil
	}

	cp.mutex.RUnlock()
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	tLogger := logur.LoggerToKV(zlogadapter.New(zerolog.New(os.Stderr)))

	tc, err := tsdkClient.NewLazyClient(tsdkClient.Options{
		HostPort:  fmt.Sprintf("%v:%v", cp.tcfg.Host, cp.tcfg.Port),
		Namespace: siteID.String(),
		ConnectionOptions: tsdkClient.ConnectionOptions{
			TLS: cp.tcfg.ClientTLSCfg,
		},
		DataConverter: tsdkConverter.NewCompositeDataConverter(
			tsdkConverter.NewNilPayloadConverter(),
			tsdkConverter.NewByteSlicePayloadConverter(),
			tsdkConverter.NewProtoJSONPayloadConverterWithOptions(tsdkConverter.ProtoJSONPayloadConverterOptions{
				AllowUnknownFields: true,
			}),
			tsdkConverter.NewProtoPayloadConverter(),
			tsdkConverter.NewJSONPayloadConverter(),
		),
		Logger: tLogger,
	})

	if err != nil {
		log.Panic().Err(err).Str("Temporal Namespace", siteID.String()).
			Msg("failed to create Temporal client for site")
		return nil, err
	}

	cp.IDClientMap[siteID.String()] = tc

	return tc, nil
}

// NewClientPool initializes and returns a new client pool
func NewClientPool(tcfg *config.TemporalConfig) *ClientPool {
	return &ClientPool{
		tcfg:        tcfg,
		IDClientMap: map[string]tsdkClient.Client{},
	}
}
