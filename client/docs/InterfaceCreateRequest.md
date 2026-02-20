# InterfaceCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SubnetId** | Pointer to **string** |  | [optional] 
**VpcPrefixId** | Pointer to **string** |  | [optional] 
**IsPhysical** | Pointer to **bool** | Specifies whether this Subnet or VPC Prefix should be attached to the Instance over physical interface. | [optional] 
**Device** | Pointer to **string** | Name of the device to use | [optional] 
**DeviceInstance** | Pointer to **int32** | Index of the device, used to identify which interface card to attache the Partition to | [optional] 
**VirtualFunctionId** | Pointer to **NullableInt32** | Index of the virtual function to use, must be specified if isPhysical is false | [optional] 

## Methods

### NewInterfaceCreateRequest

`func NewInterfaceCreateRequest() *InterfaceCreateRequest`

NewInterfaceCreateRequest instantiates a new InterfaceCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInterfaceCreateRequestWithDefaults

`func NewInterfaceCreateRequestWithDefaults() *InterfaceCreateRequest`

NewInterfaceCreateRequestWithDefaults instantiates a new InterfaceCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSubnetId

`func (o *InterfaceCreateRequest) GetSubnetId() string`

GetSubnetId returns the SubnetId field if non-nil, zero value otherwise.

### GetSubnetIdOk

`func (o *InterfaceCreateRequest) GetSubnetIdOk() (*string, bool)`

GetSubnetIdOk returns a tuple with the SubnetId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubnetId

`func (o *InterfaceCreateRequest) SetSubnetId(v string)`

SetSubnetId sets SubnetId field to given value.

### HasSubnetId

`func (o *InterfaceCreateRequest) HasSubnetId() bool`

HasSubnetId returns a boolean if a field has been set.

### GetVpcPrefixId

`func (o *InterfaceCreateRequest) GetVpcPrefixId() string`

GetVpcPrefixId returns the VpcPrefixId field if non-nil, zero value otherwise.

### GetVpcPrefixIdOk

`func (o *InterfaceCreateRequest) GetVpcPrefixIdOk() (*string, bool)`

GetVpcPrefixIdOk returns a tuple with the VpcPrefixId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpcPrefixId

`func (o *InterfaceCreateRequest) SetVpcPrefixId(v string)`

SetVpcPrefixId sets VpcPrefixId field to given value.

### HasVpcPrefixId

`func (o *InterfaceCreateRequest) HasVpcPrefixId() bool`

HasVpcPrefixId returns a boolean if a field has been set.

### GetIsPhysical

`func (o *InterfaceCreateRequest) GetIsPhysical() bool`

GetIsPhysical returns the IsPhysical field if non-nil, zero value otherwise.

### GetIsPhysicalOk

`func (o *InterfaceCreateRequest) GetIsPhysicalOk() (*bool, bool)`

GetIsPhysicalOk returns a tuple with the IsPhysical field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPhysical

`func (o *InterfaceCreateRequest) SetIsPhysical(v bool)`

SetIsPhysical sets IsPhysical field to given value.

### HasIsPhysical

`func (o *InterfaceCreateRequest) HasIsPhysical() bool`

HasIsPhysical returns a boolean if a field has been set.

### GetDevice

`func (o *InterfaceCreateRequest) GetDevice() string`

GetDevice returns the Device field if non-nil, zero value otherwise.

### GetDeviceOk

`func (o *InterfaceCreateRequest) GetDeviceOk() (*string, bool)`

GetDeviceOk returns a tuple with the Device field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevice

`func (o *InterfaceCreateRequest) SetDevice(v string)`

SetDevice sets Device field to given value.

### HasDevice

`func (o *InterfaceCreateRequest) HasDevice() bool`

HasDevice returns a boolean if a field has been set.

### GetDeviceInstance

`func (o *InterfaceCreateRequest) GetDeviceInstance() int32`

GetDeviceInstance returns the DeviceInstance field if non-nil, zero value otherwise.

### GetDeviceInstanceOk

`func (o *InterfaceCreateRequest) GetDeviceInstanceOk() (*int32, bool)`

GetDeviceInstanceOk returns a tuple with the DeviceInstance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceInstance

`func (o *InterfaceCreateRequest) SetDeviceInstance(v int32)`

SetDeviceInstance sets DeviceInstance field to given value.

### HasDeviceInstance

`func (o *InterfaceCreateRequest) HasDeviceInstance() bool`

HasDeviceInstance returns a boolean if a field has been set.

### GetVirtualFunctionId

`func (o *InterfaceCreateRequest) GetVirtualFunctionId() int32`

GetVirtualFunctionId returns the VirtualFunctionId field if non-nil, zero value otherwise.

### GetVirtualFunctionIdOk

`func (o *InterfaceCreateRequest) GetVirtualFunctionIdOk() (*int32, bool)`

GetVirtualFunctionIdOk returns a tuple with the VirtualFunctionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVirtualFunctionId

`func (o *InterfaceCreateRequest) SetVirtualFunctionId(v int32)`

SetVirtualFunctionId sets VirtualFunctionId field to given value.

### HasVirtualFunctionId

`func (o *InterfaceCreateRequest) HasVirtualFunctionId() bool`

HasVirtualFunctionId returns a boolean if a field has been set.

### SetVirtualFunctionIdNil

`func (o *InterfaceCreateRequest) SetVirtualFunctionIdNil(b bool)`

 SetVirtualFunctionIdNil sets the value for VirtualFunctionId to be an explicit nil

### UnsetVirtualFunctionId
`func (o *InterfaceCreateRequest) UnsetVirtualFunctionId()`

UnsetVirtualFunctionId ensures that no value is present for VirtualFunctionId, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


