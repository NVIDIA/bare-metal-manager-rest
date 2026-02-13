# Instance

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier for the Instance | [optional] [readonly] 
**Name** | Pointer to **string** | Name for the Instance | [optional] 
**Description** | Pointer to **string** | Description for the Instance | [optional] 
**TenantId** | Pointer to **string** | ID of the Tenant the Instance belongs to | [optional] 
**InfrastructureProviderId** | Pointer to **string** | ID of the Infrastructure Provider that owns the Site where the Instance is located | [optional] 
**SiteId** | Pointer to **string** | ID of the Site where the Instance is located | [optional] 
**InstanceTypeId** | Pointer to **string** |  | [optional] 
**VpcId** | Pointer to **string** |  | [optional] 
**MachineId** | Pointer to **NullableString** |  | [optional] 
**OperatingSystemId** | Pointer to **string** |  | [optional] 
**NetworkSecurityGroupId** | Pointer to **NullableString** |  | [optional] 
**NetworkSecurityGroupPropagationDetails** | Pointer to [**NetworkSecurityGroupPropagationDetails**](NetworkSecurityGroupPropagationDetails.md) |  | [optional] 
**NetworkSecurityGroupInherited** | Pointer to **bool** | Indicates if the Network Security Group is inherited from VPC | [optional] 
**ControllerInstanceId** | Pointer to **NullableString** |  | [optional] 
**IpxeScript** | Pointer to **NullableString** |  | [optional] 
**AlwaysBootWithCustomIpxe** | Pointer to **bool** | Indicates whether the Instance should always execute custom iPXE script when rebooting | [optional] 
**PhoneHomeEnabled** | Pointer to **bool** | Indicates whether the Phone Home service should be enabled or disabled for Instance | [optional] 
**UserData** | Pointer to **NullableString** |  | [optional] 
**Labels** | Pointer to **map[string]string** |  | [optional] 
**IsUpdatePending** | Pointer to **bool** | Indicates whether an update is available for the Instance. Updates can be applied on reboot | [optional] 
**SerialConsoleUrl** | Pointer to **NullableString** | Serial Console URL for the Instance. Format: ssh://&lt;id&gt;@siteSerialConsoleHostname | [optional] 
**Interfaces** | Pointer to [**[]Interface**](Interface.md) |  | [optional] 
**InfinibandInterfaces** | Pointer to [**[]InfiniBandInterface**](InfiniBandInterface.md) |  | [optional] 
**NvLinkInterfaces** | Pointer to [**[]NVLinkInterface**](NVLinkInterface.md) |  | [optional] 
**DpuExtensionServiceDeployments** | Pointer to [**[]DpuExtensionServiceDeployment**](DpuExtensionServiceDeployment.md) | DPU Extension Services deployed on DPUs of this Instance | [optional] 
**SshKeyGroupIds** | Pointer to **[]string** | IDs of SSH Key Groups associated with this Instance | [optional] 
**SshKeyGroups** | Pointer to [**[]SshKeyGroup**](SshKeyGroup.md) | IDs of SSH Key Groups associated with this Instance | [optional] 
**TpmEkCertificate** | Pointer to **NullableString** | base64 encoded TPM EK Certificate associated with this Instance | [optional] 
**Status** | Pointer to [**InstanceStatus**](InstanceStatus.md) |  | [optional] [readonly] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**Deprecations** | Pointer to [**[]Deprecation**](Deprecation.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewInstance

`func NewInstance() *Instance`

NewInstance instantiates a new Instance object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceWithDefaults

`func NewInstanceWithDefaults() *Instance`

NewInstanceWithDefaults instantiates a new Instance object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Instance) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Instance) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Instance) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Instance) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *Instance) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Instance) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Instance) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Instance) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *Instance) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Instance) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Instance) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Instance) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetTenantId

`func (o *Instance) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *Instance) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *Instance) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *Instance) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *Instance) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *Instance) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *Instance) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *Instance) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetSiteId

`func (o *Instance) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *Instance) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *Instance) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *Instance) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetInstanceTypeId

`func (o *Instance) GetInstanceTypeId() string`

GetInstanceTypeId returns the InstanceTypeId field if non-nil, zero value otherwise.

### GetInstanceTypeIdOk

`func (o *Instance) GetInstanceTypeIdOk() (*string, bool)`

GetInstanceTypeIdOk returns a tuple with the InstanceTypeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceTypeId

`func (o *Instance) SetInstanceTypeId(v string)`

SetInstanceTypeId sets InstanceTypeId field to given value.

### HasInstanceTypeId

`func (o *Instance) HasInstanceTypeId() bool`

