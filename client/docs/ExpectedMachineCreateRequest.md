# ExpectedMachineCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SiteId** | **string** | ID of the site the Expected Machine belongs to | 
**BmcMacAddress** | **string** | MAC address of the Expected Machine&#39;s BMC (Baseboard Management Controller) | 
**BmcUsername** | Pointer to **NullableString** | Username for accessing the Expected Machine&#39;s BMC | [optional] 
**BmcPassword** | Pointer to **NullableString** | Password for accessing the Expected Machine&#39;s BMC | [optional] 
**ChassisSerialNumber** | **string** | Serial number of the Expected Machine&#39;s chassis | 
**FallbackDPUSerialNumbers** | Pointer to **[]string** | Serial numbers of the Expected Machine&#39;s fallback DPUs (Data Processing Units) | [optional] 
**SkuId** | Pointer to **NullableString** | Optional ID of the SKU to associate with this Expected Machine | [optional] 
**Labels** | Pointer to **map[string]string** | User-defined key-value pairs for organizing and categorizing Expected Machines | [optional] 

## Methods

### NewExpectedMachineCreateRequest

`func NewExpectedMachineCreateRequest(siteId string, bmcMacAddress string, chassisSerialNumber string, ) *ExpectedMachineCreateRequest`

NewExpectedMachineCreateRequest instantiates a new ExpectedMachineCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExpectedMachineCreateRequestWithDefaults

`func NewExpectedMachineCreateRequestWithDefaults() *ExpectedMachineCreateRequest`

NewExpectedMachineCreateRequestWithDefaults instantiates a new ExpectedMachineCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSiteId

