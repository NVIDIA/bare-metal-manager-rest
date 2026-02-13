# InstanceCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name of the Instance | 
**Description** | Pointer to **NullableString** | Description of the Instance, optional | [optional] 
**TenantId** | **string** | ID of the Tenant creating the Instance | 
**InstanceTypeId** | Pointer to **string** | ID of the Instance Type to use for Instance | [optional] 
**MachineId** | Pointer to **string** | ID of of specific Machine to use for Instance. Requires Targeted Instance Creation capability enabled for Tenant | [optional] 
**VpcId** | **string** | ID of the VPC the Instance should belong to | 
**UserData** | Pointer to **NullableString** | Can only be specified if allowOverride is set to true in Operating System | [optional] 
**OperatingSystemId** | Pointer to **NullableString** | Must be specified if iPXE Script field is empty | [optional] 
**NetworkSecurityGroupId** | Pointer to **NullableString** |  | [optional] 
**IpxeScript** | Pointer to **NullableString** | Override iPXE script specified in OS, must be specified if Operating System is not specified | [optional] 
**AlwaysBootWithCustomIpxe** | Pointer to **bool** | When set to true, the iPXE script specified by OS or overridden here will always be run when rebooting the Instance. OS must be of iPXE type. | [optional] 
**PhoneHomeEnabled** | Pointer to **bool** | When set to true, the Instance will be enabled with the Phone Home service. | [optional] 
**Labels** | Pointer to **map[string]string** |  | [optional] 
**Interfaces** | [**[]InterfaceCreateRequest**](InterfaceCreateRequest.md) | At least one interface must be specified. Either Subnet or VPC Prefix interfaces allowed. Only one of the Subnets or VPC Prefixes can be attached over Physical interface. If only one Subnet is specified, then it will be attached over physical interface regardless of the value of isPhysical. In case of VPC Prefix, isPhysical will always be true | 
**InfinibandInterfaces** | Pointer to [**[]InfiniBandInterfaceCreateRequest**](InfiniBandInterfaceCreateRequest.md) | Associate one or more Partitions with this Instance | [optional] 
**DpuExtensionServiceDeployments** | Pointer to [**[]DpuExtensionServiceDeploymentRequest**](DpuExtensionServiceDeploymentRequest.md) | DPU Extension Services to deploy to the DPUs of this Instance | [optional] 
**NvLinkInterfaces** | Pointer to [**[]NVLinkInterfaceCreateRequest**](NVLinkInterfaceCreateRequest.md) | Associate one or more NVLink Logical Partitions with this Instance | [optional] 
**SshKeyGroupIds** | Pointer to **[]string** | Specify list of SSH Key Group IDs that will provide Serial over LAN access | [optional] 
**AllowUnhealthyMachine** | Pointer to **bool** | Must be set to true creating a targeted Instance with a Machine that is in Error status. Requires Targeted Instance Creation capability enabled for Tenant | [optional] 

## Methods

### NewInstanceCreateRequest

`func NewInstanceCreateRequest(name string, tenantId string, vpcId string, interfaces []InterfaceCreateRequest, ) *InstanceCreateRequest`

NewInstanceCreateRequest instantiates a new InstanceCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceCreateRequestWithDefaults

`func NewInstanceCreateRequestWithDefaults() *InstanceCreateRequest`

NewInstanceCreateRequestWithDefaults instantiates a new InstanceCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *InstanceCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *InstanceCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *InstanceCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *InstanceCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *InstanceCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *InstanceCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *InstanceCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *InstanceCreateRequest) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *InstanceCreateRequest) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetTenantId

