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
package processor

import (
	"context"
	"fmt"
)

type DryRuner struct {
}

func newDryruner() (*DryRuner, error) {
	return &DryRuner{}, nil
}

func (d *DryRuner) Process(ctx context.Context, item *ProcessableItem) error {
	if item.NVLDomain != nil {
		fmt.Printf(
			"******** NVL domain %s *********\n",
			item.NVLDomain.Identifier.Name,
		)
		fmt.Printf("%v\n", item.NVLDomain)
	}

	for _, r := range item.Racks {
		fmt.Printf("******** rack %s *********\n", r.Info.Name)
		fmt.Printf("%v\n", r)
	}

	return nil
}

func (d *DryRuner) Type() Type {
	return TypeDryRun
}

func (d *DryRuner) Done() error {
	return nil
}
