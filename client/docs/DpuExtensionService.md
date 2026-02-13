# DpuExtensionService

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier for the DPU Extension Service | [optional] [readonly] 
**Name** | Pointer to **string** | Name for the DPU Extension Service. Must be unique for a given Tenant | [optional] 
**Description** | Pointer to **NullableString** | Optional description for the DPU Extension Service | [optional] 
**ServiceType** | Pointer to **string** | Type of the DPU Extension Service | [optional] 
**SiteId** | Pointer to **string** | ID for the Site the DPU Extension Service belongs to | [optional] 
**TenantId** | Pointer to **string** | ID for the Tenant the DPU Extension Service belongs to | [optional] 
**Version** | Pointer to **NullableString** | Latest version of the DPU Extension Service | [optional] 
**VersionInfo** | Pointer to [**DpuExtensionServiceVersionInfo**](DpuExtensionServiceVersionInfo.md) | Details for the latest version of the DPU Extension Service | [optional] 
**ActiveVersions** | Pointer to **[]string** | Latest and past versions of this DPU Extension Service that have not been deleted and are available for deployment | [optional] 
**Status** | Pointer to [**DpuExtensionServiceStatus**](DpuExtensionServiceStatus.md) | Status of the DPU Extension Service | [optional] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) | History of the DPU Extension Service statuses | [optional] 
**Created** | Pointer to **time.Time** | Date/time when the DPU Extension Service was created | [optional] [readonly] 
**Updated** | Pointer to **time.Time** | Date/time when the DPU Extension Service was last updated | [optional] [readonly] 

## Methods

### NewDpuExtensionService

`func NewDpuExtensionService() *DpuExtensionService`

NewDpuExtensionService instantiates a new DpuExtensionService object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDpuExtensionServiceWithDefaults

`func NewDpuExtensionServiceWithDefaults() *DpuExtensionService`

NewDpuExtensionServiceWithDefaults instantiates a new DpuExtensionService object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *DpuExtensionService) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DpuExtensionService) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DpuExtensionService) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *DpuExtensionService) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *DpuExtensionService) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DpuExtensionService) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DpuExtensionService) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DpuExtensionService) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *DpuExtensionService) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *DpuExtensionService) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *DpuExtensionService) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *DpuExtensionService) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *DpuExtensionService) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *DpuExtensionService) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetServiceType

`func (o *DpuExtensionService) GetServiceType() string`

GetServiceType returns the ServiceType field if non-nil, zero value otherwise.

### GetServiceTypeOk

`func (o *DpuExtensionService) GetServiceTypeOk() (*string, bool)`

GetServiceTypeOk returns a tuple with the ServiceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceType

`func (o *DpuExtensionService) SetServiceType(v string)`

SetServiceType sets ServiceType field to given value.

### HasServiceType

`func (o *DpuExtensionService) HasServiceType() bool`

HasServiceType returns a boolean if a field has been set.

### GetSiteId

`func (o *DpuExtensionService) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *DpuExtensionService) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *DpuExtensionService) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *DpuExtensionService) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetTenantId

`func (o *DpuExtensionService) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *DpuExtensionService) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *DpuExtensionService) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *DpuExtensionService) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetVersion

`func (o *DpuExtensionService) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *DpuExtensionService) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *DpuExtensionService) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *DpuExtensionService) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### SetVersionNil

`func (o *DpuExtensionService) SetVersionNil(b bool)`

 SetVersionNil sets the value for Version to be an explicit nil

### UnsetVersion
`func (o *DpuExtensionService) UnsetVersion()`

UnsetVersion ensures that no value is present for Version, not even an explicit nil
### GetVersionInfo

`func (o *DpuExtensionService) GetVersionInfo() DpuExtensionServiceVersionInfo`

GetVersionInfo returns the VersionInfo field if non-nil, zero value otherwise.

### GetVersionInfoOk

`func (o *DpuExtensionService) GetVersionInfoOk() (*DpuExtensionServiceVersionInfo, bool)`

GetVersionInfoOk returns a tuple with the VersionInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionInfo

`func (o *DpuExtensionService) SetVersionInfo(v DpuExtensionServiceVersionInfo)`

SetVersionInfo sets VersionInfo field to given value.

### HasVersionInfo

`func (o *DpuExtensionService) HasVersionInfo() bool`

HasVersionInfo returns a boolean if a field has been set.

### GetActiveVersions

`func (o *DpuExtensionService) GetActiveVersions() []string`

GetActiveVersions returns the ActiveVersions field if non-nil, zero value otherwise.

### GetActiveVersionsOk

`func (o *DpuExtensionService) GetActiveVersionsOk() (*[]string, bool)`

GetActiveVersionsOk returns a tuple with the ActiveVersions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveVersions

`func (o *DpuExtensionService) SetActiveVersions(v []string)`

SetActiveVersions sets ActiveVersions field to given value.

### HasActiveVersions

`func (o *DpuExtensionService) HasActiveVersions() bool`

HasActiveVersions returns a boolean if a field has been set.

### GetStatus

`func (o *DpuExtensionService) GetStatus() DpuExtensionServiceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DpuExtensionService) GetStatusOk() (*DpuExtensionServiceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DpuExtensionService) SetStatus(v DpuExtensionServiceStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *DpuExtensionService) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *DpuExtensionService) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *DpuExtensionService) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *DpuExtensionService) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *DpuExtensionService) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *DpuExtensionService) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *DpuExtensionService) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *DpuExtensionService) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *DpuExtensionService) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *DpuExtensionService) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *DpuExtensionService) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *DpuExtensionService) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *DpuExtensionService) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


