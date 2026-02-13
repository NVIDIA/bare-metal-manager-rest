# NVLinkLogicalPartition

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique ID of the NVLink Logical Partition | [optional] [readonly] 
**Name** | Pointer to **string** | Name of the NVLink Logical Partition | [optional] 
**Description** | Pointer to **string** | Optional description of the NVLink Logical Partition | [optional] 
**SiteId** | Pointer to **string** | ID of the Site the NVLink Logical Partition belongs to | [optional] 
**TenantId** | Pointer to **string** | ID of the Tenant the NVLink Logical Partition belongs to | [optional] 
**Status** | Pointer to [**NVLinkLogicalPartitionStatus**](NVLinkLogicalPartitionStatus.md) | Status of the NVLink Logical Partition | [optional] [readonly] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) | Details of status changes for the NVLink Logical Partition over time | [optional] 
**Created** | Pointer to **time.Time** | Date and time the NVLink Logical Partition was created | [optional] [readonly] 
**Updated** | Pointer to **time.Time** | Date and time the NVLink Logical Partition was last updated | [optional] [readonly] 

## Methods

### NewNVLinkLogicalPartition

`func NewNVLinkLogicalPartition() *NVLinkLogicalPartition`

NewNVLinkLogicalPartition instantiates a new NVLinkLogicalPartition object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNVLinkLogicalPartitionWithDefaults

`func NewNVLinkLogicalPartitionWithDefaults() *NVLinkLogicalPartition`

NewNVLinkLogicalPartitionWithDefaults instantiates a new NVLinkLogicalPartition object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *NVLinkLogicalPartition) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *NVLinkLogicalPartition) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *NVLinkLogicalPartition) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *NVLinkLogicalPartition) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *NVLinkLogicalPartition) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NVLinkLogicalPartition) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NVLinkLogicalPartition) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *NVLinkLogicalPartition) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *NVLinkLogicalPartition) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *NVLinkLogicalPartition) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *NVLinkLogicalPartition) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *NVLinkLogicalPartition) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *NVLinkLogicalPartition) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *NVLinkLogicalPartition) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *NVLinkLogicalPartition) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *NVLinkLogicalPartition) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetTenantId

`func (o *NVLinkLogicalPartition) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *NVLinkLogicalPartition) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *NVLinkLogicalPartition) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *NVLinkLogicalPartition) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetStatus

`func (o *NVLinkLogicalPartition) GetStatus() NVLinkLogicalPartitionStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *NVLinkLogicalPartition) GetStatusOk() (*NVLinkLogicalPartitionStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *NVLinkLogicalPartition) SetStatus(v NVLinkLogicalPartitionStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *NVLinkLogicalPartition) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *NVLinkLogicalPartition) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *NVLinkLogicalPartition) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *NVLinkLogicalPartition) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *NVLinkLogicalPartition) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *NVLinkLogicalPartition) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *NVLinkLogicalPartition) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *NVLinkLogicalPartition) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *NVLinkLogicalPartition) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *NVLinkLogicalPartition) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *NVLinkLogicalPartition) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *NVLinkLogicalPartition) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *NVLinkLogicalPartition) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


