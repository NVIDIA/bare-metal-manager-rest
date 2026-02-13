# VpcUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Updated name of the VPC | [optional] 
**Description** | Pointer to **string** | Updated description of the VPC | [optional] 
**NetworkSecurityGroupId** | Pointer to **NullableString** | ID of the Network Security Group to attach to the VPC | [optional] 
**NvLinkLogicalPartitionId** | Pointer to **NullableString** | ID of the default NVLink Logical Partition that GPUs for all Instances in the VPC will attach to. Can only be updated if VPC currently has no active Instances | [optional] 
**Labels** | Pointer to **map[string]string** | Update labels of the VPC. Up to 10 key value pairs can be specified. The labels will be entirely replaced by those sent in the request. Any labels not included in the request will be removed. To retain existing labels, first fetch them and include them along with this request. | [optional] 

## Methods

### NewVpcUpdateRequest

`func NewVpcUpdateRequest() *VpcUpdateRequest`

NewVpcUpdateRequest instantiates a new VpcUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVpcUpdateRequestWithDefaults

`func NewVpcUpdateRequestWithDefaults() *VpcUpdateRequest`

NewVpcUpdateRequestWithDefaults instantiates a new VpcUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *VpcUpdateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *VpcUpdateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *VpcUpdateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *VpcUpdateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *VpcUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *VpcUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *VpcUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *VpcUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetNetworkSecurityGroupId

`func (o *VpcUpdateRequest) GetNetworkSecurityGroupId() string`

GetNetworkSecurityGroupId returns the NetworkSecurityGroupId field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupIdOk

`func (o *VpcUpdateRequest) GetNetworkSecurityGroupIdOk() (*string, bool)`

GetNetworkSecurityGroupIdOk returns a tuple with the NetworkSecurityGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupId

`func (o *VpcUpdateRequest) SetNetworkSecurityGroupId(v string)`

SetNetworkSecurityGroupId sets NetworkSecurityGroupId field to given value.

### HasNetworkSecurityGroupId

`func (o *VpcUpdateRequest) HasNetworkSecurityGroupId() bool`

HasNetworkSecurityGroupId returns a boolean if a field has been set.

### SetNetworkSecurityGroupIdNil

`func (o *VpcUpdateRequest) SetNetworkSecurityGroupIdNil(b bool)`

 SetNetworkSecurityGroupIdNil sets the value for NetworkSecurityGroupId to be an explicit nil

### UnsetNetworkSecurityGroupId
`func (o *VpcUpdateRequest) UnsetNetworkSecurityGroupId()`

UnsetNetworkSecurityGroupId ensures that no value is present for NetworkSecurityGroupId, not even an explicit nil
### GetNvLinkLogicalPartitionId

`func (o *VpcUpdateRequest) GetNvLinkLogicalPartitionId() string`

GetNvLinkLogicalPartitionId returns the NvLinkLogicalPartitionId field if non-nil, zero value otherwise.

### GetNvLinkLogicalPartitionIdOk

`func (o *VpcUpdateRequest) GetNvLinkLogicalPartitionIdOk() (*string, bool)`

GetNvLinkLogicalPartitionIdOk returns a tuple with the NvLinkLogicalPartitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinkLogicalPartitionId

`func (o *VpcUpdateRequest) SetNvLinkLogicalPartitionId(v string)`

SetNvLinkLogicalPartitionId sets NvLinkLogicalPartitionId field to given value.

### HasNvLinkLogicalPartitionId

`func (o *VpcUpdateRequest) HasNvLinkLogicalPartitionId() bool`

HasNvLinkLogicalPartitionId returns a boolean if a field has been set.

### SetNvLinkLogicalPartitionIdNil

`func (o *VpcUpdateRequest) SetNvLinkLogicalPartitionIdNil(b bool)`

 SetNvLinkLogicalPartitionIdNil sets the value for NvLinkLogicalPartitionId to be an explicit nil

### UnsetNvLinkLogicalPartitionId
`func (o *VpcUpdateRequest) UnsetNvLinkLogicalPartitionId()`

UnsetNvLinkLogicalPartitionId ensures that no value is present for NvLinkLogicalPartitionId, not even an explicit nil
### GetLabels

`func (o *VpcUpdateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *VpcUpdateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *VpcUpdateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *VpcUpdateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


