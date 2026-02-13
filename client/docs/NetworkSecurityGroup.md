# NetworkSecurityGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**SiteId** | Pointer to **string** |  | [optional] 
**TenantId** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**NetworkSecurityGroupStatus**](NetworkSecurityGroupStatus.md) |  | [optional] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**StatefulEgress** | Pointer to **bool** |  | [optional] 
**Rules** | Pointer to [**[]NetworkSecurityGroupRule**](NetworkSecurityGroupRule.md) |  | [optional] 
**Labels** | Pointer to **map[string]string** |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] 
**Updated** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewNetworkSecurityGroup

`func NewNetworkSecurityGroup() *NetworkSecurityGroup`

NewNetworkSecurityGroup instantiates a new NetworkSecurityGroup object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNetworkSecurityGroupWithDefaults

`func NewNetworkSecurityGroupWithDefaults() *NetworkSecurityGroup`

NewNetworkSecurityGroupWithDefaults instantiates a new NetworkSecurityGroup object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *NetworkSecurityGroup) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *NetworkSecurityGroup) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *NetworkSecurityGroup) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *NetworkSecurityGroup) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *NetworkSecurityGroup) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NetworkSecurityGroup) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NetworkSecurityGroup) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *NetworkSecurityGroup) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *NetworkSecurityGroup) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *NetworkSecurityGroup) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *NetworkSecurityGroup) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *NetworkSecurityGroup) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *NetworkSecurityGroup) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *NetworkSecurityGroup) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *NetworkSecurityGroup) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *NetworkSecurityGroup) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetTenantId

`func (o *NetworkSecurityGroup) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *NetworkSecurityGroup) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *NetworkSecurityGroup) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *NetworkSecurityGroup) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetStatus

`func (o *NetworkSecurityGroup) GetStatus() NetworkSecurityGroupStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *NetworkSecurityGroup) GetStatusOk() (*NetworkSecurityGroupStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *NetworkSecurityGroup) SetStatus(v NetworkSecurityGroupStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *NetworkSecurityGroup) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *NetworkSecurityGroup) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *NetworkSecurityGroup) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *NetworkSecurityGroup) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *NetworkSecurityGroup) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetStatefulEgress

`func (o *NetworkSecurityGroup) GetStatefulEgress() bool`

GetStatefulEgress returns the StatefulEgress field if non-nil, zero value otherwise.

### GetStatefulEgressOk

`func (o *NetworkSecurityGroup) GetStatefulEgressOk() (*bool, bool)`

GetStatefulEgressOk returns a tuple with the StatefulEgress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatefulEgress

`func (o *NetworkSecurityGroup) SetStatefulEgress(v bool)`

SetStatefulEgress sets StatefulEgress field to given value.

### HasStatefulEgress

`func (o *NetworkSecurityGroup) HasStatefulEgress() bool`

HasStatefulEgress returns a boolean if a field has been set.

### GetRules

`func (o *NetworkSecurityGroup) GetRules() []NetworkSecurityGroupRule`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *NetworkSecurityGroup) GetRulesOk() (*[]NetworkSecurityGroupRule, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *NetworkSecurityGroup) SetRules(v []NetworkSecurityGroupRule)`

SetRules sets Rules field to given value.

### HasRules

`func (o *NetworkSecurityGroup) HasRules() bool`

HasRules returns a boolean if a field has been set.

### GetLabels

`func (o *NetworkSecurityGroup) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *NetworkSecurityGroup) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *NetworkSecurityGroup) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *NetworkSecurityGroup) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetCreated

`func (o *NetworkSecurityGroup) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *NetworkSecurityGroup) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *NetworkSecurityGroup) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *NetworkSecurityGroup) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *NetworkSecurityGroup) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *NetworkSecurityGroup) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *NetworkSecurityGroup) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *NetworkSecurityGroup) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


