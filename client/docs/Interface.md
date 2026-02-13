# Interface

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**InstanceId** | Pointer to **string** |  | [optional] 
**SubnetId** | Pointer to **NullableString** |  | [optional] 
**VpcPrefixId** | Pointer to **NullableString** |  | [optional] 
**IsPhysical** | Pointer to **bool** |  | [optional] 
**Device** | Pointer to **NullableString** | Name of the device to use | [optional] 
**DeviceInstance** | Pointer to **NullableInt32** | Index of the device, used to identify which interface card to attache the Partition to | [optional] 
**VirtualFunctionId** | Pointer to **NullableInt32** | Must be specified if isPhysical is false | [optional] 
**MacAddress** | Pointer to **NullableString** |  | [optional] 
**IpAddresses** | Pointer to **[]string** | A list of IPv4 or IPv6 addresses | [optional] 
**Status** | Pointer to [**InterfaceStatus**](InterfaceStatus.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewInterface

`func NewInterface() *Interface`

NewInterface instantiates a new Interface object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInterfaceWithDefaults

`func NewInterfaceWithDefaults() *Interface`

NewInterfaceWithDefaults instantiates a new Interface object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Interface) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Interface) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Interface) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Interface) HasId() bool`

HasId returns a boolean if a field has been set.

### GetInstanceId

`func (o *Interface) GetInstanceId() string`

GetInstanceId returns the InstanceId field if non-nil, zero value otherwise.

### GetInstanceIdOk

`func (o *Interface) GetInstanceIdOk() (*string, bool)`

GetInstanceIdOk returns a tuple with the InstanceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceId

`func (o *Interface) SetInstanceId(v string)`

SetInstanceId sets InstanceId field to given value.

### HasInstanceId

`func (o *Interface) HasInstanceId() bool`

HasInstanceId returns a boolean if a field has been set.

### GetSubnetId

`func (o *Interface) GetSubnetId() string`

GetSubnetId returns the SubnetId field if non-nil, zero value otherwise.

### GetSubnetIdOk

`func (o *Interface) GetSubnetIdOk() (*string, bool)`

GetSubnetIdOk returns a tuple with the SubnetId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubnetId

`func (o *Interface) SetSubnetId(v string)`

SetSubnetId sets SubnetId field to given value.

### HasSubnetId

`func (o *Interface) HasSubnetId() bool`

HasSubnetId returns a boolean if a field has been set.

### SetSubnetIdNil

`func (o *Interface) SetSubnetIdNil(b bool)`

 SetSubnetIdNil sets the value for SubnetId to be an explicit nil

### UnsetSubnetId
`func (o *Interface) UnsetSubnetId()`

UnsetSubnetId ensures that no value is present for SubnetId, not even an explicit nil
### GetVpcPrefixId

`func (o *Interface) GetVpcPrefixId() string`

GetVpcPrefixId returns the VpcPrefixId field if non-nil, zero value otherwise.

### GetVpcPrefixIdOk

`func (o *Interface) GetVpcPrefixIdOk() (*string, bool)`

GetVpcPrefixIdOk returns a tuple with the VpcPrefixId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpcPrefixId

`func (o *Interface) SetVpcPrefixId(v string)`

SetVpcPrefixId sets VpcPrefixId field to given value.

### HasVpcPrefixId

`func (o *Interface) HasVpcPrefixId() bool`

HasVpcPrefixId returns a boolean if a field has been set.

### SetVpcPrefixIdNil

`func (o *Interface) SetVpcPrefixIdNil(b bool)`

 SetVpcPrefixIdNil sets the value for VpcPrefixId to be an explicit nil

### UnsetVpcPrefixId
`func (o *Interface) UnsetVpcPrefixId()`

UnsetVpcPrefixId ensures that no value is present for VpcPrefixId, not even an explicit nil
### GetIsPhysical

`func (o *Interface) GetIsPhysical() bool`

GetIsPhysical returns the IsPhysical field if non-nil, zero value otherwise.

### GetIsPhysicalOk

`func (o *Interface) GetIsPhysicalOk() (*bool, bool)`

GetIsPhysicalOk returns a tuple with the IsPhysical field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPhysical

`func (o *Interface) SetIsPhysical(v bool)`

SetIsPhysical sets IsPhysical field to given value.

### HasIsPhysical

`func (o *Interface) HasIsPhysical() bool`

HasIsPhysical returns a boolean if a field has been set.

### GetDevice

`func (o *Interface) GetDevice() string`

GetDevice returns the Device field if non-nil, zero value otherwise.

### GetDeviceOk

`func (o *Interface) GetDeviceOk() (*string, bool)`

GetDeviceOk returns a tuple with the Device field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevice

`func (o *Interface) SetDevice(v string)`

SetDevice sets Device field to given value.

### HasDevice

`func (o *Interface) HasDevice() bool`

HasDevice returns a boolean if a field has been set.

### SetDeviceNil

`func (o *Interface) SetDeviceNil(b bool)`

 SetDeviceNil sets the value for Device to be an explicit nil

### UnsetDevice
`func (o *Interface) UnsetDevice()`

UnsetDevice ensures that no value is present for Device, not even an explicit nil
### GetDeviceInstance

`func (o *Interface) GetDeviceInstance() int32`

GetDeviceInstance returns the DeviceInstance field if non-nil, zero value otherwise.

### GetDeviceInstanceOk

`func (o *Interface) GetDeviceInstanceOk() (*int32, bool)`

GetDeviceInstanceOk returns a tuple with the DeviceInstance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceInstance

`func (o *Interface) SetDeviceInstance(v int32)`

SetDeviceInstance sets DeviceInstance field to given value.

### HasDeviceInstance

`func (o *Interface) HasDeviceInstance() bool`

HasDeviceInstance returns a boolean if a field has been set.

### SetDeviceInstanceNil

`func (o *Interface) SetDeviceInstanceNil(b bool)`

 SetDeviceInstanceNil sets the value for DeviceInstance to be an explicit nil

### UnsetDeviceInstance
`func (o *Interface) UnsetDeviceInstance()`

UnsetDeviceInstance ensures that no value is present for DeviceInstance, not even an explicit nil
### GetVirtualFunctionId

`func (o *Interface) GetVirtualFunctionId() int32`

GetVirtualFunctionId returns the VirtualFunctionId field if non-nil, zero value otherwise.

### GetVirtualFunctionIdOk

`func (o *Interface) GetVirtualFunctionIdOk() (*int32, bool)`

GetVirtualFunctionIdOk returns a tuple with the VirtualFunctionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVirtualFunctionId

`func (o *Interface) SetVirtualFunctionId(v int32)`

SetVirtualFunctionId sets VirtualFunctionId field to given value.

### HasVirtualFunctionId

`func (o *Interface) HasVirtualFunctionId() bool`

HasVirtualFunctionId returns a boolean if a field has been set.

### SetVirtualFunctionIdNil

`func (o *Interface) SetVirtualFunctionIdNil(b bool)`

 SetVirtualFunctionIdNil sets the value for VirtualFunctionId to be an explicit nil

### UnsetVirtualFunctionId
`func (o *Interface) UnsetVirtualFunctionId()`

UnsetVirtualFunctionId ensures that no value is present for VirtualFunctionId, not even an explicit nil
### GetMacAddress

`func (o *Interface) GetMacAddress() string`

GetMacAddress returns the MacAddress field if non-nil, zero value otherwise.

### GetMacAddressOk

`func (o *Interface) GetMacAddressOk() (*string, bool)`

GetMacAddressOk returns a tuple with the MacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacAddress

`func (o *Interface) SetMacAddress(v string)`

SetMacAddress sets MacAddress field to given value.

### HasMacAddress

`func (o *Interface) HasMacAddress() bool`

HasMacAddress returns a boolean if a field has been set.

### SetMacAddressNil

`func (o *Interface) SetMacAddressNil(b bool)`

 SetMacAddressNil sets the value for MacAddress to be an explicit nil

### UnsetMacAddress
`func (o *Interface) UnsetMacAddress()`

UnsetMacAddress ensures that no value is present for MacAddress, not even an explicit nil
### GetIpAddresses

`func (o *Interface) GetIpAddresses() []string`

GetIpAddresses returns the IpAddresses field if non-nil, zero value otherwise.

### GetIpAddressesOk

`func (o *Interface) GetIpAddressesOk() (*[]string, bool)`

GetIpAddressesOk returns a tuple with the IpAddresses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpAddresses

`func (o *Interface) SetIpAddresses(v []string)`

SetIpAddresses sets IpAddresses field to given value.

### HasIpAddresses

`func (o *Interface) HasIpAddresses() bool`

HasIpAddresses returns a boolean if a field has been set.

### GetStatus

`func (o *Interface) GetStatus() InterfaceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Interface) GetStatusOk() (*InterfaceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Interface) SetStatus(v InterfaceStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Interface) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetCreated

`func (o *Interface) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *Interface) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *Interface) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *Interface) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *Interface) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *Interface) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *Interface) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *Interface) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


