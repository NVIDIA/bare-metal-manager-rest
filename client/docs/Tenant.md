# Tenant

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Carbide issued ID of the Tenant | [optional] [readonly] 
**Org** | Pointer to **string** | Name/external ID of Tenant&#39;s organization | [optional] 
**OrgDisplayName** | Pointer to **string** | Display name of Tenant&#39;s organization | [optional] 
**Created** | Pointer to **time.Time** | Date/time the Tenant was created | [optional] [readonly] 
**Updated** | Pointer to **time.Time** | Date/time when Tenant was last updated | [optional] 
**Capabilities** | Pointer to [**TenantCapabilities**](TenantCapabilities.md) | Features that are enabled/disabled for Tenant | [optional] 

## Methods

### NewTenant

`func NewTenant() *Tenant`

NewTenant instantiates a new Tenant object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTenantWithDefaults

`func NewTenantWithDefaults() *Tenant`

NewTenantWithDefaults instantiates a new Tenant object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Tenant) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Tenant) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Tenant) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Tenant) HasId() bool`

HasId returns a boolean if a field has been set.

### GetOrg

`func (o *Tenant) GetOrg() string`

GetOrg returns the Org field if non-nil, zero value otherwise.

### GetOrgOk

`func (o *Tenant) GetOrgOk() (*string, bool)`

GetOrgOk returns a tuple with the Org field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrg

`func (o *Tenant) SetOrg(v string)`

SetOrg sets Org field to given value.

### HasOrg

`func (o *Tenant) HasOrg() bool`

HasOrg returns a boolean if a field has been set.

### GetOrgDisplayName

`func (o *Tenant) GetOrgDisplayName() string`

GetOrgDisplayName returns the OrgDisplayName field if non-nil, zero value otherwise.

### GetOrgDisplayNameOk

`func (o *Tenant) GetOrgDisplayNameOk() (*string, bool)`

GetOrgDisplayNameOk returns a tuple with the OrgDisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrgDisplayName

`func (o *Tenant) SetOrgDisplayName(v string)`

SetOrgDisplayName sets OrgDisplayName field to given value.

### HasOrgDisplayName

`func (o *Tenant) HasOrgDisplayName() bool`

HasOrgDisplayName returns a boolean if a field has been set.

### GetCreated

`func (o *Tenant) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *Tenant) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *Tenant) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *Tenant) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *Tenant) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *Tenant) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *Tenant) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *Tenant) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.

### GetCapabilities

`func (o *Tenant) GetCapabilities() TenantCapabilities`

GetCapabilities returns the Capabilities field if non-nil, zero value otherwise.

### GetCapabilitiesOk

`func (o *Tenant) GetCapabilitiesOk() (*TenantCapabilities, bool)`

GetCapabilitiesOk returns a tuple with the Capabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapabilities

`func (o *Tenant) SetCapabilities(v TenantCapabilities)`

SetCapabilities sets Capabilities field to given value.

### HasCapabilities

`func (o *Tenant) HasCapabilities() bool`

HasCapabilities returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


