# MachineInterface

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**MachineId** | Pointer to **string** |  | [optional] 
**ControllerInterfaceId** | Pointer to **string** |  | [optional] 
**ControllerSegmentId** | Pointer to **string** |  | [optional] 
**AttachedDpuMachineID** | Pointer to **NullableString** |  | [optional] 
**SubnetId** | Pointer to **NullableString** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**IsPrimary** | Pointer to **bool** |  | [optional] 
**MacAddress** | Pointer to **string** |  | [optional] 
**IpAddresses** | Pointer to **[]string** | Array of IP addresses | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewMachineInterface

`func NewMachineInterface() *MachineInterface`

NewMachineInterface instantiates a new MachineInterface object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineInterfaceWithDefaults

`func NewMachineInterfaceWithDefaults() *MachineInterface`

NewMachineInterfaceWithDefaults instantiates a new MachineInterface object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *MachineInterface) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *MachineInterface) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *MachineInterface) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *MachineInterface) HasId() bool`

HasId returns a boolean if a field has been set.

### GetMachineId

`func (o *MachineInterface) GetMachineId() string`

GetMachineId returns the MachineId field if non-nil, zero value otherwise.

### GetMachineIdOk

`func (o *MachineInterface) GetMachineIdOk() (*string, bool)`

GetMachineIdOk returns a tuple with the MachineId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineId

`func (o *MachineInterface) SetMachineId(v string)`

SetMachineId sets MachineId field to given value.

### HasMachineId

`func (o *MachineInterface) HasMachineId() bool`

HasMachineId returns a boolean if a field has been set.

### GetControllerInterfaceId

`func (o *MachineInterface) GetControllerInterfaceId() string`

GetControllerInterfaceId returns the ControllerInterfaceId field if non-nil, zero value otherwise.

### GetControllerInterfaceIdOk

`func (o *MachineInterface) GetControllerInterfaceIdOk() (*string, bool)`

GetControllerInterfaceIdOk returns a tuple with the ControllerInterfaceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerInterfaceId

`func (o *MachineInterface) SetControllerInterfaceId(v string)`

SetControllerInterfaceId sets ControllerInterfaceId field to given value.

### HasControllerInterfaceId

`func (o *MachineInterface) HasControllerInterfaceId() bool`

HasControllerInterfaceId returns a boolean if a field has been set.

### GetControllerSegmentId

`func (o *MachineInterface) GetControllerSegmentId() string`

GetControllerSegmentId returns the ControllerSegmentId field if non-nil, zero value otherwise.

### GetControllerSegmentIdOk

`func (o *MachineInterface) GetControllerSegmentIdOk() (*string, bool)`

GetControllerSegmentIdOk returns a tuple with the ControllerSegmentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerSegmentId

`func (o *MachineInterface) SetControllerSegmentId(v string)`

SetControllerSegmentId sets ControllerSegmentId field to given value.

### HasControllerSegmentId

`func (o *MachineInterface) HasControllerSegmentId() bool`

HasControllerSegmentId returns a boolean if a field has been set.

### GetAttachedDpuMachineID

`func (o *MachineInterface) GetAttachedDpuMachineID() string`

GetAttachedDpuMachineID returns the AttachedDpuMachineID field if non-nil, zero value otherwise.

### GetAttachedDpuMachineIDOk

`func (o *MachineInterface) GetAttachedDpuMachineIDOk() (*string, bool)`

GetAttachedDpuMachineIDOk returns a tuple with the AttachedDpuMachineID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttachedDpuMachineID

`func (o *MachineInterface) SetAttachedDpuMachineID(v string)`

SetAttachedDpuMachineID sets AttachedDpuMachineID field to given value.

### HasAttachedDpuMachineID

`func (o *MachineInterface) HasAttachedDpuMachineID() bool`

HasAttachedDpuMachineID returns a boolean if a field has been set.

### SetAttachedDpuMachineIDNil

`func (o *MachineInterface) SetAttachedDpuMachineIDNil(b bool)`

 SetAttachedDpuMachineIDNil sets the value for AttachedDpuMachineID to be an explicit nil

### UnsetAttachedDpuMachineID
`func (o *MachineInterface) UnsetAttachedDpuMachineID()`

UnsetAttachedDpuMachineID ensures that no value is present for AttachedDpuMachineID, not even an explicit nil
### GetSubnetId

`func (o *MachineInterface) GetSubnetId() string`

GetSubnetId returns the SubnetId field if non-nil, zero value otherwise.

### GetSubnetIdOk

`func (o *MachineInterface) GetSubnetIdOk() (*string, bool)`

GetSubnetIdOk returns a tuple with the SubnetId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubnetId

`func (o *MachineInterface) SetSubnetId(v string)`

SetSubnetId sets SubnetId field to given value.

### HasSubnetId

`func (o *MachineInterface) HasSubnetId() bool`

HasSubnetId returns a boolean if a field has been set.

### SetSubnetIdNil

`func (o *MachineInterface) SetSubnetIdNil(b bool)`

 SetSubnetIdNil sets the value for SubnetId to be an explicit nil

### UnsetSubnetId
`func (o *MachineInterface) UnsetSubnetId()`

UnsetSubnetId ensures that no value is present for SubnetId, not even an explicit nil
### GetHostname

`func (o *MachineInterface) GetHostname() string`

GetHostname returns the Hostname field if non-nil, zero value otherwise.

### GetHostnameOk

`func (o *MachineInterface) GetHostnameOk() (*string, bool)`

GetHostnameOk returns a tuple with the Hostname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHostname

`func (o *MachineInterface) SetHostname(v string)`

SetHostname sets Hostname field to given value.

### HasHostname

`func (o *MachineInterface) HasHostname() bool`

HasHostname returns a boolean if a field has been set.

### GetIsPrimary

`func (o *MachineInterface) GetIsPrimary() bool`

GetIsPrimary returns the IsPrimary field if non-nil, zero value otherwise.

### GetIsPrimaryOk

`func (o *MachineInterface) GetIsPrimaryOk() (*bool, bool)`

GetIsPrimaryOk returns a tuple with the IsPrimary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPrimary

`func (o *MachineInterface) SetIsPrimary(v bool)`

SetIsPrimary sets IsPrimary field to given value.

### HasIsPrimary

`func (o *MachineInterface) HasIsPrimary() bool`

HasIsPrimary returns a boolean if a field has been set.

### GetMacAddress

`func (o *MachineInterface) GetMacAddress() string`

GetMacAddress returns the MacAddress field if non-nil, zero value otherwise.

### GetMacAddressOk

`func (o *MachineInterface) GetMacAddressOk() (*string, bool)`

GetMacAddressOk returns a tuple with the MacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacAddress

`func (o *MachineInterface) SetMacAddress(v string)`

SetMacAddress sets MacAddress field to given value.

### HasMacAddress

`func (o *MachineInterface) HasMacAddress() bool`

HasMacAddress returns a boolean if a field has been set.

### GetIpAddresses

`func (o *MachineInterface) GetIpAddresses() []string`

GetIpAddresses returns the IpAddresses field if non-nil, zero value otherwise.

### GetIpAddressesOk

`func (o *MachineInterface) GetIpAddressesOk() (*[]string, bool)`

GetIpAddressesOk returns a tuple with the IpAddresses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpAddresses

`func (o *MachineInterface) SetIpAddresses(v []string)`

SetIpAddresses sets IpAddresses field to given value.

### HasIpAddresses

`func (o *MachineInterface) HasIpAddresses() bool`

HasIpAddresses returns a boolean if a field has been set.

### GetCreated

`func (o *MachineInterface) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *MachineInterface) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *MachineInterface) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *MachineInterface) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *MachineInterface) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *MachineInterface) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *MachineInterface) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *MachineInterface) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


