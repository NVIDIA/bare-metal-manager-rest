# InfiniBandPartition

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**SiteId** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **string** |  | [optional] 
**ControllerIBPartitionId** | Pointer to **NullableString** |  | [optional] 
**PartitionKey** | Pointer to **NullableString** |  | [optional] [readonly] 
**PartitionName** | Pointer to **NullableString** |  | [optional] [readonly] 
**ServiceLevel** | Pointer to **NullableInt32** |  | [optional] [readonly] 
**RateLimit** | Pointer to **NullableFloat32** |  | [optional] [readonly] 
**Mtu** | Pointer to **NullableInt32** |  | [optional] [readonly] 
**EnableSharp** | Pointer to **bool** |  | [optional] [readonly] 
**Status** | Pointer to [**InfiniBandPartitionStatus**](InfiniBandPartitionStatus.md) |  | [optional] [readonly] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewInfiniBandPartition

`func NewInfiniBandPartition() *InfiniBandPartition`

NewInfiniBandPartition instantiates a new InfiniBandPartition object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInfiniBandPartitionWithDefaults

`func NewInfiniBandPartitionWithDefaults() *InfiniBandPartition`

NewInfiniBandPartitionWithDefaults instantiates a new InfiniBandPartition object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *InfiniBandPartition) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *InfiniBandPartition) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *InfiniBandPartition) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *InfiniBandPartition) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *InfiniBandPartition) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *InfiniBandPartition) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *InfiniBandPartition) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *InfiniBandPartition) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *InfiniBandPartition) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *InfiniBandPartition) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *InfiniBandPartition) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *InfiniBandPartition) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *InfiniBandPartition) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *InfiniBandPartition) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *InfiniBandPartition) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *InfiniBandPartition) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetTenantId

