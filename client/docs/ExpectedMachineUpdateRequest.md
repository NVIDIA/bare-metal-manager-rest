# ExpectedMachineUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **NullableString** | ID of the Expected Machine to update.  Optional for single updates (ignored if provided, URL path ID is used).  Required for batch update operations. | [optional] 
**BmcMacAddress** | Pointer to **NullableString** | MAC address of the Expected Machine&#39;s BMC (Baseboard Management Controller) | [optional] 
**BmcUsername** | Pointer to **NullableString** | Username for accessing the Expected Machine&#39;s BMC | [optional] 
**BmcPassword** | Pointer to **NullableString** | Password for accessing the Expected Machine&#39;s BMC | [optional] 
**ChassisSerialNumber** | Pointer to **NullableString** | Serial number of the Expected Machine&#39;s chassis | [optional] 
**FallbackDPUSerialNumbers** | Pointer to **[]string** | Serial numbers of the Expected Machine&#39;s fallback DPUs (Data Processing Units) | [optional] 
**SkuId** | Pointer to **NullableString** | Optional ID of the SKU to associate with this Expected Machine | [optional] 
**Labels** | Pointer to **map[string]string** | User-defined key-value pairs for organizing and categorizing Expected Machines | [optional] 

## Methods

### NewExpectedMachineUpdateRequest

`func NewExpectedMachineUpdateRequest() *ExpectedMachineUpdateRequest`

NewExpectedMachineUpdateRequest instantiates a new ExpectedMachineUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExpectedMachineUpdateRequestWithDefaults

`func NewExpectedMachineUpdateRequestWithDefaults() *ExpectedMachineUpdateRequest`

