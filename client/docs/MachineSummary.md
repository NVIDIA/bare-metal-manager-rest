# MachineSummary

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique ID of Machine across Carbide Cloud &amp; Site | [optional] [readonly] 
**ControllerMachineId** | Pointer to **string** | ID of the Machine at Site, now same as the primary ID | [optional] 
**ControllerMachineType** | Pointer to **NullableString** | Denotes architecture (x86 vs ARM) of the Machine | [optional] 
**HwSkuDeviceType** | Pointer to **NullableString** | SKU derived device type of the machine, e.g. cpu, gpu, cache, storage, etc. | [optional] 
**Vendor** | Pointer to **NullableString** | Name of the vendor of the Machine | [optional] 
**ProductName** | Pointer to **NullableString** | Product name of the Machine | [optional] 
**MaintenanceMessage** | Pointer to **NullableString** | If the Machine is in maintenance mode, this message will typically describe the reason and how long it is expected to be in maintenance | [optional] 
**Status** | Pointer to [**MachineStatus**](MachineStatus.md) |  | [optional] [readonly] 

## Methods

### NewMachineSummary

`func NewMachineSummary() *MachineSummary`

NewMachineSummary instantiates a new MachineSummary object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineSummaryWithDefaults

`func NewMachineSummaryWithDefaults() *MachineSummary`

NewMachineSummaryWithDefaults instantiates a new MachineSummary object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *MachineSummary) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *MachineSummary) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *MachineSummary) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *MachineSummary) HasId() bool`

HasId returns a boolean if a field has been set.

### GetControllerMachineId

`func (o *MachineSummary) GetControllerMachineId() string`

GetControllerMachineId returns the ControllerMachineId field if non-nil, zero value otherwise.

### GetControllerMachineIdOk

`func (o *MachineSummary) GetControllerMachineIdOk() (*string, bool)`

GetControllerMachineIdOk returns a tuple with the ControllerMachineId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerMachineId

`func (o *MachineSummary) SetControllerMachineId(v string)`

SetControllerMachineId sets ControllerMachineId field to given value.

### HasControllerMachineId

`func (o *MachineSummary) HasControllerMachineId() bool`

HasControllerMachineId returns a boolean if a field has been set.

### GetControllerMachineType

`func (o *MachineSummary) GetControllerMachineType() string`

GetControllerMachineType returns the ControllerMachineType field if non-nil, zero value otherwise.

### GetControllerMachineTypeOk

`func (o *MachineSummary) GetControllerMachineTypeOk() (*string, bool)`

GetControllerMachineTypeOk returns a tuple with the ControllerMachineType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerMachineType

`func (o *MachineSummary) SetControllerMachineType(v string)`

SetControllerMachineType sets ControllerMachineType field to given value.

### HasControllerMachineType

`func (o *MachineSummary) HasControllerMachineType() bool`

HasControllerMachineType returns a boolean if a field has been set.

### SetControllerMachineTypeNil

`func (o *MachineSummary) SetControllerMachineTypeNil(b bool)`

 SetControllerMachineTypeNil sets the value for ControllerMachineType to be an explicit nil

### UnsetControllerMachineType
`func (o *MachineSummary) UnsetControllerMachineType()`

UnsetControllerMachineType ensures that no value is present for ControllerMachineType, not even an explicit nil
### GetHwSkuDeviceType

`func (o *MachineSummary) GetHwSkuDeviceType() string`

GetHwSkuDeviceType returns the HwSkuDeviceType field if non-nil, zero value otherwise.

### GetHwSkuDeviceTypeOk

`func (o *MachineSummary) GetHwSkuDeviceTypeOk() (*string, bool)`

GetHwSkuDeviceTypeOk returns a tuple with the HwSkuDeviceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHwSkuDeviceType

`func (o *MachineSummary) SetHwSkuDeviceType(v string)`

SetHwSkuDeviceType sets HwSkuDeviceType field to given value.

### HasHwSkuDeviceType

`func (o *MachineSummary) HasHwSkuDeviceType() bool`

HasHwSkuDeviceType returns a boolean if a field has been set.

### SetHwSkuDeviceTypeNil

`func (o *MachineSummary) SetHwSkuDeviceTypeNil(b bool)`

 SetHwSkuDeviceTypeNil sets the value for HwSkuDeviceType to be an explicit nil

### UnsetHwSkuDeviceType
`func (o *MachineSummary) UnsetHwSkuDeviceType()`

UnsetHwSkuDeviceType ensures that no value is present for HwSkuDeviceType, not even an explicit nil
### GetVendor

`func (o *MachineSummary) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *MachineSummary) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *MachineSummary) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *MachineSummary) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### SetVendorNil

`func (o *MachineSummary) SetVendorNil(b bool)`

 SetVendorNil sets the value for Vendor to be an explicit nil

### UnsetVendor
`func (o *MachineSummary) UnsetVendor()`

UnsetVendor ensures that no value is present for Vendor, not even an explicit nil
### GetProductName

`func (o *MachineSummary) GetProductName() string`

GetProductName returns the ProductName field if non-nil, zero value otherwise.

### GetProductNameOk

`func (o *MachineSummary) GetProductNameOk() (*string, bool)`

GetProductNameOk returns a tuple with the ProductName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProductName

`func (o *MachineSummary) SetProductName(v string)`

SetProductName sets ProductName field to given value.

### HasProductName

`func (o *MachineSummary) HasProductName() bool`

HasProductName returns a boolean if a field has been set.

### SetProductNameNil

`func (o *MachineSummary) SetProductNameNil(b bool)`

 SetProductNameNil sets the value for ProductName to be an explicit nil

### UnsetProductName
`func (o *MachineSummary) UnsetProductName()`

UnsetProductName ensures that no value is present for ProductName, not even an explicit nil
### GetMaintenanceMessage

`func (o *MachineSummary) GetMaintenanceMessage() string`

GetMaintenanceMessage returns the MaintenanceMessage field if non-nil, zero value otherwise.

### GetMaintenanceMessageOk

`func (o *MachineSummary) GetMaintenanceMessageOk() (*string, bool)`

GetMaintenanceMessageOk returns a tuple with the MaintenanceMessage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaintenanceMessage

`func (o *MachineSummary) SetMaintenanceMessage(v string)`

SetMaintenanceMessage sets MaintenanceMessage field to given value.

### HasMaintenanceMessage

`func (o *MachineSummary) HasMaintenanceMessage() bool`

HasMaintenanceMessage returns a boolean if a field has been set.

### SetMaintenanceMessageNil

`func (o *MachineSummary) SetMaintenanceMessageNil(b bool)`

 SetMaintenanceMessageNil sets the value for MaintenanceMessage to be an explicit nil

### UnsetMaintenanceMessage
`func (o *MachineSummary) UnsetMaintenanceMessage()`

UnsetMaintenanceMessage ensures that no value is present for MaintenanceMessage, not even an explicit nil
### GetStatus

`func (o *MachineSummary) GetStatus() MachineStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *MachineSummary) GetStatusOk() (*MachineStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *MachineSummary) SetStatus(v MachineStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *MachineSummary) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


