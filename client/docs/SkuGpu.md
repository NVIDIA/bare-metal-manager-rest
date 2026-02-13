# SkuGpu

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Vendor** | Pointer to **string** | Vendor of the GPU | [optional] 
**Model** | Pointer to **string** | Model of the GPU | [optional] 
**TotalMemory** | Pointer to **string** | Total memory of the GPU (e.g. \&quot;80GB HBM3\&quot;) | [optional] 
**Count** | Pointer to **int32** | Number of GPUs present | [optional] 

## Methods

### NewSkuGpu

`func NewSkuGpu() *SkuGpu`

NewSkuGpu instantiates a new SkuGpu object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSkuGpuWithDefaults

`func NewSkuGpuWithDefaults() *SkuGpu`

NewSkuGpuWithDefaults instantiates a new SkuGpu object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVendor

`func (o *SkuGpu) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *SkuGpu) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *SkuGpu) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *SkuGpu) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### GetModel

`func (o *SkuGpu) GetModel() string`

GetModel returns the Model field if non-nil, zero value otherwise.

### GetModelOk

`func (o *SkuGpu) GetModelOk() (*string, bool)`

GetModelOk returns a tuple with the Model field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModel

`func (o *SkuGpu) SetModel(v string)`

SetModel sets Model field to given value.

### HasModel

`func (o *SkuGpu) HasModel() bool`

HasModel returns a boolean if a field has been set.

### GetTotalMemory

`func (o *SkuGpu) GetTotalMemory() string`

GetTotalMemory returns the TotalMemory field if non-nil, zero value otherwise.

### GetTotalMemoryOk

`func (o *SkuGpu) GetTotalMemoryOk() (*string, bool)`

GetTotalMemoryOk returns a tuple with the TotalMemory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalMemory

`func (o *SkuGpu) SetTotalMemory(v string)`

SetTotalMemory sets TotalMemory field to given value.

### HasTotalMemory

`func (o *SkuGpu) HasTotalMemory() bool`

HasTotalMemory returns a boolean if a field has been set.

### GetCount

`func (o *SkuGpu) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *SkuGpu) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *SkuGpu) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *SkuGpu) HasCount() bool`

HasCount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