NewExpectedMachineUpdateRequestWithDefaults instantiates a new ExpectedMachineUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ExpectedMachineUpdateRequest) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ExpectedMachineUpdateRequest) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ExpectedMachineUpdateRequest) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ExpectedMachineUpdateRequest) HasId() bool`

HasId returns a boolean if a field has been set.

### SetIdNil

`func (o *ExpectedMachineUpdateRequest) SetIdNil(b bool)`

 SetIdNil sets the value for Id to be an explicit nil

### UnsetId
`func (o *ExpectedMachineUpdateRequest) UnsetId()`

UnsetId ensures that no value is present for Id, not even an explicit nil
### GetBmcMacAddress

`func (o *ExpectedMachineUpdateRequest) GetBmcMacAddress() string`

GetBmcMacAddress returns the BmcMacAddress field if non-nil, zero value otherwise.

### GetBmcMacAddressOk

`func (o *ExpectedMachineUpdateRequest) GetBmcMacAddressOk() (*string, bool)`

GetBmcMacAddressOk returns a tuple with the BmcMacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcMacAddress

`func (o *ExpectedMachineUpdateRequest) SetBmcMacAddress(v string)`

SetBmcMacAddress sets BmcMacAddress field to given value.

### HasBmcMacAddress

`func (o *ExpectedMachineUpdateRequest) HasBmcMacAddress() bool`

HasBmcMacAddress returns a boolean if a field has been set.

### SetBmcMacAddressNil

`func (o *ExpectedMachineUpdateRequest) SetBmcMacAddressNil(b bool)`

 SetBmcMacAddressNil sets the value for BmcMacAddress to be an explicit nil

### UnsetBmcMacAddress
`func (o *ExpectedMachineUpdateRequest) UnsetBmcMacAddress()`

UnsetBmcMacAddress ensures that no value is present for BmcMacAddress, not even an explicit nil
### GetBmcUsername

`func (o *ExpectedMachineUpdateRequest) GetBmcUsername() string`

GetBmcUsername returns the BmcUsername field if non-nil, zero value otherwise.

### GetBmcUsernameOk

`func (o *ExpectedMachineUpdateRequest) GetBmcUsernameOk() (*string, bool)`

GetBmcUsernameOk returns a tuple with the BmcUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcUsername

`func (o *ExpectedMachineUpdateRequest) SetBmcUsername(v string)`

SetBmcUsername sets BmcUsername field to given value.

### HasBmcUsername

`func (o *ExpectedMachineUpdateRequest) HasBmcUsername() bool`

HasBmcUsername returns a boolean if a field has been set.

### SetBmcUsernameNil

`func (o *ExpectedMachineUpdateRequest) SetBmcUsernameNil(b bool)`

 SetBmcUsernameNil sets the value for BmcUsername to be an explicit nil

### UnsetBmcUsername
`func (o *ExpectedMachineUpdateRequest) UnsetBmcUsername()`

UnsetBmcUsername ensures that no value is present for BmcUsername, not even an explicit nil
### GetBmcPassword

`func (o *ExpectedMachineUpdateRequest) GetBmcPassword() string`

GetBmcPassword returns the BmcPassword field if non-nil, zero value otherwise.

### GetBmcPasswordOk

`func (o *ExpectedMachineUpdateRequest) GetBmcPasswordOk() (*string, bool)`

GetBmcPasswordOk returns a tuple with the BmcPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcPassword

`func (o *ExpectedMachineUpdateRequest) SetBmcPassword(v string)`

SetBmcPassword sets BmcPassword field to given value.

### HasBmcPassword

`func (o *ExpectedMachineUpdateRequest) HasBmcPassword() bool`

HasBmcPassword returns a boolean if a field has been set.

### SetBmcPasswordNil

`func (o *ExpectedMachineUpdateRequest) SetBmcPasswordNil(b bool)`

 SetBmcPasswordNil sets the value for BmcPassword to be an explicit nil

### UnsetBmcPassword
`func (o *ExpectedMachineUpdateRequest) UnsetBmcPassword()`

UnsetBmcPassword ensures that no value is present for BmcPassword, not even an explicit nil
### GetChassisSerialNumber

`func (o *ExpectedMachineUpdateRequest) GetChassisSerialNumber() string`

GetChassisSerialNumber returns the ChassisSerialNumber field if non-nil, zero value otherwise.

### GetChassisSerialNumberOk

`func (o *ExpectedMachineUpdateRequest) GetChassisSerialNumberOk() (*string, bool)`

GetChassisSerialNumberOk returns a tuple with the ChassisSerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChassisSerialNumber

`func (o *ExpectedMachineUpdateRequest) SetChassisSerialNumber(v string)`

SetChassisSerialNumber sets ChassisSerialNumber field to given value.

### HasChassisSerialNumber

`func (o *ExpectedMachineUpdateRequest) HasChassisSerialNumber() bool`

HasChassisSerialNumber returns a boolean if a field has been set.

### SetChassisSerialNumberNil

`func (o *ExpectedMachineUpdateRequest) SetChassisSerialNumberNil(b bool)`

 SetChassisSerialNumberNil sets the value for ChassisSerialNumber to be an explicit nil

### UnsetChassisSerialNumber
`func (o *ExpectedMachineUpdateRequest) UnsetChassisSerialNumber()`

UnsetChassisSerialNumber ensures that no value is present for ChassisSerialNumber, not even an explicit nil
### GetFallbackDPUSerialNumbers

`func (o *ExpectedMachineUpdateRequest) GetFallbackDPUSerialNumbers() []string`

GetFallbackDPUSerialNumbers returns the FallbackDPUSerialNumbers field if non-nil, zero value otherwise.

### GetFallbackDPUSerialNumbersOk

`func (o *ExpectedMachineUpdateRequest) GetFallbackDPUSerialNumbersOk() (*[]string, bool)`

GetFallbackDPUSerialNumbersOk returns a tuple with the FallbackDPUSerialNumbers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFallbackDPUSerialNumbers

`func (o *ExpectedMachineUpdateRequest) SetFallbackDPUSerialNumbers(v []string)`

SetFallbackDPUSerialNumbers sets FallbackDPUSerialNumbers field to given value.

### HasFallbackDPUSerialNumbers

`func (o *ExpectedMachineUpdateRequest) HasFallbackDPUSerialNumbers() bool`

HasFallbackDPUSerialNumbers returns a boolean if a field has been set.

### SetFallbackDPUSerialNumbersNil

`func (o *ExpectedMachineUpdateRequest) SetFallbackDPUSerialNumbersNil(b bool)`

 SetFallbackDPUSerialNumbersNil sets the value for FallbackDPUSerialNumbers to be an explicit nil

### UnsetFallbackDPUSerialNumbers
`func (o *ExpectedMachineUpdateRequest) UnsetFallbackDPUSerialNumbers()`

UnsetFallbackDPUSerialNumbers ensures that no value is present for FallbackDPUSerialNumbers, not even an explicit nil
### GetSkuId

`func (o *ExpectedMachineUpdateRequest) GetSkuId() string`

GetSkuId returns the SkuId field if non-nil, zero value otherwise.

### GetSkuIdOk

`func (o *ExpectedMachineUpdateRequest) GetSkuIdOk() (*string, bool)`

GetSkuIdOk returns a tuple with the SkuId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkuId

`func (o *ExpectedMachineUpdateRequest) SetSkuId(v string)`

SetSkuId sets SkuId field to given value.

### HasSkuId

`func (o *ExpectedMachineUpdateRequest) HasSkuId() bool`

HasSkuId returns a boolean if a field has been set.

### SetSkuIdNil

`func (o *ExpectedMachineUpdateRequest) SetSkuIdNil(b bool)`

 SetSkuIdNil sets the value for SkuId to be an explicit nil

### UnsetSkuId
`func (o *ExpectedMachineUpdateRequest) UnsetSkuId()`

UnsetSkuId ensures that no value is present for SkuId, not even an explicit nil
### GetLabels

`func (o *ExpectedMachineUpdateRequest) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *ExpectedMachineUpdateRequest) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *ExpectedMachineUpdateRequest) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *ExpectedMachineUpdateRequest) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


