# NetworkSecurityGroupPropagationDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | The VPC or Instance ID | [optional] 
**DetailedStatus** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**NetworkSecurityGroupPropagationStatus**](NetworkSecurityGroupPropagationStatus.md) |  | [optional] 
**Details** | Pointer to **NullableString** |  | [optional] 
**UnpropagatedInstanceIds** | Pointer to **[]string** |  | [optional] 
**RelatedInstanceIds** | Pointer to **[]string** |  | [optional] 

## Methods

### NewNetworkSecurityGroupPropagationDetails

`func NewNetworkSecurityGroupPropagationDetails() *NetworkSecurityGroupPropagationDetails`

NewNetworkSecurityGroupPropagationDetails instantiates a new NetworkSecurityGroupPropagationDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNetworkSecurityGroupPropagationDetailsWithDefaults

`func NewNetworkSecurityGroupPropagationDetailsWithDefaults() *NetworkSecurityGroupPropagationDetails`

NewNetworkSecurityGroupPropagationDetailsWithDefaults instantiates a new NetworkSecurityGroupPropagationDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *NetworkSecurityGroupPropagationDetails) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *NetworkSecurityGroupPropagationDetails) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *NetworkSecurityGroupPropagationDetails) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *NetworkSecurityGroupPropagationDetails) HasId() bool`

HasId returns a boolean if a field has been set.

### GetDetailedStatus

`func (o *NetworkSecurityGroupPropagationDetails) GetDetailedStatus() string`

GetDetailedStatus returns the DetailedStatus field if non-nil, zero value otherwise.

### GetDetailedStatusOk

`func (o *NetworkSecurityGroupPropagationDetails) GetDetailedStatusOk() (*string, bool)`

GetDetailedStatusOk returns a tuple with the DetailedStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetailedStatus

`func (o *NetworkSecurityGroupPropagationDetails) SetDetailedStatus(v string)`

SetDetailedStatus sets DetailedStatus field to given value.

### HasDetailedStatus

`func (o *NetworkSecurityGroupPropagationDetails) HasDetailedStatus() bool`

HasDetailedStatus returns a boolean if a field has been set.

### GetStatus

`func (o *NetworkSecurityGroupPropagationDetails) GetStatus() NetworkSecurityGroupPropagationStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *NetworkSecurityGroupPropagationDetails) GetStatusOk() (*NetworkSecurityGroupPropagationStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *NetworkSecurityGroupPropagationDetails) SetStatus(v NetworkSecurityGroupPropagationStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *NetworkSecurityGroupPropagationDetails) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetDetails

`func (o *NetworkSecurityGroupPropagationDetails) GetDetails() string`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *NetworkSecurityGroupPropagationDetails) GetDetailsOk() (*string, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetails

`func (o *NetworkSecurityGroupPropagationDetails) SetDetails(v string)`

SetDetails sets Details field to given value.

### HasDetails

`func (o *NetworkSecurityGroupPropagationDetails) HasDetails() bool`

HasDetails returns a boolean if a field has been set.

### SetDetailsNil

`func (o *NetworkSecurityGroupPropagationDetails) SetDetailsNil(b bool)`

 SetDetailsNil sets the value for Details to be an explicit nil

### UnsetDetails
`func (o *NetworkSecurityGroupPropagationDetails) UnsetDetails()`

UnsetDetails ensures that no value is present for Details, not even an explicit nil
### GetUnpropagatedInstanceIds

`func (o *NetworkSecurityGroupPropagationDetails) GetUnpropagatedInstanceIds() []string`

GetUnpropagatedInstanceIds returns the UnpropagatedInstanceIds field if non-nil, zero value otherwise.

### GetUnpropagatedInstanceIdsOk

`func (o *NetworkSecurityGroupPropagationDetails) GetUnpropagatedInstanceIdsOk() (*[]string, bool)`

GetUnpropagatedInstanceIdsOk returns a tuple with the UnpropagatedInstanceIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnpropagatedInstanceIds

`func (o *NetworkSecurityGroupPropagationDetails) SetUnpropagatedInstanceIds(v []string)`

SetUnpropagatedInstanceIds sets UnpropagatedInstanceIds field to given value.

### HasUnpropagatedInstanceIds

`func (o *NetworkSecurityGroupPropagationDetails) HasUnpropagatedInstanceIds() bool`

HasUnpropagatedInstanceIds returns a boolean if a field has been set.

### GetRelatedInstanceIds

`func (o *NetworkSecurityGroupPropagationDetails) GetRelatedInstanceIds() []string`

GetRelatedInstanceIds returns the RelatedInstanceIds field if non-nil, zero value otherwise.

### GetRelatedInstanceIdsOk

`func (o *NetworkSecurityGroupPropagationDetails) GetRelatedInstanceIdsOk() (*[]string, bool)`

GetRelatedInstanceIdsOk returns a tuple with the RelatedInstanceIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelatedInstanceIds

`func (o *NetworkSecurityGroupPropagationDetails) SetRelatedInstanceIds(v []string)`

SetRelatedInstanceIds sets RelatedInstanceIds field to given value.

### HasRelatedInstanceIds

`func (o *NetworkSecurityGroupPropagationDetails) HasRelatedInstanceIds() bool`

HasRelatedInstanceIds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