HasInstanceTypeId returns a boolean if a field has been set.

### GetVpcId

`func (o *Instance) GetVpcId() string`

GetVpcId returns the VpcId field if non-nil, zero value otherwise.

### GetVpcIdOk

`func (o *Instance) GetVpcIdOk() (*string, bool)`

GetVpcIdOk returns a tuple with the VpcId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpcId

`func (o *Instance) SetVpcId(v string)`

SetVpcId sets VpcId field to given value.

### HasVpcId

`func (o *Instance) HasVpcId() bool`

HasVpcId returns a boolean if a field has been set.

### GetMachineId

`func (o *Instance) GetMachineId() string`

GetMachineId returns the MachineId field if non-nil, zero value otherwise.

### GetMachineIdOk

`func (o *Instance) GetMachineIdOk() (*string, bool)`

GetMachineIdOk returns a tuple with the MachineId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineId

`func (o *Instance) SetMachineId(v string)`

SetMachineId sets MachineId field to given value.

### HasMachineId

`func (o *Instance) HasMachineId() bool`

HasMachineId returns a boolean if a field has been set.

### SetMachineIdNil

`func (o *Instance) SetMachineIdNil(b bool)`

 SetMachineIdNil sets the value for MachineId to be an explicit nil

### UnsetMachineId
`func (o *Instance) UnsetMachineId()`

UnsetMachineId ensures that no value is present for MachineId, not even an explicit nil
### GetOperatingSystemId

`func (o *Instance) GetOperatingSystemId() string`

GetOperatingSystemId returns the OperatingSystemId field if non-nil, zero value otherwise.

### GetOperatingSystemIdOk

`func (o *Instance) GetOperatingSystemIdOk() (*string, bool)`

GetOperatingSystemIdOk returns a tuple with the OperatingSystemId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperatingSystemId

`func (o *Instance) SetOperatingSystemId(v string)`

SetOperatingSystemId sets OperatingSystemId field to given value.

### HasOperatingSystemId

`func (o *Instance) HasOperatingSystemId() bool`

HasOperatingSystemId returns a boolean if a field has been set.

### GetNetworkSecurityGroupId

`func (o *Instance) GetNetworkSecurityGroupId() string`

GetNetworkSecurityGroupId returns the NetworkSecurityGroupId field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupIdOk

`func (o *Instance) GetNetworkSecurityGroupIdOk() (*string, bool)`

GetNetworkSecurityGroupIdOk returns a tuple with the NetworkSecurityGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupId

`func (o *Instance) SetNetworkSecurityGroupId(v string)`

SetNetworkSecurityGroupId sets NetworkSecurityGroupId field to given value.

### HasNetworkSecurityGroupId

`func (o *Instance) HasNetworkSecurityGroupId() bool`

HasNetworkSecurityGroupId returns a boolean if a field has been set.

### SetNetworkSecurityGroupIdNil

`func (o *Instance) SetNetworkSecurityGroupIdNil(b bool)`

 SetNetworkSecurityGroupIdNil sets the value for NetworkSecurityGroupId to be an explicit nil

### UnsetNetworkSecurityGroupId
`func (o *Instance) UnsetNetworkSecurityGroupId()`

UnsetNetworkSecurityGroupId ensures that no value is present for NetworkSecurityGroupId, not even an explicit nil
### GetNetworkSecurityGroupPropagationDetails

`func (o *Instance) GetNetworkSecurityGroupPropagationDetails() NetworkSecurityGroupPropagationDetails`

GetNetworkSecurityGroupPropagationDetails returns the NetworkSecurityGroupPropagationDetails field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupPropagationDetailsOk

`func (o *Instance) GetNetworkSecurityGroupPropagationDetailsOk() (*NetworkSecurityGroupPropagationDetails, bool)`

GetNetworkSecurityGroupPropagationDetailsOk returns a tuple with the NetworkSecurityGroupPropagationDetails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupPropagationDetails

`func (o *Instance) SetNetworkSecurityGroupPropagationDetails(v NetworkSecurityGroupPropagationDetails)`

SetNetworkSecurityGroupPropagationDetails sets NetworkSecurityGroupPropagationDetails field to given value.

### HasNetworkSecurityGroupPropagationDetails

`func (o *Instance) HasNetworkSecurityGroupPropagationDetails() bool`

HasNetworkSecurityGroupPropagationDetails returns a boolean if a field has been set.

### GetNetworkSecurityGroupInherited

`func (o *Instance) GetNetworkSecurityGroupInherited() bool`

GetNetworkSecurityGroupInherited returns the NetworkSecurityGroupInherited field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupInheritedOk

