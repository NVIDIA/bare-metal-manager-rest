# OperatingSystem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | ID of the Operating System | [optional] [readonly] 
**Name** | Pointer to **string** | Name of the Operating System | [optional] 
**Description** | Pointer to **string** | Optional description of the Operating System | [optional] 
**InfrastructureProviderId** | Pointer to **NullableString** | Specified if a Provider owns the Operating System | [optional] 
**TenantId** | Pointer to **NullableString** | Specified if a Tenant owns the Operating System | [optional] 
**Type** | Pointer to **string** | Type of the Operating System | [optional] 
**ImageUrl** | Pointer to **NullableString** | Original URL from where the Operating System image can be retreived from | [optional] 
**ImageSha** | Pointer to **NullableString** | SHA hash of the image file, only present for image based OS | [optional] 
**ImageAuthType** | Pointer to **NullableString** | Authentication type for image URL e.g. &#39;Basic&#39; or &#39;Bearer&#39; | [optional] 
**ImageAuthToken** | Pointer to **NullableString** | Auth token to retrieve the image from image URL | [optional] 
**ImageDisk** | Pointer to **NullableString** | Disk path where the image should be monuted | [optional] 
**RootFsId** | Pointer to **NullableString** | Root filesystem UUID, only applicable for image based Operating System | [optional] 
**RootFsLabel** | Pointer to **NullableString** | Root filesystem label, only applicable for image based Operating System | [optional] 
**IpxeScript** | Pointer to **NullableString** | iPXE script or URL, only applicable for iPXE based Operating System | [optional] 
**UserData** | Pointer to **NullableString** | User data for the Operating System | [optional] 
**IsCloudInit** | Pointer to **bool** | Specified when the Operating System is Cloud Init based | [optional] 
**PhoneHomeEnabled** | Pointer to **bool** | Indicates whether the Phone Home service should be enabled or disabled for Operating System | [optional] 
**IsActive** | Pointer to **bool** | Indicates if the Operating System is active | [optional] 
**DeactivationNote** | Pointer to **NullableString** | Optional deactivation note if OS is inactive | [optional] 
**AllowOverride** | Pointer to **bool** | Indicates if the user data can be overridden at Instance creation time | [optional] 
**SiteAssociations** | Pointer to [**[]OperatingSystemSiteAssociation**](OperatingSystemSiteAssociation.md) | Sites the Operating System is synced to | [optional] 
**Status** | Pointer to [**OperatingSystemStatus**](OperatingSystemStatus.md) | Status of the Operating System | [optional] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) | History of status changes over time | [optional] 
**Created** | Pointer to **time.Time** | Date/time when the Operating System was created | [optional] [readonly] 
**Updated** | Pointer to **time.Time** | Date/time when the Operating System was updated | [optional] [readonly] 

## Methods

### NewOperatingSystem

`func NewOperatingSystem() *OperatingSystem`

NewOperatingSystem instantiates a new OperatingSystem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOperatingSystemWithDefaults

`func NewOperatingSystemWithDefaults() *OperatingSystem`

NewOperatingSystemWithDefaults instantiates a new OperatingSystem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *OperatingSystem) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *OperatingSystem) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *OperatingSystem) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *OperatingSystem) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *OperatingSystem) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OperatingSystem) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OperatingSystem) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *OperatingSystem) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *OperatingSystem) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *OperatingSystem) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *OperatingSystem) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *OperatingSystem) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *OperatingSystem) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *OperatingSystem) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *OperatingSystem) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *OperatingSystem) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### SetInfrastructureProviderIdNil

`func (o *OperatingSystem) SetInfrastructureProviderIdNil(b bool)`

 SetInfrastructureProviderIdNil sets the value for InfrastructureProviderId to be an explicit nil

### UnsetInfrastructureProviderId
`func (o *OperatingSystem) UnsetInfrastructureProviderId()`

UnsetInfrastructureProviderId ensures that no value is present for InfrastructureProviderId, not even an explicit nil
### GetTenantId

