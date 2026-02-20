# TenantAccount

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**InfrastructureProviderId** | Pointer to **string** |  | [optional] 
**InfrastructureProviderOrg** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **NullableString** |  | [optional] 
**TenantOrg** | Pointer to **NullableString** |  | [optional] 
**TenantContact** | Pointer to [**User**](User.md) |  | [optional] 
**AllocationCount** | Pointer to **int32** |  | [optional] 
**Status** | Pointer to [**TenantAccountStatus**](TenantAccountStatus.md) |  | [optional] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewTenantAccount

`func NewTenantAccount() *TenantAccount`

NewTenantAccount instantiates a new TenantAccount object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTenantAccountWithDefaults

`func NewTenantAccountWithDefaults() *TenantAccount`

NewTenantAccountWithDefaults instantiates a new TenantAccount object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *TenantAccount) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TenantAccount) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TenantAccount) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *TenantAccount) HasId() bool`

HasId returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *TenantAccount) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *TenantAccount) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *TenantAccount) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *TenantAccount) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetInfrastructureProviderOrg

`func (o *TenantAccount) GetInfrastructureProviderOrg() string`

GetInfrastructureProviderOrg returns the InfrastructureProviderOrg field if non-nil, zero value otherwise.

### GetInfrastructureProviderOrgOk

`func (o *TenantAccount) GetInfrastructureProviderOrgOk() (*string, bool)`

GetInfrastructureProviderOrgOk returns a tuple with the InfrastructureProviderOrg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderOrg

`func (o *TenantAccount) SetInfrastructureProviderOrg(v string)`

SetInfrastructureProviderOrg sets InfrastructureProviderOrg field to given value.

### HasInfrastructureProviderOrg

`func (o *TenantAccount) HasInfrastructureProviderOrg() bool`

HasInfrastructureProviderOrg returns a boolean if a field has been set.

### GetTenantId

`func (o *TenantAccount) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *TenantAccount) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *TenantAccount) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *TenantAccount) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### SetTenantIdNil

`func (o *TenantAccount) SetTenantIdNil(b bool)`

 SetTenantIdNil sets the value for TenantId to be an explicit nil

### UnsetTenantId
`func (o *TenantAccount) UnsetTenantId()`

UnsetTenantId ensures that no value is present for TenantId, not even an explicit nil
### GetTenantOrg

`func (o *TenantAccount) GetTenantOrg() string`

GetTenantOrg returns the TenantOrg field if non-nil, zero value otherwise.

### GetTenantOrgOk

`func (o *TenantAccount) GetTenantOrgOk() (*string, bool)`

GetTenantOrgOk returns a tuple with the TenantOrg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantOrg

`func (o *TenantAccount) SetTenantOrg(v string)`

SetTenantOrg sets TenantOrg field to given value.

### HasTenantOrg

`func (o *TenantAccount) HasTenantOrg() bool`

HasTenantOrg returns a boolean if a field has been set.

### SetTenantOrgNil

`func (o *TenantAccount) SetTenantOrgNil(b bool)`

 SetTenantOrgNil sets the value for TenantOrg to be an explicit nil

### UnsetTenantOrg
`func (o *TenantAccount) UnsetTenantOrg()`

UnsetTenantOrg ensures that no value is present for TenantOrg, not even an explicit nil
### GetTenantContact

`func (o *TenantAccount) GetTenantContact() User`

GetTenantContact returns the TenantContact field if non-nil, zero value otherwise.

### GetTenantContactOk

`func (o *TenantAccount) GetTenantContactOk() (*User, bool)`

GetTenantContactOk returns a tuple with the TenantContact field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantContact

`func (o *TenantAccount) SetTenantContact(v User)`

SetTenantContact sets TenantContact field to given value.

### HasTenantContact

`func (o *TenantAccount) HasTenantContact() bool`

HasTenantContact returns a boolean if a field has been set.

### GetAllocationCount

`func (o *TenantAccount) GetAllocationCount() int32`

GetAllocationCount returns the AllocationCount field if non-nil, zero value otherwise.

### GetAllocationCountOk

`func (o *TenantAccount) GetAllocationCountOk() (*int32, bool)`

GetAllocationCountOk returns a tuple with the AllocationCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocationCount

`func (o *TenantAccount) SetAllocationCount(v int32)`

SetAllocationCount sets AllocationCount field to given value.

### HasAllocationCount

`func (o *TenantAccount) HasAllocationCount() bool`

HasAllocationCount returns a boolean if a field has been set.

### GetStatus

`func (o *TenantAccount) GetStatus() TenantAccountStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *TenantAccount) GetStatusOk() (*TenantAccountStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *TenantAccount) SetStatus(v TenantAccountStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *TenantAccount) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *TenantAccount) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *TenantAccount) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *TenantAccount) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *TenantAccount) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *TenantAccount) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *TenantAccount) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *TenantAccount) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *TenantAccount) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *TenantAccount) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *TenantAccount) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *TenantAccount) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *TenantAccount) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


