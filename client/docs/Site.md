# Site

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** | Name of the Site | [optional] 
**Description** | Pointer to **string** | Optional description for the Site | [optional] 
**Org** | Pointer to **string** |  | [optional] 
**InfrastructureProviderId** | Pointer to **string** |  | [optional] 
**SiteControllerVersion** | Pointer to **string** | Version of the Site Controller software | [optional] 
**SiteAgentVersion** | Pointer to **string** | Version of the Site Agent software | [optional] 
**RegistrationToken** | Pointer to **string** | Token that can be used to register a Site. Value only exposed to Provider | [optional] 
**RegistrationTokenExpiration** | Pointer to **time.Time** | Date/time when registration token  expires. Value only exposed to Provider | [optional] [readonly] 
**SerialConsoleHostname** | Pointer to **string** |  | [optional] 
**IsSerialConsoleEnabled** | Pointer to **bool** | Indicates if Serial Console is enabled for the Site by the Provider | [optional] 
**SerialConsoleIdleTimeout** | Pointer to **NullableInt32** | Maximum idle time in seconds before Serial Console is disconnected | [optional] 
**SerialConsoleMaxSessionLength** | Pointer to **NullableInt32** | Maximum length of Serial Console session in seconds | [optional] 
**IsSerialConsoleSSHKeysEnabled** | Pointer to **bool** | Only visible to Tenant retrieving the Site. Indicates if Serial Console access using SSH Keys is enabled by Tenant | [optional] 
**IsOnline** | Pointer to **bool** | Indicates if the Site is currently reachable from Cloud | [optional] 
**Status** | Pointer to [**SiteStatus**](SiteStatus.md) |  | [optional] [readonly] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] 
**Updated** | Pointer to **time.Time** |  | [optional] 
**Location** | Pointer to [**SiteLocation**](SiteLocation.md) |  | [optional] 
**Contact** | Pointer to [**SiteContact**](SiteContact.md) |  | [optional] 
**Capabilities** | Pointer to [**SiteCapabilities**](SiteCapabilities.md) |  | [optional] 

## Methods

### NewSite

`func NewSite() *Site`

NewSite instantiates a new Site object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSiteWithDefaults

`func NewSiteWithDefaults() *Site`

NewSiteWithDefaults instantiates a new Site object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Site) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Site) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Site) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Site) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *Site) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Site) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Site) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Site) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *Site) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Site) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Site) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Site) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetOrg

`func (o *Site) GetOrg() string`

GetOrg returns the Org field if non-nil, zero value otherwise.

### GetOrgOk

`func (o *Site) GetOrgOk() (*string, bool)`

GetOrgOk returns a tuple with the Org field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrg

`func (o *Site) SetOrg(v string)`

SetOrg sets Org field to given value.

### HasOrg

`func (o *Site) HasOrg() bool`

HasOrg returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *Site) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *Site) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *Site) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *Site) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetSiteControllerVersion

`func (o *Site) GetSiteControllerVersion() string`

GetSiteControllerVersion returns the SiteControllerVersion field if non-nil, zero value otherwise.

### GetSiteControllerVersionOk

`func (o *Site) GetSiteControllerVersionOk() (*string, bool)`

GetSiteControllerVersionOk returns a tuple with the SiteControllerVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteControllerVersion

`func (o *Site) SetSiteControllerVersion(v string)`

SetSiteControllerVersion sets SiteControllerVersion field to given value.

### HasSiteControllerVersion

`func (o *Site) HasSiteControllerVersion() bool`

HasSiteControllerVersion returns a boolean if a field has been set.

### GetSiteAgentVersion

`func (o *Site) GetSiteAgentVersion() string`

GetSiteAgentVersion returns the SiteAgentVersion field if non-nil, zero value otherwise.

### GetSiteAgentVersionOk

`func (o *Site) GetSiteAgentVersionOk() (*string, bool)`

GetSiteAgentVersionOk returns a tuple with the SiteAgentVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteAgentVersion

`func (o *Site) SetSiteAgentVersion(v string)`

SetSiteAgentVersion sets SiteAgentVersion field to given value.

### HasSiteAgentVersion

`func (o *Site) HasSiteAgentVersion() bool`

HasSiteAgentVersion returns a boolean if a field has been set.

### GetRegistrationToken

`func (o *Site) GetRegistrationToken() string`

GetRegistrationToken returns the RegistrationToken field if non-nil, zero value otherwise.

### GetRegistrationTokenOk

`func (o *Site) GetRegistrationTokenOk() (*string, bool)`

