# InfrastructureProviderStats

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Machine** | Pointer to [**MachineCountByStatus**](MachineCountByStatus.md) |  | [optional] 
**IpBlock** | Pointer to [**IpBlockCountByStatus**](IpBlockCountByStatus.md) |  | [optional] 
**TenantAccount** | Pointer to [**TenantAccountCountByStatus**](TenantAccountCountByStatus.md) |  | [optional] 

## Methods

### NewInfrastructureProviderStats

`func NewInfrastructureProviderStats() *InfrastructureProviderStats`

NewInfrastructureProviderStats instantiates a new InfrastructureProviderStats object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInfrastructureProviderStatsWithDefaults

`func NewInfrastructureProviderStatsWithDefaults() *InfrastructureProviderStats`

NewInfrastructureProviderStatsWithDefaults instantiates a new InfrastructureProviderStats object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMachine

`func (o *InfrastructureProviderStats) GetMachine() MachineCountByStatus`

GetMachine returns the Machine field if non-nil, zero value otherwise.

### GetMachineOk

`func (o *InfrastructureProviderStats) GetMachineOk() (*MachineCountByStatus, bool)`

GetMachineOk returns a tuple with the Machine field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachine

`func (o *InfrastructureProviderStats) SetMachine(v MachineCountByStatus)`

SetMachine sets Machine field to given value.

### HasMachine

`func (o *InfrastructureProviderStats) HasMachine() bool`

HasMachine returns a boolean if a field has been set.

### GetIpBlock

`func (o *InfrastructureProviderStats) GetIpBlock() IpBlockCountByStatus`

GetIpBlock returns the IpBlock field if non-nil, zero value otherwise.

### GetIpBlockOk

`func (o *InfrastructureProviderStats) GetIpBlockOk() (*IpBlockCountByStatus, bool)`

GetIpBlockOk returns a tuple with the IpBlock field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpBlock

`func (o *InfrastructureProviderStats) SetIpBlock(v IpBlockCountByStatus)`

SetIpBlock sets IpBlock field to given value.

### HasIpBlock

`func (o *InfrastructureProviderStats) HasIpBlock() bool`

HasIpBlock returns a boolean if a field has been set.

### GetTenantAccount

`func (o *InfrastructureProviderStats) GetTenantAccount() TenantAccountCountByStatus`

GetTenantAccount returns the TenantAccount field if non-nil, zero value otherwise.

### GetTenantAccountOk

`func (o *InfrastructureProviderStats) GetTenantAccountOk() (*TenantAccountCountByStatus, bool)`

GetTenantAccountOk returns a tuple with the TenantAccount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantAccount

`func (o *InfrastructureProviderStats) SetTenantAccount(v TenantAccountCountByStatus)`

SetTenantAccount sets TenantAccount field to given value.

### HasTenantAccount

`func (o *InfrastructureProviderStats) HasTenantAccount() bool`

HasTenantAccount returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


