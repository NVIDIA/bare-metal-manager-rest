# DpuExtensionServiceVersionInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Version** | Pointer to **NullableString** | Current version of the DPU Extension Service | [optional] 
**Data** | Pointer to **string** | Deployment spec for the DPU Extension Service | [optional] 
**HasCredentials** | Pointer to **bool** | Indicates whether this version was created with credentials | [optional] 
**Created** | Pointer to **time.Time** | Date/time when this version of the DPU Extension Service was created | [optional] [readonly] 

## Methods

### NewDpuExtensionServiceVersionInfo

`func NewDpuExtensionServiceVersionInfo() *DpuExtensionServiceVersionInfo`

NewDpuExtensionServiceVersionInfo instantiates a new DpuExtensionServiceVersionInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDpuExtensionServiceVersionInfoWithDefaults

`func NewDpuExtensionServiceVersionInfoWithDefaults() *DpuExtensionServiceVersionInfo`

NewDpuExtensionServiceVersionInfoWithDefaults instantiates a new DpuExtensionServiceVersionInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersion

`func (o *DpuExtensionServiceVersionInfo) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *DpuExtensionServiceVersionInfo) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *DpuExtensionServiceVersionInfo) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *DpuExtensionServiceVersionInfo) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### SetVersionNil

`func (o *DpuExtensionServiceVersionInfo) SetVersionNil(b bool)`

 SetVersionNil sets the value for Version to be an explicit nil

### UnsetVersion
`func (o *DpuExtensionServiceVersionInfo) UnsetVersion()`

UnsetVersion ensures that no value is present for Version, not even an explicit nil
### GetData

`func (o *DpuExtensionServiceVersionInfo) GetData() string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DpuExtensionServiceVersionInfo) GetDataOk() (*string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DpuExtensionServiceVersionInfo) SetData(v string)`

SetData sets Data field to given value.

### HasData

`func (o *DpuExtensionServiceVersionInfo) HasData() bool`

HasData returns a boolean if a field has been set.

### GetHasCredentials

`func (o *DpuExtensionServiceVersionInfo) GetHasCredentials() bool`

GetHasCredentials returns the HasCredentials field if non-nil, zero value otherwise.

### GetHasCredentialsOk

`func (o *DpuExtensionServiceVersionInfo) GetHasCredentialsOk() (*bool, bool)`

GetHasCredentialsOk returns a tuple with the HasCredentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasCredentials

`func (o *DpuExtensionServiceVersionInfo) SetHasCredentials(v bool)`

SetHasCredentials sets HasCredentials field to given value.

### HasHasCredentials

`func (o *DpuExtensionServiceVersionInfo) HasHasCredentials() bool`

HasHasCredentials returns a boolean if a field has been set.

### GetCreated

`func (o *DpuExtensionServiceVersionInfo) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *DpuExtensionServiceVersionInfo) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *DpuExtensionServiceVersionInfo) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *DpuExtensionServiceVersionInfo) HasCreated() bool`

HasCreated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


