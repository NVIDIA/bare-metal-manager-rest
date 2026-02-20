# AllocationConstraintCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ResourceType** | Pointer to **string** |  | [optional] 
**ResourceTypeId** | **string** |  | 
**ConstraintType** | **string** |  | 
**ConstraintValue** | **int32** |  | 

## Methods

### NewAllocationConstraintCreateRequest

`func NewAllocationConstraintCreateRequest(resourceTypeId string, constraintType string, constraintValue int32, ) *AllocationConstraintCreateRequest`

NewAllocationConstraintCreateRequest instantiates a new AllocationConstraintCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAllocationConstraintCreateRequestWithDefaults

`func NewAllocationConstraintCreateRequestWithDefaults() *AllocationConstraintCreateRequest`

NewAllocationConstraintCreateRequestWithDefaults instantiates a new AllocationConstraintCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetResourceType

`func (o *AllocationConstraintCreateRequest) GetResourceType() string`

GetResourceType returns the ResourceType field if non-nil, zero value otherwise.

### GetResourceTypeOk

`func (o *AllocationConstraintCreateRequest) GetResourceTypeOk() (*string, bool)`

GetResourceTypeOk returns a tuple with the ResourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceType

`func (o *AllocationConstraintCreateRequest) SetResourceType(v string)`

SetResourceType sets ResourceType field to given value.

### HasResourceType

`func (o *AllocationConstraintCreateRequest) HasResourceType() bool`

HasResourceType returns a boolean if a field has been set.

### GetResourceTypeId

`func (o *AllocationConstraintCreateRequest) GetResourceTypeId() string`

GetResourceTypeId returns the ResourceTypeId field if non-nil, zero value otherwise.

### GetResourceTypeIdOk

`func (o *AllocationConstraintCreateRequest) GetResourceTypeIdOk() (*string, bool)`

GetResourceTypeIdOk returns a tuple with the ResourceTypeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceTypeId

`func (o *AllocationConstraintCreateRequest) SetResourceTypeId(v string)`

SetResourceTypeId sets ResourceTypeId field to given value.


### GetConstraintType

`func (o *AllocationConstraintCreateRequest) GetConstraintType() string`

GetConstraintType returns the ConstraintType field if non-nil, zero value otherwise.

### GetConstraintTypeOk

`func (o *AllocationConstraintCreateRequest) GetConstraintTypeOk() (*string, bool)`

GetConstraintTypeOk returns a tuple with the ConstraintType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConstraintType

`func (o *AllocationConstraintCreateRequest) SetConstraintType(v string)`

SetConstraintType sets ConstraintType field to given value.


### GetConstraintValue

`func (o *AllocationConstraintCreateRequest) GetConstraintValue() int32`

GetConstraintValue returns the ConstraintValue field if non-nil, zero value otherwise.

### GetConstraintValueOk

`func (o *AllocationConstraintCreateRequest) GetConstraintValueOk() (*int32, bool)`

GetConstraintValueOk returns a tuple with the ConstraintValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConstraintValue

`func (o *AllocationConstraintCreateRequest) SetConstraintValue(v int32)`

SetConstraintValue sets ConstraintValue field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


