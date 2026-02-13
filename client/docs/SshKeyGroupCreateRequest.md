# SshKeyGroupCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**SiteIds** | Pointer to **[]string** |  | [optional] 
**SshKeyIds** | Pointer to **[]string** |  | [optional] 

## Methods

### NewSshKeyGroupCreateRequest

`func NewSshKeyGroupCreateRequest(name string, ) *SshKeyGroupCreateRequest`

NewSshKeyGroupCreateRequest instantiates a new SshKeyGroupCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSshKeyGroupCreateRequestWithDefaults

`func NewSshKeyGroupCreateRequestWithDefaults() *SshKeyGroupCreateRequest`

NewSshKeyGroupCreateRequestWithDefaults instantiates a new SshKeyGroupCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SshKeyGroupCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SshKeyGroupCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SshKeyGroupCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *SshKeyGroupCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *SshKeyGroupCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *SshKeyGroupCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *SshKeyGroupCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteIds

`func (o *SshKeyGroupCreateRequest) GetSiteIds() []string`

GetSiteIds returns the SiteIds field if non-nil, zero value otherwise.

### GetSiteIdsOk

`func (o *SshKeyGroupCreateRequest) GetSiteIdsOk() (*[]string, bool)`

GetSiteIdsOk returns a tuple with the SiteIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteIds

`func (o *SshKeyGroupCreateRequest) SetSiteIds(v []string)`

SetSiteIds sets SiteIds field to given value.

### HasSiteIds

`func (o *SshKeyGroupCreateRequest) HasSiteIds() bool`

HasSiteIds returns a boolean if a field has been set.

### GetSshKeyIds

`func (o *SshKeyGroupCreateRequest) GetSshKeyIds() []string`

GetSshKeyIds returns the SshKeyIds field if non-nil, zero value otherwise.

### GetSshKeyIdsOk

`func (o *SshKeyGroupCreateRequest) GetSshKeyIdsOk() (*[]string, bool)`

GetSshKeyIdsOk returns a tuple with the SshKeyIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyIds

`func (o *SshKeyGroupCreateRequest) SetSshKeyIds(v []string)`

SetSshKeyIds sets SshKeyIds field to given value.

### HasSshKeyIds

`func (o *SshKeyGroupCreateRequest) HasSshKeyIds() bool`

HasSshKeyIds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


