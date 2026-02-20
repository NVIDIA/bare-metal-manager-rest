# ServiceAccount

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** | Indicates whether Service Account is enabled | [optional] 
**InfrastructureProviderId** | Pointer to **string** | ID of the Infrastructure Provider associated with Service Account | [optional] 
**TenantId** | Pointer to **string** | ID of the Tenant associated with Service Account | [optional] 

## Methods

### NewServiceAccount

`func NewServiceAccount() *ServiceAccount`

NewServiceAccount instantiates a new ServiceAccount object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServiceAccountWithDefaults

`func NewServiceAccountWithDefaults() *ServiceAccount`

NewServiceAccountWithDefaults instantiates a new ServiceAccount object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *ServiceAccount) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *ServiceAccount) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *ServiceAccount) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *ServiceAccount) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *ServiceAccount) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *ServiceAccount) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *ServiceAccount) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *ServiceAccount) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetTenantId

`func (o *ServiceAccount) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *ServiceAccount) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *ServiceAccount) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *ServiceAccount) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


