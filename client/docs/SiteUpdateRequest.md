# SiteUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Update name of the Site. Can only be updated by Provider | [optional] 
**Description** | Pointer to **string** | Update description for the Site. Can only be updated by Provider | [optional] 
**RenewRegistrationToken** | Pointer to **bool** | Set to true to issue a new registration token. Can only be updated by Provider | [optional] 
**SerialConsoleHostname** | Pointer to **string** | Hostname to reach Serial Console for the Site. Can only be updated by Provider | [optional] 
**IsSerialConsoleEnabled** | Pointer to **bool** | Enable/disable Serial Console. Can only be updated by Provider | [optional] 
**SerialConsoleIdleTimeout** | Pointer to **int32** | Maximum idle time in seconds before Serial Console is disconnected. Can only be updated by Provider | [optional] 
**SerialConsoleMaxSessionLength** | Pointer to **int32** | Maximum length of Serial Console session in seconds. Can only be updated by Provider | [optional] 
**IsSerialConsoleSSHKeysEnabled** | Pointer to **bool** | Enable/disable Serial Console access using SSH Keys. Can only be updated by Tenant | [optional] 
**Location** | Pointer to [**SiteLocation**](SiteLocation.md) |  | [optional] 
**Contact** | Pointer to [**SiteContact**](SiteContact.md) |  | [optional] 

## Methods

### NewSiteUpdateRequest

`func NewSiteUpdateRequest() *SiteUpdateRequest`

NewSiteUpdateRequest instantiates a new SiteUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSiteUpdateRequestWithDefaults

`func NewSiteUpdateRequestWithDefaults() *SiteUpdateRequest`

NewSiteUpdateRequestWithDefaults instantiates a new SiteUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SiteUpdateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SiteUpdateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SiteUpdateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SiteUpdateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *SiteUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *SiteUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *SiteUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *SiteUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetRenewRegistrationToken

`func (o *SiteUpdateRequest) GetRenewRegistrationToken() bool`

GetRenewRegistrationToken returns the RenewRegistrationToken field if non-nil, zero value otherwise.

### GetRenewRegistrationTokenOk

`func (o *SiteUpdateRequest) GetRenewRegistrationTokenOk() (*bool, bool)`

GetRenewRegistrationTokenOk returns a tuple with the RenewRegistrationToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRenewRegistrationToken

`func (o *SiteUpdateRequest) SetRenewRegistrationToken(v bool)`

SetRenewRegistrationToken sets RenewRegistrationToken field to given value.

### HasRenewRegistrationToken

`func (o *SiteUpdateRequest) HasRenewRegistrationToken() bool`

HasRenewRegistrationToken returns a boolean if a field has been set.

### GetSerialConsoleHostname

`func (o *SiteUpdateRequest) GetSerialConsoleHostname() string`

GetSerialConsoleHostname returns the SerialConsoleHostname field if non-nil, zero value otherwise.

### GetSerialConsoleHostnameOk

`func (o *SiteUpdateRequest) GetSerialConsoleHostnameOk() (*string, bool)`

GetSerialConsoleHostnameOk returns a tuple with the SerialConsoleHostname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialConsoleHostname

`func (o *SiteUpdateRequest) SetSerialConsoleHostname(v string)`

SetSerialConsoleHostname sets SerialConsoleHostname field to given value.

### HasSerialConsoleHostname

`func (o *SiteUpdateRequest) HasSerialConsoleHostname() bool`

HasSerialConsoleHostname returns a boolean if a field has been set.

### GetIsSerialConsoleEnabled

`func (o *SiteUpdateRequest) GetIsSerialConsoleEnabled() bool`

GetIsSerialConsoleEnabled returns the IsSerialConsoleEnabled field if non-nil, zero value otherwise.

### GetIsSerialConsoleEnabledOk

`func (o *SiteUpdateRequest) GetIsSerialConsoleEnabledOk() (*bool, bool)`

GetIsSerialConsoleEnabledOk returns a tuple with the IsSerialConsoleEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsSerialConsoleEnabled

`func (o *SiteUpdateRequest) SetIsSerialConsoleEnabled(v bool)`

SetIsSerialConsoleEnabled sets IsSerialConsoleEnabled field to given value.

### HasIsSerialConsoleEnabled

