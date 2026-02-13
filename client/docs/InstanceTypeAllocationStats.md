# InstanceTypeAllocationStats

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Total** | Pointer to **int32** | Total number of Machines allocated to different Tenants for this Instance Type | [optional] 
**Used** | Pointer to **int32** | Total number of allocated Machines of this Instance Type currently being used by Tenants | [optional] 
**Unused** | Pointer to **int32** | Total number of allocated Machines of this Instance Type that is currently not being used by Tenants | [optional] 
**UnusedUsable** | Pointer to **int32** | Total number of allocated Machines of this Instance Type that is currently not in use but in Ready state, therefore can be provisioned by Tenant | [optional] 
**MaxAllocatable** | Pointer to **int32** | Maximum number of Machines of this Instance Type that can be allocated to a Tenant | [optional] 

## Methods

### NewInstanceTypeAllocationStats

`func NewInstanceTypeAllocationStats() *InstanceTypeAllocationStats`

NewInstanceTypeAllocationStats instantiates a new InstanceTypeAllocationStats object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceTypeAllocationStatsWithDefaults

`func NewInstanceTypeAllocationStatsWithDefaults() *InstanceTypeAllocationStats`

NewInstanceTypeAllocationStatsWithDefaults instantiates a new InstanceTypeAllocationStats object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTotal

`func (o *InstanceTypeAllocationStats) GetTotal() int32`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *InstanceTypeAllocationStats) GetTotalOk() (*int32, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *InstanceTypeAllocationStats) SetTotal(v int32)`

SetTotal sets Total field to given value.

### HasTotal

`func (o *InstanceTypeAllocationStats) HasTotal() bool`

HasTotal returns a boolean if a field has been set.

### GetUsed

`func (o *InstanceTypeAllocationStats) GetUsed() int32`

GetUsed returns the Used field if non-nil, zero value otherwise.

### GetUsedOk

`func (o *InstanceTypeAllocationStats) GetUsedOk() (*int32, bool)`

GetUsedOk returns a tuple with the Used field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsed

`func (o *InstanceTypeAllocationStats) SetUsed(v int32)`

SetUsed sets Used field to given value.

### HasUsed

`func (o *InstanceTypeAllocationStats) HasUsed() bool`

HasUsed returns a boolean if a field has been set.

### GetUnused

`func (o *InstanceTypeAllocationStats) GetUnused() int32`

GetUnused returns the Unused field if non-nil, zero value otherwise.

### GetUnusedOk

`func (o *InstanceTypeAllocationStats) GetUnusedOk() (*int32, bool)`

GetUnusedOk returns a tuple with the Unused field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnused

`func (o *InstanceTypeAllocationStats) SetUnused(v int32)`

SetUnused sets Unused field to given value.

### HasUnused

`func (o *InstanceTypeAllocationStats) HasUnused() bool`

HasUnused returns a boolean if a field has been set.

### GetUnusedUsable

`func (o *InstanceTypeAllocationStats) GetUnusedUsable() int32`

GetUnusedUsable returns the UnusedUsable field if non-nil, zero value otherwise.

### GetUnusedUsableOk

`func (o *InstanceTypeAllocationStats) GetUnusedUsableOk() (*int32, bool)`

GetUnusedUsableOk returns a tuple with the UnusedUsable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnusedUsable

`func (o *InstanceTypeAllocationStats) SetUnusedUsable(v int32)`

SetUnusedUsable sets UnusedUsable field to given value.

### HasUnusedUsable

`func (o *InstanceTypeAllocationStats) HasUnusedUsable() bool`

HasUnusedUsable returns a boolean if a field has been set.

### GetMaxAllocatable

`func (o *InstanceTypeAllocationStats) GetMaxAllocatable() int32`

GetMaxAllocatable returns the MaxAllocatable field if non-nil, zero value otherwise.

### GetMaxAllocatableOk

`func (o *InstanceTypeAllocationStats) GetMaxAllocatableOk() (*int32, bool)`

GetMaxAllocatableOk returns a tuple with the MaxAllocatable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxAllocatable

`func (o *InstanceTypeAllocationStats) SetMaxAllocatable(v int32)`

SetMaxAllocatable sets MaxAllocatable field to given value.

### HasMaxAllocatable

`func (o *InstanceTypeAllocationStats) HasMaxAllocatable() bool`

HasMaxAllocatable returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