`func (o *InstanceCreateRequest) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *InstanceCreateRequest) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *InstanceCreateRequest) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.


### GetInstanceTypeId

`func (o *InstanceCreateRequest) GetInstanceTypeId() string`

GetInstanceTypeId returns the InstanceTypeId field if non-nil, zero value otherwise.

### GetInstanceTypeIdOk

`func (o *InstanceCreateRequest) GetInstanceTypeIdOk() (*string, bool)`

GetInstanceTypeIdOk returns a tuple with the InstanceTypeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceTypeId

`func (o *InstanceCreateRequest) SetInstanceTypeId(v string)`

SetInstanceTypeId sets InstanceTypeId field to given value.

### HasInstanceTypeId

`func (o *InstanceCreateRequest) HasInstanceTypeId() bool`

HasInstanceTypeId returns a boolean if a field has been set.

### GetMachineId

`func (o *InstanceCreateRequest) GetMachineId() string`

GetMachineId returns the MachineId field if non-nil, zero value otherwise.

### GetMachineIdOk

`func (o *InstanceCreateRequest) GetMachineIdOk() (*string, bool)`

GetMachineIdOk returns a tuple with the MachineId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineId

`func (o *InstanceCreateRequest) SetMachineId(v string)`

SetMachineId sets MachineId field to given value.

### HasMachineId

`func (o *InstanceCreateRequest) HasMachineId() bool`

HasMachineId returns a boolean if a field has been set.

### GetVpcId

`func (o *InstanceCreateRequest) GetVpcId() string`

GetVpcId returns the VpcId field if non-nil, zero value otherwise.

### GetVpcIdOk

`func (o *InstanceCreateRequest) GetVpcIdOk() (*string, bool)`

GetVpcIdOk returns a tuple with the VpcId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpcId

`func (o *InstanceCreateRequest) SetVpcId(v string)`

SetVpcId sets VpcId field to given value.


### GetUserData

`func (o *InstanceCreateRequest) GetUserData() string`

GetUserData returns the UserData field if non-nil, zero value otherwise.

### GetUserDataOk

`func (o *InstanceCreateRequest) GetUserDataOk() (*string, bool)`

GetUserDataOk returns a tuple with the UserData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserData

`func (o *InstanceCreateRequest) SetUserData(v string)`

SetUserData sets UserData field to given value.

### HasUserData

`func (o *InstanceCreateRequest) HasUserData() bool`

HasUserData returns a boolean if a field has been set.

### SetUserDataNil

`func (o *InstanceCreateRequest) SetUserDataNil(b bool)`

 SetUserDataNil sets the value for UserData to be an explicit nil

### UnsetUserData
`func (o *InstanceCreateRequest) UnsetUserData()`

UnsetUserData ensures that no value is present for UserData, not even an explicit nil
### GetOperatingSystemId

`func (o *InstanceCreateRequest) GetOperatingSystemId() string`

GetOperatingSystemId returns the OperatingSystemId field if non-nil, zero value otherwise.

### GetOperatingSystemIdOk

`func (o *InstanceCreateRequest) GetOperatingSystemIdOk() (*string, bool)`

GetOperatingSystemIdOk returns a tuple with the OperatingSystemId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperatingSystemId

`func (o *InstanceCreateRequest) SetOperatingSystemId(v string)`

SetOperatingSystemId sets OperatingSystemId field to given value.

### HasOperatingSystemId

`func (o *InstanceCreateRequest) HasOperatingSystemId() bool`

HasOperatingSystemId returns a boolean if a field has been set.

### SetOperatingSystemIdNil

`func (o *InstanceCreateRequest) SetOperatingSystemIdNil(b bool)`

 SetOperatingSystemIdNil sets the value for OperatingSystemId to be an explicit nil

### UnsetOperatingSystemId
`func (o *InstanceCreateRequest) UnsetOperatingSystemId()`

UnsetOperatingSystemId ensures that no value is present for OperatingSystemId, not even an explicit nil
### GetNetworkSecurityGroupId

`func (o *InstanceCreateRequest) GetNetworkSecurityGroupId() string`

GetNetworkSecurityGroupId returns the NetworkSecurityGroupId field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupIdOk

`func (o *InstanceCreateRequest) GetNetworkSecurityGroupIdOk() (*string, bool)`

GetNetworkSecurityGroupIdOk returns a tuple with the NetworkSecurityGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupId

`func (o *InstanceCreateRequest) SetNetworkSecurityGroupId(v string)`

SetNetworkSecurityGroupId sets NetworkSecurityGroupId field to given value.

### HasNetworkSecurityGroupId

`func (o *InstanceCreateRequest) HasNetworkSecurityGroupId() bool`

HasNetworkSecurityGroupId returns a boolean if a field has been set.

### SetNetworkSecurityGroupIdNil

`func (o *InstanceCreateRequest) SetNetworkSecurityGroupIdNil(b bool)`

 SetNetworkSecurityGroupIdNil sets the value for NetworkSecurityGroupId to be an explicit nil

### UnsetNetworkSecurityGroupId
`func (o *InstanceCreateRequest) UnsetNetworkSecurityGroupId()`

UnsetNetworkSecurityGroupId ensures that no value is present for NetworkSecurityGroupId, not even an explicit nil
### GetIpxeScript

`func (o *InstanceCreateRequest) GetIpxeScript() string`

GetIpxeScript returns the IpxeScript field if non-nil, zero value otherwise.

### GetIpxeScriptOk

`func (o *InstanceCreateRequest) GetIpxeScriptOk() (*string, bool)`

GetIpxeScriptOk returns a tuple with the IpxeScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpxeScript

`func (o *InstanceCreateRequest) SetIpxeScript(v string)`

SetIpxeScript sets IpxeScript field to given value.

### HasIpxeScript

`func (o *InstanceCreateRequest) HasIpxeScript() bool`

HasIpxeScript returns a boolean if a field has been set.

### SetIpxeScriptNil

`func (o *InstanceCreateRequest) SetIpxeScriptNil(b bool)`

 SetIpxeScriptNil sets the value for IpxeScript to be an explicit nil

### UnsetIpxeScript
`func (o *InstanceCreateRequest) UnsetIpxeScript()`

UnsetIpxeScript ensures that no value is present for IpxeScript, not even an explicit nil
### GetAlwaysBootWithCustomIpxe

`func (o *InstanceCreateRequest) GetAlwaysBootWithCustomIpxe() bool`

GetAlwaysBootWithCustomIpxe returns the AlwaysBootWithCustomIpxe field if non-nil, zero value otherwise.

### GetAlwaysBootWithCustomIpxeOk

`func (o *InstanceCreateRequest) GetAlwaysBootWithCustomIpxeOk() (*bool, bool)`

GetAlwaysBootWithCustomIpxeOk returns a tuple with the AlwaysBootWithCustomIpxe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlwaysBootWithCustomIpxe

`func (o *InstanceCreateRequest) SetAlwaysBootWithCustomIpxe(v bool)`

SetAlwaysBootWithCustomIpxe sets AlwaysBootWithCustomIpxe field to given value.

### HasAlwaysBootWithCustomIpxe

`func (o *InstanceCreateRequest) HasAlwaysBootWithCustomIpxe() bool`

HasAlwaysBootWithCustomIpxe returns a boolean if a field has been set.

### GetPhoneHomeEnabled

`func (o *InstanceCreateRequest) GetPhoneHomeEnabled() bool`

GetPhoneHomeEnabled returns the PhoneHomeEnabled field if non-nil, zero value otherwise.

### GetPhoneHomeEnabledOk

`func (o *InstanceCreateRequest) GetPhoneHomeEnabledOk() (*bool, bool)`

GetPhoneHomeEnabledOk returns a tuple with the PhoneHomeEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneHomeEnabled

`func (o *InstanceCreateRequest) SetPhoneHomeEnabled(v bool)`

SetPhoneHomeEnabled sets PhoneHomeEnabled field to given value.

### HasPhoneHomeEnabled

`func (o *InstanceCreateRequest) HasPhoneHomeEnabled() bool`

HasPhoneHomeEnabled returns a boolean if a field has been set.

### GetLabels

`func (o *InstanceCreateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *InstanceCreateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *InstanceCreateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *InstanceCreateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetInterfaces

`func (o *InstanceCreateRequest) GetInterfaces() []InterfaceCreateRequest`

GetInterfaces returns the Interfaces field if non-nil, zero value otherwise.

### GetInterfacesOk

`func (o *InstanceCreateRequest) GetInterfacesOk() (*[]InterfaceCreateRequest, bool)`

GetInterfacesOk returns a tuple with the Interfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterfaces

`func (o *InstanceCreateRequest) SetInterfaces(v []InterfaceCreateRequest)`

SetInterfaces sets Interfaces field to given value.


### GetInfinibandInterfaces

`func (o *InstanceCreateRequest) GetInfinibandInterfaces() []InfiniBandInterfaceCreateRequest`

GetInfinibandInterfaces returns the InfinibandInterfaces field if non-nil, zero value otherwise.

### GetInfinibandInterfacesOk

`func (o *InstanceCreateRequest) GetInfinibandInterfacesOk() (*[]InfiniBandInterfaceCreateRequest, bool)`

GetInfinibandInterfacesOk returns a tuple with the InfinibandInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfinibandInterfaces

`func (o *InstanceCreateRequest) SetInfinibandInterfaces(v []InfiniBandInterfaceCreateRequest)`

SetInfinibandInterfaces sets InfinibandInterfaces field to given value.

### HasInfinibandInterfaces

`func (o *InstanceCreateRequest) HasInfinibandInterfaces() bool`

HasInfinibandInterfaces returns a boolean if a field has been set.

### GetDpuExtensionServiceDeployments

`func (o *InstanceCreateRequest) GetDpuExtensionServiceDeployments() []DpuExtensionServiceDeploymentRequest`

GetDpuExtensionServiceDeployments returns the DpuExtensionServiceDeployments field if non-nil, zero value otherwise.

### GetDpuExtensionServiceDeploymentsOk

`func (o *InstanceCreateRequest) GetDpuExtensionServiceDeploymentsOk() (*[]DpuExtensionServiceDeploymentRequest, bool)`

GetDpuExtensionServiceDeploymentsOk returns a tuple with the DpuExtensionServiceDeployments field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDpuExtensionServiceDeployments

`func (o *InstanceCreateRequest) SetDpuExtensionServiceDeployments(v []DpuExtensionServiceDeploymentRequest)`

SetDpuExtensionServiceDeployments sets DpuExtensionServiceDeployments field to given value.

### HasDpuExtensionServiceDeployments

`func (o *InstanceCreateRequest) HasDpuExtensionServiceDeployments() bool`

HasDpuExtensionServiceDeployments returns a boolean if a field has been set.

### GetNvLinkInterfaces

`func (o *InstanceCreateRequest) GetNvLinkInterfaces() []NVLinkInterfaceCreateRequest`

GetNvLinkInterfaces returns the NvLinkInterfaces field if non-nil, zero value otherwise.

### GetNvLinkInterfacesOk

`func (o *InstanceCreateRequest) GetNvLinkInterfacesOk() (*[]NVLinkInterfaceCreateRequest, bool)`

GetNvLinkInterfacesOk returns a tuple with the NvLinkInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinkInterfaces

`func (o *InstanceCreateRequest) SetNvLinkInterfaces(v []NVLinkInterfaceCreateRequest)`

SetNvLinkInterfaces sets NvLinkInterfaces field to given value.

### HasNvLinkInterfaces

`func (o *InstanceCreateRequest) HasNvLinkInterfaces() bool`

HasNvLinkInterfaces returns a boolean if a field has been set.

### GetSshKeyGroupIds

`func (o *InstanceCreateRequest) GetSshKeyGroupIds() []string`

GetSshKeyGroupIds returns the SshKeyGroupIds field if non-nil, zero value otherwise.

### GetSshKeyGroupIdsOk

`func (o *InstanceCreateRequest) GetSshKeyGroupIdsOk() (*[]string, bool)`

GetSshKeyGroupIdsOk returns a tuple with the SshKeyGroupIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyGroupIds

`func (o *InstanceCreateRequest) SetSshKeyGroupIds(v []string)`

SetSshKeyGroupIds sets SshKeyGroupIds field to given value.

### HasSshKeyGroupIds

`func (o *InstanceCreateRequest) HasSshKeyGroupIds() bool`

HasSshKeyGroupIds returns a boolean if a field has been set.

### GetAllowUnhealthyMachine

`func (o *InstanceCreateRequest) GetAllowUnhealthyMachine() bool`

GetAllowUnhealthyMachine returns the AllowUnhealthyMachine field if non-nil, zero value otherwise.

### GetAllowUnhealthyMachineOk

`func (o *InstanceCreateRequest) GetAllowUnhealthyMachineOk() (*bool, bool)`

GetAllowUnhealthyMachineOk returns a tuple with the AllowUnhealthyMachine field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowUnhealthyMachine

`func (o *InstanceCreateRequest) SetAllowUnhealthyMachine(v bool)`

SetAllowUnhealthyMachine sets AllowUnhealthyMachine field to given value.

### HasAllowUnhealthyMachine

`func (o *InstanceCreateRequest) HasAllowUnhealthyMachine() bool`

HasAllowUnhealthyMachine returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


