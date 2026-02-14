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
// Package store provides the storage layer for task management.
// It defines the TaskStore interface for persisting and retrieving task data.
package store

import (
	"context"

	"github.com/google/uuid"

	taskcommon "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/common"
	taskdef "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/task"
	dbquery "github.com/nvidia/bare-metal-manager-rest/rla/pkg/db/query"
)

// Store defines the interface for task data persistence.
// It provides operations for creating, retrieving, and updating tasks.
type Store interface {
	// CreateTask creates a new task record.
	CreateTask(ctx context.Context, task *taskdef.Task) error

	// GetTasks retrieves tasks by their IDs.
	GetTasks(ctx context.Context, ids []uuid.UUID) ([]*taskdef.Task, error)

	// ListTasks lists tasks matching the given criteria.
	ListTasks(ctx context.Context, options *taskcommon.TaskListOptions, pagination *dbquery.Pagination) ([]*taskdef.Task, int32, error)

	// UpdateScheduledTask updates task scheduling information (execution ID, executor type).
	UpdateScheduledTask(ctx context.Context, task *taskdef.Task) error

	// UpdateTaskStatus updates the status and message of a task.
	UpdateTaskStatus(ctx context.Context, arg *taskdef.TaskStatusUpdate) error
}
