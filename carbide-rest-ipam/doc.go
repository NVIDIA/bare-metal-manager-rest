// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

/*
Package ipam is a ip address management library for ip's and prefixes (networks).

It uses either memory or postgresql database to store the ip's and prefixes.
You can also bring you own Storage implementation as you need.

Example usage:

	import (
		"fmt"
		goipam "github.com/metal-stack/go-ipam"
	)


	func main() {
		// create a ipamer with in memory storage
		ipam := goipam.New()

		prefix, err := ipam.NewPrefix("192.168.0.0/24")
		if err != nil {
			panic(err)
		}

		ip, err := ipam.AcquireIP(prefix)
		if err != nil {
			panic(err)
		}
		fmt.Printf("got IP: %s", ip.IP)

		err = ipam.ReleaseIP(ip)
		if err != nil {
			panic(err)
		}
		fmt.Printf("IP: %s released.", ip.IP)
	}
*/
package ipam
