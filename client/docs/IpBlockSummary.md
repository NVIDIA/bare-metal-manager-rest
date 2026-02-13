# IpBlockSummary

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** |  | [optional] 
**RoutingType** | Pointer to **string** |  | [optional] 
**Prefix** | Pointer to **string** | Either IPv4 or IPv6 address | [optional] 
**PrefixLength** | Pointer to **int32** | Min: 1, Max: 32 for ipv4, 128 for ipv6 | [optional] 
**Status** | Pointer to [**IpBlockStatus**](IpBlockStatus.md) |  | [optional] 

## Methods

### NewIpBlockSummary

`func NewIpBlockSummary() *IpBlockSummary`

NewIpBlockSummary instantiates a new IpBlockSummary object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIpBlockSummaryWithDefaults

`func NewIpBlockSummaryWithDefaults() *IpBlockSummary`

NewIpBlockSummaryWithDefaults instantiates a new IpBlockSummary object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *IpBlockSummary) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *IpBlockSummary) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *IpBlockSummary) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *IpBlockSummary) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *IpBlockSummary) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *IpBlockSummary) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *IpBlockSummary) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *IpBlockSummary) HasName() bool`

HasName returns a boolean if a field has been set.

### GetRoutingType

`func (o *IpBlockSummary) GetRoutingType() string`

GetRoutingType returns the RoutingType field if non-nil, zero value otherwise.

### GetRoutingTypeOk

`func (o *IpBlockSummary) GetRoutingTypeOk() (*string, bool)`

GetRoutingTypeOk returns a tuple with the RoutingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingType

`func (o *IpBlockSummary) SetRoutingType(v string)`

SetRoutingType sets RoutingType field to given value.

### HasRoutingType

`func (o *IpBlockSummary) HasRoutingType() bool`

HasRoutingType returns a boolean if a field has been set.

### GetPrefix

`func (o *IpBlockSummary) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *IpBlockSummary) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *IpBlockSummary) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *IpBlockSummary) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.

### GetPrefixLength

`func (o *IpBlockSummary) GetPrefixLength() int32`

GetPrefixLength returns the PrefixLength field if non-nil, zero value otherwise.

### GetPrefixLengthOk

`func (o *IpBlockSummary) GetPrefixLengthOk() (*int32, bool)`

GetPrefixLengthOk returns a tuple with the PrefixLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefixLength

`func (o *IpBlockSummary) SetPrefixLength(v int32)`

SetPrefixLength sets PrefixLength field to given value.

### HasPrefixLength

`func (o *IpBlockSummary) HasPrefixLength() bool`

HasPrefixLength returns a boolean if a field has been set.

### GetStatus

`func (o *IpBlockSummary) GetStatus() IpBlockStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *IpBlockSummary) GetStatusOk() (*IpBlockStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *IpBlockSummary) SetStatus(v IpBlockStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *IpBlockSummary) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