`func (o *SiteUpdateRequest) HasIsSerialConsoleEnabled() bool`

HasIsSerialConsoleEnabled returns a boolean if a field has been set.

### GetSerialConsoleIdleTimeout

`func (o *SiteUpdateRequest) GetSerialConsoleIdleTimeout() int32`

GetSerialConsoleIdleTimeout returns the SerialConsoleIdleTimeout field if non-nil, zero value otherwise.

### GetSerialConsoleIdleTimeoutOk

`func (o *SiteUpdateRequest) GetSerialConsoleIdleTimeoutOk() (*int32, bool)`

GetSerialConsoleIdleTimeoutOk returns a tuple with the SerialConsoleIdleTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialConsoleIdleTimeout

`func (o *SiteUpdateRequest) SetSerialConsoleIdleTimeout(v int32)`

SetSerialConsoleIdleTimeout sets SerialConsoleIdleTimeout field to given value.

### HasSerialConsoleIdleTimeout

`func (o *SiteUpdateRequest) HasSerialConsoleIdleTimeout() bool`

HasSerialConsoleIdleTimeout returns a boolean if a field has been set.

### GetSerialConsoleMaxSessionLength

`func (o *SiteUpdateRequest) GetSerialConsoleMaxSessionLength() int32`

GetSerialConsoleMaxSessionLength returns the SerialConsoleMaxSessionLength field if non-nil, zero value otherwise.

### GetSerialConsoleMaxSessionLengthOk

`func (o *SiteUpdateRequest) GetSerialConsoleMaxSessionLengthOk() (*int32, bool)`

GetSerialConsoleMaxSessionLengthOk returns a tuple with the SerialConsoleMaxSessionLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialConsoleMaxSessionLength

`func (o *SiteUpdateRequest) SetSerialConsoleMaxSessionLength(v int32)`

SetSerialConsoleMaxSessionLength sets SerialConsoleMaxSessionLength field to given value.

### HasSerialConsoleMaxSessionLength

`func (o *SiteUpdateRequest) HasSerialConsoleMaxSessionLength() bool`

HasSerialConsoleMaxSessionLength returns a boolean if a field has been set.

### GetIsSerialConsoleSSHKeysEnabled

`func (o *SiteUpdateRequest) GetIsSerialConsoleSSHKeysEnabled() bool`

GetIsSerialConsoleSSHKeysEnabled returns the IsSerialConsoleSSHKeysEnabled field if non-nil, zero value otherwise.

### GetIsSerialConsoleSSHKeysEnabledOk

`func (o *SiteUpdateRequest) GetIsSerialConsoleSSHKeysEnabledOk() (*bool, bool)`

GetIsSerialConsoleSSHKeysEnabledOk returns a tuple with the IsSerialConsoleSSHKeysEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsSerialConsoleSSHKeysEnabled

`func (o *SiteUpdateRequest) SetIsSerialConsoleSSHKeysEnabled(v bool)`

SetIsSerialConsoleSSHKeysEnabled sets IsSerialConsoleSSHKeysEnabled field to given value.

### HasIsSerialConsoleSSHKeysEnabled

`func (o *SiteUpdateRequest) HasIsSerialConsoleSSHKeysEnabled() bool`

HasIsSerialConsoleSSHKeysEnabled returns a boolean if a field has been set.

### GetLocation

`func (o *SiteUpdateRequest) GetLocation() SiteLocation`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *SiteUpdateRequest) GetLocationOk() (*SiteLocation, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *SiteUpdateRequest) SetLocation(v SiteLocation)`

SetLocation sets Location field to given value.

### HasLocation

`func (o *SiteUpdateRequest) HasLocation() bool`

HasLocation returns a boolean if a field has been set.

### GetContact

`func (o *SiteUpdateRequest) GetContact() SiteContact`

GetContact returns the Contact field if non-nil, zero value otherwise.

### GetContactOk

`func (o *SiteUpdateRequest) GetContactOk() (*SiteContact, bool)`

GetContactOk returns a tuple with the Contact field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContact

`func (o *SiteUpdateRequest) SetContact(v SiteContact)`

SetContact sets Contact field to given value.

### HasContact

`func (o *SiteUpdateRequest) HasContact() bool`

HasContact returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


