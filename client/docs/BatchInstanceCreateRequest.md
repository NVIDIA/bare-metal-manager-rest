# BatchInstanceCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NamePrefix** | **string** | Prefix for instance names. Instances will be named with this prefix followed by a random 6-character suffix (e.g., \&quot;worker\&quot; becomes \&quot;worker-abc123\&quot;) | 
**Count** | **int32** | Number of instances to create in this batch. Minimum 2, maximum 18 (limited by topology domain size) | 
**Description** | Pointer to **NullableString** | Description applied to all instances in the batch, optional | [optional] 
**TenantId** | **string** | ID of the Tenant creating the Instances | 
**InstanceTypeId** | **string** | ID of the Instance Type to use for all Instances in the batch | 
**VpcId** | **string** | ID of the VPC the Instances should belong to | 
**UserData** | Pointer to **NullableString** | User data applied to all instances. Can only be specified if allowOverride is set to true in Operating System | [optional] 
**OperatingSystemId** | Pointer to **NullableString** | Must be specified if iPXE Script field is empty | [optional] 
**NetworkSecurityGroupId** | Pointer to **NullableString** | ID of a Network Security Group to attach to all instances | [optional] 
**IpxeScript** | Pointer to **NullableString** | Override iPXE script specified in OS, must be specified if Operating System is not specified | [optional] 
**AlwaysBootWithCustomIpxe** | Pointer to **bool** | When set to true, the iPXE script specified by OS or overridden here will always be run when rebooting the Instances. OS must be of iPXE type. | [optional] 
**PhoneHomeEnabled** | Pointer to **bool** | When set to true, the Instances will be enabled with the Phone Home service. | [optional] 
**Labels** | Pointer to **map[string]string** |  | [optional] 
**Interfaces** | [**[]InterfaceCreateRequest**](InterfaceCreateRequest.md) | Interface configuration shared across all instances. At least one interface must be specified. Either Subnet or VPC Prefix interfaces allowed, only one of the Subnets or VPC Prefixes can be attached over Physical interface. | 
**InfinibandInterfaces** | Pointer to [**[]InfiniBandInterfaceCreateRequest**](InfiniBandInterfaceCreateRequest.md) | InfiniBand interface configuration shared across all instances | [optional] 
**DpuExtensionServiceDeployments** | Pointer to [**[]DpuExtensionServiceDeploymentRequest**](DpuExtensionServiceDeploymentRequest.md) | DPU Extension Services to deploy to all instances in the batch | [optional] 
**NvLinkInterfaces** | Pointer to [**[]NVLinkInterfaceCreateRequest**](NVLinkInterfaceCreateRequest.md) | NVLink interface configuration shared across all instances | [optional] 
**SshKeyGroupIds** | Pointer to **[]string** | SSH Key Group IDs that will provide Serial over LAN access to all instances | [optional] 
**TopologyOptimized** | Pointer to **bool** | When true (default), all instances must be allocated on machines within the same NVLink domain. When false, instances can be spread across different NVLink domains. | [optional] [default to true]

## Methods

### NewBatchInstanceCreateRequest

`func NewBatchInstanceCreateRequest(namePrefix string, count int32, tenantId string, instanceTypeId string, vpcId string, interfaces []InterfaceCreateRequest, ) *BatchInstanceCreateRequest`

NewBatchInstanceCreateRequest instantiates a new BatchInstanceCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBatchInstanceCreateRequestWithDefaults

`func NewBatchInstanceCreateRequestWithDefaults() *BatchInstanceCreateRequest`

NewBatchInstanceCreateRequestWithDefaults instantiates a new BatchInstanceCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNamePrefix

`func (o *BatchInstanceCreateRequest) GetNamePrefix() string`

GetNamePrefix returns the NamePrefix field if non-nil, zero value otherwise.

### GetNamePrefixOk

`func (o *BatchInstanceCreateRequest) GetNamePrefixOk() (*string, bool)`

GetNamePrefixOk returns a tuple with the NamePrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNamePrefix

`func (o *BatchInstanceCreateRequest) SetNamePrefix(v string)`

SetNamePrefix sets NamePrefix field to given value.


### GetCount

`func (o *BatchInstanceCreateRequest) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *BatchInstanceCreateRequest) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *BatchInstanceCreateRequest) SetCount(v int32)`

SetCount sets Count field to given value.


### GetDescription

