# MachineUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InstanceTypeId** | Pointer to **NullableString** | Update the Instance Type of the Machine. Cannot be specified when clearing Instance Type. Can only be updated by Provider. | [optional] 
**ClearInstanceType** | Pointer to **NullableBool** | Set to true to clear the existing Instance Type. Cannot be specified if Instance Type ID is specified. Can only be set by Provider. | [optional] 
**SetMaintenanceMode** | Pointer to **NullableBool** | Set to &#x60;true&#x60; to enable maintenance mode and to &#x60;false&#x60; to disable maintenance mode. Can be set by Provider or Privileged Tenant. | [optional] 
**MaintenanceMessage** | Pointer to **NullableString** | Optional message describing the reason for moving Machine into maintenance mode. Can be updated by Provider or Privileged Tenant. | [optional] 
**Labels** | Pointer to **map[string]string** | Machine labels will be overwritten, include existing labels to preserve them. Can be updated by Provider or Privileged Tenant. | [optional] 

## Methods

### NewMachineUpdateRequest

`func NewMachineUpdateRequest() *MachineUpdateRequest`

NewMachineUpdateRequest instantiates a new MachineUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineUpdateRequestWithDefaults

`func NewMachineUpdateRequestWithDefaults() *MachineUpdateRequest`

NewMachineUpdateRequestWithDefaults instantiates a new MachineUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInstanceTypeId

`func (o *MachineUpdateRequest) GetInstanceTypeId() string`

GetInstanceTypeId returns the InstanceTypeId field if non-nil, zero value otherwise.

### GetInstanceTypeIdOk

`func (o *MachineUpdateRequest) GetInstanceTypeIdOk() (*string, bool)`

GetInstanceTypeIdOk returns a tuple with the InstanceTypeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceTypeId

`func (o *MachineUpdateRequest) SetInstanceTypeId(v string)`

SetInstanceTypeId sets InstanceTypeId field to given value.

### HasInstanceTypeId

`func (o *MachineUpdateRequest) HasInstanceTypeId() bool`

HasInstanceTypeId returns a boolean if a field has been set.

### SetInstanceTypeIdNil

`func (o *MachineUpdateRequest) SetInstanceTypeIdNil(b bool)`

 SetInstanceTypeIdNil sets the value for InstanceTypeId to be an explicit nil

### UnsetInstanceTypeId
`func (o *MachineUpdateRequest) UnsetInstanceTypeId()`

UnsetInstanceTypeId ensures that no value is present for InstanceTypeId, not even an explicit nil
### GetClearInstanceType

`func (o *MachineUpdateRequest) GetClearInstanceType() bool`

GetClearInstanceType returns the ClearInstanceType field if non-nil, zero value otherwise.

### GetClearInstanceTypeOk

`func (o *MachineUpdateRequest) GetClearInstanceTypeOk() (*bool, bool)`

GetClearInstanceTypeOk returns a tuple with the ClearInstanceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClearInstanceType

`func (o *MachineUpdateRequest) SetClearInstanceType(v bool)`

SetClearInstanceType sets ClearInstanceType field to given value.

### HasClearInstanceType

`func (o *MachineUpdateRequest) HasClearInstanceType() bool`

HasClearInstanceType returns a boolean if a field has been set.

### SetClearInstanceTypeNil

`func (o *MachineUpdateRequest) SetClearInstanceTypeNil(b bool)`

 SetClearInstanceTypeNil sets the value for ClearInstanceType to be an explicit nil

### UnsetClearInstanceType
`func (o *MachineUpdateRequest) UnsetClearInstanceType()`

UnsetClearInstanceType ensures that no value is present for ClearInstanceType, not even an explicit nil
### GetSetMaintenanceMode

`func (o *MachineUpdateRequest) GetSetMaintenanceMode() bool`

GetSetMaintenanceMode returns the SetMaintenanceMode field if non-nil, zero value otherwise.

### GetSetMaintenanceModeOk

`func (o *MachineUpdateRequest) GetSetMaintenanceModeOk() (*bool, bool)`

GetSetMaintenanceModeOk returns a tuple with the SetMaintenanceMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSetMaintenanceMode

`func (o *MachineUpdateRequest) SetSetMaintenanceMode(v bool)`

SetSetMaintenanceMode sets SetMaintenanceMode field to given value.

### HasSetMaintenanceMode

`func (o *MachineUpdateRequest) HasSetMaintenanceMode() bool`

HasSetMaintenanceMode returns a boolean if a field has been set.

### SetSetMaintenanceModeNil

`func (o *MachineUpdateRequest) SetSetMaintenanceModeNil(b bool)`

 SetSetMaintenanceModeNil sets the value for SetMaintenanceMode to be an explicit nil

### UnsetSetMaintenanceMode
`func (o *MachineUpdateRequest) UnsetSetMaintenanceMode()`

UnsetSetMaintenanceMode ensures that no value is present for SetMaintenanceMode, not even an explicit nil
### GetMaintenanceMessage

`func (o *MachineUpdateRequest) GetMaintenanceMessage() string`

GetMaintenanceMessage returns the MaintenanceMessage field if non-nil, zero value otherwise.

### GetMaintenanceMessageOk

`func (o *MachineUpdateRequest) GetMaintenanceMessageOk() (*string, bool)`

GetMaintenanceMessageOk returns a tuple with the MaintenanceMessage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaintenanceMessage

`func (o *MachineUpdateRequest) SetMaintenanceMessage(v string)`

SetMaintenanceMessage sets MaintenanceMessage field to given value.

### HasMaintenanceMessage

`func (o *MachineUpdateRequest) HasMaintenanceMessage() bool`

HasMaintenanceMessage returns a boolean if a field has been set.

### SetMaintenanceMessageNil

`func (o *MachineUpdateRequest) SetMaintenanceMessageNil(b bool)`

 SetMaintenanceMessageNil sets the value for MaintenanceMessage to be an explicit nil

### UnsetMaintenanceMessage
`func (o *MachineUpdateRequest) UnsetMaintenanceMessage()`

UnsetMaintenanceMessage ensures that no value is present for MaintenanceMessage, not even an explicit nil
### GetLabels

`func (o *MachineUpdateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *MachineUpdateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *MachineUpdateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *MachineUpdateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


