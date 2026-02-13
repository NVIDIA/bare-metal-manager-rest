# SkuStorage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Vendor** | Pointer to **string** | Vendor of the storage device | [optional] 
**Model** | Pointer to **string** | Model of the storage device | [optional] 
**CapacityMb** | Pointer to **int32** | Capacity in megabytes | [optional] 
**Count** | Pointer to **int32** | Number of storage devices present | [optional] 

## Methods

### NewSkuStorage

`func NewSkuStorage() *SkuStorage`

NewSkuStorage instantiates a new SkuStorage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSkuStorageWithDefaults

`func NewSkuStorageWithDefaults() *SkuStorage`

NewSkuStorageWithDefaults instantiates a new SkuStorage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVendor

`func (o *SkuStorage) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *SkuStorage) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *SkuStorage) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *SkuStorage) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### GetModel

`func (o *SkuStorage) GetModel() string`

GetModel returns the Model field if non-nil, zero value otherwise.

### GetModelOk

`func (o *SkuStorage) GetModelOk() (*string, bool)`

GetModelOk returns a tuple with the Model field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModel

`func (o *SkuStorage) SetModel(v string)`

SetModel sets Model field to given value.

### HasModel

`func (o *SkuStorage) HasModel() bool`

HasModel returns a boolean if a field has been set.

### GetCapacityMb

`func (o *SkuStorage) GetCapacityMb() int32`

GetCapacityMb returns the CapacityMb field if non-nil, zero value otherwise.

### GetCapacityMbOk

`func (o *SkuStorage) GetCapacityMbOk() (*int32, bool)`

GetCapacityMbOk returns a tuple with the CapacityMb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacityMb

`func (o *SkuStorage) SetCapacityMb(v int32)`

SetCapacityMb sets CapacityMb field to given value.

### HasCapacityMb

`func (o *SkuStorage) HasCapacityMb() bool`

HasCapacityMb returns a boolean if a field has been set.

### GetCount

`func (o *SkuStorage) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *SkuStorage) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *SkuStorage) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *SkuStorage) HasCount() bool`

HasCount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


