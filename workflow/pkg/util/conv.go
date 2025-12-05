// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package util

import (
	"errors"
	"math"
)

// A convenience function for converting a pointer to
// a native Go integer to a pointer to a uint32 for
// use with a protobuf message. Accepts a pointer to
// an int and returns a uint32 pointer.
//
// If the input is nil, nil will be returned.
// If a pointer to a value greater than
// uint32 max is submitted, an error will be returned.
func GetIntPtrToUint32Ptr(i *int) (*uint32, error) {
	if i == nil {
		return nil, nil
	}

	if *i > math.MaxUint32 {
		return nil, errors.New("conversion to uint32 pointer would exceed uint32 max")
	}

	i32 := uint32(*i)

	return &i32, nil
}

// A convenience function for converting a pointer to
// a uint32 to a pointer to a an int.
//
// If the input is nil, nil will be returned.
func GetUint32PtrToIntPtr(u32 *uint32) *int {
	if u32 == nil {
		return nil
	}

	i := int(*u32)

	return &i
}

func GetUint32Ptr(i uint32) *uint32 {
	return (&i)
}
