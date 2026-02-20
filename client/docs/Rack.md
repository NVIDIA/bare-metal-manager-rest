# Rack

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier of the Rack | [optional] 
**Name** | Pointer to **string** | Name of the Rack | [optional] 
**Manufacturer** | Pointer to **string** | Manufacturer of the Rack | [optional] 
**Model** | Pointer to **string** | Model of the Rack | [optional] 
**SerialNumber** | Pointer to **string** | Serial number of the Rack | [optional] 
**Description** | Pointer to **string** | Description of the Rack | [optional] 
**Location** | Pointer to [**RackLocation**](RackLocation.md) |  | [optional] 
**Components** | Pointer to [**[]RackComponent**](RackComponent.md) | Components within the Rack. Only returned when includeComponents is true. | [optional] 

## Methods

### NewRack

`func NewRack() *Rack`

NewRack instantiates a new Rack object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRackWithDefaults

`func NewRackWithDefaults() *Rack`

NewRackWithDefaults instantiates a new Rack object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Rack) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Rack) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Rack) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Rack) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *Rack) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Rack) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Rack) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Rack) HasName() bool`

HasName returns a boolean if a field has been set.

### GetManufacturer

`func (o *Rack) GetManufacturer() string`

GetManufacturer returns the Manufacturer field if non-nil, zero value otherwise.

### GetManufacturerOk

`func (o *Rack) GetManufacturerOk() (*string, bool)`

GetManufacturerOk returns a tuple with the Manufacturer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManufacturer

`func (o *Rack) SetManufacturer(v string)`

SetManufacturer sets Manufacturer field to given value.

### HasManufacturer

`func (o *Rack) HasManufacturer() bool`

HasManufacturer returns a boolean if a field has been set.

### GetModel

`func (o *Rack) GetModel() string`

GetModel returns the Model field if non-nil, zero value otherwise.

### GetModelOk

`func (o *Rack) GetModelOk() (*string, bool)`

GetModelOk returns a tuple with the Model field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModel

`func (o *Rack) SetModel(v string)`

SetModel sets Model field to given value.

### HasModel

`func (o *Rack) HasModel() bool`

HasModel returns a boolean if a field has been set.

### GetSerialNumber

`func (o *Rack) GetSerialNumber() string`

GetSerialNumber returns the SerialNumber field if non-nil, zero value otherwise.

### GetSerialNumberOk

`func (o *Rack) GetSerialNumberOk() (*string, bool)`

GetSerialNumberOk returns a tuple with the SerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialNumber

`func (o *Rack) SetSerialNumber(v string)`

SetSerialNumber sets SerialNumber field to given value.

### HasSerialNumber

`func (o *Rack) HasSerialNumber() bool`

HasSerialNumber returns a boolean if a field has been set.

### GetDescription

`func (o *Rack) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Rack) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Rack) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Rack) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetLocation

`func (o *Rack) GetLocation() RackLocation`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *Rack) GetLocationOk() (*RackLocation, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *Rack) SetLocation(v RackLocation)`

SetLocation sets Location field to given value.

### HasLocation

`func (o *Rack) HasLocation() bool`

HasLocation returns a boolean if a field has been set.

### GetComponents

`func (o *Rack) GetComponents() []RackComponent`

GetComponents returns the Components field if non-nil, zero value otherwise.

### GetComponentsOk

`func (o *Rack) GetComponentsOk() (*[]RackComponent, bool)`

GetComponentsOk returns a tuple with the Components field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComponents

`func (o *Rack) SetComponents(v []RackComponent)`

SetComponents sets Components field to given value.

### HasComponents

`func (o *Rack) HasComponents() bool`

HasComponents returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


