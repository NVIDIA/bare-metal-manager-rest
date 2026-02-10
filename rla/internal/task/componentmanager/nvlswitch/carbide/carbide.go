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
package carbide

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/rla/internal/carbideapi"
	"github.com/nvidia/carbide-rest/rla/internal/task/componentmanager"
	carbideprovider "github.com/nvidia/carbide-rest/rla/internal/task/componentmanager/providers/carbide"
	"github.com/nvidia/carbide-rest/rla/internal/task/executor/temporalworkflow/common"
	"github.com/nvidia/carbide-rest/rla/internal/task/operations"
	"github.com/nvidia/carbide-rest/rla/pkg/common/devicetypes"
)

const (
	// ImplementationName is the name used to identify this implementation.
	ImplementationName = "carbide"
)

// Manager manages NVLink switch components via the Carbide API.
type Manager struct {
	carbideClient carbideapi.Client
}

// New creates a new Carbide-based NVLSwitch Manager instance.
func New(carbideClient carbideapi.Client) *Manager {
	return &Manager{
		carbideClient: carbideClient,
	}
}

// Factory creates a new Manager from the provided providers.
// It retrieves the CarbideProvider from the registry and uses its client.
func Factory(providerRegistry *componentmanager.ProviderRegistry) (componentmanager.ComponentManager, error) {
	provider, err := componentmanager.GetTyped[*carbideprovider.Provider](
		providerRegistry,
		carbideprovider.ProviderName,
	)
	if err != nil {
		return nil, fmt.Errorf("nvlswitch/carbide requires carbide provider: %w", err)
	}

	return New(provider.Client()), nil
}

// Register registers the Carbide NVLSwitch manager factory with the given registry.
func Register(registry *componentmanager.Registry) {
	registry.RegisterFactory(devicetypes.ComponentTypeNVLSwitch, ImplementationName, Factory)
}

// Type returns the component type this manager handles.
func (m *Manager) Type() devicetypes.ComponentType {
	return devicetypes.ComponentTypeNVLSwitch
}

// InjectExpectation injects expected configuration or state information for an NVLink switch.
func (m *Manager) InjectExpectation(
	ctx context.Context,
	target common.Target,
	info operations.InjectExpectationTaskInfo,
) error {
	// TODO: Implement logic to inject expected configuration/state
	// This might involve validating and storing expected firmware versions,
	// link configurations, etc.
	return fmt.Errorf("InjectExpectation not yet implemented for NVLSwitch")
}

// PowerControl performs power operations on an NVLink switch via Carbide API.
func (m *Manager) PowerControl(
	ctx context.Context,
	target common.Target,
	info operations.PowerControlTaskInfo,
) error {
	log.Debug().Msgf(
		"NVLSwitch power control %s op %s activity received",
		target.String(),
		info.Operation.String(),
	)

	if m.carbideClient == nil {
		return fmt.Errorf("carbide client is not configured")
	}

	if err := target.Validate(); err != nil {
		return fmt.Errorf("target is invalid: %w", err)
	}

	// Map common.PowerOperation to carbideapi.SystemPowerControl
	var action carbideapi.SystemPowerControl
	switch info.Operation {
	// Power On
	case operations.PowerOperationPowerOn:
		action = carbideapi.PowerControlOn
	case operations.PowerOperationForcePowerOn:
		action = carbideapi.PowerControlForceOn
	// Power Off
	case operations.PowerOperationPowerOff:
		action = carbideapi.PowerControlGracefulShutdown
	case operations.PowerOperationForcePowerOff:
		action = carbideapi.PowerControlForceOff
	// Restart (OS level)
	case operations.PowerOperationRestart:
		action = carbideapi.PowerControlGracefulRestart
	case operations.PowerOperationForceRestart:
		action = carbideapi.PowerControlForceRestart
	// Reset (hardware level)
	case operations.PowerOperationWarmReset:
		action = carbideapi.PowerControlWarmReset
	case operations.PowerOperationColdReset:
		action = carbideapi.PowerControlColdReset
	default:
		return fmt.Errorf("unknown power operation: %v", info.Operation)
	}

	for _, componentID := range target.ComponentIDs {
		err := m.carbideClient.AdminPowerControl(ctx, componentID, action)
		if err != nil {
			return fmt.Errorf("failed to perform power control on %s: %w",
				componentID, err)
		}
	}

	log.Info().Msgf("power control %s on %s completed",
		info.Operation.String(), target.String())

	return nil
}

// FirmwareControl performs firmware operations on an NVLink switch.
func (m *Manager) FirmwareControl(
	ctx context.Context,
	target common.Target,
	info operations.FirmwareControlTaskInfo,
) error {
	// TODO: Implement firmware control
	switch info.Operation {
	case operations.FirmwareOperationUpgrade:
		// Implement firmware upgrade
		return fmt.Errorf("firmware upgrade not yet implemented for NVLSwitch")
	case operations.FirmwareOperationDowngrade:
		// Implement firmware downgrade
		return fmt.Errorf("firmware downgrade not yet implemented for NVLSwitch")
	case operations.FirmwareOperationRollback:
		// Implement firmware rollback
		return fmt.Errorf("firmware rollback not yet implemented for NVLSwitch")
	case operations.FirmwareOperationVersion:
		// Implement firmware version
		return fmt.Errorf("firmware version not yet implemented for NVLSwitch")
	default:
		return fmt.Errorf("unknown firmware operation: %v", info.Operation)
	}
}
