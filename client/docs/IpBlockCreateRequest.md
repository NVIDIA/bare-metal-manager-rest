# IpBlockCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**SiteId** | **string** |  | 
**RoutingType** | **string** |  | 
**Prefix** | **string** | Either IPv4 or IPv6 address | 
**PrefixLength** | **int32** | Min: 1, Max: 32 for IPv4, 128 for IPv6 | 
**ProtocolVersion** | **string** |  | 

## Methods

### NewIpBlockCreateRequest

`func NewIpBlockCreateRequest(name string, siteId string, routingType string, prefix string, prefixLength int32, protocolVersion string, ) *IpBlockCreateRequest`

NewIpBlockCreateRequest instantiates a new IpBlockCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIpBlockCreateRequestWithDefaults

`func NewIpBlockCreateRequestWithDefaults() *IpBlockCreateRequest`

NewIpBlockCreateRequestWithDefaults instantiates a new IpBlockCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *IpBlockCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *IpBlockCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *IpBlockCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *IpBlockCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *IpBlockCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *IpBlockCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *IpBlockCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *IpBlockCreateRequest) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *IpBlockCreateRequest) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *IpBlockCreateRequest) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.


### GetRoutingType

`func (o *IpBlockCreateRequest) GetRoutingType() string`

GetRoutingType returns the RoutingType field if non-nil, zero value otherwise.

### GetRoutingTypeOk

`func (o *IpBlockCreateRequest) GetRoutingTypeOk() (*string, bool)`

GetRoutingTypeOk returns a tuple with the RoutingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingType

`func (o *IpBlockCreateRequest) SetRoutingType(v string)`

SetRoutingType sets RoutingType field to given value.


### GetPrefix

`func (o *IpBlockCreateRequest) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *IpBlockCreateRequest) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *IpBlockCreateRequest) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.


### GetPrefixLength

`func (o *IpBlockCreateRequest) GetPrefixLength() int32`

GetPrefixLength returns the PrefixLength field if non-nil, zero value otherwise.

### GetPrefixLengthOk

`func (o *IpBlockCreateRequest) GetPrefixLengthOk() (*int32, bool)`

GetPrefixLengthOk returns a tuple with the PrefixLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefixLength

`func (o *IpBlockCreateRequest) SetPrefixLength(v int32)`

SetPrefixLength sets PrefixLength field to given value.


### GetProtocolVersion

`func (o *IpBlockCreateRequest) GetProtocolVersion() string`

GetProtocolVersion returns the ProtocolVersion field if non-nil, zero value otherwise.

### GetProtocolVersionOk

`func (o *IpBlockCreateRequest) GetProtocolVersionOk() (*string, bool)`

GetProtocolVersionOk returns a tuple with the ProtocolVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocolVersion

`func (o *IpBlockCreateRequest) SetProtocolVersion(v string)`

SetProtocolVersion sets ProtocolVersion field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


