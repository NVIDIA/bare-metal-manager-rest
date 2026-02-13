# VPC

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | ID of the VPC | [optional] [readonly] 
**Name** | Pointer to **string** | Name of the VPC | [optional] 
**Description** | Pointer to **string** | Description of the VPC, can be empty | [optional] 
**Org** | Pointer to **string** | Organization the VPC belongs to | [optional] [readonly] 
**TenantId** | Pointer to **string** | ID of the Tenant the VPC belongs to | [optional] [readonly] 
**SiteId** | Pointer to **string** | ID of the Site the VPC belongs to | [optional] 
**ControllerVpcId** | Pointer to **NullableString** | Legacy attribute, contains the same value as ID | [optional] 
**NetworkVirtualizationType** | Pointer to **string** | Network virtualization type of the VPC | [optional] 
**NetworkSecurityGroupId** | Pointer to **NullableString** | ID of the Network Security Group attached to the VPC | [optional] 
**NetworkSecurityGroupPropagationDetails** | Pointer to [**NetworkSecurityGroupPropagationDetails**](NetworkSecurityGroupPropagationDetails.md) | Propagation details for the attached Network Security Group | [optional] 
**NvLinkLogicalPartitionId** | Pointer to **NullableString** | ID of the default NVLink Logical Partition that GPUs for all Instances in the VPC will attach to | [optional] 
**Labels** | Pointer to **map[string]string** | String key value pairs describing VPC labels | [optional] 
**Status** | Pointer to [**VpcStatus**](VpcStatus.md) | Status of the VPC | [optional] [readonly] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) | History of status changes for the VPC | [optional] 
**Created** | Pointer to **time.Time** | Date/time when VPC was created | [optional] [readonly] 
**Updated** | Pointer to **time.Time** | Date/time when VPC was last updated | [optional] [readonly] 

## Methods

### NewVPC

`func NewVPC() *VPC`

NewVPC instantiates a new VPC object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVPCWithDefaults

`func NewVPCWithDefaults() *VPC`

NewVPCWithDefaults instantiates a new VPC object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *VPC) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *VPC) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *VPC) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *VPC) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *VPC) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *VPC) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *VPC) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *VPC) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *VPC) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *VPC) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *VPC) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *VPC) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetOrg

`func (o *VPC) GetOrg() string`

GetOrg returns the Org field if non-nil, zero value otherwise.

### GetOrgOk

`func (o *VPC) GetOrgOk() (*string, bool)`

GetOrgOk returns a tuple with the Org field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrg

`func (o *VPC) SetOrg(v string)`

SetOrg sets Org field to given value.

### HasOrg

`func (o *VPC) HasOrg() bool`

HasOrg returns a boolean if a field has been set.

### GetTenantId

`func (o *VPC) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *VPC) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *VPC) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *VPC) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetSiteId

`func (o *VPC) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *VPC) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *VPC) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *VPC) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetControllerVpcId

`func (o *VPC) GetControllerVpcId() string`

GetControllerVpcId returns the ControllerVpcId field if non-nil, zero value otherwise.

### GetControllerVpcIdOk

`func (o *VPC) GetControllerVpcIdOk() (*string, bool)`

GetControllerVpcIdOk returns a tuple with the ControllerVpcId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerVpcId

`func (o *VPC) SetControllerVpcId(v string)`

SetControllerVpcId sets ControllerVpcId field to given value.

### HasControllerVpcId

`func (o *VPC) HasControllerVpcId() bool`

HasControllerVpcId returns a boolean if a field has been set.

### SetControllerVpcIdNil

`func (o *VPC) SetControllerVpcIdNil(b bool)`

 SetControllerVpcIdNil sets the value for ControllerVpcId to be an explicit nil

### UnsetControllerVpcId
`func (o *VPC) UnsetControllerVpcId()`

UnsetControllerVpcId ensures that no value is present for ControllerVpcId, not even an explicit nil
### GetNetworkVirtualizationType

`func (o *VPC) GetNetworkVirtualizationType() string`

GetNetworkVirtualizationType returns the NetworkVirtualizationType field if non-nil, zero value otherwise.

### GetNetworkVirtualizationTypeOk

`func (o *VPC) GetNetworkVirtualizationTypeOk() (*string, bool)`

GetNetworkVirtualizationTypeOk returns a tuple with the NetworkVirtualizationType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkVirtualizationType

`func (o *VPC) SetNetworkVirtualizationType(v string)`

