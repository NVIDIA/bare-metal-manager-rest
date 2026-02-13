# OperatingSystemCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name of the Operating System | 
**Description** | Pointer to **string** | Optional description of the Operating System | [optional] 
**InfrastructureProviderId** | Pointer to **NullableString** | Specified if a Provider owns the Operating System | [optional] 
**TenantId** | Pointer to **NullableString** | Specified if a Tenant owns the Operating System | [optional] 
**SiteIds** | Pointer to **[]string** | Specified only one Site if an Operating System is Image based, more than one Site is not supported\&quot; | [optional] 
**IpxeScript** | Pointer to **NullableString** | iPXE script or URL, only applicable for iPXE based OS. Cannot be specified if imageUrl is specified | [optional] 
**ImageUrl** | Pointer to **NullableString** | Original URL from where the Operating System image can be retreived from, required for image based OS. Cannot be specified if ipxeScript is specified | [optional] 
**ImageSha** | Pointer to **NullableString** | SHA hash of the image file, required for image based OS | [optional] 
**ImageAuthType** | Pointer to **NullableString** | Authentication type for image URL if needed e.g. basic/bearer/token, required is imageAuthToken is specified | [optional] 
**ImageAuthToken** | Pointer to **NullableString** | Auth token to retrieve the image from image URL, required if imageAuthType is specified | [optional] 
**ImageDisk** | Pointer to **NullableString** | Disk path where the image should be mounted, optional | [optional] 
**RootFsId** | Pointer to **NullableString** | Root filesystem UUID, this or &#x60;rootFsLabel&#x60; required for image based OS | [optional] 
**RootFsLabel** | Pointer to **NullableString** | Root filesystem label, this or &#x60;rootFsId&#x60; required for image based OS | [optional] 
**PhoneHomeEnabled** | Pointer to **NullableBool** | Indicates whether the Phone Home service should be enabled or disabled for Operating System | [optional] 
**UserData** | Pointer to **NullableString** | User data for the Operating System | [optional] 
**IsCloudInit** | Pointer to **bool** | Specified when the Operating System is Cloud Init based | [optional] 
**AllowOverride** | Pointer to **bool** | Indicates if the user data can be overridden at Instance creation time | [optional] 

## Methods

### NewOperatingSystemCreateRequest

`func NewOperatingSystemCreateRequest(name string, ) *OperatingSystemCreateRequest`

NewOperatingSystemCreateRequest instantiates a new OperatingSystemCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOperatingSystemCreateRequestWithDefaults

`func NewOperatingSystemCreateRequestWithDefaults() *OperatingSystemCreateRequest`

NewOperatingSystemCreateRequestWithDefaults instantiates a new OperatingSystemCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *OperatingSystemCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OperatingSystemCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OperatingSystemCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *OperatingSystemCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *OperatingSystemCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *OperatingSystemCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *OperatingSystemCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *OperatingSystemCreateRequest) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *OperatingSystemCreateRequest) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *OperatingSystemCreateRequest) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *OperatingSystemCreateRequest) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### SetInfrastructureProviderIdNil

`func (o *OperatingSystemCreateRequest) SetInfrastructureProviderIdNil(b bool)`

 SetInfrastructureProviderIdNil sets the value for InfrastructureProviderId to be an explicit nil

### UnsetInfrastructureProviderId
`func (o *OperatingSystemCreateRequest) UnsetInfrastructureProviderId()`

UnsetInfrastructureProviderId ensures that no value is present for InfrastructureProviderId, not even an explicit nil
### GetTenantId

