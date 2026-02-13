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
package inventorysync

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/carbideapi"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/config"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/psmapi"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/devicetypes"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/db"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/db/model"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/db/postgres"
)

// RunInventory will loop and handle various inventory monitoring tasks
func RunInventory(ctx context.Context, dbConf *db.Config) {
	config := config.ReadConfig()
	if config.DisableInventory {
		log.Info().Msg("Inventory disabled by configuration")
		return
	}

	carbideClient, err := carbideapi.NewClient(config.GRPCTimeout)
	if err != nil {
		// Use whether CARBIDE_API_URL is set to determine if we're running in a production environment (fail hard) or not (just complain and do nothing)
		// Note that this doesn't actually create a connection immediately, so it won't fail just because carbide-api hasn't started yet.
		msg := fmt.Sprintf("Unable to create GRPC client (pre-connect): %v", err)
		if os.Getenv("CARBIDE_API_URL") == "" {
			log.Error().Msg(msg)
			return
		} else {
			log.Fatal().Msg(msg)
		}
	}

	psmClient, err := psmapi.NewClient(config.GRPCTimeout)
	if err != nil {
		log.Error().Msgf("Unable to create PSM GRPC client (PSM_API_URL: %v): %v", os.Getenv("PSM_API_URL"), err)
		return
	}

	if psmClient != nil {
		defer psmClient.Close()
	}

	pool, err := postgres.New(ctx, *dbConf)
	if err != nil {
		log.Fatal().Msgf("Unable to create database pool: %v", err)
	}

	log.Info().Msg("Starting inventory monitoring loop")

	for {
		runInventoryOne(ctx, &config, pool, carbideClient, psmClient)
	}
}

// runInventoryOne is a single iteration for RunInventory
func runInventoryOne(ctx context.Context, config *config.Config, pool *postgres.Postgres, carbideClient carbideapi.Client, psmClient psmapi.Client) {
	if config.UpdateMachineIDsFrequency == 0 {
		// A frequency of zero means to do it only once on startup
		if lastUpdateMachineIDs.IsZero() {
			updateMachineIDs(ctx, pool, carbideClient)
		}
	} else {
		if lastUpdateMachineIDs.Before(time.Now().Add(-config.UpdateMachineIDsFrequency)) {
			updateMachineIDs(ctx, pool, carbideClient)
		}
	}

	updatePowerStates(ctx, pool, carbideClient)

	// Handle expected powershelves - register with PSM if DHCPed, update firmware versions
	handleExpectedPowershelves(ctx, pool, carbideClient, psmClient)

	time.Sleep(config.InventoryRunFrequency)
}

var lastUpdateMachineIDs time.Time

// updatePowerStates requests the current power state for all components that we have recorded a machine ID for
// and updates our databse with it
func updatePowerStates(ctx context.Context, pool *postgres.Postgres, carbideClient carbideapi.Client) {
	containers, err := model.GetAllComponents(ctx, pool.DB())
	if err != nil {
		log.Error().Msgf("Unable to retrieve components from db: %v", err)
		return
	}

	machineIDs := []string{}
	for _, cur := range containers {
		if cur.ComponentID != nil {
			machineIDs = append(machineIDs, *cur.ComponentID)
		}
	}

	machines, err := carbideClient.GetPowerStates(ctx, machineIDs)
	if err != nil {
		log.Error().Msgf("Unable to retreive machines from carbide-api: %v", err)
		return
	}

	containersByMachineID := map[string]model.Component{}
	for _, cur := range containers {
		if cur.ComponentID != nil {
			containersByMachineID[*cur.ComponentID] = cur
		}
	}

	var toUpdate []model.Component
	for _, cur := range machines {
		if container, ok := containersByMachineID[cur.MachineID]; ok {
			if container.PowerState == nil || *container.PowerState != cur.PowerState {
				powerState := cur.PowerState
				container.PowerState = &powerState
				toUpdate = append(toUpdate, container)
			}
		}
	}
	if len(toUpdate) > 0 {
		if err := pool.RunInTx(ctx, func(ctx context.Context, tx bun.Tx) error {
			for _, cur := range toUpdate {
				if err := cur.SetPowerStateByComponentID(ctx, tx); err != nil {
					return fmt.Errorf("Unable to update power state: %v", err)
				}
			}
			return nil
		}); err != nil {
			log.Error().Msgf("Unable to update components with power state: %v", err)
			return
		}
	}

}

