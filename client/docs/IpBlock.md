# IpBlock

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**SiteId** | Pointer to **string** |  | [optional] 
**InfrastructureProviderId** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **NullableString** |  | [optional] 
**RoutingType** | Pointer to **string** |  | [optional] 
**Prefix** | Pointer to **string** | Either IPv4 or IPv6 address | [optional] 
**PrefixLength** | Pointer to **int32** | Min: 1, Max: 32 for ipv4, 128 for ipv6 | [optional] 
**ProtocolVersion** | Pointer to **string** |  | [optional] 
**UsageStats** | Pointer to [**IpBlockUsageStats**](IpBlockUsageStats.md) |  | [optional] 
**Status** | Pointer to [**IpBlockStatus**](IpBlockStatus.md) |  | [optional] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**Deprecations** | Pointer to [**[]Deprecation**](Deprecation.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewIpBlock

`func NewIpBlock() *IpBlock`

NewIpBlock instantiates a new IpBlock object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIpBlockWithDefaults

`func NewIpBlockWithDefaults() *IpBlock`

NewIpBlockWithDefaults instantiates a new IpBlock object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *IpBlock) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *IpBlock) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *IpBlock) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *IpBlock) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *IpBlock) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *IpBlock) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *IpBlock) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *IpBlock) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *IpBlock) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *IpBlock) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *IpBlock) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *IpBlock) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *IpBlock) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *IpBlock) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *IpBlock) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *IpBlock) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *IpBlock) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *IpBlock) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *IpBlock) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *IpBlock) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetTenantId

`func (o *IpBlock) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *IpBlock) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *IpBlock) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *IpBlock) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### SetTenantIdNil

`func (o *IpBlock) SetTenantIdNil(b bool)`

 SetTenantIdNil sets the value for TenantId to be an explicit nil

### UnsetTenantId
`func (o *IpBlock) UnsetTenantId()`

UnsetTenantId ensures that no value is present for TenantId, not even an explicit nil
### GetRoutingType

`func (o *IpBlock) GetRoutingType() string`

GetRoutingType returns the RoutingType field if non-nil, zero value otherwise.

### GetRoutingTypeOk

`func (o *IpBlock) GetRoutingTypeOk() (*string, bool)`

GetRoutingTypeOk returns a tuple with the RoutingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingType

`func (o *IpBlock) SetRoutingType(v string)`

SetRoutingType sets RoutingType field to given value.

### HasRoutingType

`func (o *IpBlock) HasRoutingType() bool`

HasRoutingType returns a boolean if a field has been set.

### GetPrefix

`func (o *IpBlock) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *IpBlock) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *IpBlock) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *IpBlock) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.

### GetPrefixLength

`func (o *IpBlock) GetPrefixLength() int32`

GetPrefixLength returns the PrefixLength field if non-nil, zero value otherwise.

### GetPrefixLengthOk

`func (o *IpBlock) GetPrefixLengthOk() (*int32, bool)`

GetPrefixLengthOk returns a tuple with the PrefixLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefixLength

`func (o *IpBlock) SetPrefixLength(v int32)`

SetPrefixLength sets PrefixLength field to given value.

### HasPrefixLength

`func (o *IpBlock) HasPrefixLength() bool`

HasPrefixLength returns a boolean if a field has been set.

### GetProtocolVersion

`func (o *IpBlock) GetProtocolVersion() string`

GetProtocolVersion returns the ProtocolVersion field if non-nil, zero value otherwise.

### GetProtocolVersionOk

`func (o *IpBlock) GetProtocolVersionOk() (*string, bool)`

GetProtocolVersionOk returns a tuple with the ProtocolVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocolVersion

`func (o *IpBlock) SetProtocolVersion(v string)`

SetProtocolVersion sets ProtocolVersion field to given value.

### HasProtocolVersion

`func (o *IpBlock) HasProtocolVersion() bool`

HasProtocolVersion returns a boolean if a field has been set.

### GetUsageStats

`func (o *IpBlock) GetUsageStats() IpBlockUsageStats`

GetUsageStats returns the UsageStats field if non-nil, zero value otherwise.

### GetUsageStatsOk

`func (o *IpBlock) GetUsageStatsOk() (*IpBlockUsageStats, bool)`

GetUsageStatsOk returns a tuple with the UsageStats field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsageStats

`func (o *IpBlock) SetUsageStats(v IpBlockUsageStats)`

SetUsageStats sets UsageStats field to given value.

### HasUsageStats

`func (o *IpBlock) HasUsageStats() bool`

HasUsageStats returns a boolean if a field has been set.

### GetStatus

`func (o *IpBlock) GetStatus() IpBlockStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *IpBlock) GetStatusOk() (*IpBlockStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *IpBlock) SetStatus(v IpBlockStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *IpBlock) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *IpBlock) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *IpBlock) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *IpBlock) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *IpBlock) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetDeprecations

`func (o *IpBlock) GetDeprecations() []Deprecation`

GetDeprecations returns the Deprecations field if non-nil, zero value otherwise.

### GetDeprecationsOk

`func (o *IpBlock) GetDeprecationsOk() (*[]Deprecation, bool)`

GetDeprecationsOk returns a tuple with the Deprecations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeprecations

`func (o *IpBlock) SetDeprecations(v []Deprecation)`

SetDeprecations sets Deprecations field to given value.

### HasDeprecations

`func (o *IpBlock) HasDeprecations() bool`

HasDeprecations returns a boolean if a field has been set.

### GetCreated

`func (o *IpBlock) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *IpBlock) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *IpBlock) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *IpBlock) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *IpBlock) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *IpBlock) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *IpBlock) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *IpBlock) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


