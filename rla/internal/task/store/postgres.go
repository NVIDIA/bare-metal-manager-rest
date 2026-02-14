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
package store

import (
	"context"

	"github.com/google/uuid"

	taskcommon "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/common"
	taskdef "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/task"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/errors"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/converter/dao"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/db/model"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/db/postgres"
	dbquery "github.com/nvidia/bare-metal-manager-rest/rla/pkg/db/query"
)

// PostgresStore implements the Store interface using PostgreSQL.
type PostgresStore struct {
	pg *postgres.Postgres
}

// NewPostgres creates a new PostgreSQL-backed task store.
func NewPostgres(pg *postgres.Postgres) *PostgresStore {
	return &PostgresStore{pg: pg}
}

// CreateTask creates a new task record.
func (s *PostgresStore) CreateTask(
	ctx context.Context,
	task *taskdef.Task,
) error {
	if err := dao.TaskTo(task).Create(ctx, s.pg.DB()); err != nil {
		return errors.GRPCErrorInternal(err.Error())
	}

	return nil
}

// GetTasks retrieves tasks by their IDs.
func (s *PostgresStore) GetTasks(
	ctx context.Context,
	ids []uuid.UUID,
) ([]*taskdef.Task, error) {
	results := make([]*taskdef.Task, 0, len(ids))
	for _, id := range ids {
		if id == uuid.Nil {
			results = append(results, nil)
			continue
		}

		taskDao, err := model.GetTask(ctx, s.pg.DB(), id)
		if err != nil {
			return nil, errors.GRPCErrorInternal(err.Error())
		}

		results = append(results, dao.TaskFrom(taskDao))
	}

	return results, nil
}

// ListTasks lists tasks matching the given criteria.
func (s *PostgresStore) ListTasks(
	ctx context.Context,
	options *taskcommon.TaskListOptions,
	pagination *dbquery.Pagination,
) ([]*taskdef.Task, int32, error) {
	taskDaos, total, err := model.ListTasks(ctx, s.pg.DB(), options, pagination)
	if err != nil {
		return nil, 0, errors.GRPCErrorInternal(err.Error())
	}

	results := make([]*taskdef.Task, 0, len(taskDaos))
	for _, taskDao := range taskDaos {
		results = append(results, dao.TaskFrom(&taskDao))
	}

	return results, total, nil
}

// UpdateScheduledTask updates task scheduling information.
func (s *PostgresStore) UpdateScheduledTask(
	ctx context.Context,
	task *taskdef.Task,
) error {
	taskDao := dao.TaskTo(task)
	if err := taskDao.UpdateScheduledTask(ctx, s.pg.DB()); err != nil {
		return errors.GRPCErrorInternal(err.Error())
	}

	return nil
}

// UpdateTaskStatus updates the status and message of a task.
func (s *PostgresStore) UpdateTaskStatus(
	ctx context.Context,
	arg *taskdef.TaskStatusUpdate,
) error {
	taskDao := &model.Task{
		ID: arg.ID,
	}

	err := taskDao.UpdateTaskStatus(ctx, s.pg.DB(), arg.Status, arg.Message)
	if err != nil {
		return errors.GRPCErrorInternal(err.Error())
	}

	return nil
}
