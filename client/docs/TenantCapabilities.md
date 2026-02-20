# TenantCapabilities

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TargetedInstanceCreation** | Pointer to **bool** | Indicates whether Tenant can create Instances by specifying Machine ID | [optional] [readonly] 

## Methods

### NewTenantCapabilities

`func NewTenantCapabilities() *TenantCapabilities`

NewTenantCapabilities instantiates a new TenantCapabilities object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTenantCapabilitiesWithDefaults

`func NewTenantCapabilitiesWithDefaults() *TenantCapabilities`

NewTenantCapabilitiesWithDefaults instantiates a new TenantCapabilities object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTargetedInstanceCreation

`func (o *TenantCapabilities) GetTargetedInstanceCreation() bool`

GetTargetedInstanceCreation returns the TargetedInstanceCreation field if non-nil, zero value otherwise.

### GetTargetedInstanceCreationOk

`func (o *TenantCapabilities) GetTargetedInstanceCreationOk() (*bool, bool)`

GetTargetedInstanceCreationOk returns a tuple with the TargetedInstanceCreation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetedInstanceCreation

`func (o *TenantCapabilities) SetTargetedInstanceCreation(v bool)`

SetTargetedInstanceCreation sets TargetedInstanceCreation field to given value.

### HasTargetedInstanceCreation

`func (o *TenantCapabilities) HasTargetedInstanceCreation() bool`

HasTargetedInstanceCreation returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


