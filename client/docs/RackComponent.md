# RackComponent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier of the component in RLA | [optional] 
**ComponentId** | Pointer to **string** | ID of the component in its source system (e.g. Carbide machine ID for compute nodes) | [optional] 
**Type** | Pointer to **string** | Type of the component (e.g. COMPONENT_TYPE_COMPUTE, COMPONENT_TYPE_SWITCH) | [optional] 
**Name** | Pointer to **string** | Name of the component | [optional] 
**SerialNumber** | Pointer to **string** | Serial number of the component | [optional] 
**Manufacturer** | Pointer to **string** | Manufacturer of the component | [optional] 
**FirmwareVersion** | Pointer to **string** | Firmware version of the component | [optional] 
**Position** | Pointer to **int32** | Slot position of the component within the rack | [optional] 

## Methods

### NewRackComponent

`func NewRackComponent() *RackComponent`

NewRackComponent instantiates a new RackComponent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRackComponentWithDefaults

`func NewRackComponentWithDefaults() *RackComponent`

NewRackComponentWithDefaults instantiates a new RackComponent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RackComponent) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RackComponent) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RackComponent) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *RackComponent) HasId() bool`

HasId returns a boolean if a field has been set.

### GetComponentId

`func (o *RackComponent) GetComponentId() string`

GetComponentId returns the ComponentId field if non-nil, zero value otherwise.

### GetComponentIdOk

`func (o *RackComponent) GetComponentIdOk() (*string, bool)`

GetComponentIdOk returns a tuple with the ComponentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComponentId

`func (o *RackComponent) SetComponentId(v string)`

SetComponentId sets ComponentId field to given value.

### HasComponentId

`func (o *RackComponent) HasComponentId() bool`

HasComponentId returns a boolean if a field has been set.

### GetType

`func (o *RackComponent) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *RackComponent) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *RackComponent) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *RackComponent) HasType() bool`

HasType returns a boolean if a field has been set.

### GetName

`func (o *RackComponent) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *RackComponent) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *RackComponent) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *RackComponent) HasName() bool`

HasName returns a boolean if a field has been set.

### GetSerialNumber

`func (o *RackComponent) GetSerialNumber() string`

GetSerialNumber returns the SerialNumber field if non-nil, zero value otherwise.

### GetSerialNumberOk

`func (o *RackComponent) GetSerialNumberOk() (*string, bool)`

GetSerialNumberOk returns a tuple with the SerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialNumber

`func (o *RackComponent) SetSerialNumber(v string)`

SetSerialNumber sets SerialNumber field to given value.

### HasSerialNumber

`func (o *RackComponent) HasSerialNumber() bool`

HasSerialNumber returns a boolean if a field has been set.

### GetManufacturer

`func (o *RackComponent) GetManufacturer() string`

GetManufacturer returns the Manufacturer field if non-nil, zero value otherwise.

### GetManufacturerOk

`func (o *RackComponent) GetManufacturerOk() (*string, bool)`

GetManufacturerOk returns a tuple with the Manufacturer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManufacturer

`func (o *RackComponent) SetManufacturer(v string)`

SetManufacturer sets Manufacturer field to given value.

### HasManufacturer

`func (o *RackComponent) HasManufacturer() bool`

HasManufacturer returns a boolean if a field has been set.

### GetFirmwareVersion

`func (o *RackComponent) GetFirmwareVersion() string`

GetFirmwareVersion returns the FirmwareVersion field if non-nil, zero value otherwise.

### GetFirmwareVersionOk

`func (o *RackComponent) GetFirmwareVersionOk() (*string, bool)`

GetFirmwareVersionOk returns a tuple with the FirmwareVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirmwareVersion

`func (o *RackComponent) SetFirmwareVersion(v string)`

SetFirmwareVersion sets FirmwareVersion field to given value.

### HasFirmwareVersion

`func (o *RackComponent) HasFirmwareVersion() bool`

HasFirmwareVersion returns a boolean if a field has been set.

### GetPosition

`func (o *RackComponent) GetPosition() int32`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *RackComponent) GetPositionOk() (*int32, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *RackComponent) SetPosition(v int32)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *RackComponent) HasPosition() bool`

HasPosition returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


