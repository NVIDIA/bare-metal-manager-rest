# OperatingSystemUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **NullableString** | Name of the Operating System | [optional] 
**Description** | Pointer to **NullableString** | Optional description of the Operating System | [optional] 
**IpxeScript** | Pointer to **NullableString** | iPXE script or URL, only applicable for iPXE based OS. Cannot be specified if imageUrl is specified | [optional] 
**ImageUrl** | Pointer to **NullableString** | Original URL from where the Operating System image can be retreived from, required for image based OS | [optional] 
**ImageSha** | Pointer to **NullableString** | SHA hash of the image file, required for image based OS | [optional] 
**ImageAuthType** | Pointer to **NullableString** | Authentication type for image URL if needed e.g. basic/bearer/token, required is imageAuthToken is specified | [optional] 
**ImageAuthToken** | Pointer to **NullableString** | Auth token to retrieve the image from image URL, required if imageAuthType is specified | [optional] 
**ImageDisk** | Pointer to **NullableString** | Disk path where the image should be mounted, optional | [optional] 
**RootFsId** | Pointer to **NullableString** | Root filesystem UUID, this or &#x60;rootFsLabel&#x60; required for image based OS | [optional] 
**RootFsLabel** | Pointer to **NullableString** | Root filesystem label, this or &#x60;rootFsId&#x60; required for image based OS | [optional] 
**PhoneHomeEnabled** | Pointer to **NullableBool** | Indicates whether the Phone Home service should be enabled or disabled for Operating System | [optional] 
**UserData** | Pointer to **NullableString** | User data for the Operating System | [optional] 
**IsCloudInit** | Pointer to **NullableBool** | Specified when the Operating System is Cloud Init based | [optional] 
**AllowOverride** | Pointer to **NullableBool** | Indicates if the user data can be overridden at Instance creation time | [optional] 
**IsActive** | Pointer to **NullableBool** | Indicates if the Operating System is active | [optional] 
**DeactivationNote** | Pointer to **NullableString** | Optional deactivation note if OS is inactive | [optional] 

## Methods

### NewOperatingSystemUpdateRequest

`func NewOperatingSystemUpdateRequest() *OperatingSystemUpdateRequest`

NewOperatingSystemUpdateRequest instantiates a new OperatingSystemUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOperatingSystemUpdateRequestWithDefaults

`func NewOperatingSystemUpdateRequestWithDefaults() *OperatingSystemUpdateRequest`

NewOperatingSystemUpdateRequestWithDefaults instantiates a new OperatingSystemUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *OperatingSystemUpdateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OperatingSystemUpdateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OperatingSystemUpdateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *OperatingSystemUpdateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### SetNameNil

`func (o *OperatingSystemUpdateRequest) SetNameNil(b bool)`

 SetNameNil sets the value for Name to be an explicit nil

### UnsetName
`func (o *OperatingSystemUpdateRequest) UnsetName()`

UnsetName ensures that no value is present for Name, not even an explicit nil
### GetDescription

