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

package model

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db"
	"github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/extra/bundebug"

	stracer "github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/tracer"
	"go.opentelemetry.io/otel/trace"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func testGenerateRandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func testInitDB(t *testing.T) *db.Session {
	dbSession := util.GetTestDBSession(t, false)
	dbSession.DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv(""),
	))

	rand.Seed(time.Now().UnixNano())

	return dbSession
}

// reset the tables needed for common tests
func testSetupSchema(t *testing.T, dbSession *db.Session) {
	// create Tenant table
	err := dbSession.DB.ResetModel(context.Background(), (*Tenant)(nil))
	assert.Nil(t, err)
	// create Site table
	err = dbSession.DB.ResetModel(context.Background(), (*Site)(nil))
	assert.Nil(t, err)
	// create Network Security Group table
	err = dbSession.DB.ResetModel(context.Background(), (*NetworkSecurityGroup)(nil))
	assert.Nil(t, err)
	// create Infrastructure Provider table
	err = dbSession.DB.ResetModel(context.Background(), (*InfrastructureProvider)(nil))
	assert.Nil(t, err)
	// create Vpc table
	err = dbSession.DB.ResetModel(context.Background(), (*Vpc)(nil))
	assert.Nil(t, err)
	// create IPBlock table
	err = dbSession.DB.ResetModel(context.Background(), (*IPBlock)(nil))
	assert.Nil(t, err)
	// create Domain table
	err = dbSession.DB.ResetModel(context.Background(), (*Domain)(nil))
	assert.Nil(t, err)
	// create Subnet table
	err = dbSession.DB.ResetModel(context.Background(), (*Subnet)(nil))
	assert.Nil(t, err)
	// create OperatingSystem table
	err = dbSession.DB.ResetModel(context.Background(), (*OperatingSystem)(nil))
	assert.Nil(t, err)
	// create OperatingSystem table
	err = dbSession.DB.ResetModel(context.Background(), (*OperatingSystemSiteAssociation)(nil))
	assert.Nil(t, err)
	// create User table
	err = dbSession.DB.ResetModel(context.Background(), (*User)(nil))
	assert.Nil(t, err)
	// create Allocation table
	err = dbSession.DB.ResetModel(context.Background(), (*Allocation)(nil))
	assert.Nil(t, err)
	// create AllocationConstraint table
	err = dbSession.DB.ResetModel(context.Background(), (*AllocationConstraint)(nil))
	assert.Nil(t, err)
	// create Machine table
	err = dbSession.DB.ResetModel(context.Background(), (*Machine)(nil))
	assert.Nil(t, err)
	// create InstanceType table
	err = dbSession.DB.ResetModel(context.Background(), (*InstanceType)(nil))
	assert.Nil(t, err)
	// create Instance table
	err = dbSession.DB.ResetModel(context.Background(), (*Instance)(nil))
	assert.Nil(t, err)
	// create Interface table
	err = dbSession.DB.ResetModel(context.Background(), (*Interface)(nil))
	assert.Nil(t, err)
	// create MachineCapability table
	err = dbSession.DB.ResetModel(context.Background(), (*MachineCapability)(nil))
	assert.Nil(t, err)
	// create MachineInterface table
	err = dbSession.DB.ResetModel(context.Background(), (*MachineInterface)(nil))
	assert.Nil(t, err)
	// create StatusDetail table
	err = dbSession.DB.ResetModel(context.Background(), (*StatusDetail)(nil))
	assert.Nil(t, err)
	// create Security Policy table
	err = dbSession.DB.ResetModel(context.Background(), (*SSHKey)(nil))
	assert.Nil(t, err)
	// create Security Policy table
	err = dbSession.DB.ResetModel(context.Background(), (*SSHKeyAssociation)(nil))
	assert.Nil(t, err)
	// create Security Policy table
	err = dbSession.DB.ResetModel(context.Background(), (*SSHKeyGroup)(nil))
	assert.Nil(t, err)
	// create Security Policy table
	err = dbSession.DB.ResetModel(context.Background(), (*SSHKeyGroupSiteAssociation)(nil))
	assert.Nil(t, err)
}

func testBuildInfrastructureProvider(t *testing.T, dbSession *db.Session, id *uuid.UUID, name, org string, createdBy uuid.UUID) *InfrastructureProvider {
	ipid := uuid.New()
	if id != nil {
		ipid = *id
	}

	provider := &InfrastructureProvider{
		ID:        ipid,
		Name:      name,
		Org:       org,
		CreatedBy: createdBy,
	}
	_, err := dbSession.DB.NewInsert().Model(provider).Exec(context.Background())
	require.NoError(t, err)

	return provider
}

