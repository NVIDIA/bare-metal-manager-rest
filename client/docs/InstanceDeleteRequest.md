# InstanceDeleteRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MachineHealthIssue** | Pointer to [**MachineHealthIssue**](MachineHealthIssue.md) | Information regarding issue with the underlying Machine experienced by Tenant | [optional] 
**IsRepairTenant** | Pointer to **NullableBool** | Should be set to true for Tenants who are performing investigation/repairing the Machine. Otherwise omit or set to false | [optional] 

## Methods

### NewInstanceDeleteRequest

`func NewInstanceDeleteRequest() *InstanceDeleteRequest`

NewInstanceDeleteRequest instantiates a new InstanceDeleteRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceDeleteRequestWithDefaults

`func NewInstanceDeleteRequestWithDefaults() *InstanceDeleteRequest`

NewInstanceDeleteRequestWithDefaults instantiates a new InstanceDeleteRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMachineHealthIssue

`func (o *InstanceDeleteRequest) GetMachineHealthIssue() MachineHealthIssue`

GetMachineHealthIssue returns the MachineHealthIssue field if non-nil, zero value otherwise.

### GetMachineHealthIssueOk

`func (o *InstanceDeleteRequest) GetMachineHealthIssueOk() (*MachineHealthIssue, bool)`

GetMachineHealthIssueOk returns a tuple with the MachineHealthIssue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineHealthIssue

`func (o *InstanceDeleteRequest) SetMachineHealthIssue(v MachineHealthIssue)`

SetMachineHealthIssue sets MachineHealthIssue field to given value.

### HasMachineHealthIssue

`func (o *InstanceDeleteRequest) HasMachineHealthIssue() bool`

HasMachineHealthIssue returns a boolean if a field has been set.

### GetIsRepairTenant

`func (o *InstanceDeleteRequest) GetIsRepairTenant() bool`

GetIsRepairTenant returns the IsRepairTenant field if non-nil, zero value otherwise.

### GetIsRepairTenantOk

`func (o *InstanceDeleteRequest) GetIsRepairTenantOk() (*bool, bool)`

GetIsRepairTenantOk returns a tuple with the IsRepairTenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsRepairTenant

`func (o *InstanceDeleteRequest) SetIsRepairTenant(v bool)`

SetIsRepairTenant sets IsRepairTenant field to given value.

### HasIsRepairTenant

`func (o *InstanceDeleteRequest) HasIsRepairTenant() bool`

HasIsRepairTenant returns a boolean if a field has been set.

### SetIsRepairTenantNil

`func (o *InstanceDeleteRequest) SetIsRepairTenantNil(b bool)`

 SetIsRepairTenantNil sets the value for IsRepairTenant to be an explicit nil

### UnsetIsRepairTenant
`func (o *InstanceDeleteRequest) UnsetIsRepairTenant()`

UnsetIsRepairTenant ensures that no value is present for IsRepairTenant, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


