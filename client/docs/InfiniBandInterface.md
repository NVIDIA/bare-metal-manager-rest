# InfiniBandInterface

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**InstanceId** | Pointer to **string** |  | [optional] 
**SiteId** | Pointer to **string** |  | [optional] 
**PartitionId** | Pointer to **string** | ID of the InfiniBand Partition associated with this interface | [optional] 
**Device** | Pointer to **string** | Name of the InfiniBand device associated with this interface | [optional] 
**DeviceInstance** | Pointer to **int32** |  | [optional] 
**IsPhysical** | Pointer to **bool** | Indicates whether this is a physical interface | [optional] 
**VirtualFunctionId** | Pointer to **NullableInt32** |  | [optional] 
**Guid** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to [**InfiniBandInterfaceStatus**](InfiniBandInterfaceStatus.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewInfiniBandInterface

`func NewInfiniBandInterface() *InfiniBandInterface`

NewInfiniBandInterface instantiates a new InfiniBandInterface object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInfiniBandInterfaceWithDefaults

`func NewInfiniBandInterfaceWithDefaults() *InfiniBandInterface`

NewInfiniBandInterfaceWithDefaults instantiates a new InfiniBandInterface object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *InfiniBandInterface) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *InfiniBandInterface) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *InfiniBandInterface) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *InfiniBandInterface) HasId() bool`

HasId returns a boolean if a field has been set.

### GetInstanceId

`func (o *InfiniBandInterface) GetInstanceId() string`

GetInstanceId returns the InstanceId field if non-nil, zero value otherwise.

### GetInstanceIdOk

`func (o *InfiniBandInterface) GetInstanceIdOk() (*string, bool)`

GetInstanceIdOk returns a tuple with the InstanceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceId

`func (o *InfiniBandInterface) SetInstanceId(v string)`

SetInstanceId sets InstanceId field to given value.

### HasInstanceId

`func (o *InfiniBandInterface) HasInstanceId() bool`

HasInstanceId returns a boolean if a field has been set.

### GetSiteId

`func (o *InfiniBandInterface) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *InfiniBandInterface) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *InfiniBandInterface) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *InfiniBandInterface) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetPartitionId

`func (o *InfiniBandInterface) GetPartitionId() string`

GetPartitionId returns the PartitionId field if non-nil, zero value otherwise.

### GetPartitionIdOk

`func (o *InfiniBandInterface) GetPartitionIdOk() (*string, bool)`

GetPartitionIdOk returns a tuple with the PartitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartitionId

`func (o *InfiniBandInterface) SetPartitionId(v string)`

SetPartitionId sets PartitionId field to given value.

### HasPartitionId

`func (o *InfiniBandInterface) HasPartitionId() bool`

HasPartitionId returns a boolean if a field has been set.

### GetDevice

`func (o *InfiniBandInterface) GetDevice() string`

GetDevice returns the Device field if non-nil, zero value otherwise.

### GetDeviceOk

`func (o *InfiniBandInterface) GetDeviceOk() (*string, bool)`

GetDeviceOk returns a tuple with the Device field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevice

`func (o *InfiniBandInterface) SetDevice(v string)`

SetDevice sets Device field to given value.

### HasDevice

`func (o *InfiniBandInterface) HasDevice() bool`

HasDevice returns a boolean if a field has been set.

### GetDeviceInstance

`func (o *InfiniBandInterface) GetDeviceInstance() int32`

GetDeviceInstance returns the DeviceInstance field if non-nil, zero value otherwise.

### GetDeviceInstanceOk

`func (o *InfiniBandInterface) GetDeviceInstanceOk() (*int32, bool)`

GetDeviceInstanceOk returns a tuple with the DeviceInstance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceInstance

`func (o *InfiniBandInterface) SetDeviceInstance(v int32)`

SetDeviceInstance sets DeviceInstance field to given value.

### HasDeviceInstance

`func (o *InfiniBandInterface) HasDeviceInstance() bool`

HasDeviceInstance returns a boolean if a field has been set.

### GetIsPhysical

`func (o *InfiniBandInterface) GetIsPhysical() bool`

GetIsPhysical returns the IsPhysical field if non-nil, zero value otherwise.

### GetIsPhysicalOk

`func (o *InfiniBandInterface) GetIsPhysicalOk() (*bool, bool)`

GetIsPhysicalOk returns a tuple with the IsPhysical field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPhysical

`func (o *InfiniBandInterface) SetIsPhysical(v bool)`

SetIsPhysical sets IsPhysical field to given value.

### HasIsPhysical

`func (o *InfiniBandInterface) HasIsPhysical() bool`

HasIsPhysical returns a boolean if a field has been set.

### GetVirtualFunctionId

`func (o *InfiniBandInterface) GetVirtualFunctionId() int32`

GetVirtualFunctionId returns the VirtualFunctionId field if non-nil, zero value otherwise.

### GetVirtualFunctionIdOk

`func (o *InfiniBandInterface) GetVirtualFunctionIdOk() (*int32, bool)`

GetVirtualFunctionIdOk returns a tuple with the VirtualFunctionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVirtualFunctionId

`func (o *InfiniBandInterface) SetVirtualFunctionId(v int32)`

SetVirtualFunctionId sets VirtualFunctionId field to given value.

### HasVirtualFunctionId

`func (o *InfiniBandInterface) HasVirtualFunctionId() bool`

HasVirtualFunctionId returns a boolean if a field has been set.

### SetVirtualFunctionIdNil

`func (o *InfiniBandInterface) SetVirtualFunctionIdNil(b bool)`

 SetVirtualFunctionIdNil sets the value for VirtualFunctionId to be an explicit nil

### UnsetVirtualFunctionId
`func (o *InfiniBandInterface) UnsetVirtualFunctionId()`

UnsetVirtualFunctionId ensures that no value is present for VirtualFunctionId, not even an explicit nil
### GetGuid

`func (o *InfiniBandInterface) GetGuid() string`

GetGuid returns the Guid field if non-nil, zero value otherwise.

### GetGuidOk

`func (o *InfiniBandInterface) GetGuidOk() (*string, bool)`

GetGuidOk returns a tuple with the Guid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGuid

`func (o *InfiniBandInterface) SetGuid(v string)`

SetGuid sets Guid field to given value.

### HasGuid

`func (o *InfiniBandInterface) HasGuid() bool`

HasGuid returns a boolean if a field has been set.

### SetGuidNil

`func (o *InfiniBandInterface) SetGuidNil(b bool)`

 SetGuidNil sets the value for Guid to be an explicit nil

### UnsetGuid
`func (o *InfiniBandInterface) UnsetGuid()`

UnsetGuid ensures that no value is present for Guid, not even an explicit nil
### GetStatus

`func (o *InfiniBandInterface) GetStatus() InfiniBandInterfaceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *InfiniBandInterface) GetStatusOk() (*InfiniBandInterfaceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *InfiniBandInterface) SetStatus(v InfiniBandInterfaceStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *InfiniBandInterface) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetCreated

`func (o *InfiniBandInterface) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *InfiniBandInterface) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *InfiniBandInterface) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *InfiniBandInterface) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *InfiniBandInterface) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *InfiniBandInterface) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *InfiniBandInterface) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *InfiniBandInterface) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