SetNetworkVirtualizationType sets NetworkVirtualizationType field to given value.

### HasNetworkVirtualizationType

`func (o *VPC) HasNetworkVirtualizationType() bool`

HasNetworkVirtualizationType returns a boolean if a field has been set.

### GetNetworkSecurityGroupId

`func (o *VPC) GetNetworkSecurityGroupId() string`

GetNetworkSecurityGroupId returns the NetworkSecurityGroupId field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupIdOk

`func (o *VPC) GetNetworkSecurityGroupIdOk() (*string, bool)`

GetNetworkSecurityGroupIdOk returns a tuple with the NetworkSecurityGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupId

`func (o *VPC) SetNetworkSecurityGroupId(v string)`

SetNetworkSecurityGroupId sets NetworkSecurityGroupId field to given value.

### HasNetworkSecurityGroupId

`func (o *VPC) HasNetworkSecurityGroupId() bool`

HasNetworkSecurityGroupId returns a boolean if a field has been set.

### SetNetworkSecurityGroupIdNil

`func (o *VPC) SetNetworkSecurityGroupIdNil(b bool)`

 SetNetworkSecurityGroupIdNil sets the value for NetworkSecurityGroupId to be an explicit nil

### UnsetNetworkSecurityGroupId
`func (o *VPC) UnsetNetworkSecurityGroupId()`

UnsetNetworkSecurityGroupId ensures that no value is present for NetworkSecurityGroupId, not even an explicit nil
### GetNetworkSecurityGroupPropagationDetails

`func (o *VPC) GetNetworkSecurityGroupPropagationDetails() NetworkSecurityGroupPropagationDetails`

GetNetworkSecurityGroupPropagationDetails returns the NetworkSecurityGroupPropagationDetails field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupPropagationDetailsOk

`func (o *VPC) GetNetworkSecurityGroupPropagationDetailsOk() (*NetworkSecurityGroupPropagationDetails, bool)`

GetNetworkSecurityGroupPropagationDetailsOk returns a tuple with the NetworkSecurityGroupPropagationDetails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupPropagationDetails

`func (o *VPC) SetNetworkSecurityGroupPropagationDetails(v NetworkSecurityGroupPropagationDetails)`

SetNetworkSecurityGroupPropagationDetails sets NetworkSecurityGroupPropagationDetails field to given value.

### HasNetworkSecurityGroupPropagationDetails

`func (o *VPC) HasNetworkSecurityGroupPropagationDetails() bool`

HasNetworkSecurityGroupPropagationDetails returns a boolean if a field has been set.

### GetNvLinkLogicalPartitionId

`func (o *VPC) GetNvLinkLogicalPartitionId() string`

GetNvLinkLogicalPartitionId returns the NvLinkLogicalPartitionId field if non-nil, zero value otherwise.

### GetNvLinkLogicalPartitionIdOk

`func (o *VPC) GetNvLinkLogicalPartitionIdOk() (*string, bool)`

GetNvLinkLogicalPartitionIdOk returns a tuple with the NvLinkLogicalPartitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinkLogicalPartitionId

`func (o *VPC) SetNvLinkLogicalPartitionId(v string)`

SetNvLinkLogicalPartitionId sets NvLinkLogicalPartitionId field to given value.

### HasNvLinkLogicalPartitionId

`func (o *VPC) HasNvLinkLogicalPartitionId() bool`

HasNvLinkLogicalPartitionId returns a boolean if a field has been set.

### SetNvLinkLogicalPartitionIdNil

`func (o *VPC) SetNvLinkLogicalPartitionIdNil(b bool)`

 SetNvLinkLogicalPartitionIdNil sets the value for NvLinkLogicalPartitionId to be an explicit nil

### UnsetNvLinkLogicalPartitionId
`func (o *VPC) UnsetNvLinkLogicalPartitionId()`

UnsetNvLinkLogicalPartitionId ensures that no value is present for NvLinkLogicalPartitionId, not even an explicit nil
### GetLabels

`func (o *VPC) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *VPC) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *VPC) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *VPC) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetStatus

`func (o *VPC) GetStatus() VpcStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *VPC) GetStatusOk() (*VpcStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *VPC) SetStatus(v VpcStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *VPC) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *VPC) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *VPC) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *VPC) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *VPC) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *VPC) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *VPC) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *VPC) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *VPC) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *VPC) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *VPC) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *VPC) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *VPC) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


