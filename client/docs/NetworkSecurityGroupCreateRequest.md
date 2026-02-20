# NetworkSecurityGroupCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**SiteId** | **string** |  | 
**StatefulEgress** | Pointer to **bool** | Egress rules with protocol and destination ports defined but without source ports defined should automatically be made stateful. | [optional] 
**Rules** | Pointer to [**[]NetworkSecurityGroupRule**](NetworkSecurityGroupRule.md) |  | [optional] 
**Labels** | Pointer to **map[string]string** |  | [optional] 

## Methods

### NewNetworkSecurityGroupCreateRequest

`func NewNetworkSecurityGroupCreateRequest(name string, siteId string, ) *NetworkSecurityGroupCreateRequest`

NewNetworkSecurityGroupCreateRequest instantiates a new NetworkSecurityGroupCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNetworkSecurityGroupCreateRequestWithDefaults

`func NewNetworkSecurityGroupCreateRequestWithDefaults() *NetworkSecurityGroupCreateRequest`

NewNetworkSecurityGroupCreateRequestWithDefaults instantiates a new NetworkSecurityGroupCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *NetworkSecurityGroupCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NetworkSecurityGroupCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NetworkSecurityGroupCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *NetworkSecurityGroupCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *NetworkSecurityGroupCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *NetworkSecurityGroupCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *NetworkSecurityGroupCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *NetworkSecurityGroupCreateRequest) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *NetworkSecurityGroupCreateRequest) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *NetworkSecurityGroupCreateRequest) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.


### GetStatefulEgress

`func (o *NetworkSecurityGroupCreateRequest) GetStatefulEgress() bool`

GetStatefulEgress returns the StatefulEgress field if non-nil, zero value otherwise.

### GetStatefulEgressOk

`func (o *NetworkSecurityGroupCreateRequest) GetStatefulEgressOk() (*bool, bool)`

GetStatefulEgressOk returns a tuple with the StatefulEgress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatefulEgress

`func (o *NetworkSecurityGroupCreateRequest) SetStatefulEgress(v bool)`

SetStatefulEgress sets StatefulEgress field to given value.

### HasStatefulEgress

`func (o *NetworkSecurityGroupCreateRequest) HasStatefulEgress() bool`

HasStatefulEgress returns a boolean if a field has been set.

### GetRules

`func (o *NetworkSecurityGroupCreateRequest) GetRules() []NetworkSecurityGroupRule`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *NetworkSecurityGroupCreateRequest) GetRulesOk() (*[]NetworkSecurityGroupRule, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *NetworkSecurityGroupCreateRequest) SetRules(v []NetworkSecurityGroupRule)`

SetRules sets Rules field to given value.

### HasRules

`func (o *NetworkSecurityGroupCreateRequest) HasRules() bool`

HasRules returns a boolean if a field has been set.

### GetLabels

`func (o *NetworkSecurityGroupCreateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *NetworkSecurityGroupCreateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *NetworkSecurityGroupCreateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *NetworkSecurityGroupCreateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


