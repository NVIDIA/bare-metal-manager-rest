# MachineHealthProbeAlert

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Health probe identifier | [optional] 
**Target** | Pointer to **NullableString** | Specific component targeted by health probe | [optional] 
**InAlertSince** | Pointer to **NullableString** | Date/time since the alert has been in effect | [optional] 
**Message** | Pointer to **string** | Details of the failed health probe result | [optional] 
**TenantMessage** | Pointer to **NullableString** | Information provided by Tenant, if any | [optional] 
**Classifications** | Pointer to **[]string** | Classifications for this alert, category or impact | [optional] 

## Methods

### NewMachineHealthProbeAlert

`func NewMachineHealthProbeAlert() *MachineHealthProbeAlert`

NewMachineHealthProbeAlert instantiates a new MachineHealthProbeAlert object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineHealthProbeAlertWithDefaults

`func NewMachineHealthProbeAlertWithDefaults() *MachineHealthProbeAlert`

NewMachineHealthProbeAlertWithDefaults instantiates a new MachineHealthProbeAlert object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *MachineHealthProbeAlert) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *MachineHealthProbeAlert) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *MachineHealthProbeAlert) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *MachineHealthProbeAlert) HasId() bool`

HasId returns a boolean if a field has been set.

### GetTarget

`func (o *MachineHealthProbeAlert) GetTarget() string`

GetTarget returns the Target field if non-nil, zero value otherwise.

### GetTargetOk

`func (o *MachineHealthProbeAlert) GetTargetOk() (*string, bool)`

GetTargetOk returns a tuple with the Target field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTarget

`func (o *MachineHealthProbeAlert) SetTarget(v string)`

SetTarget sets Target field to given value.

### HasTarget

`func (o *MachineHealthProbeAlert) HasTarget() bool`

HasTarget returns a boolean if a field has been set.

### SetTargetNil

`func (o *MachineHealthProbeAlert) SetTargetNil(b bool)`

 SetTargetNil sets the value for Target to be an explicit nil

### UnsetTarget
`func (o *MachineHealthProbeAlert) UnsetTarget()`

UnsetTarget ensures that no value is present for Target, not even an explicit nil
### GetInAlertSince

`func (o *MachineHealthProbeAlert) GetInAlertSince() string`

GetInAlertSince returns the InAlertSince field if non-nil, zero value otherwise.

### GetInAlertSinceOk

`func (o *MachineHealthProbeAlert) GetInAlertSinceOk() (*string, bool)`

GetInAlertSinceOk returns a tuple with the InAlertSince field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInAlertSince

`func (o *MachineHealthProbeAlert) SetInAlertSince(v string)`

SetInAlertSince sets InAlertSince field to given value.

### HasInAlertSince

`func (o *MachineHealthProbeAlert) HasInAlertSince() bool`

HasInAlertSince returns a boolean if a field has been set.

### SetInAlertSinceNil

`func (o *MachineHealthProbeAlert) SetInAlertSinceNil(b bool)`

 SetInAlertSinceNil sets the value for InAlertSince to be an explicit nil

### UnsetInAlertSince
`func (o *MachineHealthProbeAlert) UnsetInAlertSince()`

UnsetInAlertSince ensures that no value is present for InAlertSince, not even an explicit nil
### GetMessage

`func (o *MachineHealthProbeAlert) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *MachineHealthProbeAlert) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *MachineHealthProbeAlert) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *MachineHealthProbeAlert) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetTenantMessage

`func (o *MachineHealthProbeAlert) GetTenantMessage() string`

GetTenantMessage returns the TenantMessage field if non-nil, zero value otherwise.

### GetTenantMessageOk

`func (o *MachineHealthProbeAlert) GetTenantMessageOk() (*string, bool)`

GetTenantMessageOk returns a tuple with the TenantMessage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantMessage

`func (o *MachineHealthProbeAlert) SetTenantMessage(v string)`

SetTenantMessage sets TenantMessage field to given value.

### HasTenantMessage

`func (o *MachineHealthProbeAlert) HasTenantMessage() bool`

HasTenantMessage returns a boolean if a field has been set.

### SetTenantMessageNil

`func (o *MachineHealthProbeAlert) SetTenantMessageNil(b bool)`

 SetTenantMessageNil sets the value for TenantMessage to be an explicit nil

### UnsetTenantMessage
`func (o *MachineHealthProbeAlert) UnsetTenantMessage()`

UnsetTenantMessage ensures that no value is present for TenantMessage, not even an explicit nil
### GetClassifications

`func (o *MachineHealthProbeAlert) GetClassifications() []string`

GetClassifications returns the Classifications field if non-nil, zero value otherwise.

### GetClassificationsOk

`func (o *MachineHealthProbeAlert) GetClassificationsOk() (*[]string, bool)`

GetClassificationsOk returns a tuple with the Classifications field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClassifications

`func (o *MachineHealthProbeAlert) SetClassifications(v []string)`

SetClassifications sets Classifications field to given value.

### HasClassifications

`func (o *MachineHealthProbeAlert) HasClassifications() bool`

HasClassifications returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


