# VpcPrefixCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Human readable name for the VPC Prefix | 
**VpcId** | **string** | ID of the VPC | 
**IpBlockId** | Pointer to **string** |  | [optional] 
**PrefixLength** | **int32** |  | 

## Methods

### NewVpcPrefixCreateRequest

`func NewVpcPrefixCreateRequest(name string, vpcId string, prefixLength int32, ) *VpcPrefixCreateRequest`

NewVpcPrefixCreateRequest instantiates a new VpcPrefixCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVpcPrefixCreateRequestWithDefaults

`func NewVpcPrefixCreateRequestWithDefaults() *VpcPrefixCreateRequest`

NewVpcPrefixCreateRequestWithDefaults instantiates a new VpcPrefixCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *VpcPrefixCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *VpcPrefixCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *VpcPrefixCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetVpcId

`func (o *VpcPrefixCreateRequest) GetVpcId() string`

GetVpcId returns the VpcId field if non-nil, zero value otherwise.

### GetVpcIdOk

`func (o *VpcPrefixCreateRequest) GetVpcIdOk() (*string, bool)`

GetVpcIdOk returns a tuple with the VpcId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpcId

`func (o *VpcPrefixCreateRequest) SetVpcId(v string)`

SetVpcId sets VpcId field to given value.


### GetIpBlockId

`func (o *VpcPrefixCreateRequest) GetIpBlockId() string`

GetIpBlockId returns the IpBlockId field if non-nil, zero value otherwise.

### GetIpBlockIdOk

`func (o *VpcPrefixCreateRequest) GetIpBlockIdOk() (*string, bool)`

GetIpBlockIdOk returns a tuple with the IpBlockId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpBlockId

`func (o *VpcPrefixCreateRequest) SetIpBlockId(v string)`

SetIpBlockId sets IpBlockId field to given value.

### HasIpBlockId

`func (o *VpcPrefixCreateRequest) HasIpBlockId() bool`

HasIpBlockId returns a boolean if a field has been set.

### GetPrefixLength

`func (o *VpcPrefixCreateRequest) GetPrefixLength() int32`

GetPrefixLength returns the PrefixLength field if non-nil, zero value otherwise.

### GetPrefixLengthOk

`func (o *VpcPrefixCreateRequest) GetPrefixLengthOk() (*int32, bool)`

GetPrefixLengthOk returns a tuple with the PrefixLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefixLength

`func (o *VpcPrefixCreateRequest) SetPrefixLength(v int32)`

SetPrefixLength sets PrefixLength field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


