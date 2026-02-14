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
package cerebro

import "fmt"

const (
	DefaultLocationsLimit = 100
	MaxLocationsSupported = 1000
)

type NVLDomainInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NVLDomain struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Rack    Rack     `json:"rack"`
	Members []Device `json:"members"`
}

type Rack struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	AssetTag string   `json:"asset_tag"`
	Serial   string   `json:"serial"`
	Location Location `json:"location"`
	Devices  []Device `json:"devices"`
}

type Location struct {
	Name   string         `json:"name"`
	Type   LocationType   `json:"location_type"`
	Parent ParentLocation `json:"parent"`
}

type ParentLocation struct {
	Name string `json:"name"`
}

type LocationType struct {
	Name string `json:"name"`
}

type DeviceType struct {
	Model        string       `json:"model"`
	Manufacturer Manufacturer `json:"manufacturer"`
	Family       DeviceFamily `json:"device_family"`
}

type Manufacturer struct {
	Name string `json:"name"`
}

type DeviceFamily struct {
	Name string `json:"name"`
}

type DeviceRole struct {
	Name string `json:"name"`
}

type DeviceBay struct {
	InstalledDevice InstalledDevice `json:"installed_device"`
}

type InstalledDevice struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Device struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Serial     string      `json:"serial"`
	Position   int         `json:"position"`
	Role       DeviceRole  `json:"role"`
	DeviceBays []DeviceBay `json:"device_bays"`
	Type       DeviceType  `json:"device_type"`
	Interfaces []Interface `json:"interfaces"`
}

type Interface struct {
	Name        string      `json:"name"`
	IPAddresses []IPAddress `json:"ip_addresses"`
	MACAddress  string      `json:"mac_address"`
}

type IPAddress struct {
	Address string `json:"address"`
}

type Request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type RackResponse struct {
	Data RackData `json:"data"`
}

type RackData struct {
	Racks []Rack `json:"racks"`
}

type LocationsResponse struct {
	Data LocationsData `json:"data"`
}

type LocationsData struct {
	Locations []Location `json:"locations"`
}

type NVLDomainResponse struct {
	Data NVLDomainData `json:"data"`
}

type NVLDomainData struct {
	NVLDomains []NVLDomain `json:"nvlink_domains"`
}

type AllNVLDomainResponse struct {
	Data AllNVLDomainData `json:"data"`
}

type AllNVLDomainData struct {
	NVLDomains []NVLDomainInfo `json:"nvlink_domains"`
}

type CError map[string][]string

func (e CError) Error() string {
	if len(e) == 0 {
		return ""
	}

	for k, v := range e {
		if len(v) > 0 {
			return fmt.Sprintf("%s: %s", k, v[0])
		}
	}
	return ""
}
