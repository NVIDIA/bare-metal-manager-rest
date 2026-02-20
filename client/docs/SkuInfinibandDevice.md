# SkuInfinibandDevice

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Vendor** | Pointer to **string** | Vendor of the infiniband device | [optional] 
**Model** | Pointer to **string** | Model of the infiniband device | [optional] 
**Count** | Pointer to **int32** | Number of infiniband devices present | [optional] 

## Methods

### NewSkuInfinibandDevice

`func NewSkuInfinibandDevice() *SkuInfinibandDevice`

NewSkuInfinibandDevice instantiates a new SkuInfinibandDevice object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSkuInfinibandDeviceWithDefaults

`func NewSkuInfinibandDeviceWithDefaults() *SkuInfinibandDevice`

NewSkuInfinibandDeviceWithDefaults instantiates a new SkuInfinibandDevice object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVendor

`func (o *SkuInfinibandDevice) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *SkuInfinibandDevice) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *SkuInfinibandDevice) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *SkuInfinibandDevice) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### GetModel

`func (o *SkuInfinibandDevice) GetModel() string`

GetModel returns the Model field if non-nil, zero value otherwise.

### GetModelOk

`func (o *SkuInfinibandDevice) GetModelOk() (*string, bool)`

GetModelOk returns a tuple with the Model field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModel

`func (o *SkuInfinibandDevice) SetModel(v string)`

SetModel sets Model field to given value.

### HasModel

`func (o *SkuInfinibandDevice) HasModel() bool`

HasModel returns a boolean if a field has been set.

### GetCount

`func (o *SkuInfinibandDevice) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *SkuInfinibandDevice) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *SkuInfinibandDevice) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *SkuInfinibandDevice) HasCount() bool`

HasCount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


