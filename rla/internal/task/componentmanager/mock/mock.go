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
package mock

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/rla/internal/task/componentmanager"
	"github.com/nvidia/carbide-rest/rla/internal/task/executor/temporalworkflow/common"
	"github.com/nvidia/carbide-rest/rla/internal/task/operations"
	"github.com/nvidia/carbide-rest/rla/pkg/common/devicetypes"
)

const (
	// ImplementationName is the name used to identify mock implementations.
	ImplementationName = "mock"

	// DefaultDelay is the simulated delay for mock operations.
	DefaultDelay = time.Second
)

// Manager is a mock component manager for testing and development.
type Manager struct {
	componentType devicetypes.ComponentType
	delay         time.Duration
}

// New creates a new mock Manager for the specified component type.
func New(componentType devicetypes.ComponentType) *Manager {
	return &Manager{
		componentType: componentType,
		delay:         DefaultDelay,
	}
}

// NewWithDelay creates a new mock Manager with a custom delay.
func NewWithDelay(componentType devicetypes.ComponentType, delay time.Duration) *Manager {
	return &Manager{
		componentType: componentType,
		delay:         delay,
	}
}

// FactoryFor creates a factory function for the specified component type.
func FactoryFor(componentType devicetypes.ComponentType) componentmanager.ManagerFactory {
	return func(providers *componentmanager.ProviderRegistry) (componentmanager.ComponentManager, error) {
		return New(componentType), nil
	}
}

// RegisterAll registers mock factories for all component types.
func RegisterAll(registry *componentmanager.Registry) {
	for _, ct := range []devicetypes.ComponentType{
		devicetypes.ComponentTypeCompute,
		devicetypes.ComponentTypeNVLSwitch,
		devicetypes.ComponentTypePowerShelf,
	} {
		registry.RegisterFactory(ct, ImplementationName, FactoryFor(ct))
	}
}

// Type returns the component type this manager handles.
func (m *Manager) Type() devicetypes.ComponentType {
	return m.componentType
}

// InjectExpectation simulates injecting expected configuration.
func (m *Manager) InjectExpectation(
	ctx context.Context,
	target common.Target,
	info operations.InjectExpectationTaskInfo,
) error {
	log.Debug().
		Str("component_type", m.componentType.String()).
		Str("target", target.String()).
		Msg("Mock: InjectExpectation")

	time.Sleep(m.delay)
	return nil
}

// PowerControl simulates power operations.
func (m *Manager) PowerControl(
	ctx context.Context,
	target common.Target,
	info operations.PowerControlTaskInfo,
) error {
	log.Debug().
		Str("component_type", m.componentType.String()).
		Str("target", target.String()).
		Str("operation", info.Operation.String()).
		Msg("Mock: PowerControl")

	time.Sleep(m.delay)

	log.Info().
		Str("component_type", m.componentType.String()).
		Str("target", target.String()).
		Str("operation", info.Operation.String()).
		Msg("Mock: PowerControl completed")

	return nil
}

// FirmwareControl simulates firmware operations.
func (m *Manager) FirmwareControl(
	ctx context.Context,
	target common.Target,
	info operations.FirmwareControlTaskInfo,
) error {
	log.Debug().
		Str("component_type", m.componentType.String()).
		Str("target", target.String()).
		Str("operation", info.Operation.String()).
		Str("target_version", info.TargetVersion).
		Msg("Mock: FirmwareControl")

	time.Sleep(m.delay)

	log.Info().
		Str("component_type", m.componentType.String()).
		Str("target", target.String()).
		Str("operation", info.Operation.String()).
		Msg("Mock: FirmwareControl completed")

	return nil
}