`func (o *BatchInstanceCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *BatchInstanceCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *BatchInstanceCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *BatchInstanceCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *BatchInstanceCreateRequest) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *BatchInstanceCreateRequest) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetTenantId

`func (o *BatchInstanceCreateRequest) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *BatchInstanceCreateRequest) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *BatchInstanceCreateRequest) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.


### GetInstanceTypeId

`func (o *BatchInstanceCreateRequest) GetInstanceTypeId() string`

GetInstanceTypeId returns the InstanceTypeId field if non-nil, zero value otherwise.

### GetInstanceTypeIdOk

`func (o *BatchInstanceCreateRequest) GetInstanceTypeIdOk() (*string, bool)`

GetInstanceTypeIdOk returns a tuple with the InstanceTypeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceTypeId

`func (o *BatchInstanceCreateRequest) SetInstanceTypeId(v string)`

SetInstanceTypeId sets InstanceTypeId field to given value.


### GetVpcId

`func (o *BatchInstanceCreateRequest) GetVpcId() string`

GetVpcId returns the VpcId field if non-nil, zero value otherwise.

### GetVpcIdOk

`func (o *BatchInstanceCreateRequest) GetVpcIdOk() (*string, bool)`

GetVpcIdOk returns a tuple with the VpcId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpcId

`func (o *BatchInstanceCreateRequest) SetVpcId(v string)`

SetVpcId sets VpcId field to given value.


### GetUserData

`func (o *BatchInstanceCreateRequest) GetUserData() string`

GetUserData returns the UserData field if non-nil, zero value otherwise.

### GetUserDataOk

`func (o *BatchInstanceCreateRequest) GetUserDataOk() (*string, bool)`

GetUserDataOk returns a tuple with the UserData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserData

`func (o *BatchInstanceCreateRequest) SetUserData(v string)`

SetUserData sets UserData field to given value.

### HasUserData

`func (o *BatchInstanceCreateRequest) HasUserData() bool`

HasUserData returns a boolean if a field has been set.

### SetUserDataNil

`func (o *BatchInstanceCreateRequest) SetUserDataNil(b bool)`

 SetUserDataNil sets the value for UserData to be an explicit nil

### UnsetUserData
`func (o *BatchInstanceCreateRequest) UnsetUserData()`

UnsetUserData ensures that no value is present for UserData, not even an explicit nil
### GetOperatingSystemId

`func (o *BatchInstanceCreateRequest) GetOperatingSystemId() string`

GetOperatingSystemId returns the OperatingSystemId field if non-nil, zero value otherwise.

### GetOperatingSystemIdOk

`func (o *BatchInstanceCreateRequest) GetOperatingSystemIdOk() (*string, bool)`

GetOperatingSystemIdOk returns a tuple with the OperatingSystemId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperatingSystemId

`func (o *BatchInstanceCreateRequest) SetOperatingSystemId(v string)`

SetOperatingSystemId sets OperatingSystemId field to given value.

### HasOperatingSystemId

`func (o *BatchInstanceCreateRequest) HasOperatingSystemId() bool`

HasOperatingSystemId returns a boolean if a field has been set.

### SetOperatingSystemIdNil

`func (o *BatchInstanceCreateRequest) SetOperatingSystemIdNil(b bool)`

 SetOperatingSystemIdNil sets the value for OperatingSystemId to be an explicit nil

### UnsetOperatingSystemId
`func (o *BatchInstanceCreateRequest) UnsetOperatingSystemId()`

UnsetOperatingSystemId ensures that no value is present for OperatingSystemId, not even an explicit nil
### GetNetworkSecurityGroupId

`func (o *BatchInstanceCreateRequest) GetNetworkSecurityGroupId() string`

GetNetworkSecurityGroupId returns the NetworkSecurityGroupId field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupIdOk

`func (o *BatchInstanceCreateRequest) GetNetworkSecurityGroupIdOk() (*string, bool)`

GetNetworkSecurityGroupIdOk returns a tuple with the NetworkSecurityGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupId

`func (o *BatchInstanceCreateRequest) SetNetworkSecurityGroupId(v string)`

SetNetworkSecurityGroupId sets NetworkSecurityGroupId field to given value.

### HasNetworkSecurityGroupId

`func (o *BatchInstanceCreateRequest) HasNetworkSecurityGroupId() bool`

HasNetworkSecurityGroupId returns a boolean if a field has been set.

### SetNetworkSecurityGroupIdNil

