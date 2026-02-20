# VpcCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name of the VPC | 
**Description** | Pointer to **string** | Optional description for the VPC | [optional] 
**SiteId** | **string** | ID of the Site where the VPC should be created | 
**NetworkVirtualizationType** | Pointer to **NullableString** | Network virtualization type of the VPC. If no value is specified, then defaults to &#x60;FNN&#x60; if Site has native networking enabled, or ETHERNET_VIRTUALIZER if native networking is disabled | [optional] 
**NetworkSecurityGroupId** | Pointer to **NullableString** | ID of the Network Security Group to attach to the VPC | [optional] 
**NvLinkLogicalPartitionId** | Pointer to **NullableString** | ID of the default NVLink Logical Partition that GPUs for all Instances in the VPC will attach to | [optional] 
**Labels** | Pointer to **map[string]string** | String key value pairs describing VPC labels. Up to 10 key value pairs can be specified | [optional] 

## Methods

### NewVpcCreateRequest

`func NewVpcCreateRequest(name string, siteId string, ) *VpcCreateRequest`

NewVpcCreateRequest instantiates a new VpcCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVpcCreateRequestWithDefaults

`func NewVpcCreateRequestWithDefaults() *VpcCreateRequest`

NewVpcCreateRequestWithDefaults instantiates a new VpcCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *VpcCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *VpcCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *VpcCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *VpcCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *VpcCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *VpcCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *VpcCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *VpcCreateRequest) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *VpcCreateRequest) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *VpcCreateRequest) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.


### GetNetworkVirtualizationType

`func (o *VpcCreateRequest) GetNetworkVirtualizationType() string`

GetNetworkVirtualizationType returns the NetworkVirtualizationType field if non-nil, zero value otherwise.

### GetNetworkVirtualizationTypeOk

`func (o *VpcCreateRequest) GetNetworkVirtualizationTypeOk() (*string, bool)`

GetNetworkVirtualizationTypeOk returns a tuple with the NetworkVirtualizationType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkVirtualizationType

`func (o *VpcCreateRequest) SetNetworkVirtualizationType(v string)`

SetNetworkVirtualizationType sets NetworkVirtualizationType field to given value.

### HasNetworkVirtualizationType

`func (o *VpcCreateRequest) HasNetworkVirtualizationType() bool`

HasNetworkVirtualizationType returns a boolean if a field has been set.

### SetNetworkVirtualizationTypeNil

`func (o *VpcCreateRequest) SetNetworkVirtualizationTypeNil(b bool)`

 SetNetworkVirtualizationTypeNil sets the value for NetworkVirtualizationType to be an explicit nil

### UnsetNetworkVirtualizationType
`func (o *VpcCreateRequest) UnsetNetworkVirtualizationType()`

UnsetNetworkVirtualizationType ensures that no value is present for NetworkVirtualizationType, not even an explicit nil
### GetNetworkSecurityGroupId

`func (o *VpcCreateRequest) GetNetworkSecurityGroupId() string`

GetNetworkSecurityGroupId returns the NetworkSecurityGroupId field if non-nil, zero value otherwise.

### GetNetworkSecurityGroupIdOk

`func (o *VpcCreateRequest) GetNetworkSecurityGroupIdOk() (*string, bool)`

GetNetworkSecurityGroupIdOk returns a tuple with the NetworkSecurityGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkSecurityGroupId

`func (o *VpcCreateRequest) SetNetworkSecurityGroupId(v string)`

SetNetworkSecurityGroupId sets NetworkSecurityGroupId field to given value.

### HasNetworkSecurityGroupId

`func (o *VpcCreateRequest) HasNetworkSecurityGroupId() bool`

HasNetworkSecurityGroupId returns a boolean if a field has been set.

### SetNetworkSecurityGroupIdNil

`func (o *VpcCreateRequest) SetNetworkSecurityGroupIdNil(b bool)`

 SetNetworkSecurityGroupIdNil sets the value for NetworkSecurityGroupId to be an explicit nil

### UnsetNetworkSecurityGroupId
`func (o *VpcCreateRequest) UnsetNetworkSecurityGroupId()`

UnsetNetworkSecurityGroupId ensures that no value is present for NetworkSecurityGroupId, not even an explicit nil
### GetNvLinkLogicalPartitionId

`func (o *VpcCreateRequest) GetNvLinkLogicalPartitionId() string`

GetNvLinkLogicalPartitionId returns the NvLinkLogicalPartitionId field if non-nil, zero value otherwise.

### GetNvLinkLogicalPartitionIdOk

`func (o *VpcCreateRequest) GetNvLinkLogicalPartitionIdOk() (*string, bool)`

GetNvLinkLogicalPartitionIdOk returns a tuple with the NvLinkLogicalPartitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinkLogicalPartitionId

`func (o *VpcCreateRequest) SetNvLinkLogicalPartitionId(v string)`

SetNvLinkLogicalPartitionId sets NvLinkLogicalPartitionId field to given value.

### HasNvLinkLogicalPartitionId

`func (o *VpcCreateRequest) HasNvLinkLogicalPartitionId() bool`

HasNvLinkLogicalPartitionId returns a boolean if a field has been set.

### SetNvLinkLogicalPartitionIdNil

`func (o *VpcCreateRequest) SetNvLinkLogicalPartitionIdNil(b bool)`

 SetNvLinkLogicalPartitionIdNil sets the value for NvLinkLogicalPartitionId to be an explicit nil

### UnsetNvLinkLogicalPartitionId
`func (o *VpcCreateRequest) UnsetNvLinkLogicalPartitionId()`

UnsetNvLinkLogicalPartitionId ensures that no value is present for NvLinkLogicalPartitionId, not even an explicit nil
### GetLabels

`func (o *VpcCreateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *VpcCreateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *VpcCreateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *VpcCreateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


