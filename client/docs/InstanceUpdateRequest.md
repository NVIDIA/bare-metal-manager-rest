# InstanceUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **NullableString** | Updated name for the Instance | [optional] 
**Description** | Pointer to **NullableString** | Updated description of the Instance | [optional] 
**TriggerReboot** | Pointer to **NullableBool** | Trigger power cycle for Instance | [optional] 
**RebootWithCustomIpxe** | Pointer to **NullableBool** | When specified along with triggerReboot, the Instance will boot using the custom iPXE specified by OS. If Instance has alwaysBootWithCustomIpxe flag set then this value will be ignored. | [optional] 
**ApplyUpdatesOnReboot** | Pointer to **NullableBool** | When specified, any updates pending for the Instance e.g. DPU reprovisioning, will be applied on reboot | [optional] 
**OperatingSystemId** | Pointer to **NullableString** | The UUID of the desired operating system. | [optional] 
**IpxeScript** | Pointer to **NullableString** | The iPXE script content to be used for booting. | [optional] 
**SshKeyGroupIds** | Pointer to **[]string** | Specify a new list of SSH Key Group IDs that will provide Serial over LAN and SSH access.  This will overwrite an existing list. | [optional] 
**NetworkSecurityGroupId** | Pointer to **NullableString** |  | [optional] 
**UserData** | Pointer to **NullableString** | Any user-data to be sent to the booting OS.  For example, cloud-init data. | [optional] 
**AlwaysBootWithCustomIpxe** | Pointer to **NullableBool** | Whether the custom iPXE data should be used for every boot. | [optional] 
**PhoneHomeEnabled** | Pointer to **NullableBool** | Indicates whether the Phone Home service should be enabled or disabled for Instance | [optional] 
**Labels** | Pointer to **map[string]string** | Update labels of the Instance. The labels will be entirely replaced by those sent in the request. Any labels not included in the request will be removed. To retain existing labels, first fetch them and include them along with this request. | [optional] 
**Interfaces** | Pointer to [**[]InterfaceCreateRequest**](InterfaceCreateRequest.md) | Update Interfaces of the Instance | [optional] 
**InfinibandInterfaces** | Pointer to [**[]InfiniBandInterfaceCreateRequest**](InfiniBandInterfaceCreateRequest.md) | Update InfiniBand Interfaces of the Instance | [optional] 
**NvLinkInterfaces** | Pointer to [**[]NVLinkInterfaceCreateRequest**](NVLinkInterfaceCreateRequest.md) | Update NVLink Interfaces of the Instance | [optional] 
**DpuExtensionServiceDeployments** | Pointer to [**[]DpuExtensionServiceDeploymentRequest**](DpuExtensionServiceDeploymentRequest.md) | Updated set of DPU Extension Services to deploy to the DPUs of this Instance | [optional] 

## Methods

### NewInstanceUpdateRequest

`func NewInstanceUpdateRequest() *InstanceUpdateRequest`

NewInstanceUpdateRequest instantiates a new InstanceUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceUpdateRequestWithDefaults

`func NewInstanceUpdateRequestWithDefaults() *InstanceUpdateRequest`

NewInstanceUpdateRequestWithDefaults instantiates a new InstanceUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *InstanceUpdateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *InstanceUpdateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *InstanceUpdateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *InstanceUpdateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### SetNameNil

`func (o *InstanceUpdateRequest) SetNameNil(b bool)`

 SetNameNil sets the value for Name to be an explicit nil

### UnsetName
`func (o *InstanceUpdateRequest) UnsetName()`

UnsetName ensures that no value is present for Name, not even an explicit nil
### GetDescription

