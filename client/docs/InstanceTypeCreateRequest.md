# InstanceTypeCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**SiteId** | **string** |  | 
**ControllerMachineType** | Pointer to **string** |  | [optional] 
**MachineCapabilities** | Pointer to [**[]InstanceTypeCapabilityCreateRequest**](InstanceTypeCapabilityCreateRequest.md) |  | [optional] 

## Methods

### NewInstanceTypeCreateRequest

`func NewInstanceTypeCreateRequest(name string, siteId string, ) *InstanceTypeCreateRequest`

NewInstanceTypeCreateRequest instantiates a new InstanceTypeCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceTypeCreateRequestWithDefaults

`func NewInstanceTypeCreateRequestWithDefaults() *InstanceTypeCreateRequest`

NewInstanceTypeCreateRequestWithDefaults instantiates a new InstanceTypeCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *InstanceTypeCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *InstanceTypeCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *InstanceTypeCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *InstanceTypeCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *InstanceTypeCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *InstanceTypeCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *InstanceTypeCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *InstanceTypeCreateRequest) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *InstanceTypeCreateRequest) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *InstanceTypeCreateRequest) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.


### GetControllerMachineType

`func (o *InstanceTypeCreateRequest) GetControllerMachineType() string`

GetControllerMachineType returns the ControllerMachineType field if non-nil, zero value otherwise.

### GetControllerMachineTypeOk

`func (o *InstanceTypeCreateRequest) GetControllerMachineTypeOk() (*string, bool)`

GetControllerMachineTypeOk returns a tuple with the ControllerMachineType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerMachineType

`func (o *InstanceTypeCreateRequest) SetControllerMachineType(v string)`

SetControllerMachineType sets ControllerMachineType field to given value.

### HasControllerMachineType

`func (o *InstanceTypeCreateRequest) HasControllerMachineType() bool`

HasControllerMachineType returns a boolean if a field has been set.

### GetMachineCapabilities

`func (o *InstanceTypeCreateRequest) GetMachineCapabilities() []InstanceTypeCapabilityCreateRequest`

GetMachineCapabilities returns the MachineCapabilities field if non-nil, zero value otherwise.

### GetMachineCapabilitiesOk

`func (o *InstanceTypeCreateRequest) GetMachineCapabilitiesOk() (*[]InstanceTypeCapabilityCreateRequest, bool)`

GetMachineCapabilitiesOk returns a tuple with the MachineCapabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineCapabilities

`func (o *InstanceTypeCreateRequest) SetMachineCapabilities(v []InstanceTypeCapabilityCreateRequest)`

SetMachineCapabilities sets MachineCapabilities field to given value.

### HasMachineCapabilities

`func (o *InstanceTypeCreateRequest) HasMachineCapabilities() bool`

HasMachineCapabilities returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


