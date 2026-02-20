# DpuExtensionServiceSummary

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier for the DPU Extension Service | [optional] [readonly] 
**Name** | Pointer to **string** | Name for the DPU Extension Service. Must be unique for a given Tenant | [optional] 
**ServiceType** | Pointer to **string** | Type of the DPU Extension Service | [optional] 
**LatestVersion** | Pointer to **NullableString** | Latest version of the DPU Extension Service | [optional] 
**Status** | Pointer to [**DpuExtensionServiceStatus**](DpuExtensionServiceStatus.md) | Status of the DPU Extension Service | [optional] 

## Methods

### NewDpuExtensionServiceSummary

`func NewDpuExtensionServiceSummary() *DpuExtensionServiceSummary`

NewDpuExtensionServiceSummary instantiates a new DpuExtensionServiceSummary object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDpuExtensionServiceSummaryWithDefaults

`func NewDpuExtensionServiceSummaryWithDefaults() *DpuExtensionServiceSummary`

NewDpuExtensionServiceSummaryWithDefaults instantiates a new DpuExtensionServiceSummary object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *DpuExtensionServiceSummary) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DpuExtensionServiceSummary) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DpuExtensionServiceSummary) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *DpuExtensionServiceSummary) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *DpuExtensionServiceSummary) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DpuExtensionServiceSummary) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DpuExtensionServiceSummary) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DpuExtensionServiceSummary) HasName() bool`

HasName returns a boolean if a field has been set.

### GetServiceType

`func (o *DpuExtensionServiceSummary) GetServiceType() string`

GetServiceType returns the ServiceType field if non-nil, zero value otherwise.

### GetServiceTypeOk

`func (o *DpuExtensionServiceSummary) GetServiceTypeOk() (*string, bool)`

GetServiceTypeOk returns a tuple with the ServiceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceType

`func (o *DpuExtensionServiceSummary) SetServiceType(v string)`

SetServiceType sets ServiceType field to given value.

### HasServiceType

`func (o *DpuExtensionServiceSummary) HasServiceType() bool`

HasServiceType returns a boolean if a field has been set.

### GetLatestVersion

`func (o *DpuExtensionServiceSummary) GetLatestVersion() string`

GetLatestVersion returns the LatestVersion field if non-nil, zero value otherwise.

### GetLatestVersionOk

`func (o *DpuExtensionServiceSummary) GetLatestVersionOk() (*string, bool)`

GetLatestVersionOk returns a tuple with the LatestVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLatestVersion

`func (o *DpuExtensionServiceSummary) SetLatestVersion(v string)`

SetLatestVersion sets LatestVersion field to given value.

### HasLatestVersion

`func (o *DpuExtensionServiceSummary) HasLatestVersion() bool`

HasLatestVersion returns a boolean if a field has been set.

### SetLatestVersionNil

`func (o *DpuExtensionServiceSummary) SetLatestVersionNil(b bool)`

 SetLatestVersionNil sets the value for LatestVersion to be an explicit nil

### UnsetLatestVersion
`func (o *DpuExtensionServiceSummary) UnsetLatestVersion()`

UnsetLatestVersion ensures that no value is present for LatestVersion, not even an explicit nil
### GetStatus

`func (o *DpuExtensionServiceSummary) GetStatus() DpuExtensionServiceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DpuExtensionServiceSummary) GetStatusOk() (*DpuExtensionServiceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DpuExtensionServiceSummary) SetStatus(v DpuExtensionServiceStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *DpuExtensionServiceSummary) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