// updateMachineIDs requests all of the machines from carbide-api and matches up any components that are missing machine IDs
func updateMachineIDs(ctx context.Context, pool *postgres.Postgres, carbideClient carbideapi.Client) {
	containers, err := model.GetAllComponents(ctx, pool.DB())
	if err != nil {
		log.Error().Msgf("Unable to retrieve components from db: %v", err)
		return
	}

	// If we already found everything, don't bother to recheck
	missingMachine := false
	for _, cur := range containers {
		if cur.ComponentID == nil {
			missingMachine = true
			break
		}
	}
	if !missingMachine {
		return
	}

	machines, err := carbideClient.GetMachines(ctx)
	if err != nil {
		log.Error().Msgf("Unable to retreive machines from carbide-api: %v", err)
		return
	}

	containersBySerial := map[string]model.Component{}
	for _, cur := range containers {
		containersBySerial[cur.SerialNumber] = cur
	}
	var toUpdate []model.Component
	for _, cur := range machines {
		if cur.ChassisSerial == nil {
			continue
		}
		if container, ok := containersBySerial[*cur.ChassisSerial]; ok {
			if container.ComponentID == nil || *container.ComponentID != cur.MachineID {
				componentID := cur.MachineID
				container.ComponentID = &componentID
				toUpdate = append(toUpdate, container)
			}
		}
	}

	if len(toUpdate) > 0 {
		if err := pool.RunInTx(ctx, func(ctx context.Context, tx bun.Tx) error {
			for _, cur := range toUpdate {
				if err := cur.SetComponentIDBySerial(ctx, tx); err != nil {
					return fmt.Errorf("Unable to update machine ID: %v", err)
				}
			}
			return nil
		}); err != nil {
			log.Error().Msgf("Unable to update components with serial: %v", err)
			return
		}

		log.Info().Msgf("Updated %d machine ID(s)", len(toUpdate))
	}

	// lastUpdateMachineIDs is the last time we ran successfully, not necessarily when we last actually changed something
	lastUpdateMachineIDs = time.Now()
}

// Default factory credentials for powershelf BMCs
const (
	powershelfDefaultUsername = "root"
	powershelfDefaultPassword = "0penBmc"
)