func testBuildTenant(t *testing.T, dbSession *db.Session, id *uuid.UUID, name, org string, createdBy uuid.UUID) *Tenant {
	tid := uuid.New()
	if id != nil {
		tid = *id
	}

	tenant := &Tenant{
		ID:        tid,
		Name:      name,
		Org:       org,
		CreatedBy: createdBy,
	}
	_, err := dbSession.DB.NewInsert().Model(tenant).Exec(context.Background())
	require.NoError(t, err)

	return tenant
}

func testBuildAllocation(t *testing.T, dbSession *db.Session, ip *InfrastructureProvider, tenant *Tenant, site *Site, name string) *Allocation {
	allocation := &Allocation{
		ID:                       uuid.New(),
		Name:                     name,
		InfrastructureProviderID: ip.ID,
		TenantID:                 tenant.ID,
		SiteID:                   site.ID,
		Status:                   AllocationStatusPending,
		CreatedBy:                uuid.New(),
	}
	_, err := dbSession.DB.NewInsert().Model(allocation).Exec(context.Background())
	assert.Nil(t, err)

	return allocation
}

func testBuildAllocationConstraint(t *testing.T, dbSession *db.Session, allocation *Allocation, resourceType string, resourceTypeID uuid.UUID, constraintType string, constraintValue int, createdBy uuid.UUID) *AllocationConstraint {
	constraint := &AllocationConstraint{
		ID:              uuid.New(),
		AllocationID:    allocation.ID,
		ResourceType:    resourceType,
		ResourceTypeID:  resourceTypeID,
		ConstraintType:  constraintType,
		ConstraintValue: constraintValue,
		CreatedBy:       createdBy,
	}

	_, err := dbSession.DB.NewInsert().Model(constraint).Exec(context.Background())
	assert.Nil(t, err)

	return constraint
}

func testBuildSite(t *testing.T, dbSession *db.Session, id *uuid.UUID, ipID uuid.UUID, name, displayName, org string, createdBy uuid.UUID) *Site {
	sid := uuid.New()
	if id != nil {
		sid = *id
	}

	site := &Site{
		ID:                       sid,
		InfrastructureProviderID: ipID,
		Name:                     name,
		DisplayName:              &displayName,
		Org:                      org,
		CreatedBy:                uuid.New(),
	}

	_, err := dbSession.DB.NewInsert().Model(site).Exec(context.Background())
	require.NoError(t, err)

	return site
}

func testBuildVpc(t *testing.T, dbSession *db.Session, id *uuid.UUID, name string, description *string, org string, infrastructureProviderID uuid.UUID, tenantID uuid.UUID, siteID uuid.UUID, networkVirtualizationType *string, controllerVpcID *uuid.UUID, labels map[string]string, status *string, createdBy uuid.UUID, networkSecurityGroupId *string) *Vpc {
	vid := uuid.New()
	if id != nil {
		vid = *id
	}

	vstatus := VpcStatusPending
	if status != nil {
		vstatus = *status
	}

	vpc := &Vpc{
		ID:                        vid,
		Name:                      name,
		Description:               description,
		Org:                       org,
		InfrastructureProviderID:  infrastructureProviderID,
		TenantID:                  tenantID,
		SiteID:                    siteID,
		NetworkVirtualizationType: networkVirtualizationType,
		ControllerVpcID:           controllerVpcID,
		NetworkSecurityGroupID:    networkSecurityGroupId,
		Labels:                    labels,
		Status:                    vstatus,
		CreatedBy:                 createdBy,
	}

	_, err := dbSession.DB.NewInsert().Model(vpc).Exec(context.Background())
	require.NoError(t, err)

	return vpc
}

func testUpdateVpc(t *testing.T, dbSession *db.Session, v *Vpc) {

	_, err := dbSession.DB.NewUpdate().Where("id = ?", v.ID).Model(v).Exec(context.Background())
	require.NoError(t, err)
}

func testBuildUser(t *testing.T, dbSession *db.Session, id *uuid.UUID, starfleetID string, email, firstName, lastName *string) *User {
	uid := uuid.New()
	if id != nil {
		uid = *id
	}

	user := &User{
		ID:          uid,
		StarfleetID: db.GetStrPtr(starfleetID),
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
	}

	_, err := dbSession.DB.NewInsert().Model(user).Exec(context.Background())
	require.NoError(t, err)

	return user
}

