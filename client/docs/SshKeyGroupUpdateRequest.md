# SshKeyGroupUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**SiteIds** | Pointer to **[]string** | When specified, replaces existing Site associations | [optional] 
**SshKeyIds** | Pointer to **[]string** | When specified, replaces existing SSH Key associations | [optional] 
**Version** | **string** | Version of the SSH Key Group being modified must be provided | 

## Methods

### NewSshKeyGroupUpdateRequest

`func NewSshKeyGroupUpdateRequest(version string, ) *SshKeyGroupUpdateRequest`

NewSshKeyGroupUpdateRequest instantiates a new SshKeyGroupUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSshKeyGroupUpdateRequestWithDefaults

`func NewSshKeyGroupUpdateRequestWithDefaults() *SshKeyGroupUpdateRequest`

NewSshKeyGroupUpdateRequestWithDefaults instantiates a new SshKeyGroupUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SshKeyGroupUpdateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SshKeyGroupUpdateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SshKeyGroupUpdateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SshKeyGroupUpdateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### SetNameNil

`func (o *SshKeyGroupUpdateRequest) SetNameNil(b bool)`

 SetNameNil sets the value for Name to be an explicit nil

### UnsetName
`func (o *SshKeyGroupUpdateRequest) UnsetName()`

UnsetName ensures that no value is present for Name, not even an explicit nil
### GetDescription

`func (o *SshKeyGroupUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *SshKeyGroupUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *SshKeyGroupUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *SshKeyGroupUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *SshKeyGroupUpdateRequest) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *SshKeyGroupUpdateRequest) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetSiteIds

`func (o *SshKeyGroupUpdateRequest) GetSiteIds() []string`

GetSiteIds returns the SiteIds field if non-nil, zero value otherwise.

### GetSiteIdsOk

`func (o *SshKeyGroupUpdateRequest) GetSiteIdsOk() (*[]string, bool)`

GetSiteIdsOk returns a tuple with the SiteIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteIds

`func (o *SshKeyGroupUpdateRequest) SetSiteIds(v []string)`

SetSiteIds sets SiteIds field to given value.

### HasSiteIds

`func (o *SshKeyGroupUpdateRequest) HasSiteIds() bool`

HasSiteIds returns a boolean if a field has been set.

### GetSshKeyIds

`func (o *SshKeyGroupUpdateRequest) GetSshKeyIds() []string`

GetSshKeyIds returns the SshKeyIds field if non-nil, zero value otherwise.

### GetSshKeyIdsOk

`func (o *SshKeyGroupUpdateRequest) GetSshKeyIdsOk() (*[]string, bool)`

GetSshKeyIdsOk returns a tuple with the SshKeyIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyIds

`func (o *SshKeyGroupUpdateRequest) SetSshKeyIds(v []string)`

SetSshKeyIds sets SshKeyIds field to given value.

### HasSshKeyIds

`func (o *SshKeyGroupUpdateRequest) HasSshKeyIds() bool`

HasSshKeyIds returns a boolean if a field has been set.

### GetVersion

`func (o *SshKeyGroupUpdateRequest) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *SshKeyGroupUpdateRequest) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *SshKeyGroupUpdateRequest) SetVersion(v string)`

SetVersion sets Version field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