`func (o *ExpectedMachineCreateRequest) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *ExpectedMachineCreateRequest) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *ExpectedMachineCreateRequest) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.


### GetBmcMacAddress

`func (o *ExpectedMachineCreateRequest) GetBmcMacAddress() string`

GetBmcMacAddress returns the BmcMacAddress field if non-nil, zero value otherwise.

### GetBmcMacAddressOk

`func (o *ExpectedMachineCreateRequest) GetBmcMacAddressOk() (*string, bool)`

GetBmcMacAddressOk returns a tuple with the BmcMacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcMacAddress

`func (o *ExpectedMachineCreateRequest) SetBmcMacAddress(v string)`

SetBmcMacAddress sets BmcMacAddress field to given value.


### GetBmcUsername

`func (o *ExpectedMachineCreateRequest) GetBmcUsername() string`

GetBmcUsername returns the BmcUsername field if non-nil, zero value otherwise.

### GetBmcUsernameOk

`func (o *ExpectedMachineCreateRequest) GetBmcUsernameOk() (*string, bool)`

GetBmcUsernameOk returns a tuple with the BmcUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcUsername

`func (o *ExpectedMachineCreateRequest) SetBmcUsername(v string)`

SetBmcUsername sets BmcUsername field to given value.

### HasBmcUsername

`func (o *ExpectedMachineCreateRequest) HasBmcUsername() bool`

HasBmcUsername returns a boolean if a field has been set.

### SetBmcUsernameNil

`func (o *ExpectedMachineCreateRequest) SetBmcUsernameNil(b bool)`

 SetBmcUsernameNil sets the value for BmcUsername to be an explicit nil

### UnsetBmcUsername
`func (o *ExpectedMachineCreateRequest) UnsetBmcUsername()`

UnsetBmcUsername ensures that no value is present for BmcUsername, not even an explicit nil
### GetBmcPassword

`func (o *ExpectedMachineCreateRequest) GetBmcPassword() string`

GetBmcPassword returns the BmcPassword field if non-nil, zero value otherwise.

### GetBmcPasswordOk

`func (o *ExpectedMachineCreateRequest) GetBmcPasswordOk() (*string, bool)`

GetBmcPasswordOk returns a tuple with the BmcPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcPassword

`func (o *ExpectedMachineCreateRequest) SetBmcPassword(v string)`

SetBmcPassword sets BmcPassword field to given value.

### HasBmcPassword

`func (o *ExpectedMachineCreateRequest) HasBmcPassword() bool`

HasBmcPassword returns a boolean if a field has been set.

### SetBmcPasswordNil

`func (o *ExpectedMachineCreateRequest) SetBmcPasswordNil(b bool)`

 SetBmcPasswordNil sets the value for BmcPassword to be an explicit nil

### UnsetBmcPassword
`func (o *ExpectedMachineCreateRequest) UnsetBmcPassword()`

UnsetBmcPassword ensures that no value is present for BmcPassword, not even an explicit nil
### GetChassisSerialNumber

`func (o *ExpectedMachineCreateRequest) GetChassisSerialNumber() string`

GetChassisSerialNumber returns the ChassisSerialNumber field if non-nil, zero value otherwise.

### GetChassisSerialNumberOk

`func (o *ExpectedMachineCreateRequest) GetChassisSerialNumberOk() (*string, bool)`

GetChassisSerialNumberOk returns a tuple with the ChassisSerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChassisSerialNumber

`func (o *ExpectedMachineCreateRequest) SetChassisSerialNumber(v string)`

SetChassisSerialNumber sets ChassisSerialNumber field to given value.


### GetFallbackDPUSerialNumbers

`func (o *ExpectedMachineCreateRequest) GetFallbackDPUSerialNumbers() []string`

GetFallbackDPUSerialNumbers returns the FallbackDPUSerialNumbers field if non-nil, zero value otherwise.

### GetFallbackDPUSerialNumbersOk

`func (o *ExpectedMachineCreateRequest) GetFallbackDPUSerialNumbersOk() (*[]string, bool)`

GetFallbackDPUSerialNumbersOk returns a tuple with the FallbackDPUSerialNumbers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFallbackDPUSerialNumbers

`func (o *ExpectedMachineCreateRequest) SetFallbackDPUSerialNumbers(v []string)`

SetFallbackDPUSerialNumbers sets FallbackDPUSerialNumbers field to given value.

### HasFallbackDPUSerialNumbers

`func (o *ExpectedMachineCreateRequest) HasFallbackDPUSerialNumbers() bool`

HasFallbackDPUSerialNumbers returns a boolean if a field has been set.

### SetFallbackDPUSerialNumbersNil

`func (o *ExpectedMachineCreateRequest) SetFallbackDPUSerialNumbersNil(b bool)`

 SetFallbackDPUSerialNumbersNil sets the value for FallbackDPUSerialNumbers to be an explicit nil

### UnsetFallbackDPUSerialNumbers
`func (o *ExpectedMachineCreateRequest) UnsetFallbackDPUSerialNumbers()`

UnsetFallbackDPUSerialNumbers ensures that no value is present for FallbackDPUSerialNumbers, not even an explicit nil
### GetSkuId

`func (o *ExpectedMachineCreateRequest) GetSkuId() string`

GetSkuId returns the SkuId field if non-nil, zero value otherwise.

### GetSkuIdOk

`func (o *ExpectedMachineCreateRequest) GetSkuIdOk() (*string, bool)`

GetSkuIdOk returns a tuple with the SkuId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkuId

`func (o *ExpectedMachineCreateRequest) SetSkuId(v string)`

SetSkuId sets SkuId field to given value.

### HasSkuId

`func (o *ExpectedMachineCreateRequest) HasSkuId() bool`

HasSkuId returns a boolean if a field has been set.

### SetSkuIdNil

`func (o *ExpectedMachineCreateRequest) SetSkuIdNil(b bool)`

 SetSkuIdNil sets the value for SkuId to be an explicit nil

### UnsetSkuId
`func (o *ExpectedMachineCreateRequest) UnsetSkuId()`

UnsetSkuId ensures that no value is present for SkuId, not even an explicit nil
### GetLabels

`func (o *ExpectedMachineCreateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *ExpectedMachineCreateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *ExpectedMachineCreateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *ExpectedMachineCreateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


