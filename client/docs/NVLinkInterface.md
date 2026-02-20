# NVLinkInterface

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**InstanceId** | Pointer to **string** |  | [optional] 
**NvLinkLogicalPartitionId** | Pointer to **string** | ID of the NVLink Logical Partition associated with this interface | [optional] 
**NvLinkDomainId** | Pointer to **string** | ID of the NVLink Domain associated with this Interface | [optional] 
**DeviceInstance** | Pointer to **int32** | Index of the device, used to identify the GPU associated with this Interface | [optional] 
**GpuGuid** | Pointer to **string** | Unique ID of the GPU | [optional] 
**Status** | Pointer to [**NVLinkInterfaceStatus**](NVLinkInterfaceStatus.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewNVLinkInterface

`func NewNVLinkInterface() *NVLinkInterface`

NewNVLinkInterface instantiates a new NVLinkInterface object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNVLinkInterfaceWithDefaults

`func NewNVLinkInterfaceWithDefaults() *NVLinkInterface`

NewNVLinkInterfaceWithDefaults instantiates a new NVLinkInterface object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *NVLinkInterface) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *NVLinkInterface) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *NVLinkInterface) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *NVLinkInterface) HasId() bool`

HasId returns a boolean if a field has been set.

### GetInstanceId

`func (o *NVLinkInterface) GetInstanceId() string`

GetInstanceId returns the InstanceId field if non-nil, zero value otherwise.

### GetInstanceIdOk

`func (o *NVLinkInterface) GetInstanceIdOk() (*string, bool)`

GetInstanceIdOk returns a tuple with the InstanceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceId

`func (o *NVLinkInterface) SetInstanceId(v string)`

SetInstanceId sets InstanceId field to given value.

### HasInstanceId

`func (o *NVLinkInterface) HasInstanceId() bool`

HasInstanceId returns a boolean if a field has been set.

### GetNvLinkLogicalPartitionId

`func (o *NVLinkInterface) GetNvLinkLogicalPartitionId() string`

GetNvLinkLogicalPartitionId returns the NvLinkLogicalPartitionId field if non-nil, zero value otherwise.

### GetNvLinkLogicalPartitionIdOk

`func (o *NVLinkInterface) GetNvLinkLogicalPartitionIdOk() (*string, bool)`

GetNvLinkLogicalPartitionIdOk returns a tuple with the NvLinkLogicalPartitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinkLogicalPartitionId

`func (o *NVLinkInterface) SetNvLinkLogicalPartitionId(v string)`

SetNvLinkLogicalPartitionId sets NvLinkLogicalPartitionId field to given value.

### HasNvLinkLogicalPartitionId

`func (o *NVLinkInterface) HasNvLinkLogicalPartitionId() bool`

HasNvLinkLogicalPartitionId returns a boolean if a field has been set.

### GetNvLinkDomainId

`func (o *NVLinkInterface) GetNvLinkDomainId() string`

GetNvLinkDomainId returns the NvLinkDomainId field if non-nil, zero value otherwise.

### GetNvLinkDomainIdOk

`func (o *NVLinkInterface) GetNvLinkDomainIdOk() (*string, bool)`

GetNvLinkDomainIdOk returns a tuple with the NvLinkDomainId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinkDomainId

`func (o *NVLinkInterface) SetNvLinkDomainId(v string)`

SetNvLinkDomainId sets NvLinkDomainId field to given value.

### HasNvLinkDomainId

`func (o *NVLinkInterface) HasNvLinkDomainId() bool`

HasNvLinkDomainId returns a boolean if a field has been set.

### GetDeviceInstance

`func (o *NVLinkInterface) GetDeviceInstance() int32`

GetDeviceInstance returns the DeviceInstance field if non-nil, zero value otherwise.

### GetDeviceInstanceOk

`func (o *NVLinkInterface) GetDeviceInstanceOk() (*int32, bool)`

GetDeviceInstanceOk returns a tuple with the DeviceInstance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceInstance

`func (o *NVLinkInterface) SetDeviceInstance(v int32)`

SetDeviceInstance sets DeviceInstance field to given value.

### HasDeviceInstance

`func (o *NVLinkInterface) HasDeviceInstance() bool`

HasDeviceInstance returns a boolean if a field has been set.

### GetGpuGuid

`func (o *NVLinkInterface) GetGpuGuid() string`

GetGpuGuid returns the GpuGuid field if non-nil, zero value otherwise.

### GetGpuGuidOk

`func (o *NVLinkInterface) GetGpuGuidOk() (*string, bool)`

GetGpuGuidOk returns a tuple with the GpuGuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGpuGuid

`func (o *NVLinkInterface) SetGpuGuid(v string)`

SetGpuGuid sets GpuGuid field to given value.

### HasGpuGuid

`func (o *NVLinkInterface) HasGpuGuid() bool`

HasGpuGuid returns a boolean if a field has been set.

### GetStatus

`func (o *NVLinkInterface) GetStatus() NVLinkInterfaceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *NVLinkInterface) GetStatusOk() (*NVLinkInterfaceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *NVLinkInterface) SetStatus(v NVLinkInterfaceStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *NVLinkInterface) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetCreated

`func (o *NVLinkInterface) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *NVLinkInterface) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *NVLinkInterface) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *NVLinkInterface) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *NVLinkInterface) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *NVLinkInterface) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *NVLinkInterface) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *NVLinkInterface) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


