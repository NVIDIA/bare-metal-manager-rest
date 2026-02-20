# InfiniBandInterfaceCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PartitionId** | Pointer to **string** | ID of the Partition the Interface should attach to | [optional] 
**Device** | Pointer to **string** | Name of the InfiniBand device to use | [optional] 
**Vendor** | Pointer to **NullableString** | Name of the InfiniBand device vendor, optional | [optional] 
**DeviceInstance** | Pointer to **int32** | Index of the device, used to identify which interface card to attache the Partition to | [optional] 
**IsPhysical** | Pointer to **bool** | Specifies whether this Partition should be attached to the Instance over physical interface | [optional] 
**VirtualFunctionId** | Pointer to **NullableInt32** | Must be specified if isPhysical is false | [optional] 

## Methods

### NewInfiniBandInterfaceCreateRequest

`func NewInfiniBandInterfaceCreateRequest() *InfiniBandInterfaceCreateRequest`

NewInfiniBandInterfaceCreateRequest instantiates a new InfiniBandInterfaceCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInfiniBandInterfaceCreateRequestWithDefaults

`func NewInfiniBandInterfaceCreateRequestWithDefaults() *InfiniBandInterfaceCreateRequest`

NewInfiniBandInterfaceCreateRequestWithDefaults instantiates a new InfiniBandInterfaceCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPartitionId

`func (o *InfiniBandInterfaceCreateRequest) GetPartitionId() string`

GetPartitionId returns the PartitionId field if non-nil, zero value otherwise.

### GetPartitionIdOk

`func (o *InfiniBandInterfaceCreateRequest) GetPartitionIdOk() (*string, bool)`

GetPartitionIdOk returns a tuple with the PartitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartitionId

`func (o *InfiniBandInterfaceCreateRequest) SetPartitionId(v string)`

SetPartitionId sets PartitionId field to given value.

### HasPartitionId

`func (o *InfiniBandInterfaceCreateRequest) HasPartitionId() bool`

HasPartitionId returns a boolean if a field has been set.

### GetDevice

`func (o *InfiniBandInterfaceCreateRequest) GetDevice() string`

GetDevice returns the Device field if non-nil, zero value otherwise.

### GetDeviceOk

`func (o *InfiniBandInterfaceCreateRequest) GetDeviceOk() (*string, bool)`

GetDeviceOk returns a tuple with the Device field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevice

`func (o *InfiniBandInterfaceCreateRequest) SetDevice(v string)`

SetDevice sets Device field to given value.

### HasDevice

`func (o *InfiniBandInterfaceCreateRequest) HasDevice() bool`

HasDevice returns a boolean if a field has been set.

### GetVendor

`func (o *InfiniBandInterfaceCreateRequest) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *InfiniBandInterfaceCreateRequest) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *InfiniBandInterfaceCreateRequest) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *InfiniBandInterfaceCreateRequest) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### SetVendorNil

`func (o *InfiniBandInterfaceCreateRequest) SetVendorNil(b bool)`

 SetVendorNil sets the value for Vendor to be an explicit nil

### UnsetVendor
`func (o *InfiniBandInterfaceCreateRequest) UnsetVendor()`

UnsetVendor ensures that no value is present for Vendor, not even an explicit nil
### GetDeviceInstance

`func (o *InfiniBandInterfaceCreateRequest) GetDeviceInstance() int32`

GetDeviceInstance returns the DeviceInstance field if non-nil, zero value otherwise.

### GetDeviceInstanceOk

`func (o *InfiniBandInterfaceCreateRequest) GetDeviceInstanceOk() (*int32, bool)`

GetDeviceInstanceOk returns a tuple with the DeviceInstance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceInstance

`func (o *InfiniBandInterfaceCreateRequest) SetDeviceInstance(v int32)`

SetDeviceInstance sets DeviceInstance field to given value.

### HasDeviceInstance

`func (o *InfiniBandInterfaceCreateRequest) HasDeviceInstance() bool`

HasDeviceInstance returns a boolean if a field has been set.

### GetIsPhysical

`func (o *InfiniBandInterfaceCreateRequest) GetIsPhysical() bool`

GetIsPhysical returns the IsPhysical field if non-nil, zero value otherwise.

### GetIsPhysicalOk

`func (o *InfiniBandInterfaceCreateRequest) GetIsPhysicalOk() (*bool, bool)`

GetIsPhysicalOk returns a tuple with the IsPhysical field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPhysical

`func (o *InfiniBandInterfaceCreateRequest) SetIsPhysical(v bool)`

SetIsPhysical sets IsPhysical field to given value.

### HasIsPhysical

`func (o *InfiniBandInterfaceCreateRequest) HasIsPhysical() bool`

HasIsPhysical returns a boolean if a field has been set.

### GetVirtualFunctionId

`func (o *InfiniBandInterfaceCreateRequest) GetVirtualFunctionId() int32`

GetVirtualFunctionId returns the VirtualFunctionId field if non-nil, zero value otherwise.

### GetVirtualFunctionIdOk

`func (o *InfiniBandInterfaceCreateRequest) GetVirtualFunctionIdOk() (*int32, bool)`

GetVirtualFunctionIdOk returns a tuple with the VirtualFunctionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVirtualFunctionId

`func (o *InfiniBandInterfaceCreateRequest) SetVirtualFunctionId(v int32)`

SetVirtualFunctionId sets VirtualFunctionId field to given value.

### HasVirtualFunctionId

`func (o *InfiniBandInterfaceCreateRequest) HasVirtualFunctionId() bool`

HasVirtualFunctionId returns a boolean if a field has been set.

### SetVirtualFunctionIdNil

`func (o *InfiniBandInterfaceCreateRequest) SetVirtualFunctionIdNil(b bool)`

 SetVirtualFunctionIdNil sets the value for VirtualFunctionId to be an explicit nil

### UnsetVirtualFunctionId
`func (o *InfiniBandInterfaceCreateRequest) UnsetVirtualFunctionId()`

UnsetVirtualFunctionId ensures that no value is present for VirtualFunctionId, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


