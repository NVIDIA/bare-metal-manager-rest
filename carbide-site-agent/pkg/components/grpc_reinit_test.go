/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package elektra

import (
	"testing"
	"time"
)

func TestCarbideClientReinitializationOnCertRenewal(t *testing.T) {
	// Initial setup with TestInitElektra which configures the Carbide client with initial certificates
	TestInitElektra(t)
	initialVersion := testElektra.manager.API.Carbide.GetGRPCClientVersion()

	// Paths for CA and client certificates
	caKeyPath := "/tmp/ca.key"
	caCertPath := "/tmp/ca.crt"

	// Regenerate and replace the certificates to simulate renewal
	SetupTestCerts(t, caCertPath, caKeyPath)

	// Wait a few seconds to allow any background processes to complete
	time.Sleep(time.Second * 5)
	renewedVersion := testElektra.manager.API.Carbide.GetGRPCClientVersion()

	if renewedVersion > initialVersion {
		t.Logf("The Carbide client was successfully reinitialized from version %d to %d.", initialVersion, renewedVersion)
	} else {
		t.Errorf("The Carbide client was not reinitialized as expected. It remains at version %d.", initialVersion)
	}
}
