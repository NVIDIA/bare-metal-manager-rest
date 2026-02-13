# MachineHealthProbeSuccess

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Health probe identifier | [optional] 
**Target** | Pointer to **NullableString** | Specific component targeted by health probe | [optional] 

## Methods

### NewMachineHealthProbeSuccess

`func NewMachineHealthProbeSuccess() *MachineHealthProbeSuccess`

NewMachineHealthProbeSuccess instantiates a new MachineHealthProbeSuccess object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineHealthProbeSuccessWithDefaults

`func NewMachineHealthProbeSuccessWithDefaults() *MachineHealthProbeSuccess`

NewMachineHealthProbeSuccessWithDefaults instantiates a new MachineHealthProbeSuccess object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *MachineHealthProbeSuccess) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *MachineHealthProbeSuccess) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *MachineHealthProbeSuccess) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *MachineHealthProbeSuccess) HasId() bool`

HasId returns a boolean if a field has been set.

### GetTarget

`func (o *MachineHealthProbeSuccess) GetTarget() string`

GetTarget returns the Target field if non-nil, zero value otherwise.

### GetTargetOk

`func (o *MachineHealthProbeSuccess) GetTargetOk() (*string, bool)`

GetTargetOk returns a tuple with the Target field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTarget

`func (o *MachineHealthProbeSuccess) SetTarget(v string)`

SetTarget sets Target field to given value.

### HasTarget

`func (o *MachineHealthProbeSuccess) HasTarget() bool`

HasTarget returns a boolean if a field has been set.

### SetTargetNil

`func (o *MachineHealthProbeSuccess) SetTargetNil(b bool)`

 SetTargetNil sets the value for Target to be an explicit nil

### UnsetTarget
`func (o *MachineHealthProbeSuccess) UnsetTarget()`

UnsetTarget ensures that no value is present for Target, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


