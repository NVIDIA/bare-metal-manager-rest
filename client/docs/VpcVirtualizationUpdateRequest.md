# VpcVirtualizationUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NetworkVirtualizationType** | Pointer to **string** | Network virtualization type of the VPC. Can only be updated to &#x60;FNN&#x60; | [optional] 

## Methods

### NewVpcVirtualizationUpdateRequest

`func NewVpcVirtualizationUpdateRequest() *VpcVirtualizationUpdateRequest`

NewVpcVirtualizationUpdateRequest instantiates a new VpcVirtualizationUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVpcVirtualizationUpdateRequestWithDefaults

`func NewVpcVirtualizationUpdateRequestWithDefaults() *VpcVirtualizationUpdateRequest`

NewVpcVirtualizationUpdateRequestWithDefaults instantiates a new VpcVirtualizationUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNetworkVirtualizationType

`func (o *VpcVirtualizationUpdateRequest) GetNetworkVirtualizationType() string`

GetNetworkVirtualizationType returns the NetworkVirtualizationType field if non-nil, zero value otherwise.

### GetNetworkVirtualizationTypeOk

`func (o *VpcVirtualizationUpdateRequest) GetNetworkVirtualizationTypeOk() (*string, bool)`

GetNetworkVirtualizationTypeOk returns a tuple with the NetworkVirtualizationType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkVirtualizationType

`func (o *VpcVirtualizationUpdateRequest) SetNetworkVirtualizationType(v string)`

SetNetworkVirtualizationType sets NetworkVirtualizationType field to given value.

### HasNetworkVirtualizationType

`func (o *VpcVirtualizationUpdateRequest) HasNetworkVirtualizationType() bool`

HasNetworkVirtualizationType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


