# SubnetCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**VpcId** | **string** |  | 
**Ipv4BlockId** | Pointer to **string** |  | [optional] 
**Ipv6BlockId** | Pointer to **string** |  | [optional] 
**PrefixLength** | **int32** |  | 

## Methods

### NewSubnetCreateRequest

`func NewSubnetCreateRequest(name string, vpcId string, prefixLength int32, ) *SubnetCreateRequest`

NewSubnetCreateRequest instantiates a new SubnetCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSubnetCreateRequestWithDefaults

`func NewSubnetCreateRequestWithDefaults() *SubnetCreateRequest`

NewSubnetCreateRequestWithDefaults instantiates a new SubnetCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SubnetCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SubnetCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SubnetCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *SubnetCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *SubnetCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *SubnetCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *SubnetCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetVpcId

`func (o *SubnetCreateRequest) GetVpcId() string`

GetVpcId returns the VpcId field if non-nil, zero value otherwise.

### GetVpcIdOk

`func (o *SubnetCreateRequest) GetVpcIdOk() (*string, bool)`

GetVpcIdOk returns a tuple with the VpcId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpcId

`func (o *SubnetCreateRequest) SetVpcId(v string)`

SetVpcId sets VpcId field to given value.


### GetIpv4BlockId

`func (o *SubnetCreateRequest) GetIpv4BlockId() string`

GetIpv4BlockId returns the Ipv4BlockId field if non-nil, zero value otherwise.

### GetIpv4BlockIdOk

`func (o *SubnetCreateRequest) GetIpv4BlockIdOk() (*string, bool)`

GetIpv4BlockIdOk returns a tuple with the Ipv4BlockId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4BlockId

`func (o *SubnetCreateRequest) SetIpv4BlockId(v string)`

SetIpv4BlockId sets Ipv4BlockId field to given value.

### HasIpv4BlockId

`func (o *SubnetCreateRequest) HasIpv4BlockId() bool`

HasIpv4BlockId returns a boolean if a field has been set.

### GetIpv6BlockId

`func (o *SubnetCreateRequest) GetIpv6BlockId() string`

GetIpv6BlockId returns the Ipv6BlockId field if non-nil, zero value otherwise.

### GetIpv6BlockIdOk

`func (o *SubnetCreateRequest) GetIpv6BlockIdOk() (*string, bool)`

GetIpv6BlockIdOk returns a tuple with the Ipv6BlockId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6BlockId

`func (o *SubnetCreateRequest) SetIpv6BlockId(v string)`

SetIpv6BlockId sets Ipv6BlockId field to given value.

### HasIpv6BlockId

`func (o *SubnetCreateRequest) HasIpv6BlockId() bool`

HasIpv6BlockId returns a boolean if a field has been set.

### GetPrefixLength

`func (o *SubnetCreateRequest) GetPrefixLength() int32`

GetPrefixLength returns the PrefixLength field if non-nil, zero value otherwise.

### GetPrefixLengthOk

`func (o *SubnetCreateRequest) GetPrefixLengthOk() (*int32, bool)`

GetPrefixLengthOk returns a tuple with the PrefixLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefixLength

`func (o *SubnetCreateRequest) SetPrefixLength(v int32)`

SetPrefixLength sets PrefixLength field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


