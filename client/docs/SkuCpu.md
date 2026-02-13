# SkuCpu

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Vendor** | Pointer to **string** | Vendor of the CPU | [optional] 
**Model** | Pointer to **string** | Model of the CPU | [optional] 
**ThreadCount** | Pointer to **int32** | Number of threads for the CPU | [optional] 
**Count** | Pointer to **int32** | Number of CPUs present | [optional] 

## Methods

### NewSkuCpu

`func NewSkuCpu() *SkuCpu`

NewSkuCpu instantiates a new SkuCpu object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSkuCpuWithDefaults

`func NewSkuCpuWithDefaults() *SkuCpu`

NewSkuCpuWithDefaults instantiates a new SkuCpu object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVendor

`func (o *SkuCpu) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *SkuCpu) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *SkuCpu) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *SkuCpu) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### GetModel

`func (o *SkuCpu) GetModel() string`

GetModel returns the Model field if non-nil, zero value otherwise.

### GetModelOk

`func (o *SkuCpu) GetModelOk() (*string, bool)`

GetModelOk returns a tuple with the Model field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModel

`func (o *SkuCpu) SetModel(v string)`

SetModel sets Model field to given value.

### HasModel

`func (o *SkuCpu) HasModel() bool`

HasModel returns a boolean if a field has been set.

### GetThreadCount

`func (o *SkuCpu) GetThreadCount() int32`

GetThreadCount returns the ThreadCount field if non-nil, zero value otherwise.

### GetThreadCountOk

`func (o *SkuCpu) GetThreadCountOk() (*int32, bool)`

GetThreadCountOk returns a tuple with the ThreadCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThreadCount

`func (o *SkuCpu) SetThreadCount(v int32)`

SetThreadCount sets ThreadCount field to given value.

### HasThreadCount

`func (o *SkuCpu) HasThreadCount() bool`

HasThreadCount returns a boolean if a field has been set.

### GetCount

`func (o *SkuCpu) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *SkuCpu) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *SkuCpu) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *SkuCpu) HasCount() bool`

HasCount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


