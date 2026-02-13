# Subnet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**SiteId** | Pointer to **string** |  | [optional] 
**VpcId** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **string** |  | [optional] 
**ControllerNetworkSegmentId** | Pointer to **NullableString** |  | [optional] 
**Ipv4Prefix** | Pointer to **NullableString** | The prefix that gets assigned to the subnet if ipv4 block is chosen | [optional] 
**Ipv4BlockId** | Pointer to **NullableString** |  | [optional] 
**Ipv4Gateway** | Pointer to **NullableString** |  | [optional] 
**Ipv6Prefix** | Pointer to **NullableString** |  | [optional] 
**Ipv6BlockId** | Pointer to **NullableString** |  | [optional] 
**Ipv6Gateway** | Pointer to **NullableString** |  | [optional] 
**Mtu** | Pointer to **int32** | Maximum Transmission Unit size in bytes. This property is system-determined and read-only. | [optional] [readonly] 
**PrefixLength** | Pointer to **int32** | Max value depends on prefix length of parent IP Block | [optional] 
**RoutingType** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**SubnetStatus**](SubnetStatus.md) |  | [optional] [readonly] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 
**Deprecations** | Pointer to [**[]Deprecation**](Deprecation.md) |  | [optional] 

## Methods

### NewSubnet

`func NewSubnet() *Subnet`

NewSubnet instantiates a new Subnet object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSubnetWithDefaults

`func NewSubnetWithDefaults() *Subnet`

NewSubnetWithDefaults instantiates a new Subnet object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Subnet) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Subnet) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Subnet) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Subnet) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *Subnet) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Subnet) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Subnet) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Subnet) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *Subnet) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Subnet) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Subnet) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Subnet) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *Subnet) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *Subnet) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *Subnet) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *Subnet) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetVpcId

`func (o *Subnet) GetVpcId() string`

GetVpcId returns the VpcId field if non-nil, zero value otherwise.

### GetVpcIdOk

`func (o *Subnet) GetVpcIdOk() (*string, bool)`

GetVpcIdOk returns a tuple with the VpcId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpcId

`func (o *Subnet) SetVpcId(v string)`

SetVpcId sets VpcId field to given value.

### HasVpcId

`func (o *Subnet) HasVpcId() bool`

HasVpcId returns a boolean if a field has been set.

### GetTenantId