`func (o *OperatingSystem) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *OperatingSystem) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *OperatingSystem) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *OperatingSystem) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### SetTenantIdNil

`func (o *OperatingSystem) SetTenantIdNil(b bool)`

 SetTenantIdNil sets the value for TenantId to be an explicit nil

### UnsetTenantId
`func (o *OperatingSystem) UnsetTenantId()`

UnsetTenantId ensures that no value is present for TenantId, not even an explicit nil
### GetType

`func (o *OperatingSystem) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *OperatingSystem) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *OperatingSystem) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *OperatingSystem) HasType() bool`

HasType returns a boolean if a field has been set.

### GetImageUrl

`func (o *OperatingSystem) GetImageUrl() string`

GetImageUrl returns the ImageUrl field if non-nil, zero value otherwise.

### GetImageUrlOk

`func (o *OperatingSystem) GetImageUrlOk() (*string, bool)`

GetImageUrlOk returns a tuple with the ImageUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageUrl

`func (o *OperatingSystem) SetImageUrl(v string)`

SetImageUrl sets ImageUrl field to given value.

### HasImageUrl

`func (o *OperatingSystem) HasImageUrl() bool`

HasImageUrl returns a boolean if a field has been set.

### SetImageUrlNil

`func (o *OperatingSystem) SetImageUrlNil(b bool)`

 SetImageUrlNil sets the value for ImageUrl to be an explicit nil

### UnsetImageUrl
`func (o *OperatingSystem) UnsetImageUrl()`

UnsetImageUrl ensures that no value is present for ImageUrl, not even an explicit nil
### GetImageSha

`func (o *OperatingSystem) GetImageSha() string`

GetImageSha returns the ImageSha field if non-nil, zero value otherwise.

### GetImageShaOk

`func (o *OperatingSystem) GetImageShaOk() (*string, bool)`

GetImageShaOk returns a tuple with the ImageSha field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageSha

`func (o *OperatingSystem) SetImageSha(v string)`

SetImageSha sets ImageSha field to given value.

### HasImageSha

`func (o *OperatingSystem) HasImageSha() bool`

HasImageSha returns a boolean if a field has been set.

### SetImageShaNil

`func (o *OperatingSystem) SetImageShaNil(b bool)`

 SetImageShaNil sets the value for ImageSha to be an explicit nil

### UnsetImageSha
`func (o *OperatingSystem) UnsetImageSha()`

UnsetImageSha ensures that no value is present for ImageSha, not even an explicit nil
### GetImageAuthType

`func (o *OperatingSystem) GetImageAuthType() string`

GetImageAuthType returns the ImageAuthType field if non-nil, zero value otherwise.

### GetImageAuthTypeOk

`func (o *OperatingSystem) GetImageAuthTypeOk() (*string, bool)`

GetImageAuthTypeOk returns a tuple with the ImageAuthType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageAuthType

`func (o *OperatingSystem) SetImageAuthType(v string)`

SetImageAuthType sets ImageAuthType field to given value.

### HasImageAuthType

`func (o *OperatingSystem) HasImageAuthType() bool`

HasImageAuthType returns a boolean if a field has been set.

### SetImageAuthTypeNil

`func (o *OperatingSystem) SetImageAuthTypeNil(b bool)`

 SetImageAuthTypeNil sets the value for ImageAuthType to be an explicit nil

### UnsetImageAuthType
`func (o *OperatingSystem) UnsetImageAuthType()`

UnsetImageAuthType ensures that no value is present for ImageAuthType, not even an explicit nil
### GetImageAuthToken

`func (o *OperatingSystem) GetImageAuthToken() string`

GetImageAuthToken returns the ImageAuthToken field if non-nil, zero value otherwise.

### GetImageAuthTokenOk

`func (o *OperatingSystem) GetImageAuthTokenOk() (*string, bool)`

GetImageAuthTokenOk returns a tuple with the ImageAuthToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageAuthToken

`func (o *OperatingSystem) SetImageAuthToken(v string)`

SetImageAuthToken sets ImageAuthToken field to given value.

### HasImageAuthToken

`func (o *OperatingSystem) HasImageAuthToken() bool`

HasImageAuthToken returns a boolean if a field has been set.

### SetImageAuthTokenNil

`func (o *OperatingSystem) SetImageAuthTokenNil(b bool)`

 SetImageAuthTokenNil sets the value for ImageAuthToken to be an explicit nil

### UnsetImageAuthToken
`func (o *OperatingSystem) UnsetImageAuthToken()`

UnsetImageAuthToken ensures that no value is present for ImageAuthToken, not even an explicit nil
### GetImageDisk

`func (o *OperatingSystem) GetImageDisk() string`

GetImageDisk returns the ImageDisk field if non-nil, zero value otherwise.

### GetImageDiskOk

`func (o *OperatingSystem) GetImageDiskOk() (*string, bool)`

GetImageDiskOk returns a tuple with the ImageDisk field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageDisk

`func (o *OperatingSystem) SetImageDisk(v string)`

SetImageDisk sets ImageDisk field to given value.

### HasImageDisk

`func (o *OperatingSystem) HasImageDisk() bool`

HasImageDisk returns a boolean if a field has been set.

### SetImageDiskNil

`func (o *OperatingSystem) SetImageDiskNil(b bool)`

 SetImageDiskNil sets the value for ImageDisk to be an explicit nil

### UnsetImageDisk
`func (o *OperatingSystem) UnsetImageDisk()`

UnsetImageDisk ensures that no value is present for ImageDisk, not even an explicit nil
### GetRootFsId

`func (o *OperatingSystem) GetRootFsId() string`

GetRootFsId returns the RootFsId field if non-nil, zero value otherwise.

### GetRootFsIdOk

`func (o *OperatingSystem) GetRootFsIdOk() (*string, bool)`

GetRootFsIdOk returns a tuple with the RootFsId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRootFsId

`func (o *OperatingSystem) SetRootFsId(v string)`

SetRootFsId sets RootFsId field to given value.

### HasRootFsId

`func (o *OperatingSystem) HasRootFsId() bool`

HasRootFsId returns a boolean if a field has been set.

### SetRootFsIdNil

`func (o *OperatingSystem) SetRootFsIdNil(b bool)`

 SetRootFsIdNil sets the value for RootFsId to be an explicit nil

### UnsetRootFsId
`func (o *OperatingSystem) UnsetRootFsId()`

UnsetRootFsId ensures that no value is present for RootFsId, not even an explicit nil
### GetRootFsLabel

`func (o *OperatingSystem) GetRootFsLabel() string`

GetRootFsLabel returns the RootFsLabel field if non-nil, zero value otherwise.

### GetRootFsLabelOk

`func (o *OperatingSystem) GetRootFsLabelOk() (*string, bool)`

GetRootFsLabelOk returns a tuple with the RootFsLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRootFsLabel

`func (o *OperatingSystem) SetRootFsLabel(v string)`

SetRootFsLabel sets RootFsLabel field to given value.

### HasRootFsLabel

`func (o *OperatingSystem) HasRootFsLabel() bool`

HasRootFsLabel returns a boolean if a field has been set.

### SetRootFsLabelNil

`func (o *OperatingSystem) SetRootFsLabelNil(b bool)`

 SetRootFsLabelNil sets the value for RootFsLabel to be an explicit nil

### UnsetRootFsLabel
`func (o *OperatingSystem) UnsetRootFsLabel()`

UnsetRootFsLabel ensures that no value is present for RootFsLabel, not even an explicit nil
### GetIpxeScript

`func (o *OperatingSystem) GetIpxeScript() string`

GetIpxeScript returns the IpxeScript field if non-nil, zero value otherwise.

### GetIpxeScriptOk

`func (o *OperatingSystem) GetIpxeScriptOk() (*string, bool)`

GetIpxeScriptOk returns a tuple with the IpxeScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpxeScript

`func (o *OperatingSystem) SetIpxeScript(v string)`

SetIpxeScript sets IpxeScript field to given value.

### HasIpxeScript

`func (o *OperatingSystem) HasIpxeScript() bool`

HasIpxeScript returns a boolean if a field has been set.

### SetIpxeScriptNil

`func (o *OperatingSystem) SetIpxeScriptNil(b bool)`

 SetIpxeScriptNil sets the value for IpxeScript to be an explicit nil

### UnsetIpxeScript
`func (o *OperatingSystem) UnsetIpxeScript()`

UnsetIpxeScript ensures that no value is present for IpxeScript, not even an explicit nil
### GetUserData

`func (o *OperatingSystem) GetUserData() string`

GetUserData returns the UserData field if non-nil, zero value otherwise.

### GetUserDataOk

`func (o *OperatingSystem) GetUserDataOk() (*string, bool)`

GetUserDataOk returns a tuple with the UserData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserData

`func (o *OperatingSystem) SetUserData(v string)`

SetUserData sets UserData field to given value.

### HasUserData

`func (o *OperatingSystem) HasUserData() bool`

HasUserData returns a boolean if a field has been set.

### SetUserDataNil

`func (o *OperatingSystem) SetUserDataNil(b bool)`

 SetUserDataNil sets the value for UserData to be an explicit nil

### UnsetUserData
`func (o *OperatingSystem) UnsetUserData()`

UnsetUserData ensures that no value is present for UserData, not even an explicit nil
### GetIsCloudInit

`func (o *OperatingSystem) GetIsCloudInit() bool`

GetIsCloudInit returns the IsCloudInit field if non-nil, zero value otherwise.

### GetIsCloudInitOk

`func (o *OperatingSystem) GetIsCloudInitOk() (*bool, bool)`

GetIsCloudInitOk returns a tuple with the IsCloudInit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsCloudInit

`func (o *OperatingSystem) SetIsCloudInit(v bool)`

SetIsCloudInit sets IsCloudInit field to given value.

### HasIsCloudInit

`func (o *OperatingSystem) HasIsCloudInit() bool`

HasIsCloudInit returns a boolean if a field has been set.

### GetPhoneHomeEnabled

`func (o *OperatingSystem) GetPhoneHomeEnabled() bool`

GetPhoneHomeEnabled returns the PhoneHomeEnabled field if non-nil, zero value otherwise.

### GetPhoneHomeEnabledOk

`func (o *OperatingSystem) GetPhoneHomeEnabledOk() (*bool, bool)`

GetPhoneHomeEnabledOk returns a tuple with the PhoneHomeEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneHomeEnabled

`func (o *OperatingSystem) SetPhoneHomeEnabled(v bool)`

SetPhoneHomeEnabled sets PhoneHomeEnabled field to given value.

### HasPhoneHomeEnabled

`func (o *OperatingSystem) HasPhoneHomeEnabled() bool`

HasPhoneHomeEnabled returns a boolean if a field has been set.

### GetIsActive

`func (o *OperatingSystem) GetIsActive() bool`

GetIsActive returns the IsActive field if non-nil, zero value otherwise.

### GetIsActiveOk

`func (o *OperatingSystem) GetIsActiveOk() (*bool, bool)`

GetIsActiveOk returns a tuple with the IsActive field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsActive

`func (o *OperatingSystem) SetIsActive(v bool)`

SetIsActive sets IsActive field to given value.

### HasIsActive

`func (o *OperatingSystem) HasIsActive() bool`

HasIsActive returns a boolean if a field has been set.

### GetDeactivationNote

`func (o *OperatingSystem) GetDeactivationNote() string`

GetDeactivationNote returns the DeactivationNote field if non-nil, zero value otherwise.

### GetDeactivationNoteOk

`func (o *OperatingSystem) GetDeactivationNoteOk() (*string, bool)`

GetDeactivationNoteOk returns a tuple with the DeactivationNote field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeactivationNote

`func (o *OperatingSystem) SetDeactivationNote(v string)`

SetDeactivationNote sets DeactivationNote field to given value.

### HasDeactivationNote

`func (o *OperatingSystem) HasDeactivationNote() bool`

HasDeactivationNote returns a boolean if a field has been set.

### SetDeactivationNoteNil

`func (o *OperatingSystem) SetDeactivationNoteNil(b bool)`

 SetDeactivationNoteNil sets the value for DeactivationNote to be an explicit nil

### UnsetDeactivationNote
`func (o *OperatingSystem) UnsetDeactivationNote()`

UnsetDeactivationNote ensures that no value is present for DeactivationNote, not even an explicit nil
### GetAllowOverride

`func (o *OperatingSystem) GetAllowOverride() bool`

GetAllowOverride returns the AllowOverride field if non-nil, zero value otherwise.

### GetAllowOverrideOk

`func (o *OperatingSystem) GetAllowOverrideOk() (*bool, bool)`

GetAllowOverrideOk returns a tuple with the AllowOverride field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowOverride

`func (o *OperatingSystem) SetAllowOverride(v bool)`

SetAllowOverride sets AllowOverride field to given value.

### HasAllowOverride

`func (o *OperatingSystem) HasAllowOverride() bool`

HasAllowOverride returns a boolean if a field has been set.

### GetSiteAssociations

`func (o *OperatingSystem) GetSiteAssociations() []OperatingSystemSiteAssociation`

GetSiteAssociations returns the SiteAssociations field if non-nil, zero value otherwise.

### GetSiteAssociationsOk

`func (o *OperatingSystem) GetSiteAssociationsOk() (*[]OperatingSystemSiteAssociation, bool)`

GetSiteAssociationsOk returns a tuple with the SiteAssociations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteAssociations

`func (o *OperatingSystem) SetSiteAssociations(v []OperatingSystemSiteAssociation)`

SetSiteAssociations sets SiteAssociations field to given value.

### HasSiteAssociations

`func (o *OperatingSystem) HasSiteAssociations() bool`

HasSiteAssociations returns a boolean if a field has been set.

### GetStatus

`func (o *OperatingSystem) GetStatus() OperatingSystemStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *OperatingSystem) GetStatusOk() (*OperatingSystemStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *OperatingSystem) SetStatus(v OperatingSystemStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *OperatingSystem) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *OperatingSystem) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *OperatingSystem) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *OperatingSystem) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *OperatingSystem) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *OperatingSystem) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *OperatingSystem) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *OperatingSystem) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *OperatingSystem) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *OperatingSystem) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *OperatingSystem) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *OperatingSystem) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *OperatingSystem) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


