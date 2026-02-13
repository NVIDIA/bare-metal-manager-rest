# SiteSummary

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** | Name of the Site | [optional] 
**InfrastructureProviderId** | Pointer to **string** |  | [optional] 
**IsSerialConsoleEnabled** | Pointer to **bool** | Indicates if Serial Console is enabled for the Site by the Provider | [optional] 
**IsOnline** | Pointer to **bool** | Indicates if the Site is currently reachable from Cloud | [optional] 
**Capabilities** | Pointer to [**SiteCapabilities**](SiteCapabilities.md) |  | [optional] 
**Status** | Pointer to [**SiteStatus**](SiteStatus.md) |  | [optional] [readonly] 

## Methods

### NewSiteSummary

`func NewSiteSummary() *SiteSummary`

NewSiteSummary instantiates a new SiteSummary object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSiteSummaryWithDefaults

`func NewSiteSummaryWithDefaults() *SiteSummary`

NewSiteSummaryWithDefaults instantiates a new SiteSummary object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *SiteSummary) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SiteSummary) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SiteSummary) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *SiteSummary) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *SiteSummary) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SiteSummary) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SiteSummary) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SiteSummary) HasName() bool`

HasName returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *SiteSummary) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *SiteSummary) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *SiteSummary) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *SiteSummary) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetIsSerialConsoleEnabled

`func (o *SiteSummary) GetIsSerialConsoleEnabled() bool`

GetIsSerialConsoleEnabled returns the IsSerialConsoleEnabled field if non-nil, zero value otherwise.

### GetIsSerialConsoleEnabledOk

`func (o *SiteSummary) GetIsSerialConsoleEnabledOk() (*bool, bool)`

GetIsSerialConsoleEnabledOk returns a tuple with the IsSerialConsoleEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsSerialConsoleEnabled

`func (o *SiteSummary) SetIsSerialConsoleEnabled(v bool)`

SetIsSerialConsoleEnabled sets IsSerialConsoleEnabled field to given value.

### HasIsSerialConsoleEnabled

`func (o *SiteSummary) HasIsSerialConsoleEnabled() bool`

HasIsSerialConsoleEnabled returns a boolean if a field has been set.

### GetIsOnline

`func (o *SiteSummary) GetIsOnline() bool`

GetIsOnline returns the IsOnline field if non-nil, zero value otherwise.

### GetIsOnlineOk

`func (o *SiteSummary) GetIsOnlineOk() (*bool, bool)`

GetIsOnlineOk returns a tuple with the IsOnline field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsOnline

`func (o *SiteSummary) SetIsOnline(v bool)`

SetIsOnline sets IsOnline field to given value.

### HasIsOnline

`func (o *SiteSummary) HasIsOnline() bool`

HasIsOnline returns a boolean if a field has been set.

### GetCapabilities

`func (o *SiteSummary) GetCapabilities() SiteCapabilities`

GetCapabilities returns the Capabilities field if non-nil, zero value otherwise.

### GetCapabilitiesOk

`func (o *SiteSummary) GetCapabilitiesOk() (*SiteCapabilities, bool)`

GetCapabilitiesOk returns a tuple with the Capabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapabilities

`func (o *SiteSummary) SetCapabilities(v SiteCapabilities)`

SetCapabilities sets Capabilities field to given value.

### HasCapabilities

`func (o *SiteSummary) HasCapabilities() bool`

HasCapabilities returns a boolean if a field has been set.

### GetStatus

`func (o *SiteSummary) GetStatus() SiteStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *SiteSummary) GetStatusOk() (*SiteStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *SiteSummary) SetStatus(v SiteStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *SiteSummary) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


