# MachineMetadata

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DmiData** | Pointer to [**MachineDMIData**](MachineDMIData.md) |  | [optional] 
**BmcInfo** | Pointer to [**MachineBMCInfo**](MachineBMCInfo.md) |  | [optional] 
**Gpus** | Pointer to [**[]MachineGPUInfo**](MachineGPUInfo.md) |  | [optional] 
**NetworkInterfaces** | Pointer to [**[]MachineNetworkInterface**](MachineNetworkInterface.md) |  | [optional] 
**InfinibandInterfaces** | Pointer to [**[]MachineInfiniBandInterface**](MachineInfiniBandInterface.md) |  | [optional] 

## Methods

### NewMachineMetadata

`func NewMachineMetadata() *MachineMetadata`

NewMachineMetadata instantiates a new MachineMetadata object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineMetadataWithDefaults

`func NewMachineMetadataWithDefaults() *MachineMetadata`

NewMachineMetadataWithDefaults instantiates a new MachineMetadata object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDmiData

`func (o *MachineMetadata) GetDmiData() MachineDMIData`

GetDmiData returns the DmiData field if non-nil, zero value otherwise.

### GetDmiDataOk

`func (o *MachineMetadata) GetDmiDataOk() (*MachineDMIData, bool)`

GetDmiDataOk returns a tuple with the DmiData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDmiData

`func (o *MachineMetadata) SetDmiData(v MachineDMIData)`

SetDmiData sets DmiData field to given value.

### HasDmiData

`func (o *MachineMetadata) HasDmiData() bool`

HasDmiData returns a boolean if a field has been set.

### GetBmcInfo

`func (o *MachineMetadata) GetBmcInfo() MachineBMCInfo`

GetBmcInfo returns the BmcInfo field if non-nil, zero value otherwise.

### GetBmcInfoOk

`func (o *MachineMetadata) GetBmcInfoOk() (*MachineBMCInfo, bool)`

GetBmcInfoOk returns a tuple with the BmcInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBmcInfo

`func (o *MachineMetadata) SetBmcInfo(v MachineBMCInfo)`

SetBmcInfo sets BmcInfo field to given value.

### HasBmcInfo

`func (o *MachineMetadata) HasBmcInfo() bool`

HasBmcInfo returns a boolean if a field has been set.

### GetGpus

`func (o *MachineMetadata) GetGpus() []MachineGPUInfo`

GetGpus returns the Gpus field if non-nil, zero value otherwise.

### GetGpusOk

`func (o *MachineMetadata) GetGpusOk() (*[]MachineGPUInfo, bool)`

GetGpusOk returns a tuple with the Gpus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGpus

`func (o *MachineMetadata) SetGpus(v []MachineGPUInfo)`

SetGpus sets Gpus field to given value.

### HasGpus

`func (o *MachineMetadata) HasGpus() bool`

HasGpus returns a boolean if a field has been set.

### GetNetworkInterfaces

`func (o *MachineMetadata) GetNetworkInterfaces() []MachineNetworkInterface`

GetNetworkInterfaces returns the NetworkInterfaces field if non-nil, zero value otherwise.

### GetNetworkInterfacesOk

`func (o *MachineMetadata) GetNetworkInterfacesOk() (*[]MachineNetworkInterface, bool)`

GetNetworkInterfacesOk returns a tuple with the NetworkInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkInterfaces

`func (o *MachineMetadata) SetNetworkInterfaces(v []MachineNetworkInterface)`

SetNetworkInterfaces sets NetworkInterfaces field to given value.

### HasNetworkInterfaces

`func (o *MachineMetadata) HasNetworkInterfaces() bool`

HasNetworkInterfaces returns a boolean if a field has been set.

### GetInfinibandInterfaces

`func (o *MachineMetadata) GetInfinibandInterfaces() []MachineInfiniBandInterface`

GetInfinibandInterfaces returns the InfinibandInterfaces field if non-nil, zero value otherwise.

### GetInfinibandInterfacesOk

`func (o *MachineMetadata) GetInfinibandInterfacesOk() (*[]MachineInfiniBandInterface, bool)`

GetInfinibandInterfacesOk returns a tuple with the InfinibandInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfinibandInterfaces

`func (o *MachineMetadata) SetInfinibandInterfaces(v []MachineInfiniBandInterface)`

SetInfinibandInterfaces sets InfinibandInterfaces field to given value.

### HasInfinibandInterfaces

`func (o *MachineMetadata) HasInfinibandInterfaces() bool`

HasInfinibandInterfaces returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


