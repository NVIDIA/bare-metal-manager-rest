# SkuMemory

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CapacityMb** | Pointer to **int32** | Capacity in megabytes | [optional] 
**MemoryType** | Pointer to **string** | Type of memory (e.g. \&quot;DDR4\&quot;, \&quot;DDR5\&quot;) | [optional] 
**Count** | Pointer to **int32** | Number of memory modules present | [optional] 

## Methods

### NewSkuMemory

`func NewSkuMemory() *SkuMemory`

NewSkuMemory instantiates a new SkuMemory object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSkuMemoryWithDefaults

`func NewSkuMemoryWithDefaults() *SkuMemory`

NewSkuMemoryWithDefaults instantiates a new SkuMemory object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCapacityMb

`func (o *SkuMemory) GetCapacityMb() int32`

GetCapacityMb returns the CapacityMb field if non-nil, zero value otherwise.

### GetCapacityMbOk

`func (o *SkuMemory) GetCapacityMbOk() (*int32, bool)`

GetCapacityMbOk returns a tuple with the CapacityMb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacityMb

`func (o *SkuMemory) SetCapacityMb(v int32)`

SetCapacityMb sets CapacityMb field to given value.

### HasCapacityMb

`func (o *SkuMemory) HasCapacityMb() bool`

HasCapacityMb returns a boolean if a field has been set.

### GetMemoryType

`func (o *SkuMemory) GetMemoryType() string`

GetMemoryType returns the MemoryType field if non-nil, zero value otherwise.

### GetMemoryTypeOk

`func (o *SkuMemory) GetMemoryTypeOk() (*string, bool)`

GetMemoryTypeOk returns a tuple with the MemoryType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemoryType

`func (o *SkuMemory) SetMemoryType(v string)`

SetMemoryType sets MemoryType field to given value.

### HasMemoryType

`func (o *SkuMemory) HasMemoryType() bool`

HasMemoryType returns a boolean if a field has been set.

### GetCount

`func (o *SkuMemory) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *SkuMemory) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *SkuMemory) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *SkuMemory) HasCount() bool`

HasCount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


