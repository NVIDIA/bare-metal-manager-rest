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
package dumper

import (
	"encoding/json"
)

type Summary struct {
	SourceType     SourceType `json:"source_type"`
	RackNames      []string   `json:"rack_names"`
	NVLDomainNames []string   `json:"nvl_domain_names"`
}

func NewSummary(bytes []byte) (*Summary, error) {
	var summary Summary

	if err := json.Unmarshal(bytes, &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

func (s *Summary) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}