func testBuildSSHKeyGroup(t *testing.T, dbSession *db.Session, name string, description *string, org string, tenantID uuid.UUID, version *string, status string, createdBy uuid.UUID) *SSHKeyGroup {
	sshKeyGroup := &SSHKeyGroup{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Org:         org,
		TenantID:    tenantID,
		Version:     version,
		Status:      status,
		CreatedBy:   createdBy,
	}
	_, err := dbSession.DB.NewInsert().Model(sshKeyGroup).Exec(context.Background())
	assert.Nil(t, err)
	return sshKeyGroup
}

func testBuildSSHKeyGroupSiteAssociation(t *testing.T, dbSession *db.Session, sshKeyGroupID uuid.UUID, siteID uuid.UUID, version *string, status string, createdBy uuid.UUID) *SSHKeyGroupSiteAssociation {
	SSHKeyGroupSiteAssociation := &SSHKeyGroupSiteAssociation{
		ID:            uuid.New(),
		SSHKeyGroupID: sshKeyGroupID,
		SiteID:        siteID,
		Version:       version,
		Status:        status,
		CreatedBy:     createdBy,
	}
	_, err := dbSession.DB.NewInsert().Model(SSHKeyGroupSiteAssociation).Exec(context.Background())
	assert.Nil(t, err)
	return SSHKeyGroupSiteAssociation
}

func testBuildSSHKey(t *testing.T, dbSession *db.Session, name, org string, tenantID uuid.UUID, publicKey string, fingerprint *string, expires *time.Time, createdBy uuid.UUID) *SSHKey {
	sshkey := &SSHKey{
		ID:          uuid.New(),
		Name:        name,
		Org:         org,
		TenantID:    tenantID,
		PublicKey:   publicKey,
		Fingerprint: fingerprint,
		Expires:     expires,
		CreatedBy:   createdBy,
	}
	_, err := dbSession.DB.NewInsert().Model(sshkey).Exec(context.Background())
	assert.Nil(t, err)
	return sshkey
}

func testBuildSSHKeyAssociation(t *testing.T, dbSession *db.Session, sshKeyID uuid.UUID, sshKeyGroupID uuid.UUID, createdBy uuid.UUID) *SSHKeyAssociation {
	sshkeyassociation := &SSHKeyAssociation{
		ID:            uuid.New(),
		SSHKeyID:      sshKeyID,
		SSHKeyGroupID: sshKeyGroupID,
		CreatedBy:     createdBy,
	}
	_, err := dbSession.DB.NewInsert().Model(sshkeyassociation).Exec(context.Background())
	assert.Nil(t, err)
	return sshkeyassociation
}

func testBuildFabric(t *testing.T, dbSession *db.Session, id *string, org string, siteID uuid.UUID, ipID uuid.UUID, isMissingOnSite bool, status *string) *Fabric {

	fb := &Fabric{
		ID:                       *id,
		Org:                      org,
		SiteID:                   siteID,
		InfrastructureProviderID: ipID,
		IsMissingOnSite:          isMissingOnSite,
		Status:                   *status,
	}

	_, err := dbSession.DB.NewInsert().Model(fb).Exec(context.Background())
	require.NoError(t, err)

	return fb
}

func testBuildInfiniBandPartition(t *testing.T, dbSession *db.Session, id *uuid.UUID, name string, description *string, org string, tenantID uuid.UUID, siteID uuid.UUID, controllerIBInfiniBandPartitionID *uuid.UUID, partitionKey *string, partitionName *string, serviceLevel *int, rateLimit *float32, mtu *int, enableSharp *bool, status *string, createdBy uuid.UUID) *InfiniBandPartition {
	pid := uuid.New()
	if id != nil {
		pid = *id
	}

	pstatus := InfiniBandPartitionStatusPending
	if status != nil {
		pstatus = *status
	}

	pt := &InfiniBandPartition{
		ID:                      pid,
		Name:                    name,
		Description:             description,
		Org:                     org,
		SiteID:                  siteID,
		TenantID:                tenantID,
		ControllerIBPartitionID: controllerIBInfiniBandPartitionID,
		PartitionKey:            partitionKey,
		PartitionName:           partitionName,
		ServiceLevel:            serviceLevel,
		RateLimit:               rateLimit,
		Mtu:                     mtu,
		EnableSharp:             enableSharp,
		Status:                  pstatus,
		CreatedBy:               createdBy,
	}

	_, err := dbSession.DB.NewInsert().Model(pt).Exec(context.Background())
	require.NoError(t, err)

	return pt
}