// handleExpectedPowershelves checks if any powershelf PMCs have received DHCP leases
// and registers them with the Powershelf Manager (PSM) service
func handleExpectedPowershelves(ctx context.Context, pool *postgres.Postgres, carbideClient carbideapi.Client, psmClient psmapi.Client) {
	if psmClient == nil {
		log.Debug().Msg("PSM client not available, skipping powershelf handling")
		return
	}

	// Get all PowerShelf components with their PMCs
	expectedPowershelves, err := model.GetComponentsByType(ctx, pool.DB(), devicetypes.ComponentTypePowerShelf)
	if err != nil {
		log.Error().Msgf("Unable to retrieve powershelf components from db: %v", err)
		return
	}

	if len(expectedPowershelves) == 0 {
		return
	}

	// Build map from PMC MAC to component
	// Each powershelf should have exactly one PMC (BMC)
	expectedPowershelvesByPmcMac := make(map[string]*model.Component)
	for i := range expectedPowershelves {
		ps := &expectedPowershelves[i]
		if len(ps.BMCs) != 1 {
			log.Error().Msgf("Powershelf %s has %d BMCs, expected exactly 1; skipping", ps.SerialNumber, len(ps.BMCs))
			continue
		}

		// Validate PMC MAC address
		pmcMacAddr, err := net.ParseMAC(ps.BMCs[0].MacAddress)
		if err != nil || pmcMacAddr == nil {
			log.Error().Msgf("Powershelf %s has invalid BMC MAC address %s; skipping", ps.SerialNumber, ps.BMCs[0].MacAddress)
			continue
		}

		expectedPowershelvesByPmcMac[ps.BMCs[0].MacAddress] = ps
	}

	// Get list of expected PMC MACs
	expectedPmcMacs := make([]string, 0, len(expectedPowershelvesByPmcMac))
	for mac := range expectedPowershelvesByPmcMac {
		expectedPmcMacs = append(expectedPmcMacs, mac)
	}

	// Get currently registered powershelves from PSM
	registeredPowershelves, err := psmClient.GetPowershelves(ctx, expectedPmcMacs)
	if err != nil {
		log.Error().Msgf("Unable to retrieve registered powershelves from PSM: %v", err)
		return
	}

	// Build a map of MAC address -> registered PowerShelf for quick lookup
	registeredPmcsByMac := make(map[string]psmapi.PowerShelf)
	for _, ps := range registeredPowershelves {
		registeredPmcsByMac[ps.PMC.MACAddress] = ps
	}

	// Get all machine interfaces from CarbideAPI to find which PMCs have DHCPed
	interfacesByMac, err := carbideClient.FindInterfaces(ctx)
	if err != nil {
		log.Error().Msgf("Unable to retrieve interfaces from carbide-api: %v", err)
		return
	}

	// Collect powershelves that need to be registered
	var toRegister []psmapi.RegisterPowershelfRequest

	for pmcMac, powershelf := range expectedPowershelvesByPmcMac {

		// Check if already registered with PSM
		if registeredPS, isRegistered := registeredPmcsByMac[pmcMac]; isRegistered {
			needsUpdate := false

			// Update firmware version if it differs
			if registeredPS.PMC.FirmwareVersion != "" && powershelf.FirmwareVersion != registeredPS.PMC.FirmwareVersion {
				powershelf.FirmwareVersion = registeredPS.PMC.FirmwareVersion
				needsUpdate = true
				log.Info().Msgf("Updating firmware version for powershelf %s to %s", pmcMac, registeredPS.PMC.FirmwareVersion)
			}

			// Derive power state from PSUs
			// All on → On, All off → Off, Mix or no PSUs → Unknown
			allOn := len(registeredPS.PSUs) > 0
			allOff := len(registeredPS.PSUs) > 0
			for _, psu := range registeredPS.PSUs {
				if psu.PowerState {
					allOff = false
				} else {
					allOn = false
				}
			}
			psuPowerState := carbideapi.PowerStateUnknown
			if allOn {
				psuPowerState = carbideapi.PowerStateOn
			} else if allOff {
				psuPowerState = carbideapi.PowerStateOff
			}
			if powershelf.PowerState == nil || *powershelf.PowerState != psuPowerState {
				powershelf.PowerState = &psuPowerState
				needsUpdate = true
				log.Info().Msgf("Updating power state for powershelf %s to %v", pmcMac, psuPowerState)
			}

			if needsUpdate {
				if err := powershelf.Patch(ctx, pool.DB()); err != nil {
					log.Error().Msgf("Unable to update powershelf %s: %v", pmcMac, err)
				}
			}
			continue
		}

		// Check if this PMC has an interface entry with an IP address
		iface, found := interfacesByMac[pmcMac]
		if !found || len(iface.Addresses) == 0 {
			// PMC hasn't DHCPed yet
			log.Warn().Msgf("Powershelf PMC %s has not DHCPed yet", pmcMac)
			continue
		}

		// Check for unexpected multiple IP addresses
		if len(iface.Addresses) > 1 {
			log.Error().Msgf("Powershelf PMC %s has multiple IP addresses assigned (%v), skipping registration", pmcMac, iface.Addresses)
			continue
		}

		ipAddress := iface.Addresses[0]

		log.Info().Msgf("Powershelf PMC %s has DHCPed with IP %s, registering with PSM", pmcMac, ipAddress)

		toRegister = append(toRegister, psmapi.RegisterPowershelfRequest{
			PMCMACAddress:  pmcMac,
			PMCIPAddress:   ipAddress,
			PMCVendor:      psmapi.PMCVendorLiteon, // Default vendor for powershelves
			PMCCredentials: psmapi.Credentials{Username: powershelfDefaultUsername, Password: powershelfDefaultPassword},
		})
	}

	if len(toRegister) == 0 {
		return
	}

	// Register all powershelves with PSM
	responses, err := psmClient.RegisterPowershelves(ctx, toRegister)
	if err != nil {
		log.Error().Msgf("Unable to register powershelves with PSM: %v", err)
		return
	}

	// Log the results
	for _, resp := range responses {
		if resp.Status != psmapi.StatusSuccess {
			log.Error().Msgf("Failed to register powershelf %s with PSM: %s", resp.PMCMACAddress, resp.Error)
		} else if resp.IsNew {
			log.Info().Msgf("Successfully registered new powershelf %s with PSM", resp.PMCMACAddress)
		} else {
			log.Debug().Msgf("Powershelf %s was already registered with PSM", resp.PMCMACAddress)
		}
	}
}
