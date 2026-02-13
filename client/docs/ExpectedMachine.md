# ExpectedMachine

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier for the Expected Machine | [optional] [readonly] 
**SiteId** | Pointer to **string** | ID of the site the Expected Machine belongs to | [optional] [readonly] 
**BmcMacAddress** | Pointer to **string** | MAC address of the Expected Machine&#39;s BMC (Baseboard Management Controller) | [optional] 
**ChassisSerialNumber** | Pointer to **string** | Serial number of the Expected Machine&#39;s chassis | [optional] 
**FallbackDPUSerialNumbers** | Pointer to **[]string** | Serial numbers of the Expected Machine&#39;s fallback DPUs (Data Processing Units) | [optional] 
**SkuId** | Pointer to **NullableString** | Optional ID of the SKU associated with this Expected Machine | [optional] 
**Sku** | Pointer to [**Sku**](Sku.md) | SKU information for this Expected Machine (populated when includeRelation&#x3D;Sku is specified) | [optional] 
**MachineId** | Pointer to **NullableString** | Optional ID of the Machine associated with this Expected Machine | [optional] 
**Machine** | Pointer to [**MachineSummary**](MachineSummary.md) | Machine information for this Expected Machine (populated when includeRelation&#x3D;Machine is specified) | [optional] 
**Labels** | Pointer to **map[string]string** | User-defined key-value pairs for organizing and categorizing Expected Machines | [optional] 
**Created** | Pointer to **time.Time** | ISO 8601 datetime when the Expected Machine was created | [optional] [readonly] 
**Updated** | Pointer to **time.Time** | ISO 8601 datetime when the Expected Machine was last updated | [optional] [readonly] 

## Methods

### NewExpectedMachine

`func NewExpectedMachine() *ExpectedMachine`

NewExpectedMachine instantiates a new ExpectedMachine object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExpectedMachineWithDefaults

`func NewExpectedMachineWithDefaults() *ExpectedMachine`

NewExpectedMachineWithDefaults instantiates a new ExpectedMachine object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ExpectedMachine) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ExpectedMachine) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ExpectedMachine) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ExpectedMachine) HasId() bool`

HasId returns a boolean if a field has been set.

### GetSiteId

`func (o *ExpectedMachine) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *ExpectedMachine) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *ExpectedMachine) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *ExpectedMachine) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetBmcMacAddress

`func (o *ExpectedMachine) GetBmcMacAddress() string`

GetBmcMacAddress returns the BmcMacAddress field if non-nil, zero value otherwise.

### GetBmcMacAddressOk

`func (o *ExpectedMachine) GetBmcMacAddressOk() (*string, bool)`

GetBmcMacAddressOk returns a tuple with the BmcMacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcMacAddress

`func (o *ExpectedMachine) SetBmcMacAddress(v string)`

SetBmcMacAddress sets BmcMacAddress field to given value.

### HasBmcMacAddress

`func (o *ExpectedMachine) HasBmcMacAddress() bool`

HasBmcMacAddress returns a boolean if a field has been set.

### GetChassisSerialNumber

`func (o *ExpectedMachine) GetChassisSerialNumber() string`

GetChassisSerialNumber returns the ChassisSerialNumber field if non-nil, zero value otherwise.

### GetChassisSerialNumberOk

`func (o *ExpectedMachine) GetChassisSerialNumberOk() (*string, bool)`

GetChassisSerialNumberOk returns a tuple with the ChassisSerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChassisSerialNumber

`func (o *ExpectedMachine) SetChassisSerialNumber(v string)`

SetChassisSerialNumber sets ChassisSerialNumber field to given value.

### HasChassisSerialNumber

`func (o *ExpectedMachine) HasChassisSerialNumber() bool`

HasChassisSerialNumber returns a boolean if a field has been set.

### GetFallbackDPUSerialNumbers

`func (o *ExpectedMachine) GetFallbackDPUSerialNumbers() []string`

GetFallbackDPUSerialNumbers returns the FallbackDPUSerialNumbers field if non-nil, zero value otherwise.

