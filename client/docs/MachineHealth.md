# MachineHealth

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Source** | Pointer to **string** | Describes how the health report was generated | [optional] 
**ObservedAt** | Pointer to **NullableString** | Date/time when health report was generated | [optional] 
**Successes** | Pointer to [**[]MachineHealthProbeSuccess**](MachineHealthProbeSuccess.md) | Results from successful health probes for the Machine | [optional] 
**Alerts** | Pointer to [**[]MachineHealthProbeAlert**](MachineHealthProbeAlert.md) | Results from failed health probes for the Machine | [optional] 

## Methods

### NewMachineHealth

`func NewMachineHealth() *MachineHealth`

NewMachineHealth instantiates a new MachineHealth object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineHealthWithDefaults

`func NewMachineHealthWithDefaults() *MachineHealth`

NewMachineHealthWithDefaults instantiates a new MachineHealth object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSource

`func (o *MachineHealth) GetSource() string`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *MachineHealth) GetSourceOk() (*string, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *MachineHealth) SetSource(v string)`

SetSource sets Source field to given value.

### HasSource

`func (o *MachineHealth) HasSource() bool`

HasSource returns a boolean if a field has been set.

### GetObservedAt

`func (o *MachineHealth) GetObservedAt() string`

GetObservedAt returns the ObservedAt field if non-nil, zero value otherwise.

### GetObservedAtOk

`func (o *MachineHealth) GetObservedAtOk() (*string, bool)`

GetObservedAtOk returns a tuple with the ObservedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObservedAt

`func (o *MachineHealth) SetObservedAt(v string)`

SetObservedAt sets ObservedAt field to given value.

### HasObservedAt

`func (o *MachineHealth) HasObservedAt() bool`

HasObservedAt returns a boolean if a field has been set.

### SetObservedAtNil

`func (o *MachineHealth) SetObservedAtNil(b bool)`

 SetObservedAtNil sets the value for ObservedAt to be an explicit nil

### UnsetObservedAt
`func (o *MachineHealth) UnsetObservedAt()`

UnsetObservedAt ensures that no value is present for ObservedAt, not even an explicit nil
### GetSuccesses

`func (o *MachineHealth) GetSuccesses() []MachineHealthProbeSuccess`

GetSuccesses returns the Successes field if non-nil, zero value otherwise.

### GetSuccessesOk

`func (o *MachineHealth) GetSuccessesOk() (*[]MachineHealthProbeSuccess, bool)`

GetSuccessesOk returns a tuple with the Successes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccesses

`func (o *MachineHealth) SetSuccesses(v []MachineHealthProbeSuccess)`

SetSuccesses sets Successes field to given value.

### HasSuccesses

`func (o *MachineHealth) HasSuccesses() bool`

HasSuccesses returns a boolean if a field has been set.

### GetAlerts

`func (o *MachineHealth) GetAlerts() []MachineHealthProbeAlert`

GetAlerts returns the Alerts field if non-nil, zero value otherwise.

### GetAlertsOk

`func (o *MachineHealth) GetAlertsOk() (*[]MachineHealthProbeAlert, bool)`

GetAlertsOk returns a tuple with the Alerts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlerts

`func (o *MachineHealth) SetAlerts(v []MachineHealthProbeAlert)`

SetAlerts sets Alerts field to given value.

### HasAlerts

`func (o *MachineHealth) HasAlerts() bool`

HasAlerts returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


