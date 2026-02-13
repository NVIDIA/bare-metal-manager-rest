# NVLinkInterfaceCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NvLinklogicalPartitionId** | Pointer to **string** | ID of the NVLink Logical Partition the Interface should attach to | [optional] 
**DeviceInstance** | Pointer to **int32** | Index of the device, used to identify which GPU to attache the Partition to | [optional] 

## Methods

### NewNVLinkInterfaceCreateRequest

`func NewNVLinkInterfaceCreateRequest() *NVLinkInterfaceCreateRequest`

NewNVLinkInterfaceCreateRequest instantiates a new NVLinkInterfaceCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNVLinkInterfaceCreateRequestWithDefaults

`func NewNVLinkInterfaceCreateRequestWithDefaults() *NVLinkInterfaceCreateRequest`

NewNVLinkInterfaceCreateRequestWithDefaults instantiates a new NVLinkInterfaceCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNvLinklogicalPartitionId

`func (o *NVLinkInterfaceCreateRequest) GetNvLinklogicalPartitionId() string`

GetNvLinklogicalPartitionId returns the NvLinklogicalPartitionId field if non-nil, zero value otherwise.

### GetNvLinklogicalPartitionIdOk

`func (o *NVLinkInterfaceCreateRequest) GetNvLinklogicalPartitionIdOk() (*string, bool)`

GetNvLinklogicalPartitionIdOk returns a tuple with the NvLinklogicalPartitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNvLinklogicalPartitionId

`func (o *NVLinkInterfaceCreateRequest) SetNvLinklogicalPartitionId(v string)`

SetNvLinklogicalPartitionId sets NvLinklogicalPartitionId field to given value.

### HasNvLinklogicalPartitionId

`func (o *NVLinkInterfaceCreateRequest) HasNvLinklogicalPartitionId() bool`

HasNvLinklogicalPartitionId returns a boolean if a field has been set.

### GetDeviceInstance

`func (o *NVLinkInterfaceCreateRequest) GetDeviceInstance() int32`

GetDeviceInstance returns the DeviceInstance field if non-nil, zero value otherwise.

### GetDeviceInstanceOk

`func (o *NVLinkInterfaceCreateRequest) GetDeviceInstanceOk() (*int32, bool)`

GetDeviceInstanceOk returns a tuple with the DeviceInstance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceInstance

`func (o *NVLinkInterfaceCreateRequest) SetDeviceInstance(v int32)`

SetDeviceInstance sets DeviceInstance field to given value.

### HasDeviceInstance

`func (o *NVLinkInterfaceCreateRequest) HasDeviceInstance() bool`

HasDeviceInstance returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


