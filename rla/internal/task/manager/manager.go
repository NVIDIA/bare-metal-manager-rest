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
package manager

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/rack"
	inventorystore "github.com/nvidia/carbide-rest/rla/internal/inventory/store"
	"github.com/nvidia/carbide-rest/rla/internal/operation"
	taskcommon "github.com/nvidia/carbide-rest/rla/internal/task/common"
	"github.com/nvidia/carbide-rest/rla/internal/task/executor"
	"github.com/nvidia/carbide-rest/rla/internal/task/executor/temporalworkflow/activity"
	"github.com/nvidia/carbide-rest/rla/internal/task/operations"
	taskstore "github.com/nvidia/carbide-rest/rla/internal/task/store"
	taskdef "github.com/nvidia/carbide-rest/rla/internal/task/task"
)

// Config holds the configuration for the task manager.
type Config struct {
	InventoryStore inventorystore.Store // For rack/component lookups (read-only)
	TaskStore      taskstore.Store      // For task persistence
	ExecutorConfig executor.ExecutorConfig
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("configuration is nil")
	}

	if c.InventoryStore == nil {
		return fmt.Errorf("inventory store is required")
	}

	if c.TaskStore == nil {
		return fmt.Errorf("task store is required")
	}

	if c.ExecutorConfig == nil {
		return fmt.Errorf("executor config is required")
	}

	return c.ExecutorConfig.Validate()
}

// Manager maintains unfinished tasks, schedules them via temporal workflows,
// and monitors their progress.
type Manager struct {
	inventoryStore inventorystore.Store // For rack/component lookups
	taskStore      taskstore.Store      // For task persistence
	executor       executor.Executor

	ctx       context.Context
	cancel    context.CancelFunc
	startOnce sync.Once
	stopOnce  sync.Once
}

// New creates a new task manager.
func New(ctx context.Context, conf *Config) (*Manager, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	// Set task store as the status updater for workflow activities
	activity.SetTaskStatusUpdater(conf.TaskStore)

	executor, err := executor.New(ctx, conf.ExecutorConfig)
	if err != nil {
		return nil, err
	}

	return &Manager{
		inventoryStore: conf.InventoryStore,
		taskStore:      conf.TaskStore,
		executor:       executor,
	}, nil
}

// Start starts the task manager to make it ready to accept tasks.
func (m *Manager) Start(ctx context.Context) error {
	var startErr error

	m.startOnce.Do(func() {
		if m.executor == nil {
			startErr = fmt.Errorf("executor is required")
			return
		}

		if err := m.executor.Start(ctx); err != nil {
			startErr = fmt.Errorf("failed to start executor: %w", err)
			return
		}

		startCtx, cancel := context.WithCancel(ctx)
		m.ctx = startCtx
		m.cancel = cancel
	})

	return startErr
}

// Stop shuts down the manager and waits for all routines to finish.
func (m *Manager) Stop(ctx context.Context) {
	m.stopOnce.Do(func() {
		if m.cancel != nil {
			m.cancel()
		}

		if m.executor != nil {
			if err := m.executor.Stop(ctx); err != nil {
				log.Warn().Err(err).Msg("failed to stop executor")
			}
		}
	})
}

// SubmitTask submits a task to the task manager.
// operation.Request can contain multiple racks. Task Manager resolves identifiers,
// splits by rack, and creates one Task per rack.
// Returns a list of created task IDs.
func (m *Manager) SubmitTask(
	ctx context.Context,
	req *operation.Request,
) ([]uuid.UUID, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Resolve targets to racks with components
	rackMap, err := resolveTargetSpecToRacks(ctx, m.inventoryStore, &req.TargetSpec)
	if err != nil {
		return nil, err
	}

	if len(rackMap) == 0 {
		return nil, fmt.Errorf("no valid racks found for request")
	}

	// Create and execute task for each rack
	var taskIDs []uuid.UUID
	for _, targetRack := range rackMap {
		taskID, err := m.createAndExecuteTask(ctx, req, targetRack)
		if err != nil {
			log.Error().Err(err).Str("rack_id", targetRack.Info.ID.String()).Msg("failed to create task for rack")
			continue
		}
		taskIDs = append(taskIDs, taskID)
	}

	return taskIDs, nil
}

// createAndExecuteTask creates a task for a single rack and executes it.
func (m *Manager) createAndExecuteTask(
	ctx context.Context,
	req *operation.Request,
	targetRack *rack.Rack,
) (uuid.UUID, error) {
	// Extract component UUIDs for tracking
	componentUUIDs := make([]uuid.UUID, 0, len(targetRack.Components))
	for _, c := range targetRack.Components {
		componentUUIDs = append(componentUUIDs, c.Info.ID)
	}

	// Create task record
	task := taskdef.Task{
		ID:             uuid.New(),
		Operation:      req.Operation,
		RackID:         targetRack.Info.ID,
		ComponentUUIDs: componentUUIDs,
		Description:    req.Description,
		ExecutorType:   taskcommon.ExecutorTypeUnknown,
		ExecutionID:    "",
		Status:         taskcommon.TaskStatusPending,
		Message:        "Created",
	}

	if err := m.taskStore.CreateTask(ctx, &task); err != nil {
		return uuid.Nil, err
	}

	// Execute the task
	resp, err := m.executeTask(ctx, &task, targetRack)
	if err != nil {
		lerr := m.taskStore.UpdateTaskStatus(
			ctx,
			&taskdef.TaskStatusUpdate{
				ID:      task.ID,
				Status:  taskcommon.TaskStatusFailed,
				Message: err.Error(),
			},
		)
		if lerr != nil {
			log.Error().Err(err).Msgf("failed to update task %s status to failed", task.ID)
			err = errors.Join(err, lerr)
		}
		return uuid.Nil, err
	}

	// Update task with execution info
	task.ExecutionID = resp.ExecutionID
	task.ExecutorType = m.executor.Type()

	if err := m.taskStore.UpdateScheduledTask(ctx, &task); err != nil {
		log.Error().Err(err).Msgf("failed to update task %s scheduled", task.ID)
	}

	return task.ID, nil
}

func (m *Manager) executeTask(
	ctx context.Context,
	task *taskdef.Task,
	targetRack *rack.Rack,
) (*taskdef.ExecutionResponse, error) {
	if task == nil {
		return nil, fmt.Errorf("task is nil")
	}

	req := taskdef.ExecutionRequest{
		Info: taskdef.ExecutionInfo{
			TaskID: task.ID,
			Rack:   targetRack,
		},
		Async: false,
	}

	switch task.Operation.Type {
	case taskcommon.TaskTypePowerControl:
		var info operations.PowerControlTaskInfo
		if err := info.Unmarshal(task.Operation.Info); err != nil {
			return nil, err
		}
		return m.executor.PowerControl(ctx, &req, info)
	case taskcommon.TaskTypeFirmwareControl:
		var info operations.FirmwareControlTaskInfo
		if err := info.Unmarshal(task.Operation.Info); err != nil {
			return nil, err
		}
		return m.executor.FirmwareControl(ctx, &req, info)
	case taskcommon.TaskTypeInjectExpectation:
		var info operations.InjectExpectationTaskInfo
		if err := info.Unmarshal(task.Operation.Info); err != nil {
			return nil, err
		}
		return m.executor.InjectExpectation(ctx, &req, info)
	default:
		return nil, fmt.Errorf("unsupported task type: %s", task.Operation.Type)
	}
}