`func (o *Subnet) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *Subnet) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *Subnet) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *Subnet) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetControllerNetworkSegmentId

`func (o *Subnet) GetControllerNetworkSegmentId() string`

GetControllerNetworkSegmentId returns the ControllerNetworkSegmentId field if non-nil, zero value otherwise.

### GetControllerNetworkSegmentIdOk

`func (o *Subnet) GetControllerNetworkSegmentIdOk() (*string, bool)`

GetControllerNetworkSegmentIdOk returns a tuple with the ControllerNetworkSegmentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerNetworkSegmentId

`func (o *Subnet) SetControllerNetworkSegmentId(v string)`

SetControllerNetworkSegmentId sets ControllerNetworkSegmentId field to given value.

### HasControllerNetworkSegmentId

`func (o *Subnet) HasControllerNetworkSegmentId() bool`

HasControllerNetworkSegmentId returns a boolean if a field has been set.

### SetControllerNetworkSegmentIdNil

`func (o *Subnet) SetControllerNetworkSegmentIdNil(b bool)`

 SetControllerNetworkSegmentIdNil sets the value for ControllerNetworkSegmentId to be an explicit nil

### UnsetControllerNetworkSegmentId
`func (o *Subnet) UnsetControllerNetworkSegmentId()`

UnsetControllerNetworkSegmentId ensures that no value is present for ControllerNetworkSegmentId, not even an explicit nil
### GetIpv4Prefix

`func (o *Subnet) GetIpv4Prefix() string`

GetIpv4Prefix returns the Ipv4Prefix field if non-nil, zero value otherwise.

### GetIpv4PrefixOk

`func (o *Subnet) GetIpv4PrefixOk() (*string, bool)`

GetIpv4PrefixOk returns a tuple with the Ipv4Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Prefix

`func (o *Subnet) SetIpv4Prefix(v string)`

SetIpv4Prefix sets Ipv4Prefix field to given value.

### HasIpv4Prefix

`func (o *Subnet) HasIpv4Prefix() bool`

HasIpv4Prefix returns a boolean if a field has been set.

### SetIpv4PrefixNil

`func (o *Subnet) SetIpv4PrefixNil(b bool)`

 SetIpv4PrefixNil sets the value for Ipv4Prefix to be an explicit nil

### UnsetIpv4Prefix
`func (o *Subnet) UnsetIpv4Prefix()`

UnsetIpv4Prefix ensures that no value is present for Ipv4Prefix, not even an explicit nil
### GetIpv4BlockId

`func (o *Subnet) GetIpv4BlockId() string`

GetIpv4BlockId returns the Ipv4BlockId field if non-nil, zero value otherwise.

### GetIpv4BlockIdOk

`func (o *Subnet) GetIpv4BlockIdOk() (*string, bool)`

GetIpv4BlockIdOk returns a tuple with the Ipv4BlockId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4BlockId

`func (o *Subnet) SetIpv4BlockId(v string)`

SetIpv4BlockId sets Ipv4BlockId field to given value.

### HasIpv4BlockId

`func (o *Subnet) HasIpv4BlockId() bool`

HasIpv4BlockId returns a boolean if a field has been set.

### SetIpv4BlockIdNil

`func (o *Subnet) SetIpv4BlockIdNil(b bool)`

 SetIpv4BlockIdNil sets the value for Ipv4BlockId to be an explicit nil

### UnsetIpv4BlockId
`func (o *Subnet) UnsetIpv4BlockId()`

UnsetIpv4BlockId ensures that no value is present for Ipv4BlockId, not even an explicit nil
### GetIpv4Gateway

`func (o *Subnet) GetIpv4Gateway() string`

GetIpv4Gateway returns the Ipv4Gateway field if non-nil, zero value otherwise.

### GetIpv4GatewayOk

`func (o *Subnet) GetIpv4GatewayOk() (*string, bool)`

GetIpv4GatewayOk returns a tuple with the Ipv4Gateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Gateway

`func (o *Subnet) SetIpv4Gateway(v string)`

SetIpv4Gateway sets Ipv4Gateway field to given value.

### HasIpv4Gateway

`func (o *Subnet) HasIpv4Gateway() bool`

HasIpv4Gateway returns a boolean if a field has been set.

### SetIpv4GatewayNil

`func (o *Subnet) SetIpv4GatewayNil(b bool)`

 SetIpv4GatewayNil sets the value for Ipv4Gateway to be an explicit nil

### UnsetIpv4Gateway
`func (o *Subnet) UnsetIpv4Gateway()`

UnsetIpv4Gateway ensures that no value is present for Ipv4Gateway, not even an explicit nil
### GetIpv6Prefix

`func (o *Subnet) GetIpv6Prefix() string`

GetIpv6Prefix returns the Ipv6Prefix field if non-nil, zero value otherwise.

### GetIpv6PrefixOk

`func (o *Subnet) GetIpv6PrefixOk() (*string, bool)`

GetIpv6PrefixOk returns a tuple with the Ipv6Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Prefix

`func (o *Subnet) SetIpv6Prefix(v string)`

SetIpv6Prefix sets Ipv6Prefix field to given value.

### HasIpv6Prefix

`func (o *Subnet) HasIpv6Prefix() bool`

HasIpv6Prefix returns a boolean if a field has been set.

### SetIpv6PrefixNil

`func (o *Subnet) SetIpv6PrefixNil(b bool)`

 SetIpv6PrefixNil sets the value for Ipv6Prefix to be an explicit nil

### UnsetIpv6Prefix
`func (o *Subnet) UnsetIpv6Prefix()`

UnsetIpv6Prefix ensures that no value is present for Ipv6Prefix, not even an explicit nil
### GetIpv6BlockId

`func (o *Subnet) GetIpv6BlockId() string`

GetIpv6BlockId returns the Ipv6BlockId field if non-nil, zero value otherwise.

### GetIpv6BlockIdOk

`func (o *Subnet) GetIpv6BlockIdOk() (*string, bool)`

GetIpv6BlockIdOk returns a tuple with the Ipv6BlockId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6BlockId

`func (o *Subnet) SetIpv6BlockId(v string)`

SetIpv6BlockId sets Ipv6BlockId field to given value.

### HasIpv6BlockId

`func (o *Subnet) HasIpv6BlockId() bool`

HasIpv6BlockId returns a boolean if a field has been set.

### SetIpv6BlockIdNil

`func (o *Subnet) SetIpv6BlockIdNil(b bool)`

 SetIpv6BlockIdNil sets the value for Ipv6BlockId to be an explicit nil

### UnsetIpv6BlockId
`func (o *Subnet) UnsetIpv6BlockId()`

UnsetIpv6BlockId ensures that no value is present for Ipv6BlockId, not even an explicit nil
### GetIpv6Gateway

`func (o *Subnet) GetIpv6Gateway() string`

GetIpv6Gateway returns the Ipv6Gateway field if non-nil, zero value otherwise.

### GetIpv6GatewayOk

`func (o *Subnet) GetIpv6GatewayOk() (*string, bool)`

GetIpv6GatewayOk returns a tuple with the Ipv6Gateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Gateway

`func (o *Subnet) SetIpv6Gateway(v string)`

SetIpv6Gateway sets Ipv6Gateway field to given value.

### HasIpv6Gateway

`func (o *Subnet) HasIpv6Gateway() bool`

HasIpv6Gateway returns a boolean if a field has been set.

### SetIpv6GatewayNil

`func (o *Subnet) SetIpv6GatewayNil(b bool)`

 SetIpv6GatewayNil sets the value for Ipv6Gateway to be an explicit nil

### UnsetIpv6Gateway
`func (o *Subnet) UnsetIpv6Gateway()`

UnsetIpv6Gateway ensures that no value is present for Ipv6Gateway, not even an explicit nil
### GetMtu

`func (o *Subnet) GetMtu() int32`

GetMtu returns the Mtu field if non-nil, zero value otherwise.

### GetMtuOk

`func (o *Subnet) GetMtuOk() (*int32, bool)`

GetMtuOk returns a tuple with the Mtu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMtu

`func (o *Subnet) SetMtu(v int32)`

SetMtu sets Mtu field to given value.

### HasMtu

`func (o *Subnet) HasMtu() bool`

HasMtu returns a boolean if a field has been set.

### GetPrefixLength

`func (o *Subnet) GetPrefixLength() int32`

GetPrefixLength returns the PrefixLength field if non-nil, zero value otherwise.

### GetPrefixLengthOk

`func (o *Subnet) GetPrefixLengthOk() (*int32, bool)`

GetPrefixLengthOk returns a tuple with the PrefixLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefixLength

`func (o *Subnet) SetPrefixLength(v int32)`

SetPrefixLength sets PrefixLength field to given value.

### HasPrefixLength

`func (o *Subnet) HasPrefixLength() bool`

HasPrefixLength returns a boolean if a field has been set.

### GetRoutingType

`func (o *Subnet) GetRoutingType() string`

GetRoutingType returns the RoutingType field if non-nil, zero value otherwise.

### GetRoutingTypeOk

`func (o *Subnet) GetRoutingTypeOk() (*string, bool)`

GetRoutingTypeOk returns a tuple with the RoutingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingType

`func (o *Subnet) SetRoutingType(v string)`

SetRoutingType sets RoutingType field to given value.

### HasRoutingType

`func (o *Subnet) HasRoutingType() bool`

HasRoutingType returns a boolean if a field has been set.

### GetStatus

`func (o *Subnet) GetStatus() SubnetStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Subnet) GetStatusOk() (*SubnetStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Subnet) SetStatus(v SubnetStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Subnet) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *Subnet) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *Subnet) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *Subnet) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *Subnet) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *Subnet) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *Subnet) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *Subnet) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *Subnet) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *Subnet) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *Subnet) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *Subnet) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *Subnet) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.

### GetDeprecations

`func (o *Subnet) GetDeprecations() []Deprecation`

GetDeprecations returns the Deprecations field if non-nil, zero value otherwise.

### GetDeprecationsOk

`func (o *Subnet) GetDeprecationsOk() (*[]Deprecation, bool)`

GetDeprecationsOk returns a tuple with the Deprecations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeprecations

`func (o *Subnet) SetDeprecations(v []Deprecation)`

SetDeprecations sets Deprecations field to given value.

### HasDeprecations

`func (o *Subnet) HasDeprecations() bool`

HasDeprecations returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


