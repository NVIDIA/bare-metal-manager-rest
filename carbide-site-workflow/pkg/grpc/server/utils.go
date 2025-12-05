// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package server

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	mrand "math/rand"
)

func generateMacAddress() string {
	buf := make([]byte, 6)
	rand.Read(buf)

	// Set the local bit
	buf[0] |= 2
	maca := fmt.Sprintf("Random MAC address: %02x:%02x:%02x:%02x:%02x:%02x\n", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return maca
}

func generateInteger(max int) int {
	s := mrand.NewSource(time.Now().UnixNano())
	r := mrand.New(s)

	return r.Intn(max)
}

func generateIPAddress() string {
	buf := make([]byte, 4)

	ip := mrand.Uint32()
	binary.LittleEndian.PutUint32(buf, ip)

	return fmt.Sprintf("%s\n", net.IP(buf))
}

func getStrPtr(s string) *string {
	sp := s
	return &sp
}
