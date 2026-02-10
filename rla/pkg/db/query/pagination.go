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
package Query

import (
	"errors"
)

const (
	SortAscending  SortDirection = "ASC"
	SortDescending SortDirection = "DESC"
)

type SortDirection string

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

func (p *Pagination) Validate() error {
	if p == nil {
		return nil
	}

	if p.Offset < 0 {
		return errors.New("offset must not be negative")
	}

	if p.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}

	return nil
}

type SortBy struct {
	Column    string        `json:"column"`
	Direction SortDirection `json:"direction"`
}

func (sb *SortBy) Validate() error {
	if sb == nil {
		return nil
	}

	if sb.Column == "" {
		return errors.New("column is required")
	}

	if sb.Direction != SortAscending && sb.Direction != SortDescending {
		return errors.New("direction must be ASC or DESC")
	}

	return nil
}

func (sb *SortBy) String() string {
	if sb.Direction == "" {
		return sb.Column
	}
	return sb.Column + " " + string(sb.Direction)
}
