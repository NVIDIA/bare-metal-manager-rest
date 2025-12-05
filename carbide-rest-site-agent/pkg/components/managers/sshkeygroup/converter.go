// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package sshkeygroup

import (
	wflows "github.com/nvidia/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SSHKeyGroupReqTransformer struct {
	FromVersion string
	ToVersion   string
	Op          string
	Request     interface{}
}

type SSHKeyGroupRespTransformer struct {
	FromVersion string
	ToVersion   string
	Op          string
	Response    interface{}
}

type SitesshkeygroupV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           *wflows.UUID           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name         string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Organization string                 `protobuf:"bytes,3,opt,name=organization,proto3" json:"organization,omitempty"`
	Version      string                 `protobuf:"bytes,99,opt,name=version,proto3" json:"version,omitempty"`
	Created      *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created,proto3" json:"created,omitempty"`
	Updated      *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=updated,proto3" json:"updated,omitempty"`
	Deleted      *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

type SitesshkeygroupInfoV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status       wflows.WorkflowStatus `protobuf:"varint,1,opt,name=status,proto3,enum=workflows.v1.common.WorkflowStatus" json:"status,omitempty"`
	ObjectStatus wflows.ObjectStatus   `protobuf:"varint,2,opt,name=object_status,json=objectStatus,proto3,enum=workflows.v1.common.ObjectStatus" json:"object_status,omitempty"`
	StatusMsg    string                `protobuf:"bytes,3,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	sshkeygroup  *SitesshkeygroupV1    `protobuf:"bytes,4,opt,name=sshkeygroup,proto3" json:"sshkeygroup,omitempty"`
}