`func (o *OperatingSystemUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *OperatingSystemUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *OperatingSystemUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *OperatingSystemUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *OperatingSystemUpdateRequest) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *OperatingSystemUpdateRequest) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetIpxeScript

`func (o *OperatingSystemUpdateRequest) GetIpxeScript() string`

GetIpxeScript returns the IpxeScript field if non-nil, zero value otherwise.

### GetIpxeScriptOk

`func (o *OperatingSystemUpdateRequest) GetIpxeScriptOk() (*string, bool)`

GetIpxeScriptOk returns a tuple with the IpxeScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpxeScript

`func (o *OperatingSystemUpdateRequest) SetIpxeScript(v string)`

SetIpxeScript sets IpxeScript field to given value.

### HasIpxeScript

`func (o *OperatingSystemUpdateRequest) HasIpxeScript() bool`

HasIpxeScript returns a boolean if a field has been set.

### SetIpxeScriptNil

`func (o *OperatingSystemUpdateRequest) SetIpxeScriptNil(b bool)`

 SetIpxeScriptNil sets the value for IpxeScript to be an explicit nil

### UnsetIpxeScript
`func (o *OperatingSystemUpdateRequest) UnsetIpxeScript()`

UnsetIpxeScript ensures that no value is present for IpxeScript, not even an explicit nil
### GetImageUrl

`func (o *OperatingSystemUpdateRequest) GetImageUrl() string`

GetImageUrl returns the ImageUrl field if non-nil, zero value otherwise.

### GetImageUrlOk

`func (o *OperatingSystemUpdateRequest) GetImageUrlOk() (*string, bool)`

GetImageUrlOk returns a tuple with the ImageUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageUrl

`func (o *OperatingSystemUpdateRequest) SetImageUrl(v string)`

SetImageUrl sets ImageUrl field to given value.

### HasImageUrl

`func (o *OperatingSystemUpdateRequest) HasImageUrl() bool`

HasImageUrl returns a boolean if a field has been set.

### SetImageUrlNil

`func (o *OperatingSystemUpdateRequest) SetImageUrlNil(b bool)`

 SetImageUrlNil sets the value for ImageUrl to be an explicit nil

### UnsetImageUrl
`func (o *OperatingSystemUpdateRequest) UnsetImageUrl()`

UnsetImageUrl ensures that no value is present for ImageUrl, not even an explicit nil
### GetImageSha

`func (o *OperatingSystemUpdateRequest) GetImageSha() string`

GetImageSha returns the ImageSha field if non-nil, zero value otherwise.

### GetImageShaOk

`func (o *OperatingSystemUpdateRequest) GetImageShaOk() (*string, bool)`

GetImageShaOk returns a tuple with the ImageSha field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageSha

`func (o *OperatingSystemUpdateRequest) SetImageSha(v string)`

SetImageSha sets ImageSha field to given value.

### HasImageSha

`func (o *OperatingSystemUpdateRequest) HasImageSha() bool`

HasImageSha returns a boolean if a field has been set.

### SetImageShaNil

`func (o *OperatingSystemUpdateRequest) SetImageShaNil(b bool)`

 SetImageShaNil sets the value for ImageSha to be an explicit nil

### UnsetImageSha
`func (o *OperatingSystemUpdateRequest) UnsetImageSha()`

UnsetImageSha ensures that no value is present for ImageSha, not even an explicit nil
### GetImageAuthType

`func (o *OperatingSystemUpdateRequest) GetImageAuthType() string`

GetImageAuthType returns the ImageAuthType field if non-nil, zero value otherwise.

### GetImageAuthTypeOk

`func (o *OperatingSystemUpdateRequest) GetImageAuthTypeOk() (*string, bool)`

GetImageAuthTypeOk returns a tuple with the ImageAuthType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageAuthType

`func (o *OperatingSystemUpdateRequest) SetImageAuthType(v string)`

SetImageAuthType sets ImageAuthType field to given value.

### HasImageAuthType

`func (o *OperatingSystemUpdateRequest) HasImageAuthType() bool`

HasImageAuthType returns a boolean if a field has been set.

### SetImageAuthTypeNil

`func (o *OperatingSystemUpdateRequest) SetImageAuthTypeNil(b bool)`

 SetImageAuthTypeNil sets the value for ImageAuthType to be an explicit nil

### UnsetImageAuthType
`func (o *OperatingSystemUpdateRequest) UnsetImageAuthType()`

UnsetImageAuthType ensures that no value is present for ImageAuthType, not even an explicit nil
### GetImageAuthToken

`func (o *OperatingSystemUpdateRequest) GetImageAuthToken() string`

GetImageAuthToken returns the ImageAuthToken field if non-nil, zero value otherwise.

### GetImageAuthTokenOk

`func (o *OperatingSystemUpdateRequest) GetImageAuthTokenOk() (*string, bool)`

GetImageAuthTokenOk returns a tuple with the ImageAuthToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageAuthToken

`func (o *OperatingSystemUpdateRequest) SetImageAuthToken(v string)`

SetImageAuthToken sets ImageAuthToken field to given value.

### HasImageAuthToken

`func (o *OperatingSystemUpdateRequest) HasImageAuthToken() bool`

HasImageAuthToken returns a boolean if a field has been set.

### SetImageAuthTokenNil

`func (o *OperatingSystemUpdateRequest) SetImageAuthTokenNil(b bool)`

 SetImageAuthTokenNil sets the value for ImageAuthToken to be an explicit nil

### UnsetImageAuthToken
`func (o *OperatingSystemUpdateRequest) UnsetImageAuthToken()`

UnsetImageAuthToken ensures that no value is present for ImageAuthToken, not even an explicit nil
### GetImageDisk

`func (o *OperatingSystemUpdateRequest) GetImageDisk() string`

GetImageDisk returns the ImageDisk field if non-nil, zero value otherwise.

### GetImageDiskOk

`func (o *OperatingSystemUpdateRequest) GetImageDiskOk() (*string, bool)`

GetImageDiskOk returns a tuple with the ImageDisk field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageDisk

`func (o *OperatingSystemUpdateRequest) SetImageDisk(v string)`

SetImageDisk sets ImageDisk field to given value.

### HasImageDisk

`func (o *OperatingSystemUpdateRequest) HasImageDisk() bool`

HasImageDisk returns a boolean if a field has been set.

### SetImageDiskNil

`func (o *OperatingSystemUpdateRequest) SetImageDiskNil(b bool)`

 SetImageDiskNil sets the value for ImageDisk to be an explicit nil

### UnsetImageDisk
`func (o *OperatingSystemUpdateRequest) UnsetImageDisk()`

UnsetImageDisk ensures that no value is present for ImageDisk, not even an explicit nil
### GetRootFsId

`func (o *OperatingSystemUpdateRequest) GetRootFsId() string`

GetRootFsId returns the RootFsId field if non-nil, zero value otherwise.

### GetRootFsIdOk

`func (o *OperatingSystemUpdateRequest) GetRootFsIdOk() (*string, bool)`

GetRootFsIdOk returns a tuple with the RootFsId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRootFsId

`func (o *OperatingSystemUpdateRequest) SetRootFsId(v string)`

SetRootFsId sets RootFsId field to given value.

### HasRootFsId

`func (o *OperatingSystemUpdateRequest) HasRootFsId() bool`

HasRootFsId returns a boolean if a field has been set.

### SetRootFsIdNil

`func (o *OperatingSystemUpdateRequest) SetRootFsIdNil(b bool)`

 SetRootFsIdNil sets the value for RootFsId to be an explicit nil

### UnsetRootFsId
`func (o *OperatingSystemUpdateRequest) UnsetRootFsId()`

UnsetRootFsId ensures that no value is present for RootFsId, not even an explicit nil
### GetRootFsLabel

`func (o *OperatingSystemUpdateRequest) GetRootFsLabel() string`

GetRootFsLabel returns the RootFsLabel field if non-nil, zero value otherwise.

### GetRootFsLabelOk

`func (o *OperatingSystemUpdateRequest) GetRootFsLabelOk() (*string, bool)`

GetRootFsLabelOk returns a tuple with the RootFsLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRootFsLabel

`func (o *OperatingSystemUpdateRequest) SetRootFsLabel(v string)`

SetRootFsLabel sets RootFsLabel field to given value.

### HasRootFsLabel

`func (o *OperatingSystemUpdateRequest) HasRootFsLabel() bool`

HasRootFsLabel returns a boolean if a field has been set.

### SetRootFsLabelNil

`func (o *OperatingSystemUpdateRequest) SetRootFsLabelNil(b bool)`

 SetRootFsLabelNil sets the value for RootFsLabel to be an explicit nil

### UnsetRootFsLabel
`func (o *OperatingSystemUpdateRequest) UnsetRootFsLabel()`

UnsetRootFsLabel ensures that no value is present for RootFsLabel, not even an explicit nil
### GetPhoneHomeEnabled

`func (o *OperatingSystemUpdateRequest) GetPhoneHomeEnabled() bool`

GetPhoneHomeEnabled returns the PhoneHomeEnabled field if non-nil, zero value otherwise.

### GetPhoneHomeEnabledOk

`func (o *OperatingSystemUpdateRequest) GetPhoneHomeEnabledOk() (*bool, bool)`

GetPhoneHomeEnabledOk returns a tuple with the PhoneHomeEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneHomeEnabled

`func (o *OperatingSystemUpdateRequest) SetPhoneHomeEnabled(v bool)`

SetPhoneHomeEnabled sets PhoneHomeEnabled field to given value.

### HasPhoneHomeEnabled

`func (o *OperatingSystemUpdateRequest) HasPhoneHomeEnabled() bool`

HasPhoneHomeEnabled returns a boolean if a field has been set.

### SetPhoneHomeEnabledNil

`func (o *OperatingSystemUpdateRequest) SetPhoneHomeEnabledNil(b bool)`

 SetPhoneHomeEnabledNil sets the value for PhoneHomeEnabled to be an explicit nil

### UnsetPhoneHomeEnabled
`func (o *OperatingSystemUpdateRequest) UnsetPhoneHomeEnabled()`

UnsetPhoneHomeEnabled ensures that no value is present for PhoneHomeEnabled, not even an explicit nil
### GetUserData

`func (o *OperatingSystemUpdateRequest) GetUserData() string`

GetUserData returns the UserData field if non-nil, zero value otherwise.

### GetUserDataOk

`func (o *OperatingSystemUpdateRequest) GetUserDataOk() (*string, bool)`

GetUserDataOk returns a tuple with the UserData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserData

`func (o *OperatingSystemUpdateRequest) SetUserData(v string)`

SetUserData sets UserData field to given value.

### HasUserData

`func (o *OperatingSystemUpdateRequest) HasUserData() bool`

HasUserData returns a boolean if a field has been set.

### SetUserDataNil

`func (o *OperatingSystemUpdateRequest) SetUserDataNil(b bool)`

 SetUserDataNil sets the value for UserData to be an explicit nil

### UnsetUserData
`func (o *OperatingSystemUpdateRequest) UnsetUserData()`

UnsetUserData ensures that no value is present for UserData, not even an explicit nil
### GetIsCloudInit

`func (o *OperatingSystemUpdateRequest) GetIsCloudInit() bool`

GetIsCloudInit returns the IsCloudInit field if non-nil, zero value otherwise.

### GetIsCloudInitOk

`func (o *OperatingSystemUpdateRequest) GetIsCloudInitOk() (*bool, bool)`

GetIsCloudInitOk returns a tuple with the IsCloudInit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsCloudInit

`func (o *OperatingSystemUpdateRequest) SetIsCloudInit(v bool)`

SetIsCloudInit sets IsCloudInit field to given value.

### HasIsCloudInit

`func (o *OperatingSystemUpdateRequest) HasIsCloudInit() bool`

HasIsCloudInit returns a boolean if a field has been set.

### SetIsCloudInitNil

`func (o *OperatingSystemUpdateRequest) SetIsCloudInitNil(b bool)`

 SetIsCloudInitNil sets the value for IsCloudInit to be an explicit nil

### UnsetIsCloudInit
`func (o *OperatingSystemUpdateRequest) UnsetIsCloudInit()`

UnsetIsCloudInit ensures that no value is present for IsCloudInit, not even an explicit nil
### GetAllowOverride

`func (o *OperatingSystemUpdateRequest) GetAllowOverride() bool`

GetAllowOverride returns the AllowOverride field if non-nil, zero value otherwise.

### GetAllowOverrideOk

`func (o *OperatingSystemUpdateRequest) GetAllowOverrideOk() (*bool, bool)`

GetAllowOverrideOk returns a tuple with the AllowOverride field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowOverride

`func (o *OperatingSystemUpdateRequest) SetAllowOverride(v bool)`

SetAllowOverride sets AllowOverride field to given value.

### HasAllowOverride

`func (o *OperatingSystemUpdateRequest) HasAllowOverride() bool`

HasAllowOverride returns a boolean if a field has been set.

### SetAllowOverrideNil

`func (o *OperatingSystemUpdateRequest) SetAllowOverrideNil(b bool)`

 SetAllowOverrideNil sets the value for AllowOverride to be an explicit nil

### UnsetAllowOverride
`func (o *OperatingSystemUpdateRequest) UnsetAllowOverride()`

UnsetAllowOverride ensures that no value is present for AllowOverride, not even an explicit nil
### GetIsActive

`func (o *OperatingSystemUpdateRequest) GetIsActive() bool`

GetIsActive returns the IsActive field if non-nil, zero value otherwise.

### GetIsActiveOk

`func (o *OperatingSystemUpdateRequest) GetIsActiveOk() (*bool, bool)`

GetIsActiveOk returns a tuple with the IsActive field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsActive

`func (o *OperatingSystemUpdateRequest) SetIsActive(v bool)`

SetIsActive sets IsActive field to given value.

### HasIsActive

`func (o *OperatingSystemUpdateRequest) HasIsActive() bool`

HasIsActive returns a boolean if a field has been set.

### SetIsActiveNil

`func (o *OperatingSystemUpdateRequest) SetIsActiveNil(b bool)`

 SetIsActiveNil sets the value for IsActive to be an explicit nil

### UnsetIsActive
`func (o *OperatingSystemUpdateRequest) UnsetIsActive()`

UnsetIsActive ensures that no value is present for IsActive, not even an explicit nil
### GetDeactivationNote

`func (o *OperatingSystemUpdateRequest) GetDeactivationNote() string`

GetDeactivationNote returns the DeactivationNote field if non-nil, zero value otherwise.

### GetDeactivationNoteOk

`func (o *OperatingSystemUpdateRequest) GetDeactivationNoteOk() (*string, bool)`

GetDeactivationNoteOk returns a tuple with the DeactivationNote field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeactivationNote

`func (o *OperatingSystemUpdateRequest) SetDeactivationNote(v string)`

SetDeactivationNote sets DeactivationNote field to given value.

### HasDeactivationNote

`func (o *OperatingSystemUpdateRequest) HasDeactivationNote() bool`

HasDeactivationNote returns a boolean if a field has been set.

### SetDeactivationNoteNil

`func (o *OperatingSystemUpdateRequest) SetDeactivationNoteNil(b bool)`

 SetDeactivationNoteNil sets the value for DeactivationNote to be an explicit nil

### UnsetDeactivationNote
`func (o *OperatingSystemUpdateRequest) UnsetDeactivationNote()`

UnsetDeactivationNote ensures that no value is present for DeactivationNote, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


