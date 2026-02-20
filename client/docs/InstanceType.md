# InstanceType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**ControllerMachineType** | Pointer to **string** |  | [optional] 
**InfrastructureProviderId** | Pointer to **string** |  | [optional] 
**SiteId** | Pointer to **string** |  | [optional] 
**MachineCapabilities** | Pointer to [**[]MachineCapability**](MachineCapability.md) |  | [optional] 
**MachineInstanceTypes** | Pointer to [**[]MachineInstanceType**](MachineInstanceType.md) | Available only for Providers | [optional] 
**AllocationStats** | Pointer to [**InstanceTypeAllocationStats**](InstanceTypeAllocationStats.md) | Currently only available for Tenants | [optional] 
**Status** | Pointer to [**InstanceTypeStatus**](InstanceTypeStatus.md) |  | [optional] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**Deprecations** | Pointer to [**[]Deprecation**](Deprecation.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewInstanceType

`func NewInstanceType() *InstanceType`

NewInstanceType instantiates a new InstanceType object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceTypeWithDefaults

`func NewInstanceTypeWithDefaults() *InstanceType`

NewInstanceTypeWithDefaults instantiates a new InstanceType object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *InstanceType) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *InstanceType) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *InstanceType) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *InstanceType) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *InstanceType) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *InstanceType) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *InstanceType) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *InstanceType) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *InstanceType) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *InstanceType) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *InstanceType) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *InstanceType) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetControllerMachineType

`func (o *InstanceType) GetControllerMachineType() string`

GetControllerMachineType returns the ControllerMachineType field if non-nil, zero value otherwise.

### GetControllerMachineTypeOk

`func (o *InstanceType) GetControllerMachineTypeOk() (*string, bool)`

GetControllerMachineTypeOk returns a tuple with the ControllerMachineType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerMachineType

`func (o *InstanceType) SetControllerMachineType(v string)`

SetControllerMachineType sets ControllerMachineType field to given value.

### HasControllerMachineType

`func (o *InstanceType) HasControllerMachineType() bool`

HasControllerMachineType returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *InstanceType) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *InstanceType) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *InstanceType) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *InstanceType) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetSiteId

`func (o *InstanceType) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *InstanceType) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *InstanceType) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *InstanceType) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetMachineCapabilities

`func (o *InstanceType) GetMachineCapabilities() []MachineCapability`

GetMachineCapabilities returns the MachineCapabilities field if non-nil, zero value otherwise.

### GetMachineCapabilitiesOk

`func (o *InstanceType) GetMachineCapabilitiesOk() (*[]MachineCapability, bool)`

GetMachineCapabilitiesOk returns a tuple with the MachineCapabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineCapabilities

`func (o *InstanceType) SetMachineCapabilities(v []MachineCapability)`

SetMachineCapabilities sets MachineCapabilities field to given value.

### HasMachineCapabilities

`func (o *InstanceType) HasMachineCapabilities() bool`

HasMachineCapabilities returns a boolean if a field has been set.

### GetMachineInstanceTypes

`func (o *InstanceType) GetMachineInstanceTypes() []MachineInstanceType`

GetMachineInstanceTypes returns the MachineInstanceTypes field if non-nil, zero value otherwise.

### GetMachineInstanceTypesOk

`func (o *InstanceType) GetMachineInstanceTypesOk() (*[]MachineInstanceType, bool)`

GetMachineInstanceTypesOk returns a tuple with the MachineInstanceTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineInstanceTypes

`func (o *InstanceType) SetMachineInstanceTypes(v []MachineInstanceType)`

SetMachineInstanceTypes sets MachineInstanceTypes field to given value.

### HasMachineInstanceTypes

`func (o *InstanceType) HasMachineInstanceTypes() bool`

HasMachineInstanceTypes returns a boolean if a field has been set.

### GetAllocationStats

`func (o *InstanceType) GetAllocationStats() InstanceTypeAllocationStats`

GetAllocationStats returns the AllocationStats field if non-nil, zero value otherwise.

### GetAllocationStatsOk

`func (o *InstanceType) GetAllocationStatsOk() (*InstanceTypeAllocationStats, bool)`

GetAllocationStatsOk returns a tuple with the AllocationStats field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocationStats

`func (o *InstanceType) SetAllocationStats(v InstanceTypeAllocationStats)`

SetAllocationStats sets AllocationStats field to given value.

### HasAllocationStats

`func (o *InstanceType) HasAllocationStats() bool`

HasAllocationStats returns a boolean if a field has been set.

### GetStatus

`func (o *InstanceType) GetStatus() InstanceTypeStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *InstanceType) GetStatusOk() (*InstanceTypeStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *InstanceType) SetStatus(v InstanceTypeStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *InstanceType) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *InstanceType) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *InstanceType) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *InstanceType) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *InstanceType) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetDeprecations

`func (o *InstanceType) GetDeprecations() []Deprecation`

GetDeprecations returns the Deprecations field if non-nil, zero value otherwise.

### GetDeprecationsOk

`func (o *InstanceType) GetDeprecationsOk() (*[]Deprecation, bool)`

GetDeprecationsOk returns a tuple with the Deprecations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeprecations

`func (o *InstanceType) SetDeprecations(v []Deprecation)`

SetDeprecations sets Deprecations field to given value.

### HasDeprecations

`func (o *InstanceType) HasDeprecations() bool`

HasDeprecations returns a boolean if a field has been set.

### GetCreated

`func (o *InstanceType) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *InstanceType) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *InstanceType) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *InstanceType) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *InstanceType) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *InstanceType) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *InstanceType) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *InstanceType) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