`func (o *InfiniBandPartition) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *InfiniBandPartition) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *InfiniBandPartition) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *InfiniBandPartition) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetControllerIBPartitionId

`func (o *InfiniBandPartition) GetControllerIBPartitionId() string`

GetControllerIBPartitionId returns the ControllerIBPartitionId field if non-nil, zero value otherwise.

### GetControllerIBPartitionIdOk

`func (o *InfiniBandPartition) GetControllerIBPartitionIdOk() (*string, bool)`

GetControllerIBPartitionIdOk returns a tuple with the ControllerIBPartitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerIBPartitionId

`func (o *InfiniBandPartition) SetControllerIBPartitionId(v string)`

SetControllerIBPartitionId sets ControllerIBPartitionId field to given value.

### HasControllerIBPartitionId

`func (o *InfiniBandPartition) HasControllerIBPartitionId() bool`

HasControllerIBPartitionId returns a boolean if a field has been set.

### SetControllerIBPartitionIdNil

`func (o *InfiniBandPartition) SetControllerIBPartitionIdNil(b bool)`

 SetControllerIBPartitionIdNil sets the value for ControllerIBPartitionId to be an explicit nil

### UnsetControllerIBPartitionId
`func (o *InfiniBandPartition) UnsetControllerIBPartitionId()`

UnsetControllerIBPartitionId ensures that no value is present for ControllerIBPartitionId, not even an explicit nil
### GetPartitionKey

`func (o *InfiniBandPartition) GetPartitionKey() string`

GetPartitionKey returns the PartitionKey field if non-nil, zero value otherwise.

### GetPartitionKeyOk

`func (o *InfiniBandPartition) GetPartitionKeyOk() (*string, bool)`

GetPartitionKeyOk returns a tuple with the PartitionKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartitionKey

`func (o *InfiniBandPartition) SetPartitionKey(v string)`

SetPartitionKey sets PartitionKey field to given value.

### HasPartitionKey

`func (o *InfiniBandPartition) HasPartitionKey() bool`

HasPartitionKey returns a boolean if a field has been set.

### SetPartitionKeyNil

`func (o *InfiniBandPartition) SetPartitionKeyNil(b bool)`

 SetPartitionKeyNil sets the value for PartitionKey to be an explicit nil

### UnsetPartitionKey
`func (o *InfiniBandPartition) UnsetPartitionKey()`

UnsetPartitionKey ensures that no value is present for PartitionKey, not even an explicit nil
### GetPartitionName

`func (o *InfiniBandPartition) GetPartitionName() string`

GetPartitionName returns the PartitionName field if non-nil, zero value otherwise.

### GetPartitionNameOk

`func (o *InfiniBandPartition) GetPartitionNameOk() (*string, bool)`

GetPartitionNameOk returns a tuple with the PartitionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartitionName

`func (o *InfiniBandPartition) SetPartitionName(v string)`

SetPartitionName sets PartitionName field to given value.

### HasPartitionName

`func (o *InfiniBandPartition) HasPartitionName() bool`

HasPartitionName returns a boolean if a field has been set.

### SetPartitionNameNil

`func (o *InfiniBandPartition) SetPartitionNameNil(b bool)`

 SetPartitionNameNil sets the value for PartitionName to be an explicit nil

### UnsetPartitionName
`func (o *InfiniBandPartition) UnsetPartitionName()`

UnsetPartitionName ensures that no value is present for PartitionName, not even an explicit nil
### GetServiceLevel

`func (o *InfiniBandPartition) GetServiceLevel() int32`

GetServiceLevel returns the ServiceLevel field if non-nil, zero value otherwise.

### GetServiceLevelOk

`func (o *InfiniBandPartition) GetServiceLevelOk() (*int32, bool)`

GetServiceLevelOk returns a tuple with the ServiceLevel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceLevel

`func (o *InfiniBandPartition) SetServiceLevel(v int32)`

SetServiceLevel sets ServiceLevel field to given value.

### HasServiceLevel

`func (o *InfiniBandPartition) HasServiceLevel() bool`

HasServiceLevel returns a boolean if a field has been set.

### SetServiceLevelNil

`func (o *InfiniBandPartition) SetServiceLevelNil(b bool)`

 SetServiceLevelNil sets the value for ServiceLevel to be an explicit nil

### UnsetServiceLevel
`func (o *InfiniBandPartition) UnsetServiceLevel()`

UnsetServiceLevel ensures that no value is present for ServiceLevel, not even an explicit nil
### GetRateLimit

`func (o *InfiniBandPartition) GetRateLimit() float32`

GetRateLimit returns the RateLimit field if non-nil, zero value otherwise.

### GetRateLimitOk

`func (o *InfiniBandPartition) GetRateLimitOk() (*float32, bool)`

GetRateLimitOk returns a tuple with the RateLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRateLimit

`func (o *InfiniBandPartition) SetRateLimit(v float32)`

SetRateLimit sets RateLimit field to given value.

### HasRateLimit

`func (o *InfiniBandPartition) HasRateLimit() bool`

HasRateLimit returns a boolean if a field has been set.

### SetRateLimitNil

`func (o *InfiniBandPartition) SetRateLimitNil(b bool)`

 SetRateLimitNil sets the value for RateLimit to be an explicit nil

### UnsetRateLimit
`func (o *InfiniBandPartition) UnsetRateLimit()`

UnsetRateLimit ensures that no value is present for RateLimit, not even an explicit nil
### GetMtu

`func (o *InfiniBandPartition) GetMtu() int32`

GetMtu returns the Mtu field if non-nil, zero value otherwise.

### GetMtuOk

`func (o *InfiniBandPartition) GetMtuOk() (*int32, bool)`

GetMtuOk returns a tuple with the Mtu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMtu

`func (o *InfiniBandPartition) SetMtu(v int32)`

SetMtu sets Mtu field to given value.

### HasMtu

`func (o *InfiniBandPartition) HasMtu() bool`

HasMtu returns a boolean if a field has been set.

### SetMtuNil

`func (o *InfiniBandPartition) SetMtuNil(b bool)`

 SetMtuNil sets the value for Mtu to be an explicit nil

### UnsetMtu
`func (o *InfiniBandPartition) UnsetMtu()`

UnsetMtu ensures that no value is present for Mtu, not even an explicit nil
### GetEnableSharp

`func (o *InfiniBandPartition) GetEnableSharp() bool`

GetEnableSharp returns the EnableSharp field if non-nil, zero value otherwise.

### GetEnableSharpOk

`func (o *InfiniBandPartition) GetEnableSharpOk() (*bool, bool)`

GetEnableSharpOk returns a tuple with the EnableSharp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableSharp

`func (o *InfiniBandPartition) SetEnableSharp(v bool)`

SetEnableSharp sets EnableSharp field to given value.

### HasEnableSharp

`func (o *InfiniBandPartition) HasEnableSharp() bool`

HasEnableSharp returns a boolean if a field has been set.

### GetStatus

`func (o *InfiniBandPartition) GetStatus() InfiniBandPartitionStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *InfiniBandPartition) GetStatusOk() (*InfiniBandPartitionStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *InfiniBandPartition) SetStatus(v InfiniBandPartitionStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *InfiniBandPartition) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *InfiniBandPartition) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *InfiniBandPartition) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *InfiniBandPartition) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *InfiniBandPartition) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *InfiniBandPartition) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *InfiniBandPartition) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *InfiniBandPartition) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *InfiniBandPartition) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *InfiniBandPartition) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *InfiniBandPartition) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *InfiniBandPartition) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *InfiniBandPartition) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


