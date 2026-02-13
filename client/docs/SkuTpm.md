# SkuTpm

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Vendor** | Pointer to **string** | Vendor of the TPM | [optional] 
**Model** | Pointer to **string** | Model of the TPM | [optional] 
**Count** | Pointer to **int32** | Number of TPMs present | [optional] 

## Methods

### NewSkuTpm

`func NewSkuTpm() *SkuTpm`

NewSkuTpm instantiates a new SkuTpm object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSkuTpmWithDefaults

`func NewSkuTpmWithDefaults() *SkuTpm`

NewSkuTpmWithDefaults instantiates a new SkuTpm object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVendor

`func (o *SkuTpm) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *SkuTpm) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *SkuTpm) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *SkuTpm) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### GetModel

`func (o *SkuTpm) GetModel() string`

GetModel returns the Model field if non-nil, zero value otherwise.

### GetModelOk

`func (o *SkuTpm) GetModelOk() (*string, bool)`

GetModelOk returns a tuple with the Model field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModel

`func (o *SkuTpm) SetModel(v string)`

SetModel sets Model field to given value.

### HasModel

`func (o *SkuTpm) HasModel() bool`

HasModel returns a boolean if a field has been set.

### GetCount

`func (o *SkuTpm) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *SkuTpm) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *SkuTpm) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *SkuTpm) HasCount() bool`

HasCount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