func testBuildInfiniBandInterface(t *testing.T, dbSession *db.Session, id *uuid.UUID, siteID uuid.UUID, instanceID uuid.UUID, infiniBandPartitionID uuid.UUID, deviceInstance int, isPhysical bool, virtualFunctionID *int, guid *string, isMissingOnSite bool, status *string, createdBy uuid.UUID) *InfiniBandInterface {
	pid := uuid.New()
	if id != nil {
		pid = *id
	}

	if status == nil {
		status = db.GetStrPtr(InfiniBandInterfaceStatusPending)
	}

	ibif := &InfiniBandInterface{
		ID:                    pid,
		SiteID:                siteID,
		InstanceID:            instanceID,
		InfiniBandPartitionID: infiniBandPartitionID,
		DeviceInstance:        deviceInstance,
		IsPhysical:            isPhysical,
		VirtualFunctionID:     virtualFunctionID,
		GUID:                  guid,
		IsMissingOnSite:       isMissingOnSite,
		Status:                *status,
		CreatedBy:             createdBy,
	}

	_, err := dbSession.DB.NewInsert().Model(ibif).Exec(context.Background())
	require.NoError(t, err)

	return ibif
}

func testGenerateStarfleetID() string {
	// Example: eKTl_V61-o6x2-Uyz8jqmxkRFrfP_ggyIvzp-ITQSvc
	return fmt.Sprintf("%v-%v-%v-%v", testGenerateRandomString(8), testGenerateRandomString(4), testGenerateRandomString(22), testGenerateRandomString(6))
}

func testGenerateMacAddress(t *testing.T) string {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	assert.NoError(t, err)
	// Set the local bit
	buf[0] = (buf[0] | 2) & 0xfe // Set local bit, ensure unicast address
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x\n", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

// testCommonTraceProviderSetup creates a test provider and spanner
func testCommonTraceProviderSetup(t *testing.T, ctx context.Context) (trace.Tracer, trace.SpanContext, context.Context) {
	// OTEL spanner configuration
	provider := trace.NewNoopTracerProvider()
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{0x01},
		SpanID:  trace.SpanID{0x01},
	})

	ctx = trace.ContextWithRemoteSpanContext(ctx, sc)

	tracer := provider.Tracer(stracer.TracerName)
	tracer.Start(ctx, "Test-DB-Spanner")
	ctx = context.WithValue(ctx, stracer.TracerKey, tracer)

	return tracer, sc, ctx
}

func testBuildImageOperatingSystem(t *testing.T, dbSession *db.Session, name string, description *string, org string, ipID *uuid.UUID, tenantID *uuid.UUID, version *string, enabelBS bool, status string, createdBy uuid.UUID) *OperatingSystem {
	os := &OperatingSystem{
		ID:                       uuid.New(),
		Name:                     name,
		Description:              description,
		Org:                      org,
		InfrastructureProviderID: ipID,
		TenantID:                 tenantID,
		Type:                     "Image",
		ImageURL:                 db.GetStrPtr("imageURL"),
		ImageSHA:                 db.GetStrPtr("imageSHA"),
		ImageAuthType:            db.GetStrPtr("imageAuthType"),
		ImageAuthToken:           db.GetStrPtr("imageAuthToken"),
		ImageDisk:                db.GetStrPtr("imageDisk"),
		RootFsID:                 db.GetStrPtr("rootFsId"),
		RootFsLabel:              db.GetStrPtr("rootFsLabel"),
		Version:                  version,
		EnableBlockStorage:       enabelBS,
		Status:                   status,
		CreatedBy:                createdBy,
	}

	_, err := dbSession.DB.NewInsert().Model(os).Exec(context.Background())
	assert.Nil(t, err)
	return os
}

func testBuildOperatingSystemSiteAssociation(t *testing.T, dbSession *db.Session, osID uuid.UUID, siteID uuid.UUID, version *string, status string, createdBy uuid.UUID) *OperatingSystemSiteAssociation {
	ossa := &OperatingSystemSiteAssociation{
		ID:                uuid.New(),
		OperatingSystemID: osID,
		SiteID:            siteID,
		Version:           version,
		Status:            status,
		CreatedBy:         createdBy,
	}
	_, err := dbSession.DB.NewInsert().Model(ossa).Exec(context.Background())
	assert.Nil(t, err)
	return ossa
}
