# BatchUpdateExpectedMachinesRequestInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | ID of the Expected Machine to update (required for batch operations) | 
**BmcMacAddress** | Pointer to **string** | MAC address of the Expected Machine&#39;s BMC (Baseboard Management Controller) | [optional] 
**BmcUsername** | Pointer to **string** | Username for accessing the Expected Machine&#39;s BMC | [optional] 
**BmcPassword** | Pointer to **string** | Password for accessing the Expected Machine&#39;s BMC | [optional] 
**ChassisSerialNumber** | Pointer to **string** | Serial number of the Expected Machine&#39;s chassis | [optional] 
**FallbackDPUSerialNumbers** | Pointer to **[]string** | Serial numbers of the Expected Machine&#39;s fallback DPUs (Data Processing Units) | [optional] 
**SkuId** | Pointer to **string** | Optional ID of the SKU to associate with this Expected Machine | [optional] 
**Labels** | Pointer to **map[string]string** | User-defined key-value pairs for organizing and categorizing Expected Machines | [optional] 

## Methods

### NewBatchUpdateExpectedMachinesRequestInner

`func NewBatchUpdateExpectedMachinesRequestInner(id string, ) *BatchUpdateExpectedMachinesRequestInner`

NewBatchUpdateExpectedMachinesRequestInner instantiates a new BatchUpdateExpectedMachinesRequestInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBatchUpdateExpectedMachinesRequestInnerWithDefaults

`func NewBatchUpdateExpectedMachinesRequestInnerWithDefaults() *BatchUpdateExpectedMachinesRequestInner`

NewBatchUpdateExpectedMachinesRequestInnerWithDefaults instantiates a new BatchUpdateExpectedMachinesRequestInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *BatchUpdateExpectedMachinesRequestInner) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *BatchUpdateExpectedMachinesRequestInner) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *BatchUpdateExpectedMachinesRequestInner) SetId(v string)`

SetId sets Id field to given value.


### GetBmcMacAddress

`func (o *BatchUpdateExpectedMachinesRequestInner) GetBmcMacAddress() string`

GetBmcMacAddress returns the BmcMacAddress field if non-nil, zero value otherwise.

### GetBmcMacAddressOk

`func (o *BatchUpdateExpectedMachinesRequestInner) GetBmcMacAddressOk() (*string, bool)`

GetBmcMacAddressOk returns a tuple with the BmcMacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcMacAddress

`func (o *BatchUpdateExpectedMachinesRequestInner) SetBmcMacAddress(v string)`

SetBmcMacAddress sets BmcMacAddress field to given value.

### HasBmcMacAddress

`func (o *BatchUpdateExpectedMachinesRequestInner) HasBmcMacAddress() bool`

HasBmcMacAddress returns a boolean if a field has been set.

### GetBmcUsername

`func (o *BatchUpdateExpectedMachinesRequestInner) GetBmcUsername() string`

GetBmcUsername returns the BmcUsername field if non-nil, zero value otherwise.

### GetBmcUsernameOk

`func (o *BatchUpdateExpectedMachinesRequestInner) GetBmcUsernameOk() (*string, bool)`

GetBmcUsernameOk returns a tuple with the BmcUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcUsername

`func (o *BatchUpdateExpectedMachinesRequestInner) SetBmcUsername(v string)`

SetBmcUsername sets BmcUsername field to given value.

### HasBmcUsername

`func (o *BatchUpdateExpectedMachinesRequestInner) HasBmcUsername() bool`

HasBmcUsername returns a boolean if a field has been set.

### GetBmcPassword

`func (o *BatchUpdateExpectedMachinesRequestInner) GetBmcPassword() string`

GetBmcPassword returns the BmcPassword field if non-nil, zero value otherwise.

### GetBmcPasswordOk

`func (o *BatchUpdateExpectedMachinesRequestInner) GetBmcPasswordOk() (*string, bool)`

GetBmcPasswordOk returns a tuple with the BmcPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcPassword

`func (o *BatchUpdateExpectedMachinesRequestInner) SetBmcPassword(v string)`

SetBmcPassword sets BmcPassword field to given value.

### HasBmcPassword

`func (o *BatchUpdateExpectedMachinesRequestInner) HasBmcPassword() bool`

HasBmcPassword returns a boolean if a field has been set.

### GetChassisSerialNumber

`func (o *BatchUpdateExpectedMachinesRequestInner) GetChassisSerialNumber() string`

GetChassisSerialNumber returns the ChassisSerialNumber field if non-nil, zero value otherwise.

### GetChassisSerialNumberOk

`func (o *BatchUpdateExpectedMachinesRequestInner) GetChassisSerialNumberOk() (*string, bool)`

GetChassisSerialNumberOk returns a tuple with the ChassisSerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChassisSerialNumber

`func (o *BatchUpdateExpectedMachinesRequestInner) SetChassisSerialNumber(v string)`

SetChassisSerialNumber sets ChassisSerialNumber field to given value.

### HasChassisSerialNumber

`func (o *BatchUpdateExpectedMachinesRequestInner) HasChassisSerialNumber() bool`

HasChassisSerialNumber returns a boolean if a field has been set.

### GetFallbackDPUSerialNumbers

`func (o *BatchUpdateExpectedMachinesRequestInner) GetFallbackDPUSerialNumbers() []string`

GetFallbackDPUSerialNumbers returns the FallbackDPUSerialNumbers field if non-nil, zero value otherwise.

### GetFallbackDPUSerialNumbersOk

`func (o *BatchUpdateExpectedMachinesRequestInner) GetFallbackDPUSerialNumbersOk() (*[]string, bool)`

GetFallbackDPUSerialNumbersOk returns a tuple with the FallbackDPUSerialNumbers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFallbackDPUSerialNumbers

`func (o *BatchUpdateExpectedMachinesRequestInner) SetFallbackDPUSerialNumbers(v []string)`

SetFallbackDPUSerialNumbers sets FallbackDPUSerialNumbers field to given value.

### HasFallbackDPUSerialNumbers

`func (o *BatchUpdateExpectedMachinesRequestInner) HasFallbackDPUSerialNumbers() bool`

HasFallbackDPUSerialNumbers returns a boolean if a field has been set.

### GetSkuId

`func (o *BatchUpdateExpectedMachinesRequestInner) GetSkuId() string`

GetSkuId returns the SkuId field if non-nil, zero value otherwise.

### GetSkuIdOk

`func (o *BatchUpdateExpectedMachinesRequestInner) GetSkuIdOk() (*string, bool)`

GetSkuIdOk returns a tuple with the SkuId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkuId

`func (o *BatchUpdateExpectedMachinesRequestInner) SetSkuId(v string)`

SetSkuId sets SkuId field to given value.

### HasSkuId

`func (o *BatchUpdateExpectedMachinesRequestInner) HasSkuId() bool`

HasSkuId returns a boolean if a field has been set.

### GetLabels

`func (o *BatchUpdateExpectedMachinesRequestInner) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *BatchUpdateExpectedMachinesRequestInner) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *BatchUpdateExpectedMachinesRequestInner) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *BatchUpdateExpectedMachinesRequestInner) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


