# NVLinkLogicalPartitionCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name of the NVLink Logical Partition to create | 
**Description** | Pointer to **string** | Optional description of the NVLink Logical Partition | [optional] 
**SiteId** | **string** | ID of the Site the NVLink Logical Partition should belong to | 

## Methods

### NewNVLinkLogicalPartitionCreateRequest

`func NewNVLinkLogicalPartitionCreateRequest(name string, siteId string, ) *NVLinkLogicalPartitionCreateRequest`

NewNVLinkLogicalPartitionCreateRequest instantiates a new NVLinkLogicalPartitionCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNVLinkLogicalPartitionCreateRequestWithDefaults

`func NewNVLinkLogicalPartitionCreateRequestWithDefaults() *NVLinkLogicalPartitionCreateRequest`

NewNVLinkLogicalPartitionCreateRequestWithDefaults instantiates a new NVLinkLogicalPartitionCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *NVLinkLogicalPartitionCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NVLinkLogicalPartitionCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NVLinkLogicalPartitionCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *NVLinkLogicalPartitionCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *NVLinkLogicalPartitionCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *NVLinkLogicalPartitionCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *NVLinkLogicalPartitionCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *NVLinkLogicalPartitionCreateRequest) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *NVLinkLogicalPartitionCreateRequest) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *NVLinkLogicalPartitionCreateRequest) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


