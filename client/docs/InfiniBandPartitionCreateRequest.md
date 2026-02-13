# InfiniBandPartitionCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name of the Parition to create | 
**Description** | Pointer to **string** | Optional description of the Partition | [optional] 
**SiteId** | **string** | ID of the Site the Partition should belong to | 

## Methods

### NewInfiniBandPartitionCreateRequest

`func NewInfiniBandPartitionCreateRequest(name string, siteId string, ) *InfiniBandPartitionCreateRequest`

NewInfiniBandPartitionCreateRequest instantiates a new InfiniBandPartitionCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInfiniBandPartitionCreateRequestWithDefaults

`func NewInfiniBandPartitionCreateRequestWithDefaults() *InfiniBandPartitionCreateRequest`

NewInfiniBandPartitionCreateRequestWithDefaults instantiates a new InfiniBandPartitionCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *InfiniBandPartitionCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *InfiniBandPartitionCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *InfiniBandPartitionCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *InfiniBandPartitionCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *InfiniBandPartitionCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *InfiniBandPartitionCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *InfiniBandPartitionCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSiteId

`func (o *InfiniBandPartitionCreateRequest) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *InfiniBandPartitionCreateRequest) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *InfiniBandPartitionCreateRequest) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


