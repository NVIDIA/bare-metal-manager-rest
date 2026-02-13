# SshKeyGroupSiteAssociation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Site** | Pointer to [**SiteSummary**](SiteSummary.md) |  | [optional] 
**Status** | Pointer to [**SshKeyGroupSiteAssociationStatus**](SshKeyGroupSiteAssociationStatus.md) |  | [optional] 
**Version** | Pointer to **string** | Version of the Key Group on Site | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewSshKeyGroupSiteAssociation

`func NewSshKeyGroupSiteAssociation() *SshKeyGroupSiteAssociation`

NewSshKeyGroupSiteAssociation instantiates a new SshKeyGroupSiteAssociation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSshKeyGroupSiteAssociationWithDefaults

`func NewSshKeyGroupSiteAssociationWithDefaults() *SshKeyGroupSiteAssociation`

NewSshKeyGroupSiteAssociationWithDefaults instantiates a new SshKeyGroupSiteAssociation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSite

`func (o *SshKeyGroupSiteAssociation) GetSite() SiteSummary`

GetSite returns the Site field if non-nil, zero value otherwise.

### GetSiteOk

`func (o *SshKeyGroupSiteAssociation) GetSiteOk() (*SiteSummary, bool)`

GetSiteOk returns a tuple with the Site field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSite

`func (o *SshKeyGroupSiteAssociation) SetSite(v SiteSummary)`

SetSite sets Site field to given value.

### HasSite

`func (o *SshKeyGroupSiteAssociation) HasSite() bool`

HasSite returns a boolean if a field has been set.

### GetStatus

`func (o *SshKeyGroupSiteAssociation) GetStatus() SshKeyGroupSiteAssociationStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *SshKeyGroupSiteAssociation) GetStatusOk() (*SshKeyGroupSiteAssociationStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *SshKeyGroupSiteAssociation) SetStatus(v SshKeyGroupSiteAssociationStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *SshKeyGroupSiteAssociation) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetVersion

`func (o *SshKeyGroupSiteAssociation) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *SshKeyGroupSiteAssociation) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *SshKeyGroupSiteAssociation) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *SshKeyGroupSiteAssociation) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetCreated

`func (o *SshKeyGroupSiteAssociation) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *SshKeyGroupSiteAssociation) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *SshKeyGroupSiteAssociation) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *SshKeyGroupSiteAssociation) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *SshKeyGroupSiteAssociation) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *SshKeyGroupSiteAssociation) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *SshKeyGroupSiteAssociation) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *SshKeyGroupSiteAssociation) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


