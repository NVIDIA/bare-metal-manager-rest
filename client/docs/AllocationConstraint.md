# AllocationConstraint

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**AllocationId** | Pointer to **string** |  | [optional] 
**ResourceType** | Pointer to **string** |  | [optional] 
**ResourceTypeId** | Pointer to **string** |  | [optional] 
**ConstraintType** | Pointer to **string** |  | [optional] 
**ConstraintValue** | Pointer to **int32** |  | [optional] 
**DerivedResourceId** | Pointer to **NullableString** | Populated with resulting IP Block ID when resource type is IP Block | [optional] 
**InstanceType** | Pointer to [**InstanceTypeSummary**](InstanceTypeSummary.md) |  | [optional] 
**IpBlock** | Pointer to [**IpBlockSummary**](IpBlockSummary.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewAllocationConstraint

`func NewAllocationConstraint() *AllocationConstraint`

NewAllocationConstraint instantiates a new AllocationConstraint object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAllocationConstraintWithDefaults

`func NewAllocationConstraintWithDefaults() *AllocationConstraint`

NewAllocationConstraintWithDefaults instantiates a new AllocationConstraint object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AllocationConstraint) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AllocationConstraint) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AllocationConstraint) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AllocationConstraint) HasId() bool`

HasId returns a boolean if a field has been set.

### GetAllocationId

`func (o *AllocationConstraint) GetAllocationId() string`

GetAllocationId returns the AllocationId field if non-nil, zero value otherwise.

### GetAllocationIdOk

`func (o *AllocationConstraint) GetAllocationIdOk() (*string, bool)`

GetAllocationIdOk returns a tuple with the AllocationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocationId

`func (o *AllocationConstraint) SetAllocationId(v string)`

SetAllocationId sets AllocationId field to given value.

### HasAllocationId

`func (o *AllocationConstraint) HasAllocationId() bool`

HasAllocationId returns a boolean if a field has been set.

### GetResourceType

`func (o *AllocationConstraint) GetResourceType() string`

GetResourceType returns the ResourceType field if non-nil, zero value otherwise.

### GetResourceTypeOk

`func (o *AllocationConstraint) GetResourceTypeOk() (*string, bool)`

GetResourceTypeOk returns a tuple with the ResourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceType

`func (o *AllocationConstraint) SetResourceType(v string)`

SetResourceType sets ResourceType field to given value.

### HasResourceType

`func (o *AllocationConstraint) HasResourceType() bool`

HasResourceType returns a boolean if a field has been set.

### GetResourceTypeId

`func (o *AllocationConstraint) GetResourceTypeId() string`

GetResourceTypeId returns the ResourceTypeId field if non-nil, zero value otherwise.

### GetResourceTypeIdOk

`func (o *AllocationConstraint) GetResourceTypeIdOk() (*string, bool)`

GetResourceTypeIdOk returns a tuple with the ResourceTypeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceTypeId

`func (o *AllocationConstraint) SetResourceTypeId(v string)`

SetResourceTypeId sets ResourceTypeId field to given value.

### HasResourceTypeId

`func (o *AllocationConstraint) HasResourceTypeId() bool`

HasResourceTypeId returns a boolean if a field has been set.

### GetConstraintType

`func (o *AllocationConstraint) GetConstraintType() string`

GetConstraintType returns the ConstraintType field if non-nil, zero value otherwise.

### GetConstraintTypeOk

`func (o *AllocationConstraint) GetConstraintTypeOk() (*string, bool)`

GetConstraintTypeOk returns a tuple with the ConstraintType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConstraintType

`func (o *AllocationConstraint) SetConstraintType(v string)`

SetConstraintType sets ConstraintType field to given value.

### HasConstraintType

`func (o *AllocationConstraint) HasConstraintType() bool`

HasConstraintType returns a boolean if a field has been set.

### GetConstraintValue

`func (o *AllocationConstraint) GetConstraintValue() int32`

GetConstraintValue returns the ConstraintValue field if non-nil, zero value otherwise.

### GetConstraintValueOk

`func (o *AllocationConstraint) GetConstraintValueOk() (*int32, bool)`

GetConstraintValueOk returns a tuple with the ConstraintValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConstraintValue

`func (o *AllocationConstraint) SetConstraintValue(v int32)`

SetConstraintValue sets ConstraintValue field to given value.

### HasConstraintValue

`func (o *AllocationConstraint) HasConstraintValue() bool`

HasConstraintValue returns a boolean if a field has been set.

### GetDerivedResourceId

`func (o *AllocationConstraint) GetDerivedResourceId() string`

GetDerivedResourceId returns the DerivedResourceId field if non-nil, zero value otherwise.

### GetDerivedResourceIdOk

`func (o *AllocationConstraint) GetDerivedResourceIdOk() (*string, bool)`

GetDerivedResourceIdOk returns a tuple with the DerivedResourceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDerivedResourceId

`func (o *AllocationConstraint) SetDerivedResourceId(v string)`

SetDerivedResourceId sets DerivedResourceId field to given value.

### HasDerivedResourceId

`func (o *AllocationConstraint) HasDerivedResourceId() bool`

HasDerivedResourceId returns a boolean if a field has been set.

### SetDerivedResourceIdNil

`func (o *AllocationConstraint) SetDerivedResourceIdNil(b bool)`

 SetDerivedResourceIdNil sets the value for DerivedResourceId to be an explicit nil

### UnsetDerivedResourceId
`func (o *AllocationConstraint) UnsetDerivedResourceId()`

UnsetDerivedResourceId ensures that no value is present for DerivedResourceId, not even an explicit nil
### GetInstanceType

`func (o *AllocationConstraint) GetInstanceType() InstanceTypeSummary`

GetInstanceType returns the InstanceType field if non-nil, zero value otherwise.

### GetInstanceTypeOk

`func (o *AllocationConstraint) GetInstanceTypeOk() (*InstanceTypeSummary, bool)`

GetInstanceTypeOk returns a tuple with the InstanceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceType

`func (o *AllocationConstraint) SetInstanceType(v InstanceTypeSummary)`

SetInstanceType sets InstanceType field to given value.

### HasInstanceType

`func (o *AllocationConstraint) HasInstanceType() bool`

HasInstanceType returns a boolean if a field has been set.

### GetIpBlock

`func (o *AllocationConstraint) GetIpBlock() IpBlockSummary`

GetIpBlock returns the IpBlock field if non-nil, zero value otherwise.

### GetIpBlockOk

`func (o *AllocationConstraint) GetIpBlockOk() (*IpBlockSummary, bool)`

GetIpBlockOk returns a tuple with the IpBlock field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpBlock

`func (o *AllocationConstraint) SetIpBlock(v IpBlockSummary)`

SetIpBlock sets IpBlock field to given value.

### HasIpBlock

`func (o *AllocationConstraint) HasIpBlock() bool`

HasIpBlock returns a boolean if a field has been set.

### GetCreated

`func (o *AllocationConstraint) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *AllocationConstraint) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *AllocationConstraint) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *AllocationConstraint) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *AllocationConstraint) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *AllocationConstraint) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *AllocationConstraint) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *AllocationConstraint) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


