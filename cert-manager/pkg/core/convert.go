// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package core

import (
	"encoding/json"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

// ConvertVaultResponseToStruct marshals the response
func ConvertVaultResponseToStruct(s *vault.Secret, out interface{}) error {
	if s == nil {
		return fmt.Errorf("Vault response is empty")
	}

	return ConvertToStruct(s.Data, out)
}

// ConvertToStruct reads in the map and writes to struct
func ConvertToStruct(in map[string]interface{}, out interface{}) error {
	s, err := json.Marshal(in)
	if err != nil {
		return err
	}

	err = json.Unmarshal(s, out)
	return err
}

// ConvertToInterfaceMap reads the input and writes to map
func ConvertToInterfaceMap(in interface{}, out *map[string]interface{}) error {
	s, err := json.Marshal(in)
	if err != nil {
		return err
	}

	err = json.Unmarshal(s, out)
	return err
}