`func (o *InstanceUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *InstanceUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *InstanceUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *InstanceUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *InstanceUpdateRequest) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *InstanceUpdateRequest) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetTriggerReboot

`func (o *InstanceUpdateRequest) GetTriggerReboot() bool`

GetTriggerReboot returns the TriggerReboot field if non-nil, zero value otherwise.

### GetTriggerRebootOk

`func (o *InstanceUpdateRequest) GetTriggerRebootOk() (*bool, bool)`

GetTriggerRebootOk returns a tuple with the TriggerReboot field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTriggerReboot

`func (o *InstanceUpdateRequest) SetTriggerReboot(v bool)`

SetTriggerReboot sets TriggerReboot field to given value.

### HasTriggerReboot

`func (o *InstanceUpdateRequest) HasTriggerReboot() bool`

HasTriggerReboot returns a boolean if a field has been set.

### SetTriggerRebootNil

`func (o *InstanceUpdateRequest) SetTriggerRebootNil(b bool)`

 SetTriggerRebootNil sets the value for TriggerReboot to be an explicit nil

### UnsetTriggerReboot
`func (o *InstanceUpdateRequest) UnsetTriggerReboot()`

UnsetTriggerReboot ensures that no value is present for TriggerReboot, not even an explicit nil
### GetRebootWithCustomIpxe

`func (o *InstanceUpdateRequest) GetRebootWithCustomIpxe() bool`

GetRebootWithCustomIpxe returns the RebootWithCustomIpxe field if non-nil, zero value otherwise.

### GetRebootWithCustomIpxeOk

`func (o *InstanceUpdateRequest) GetRebootWithCustomIpxeOk() (*bool, bool)`

GetRebootWithCustomIpxeOk returns a tuple with the RebootWithCustomIpxe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRebootWithCustomIpxe

`func (o *InstanceUpdateRequest) SetRebootWithCustomIpxe(v bool)`

SetRebootWithCustomIpxe sets RebootWithCustomIpxe field to given value.

### HasRebootWithCustomIpxe

`func (o *InstanceUpdateRequest) HasRebootWithCustomIpxe() bool`

HasRebootWithCustomIpxe returns a boolean if a field has been set.

### SetRebootWithCustomIpxeNil

`func (o *InstanceUpdateRequest) SetRebootWithCustomIpxeNil(b bool)`

 SetRebootWithCustomIpxeNil sets the value for RebootWithCustomIpxe to be an explicit nil

### UnsetRebootWithCustomIpxe
`func (o *InstanceUpdateRequest) UnsetRebootWithCustomIpxe()`

UnsetRebootWithCustomIpxe ensures that no value is present for RebootWithCustomIpxe, not even an explicit nil
### GetApplyUpdatesOnReboot

`func (o *InstanceUpdateRequest) GetApplyUpdatesOnReboot() bool`

GetApplyUpdatesOnReboot returns the ApplyUpdatesOnReboot field if non-nil, zero value otherwise.

### GetApplyUpdatesOnRebootOk

`func (o *InstanceUpdateRequest) GetApplyUpdatesOnRebootOk() (*bool, bool)`

GetApplyUpdatesOnRebootOk returns a tuple with the ApplyUpdatesOnReboot field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApplyUpdatesOnReboot

`func (o *InstanceUpdateRequest) SetApplyUpdatesOnReboot(v bool)`

SetApplyUpdatesOnReboot sets ApplyUpdatesOnReboot field to given value.

### HasApplyUpdatesOnReboot

`func (o *InstanceUpdateRequest) HasApplyUpdatesOnReboot() bool`

HasApplyUpdatesOnReboot returns a boolean if a field has been set.

### SetApplyUpdatesOnRebootNil

`func (o *InstanceUpdateRequest) SetApplyUpdatesOnRebootNil(b bool)`

 SetApplyUpdatesOnRebootNil sets the value for ApplyUpdatesOnReboot to be an explicit nil

### UnsetApplyUpdatesOnReboot
`func (o *InstanceUpdateRequest) UnsetApplyUpdatesOnReboot()`

UnsetApplyUpdatesOnReboot ensures that no value is present for ApplyUpdatesOnReboot, not even an explicit nil
### GetOperatingSystemId

`func (o *InstanceUpdateRequest) GetOperatingSystemId() string`

GetOperatingSystemId returns the OperatingSystemId field if non-nil, zero value otherwise.

### GetOperatingSystemIdOk

`func (o *InstanceUpdateRequest) GetOperatingSystemIdOk() (*string, bool)`

GetOperatingSystemIdOk returns a tuple with the OperatingSystemId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperatingSystemId

`func (o *InstanceUpdateRequest) SetOperatingSystemId(v string)`

SetOperatingSystemId sets OperatingSystemId field to given value.

### HasOperatingSystemId

`func (o *InstanceUpdateRequest) HasOperatingSystemId() bool`

HasOperatingSystemId returns a boolean if a field has been set.

### SetOperatingSystemIdNil

`func (o *InstanceUpdateRequest) SetOperatingSystemIdNil(b bool)`

 SetOperatingSystemIdNil sets the value for OperatingSystemId to be an explicit nil

### UnsetOperatingSystemId
`func (o *InstanceUpdateRequest) UnsetOperatingSystemId()`

UnsetOperatingSystemId ensures that no value is present for OperatingSystemId, not even an explicit nil
### GetIpxeScript

`func (o *InstanceUpdateRequest) GetIpxeScript() string`

GetIpxeScript returns the IpxeScript field if non-nil, zero value otherwise.

### GetIpxeScriptOk

`func (o *InstanceUpdateRequest) GetIpxeScriptOk() (*string, bool)`

GetIpxeScriptOk returns a tuple with the IpxeScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpxeScript

`func (o *InstanceUpdateRequest) SetIpxeScript(v string)`

SetIpxeScript sets IpxeScript field to given value.

### HasIpxeScript

`func (o *InstanceUpdateRequest) HasIpxeScript() bool`

HasIpxeScript returns a boolean if a field has been set.

### SetIpxeScriptNil

`func (o *InstanceUpdateRequest) SetIpxeScriptNil(b bool)`

 SetIpxeScriptNil sets the value for IpxeScript to be an explicit nil

### UnsetIpxeScript
`func (o *InstanceUpdateRequest) UnsetIpxeScript()`

UnsetIpxeScript ensures that no value is present for IpxeScript, not even an explicit nil
### GetSshKeyGroupIds

`func (o *InstanceUpdateRequest) GetSshKeyGroupIds() []string`

GetSshKeyGroupIds returns the SshKeyGroupIds field if non-nil, zero value otherwise.

### GetSshKeyGroupIdsOk

`func (o *InstanceUpdateRequest) GetSshKeyGroupIdsOk() (*[]string, bool)`

GetSshKeyGroupIdsOk returns a tuple with the SshKeyGroupIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyGroupIds

`func (o *InstanceUpdateRequest) SetSshKeyGroupIds(v []string)`

SetSshKeyGroupIds sets SshKeyGroupIds field to given value.

### HasSshKeyGroupIds

`func (o *InstanceUpdateRequest) HasSshKeyGroupIds() bool`

HasSshKeyGroupIds returns a boolean if a field has been set.

### GetNetworkSecurityGroupId

`func (o *InstanceUpdateRequest) GetNetworkSecurityGroupId() string`

GetNetworkSecurityGroupId returns the NetworkSecurityGroupId field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupIdOk

`func (o *InstanceUpdateRequest) GetNetworkSecurityGroupIdOk() (*string, bool)`

GetNetworkSecurityGroupIdOk returns a tuple with the NetworkSecurityGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupId

`func (o *InstanceUpdateRequest) SetNetworkSecurityGroupId(v string)`

SetNetworkSecurityGroupId sets NetworkSecurityGroupId field to given value.

### HasNetworkSecurityGroupId

`func (o *InstanceUpdateRequest) HasNetworkSecurityGroupId() bool`

HasNetworkSecurityGroupId returns a boolean if a field has been set.

### SetNetworkSecurityGroupIdNil

`func (o *InstanceUpdateRequest) SetNetworkSecurityGroupIdNil(b bool)`

 SetNetworkSecurityGroupIdNil sets the value for NetworkSecurityGroupId to be an explicit nil

### UnsetNetworkSecurityGroupId
`func (o *InstanceUpdateRequest) UnsetNetworkSecurityGroupId()`

UnsetNetworkSecurityGroupId ensures that no value is present for NetworkSecurityGroupId, not even an explicit nil
### GetUserData

`func (o *InstanceUpdateRequest) GetUserData() string`

GetUserData returns the UserData field if non-nil, zero value otherwise.

### GetUserDataOk

`func (o *InstanceUpdateRequest) GetUserDataOk() (*string, bool)`

GetUserDataOk returns a tuple with the UserData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserData

`func (o *InstanceUpdateRequest) SetUserData(v string)`

SetUserData sets UserData field to given value.

### HasUserData

`func (o *InstanceUpdateRequest) HasUserData() bool`

HasUserData returns a boolean if a field has been set.

### SetUserDataNil

`func (o *InstanceUpdateRequest) SetUserDataNil(b bool)`

 SetUserDataNil sets the value for UserData to be an explicit nil

### UnsetUserData
`func (o *InstanceUpdateRequest) UnsetUserData()`

UnsetUserData ensures that no value is present for UserData, not even an explicit nil
### GetAlwaysBootWithCustomIpxe

`func (o *InstanceUpdateRequest) GetAlwaysBootWithCustomIpxe() bool`

GetAlwaysBootWithCustomIpxe returns the AlwaysBootWithCustomIpxe field if non-nil, zero value otherwise.

### GetAlwaysBootWithCustomIpxeOk

`func (o *InstanceUpdateRequest) GetAlwaysBootWithCustomIpxeOk() (*bool, bool)`

GetAlwaysBootWithCustomIpxeOk returns a tuple with the AlwaysBootWithCustomIpxe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlwaysBootWithCustomIpxe

`func (o *InstanceUpdateRequest) SetAlwaysBootWithCustomIpxe(v bool)`

SetAlwaysBootWithCustomIpxe sets AlwaysBootWithCustomIpxe field to given value.

### HasAlwaysBootWithCustomIpxe

`func (o *InstanceUpdateRequest) HasAlwaysBootWithCustomIpxe() bool`

HasAlwaysBootWithCustomIpxe returns a boolean if a field has been set.

### SetAlwaysBootWithCustomIpxeNil

`func (o *InstanceUpdateRequest) SetAlwaysBootWithCustomIpxeNil(b bool)`

 SetAlwaysBootWithCustomIpxeNil sets the value for AlwaysBootWithCustomIpxe to be an explicit nil

### UnsetAlwaysBootWithCustomIpxe
`func (o *InstanceUpdateRequest) UnsetAlwaysBootWithCustomIpxe()`

UnsetAlwaysBootWithCustomIpxe ensures that no value is present for AlwaysBootWithCustomIpxe, not even an explicit nil
### GetPhoneHomeEnabled

`func (o *InstanceUpdateRequest) GetPhoneHomeEnabled() bool`

GetPhoneHomeEnabled returns the PhoneHomeEnabled field if non-nil, zero value otherwise.

### GetPhoneHomeEnabledOk

`func (o *InstanceUpdateRequest) GetPhoneHomeEnabledOk() (*bool, bool)`

GetPhoneHomeEnabledOk returns a tuple with the PhoneHomeEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneHomeEnabled

`func (o *InstanceUpdateRequest) SetPhoneHomeEnabled(v bool)`

SetPhoneHomeEnabled sets PhoneHomeEnabled field to given value.

### HasPhoneHomeEnabled

`func (o *InstanceUpdateRequest) HasPhoneHomeEnabled() bool`

HasPhoneHomeEnabled returns a boolean if a field has been set.

### SetPhoneHomeEnabledNil

`func (o *InstanceUpdateRequest) SetPhoneHomeEnabledNil(b bool)`

 SetPhoneHomeEnabledNil sets the value for PhoneHomeEnabled to be an explicit nil

### UnsetPhoneHomeEnabled
`func (o *InstanceUpdateRequest) UnsetPhoneHomeEnabled()`

UnsetPhoneHomeEnabled ensures that no value is present for PhoneHomeEnabled, not even an explicit nil
### GetLabels

`func (o *InstanceUpdateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *InstanceUpdateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *InstanceUpdateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *InstanceUpdateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetInterfaces

`func (o *InstanceUpdateRequest) GetInterfaces() []InterfaceCreateRequest`

GetInterfaces returns the Interfaces field if non-nil, zero value otherwise.

### GetInterfacesOk

`func (o *InstanceUpdateRequest) GetInterfacesOk() (*[]InterfaceCreateRequest, bool)`

GetInterfacesOk returns a tuple with the Interfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterfaces

`func (o *InstanceUpdateRequest) SetInterfaces(v []InterfaceCreateRequest)`

SetInterfaces sets Interfaces field to given value.

### HasInterfaces

`func (o *InstanceUpdateRequest) HasInterfaces() bool`

HasInterfaces returns a boolean if a field has been set.

### GetInfinibandInterfaces

`func (o *InstanceUpdateRequest) GetInfinibandInterfaces() []InfiniBandInterfaceCreateRequest`

GetInfinibandInterfaces returns the InfinibandInterfaces field if non-nil, zero value otherwise.

### GetInfinibandInterfacesOk

`func (o *InstanceUpdateRequest) GetInfinibandInterfacesOk() (*[]InfiniBandInterfaceCreateRequest, bool)`

GetInfinibandInterfacesOk returns a tuple with the InfinibandInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfinibandInterfaces

`func (o *InstanceUpdateRequest) SetInfinibandInterfaces(v []InfiniBandInterfaceCreateRequest)`

SetInfinibandInterfaces sets InfinibandInterfaces field to given value.

### HasInfinibandInterfaces

`func (o *InstanceUpdateRequest) HasInfinibandInterfaces() bool`

HasInfinibandInterfaces returns a boolean if a field has been set.

### GetNvLinkInterfaces

`func (o *InstanceUpdateRequest) GetNvLinkInterfaces() []NVLinkInterfaceCreateRequest`

GetNvLinkInterfaces returns the NvLinkInterfaces field if non-nil, zero value otherwise.

### GetNvLinkInterfacesOk

`func (o *InstanceUpdateRequest) GetNvLinkInterfacesOk() (*[]NVLinkInterfaceCreateRequest, bool)`

GetNvLinkInterfacesOk returns a tuple with the NvLinkInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinkInterfaces

`func (o *InstanceUpdateRequest) SetNvLinkInterfaces(v []NVLinkInterfaceCreateRequest)`

SetNvLinkInterfaces sets NvLinkInterfaces field to given value.

### HasNvLinkInterfaces

`func (o *InstanceUpdateRequest) HasNvLinkInterfaces() bool`

HasNvLinkInterfaces returns a boolean if a field has been set.

### GetDpuExtensionServiceDeployments

`func (o *InstanceUpdateRequest) GetDpuExtensionServiceDeployments() []DpuExtensionServiceDeploymentRequest`

GetDpuExtensionServiceDeployments returns the DpuExtensionServiceDeployments field if non-nil, zero value otherwise.

### GetDpuExtensionServiceDeploymentsOk

`func (o *InstanceUpdateRequest) GetDpuExtensionServiceDeploymentsOk() (*[]DpuExtensionServiceDeploymentRequest, bool)`

GetDpuExtensionServiceDeploymentsOk returns a tuple with the DpuExtensionServiceDeployments field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDpuExtensionServiceDeployments

`func (o *InstanceUpdateRequest) SetDpuExtensionServiceDeployments(v []DpuExtensionServiceDeploymentRequest)`

SetDpuExtensionServiceDeployments sets DpuExtensionServiceDeployments field to given value.

### HasDpuExtensionServiceDeployments

`func (o *InstanceUpdateRequest) HasDpuExtensionServiceDeployments() bool`

HasDpuExtensionServiceDeployments returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


