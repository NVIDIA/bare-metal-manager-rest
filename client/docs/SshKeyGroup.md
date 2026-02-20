# SshKeyGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier for the SSH Key Group | [optional] [readonly] 
**Name** | Pointer to **string** | Name of the SSH Key Group | [optional] 
**Description** | Pointer to **NullableString** | Description for the SSH Key Group, optional | [optional] 
**Org** | Pointer to **string** | Organization this SSH Key Group belongs to | [optional] 
**TenantId** | Pointer to **string** | ID of the Tenane the SSH Key Group belongs to | [optional] 
**Version** | Pointer to **string** | Version of the SSH Key Group | [optional] 
**SshKeys** | Pointer to [**[]SshKey**](SshKey.md) | SSH Keys associated with this SSH Key Group | [optional] 
**SiteAssociations** | Pointer to [**[]SshKeyGroupSiteAssociation**](SshKeyGroupSiteAssociation.md) | Sites the SSH Key Group is synced to | [optional] 
**Status** | Pointer to [**SshKeyGroupStatus**](SshKeyGroupStatus.md) | Status of the SSH Key Group | [optional] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) | History of the SSH Key Group states | [optional] 
**Created** | Pointer to **time.Time** | Date/time when the SSH key was created | [optional] [readonly] 
**Updated** | Pointer to **time.Time** | Date/time when the SSH key was last updated | [optional] [readonly] 

## Methods

### NewSshKeyGroup

`func NewSshKeyGroup() *SshKeyGroup`

NewSshKeyGroup instantiates a new SshKeyGroup object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSshKeyGroupWithDefaults

`func NewSshKeyGroupWithDefaults() *SshKeyGroup`

NewSshKeyGroupWithDefaults instantiates a new SshKeyGroup object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *SshKeyGroup) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SshKeyGroup) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SshKeyGroup) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *SshKeyGroup) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *SshKeyGroup) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SshKeyGroup) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SshKeyGroup) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SshKeyGroup) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *SshKeyGroup) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *SshKeyGroup) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *SshKeyGroup) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *SshKeyGroup) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *SshKeyGroup) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *SshKeyGroup) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetOrg

`func (o *SshKeyGroup) GetOrg() string`

GetOrg returns the Org field if non-nil, zero value otherwise.

### GetOrgOk

`func (o *SshKeyGroup) GetOrgOk() (*string, bool)`

GetOrgOk returns a tuple with the Org field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrg

`func (o *SshKeyGroup) SetOrg(v string)`

SetOrg sets Org field to given value.

### HasOrg

`func (o *SshKeyGroup) HasOrg() bool`

HasOrg returns a boolean if a field has been set.

### GetTenantId

`func (o *SshKeyGroup) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *SshKeyGroup) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *SshKeyGroup) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *SshKeyGroup) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetVersion

`func (o *SshKeyGroup) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *SshKeyGroup) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *SshKeyGroup) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *SshKeyGroup) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetSshKeys

`func (o *SshKeyGroup) GetSshKeys() []SshKey`

GetSshKeys returns the SshKeys field if non-nil, zero value otherwise.

### GetSshKeysOk

`func (o *SshKeyGroup) GetSshKeysOk() (*[]SshKey, bool)`

GetSshKeysOk returns a tuple with the SshKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeys

`func (o *SshKeyGroup) SetSshKeys(v []SshKey)`

SetSshKeys sets SshKeys field to given value.

### HasSshKeys

`func (o *SshKeyGroup) HasSshKeys() bool`

HasSshKeys returns a boolean if a field has been set.

### GetSiteAssociations

`func (o *SshKeyGroup) GetSiteAssociations() []SshKeyGroupSiteAssociation`

GetSiteAssociations returns the SiteAssociations field if non-nil, zero value otherwise.

### GetSiteAssociationsOk

`func (o *SshKeyGroup) GetSiteAssociationsOk() (*[]SshKeyGroupSiteAssociation, bool)`

GetSiteAssociationsOk returns a tuple with the SiteAssociations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteAssociations

`func (o *SshKeyGroup) SetSiteAssociations(v []SshKeyGroupSiteAssociation)`

SetSiteAssociations sets SiteAssociations field to given value.

### HasSiteAssociations

`func (o *SshKeyGroup) HasSiteAssociations() bool`

HasSiteAssociations returns a boolean if a field has been set.

### GetStatus

`func (o *SshKeyGroup) GetStatus() SshKeyGroupStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *SshKeyGroup) GetStatusOk() (*SshKeyGroupStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *SshKeyGroup) SetStatus(v SshKeyGroupStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *SshKeyGroup) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *SshKeyGroup) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *SshKeyGroup) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *SshKeyGroup) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *SshKeyGroup) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *SshKeyGroup) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *SshKeyGroup) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *SshKeyGroup) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *SshKeyGroup) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *SshKeyGroup) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *SshKeyGroup) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *SshKeyGroup) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *SshKeyGroup) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


