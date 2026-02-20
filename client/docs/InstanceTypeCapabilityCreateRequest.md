# InstanceTypeCapabilityCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to **string** | Type of the Capability | [optional] 
**Name** | Pointer to **string** | Name of the Capability component | [optional] 
**Frequency** | Pointer to **NullableString** | Frequency of the Capability component, if available | [optional] 
**Capacity** | Pointer to **NullableString** | Capacity of the Capability component, if applicable | [optional] 
**Vendor** | Pointer to **NullableString** | Vendor of the Capability component, if available | [optional] 
**Count** | Pointer to **NullableInt32** | Count of the Capability component | [optional] 
**InactiveDevices** | Pointer to **[]int32** | Indices of those devices that are inactive, only valid for InfiniBand Capability type | [optional] 
**DeviceType** | Pointer to **NullableString** | Device Type of the Capability component, if available | [optional] 

## Methods

### NewInstanceTypeCapabilityCreateRequest

`func NewInstanceTypeCapabilityCreateRequest() *InstanceTypeCapabilityCreateRequest`

NewInstanceTypeCapabilityCreateRequest instantiates a new InstanceTypeCapabilityCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceTypeCapabilityCreateRequestWithDefaults

`func NewInstanceTypeCapabilityCreateRequestWithDefaults() *InstanceTypeCapabilityCreateRequest`

NewInstanceTypeCapabilityCreateRequestWithDefaults instantiates a new InstanceTypeCapabilityCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *InstanceTypeCapabilityCreateRequest) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *InstanceTypeCapabilityCreateRequest) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *InstanceTypeCapabilityCreateRequest) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *InstanceTypeCapabilityCreateRequest) HasType() bool`

HasType returns a boolean if a field has been set.

### GetName

`func (o *InstanceTypeCapabilityCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *InstanceTypeCapabilityCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *InstanceTypeCapabilityCreateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *InstanceTypeCapabilityCreateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetFrequency

`func (o *InstanceTypeCapabilityCreateRequest) GetFrequency() string`

GetFrequency returns the Frequency field if non-nil, zero value otherwise.

### GetFrequencyOk

`func (o *InstanceTypeCapabilityCreateRequest) GetFrequencyOk() (*string, bool)`

GetFrequencyOk returns a tuple with the Frequency field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrequency

`func (o *InstanceTypeCapabilityCreateRequest) SetFrequency(v string)`

SetFrequency sets Frequency field to given value.

### HasFrequency

`func (o *InstanceTypeCapabilityCreateRequest) HasFrequency() bool`

HasFrequency returns a boolean if a field has been set.

### SetFrequencyNil

`func (o *InstanceTypeCapabilityCreateRequest) SetFrequencyNil(b bool)`

 SetFrequencyNil sets the value for Frequency to be an explicit nil

### UnsetFrequency
`func (o *InstanceTypeCapabilityCreateRequest) UnsetFrequency()`

UnsetFrequency ensures that no value is present for Frequency, not even an explicit nil
### GetCapacity

`func (o *InstanceTypeCapabilityCreateRequest) GetCapacity() string`

GetCapacity returns the Capacity field if non-nil, zero value otherwise.

### GetCapacityOk

`func (o *InstanceTypeCapabilityCreateRequest) GetCapacityOk() (*string, bool)`

GetCapacityOk returns a tuple with the Capacity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacity

`func (o *InstanceTypeCapabilityCreateRequest) SetCapacity(v string)`

SetCapacity sets Capacity field to given value.

### HasCapacity

`func (o *InstanceTypeCapabilityCreateRequest) HasCapacity() bool`

HasCapacity returns a boolean if a field has been set.

### SetCapacityNil

`func (o *InstanceTypeCapabilityCreateRequest) SetCapacityNil(b bool)`

 SetCapacityNil sets the value for Capacity to be an explicit nil

### UnsetCapacity
`func (o *InstanceTypeCapabilityCreateRequest) UnsetCapacity()`

UnsetCapacity ensures that no value is present for Capacity, not even an explicit nil
### GetVendor

`func (o *InstanceTypeCapabilityCreateRequest) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *InstanceTypeCapabilityCreateRequest) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *InstanceTypeCapabilityCreateRequest) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *InstanceTypeCapabilityCreateRequest) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### SetVendorNil

`func (o *InstanceTypeCapabilityCreateRequest) SetVendorNil(b bool)`

 SetVendorNil sets the value for Vendor to be an explicit nil

### UnsetVendor
`func (o *InstanceTypeCapabilityCreateRequest) UnsetVendor()`

UnsetVendor ensures that no value is present for Vendor, not even an explicit nil
### GetCount

`func (o *InstanceTypeCapabilityCreateRequest) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *InstanceTypeCapabilityCreateRequest) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *InstanceTypeCapabilityCreateRequest) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *InstanceTypeCapabilityCreateRequest) HasCount() bool`

HasCount returns a boolean if a field has been set.

### SetCountNil

`func (o *InstanceTypeCapabilityCreateRequest) SetCountNil(b bool)`

 SetCountNil sets the value for Count to be an explicit nil

### UnsetCount
`func (o *InstanceTypeCapabilityCreateRequest) UnsetCount()`

UnsetCount ensures that no value is present for Count, not even an explicit nil
### GetInactiveDevices

`func (o *InstanceTypeCapabilityCreateRequest) GetInactiveDevices() []int32`

GetInactiveDevices returns the InactiveDevices field if non-nil, zero value otherwise.

### GetInactiveDevicesOk

`func (o *InstanceTypeCapabilityCreateRequest) GetInactiveDevicesOk() (*[]int32, bool)`

GetInactiveDevicesOk returns a tuple with the InactiveDevices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInactiveDevices

`func (o *InstanceTypeCapabilityCreateRequest) SetInactiveDevices(v []int32)`

SetInactiveDevices sets InactiveDevices field to given value.

### HasInactiveDevices

`func (o *InstanceTypeCapabilityCreateRequest) HasInactiveDevices() bool`

HasInactiveDevices returns a boolean if a field has been set.

### GetDeviceType

`func (o *InstanceTypeCapabilityCreateRequest) GetDeviceType() string`

GetDeviceType returns the DeviceType field if non-nil, zero value otherwise.

### GetDeviceTypeOk

`func (o *InstanceTypeCapabilityCreateRequest) GetDeviceTypeOk() (*string, bool)`

GetDeviceTypeOk returns a tuple with the DeviceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceType

`func (o *InstanceTypeCapabilityCreateRequest) SetDeviceType(v string)`

SetDeviceType sets DeviceType field to given value.

### HasDeviceType

`func (o *InstanceTypeCapabilityCreateRequest) HasDeviceType() bool`

HasDeviceType returns a boolean if a field has been set.

### SetDeviceTypeNil

`func (o *InstanceTypeCapabilityCreateRequest) SetDeviceTypeNil(b bool)`

 SetDeviceTypeNil sets the value for DeviceType to be an explicit nil

### UnsetDeviceType
`func (o *InstanceTypeCapabilityCreateRequest) UnsetDeviceType()`

UnsetDeviceType ensures that no value is present for DeviceType, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


