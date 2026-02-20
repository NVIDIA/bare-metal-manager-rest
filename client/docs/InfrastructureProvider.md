# InfrastructureProvider

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Org** | Pointer to **string** |  | [optional] 
**OrgDisplayName** | Pointer to **string** |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewInfrastructureProvider

`func NewInfrastructureProvider() *InfrastructureProvider`

NewInfrastructureProvider instantiates a new InfrastructureProvider object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInfrastructureProviderWithDefaults

`func NewInfrastructureProviderWithDefaults() *InfrastructureProvider`

NewInfrastructureProviderWithDefaults instantiates a new InfrastructureProvider object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *InfrastructureProvider) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *InfrastructureProvider) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *InfrastructureProvider) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *InfrastructureProvider) HasId() bool`

HasId returns a boolean if a field has been set.

### GetOrg

`func (o *InfrastructureProvider) GetOrg() string`

GetOrg returns the Org field if non-nil, zero value otherwise.

### GetOrgOk

`func (o *InfrastructureProvider) GetOrgOk() (*string, bool)`

GetOrgOk returns a tuple with the Org field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrg

`func (o *InfrastructureProvider) SetOrg(v string)`

SetOrg sets Org field to given value.

### HasOrg

`func (o *InfrastructureProvider) HasOrg() bool`

HasOrg returns a boolean if a field has been set.

### GetOrgDisplayName

`func (o *InfrastructureProvider) GetOrgDisplayName() string`

GetOrgDisplayName returns the OrgDisplayName field if non-nil, zero value otherwise.

### GetOrgDisplayNameOk

`func (o *InfrastructureProvider) GetOrgDisplayNameOk() (*string, bool)`

GetOrgDisplayNameOk returns a tuple with the OrgDisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrgDisplayName

`func (o *InfrastructureProvider) SetOrgDisplayName(v string)`

SetOrgDisplayName sets OrgDisplayName field to given value.

### HasOrgDisplayName

`func (o *InfrastructureProvider) HasOrgDisplayName() bool`

HasOrgDisplayName returns a boolean if a field has been set.

### GetCreated

`func (o *InfrastructureProvider) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *InfrastructureProvider) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *InfrastructureProvider) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *InfrastructureProvider) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *InfrastructureProvider) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *InfrastructureProvider) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *InfrastructureProvider) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *InfrastructureProvider) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