GetRegistrationTokenOk returns a tuple with the RegistrationToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrationToken

`func (o *Site) SetRegistrationToken(v string)`

SetRegistrationToken sets RegistrationToken field to given value.

### HasRegistrationToken

`func (o *Site) HasRegistrationToken() bool`

HasRegistrationToken returns a boolean if a field has been set.

### GetRegistrationTokenExpiration

`func (o *Site) GetRegistrationTokenExpiration() time.Time`

GetRegistrationTokenExpiration returns the RegistrationTokenExpiration field if non-nil, zero value otherwise.

### GetRegistrationTokenExpirationOk

`func (o *Site) GetRegistrationTokenExpirationOk() (*time.Time, bool)`

GetRegistrationTokenExpirationOk returns a tuple with the RegistrationTokenExpiration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrationTokenExpiration

`func (o *Site) SetRegistrationTokenExpiration(v time.Time)`

SetRegistrationTokenExpiration sets RegistrationTokenExpiration field to given value.

### HasRegistrationTokenExpiration

`func (o *Site) HasRegistrationTokenExpiration() bool`

HasRegistrationTokenExpiration returns a boolean if a field has been set.

### GetSerialConsoleHostname

`func (o *Site) GetSerialConsoleHostname() string`

GetSerialConsoleHostname returns the SerialConsoleHostname field if non-nil, zero value otherwise.

### GetSerialConsoleHostnameOk

`func (o *Site) GetSerialConsoleHostnameOk() (*string, bool)`

GetSerialConsoleHostnameOk returns a tuple with the SerialConsoleHostname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialConsoleHostname

`func (o *Site) SetSerialConsoleHostname(v string)`

SetSerialConsoleHostname sets SerialConsoleHostname field to given value.

### HasSerialConsoleHostname

`func (o *Site) HasSerialConsoleHostname() bool`

HasSerialConsoleHostname returns a boolean if a field has been set.

### GetIsSerialConsoleEnabled

`func (o *Site) GetIsSerialConsoleEnabled() bool`

GetIsSerialConsoleEnabled returns the IsSerialConsoleEnabled field if non-nil, zero value otherwise.

### GetIsSerialConsoleEnabledOk

`func (o *Site) GetIsSerialConsoleEnabledOk() (*bool, bool)`

GetIsSerialConsoleEnabledOk returns a tuple with the IsSerialConsoleEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsSerialConsoleEnabled

`func (o *Site) SetIsSerialConsoleEnabled(v bool)`

SetIsSerialConsoleEnabled sets IsSerialConsoleEnabled field to given value.

### HasIsSerialConsoleEnabled

`func (o *Site) HasIsSerialConsoleEnabled() bool`

HasIsSerialConsoleEnabled returns a boolean if a field has been set.

### GetSerialConsoleIdleTimeout

`func (o *Site) GetSerialConsoleIdleTimeout() int32`

GetSerialConsoleIdleTimeout returns the SerialConsoleIdleTimeout field if non-nil, zero value otherwise.

### GetSerialConsoleIdleTimeoutOk

`func (o *Site) GetSerialConsoleIdleTimeoutOk() (*int32, bool)`

GetSerialConsoleIdleTimeoutOk returns a tuple with the SerialConsoleIdleTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialConsoleIdleTimeout

`func (o *Site) SetSerialConsoleIdleTimeout(v int32)`

SetSerialConsoleIdleTimeout sets SerialConsoleIdleTimeout field to given value.

### HasSerialConsoleIdleTimeout

`func (o *Site) HasSerialConsoleIdleTimeout() bool`

HasSerialConsoleIdleTimeout returns a boolean if a field has been set.

### SetSerialConsoleIdleTimeoutNil

`func (o *Site) SetSerialConsoleIdleTimeoutNil(b bool)`

 SetSerialConsoleIdleTimeoutNil sets the value for SerialConsoleIdleTimeout to be an explicit nil

### UnsetSerialConsoleIdleTimeout
`func (o *Site) UnsetSerialConsoleIdleTimeout()`

UnsetSerialConsoleIdleTimeout ensures that no value is present for SerialConsoleIdleTimeout, not even an explicit nil
### GetSerialConsoleMaxSessionLength

`func (o *Site) GetSerialConsoleMaxSessionLength() int32`

GetSerialConsoleMaxSessionLength returns the SerialConsoleMaxSessionLength field if non-nil, zero value otherwise.

### GetSerialConsoleMaxSessionLengthOk