`func (o *Instance) GetNetworkSecurityGroupInheritedOk() (*bool, bool)`

GetNetworkSecurityGroupInheritedOk returns a tuple with the NetworkSecurityGroupInherited field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupInherited

`func (o *Instance) SetNetworkSecurityGroupInherited(v bool)`

SetNetworkSecurityGroupInherited sets NetworkSecurityGroupInherited field to given value.

### HasNetworkSecurityGroupInherited

`func (o *Instance) HasNetworkSecurityGroupInherited() bool`

HasNetworkSecurityGroupInherited returns a boolean if a field has been set.

### GetControllerInstanceId

`func (o *Instance) GetControllerInstanceId() string`

GetControllerInstanceId returns the ControllerInstanceId field if non-nil, zero value otherwise.

### GetControllerInstanceIdOk

`func (o *Instance) GetControllerInstanceIdOk() (*string, bool)`

GetControllerInstanceIdOk returns a tuple with the ControllerInstanceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerInstanceId

`func (o *Instance) SetControllerInstanceId(v string)`

SetControllerInstanceId sets ControllerInstanceId field to given value.

### HasControllerInstanceId

`func (o *Instance) HasControllerInstanceId() bool`

HasControllerInstanceId returns a boolean if a field has been set.

### SetControllerInstanceIdNil

`func (o *Instance) SetControllerInstanceIdNil(b bool)`

 SetControllerInstanceIdNil sets the value for ControllerInstanceId to be an explicit nil

### UnsetControllerInstanceId
`func (o *Instance) UnsetControllerInstanceId()`

UnsetControllerInstanceId ensures that no value is present for ControllerInstanceId, not even an explicit nil
### GetIpxeScript

`func (o *Instance) GetIpxeScript() string`

GetIpxeScript returns the IpxeScript field if non-nil, zero value otherwise.

### GetIpxeScriptOk

`func (o *Instance) GetIpxeScriptOk() (*string, bool)`

GetIpxeScriptOk returns a tuple with the IpxeScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpxeScript

`func (o *Instance) SetIpxeScript(v string)`

SetIpxeScript sets IpxeScript field to given value.

### HasIpxeScript

`func (o *Instance) HasIpxeScript() bool`

HasIpxeScript returns a boolean if a field has been set.

### SetIpxeScriptNil

`func (o *Instance) SetIpxeScriptNil(b bool)`

 SetIpxeScriptNil sets the value for IpxeScript to be an explicit nil

### UnsetIpxeScript
`func (o *Instance) UnsetIpxeScript()`

UnsetIpxeScript ensures that no value is present for IpxeScript, not even an explicit nil
### GetAlwaysBootWithCustomIpxe

`func (o *Instance) GetAlwaysBootWithCustomIpxe() bool`

GetAlwaysBootWithCustomIpxe returns the AlwaysBootWithCustomIpxe field if non-nil, zero value otherwise.

### GetAlwaysBootWithCustomIpxeOk

`func (o *Instance) GetAlwaysBootWithCustomIpxeOk() (*bool, bool)`

GetAlwaysBootWithCustomIpxeOk returns a tuple with the AlwaysBootWithCustomIpxe field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlwaysBootWithCustomIpxe

`func (o *Instance) SetAlwaysBootWithCustomIpxe(v bool)`

SetAlwaysBootWithCustomIpxe sets AlwaysBootWithCustomIpxe field to given value.

### HasAlwaysBootWithCustomIpxe

`func (o *Instance) HasAlwaysBootWithCustomIpxe() bool`

HasAlwaysBootWithCustomIpxe returns a boolean if a field has been set.

### GetPhoneHomeEnabled

`func (o *Instance) GetPhoneHomeEnabled() bool`

GetPhoneHomeEnabled returns the PhoneHomeEnabled field if non-nil, zero value otherwise.

### GetPhoneHomeEnabledOk

`func (o *Instance) GetPhoneHomeEnabledOk() (*bool, bool)`

GetPhoneHomeEnabledOk returns a tuple with the PhoneHomeEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneHomeEnabled

`func (o *Instance) SetPhoneHomeEnabled(v bool)`

SetPhoneHomeEnabled sets PhoneHomeEnabled field to given value.

### HasPhoneHomeEnabled

`func (o *Instance) HasPhoneHomeEnabled() bool`

HasPhoneHomeEnabled returns a boolean if a field has been set.

### GetUserData

`func (o *Instance) GetUserData() string`

GetUserData returns the UserData field if non-nil, zero value otherwise.

### GetUserDataOk

`func (o *Instance) GetUserDataOk() (*string, bool)`

