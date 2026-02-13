# InstanceTypeUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**MachineCapabilities** | Pointer to [**[]MachineCapability**](MachineCapability.md) |  | [optional] 

## Methods

### NewInstanceTypeUpdateRequest

`func NewInstanceTypeUpdateRequest() *InstanceTypeUpdateRequest`

NewInstanceTypeUpdateRequest instantiates a new InstanceTypeUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceTypeUpdateRequestWithDefaults

`func NewInstanceTypeUpdateRequestWithDefaults() *InstanceTypeUpdateRequest`

NewInstanceTypeUpdateRequestWithDefaults instantiates a new InstanceTypeUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *InstanceTypeUpdateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *InstanceTypeUpdateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *InstanceTypeUpdateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *InstanceTypeUpdateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *InstanceTypeUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *InstanceTypeUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *InstanceTypeUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *InstanceTypeUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetMachineCapabilities

`func (o *InstanceTypeUpdateRequest) GetMachineCapabilities() []MachineCapability`

GetMachineCapabilities returns the MachineCapabilities field if non-nil, zero value otherwise.

### GetMachineCapabilitiesOk

`func (o *InstanceTypeUpdateRequest) GetMachineCapabilitiesOk() (*[]MachineCapability, bool)`

GetMachineCapabilitiesOk returns a tuple with the MachineCapabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineCapabilities

`func (o *InstanceTypeUpdateRequest) SetMachineCapabilities(v []MachineCapability)`

SetMachineCapabilities sets MachineCapabilities field to given value.

### HasMachineCapabilities

`func (o *InstanceTypeUpdateRequest) HasMachineCapabilities() bool`

HasMachineCapabilities returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


