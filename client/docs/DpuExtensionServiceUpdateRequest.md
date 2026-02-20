# DpuExtensionServiceUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name for the DPU Extension Service. Must be unique for a given Tenant | [optional] 
**Description** | Pointer to **string** | Optional description for the DPU Extension Service | [optional] 
**Data** | Pointer to **string** | Deployment spec for the DPU Extension Service | [optional] 
**Credentials** | Pointer to [**DpuExtensionServiceCredentials**](DpuExtensionServiceCredentials.md) | Credentials to download resources specified in DPU Extension Service data | [optional] 

## Methods

### NewDpuExtensionServiceUpdateRequest

`func NewDpuExtensionServiceUpdateRequest() *DpuExtensionServiceUpdateRequest`

NewDpuExtensionServiceUpdateRequest instantiates a new DpuExtensionServiceUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDpuExtensionServiceUpdateRequestWithDefaults

`func NewDpuExtensionServiceUpdateRequestWithDefaults() *DpuExtensionServiceUpdateRequest`

NewDpuExtensionServiceUpdateRequestWithDefaults instantiates a new DpuExtensionServiceUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *DpuExtensionServiceUpdateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DpuExtensionServiceUpdateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DpuExtensionServiceUpdateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DpuExtensionServiceUpdateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *DpuExtensionServiceUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *DpuExtensionServiceUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *DpuExtensionServiceUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *DpuExtensionServiceUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetData

`func (o *DpuExtensionServiceUpdateRequest) GetData() string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DpuExtensionServiceUpdateRequest) GetDataOk() (*string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DpuExtensionServiceUpdateRequest) SetData(v string)`

SetData sets Data field to given value.

### HasData

`func (o *DpuExtensionServiceUpdateRequest) HasData() bool`

HasData returns a boolean if a field has been set.

### GetCredentials

`func (o *DpuExtensionServiceUpdateRequest) GetCredentials() DpuExtensionServiceCredentials`

GetCredentials returns the Credentials field if non-nil, zero value otherwise.

### GetCredentialsOk

`func (o *DpuExtensionServiceUpdateRequest) GetCredentialsOk() (*DpuExtensionServiceCredentials, bool)`

GetCredentialsOk returns a tuple with the Credentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentials

`func (o *DpuExtensionServiceUpdateRequest) SetCredentials(v DpuExtensionServiceCredentials)`

SetCredentials sets Credentials field to given value.

### HasCredentials

`func (o *DpuExtensionServiceUpdateRequest) HasCredentials() bool`

HasCredentials returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


