# Allocation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**InfrastructureProviderId** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **string** |  | [optional] 
**SiteId** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**AllocationStatus**](AllocationStatus.md) |  | [optional] [readonly] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**AllocationConstraints** | Pointer to [**[]AllocationConstraint**](AllocationConstraint.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewAllocation

`func NewAllocation() *Allocation`

NewAllocation instantiates a new Allocation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAllocationWithDefaults

`func NewAllocationWithDefaults() *Allocation`

NewAllocationWithDefaults instantiates a new Allocation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Allocation) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Allocation) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Allocation) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Allocation) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *Allocation) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Allocation) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Allocation) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Allocation) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *Allocation) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Allocation) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Allocation) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Allocation) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *Allocation) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *Allocation) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *Allocation) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *Allocation) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetTenantId

`func (o *Allocation) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *Allocation) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *Allocation) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *Allocation) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetSiteId

`func (o *Allocation) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *Allocation) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *Allocation) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *Allocation) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetStatus

`func (o *Allocation) GetStatus() AllocationStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Allocation) GetStatusOk() (*AllocationStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Allocation) SetStatus(v AllocationStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Allocation) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *Allocation) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *Allocation) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *Allocation) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *Allocation) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetAllocationConstraints

`func (o *Allocation) GetAllocationConstraints() []AllocationConstraint`

GetAllocationConstraints returns the AllocationConstraints field if non-nil, zero value otherwise.

### GetAllocationConstraintsOk

`func (o *Allocation) GetAllocationConstraintsOk() (*[]AllocationConstraint, bool)`

GetAllocationConstraintsOk returns a tuple with the AllocationConstraints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocationConstraints

`func (o *Allocation) SetAllocationConstraints(v []AllocationConstraint)`

SetAllocationConstraints sets AllocationConstraints field to given value.

### HasAllocationConstraints

`func (o *Allocation) HasAllocationConstraints() bool`

HasAllocationConstraints returns a boolean if a field has been set.

### GetCreated

`func (o *Allocation) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *Allocation) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *Allocation) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *Allocation) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *Allocation) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *Allocation) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *Allocation) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *Allocation) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