`func (o *BatchInstanceCreateRequest) SetNetworkSecurityGroupIdNil(b bool)`

 SetNetworkSecurityGroupIdNil sets the value for NetworkSecurityGroupId to be an explicit nil

### UnsetNetworkSecurityGroupId
`func (o *BatchInstanceCreateRequest) UnsetNetworkSecurityGroupId()`

UnsetNetworkSecurityGroupId ensures that no value is present for NetworkSecurityGroupId, not even an explicit nil
### GetIpxeScript

`func (o *BatchInstanceCreateRequest) GetIpxeScript() string`

GetIpxeScript returns the IpxeScript field if non-nil, zero value otherwise.

### GetIpxeScriptOk

`func (o *BatchInstanceCreateRequest) GetIpxeScriptOk() (*string, bool)`

GetIpxeScriptOk returns a tuple with the IpxeScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpxeScript

`func (o *BatchInstanceCreateRequest) SetIpxeScript(v string)`

SetIpxeScript sets IpxeScript field to given value.

### HasIpxeScript

`func (o *BatchInstanceCreateRequest) HasIpxeScript() bool`

HasIpxeScript returns a boolean if a field has been set.

### SetIpxeScriptNil

`func (o *BatchInstanceCreateRequest) SetIpxeScriptNil(b bool)`

 SetIpxeScriptNil sets the value for IpxeScript to be an explicit nil

### UnsetIpxeScript
`func (o *BatchInstanceCreateRequest) UnsetIpxeScript()`

UnsetIpxeScript ensures that no value is present for IpxeScript, not even an explicit nil
### GetAlwaysBootWithCustomIpxe

`func (o *BatchInstanceCreateRequest) GetAlwaysBootWithCustomIpxe() bool`

GetAlwaysBootWithCustomIpxe returns the AlwaysBootWithCustomIpxe field if non-nil, zero value otherwise.

### GetAlwaysBootWithCustomIpxeOk

`func (o *BatchInstanceCreateRequest) GetAlwaysBootWithCustomIpxeOk() (*bool, bool)`

GetAlwaysBootWithCustomIpxeOk returns a tuple with the AlwaysBootWithCustomIpxe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlwaysBootWithCustomIpxe

`func (o *BatchInstanceCreateRequest) SetAlwaysBootWithCustomIpxe(v bool)`

SetAlwaysBootWithCustomIpxe sets AlwaysBootWithCustomIpxe field to given value.

### HasAlwaysBootWithCustomIpxe

`func (o *BatchInstanceCreateRequest) HasAlwaysBootWithCustomIpxe() bool`

HasAlwaysBootWithCustomIpxe returns a boolean if a field has been set.

### GetPhoneHomeEnabled

`func (o *BatchInstanceCreateRequest) GetPhoneHomeEnabled() bool`

GetPhoneHomeEnabled returns the PhoneHomeEnabled field if non-nil, zero value otherwise.

### GetPhoneHomeEnabledOk

`func (o *BatchInstanceCreateRequest) GetPhoneHomeEnabledOk() (*bool, bool)`

GetPhoneHomeEnabledOk returns a tuple with the PhoneHomeEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneHomeEnabled

`func (o *BatchInstanceCreateRequest) SetPhoneHomeEnabled(v bool)`

SetPhoneHomeEnabled sets PhoneHomeEnabled field to given value.

### HasPhoneHomeEnabled

`func (o *BatchInstanceCreateRequest) HasPhoneHomeEnabled() bool`

HasPhoneHomeEnabled returns a boolean if a field has been set.

### GetLabels