`func (o *OperatingSystemCreateRequest) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *OperatingSystemCreateRequest) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *OperatingSystemCreateRequest) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *OperatingSystemCreateRequest) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### SetTenantIdNil

`func (o *OperatingSystemCreateRequest) SetTenantIdNil(b bool)`

 SetTenantIdNil sets the value for TenantId to be an explicit nil

### UnsetTenantId
`func (o *OperatingSystemCreateRequest) UnsetTenantId()`

UnsetTenantId ensures that no value is present for TenantId, not even an explicit nil
### GetSiteIds

`func (o *OperatingSystemCreateRequest) GetSiteIds() []string`

GetSiteIds returns the SiteIds field if non-nil, zero value otherwise.

### GetSiteIdsOk

`func (o *OperatingSystemCreateRequest) GetSiteIdsOk() (*[]string, bool)`

GetSiteIdsOk returns a tuple with the SiteIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteIds

`func (o *OperatingSystemCreateRequest) SetSiteIds(v []string)`

SetSiteIds sets SiteIds field to given value.

### HasSiteIds

`func (o *OperatingSystemCreateRequest) HasSiteIds() bool`

HasSiteIds returns a boolean if a field has been set.

### GetIpxeScript

`func (o *OperatingSystemCreateRequest) GetIpxeScript() string`

GetIpxeScript returns the IpxeScript field if non-nil, zero value otherwise.

### GetIpxeScriptOk

`func (o *OperatingSystemCreateRequest) GetIpxeScriptOk() (*string, bool)`

GetIpxeScriptOk returns a tuple with the IpxeScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpxeScript

`func (o *OperatingSystemCreateRequest) SetIpxeScript(v string)`

SetIpxeScript sets IpxeScript field to given value.

### HasIpxeScript

`func (o *OperatingSystemCreateRequest) HasIpxeScript() bool`

HasIpxeScript returns a boolean if a field has been set.

### SetIpxeScriptNil

`func (o *OperatingSystemCreateRequest) SetIpxeScriptNil(b bool)`

 SetIpxeScriptNil sets the value for IpxeScript to be an explicit nil

### UnsetIpxeScript
`func (o *OperatingSystemCreateRequest) UnsetIpxeScript()`

UnsetIpxeScript ensures that no value is present for IpxeScript, not even an explicit nil
### GetImageUrl

`func (o *OperatingSystemCreateRequest) GetImageUrl() string`

GetImageUrl returns the ImageUrl field if non-nil, zero value otherwise.

### GetImageUrlOk

`func (o *OperatingSystemCreateRequest) GetImageUrlOk() (*string, bool)`

GetImageUrlOk returns a tuple with the ImageUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageUrl

`func (o *OperatingSystemCreateRequest) SetImageUrl(v string)`

SetImageUrl sets ImageUrl field to given value.

### HasImageUrl

`func (o *OperatingSystemCreateRequest) HasImageUrl() bool`

HasImageUrl returns a boolean if a field has been set.

### SetImageUrlNil

`func (o *OperatingSystemCreateRequest) SetImageUrlNil(b bool)`

 SetImageUrlNil sets the value for ImageUrl to be an explicit nil

### UnsetImageUrl
`func (o *OperatingSystemCreateRequest) UnsetImageUrl()`

UnsetImageUrl ensures that no value is present for ImageUrl, not even an explicit nil
### GetImageSha

`func (o *OperatingSystemCreateRequest) GetImageSha() string`

GetImageSha returns the ImageSha field if non-nil, zero value otherwise.

### GetImageShaOk

`func (o *OperatingSystemCreateRequest) GetImageShaOk() (*string, bool)`

GetImageShaOk returns a tuple with the ImageSha field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageSha

`func (o *OperatingSystemCreateRequest) SetImageSha(v string)`

SetImageSha sets ImageSha field to given value.

### HasImageSha

`func (o *OperatingSystemCreateRequest) HasImageSha() bool`

HasImageSha returns a boolean if a field has been set.

### SetImageShaNil

`func (o *OperatingSystemCreateRequest) SetImageShaNil(b bool)`

 SetImageShaNil sets the value for ImageSha to be an explicit nil

### UnsetImageSha
`func (o *OperatingSystemCreateRequest) UnsetImageSha()`

UnsetImageSha ensures that no value is present for ImageSha, not even an explicit nil
### GetImageAuthType

`func (o *OperatingSystemCreateRequest) GetImageAuthType() string`

GetImageAuthType returns the ImageAuthType field if non-nil, zero value otherwise.

### GetImageAuthTypeOk

`func (o *OperatingSystemCreateRequest) GetImageAuthTypeOk() (*string, bool)`

GetImageAuthTypeOk returns a tuple with the ImageAuthType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageAuthType

`func (o *OperatingSystemCreateRequest) SetImageAuthType(v string)`

SetImageAuthType sets ImageAuthType field to given value.

### HasImageAuthType

`func (o *OperatingSystemCreateRequest) HasImageAuthType() bool`

HasImageAuthType returns a boolean if a field has been set.

### SetImageAuthTypeNil

`func (o *OperatingSystemCreateRequest) SetImageAuthTypeNil(b bool)`

 SetImageAuthTypeNil sets the value for ImageAuthType to be an explicit nil

### UnsetImageAuthType
`func (o *OperatingSystemCreateRequest) UnsetImageAuthType()`

UnsetImageAuthType ensures that no value is present for ImageAuthType, not even an explicit nil
### GetImageAuthToken

`func (o *OperatingSystemCreateRequest) GetImageAuthToken() string`

GetImageAuthToken returns the ImageAuthToken field if non-nil, zero value otherwise.

### GetImageAuthTokenOk

`func (o *OperatingSystemCreateRequest) GetImageAuthTokenOk() (*string, bool)`

GetImageAuthTokenOk returns a tuple with the ImageAuthToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageAuthToken

`func (o *OperatingSystemCreateRequest) SetImageAuthToken(v string)`

SetImageAuthToken sets ImageAuthToken field to given value.

### HasImageAuthToken

`func (o *OperatingSystemCreateRequest) HasImageAuthToken() bool`

HasImageAuthToken returns a boolean if a field has been set.

### SetImageAuthTokenNil

`func (o *OperatingSystemCreateRequest) SetImageAuthTokenNil(b bool)`

 SetImageAuthTokenNil sets the value for ImageAuthToken to be an explicit nil

### UnsetImageAuthToken
`func (o *OperatingSystemCreateRequest) UnsetImageAuthToken()`

UnsetImageAuthToken ensures that no value is present for ImageAuthToken, not even an explicit nil
### GetImageDisk

`func (o *OperatingSystemCreateRequest) GetImageDisk() string`

GetImageDisk returns the ImageDisk field if non-nil, zero value otherwise.

### GetImageDiskOk

`func (o *OperatingSystemCreateRequest) GetImageDiskOk() (*string, bool)`

GetImageDiskOk returns a tuple with the ImageDisk field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageDisk

`func (o *OperatingSystemCreateRequest) SetImageDisk(v string)`

SetImageDisk sets ImageDisk field to given value.

### HasImageDisk

`func (o *OperatingSystemCreateRequest) HasImageDisk() bool`

HasImageDisk returns a boolean if a field has been set.

### SetImageDiskNil

`func (o *OperatingSystemCreateRequest) SetImageDiskNil(b bool)`

 SetImageDiskNil sets the value for ImageDisk to be an explicit nil

### UnsetImageDisk
`func (o *OperatingSystemCreateRequest) UnsetImageDisk()`

UnsetImageDisk ensures that no value is present for ImageDisk, not even an explicit nil
### GetRootFsId

`func (o *OperatingSystemCreateRequest) GetRootFsId() string`

GetRootFsId returns the RootFsId field if non-nil, zero value otherwise.

### GetRootFsIdOk

`func (o *OperatingSystemCreateRequest) GetRootFsIdOk() (*string, bool)`

GetRootFsIdOk returns a tuple with the RootFsId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRootFsId

`func (o *OperatingSystemCreateRequest) SetRootFsId(v string)`

SetRootFsId sets RootFsId field to given value.

### HasRootFsId

`func (o *OperatingSystemCreateRequest) HasRootFsId() bool`

HasRootFsId returns a boolean if a field has been set.

### SetRootFsIdNil

`func (o *OperatingSystemCreateRequest) SetRootFsIdNil(b bool)`

 SetRootFsIdNil sets the value for RootFsId to be an explicit nil

### UnsetRootFsId
`func (o *OperatingSystemCreateRequest) UnsetRootFsId()`

UnsetRootFsId ensures that no value is present for RootFsId, not even an explicit nil
### GetRootFsLabel

`func (o *OperatingSystemCreateRequest) GetRootFsLabel() string`

GetRootFsLabel returns the RootFsLabel field if non-nil, zero value otherwise.

### GetRootFsLabelOk

`func (o *OperatingSystemCreateRequest) GetRootFsLabelOk() (*string, bool)`

GetRootFsLabelOk returns a tuple with the RootFsLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRootFsLabel

`func (o *OperatingSystemCreateRequest) SetRootFsLabel(v string)`

SetRootFsLabel sets RootFsLabel field to given value.

### HasRootFsLabel

`func (o *OperatingSystemCreateRequest) HasRootFsLabel() bool`

HasRootFsLabel returns a boolean if a field has been set.

### SetRootFsLabelNil

`func (o *OperatingSystemCreateRequest) SetRootFsLabelNil(b bool)`

 SetRootFsLabelNil sets the value for RootFsLabel to be an explicit nil

### UnsetRootFsLabel
`func (o *OperatingSystemCreateRequest) UnsetRootFsLabel()`

UnsetRootFsLabel ensures that no value is present for RootFsLabel, not even an explicit nil
### GetPhoneHomeEnabled

`func (o *OperatingSystemCreateRequest) GetPhoneHomeEnabled() bool`

GetPhoneHomeEnabled returns the PhoneHomeEnabled field if non-nil, zero value otherwise.

### GetPhoneHomeEnabledOk

`func (o *OperatingSystemCreateRequest) GetPhoneHomeEnabledOk() (*bool, bool)`

GetPhoneHomeEnabledOk returns a tuple with the PhoneHomeEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneHomeEnabled

`func (o *OperatingSystemCreateRequest) SetPhoneHomeEnabled(v bool)`

SetPhoneHomeEnabled sets PhoneHomeEnabled field to given value.

### HasPhoneHomeEnabled

`func (o *OperatingSystemCreateRequest) HasPhoneHomeEnabled() bool`

HasPhoneHomeEnabled returns a boolean if a field has been set.

### SetPhoneHomeEnabledNil

`func (o *OperatingSystemCreateRequest) SetPhoneHomeEnabledNil(b bool)`

 SetPhoneHomeEnabledNil sets the value for PhoneHomeEnabled to be an explicit nil

### UnsetPhoneHomeEnabled
`func (o *OperatingSystemCreateRequest) UnsetPhoneHomeEnabled()`

UnsetPhoneHomeEnabled ensures that no value is present for PhoneHomeEnabled, not even an explicit nil
### GetUserData

`func (o *OperatingSystemCreateRequest) GetUserData() string`

GetUserData returns the UserData field if non-nil, zero value otherwise.

### GetUserDataOk

`func (o *OperatingSystemCreateRequest) GetUserDataOk() (*string, bool)`

GetUserDataOk returns a tuple with the UserData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserData

`func (o *OperatingSystemCreateRequest) SetUserData(v string)`

SetUserData sets UserData field to given value.

### HasUserData

`func (o *OperatingSystemCreateRequest) HasUserData() bool`

HasUserData returns a boolean if a field has been set.

### SetUserDataNil

`func (o *OperatingSystemCreateRequest) SetUserDataNil(b bool)`

 SetUserDataNil sets the value for UserData to be an explicit nil

### UnsetUserData
`func (o *OperatingSystemCreateRequest) UnsetUserData()`

UnsetUserData ensures that no value is present for UserData, not even an explicit nil
### GetIsCloudInit

`func (o *OperatingSystemCreateRequest) GetIsCloudInit() bool`

GetIsCloudInit returns the IsCloudInit field if non-nil, zero value otherwise.

### GetIsCloudInitOk

`func (o *OperatingSystemCreateRequest) GetIsCloudInitOk() (*bool, bool)`

GetIsCloudInitOk returns a tuple with the IsCloudInit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsCloudInit

`func (o *OperatingSystemCreateRequest) SetIsCloudInit(v bool)`

SetIsCloudInit sets IsCloudInit field to given value.

### HasIsCloudInit

`func (o *OperatingSystemCreateRequest) HasIsCloudInit() bool`

HasIsCloudInit returns a boolean if a field has been set.

### GetAllowOverride

`func (o *OperatingSystemCreateRequest) GetAllowOverride() bool`

GetAllowOverride returns the AllowOverride field if non-nil, zero value otherwise.

### GetAllowOverrideOk

`func (o *OperatingSystemCreateRequest) GetAllowOverrideOk() (*bool, bool)`

GetAllowOverrideOk returns a tuple with the AllowOverride field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowOverride

`func (o *OperatingSystemCreateRequest) SetAllowOverride(v bool)`

SetAllowOverride sets AllowOverride field to given value.

### HasAllowOverride

`func (o *OperatingSystemCreateRequest) HasAllowOverride() bool`

HasAllowOverride returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