`func (o *Site) GetSerialConsoleMaxSessionLengthOk() (*int32, bool)`

GetSerialConsoleMaxSessionLengthOk returns a tuple with the SerialConsoleMaxSessionLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialConsoleMaxSessionLength

`func (o *Site) SetSerialConsoleMaxSessionLength(v int32)`

SetSerialConsoleMaxSessionLength sets SerialConsoleMaxSessionLength field to given value.

### HasSerialConsoleMaxSessionLength

`func (o *Site) HasSerialConsoleMaxSessionLength() bool`

HasSerialConsoleMaxSessionLength returns a boolean if a field has been set.

### SetSerialConsoleMaxSessionLengthNil

`func (o *Site) SetSerialConsoleMaxSessionLengthNil(b bool)`

 SetSerialConsoleMaxSessionLengthNil sets the value for SerialConsoleMaxSessionLength to be an explicit nil

### UnsetSerialConsoleMaxSessionLength
`func (o *Site) UnsetSerialConsoleMaxSessionLength()`

UnsetSerialConsoleMaxSessionLength ensures that no value is present for SerialConsoleMaxSessionLength, not even an explicit nil
### GetIsSerialConsoleSSHKeysEnabled

`func (o *Site) GetIsSerialConsoleSSHKeysEnabled() bool`

GetIsSerialConsoleSSHKeysEnabled returns the IsSerialConsoleSSHKeysEnabled field if non-nil, zero value otherwise.

### GetIsSerialConsoleSSHKeysEnabledOk

`func (o *Site) GetIsSerialConsoleSSHKeysEnabledOk() (*bool, bool)`

GetIsSerialConsoleSSHKeysEnabledOk returns a tuple with the IsSerialConsoleSSHKeysEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsSerialConsoleSSHKeysEnabled

`func (o *Site) SetIsSerialConsoleSSHKeysEnabled(v bool)`

SetIsSerialConsoleSSHKeysEnabled sets IsSerialConsoleSSHKeysEnabled field to given value.

### HasIsSerialConsoleSSHKeysEnabled

`func (o *Site) HasIsSerialConsoleSSHKeysEnabled() bool`

HasIsSerialConsoleSSHKeysEnabled returns a boolean if a field has been set.

### GetIsOnline

`func (o *Site) GetIsOnline() bool`

GetIsOnline returns the IsOnline field if non-nil, zero value otherwise.

### GetIsOnlineOk

`func (o *Site) GetIsOnlineOk() (*bool, bool)`

GetIsOnlineOk returns a tuple with the IsOnline field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsOnline

`func (o *Site) SetIsOnline(v bool)`

SetIsOnline sets IsOnline field to given value.

### HasIsOnline

`func (o *Site) HasIsOnline() bool`

HasIsOnline returns a boolean if a field has been set.

### GetStatus

`func (o *Site) GetStatus() SiteStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Site) GetStatusOk() (*SiteStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Site) SetStatus(v SiteStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Site) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *Site) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *Site) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *Site) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *Site) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *Site) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *Site) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *Site) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *Site) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *Site) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *Site) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *Site) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *Site) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.

### GetLocation

`func (o *Site) GetLocation() SiteLocation`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *Site) GetLocationOk() (*SiteLocation, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *Site) SetLocation(v SiteLocation)`

SetLocation sets Location field to given value.

### HasLocation

`func (o *Site) HasLocation() bool`

HasLocation returns a boolean if a field has been set.

### GetContact

`func (o *Site) GetContact() SiteContact`

GetContact returns the Contact field if non-nil, zero value otherwise.

### GetContactOk

`func (o *Site) GetContactOk() (*SiteContact, bool)`

GetContactOk returns a tuple with the Contact field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContact

`func (o *Site) SetContact(v SiteContact)`

SetContact sets Contact field to given value.

### HasContact

`func (o *Site) HasContact() bool`

HasContact returns a boolean if a field has been set.

### GetCapabilities

`func (o *Site) GetCapabilities() SiteCapabilities`

GetCapabilities returns the Capabilities field if non-nil, zero value otherwise.

### GetCapabilitiesOk

`func (o *Site) GetCapabilitiesOk() (*SiteCapabilities, bool)`

GetCapabilitiesOk returns a tuple with the Capabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapabilities

`func (o *Site) SetCapabilities(v SiteCapabilities)`

SetCapabilities sets Capabilities field to given value.

### HasCapabilities

`func (o *Site) HasCapabilities() bool`

HasCapabilities returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