`func (o *BatchInstanceCreateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *BatchInstanceCreateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *BatchInstanceCreateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *BatchInstanceCreateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetInterfaces

`func (o *BatchInstanceCreateRequest) GetInterfaces() []InterfaceCreateRequest`

GetInterfaces returns the Interfaces field if non-nil, zero value otherwise.

### GetInterfacesOk

`func (o *BatchInstanceCreateRequest) GetInterfacesOk() (*[]InterfaceCreateRequest, bool)`

GetInterfacesOk returns a tuple with the Interfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterfaces

`func (o *BatchInstanceCreateRequest) SetInterfaces(v []InterfaceCreateRequest)`

SetInterfaces sets Interfaces field to given value.


### GetInfinibandInterfaces

`func (o *BatchInstanceCreateRequest) GetInfinibandInterfaces() []InfiniBandInterfaceCreateRequest`

GetInfinibandInterfaces returns the InfinibandInterfaces field if non-nil, zero value otherwise.

### GetInfinibandInterfacesOk

`func (o *BatchInstanceCreateRequest) GetInfinibandInterfacesOk() (*[]InfiniBandInterfaceCreateRequest, bool)`

GetInfinibandInterfacesOk returns a tuple with the InfinibandInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfinibandInterfaces

`func (o *BatchInstanceCreateRequest) SetInfinibandInterfaces(v []InfiniBandInterfaceCreateRequest)`

SetInfinibandInterfaces sets InfinibandInterfaces field to given value.

### HasInfinibandInterfaces

`func (o *BatchInstanceCreateRequest) HasInfinibandInterfaces() bool`

HasInfinibandInterfaces returns a boolean if a field has been set.

### GetDpuExtensionServiceDeployments

`func (o *BatchInstanceCreateRequest) GetDpuExtensionServiceDeployments() []DpuExtensionServiceDeploymentRequest`

GetDpuExtensionServiceDeployments returns the DpuExtensionServiceDeployments field if non-nil, zero value otherwise.

### GetDpuExtensionServiceDeploymentsOk

`func (o *BatchInstanceCreateRequest) GetDpuExtensionServiceDeploymentsOk() (*[]DpuExtensionServiceDeploymentRequest, bool)`

GetDpuExtensionServiceDeploymentsOk returns a tuple with the DpuExtensionServiceDeployments field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDpuExtensionServiceDeployments

`func (o *BatchInstanceCreateRequest) SetDpuExtensionServiceDeployments(v []DpuExtensionServiceDeploymentRequest)`

SetDpuExtensionServiceDeployments sets DpuExtensionServiceDeployments field to given value.

### HasDpuExtensionServiceDeployments

`func (o *BatchInstanceCreateRequest) HasDpuExtensionServiceDeployments() bool`

HasDpuExtensionServiceDeployments returns a boolean if a field has been set.

### GetNvLinkInterfaces

`func (o *BatchInstanceCreateRequest) GetNvLinkInterfaces() []NVLinkInterfaceCreateRequest`

GetNvLinkInterfaces returns the NvLinkInterfaces field if non-nil, zero value otherwise.

### GetNvLinkInterfacesOk

`func (o *BatchInstanceCreateRequest) GetNvLinkInterfacesOk() (*[]NVLinkInterfaceCreateRequest, bool)`

GetNvLinkInterfacesOk returns a tuple with the NvLinkInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinkInterfaces

`func (o *BatchInstanceCreateRequest) SetNvLinkInterfaces(v []NVLinkInterfaceCreateRequest)`

SetNvLinkInterfaces sets NvLinkInterfaces field to given value.

### HasNvLinkInterfaces

`func (o *BatchInstanceCreateRequest) HasNvLinkInterfaces() bool`

HasNvLinkInterfaces returns a boolean if a field has been set.

### GetSshKeyGroupIds

`func (o *BatchInstanceCreateRequest) GetSshKeyGroupIds() []string`

GetSshKeyGroupIds returns the SshKeyGroupIds field if non-nil, zero value otherwise.

### GetSshKeyGroupIdsOk

`func (o *BatchInstanceCreateRequest) GetSshKeyGroupIdsOk() (*[]string, bool)`

GetSshKeyGroupIdsOk returns a tuple with the SshKeyGroupIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyGroupIds

`func (o *BatchInstanceCreateRequest) SetSshKeyGroupIds(v []string)`

SetSshKeyGroupIds sets SshKeyGroupIds field to given value.

### HasSshKeyGroupIds

`func (o *BatchInstanceCreateRequest) HasSshKeyGroupIds() bool`

HasSshKeyGroupIds returns a boolean if a field has been set.

### GetTopologyOptimized

`func (o *BatchInstanceCreateRequest) GetTopologyOptimized() bool`

GetTopologyOptimized returns the TopologyOptimized field if non-nil, zero value otherwise.

### GetTopologyOptimizedOk

`func (o *BatchInstanceCreateRequest) GetTopologyOptimizedOk() (*bool, bool)`

GetTopologyOptimizedOk returns a tuple with the TopologyOptimized field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTopologyOptimized

`func (o *BatchInstanceCreateRequest) SetTopologyOptimized(v bool)`

SetTopologyOptimized sets TopologyOptimized field to given value.

### HasTopologyOptimized

`func (o *BatchInstanceCreateRequest) HasTopologyOptimized() bool`

HasTopologyOptimized returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