### GetFallbackDPUSerialNumbersOk

`func (o *ExpectedMachine) GetFallbackDPUSerialNumbersOk() (*[]string, bool)`

GetFallbackDPUSerialNumbersOk returns a tuple with the FallbackDPUSerialNumbers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFallbackDPUSerialNumbers

`func (o *ExpectedMachine) SetFallbackDPUSerialNumbers(v []string)`

SetFallbackDPUSerialNumbers sets FallbackDPUSerialNumbers field to given value.

### HasFallbackDPUSerialNumbers

`func (o *ExpectedMachine) HasFallbackDPUSerialNumbers() bool`

HasFallbackDPUSerialNumbers returns a boolean if a field has been set.

### GetSkuId

`func (o *ExpectedMachine) GetSkuId() string`

GetSkuId returns the SkuId field if non-nil, zero value otherwise.

### GetSkuIdOk

`func (o *ExpectedMachine) GetSkuIdOk() (*string, bool)`

GetSkuIdOk returns a tuple with the SkuId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkuId

`func (o *ExpectedMachine) SetSkuId(v string)`

SetSkuId sets SkuId field to given value.

### HasSkuId

`func (o *ExpectedMachine) HasSkuId() bool`

HasSkuId returns a boolean if a field has been set.

### SetSkuIdNil

`func (o *ExpectedMachine) SetSkuIdNil(b bool)`

 SetSkuIdNil sets the value for SkuId to be an explicit nil

### UnsetSkuId
`func (o *ExpectedMachine) UnsetSkuId()`

UnsetSkuId ensures that no value is present for SkuId, not even an explicit nil
### GetSku

`func (o *ExpectedMachine) GetSku() Sku`

GetSku returns the Sku field if non-nil, zero value otherwise.

### GetSkuOk

`func (o *ExpectedMachine) GetSkuOk() (*Sku, bool)`

GetSkuOk returns a tuple with the Sku field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSku

`func (o *ExpectedMachine) SetSku(v Sku)`

SetSku sets Sku field to given value.

### HasSku

`func (o *ExpectedMachine) HasSku() bool`

HasSku returns a boolean if a field has been set.

### GetMachineId

`func (o *ExpectedMachine) GetMachineId() string`

GetMachineId returns the MachineId field if non-nil, zero value otherwise.

### GetMachineIdOk

`func (o *ExpectedMachine) GetMachineIdOk() (*string, bool)`

GetMachineIdOk returns a tuple with the MachineId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineId

`func (o *ExpectedMachine) SetMachineId(v string)`

SetMachineId sets MachineId field to given value.

### HasMachineId

`func (o *ExpectedMachine) HasMachineId() bool`

HasMachineId returns a boolean if a field has been set.

### SetMachineIdNil

`func (o *ExpectedMachine) SetMachineIdNil(b bool)`

 SetMachineIdNil sets the value for MachineId to be an explicit nil

### UnsetMachineId
`func (o *ExpectedMachine) UnsetMachineId()`

UnsetMachineId ensures that no value is present for MachineId, not even an explicit nil
### GetMachine

`func (o *ExpectedMachine) GetMachine() MachineSummary`

GetMachine returns the Machine field if non-nil, zero value otherwise.

### GetMachineOk

`func (o *ExpectedMachine) GetMachineOk() (*MachineSummary, bool)`

GetMachineOk returns a tuple with the Machine field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachine

`func (o *ExpectedMachine) SetMachine(v MachineSummary)`

SetMachine sets Machine field to given value.

### HasMachine

`func (o *ExpectedMachine) HasMachine() bool`

HasMachine returns a boolean if a field has been set.

### GetLabels

`func (o *ExpectedMachine) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *ExpectedMachine) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *ExpectedMachine) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *ExpectedMachine) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetCreated

`func (o *ExpectedMachine) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *ExpectedMachine) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *ExpectedMachine) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *ExpectedMachine) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *ExpectedMachine) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *ExpectedMachine) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *ExpectedMachine) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *ExpectedMachine) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


