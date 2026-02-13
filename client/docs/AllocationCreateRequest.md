# AllocationCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**TenantId** | **string** |  | 
**SiteId** | **string** |  | 
**AllocationConstraints** | Pointer to [**[]AllocationConstraintCreateRequest**](AllocationConstraintCreateRequest.md) |  | [optional] 

## Methods

### NewAllocationCreateRequest

`func NewAllocationCreateRequest(name string, tenantId string, siteId string, ) *AllocationCreateRequest`

NewAllocationCreateRequest instantiates a new AllocationCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAllocationCreateRequestWithDefaults

`func NewAllocationCreateRequestWithDefaults() *AllocationCreateRequest`

NewAllocationCreateRequestWithDefaults instantiates a new AllocationCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *AllocationCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AllocationCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AllocationCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *AllocationCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *AllocationCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *AllocationCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *AllocationCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetTenantId

`func (o *AllocationCreateRequest) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *AllocationCreateRequest) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *AllocationCreateRequest) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.


### GetSiteId

`func (o *AllocationCreateRequest) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *AllocationCreateRequest) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *AllocationCreateRequest) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.


### GetAllocationConstraints

`func (o *AllocationCreateRequest) GetAllocationConstraints() []AllocationConstraintCreateRequest`

GetAllocationConstraints returns the AllocationConstraints field if non-nil, zero value otherwise.

### GetAllocationConstraintsOk

`func (o *AllocationCreateRequest) GetAllocationConstraintsOk() (*[]AllocationConstraintCreateRequest, bool)`

GetAllocationConstraintsOk returns a tuple with the AllocationConstraints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocationConstraints

`func (o *AllocationCreateRequest) SetAllocationConstraints(v []AllocationConstraintCreateRequest)`

SetAllocationConstraints sets AllocationConstraints field to given value.

### HasAllocationConstraints

`func (o *AllocationCreateRequest) HasAllocationConstraints() bool`

HasAllocationConstraints returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


