# MachineCapability

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to **string** | Type of the Capability | [optional] 
**Name** | Pointer to **string** | Name of the Capability component | [optional] 
**Frequency** | Pointer to **NullableString** | Frequency of the Capability component, if available | [optional] 
**Cores** | Pointer to **NullableInt32** | Number of Cores in the Capability component, if applicable | [optional] 
**Threads** | Pointer to **NullableInt32** | Number of Threads in the Capability component, if applicable | [optional] 
**Capacity** | Pointer to **NullableString** | Capacity of the Capability component, if applicable | [optional] 
**Vendor** | Pointer to **NullableString** | Vendor of the Capability component, if available | [optional] 
**InactiveDevices** | Pointer to **[]int32** | A list of inactive devices | [optional] 
**Count** | Pointer to **NullableInt32** | Count of the Capability component | [optional] 
**DeviceType** | Pointer to **NullableString** | Device Type of the Capability component, if available | [optional] 

## Methods

### NewMachineCapability

`func NewMachineCapability() *MachineCapability`

NewMachineCapability instantiates a new MachineCapability object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineCapabilityWithDefaults

`func NewMachineCapabilityWithDefaults() *MachineCapability`

NewMachineCapabilityWithDefaults instantiates a new MachineCapability object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *MachineCapability) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *MachineCapability) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *MachineCapability) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *MachineCapability) HasType() bool`

HasType returns a boolean if a field has been set.

### GetName

`func (o *MachineCapability) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *MachineCapability) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *MachineCapability) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *MachineCapability) HasName() bool`

HasName returns a boolean if a field has been set.

### GetFrequency

`func (o *MachineCapability) GetFrequency() string`

GetFrequency returns the Frequency field if non-nil, zero value otherwise.

### GetFrequencyOk

`func (o *MachineCapability) GetFrequencyOk() (*string, bool)`

GetFrequencyOk returns a tuple with the Frequency field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrequency

`func (o *MachineCapability) SetFrequency(v string)`

SetFrequency sets Frequency field to given value.

### HasFrequency

`func (o *MachineCapability) HasFrequency() bool`

HasFrequency returns a boolean if a field has been set.

### SetFrequencyNil

`func (o *MachineCapability) SetFrequencyNil(b bool)`

 SetFrequencyNil sets the value for Frequency to be an explicit nil

### UnsetFrequency
`func (o *MachineCapability) UnsetFrequency()`

UnsetFrequency ensures that no value is present for Frequency, not even an explicit nil
### GetCores

`func (o *MachineCapability) GetCores() int32`

GetCores returns the Cores field if non-nil, zero value otherwise.

### GetCoresOk

`func (o *MachineCapability) GetCoresOk() (*int32, bool)`

GetCoresOk returns a tuple with the Cores field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCores

`func (o *MachineCapability) SetCores(v int32)`

SetCores sets Cores field to given value.

### HasCores

`func (o *MachineCapability) HasCores() bool`

HasCores returns a boolean if a field has been set.

### SetCoresNil

`func (o *MachineCapability) SetCoresNil(b bool)`

 SetCoresNil sets the value for Cores to be an explicit nil

### UnsetCores
`func (o *MachineCapability) UnsetCores()`

UnsetCores ensures that no value is present for Cores, not even an explicit nil
### GetThreads

`func (o *MachineCapability) GetThreads() int32`

GetThreads returns the Threads field if non-nil, zero value otherwise.

### GetThreadsOk

`func (o *MachineCapability) GetThreadsOk() (*int32, bool)`

GetThreadsOk returns a tuple with the Threads field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThreads

`func (o *MachineCapability) SetThreads(v int32)`

SetThreads sets Threads field to given value.

### HasThreads

`func (o *MachineCapability) HasThreads() bool`

HasThreads returns a boolean if a field has been set.

### SetThreadsNil

`func (o *MachineCapability) SetThreadsNil(b bool)`

 SetThreadsNil sets the value for Threads to be an explicit nil

### UnsetThreads
`func (o *MachineCapability) UnsetThreads()`

UnsetThreads ensures that no value is present for Threads, not even an explicit nil
### GetCapacity

`func (o *MachineCapability) GetCapacity() string`

GetCapacity returns the Capacity field if non-nil, zero value otherwise.

### GetCapacityOk

`func (o *MachineCapability) GetCapacityOk() (*string, bool)`

GetCapacityOk returns a tuple with the Capacity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacity

`func (o *MachineCapability) SetCapacity(v string)`

SetCapacity sets Capacity field to given value.

### HasCapacity

`func (o *MachineCapability) HasCapacity() bool`

HasCapacity returns a boolean if a field has been set.

### SetCapacityNil

`func (o *MachineCapability) SetCapacityNil(b bool)`

 SetCapacityNil sets the value for Capacity to be an explicit nil

### UnsetCapacity
`func (o *MachineCapability) UnsetCapacity()`

UnsetCapacity ensures that no value is present for Capacity, not even an explicit nil
### GetVendor

`func (o *MachineCapability) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *MachineCapability) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *MachineCapability) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *MachineCapability) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### SetVendorNil

`func (o *MachineCapability) SetVendorNil(b bool)`

 SetVendorNil sets the value for Vendor to be an explicit nil

### UnsetVendor
`func (o *MachineCapability) UnsetVendor()`

UnsetVendor ensures that no value is present for Vendor, not even an explicit nil
### GetInactiveDevices

`func (o *MachineCapability) GetInactiveDevices() []int32`

GetInactiveDevices returns the InactiveDevices field if non-nil, zero value otherwise.

### GetInactiveDevicesOk

`func (o *MachineCapability) GetInactiveDevicesOk() (*[]int32, bool)`

GetInactiveDevicesOk returns a tuple with the InactiveDevices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInactiveDevices

`func (o *MachineCapability) SetInactiveDevices(v []int32)`

SetInactiveDevices sets InactiveDevices field to given value.

### HasInactiveDevices

`func (o *MachineCapability) HasInactiveDevices() bool`

HasInactiveDevices returns a boolean if a field has been set.

### GetCount

`func (o *MachineCapability) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *MachineCapability) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *MachineCapability) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *MachineCapability) HasCount() bool`

HasCount returns a boolean if a field has been set.

### SetCountNil

`func (o *MachineCapability) SetCountNil(b bool)`

 SetCountNil sets the value for Count to be an explicit nil

### UnsetCount
`func (o *MachineCapability) UnsetCount()`

UnsetCount ensures that no value is present for Count, not even an explicit nil
### GetDeviceType

`func (o *MachineCapability) GetDeviceType() string`

GetDeviceType returns the DeviceType field if non-nil, zero value otherwise.

### GetDeviceTypeOk

`func (o *MachineCapability) GetDeviceTypeOk() (*string, bool)`

GetDeviceTypeOk returns a tuple with the DeviceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceType

`func (o *MachineCapability) SetDeviceType(v string)`

SetDeviceType sets DeviceType field to given value.

### HasDeviceType

`func (o *MachineCapability) HasDeviceType() bool`

HasDeviceType returns a boolean if a field has been set.

### SetDeviceTypeNil

`func (o *MachineCapability) SetDeviceTypeNil(b bool)`

 SetDeviceTypeNil sets the value for DeviceType to be an explicit nil

### UnsetDeviceType
`func (o *MachineCapability) UnsetDeviceType()`

UnsetDeviceType ensures that no value is present for DeviceType, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


