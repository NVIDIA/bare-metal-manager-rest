# OperatingSystemSiteAssociation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Site** | Pointer to [**SiteSummary**](SiteSummary.md) |  | [optional] 
**Status** | Pointer to [**SshKeyGroupStatus**](SshKeyGroupStatus.md) |  | [optional] 
**Version** | Pointer to **string** | Version of the Key Group on Site | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewOperatingSystemSiteAssociation

`func NewOperatingSystemSiteAssociation() *OperatingSystemSiteAssociation`

NewOperatingSystemSiteAssociation instantiates a new OperatingSystemSiteAssociation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOperatingSystemSiteAssociationWithDefaults

`func NewOperatingSystemSiteAssociationWithDefaults() *OperatingSystemSiteAssociation`

NewOperatingSystemSiteAssociationWithDefaults instantiates a new OperatingSystemSiteAssociation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSite

`func (o *OperatingSystemSiteAssociation) GetSite() SiteSummary`

GetSite returns the Site field if non-nil, zero value otherwise.

### GetSiteOk

`func (o *OperatingSystemSiteAssociation) GetSiteOk() (*SiteSummary, bool)`

GetSiteOk returns a tuple with the Site field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSite

`func (o *OperatingSystemSiteAssociation) SetSite(v SiteSummary)`

SetSite sets Site field to given value.

### HasSite

`func (o *OperatingSystemSiteAssociation) HasSite() bool`

HasSite returns a boolean if a field has been set.

### GetStatus

`func (o *OperatingSystemSiteAssociation) GetStatus() SshKeyGroupStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *OperatingSystemSiteAssociation) GetStatusOk() (*SshKeyGroupStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *OperatingSystemSiteAssociation) SetStatus(v SshKeyGroupStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *OperatingSystemSiteAssociation) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetVersion

`func (o *OperatingSystemSiteAssociation) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *OperatingSystemSiteAssociation) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *OperatingSystemSiteAssociation) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *OperatingSystemSiteAssociation) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetCreated

`func (o *OperatingSystemSiteAssociation) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *OperatingSystemSiteAssociation) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *OperatingSystemSiteAssociation) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *OperatingSystemSiteAssociation) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *OperatingSystemSiteAssociation) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *OperatingSystemSiteAssociation) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *OperatingSystemSiteAssociation) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *OperatingSystemSiteAssociation) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


