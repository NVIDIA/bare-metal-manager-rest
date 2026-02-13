# InstanceTypeSummary

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** |  | [optional] 
**InfrastructureProviderId** | Pointer to **string** |  | [optional] 
**SiteId** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**InstanceTypeStatus**](InstanceTypeStatus.md) |  | [optional] 

## Methods

### NewInstanceTypeSummary

`func NewInstanceTypeSummary() *InstanceTypeSummary`

NewInstanceTypeSummary instantiates a new InstanceTypeSummary object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceTypeSummaryWithDefaults

`func NewInstanceTypeSummaryWithDefaults() *InstanceTypeSummary`

NewInstanceTypeSummaryWithDefaults instantiates a new InstanceTypeSummary object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *InstanceTypeSummary) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *InstanceTypeSummary) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *InstanceTypeSummary) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *InstanceTypeSummary) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *InstanceTypeSummary) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *InstanceTypeSummary) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *InstanceTypeSummary) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *InstanceTypeSummary) HasName() bool`

HasName returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *InstanceTypeSummary) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *InstanceTypeSummary) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *InstanceTypeSummary) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *InstanceTypeSummary) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetSiteId

`func (o *InstanceTypeSummary) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *InstanceTypeSummary) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *InstanceTypeSummary) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *InstanceTypeSummary) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetStatus

`func (o *InstanceTypeSummary) GetStatus() InstanceTypeStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *InstanceTypeSummary) GetStatusOk() (*InstanceTypeStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *InstanceTypeSummary) SetStatus(v InstanceTypeStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *InstanceTypeSummary) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


