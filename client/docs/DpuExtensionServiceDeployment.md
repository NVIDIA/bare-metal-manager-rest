# DpuExtensionServiceDeployment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier for the DPU Extension Service Deployment | [optional] [readonly] 
**DpuExtensionService** | Pointer to [**DpuExtensionServiceSummary**](DpuExtensionServiceSummary.md) | Summary of the DPU Extension Service. Deployed version may be different | [optional] 
**Version** | Pointer to **string** | Deployed version of the DPU Extension Service | [optional] 
**Status** | Pointer to [**DpuExtensionServiceDeploymentStatus**](DpuExtensionServiceDeploymentStatus.md) | Status of the DPU Extension Service Deployment | [optional] 
**Created** | Pointer to **time.Time** | Date/time when this version of the DPU Extension Service Deployment was created | [optional] 
**Updated** | Pointer to **time.Time** | Date/time when this version of the DPU Extension Service Deployment was updated | [optional] 

## Methods

### NewDpuExtensionServiceDeployment

`func NewDpuExtensionServiceDeployment() *DpuExtensionServiceDeployment`

NewDpuExtensionServiceDeployment instantiates a new DpuExtensionServiceDeployment object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDpuExtensionServiceDeploymentWithDefaults

`func NewDpuExtensionServiceDeploymentWithDefaults() *DpuExtensionServiceDeployment`

NewDpuExtensionServiceDeploymentWithDefaults instantiates a new DpuExtensionServiceDeployment object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *DpuExtensionServiceDeployment) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DpuExtensionServiceDeployment) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DpuExtensionServiceDeployment) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *DpuExtensionServiceDeployment) HasId() bool`

HasId returns a boolean if a field has been set.

### GetDpuExtensionService

`func (o *DpuExtensionServiceDeployment) GetDpuExtensionService() DpuExtensionServiceSummary`

GetDpuExtensionService returns the DpuExtensionService field if non-nil, zero value otherwise.

### GetDpuExtensionServiceOk

`func (o *DpuExtensionServiceDeployment) GetDpuExtensionServiceOk() (*DpuExtensionServiceSummary, bool)`

GetDpuExtensionServiceOk returns a tuple with the DpuExtensionService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDpuExtensionService

`func (o *DpuExtensionServiceDeployment) SetDpuExtensionService(v DpuExtensionServiceSummary)`

SetDpuExtensionService sets DpuExtensionService field to given value.

### HasDpuExtensionService

`func (o *DpuExtensionServiceDeployment) HasDpuExtensionService() bool`

HasDpuExtensionService returns a boolean if a field has been set.

### GetVersion

`func (o *DpuExtensionServiceDeployment) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *DpuExtensionServiceDeployment) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *DpuExtensionServiceDeployment) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *DpuExtensionServiceDeployment) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetStatus

`func (o *DpuExtensionServiceDeployment) GetStatus() DpuExtensionServiceDeploymentStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DpuExtensionServiceDeployment) GetStatusOk() (*DpuExtensionServiceDeploymentStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DpuExtensionServiceDeployment) SetStatus(v DpuExtensionServiceDeploymentStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *DpuExtensionServiceDeployment) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetCreated

`func (o *DpuExtensionServiceDeployment) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *DpuExtensionServiceDeployment) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *DpuExtensionServiceDeployment) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *DpuExtensionServiceDeployment) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *DpuExtensionServiceDeployment) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *DpuExtensionServiceDeployment) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *DpuExtensionServiceDeployment) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *DpuExtensionServiceDeployment) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


