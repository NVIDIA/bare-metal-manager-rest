# NetworkSecurityGroupUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**StatefulEgress** | Pointer to **bool** | Egress rules with protocol and destination ports defined but without source ports defined should automatically be made stateful. | [optional] 
**Rules** | Pointer to [**[]NetworkSecurityGroupRule**](NetworkSecurityGroupRule.md) | Update rules of the NetworkSecurityGroup. The rules will be entirely replaced by those sent in the request. Any rules not included in the request will be removed. To retain existing rules, first fetch them and include them.  | [optional] 
**Labels** | Pointer to **map[string]string** |  | [optional] 

## Methods

### NewNetworkSecurityGroupUpdateRequest

`func NewNetworkSecurityGroupUpdateRequest() *NetworkSecurityGroupUpdateRequest`

NewNetworkSecurityGroupUpdateRequest instantiates a new NetworkSecurityGroupUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNetworkSecurityGroupUpdateRequestWithDefaults

`func NewNetworkSecurityGroupUpdateRequestWithDefaults() *NetworkSecurityGroupUpdateRequest`

NewNetworkSecurityGroupUpdateRequestWithDefaults instantiates a new NetworkSecurityGroupUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *NetworkSecurityGroupUpdateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NetworkSecurityGroupUpdateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NetworkSecurityGroupUpdateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *NetworkSecurityGroupUpdateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### SetNameNil

`func (o *NetworkSecurityGroupUpdateRequest) SetNameNil(b bool)`

 SetNameNil sets the value for Name to be an explicit nil

### UnsetName
`func (o *NetworkSecurityGroupUpdateRequest) UnsetName()`

UnsetName ensures that no value is present for Name, not even an explicit nil
### GetDescription

`func (o *NetworkSecurityGroupUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *NetworkSecurityGroupUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *NetworkSecurityGroupUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *NetworkSecurityGroupUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *NetworkSecurityGroupUpdateRequest) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *NetworkSecurityGroupUpdateRequest) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetStatefulEgress

`func (o *NetworkSecurityGroupUpdateRequest) GetStatefulEgress() bool`

GetStatefulEgress returns the StatefulEgress field if non-nil, zero value otherwise.

### GetStatefulEgressOk

`func (o *NetworkSecurityGroupUpdateRequest) GetStatefulEgressOk() (*bool, bool)`

GetStatefulEgressOk returns a tuple with the StatefulEgress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatefulEgress

`func (o *NetworkSecurityGroupUpdateRequest) SetStatefulEgress(v bool)`

SetStatefulEgress sets StatefulEgress field to given value.

### HasStatefulEgress

`func (o *NetworkSecurityGroupUpdateRequest) HasStatefulEgress() bool`

HasStatefulEgress returns a boolean if a field has been set.

### GetRules

`func (o *NetworkSecurityGroupUpdateRequest) GetRules() []NetworkSecurityGroupRule`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *NetworkSecurityGroupUpdateRequest) GetRulesOk() (*[]NetworkSecurityGroupRule, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *NetworkSecurityGroupUpdateRequest) SetRules(v []NetworkSecurityGroupRule)`

SetRules sets Rules field to given value.

### HasRules

`func (o *NetworkSecurityGroupUpdateRequest) HasRules() bool`

HasRules returns a boolean if a field has been set.

### GetLabels

`func (o *NetworkSecurityGroupUpdateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *NetworkSecurityGroupUpdateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *NetworkSecurityGroupUpdateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *NetworkSecurityGroupUpdateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