GetUserDataOk returns a tuple with the UserData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserData

`func (o *Instance) SetUserData(v string)`

SetUserData sets UserData field to given value.

### HasUserData

`func (o *Instance) HasUserData() bool`

HasUserData returns a boolean if a field has been set.

### SetUserDataNil

`func (o *Instance) SetUserDataNil(b bool)`

 SetUserDataNil sets the value for UserData to be an explicit nil

### UnsetUserData
`func (o *Instance) UnsetUserData()`

UnsetUserData ensures that no value is present for UserData, not even an explicit nil
### GetLabels

`func (o *Instance) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *Instance) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *Instance) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *Instance) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetIsUpdatePending

`func (o *Instance) GetIsUpdatePending() bool`

GetIsUpdatePending returns the IsUpdatePending field if non-nil, zero value otherwise.

### GetIsUpdatePendingOk

`func (o *Instance) GetIsUpdatePendingOk() (*bool, bool)`

GetIsUpdatePendingOk returns a tuple with the IsUpdatePending field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsUpdatePending

`func (o *Instance) SetIsUpdatePending(v bool)`

SetIsUpdatePending sets IsUpdatePending field to given value.

### HasIsUpdatePending

`func (o *Instance) HasIsUpdatePending() bool`

HasIsUpdatePending returns a boolean if a field has been set.

### GetSerialConsoleUrl

`func (o *Instance) GetSerialConsoleUrl() string`

GetSerialConsoleUrl returns the SerialConsoleUrl field if non-nil, zero value otherwise.

### GetSerialConsoleUrlOk

`func (o *Instance) GetSerialConsoleUrlOk() (*string, bool)`

GetSerialConsoleUrlOk returns a tuple with the SerialConsoleUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialConsoleUrl

`func (o *Instance) SetSerialConsoleUrl(v string)`

SetSerialConsoleUrl sets SerialConsoleUrl field to given value.

### HasSerialConsoleUrl

`func (o *Instance) HasSerialConsoleUrl() bool`

HasSerialConsoleUrl returns a boolean if a field has been set.

### SetSerialConsoleUrlNil

`func (o *Instance) SetSerialConsoleUrlNil(b bool)`

 SetSerialConsoleUrlNil sets the value for SerialConsoleUrl to be an explicit nil

### UnsetSerialConsoleUrl
`func (o *Instance) UnsetSerialConsoleUrl()`

UnsetSerialConsoleUrl ensures that no value is present for SerialConsoleUrl, not even an explicit nil
### GetInterfaces

`func (o *Instance) GetInterfaces() []Interface`

GetInterfaces returns the Interfaces field if non-nil, zero value otherwise.

### GetInterfacesOk

`func (o *Instance) GetInterfacesOk() (*[]Interface, bool)`

GetInterfacesOk returns a tuple with the Interfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterfaces

`func (o *Instance) SetInterfaces(v []Interface)`

SetInterfaces sets Interfaces field to given value.

### HasInterfaces

`func (o *Instance) HasInterfaces() bool`

HasInterfaces returns a boolean if a field has been set.

### GetInfinibandInterfaces

`func (o *Instance) GetInfinibandInterfaces() []InfiniBandInterface`

GetInfinibandInterfaces returns the InfinibandInterfaces field if non-nil, zero value otherwise.

### GetInfinibandInterfacesOk

`func (o *Instance) GetInfinibandInterfacesOk() (*[]InfiniBandInterface, bool)`

GetInfinibandInterfacesOk returns a tuple with the InfinibandInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfinibandInterfaces

`func (o *Instance) SetInfinibandInterfaces(v []InfiniBandInterface)`

SetInfinibandInterfaces sets InfinibandInterfaces field to given value.

### HasInfinibandInterfaces

`func (o *Instance) HasInfinibandInterfaces() bool`

HasInfinibandInterfaces returns a boolean if a field has been set.

### GetNvLinkInterfaces

`func (o *Instance) GetNvLinkInterfaces() []NVLinkInterface`

GetNvLinkInterfaces returns the NvLinkInterfaces field if non-nil, zero value otherwise.

### GetNvLinkInterfacesOk

`func (o *Instance) GetNvLinkInterfacesOk() (*[]NVLinkInterface, bool)`

GetNvLinkInterfacesOk returns a tuple with the NvLinkInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinkInterfaces

`func (o *Instance) SetNvLinkInterfaces(v []NVLinkInterface)`

SetNvLinkInterfaces sets NvLinkInterfaces field to given value.

### HasNvLinkInterfaces

`func (o *Instance) HasNvLinkInterfaces() bool`

HasNvLinkInterfaces returns a boolean if a field has been set.

### GetDpuExtensionServiceDeployments

`func (o *Instance) GetDpuExtensionServiceDeployments() []DpuExtensionServiceDeployment`

GetDpuExtensionServiceDeployments returns the DpuExtensionServiceDeployments field if non-nil, zero value otherwise.

### GetDpuExtensionServiceDeploymentsOk

`func (o *Instance) GetDpuExtensionServiceDeploymentsOk() (*[]DpuExtensionServiceDeployment, bool)`

GetDpuExtensionServiceDeploymentsOk returns a tuple with the DpuExtensionServiceDeployments field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDpuExtensionServiceDeployments

`func (o *Instance) SetDpuExtensionServiceDeployments(v []DpuExtensionServiceDeployment)`

SetDpuExtensionServiceDeployments sets DpuExtensionServiceDeployments field to given value.

### HasDpuExtensionServiceDeployments

`func (o *Instance) HasDpuExtensionServiceDeployments() bool`

HasDpuExtensionServiceDeployments returns a boolean if a field has been set.

### GetSshKeyGroupIds

`func (o *Instance) GetSshKeyGroupIds() []string`

GetSshKeyGroupIds returns the SshKeyGroupIds field if non-nil, zero value otherwise.

### GetSshKeyGroupIdsOk

`func (o *Instance) GetSshKeyGroupIdsOk() (*[]string, bool)`

GetSshKeyGroupIdsOk returns a tuple with the SshKeyGroupIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyGroupIds

`func (o *Instance) SetSshKeyGroupIds(v []string)`

SetSshKeyGroupIds sets SshKeyGroupIds field to given value.

### HasSshKeyGroupIds

`func (o *Instance) HasSshKeyGroupIds() bool`

HasSshKeyGroupIds returns a boolean if a field has been set.

### GetSshKeyGroups

`func (o *Instance) GetSshKeyGroups() []SshKeyGroup`

GetSshKeyGroups returns the SshKeyGroups field if non-nil, zero value otherwise.

### GetSshKeyGroupsOk

`func (o *Instance) GetSshKeyGroupsOk() (*[]SshKeyGroup, bool)`

GetSshKeyGroupsOk returns a tuple with the SshKeyGroups field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyGroups

`func (o *Instance) SetSshKeyGroups(v []SshKeyGroup)`

SetSshKeyGroups sets SshKeyGroups field to given value.

### HasSshKeyGroups

`func (o *Instance) HasSshKeyGroups() bool`

HasSshKeyGroups returns a boolean if a field has been set.

### GetTpmEkCertificate

`func (o *Instance) GetTpmEkCertificate() string`

GetTpmEkCertificate returns the TpmEkCertificate field if non-nil, zero value otherwise.

### GetTpmEkCertificateOk

`func (o *Instance) GetTpmEkCertificateOk() (*string, bool)`

GetTpmEkCertificateOk returns a tuple with the TpmEkCertificate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTpmEkCertificate

`func (o *Instance) SetTpmEkCertificate(v string)`

SetTpmEkCertificate sets TpmEkCertificate field to given value.

### HasTpmEkCertificate

`func (o *Instance) HasTpmEkCertificate() bool`

HasTpmEkCertificate returns a boolean if a field has been set.

### SetTpmEkCertificateNil

`func (o *Instance) SetTpmEkCertificateNil(b bool)`

 SetTpmEkCertificateNil sets the value for TpmEkCertificate to be an explicit nil

### UnsetTpmEkCertificate
`func (o *Instance) UnsetTpmEkCertificate()`

UnsetTpmEkCertificate ensures that no value is present for TpmEkCertificate, not even an explicit nil
### GetStatus

`func (o *Instance) GetStatus() InstanceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Instance) GetStatusOk() (*InstanceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Instance) SetStatus(v InstanceStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Instance) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *Instance) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *Instance) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *Instance) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *Instance) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetDeprecations

`func (o *Instance) GetDeprecations() []Deprecation`

GetDeprecations returns the Deprecations field if non-nil, zero value otherwise.

### GetDeprecationsOk

`func (o *Instance) GetDeprecationsOk() (*[]Deprecation, bool)`

GetDeprecationsOk returns a tuple with the Deprecations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeprecations

`func (o *Instance) SetDeprecations(v []Deprecation)`

SetDeprecations sets Deprecations field to given value.

### HasDeprecations

`func (o *Instance) HasDeprecations() bool`

HasDeprecations returns a boolean if a field has been set.

### GetCreated

`func (o *Instance) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *Instance) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *Instance) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *Instance) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *Instance) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *Instance) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *Instance) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *Instance) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


