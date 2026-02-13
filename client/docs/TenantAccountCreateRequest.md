# TenantAccountCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InfrastructureProviderId** | **string** |  | 
**TenantOrg** | **string** | Must be a valid Org name | 

## Methods

### NewTenantAccountCreateRequest

`func NewTenantAccountCreateRequest(infrastructureProviderId string, tenantOrg string, ) *TenantAccountCreateRequest`

NewTenantAccountCreateRequest instantiates a new TenantAccountCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTenantAccountCreateRequestWithDefaults

`func NewTenantAccountCreateRequestWithDefaults() *TenantAccountCreateRequest`

NewTenantAccountCreateRequestWithDefaults instantiates a new TenantAccountCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInfrastructureProviderId

`func (o *TenantAccountCreateRequest) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *TenantAccountCreateRequest) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *TenantAccountCreateRequest) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.


### GetTenantOrg

`func (o *TenantAccountCreateRequest) GetTenantOrg() string`

GetTenantOrg returns the TenantOrg field if non-nil, zero value otherwise.

### GetTenantOrgOk

`func (o *TenantAccountCreateRequest) GetTenantOrgOk() (*string, bool)`

GetTenantOrgOk returns a tuple with the TenantOrg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantOrg

`func (o *TenantAccountCreateRequest) SetTenantOrg(v string)`

SetTenantOrg sets TenantOrg field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


