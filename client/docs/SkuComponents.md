# SkuComponents

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cpus** | Pointer to [**[]SkuCpu**](SkuCpu.md) | CPU components | [optional] 
**Gpus** | Pointer to [**[]SkuGpu**](SkuGpu.md) | GPU components | [optional] 
**Memory** | Pointer to [**[]SkuMemory**](SkuMemory.md) | Memory components | [optional] 
**Storage** | Pointer to [**[]SkuStorage**](SkuStorage.md) | Storage components | [optional] 
**Chassis** | Pointer to [**SkuChassis**](SkuChassis.md) | Chassis component | [optional] 
**EthernetDevices** | Pointer to [**[]SkuEthernetDevice**](SkuEthernetDevice.md) | Ethernet device components | [optional] 
**InfinibandDevices** | Pointer to [**[]SkuInfinibandDevice**](SkuInfinibandDevice.md) | Infiniband device components | [optional] 
**Tpm** | Pointer to [**[]SkuTpm**](SkuTpm.md) | TPM components | [optional] 

## Methods

### NewSkuComponents

`func NewSkuComponents() *SkuComponents`

NewSkuComponents instantiates a new SkuComponents object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSkuComponentsWithDefaults

`func NewSkuComponentsWithDefaults() *SkuComponents`

NewSkuComponentsWithDefaults instantiates a new SkuComponents object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCpus

`func (o *SkuComponents) GetCpus() []SkuCpu`

GetCpus returns the Cpus field if non-nil, zero value otherwise.

### GetCpusOk

`func (o *SkuComponents) GetCpusOk() (*[]SkuCpu, bool)`

GetCpusOk returns a tuple with the Cpus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpus

`func (o *SkuComponents) SetCpus(v []SkuCpu)`

SetCpus sets Cpus field to given value.

### HasCpus

`func (o *SkuComponents) HasCpus() bool`

HasCpus returns a boolean if a field has been set.

### GetGpus

`func (o *SkuComponents) GetGpus() []SkuGpu`

GetGpus returns the Gpus field if non-nil, zero value otherwise.

### GetGpusOk

`func (o *SkuComponents) GetGpusOk() (*[]SkuGpu, bool)`

GetGpusOk returns a tuple with the Gpus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGpus

`func (o *SkuComponents) SetGpus(v []SkuGpu)`

SetGpus sets Gpus field to given value.

### HasGpus

`func (o *SkuComponents) HasGpus() bool`

HasGpus returns a boolean if a field has been set.

### GetMemory

`func (o *SkuComponents) GetMemory() []SkuMemory`

GetMemory returns the Memory field if non-nil, zero value otherwise.

### GetMemoryOk

`func (o *SkuComponents) GetMemoryOk() (*[]SkuMemory, bool)`

GetMemoryOk returns a tuple with the Memory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemory

`func (o *SkuComponents) SetMemory(v []SkuMemory)`

SetMemory sets Memory field to given value.

### HasMemory

`func (o *SkuComponents) HasMemory() bool`

HasMemory returns a boolean if a field has been set.

### GetStorage

`func (o *SkuComponents) GetStorage() []SkuStorage`

GetStorage returns the Storage field if non-nil, zero value otherwise.

### GetStorageOk

`func (o *SkuComponents) GetStorageOk() (*[]SkuStorage, bool)`

GetStorageOk returns a tuple with the Storage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStorage

`func (o *SkuComponents) SetStorage(v []SkuStorage)`

SetStorage sets Storage field to given value.

### HasStorage

`func (o *SkuComponents) HasStorage() bool`

HasStorage returns a boolean if a field has been set.

### GetChassis

`func (o *SkuComponents) GetChassis() SkuChassis`

GetChassis returns the Chassis field if non-nil, zero value otherwise.

### GetChassisOk

`func (o *SkuComponents) GetChassisOk() (*SkuChassis, bool)`

GetChassisOk returns a tuple with the Chassis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChassis

`func (o *SkuComponents) SetChassis(v SkuChassis)`

SetChassis sets Chassis field to given value.

### HasChassis

`func (o *SkuComponents) HasChassis() bool`

HasChassis returns a boolean if a field has been set.

### GetEthernetDevices

`func (o *SkuComponents) GetEthernetDevices() []SkuEthernetDevice`

GetEthernetDevices returns the EthernetDevices field if non-nil, zero value otherwise.

### GetEthernetDevicesOk

`func (o *SkuComponents) GetEthernetDevicesOk() (*[]SkuEthernetDevice, bool)`

GetEthernetDevicesOk returns a tuple with the EthernetDevices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthernetDevices

`func (o *SkuComponents) SetEthernetDevices(v []SkuEthernetDevice)`

SetEthernetDevices sets EthernetDevices field to given value.

### HasEthernetDevices

`func (o *SkuComponents) HasEthernetDevices() bool`

HasEthernetDevices returns a boolean if a field has been set.

### GetInfinibandDevices

`func (o *SkuComponents) GetInfinibandDevices() []SkuInfinibandDevice`

GetInfinibandDevices returns the InfinibandDevices field if non-nil, zero value otherwise.

### GetInfinibandDevicesOk

`func (o *SkuComponents) GetInfinibandDevicesOk() (*[]SkuInfinibandDevice, bool)`

GetInfinibandDevicesOk returns a tuple with the InfinibandDevices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfinibandDevices

`func (o *SkuComponents) SetInfinibandDevices(v []SkuInfinibandDevice)`

SetInfinibandDevices sets InfinibandDevices field to given value.

### HasInfinibandDevices

`func (o *SkuComponents) HasInfinibandDevices() bool`

HasInfinibandDevices returns a boolean if a field has been set.

### GetTpm

`func (o *SkuComponents) GetTpm() []SkuTpm`

GetTpm returns the Tpm field if non-nil, zero value otherwise.

### GetTpmOk

`func (o *SkuComponents) GetTpmOk() (*[]SkuTpm, bool)`

GetTpmOk returns a tuple with the Tpm field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTpm

`func (o *SkuComponents) SetTpm(v []SkuTpm)`

SetTpm sets Tpm field to given value.

### HasTpm

`func (o *SkuComponents) HasTpm() bool`

HasTpm returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


