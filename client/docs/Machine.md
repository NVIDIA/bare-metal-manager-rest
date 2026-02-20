# Machine

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique ID of Machine | [optional] [readonly] 
**InfrastructureProviderId** | Pointer to **string** | ID of the Provider that owns the Machine | [optional] [readonly] 
**SiteId** | Pointer to **string** | ID of the Site the Machine belongs to | [optional] 
**InstanceTypeId** | Pointer to **NullableString** | ID of the Instance Type, if assigned | [optional] 
**InstanceId** | Pointer to **NullableString** | ID of the Instance if this Machine is assigned to one | [optional] 
**TenantId** | Pointer to **NullableString** | ID of the Tenant that owns the Instance if the Machine is assigned to one | [optional] 
**ControllerMachineId** | Pointer to **string** | ID of the Machine at Site, now same as the primary ID | [optional] 
**ControllerMachineType** | Pointer to **string** | Denotes architecture (x86 vs ARM) of the Machine | [optional] 
**HwSkuDeviceType** | Pointer to **NullableString** | SKU derived device type of the machine, e.g. cpu, gpu, cache, storage, etc. | [optional] 
**Vendor** | Pointer to **string** | Name of the vendor of the Machine | [optional] 
**ProductName** | Pointer to **string** | Product name of the Machine | [optional] 
**SerialNumber** | Pointer to **string** | Serial number of the Machine, only visible to Provider | [optional] 
**MachineCapabilities** | Pointer to [**[]MachineCapability**](MachineCapability.md) |  | [optional] 
**MachineInterfaces** | Pointer to [**[]MachineInterface**](MachineInterface.md) |  | [optional] 
**MaintenanceMessage** | Pointer to **NullableString** | If the Machine is in maintenance mode, this message will typically describe the reason and how long it is expected to be in maintenance | [optional] 
**Health** | Pointer to [**MachineHealth**](MachineHealth.md) |  | [optional] 
**Metadata** | Pointer to [**MachineMetadata**](MachineMetadata.md) | Only available to Providers. Returned if includeMetadata query param is specified. Otherwise attribute is omitted from response. | [optional] 
**Labels** | Pointer to **map[string]string** |  | [optional] 
**Status** | Pointer to [**MachineStatus**](MachineStatus.md) |  | [optional] [readonly] 
**IsUsableByTenant** | Pointer to **bool** | Indicates whether the machine is usable by or currently in use by a tenant. | [optional] [readonly] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) |  | [optional] 
**Created** | Pointer to **time.Time** |  | [optional] [readonly] 
**Updated** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewMachine

`func NewMachine() *Machine`

NewMachine instantiates a new Machine object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineWithDefaults

`func NewMachineWithDefaults() *Machine`

NewMachineWithDefaults instantiates a new Machine object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Machine) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Machine) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Machine) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Machine) HasId() bool`

HasId returns a boolean if a field has been set.

### GetInfrastructureProviderId

`func (o *Machine) GetInfrastructureProviderId() string`

GetInfrastructureProviderId returns the InfrastructureProviderId field if non-nil, zero value otherwise.

### GetInfrastructureProviderIdOk

`func (o *Machine) GetInfrastructureProviderIdOk() (*string, bool)`

GetInfrastructureProviderIdOk returns a tuple with the InfrastructureProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfrastructureProviderId

`func (o *Machine) SetInfrastructureProviderId(v string)`

SetInfrastructureProviderId sets InfrastructureProviderId field to given value.

### HasInfrastructureProviderId

`func (o *Machine) HasInfrastructureProviderId() bool`

HasInfrastructureProviderId returns a boolean if a field has been set.

### GetSiteId

`func (o *Machine) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *Machine) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *Machine) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *Machine) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetInstanceTypeId

`func (o *Machine) GetInstanceTypeId() string`

GetInstanceTypeId returns the InstanceTypeId field if non-nil, zero value otherwise.

### GetInstanceTypeIdOk

`func (o *Machine) GetInstanceTypeIdOk() (*string, bool)`

GetInstanceTypeIdOk returns a tuple with the InstanceTypeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceTypeId

`func (o *Machine) SetInstanceTypeId(v string)`

SetInstanceTypeId sets InstanceTypeId field to given value.

### HasInstanceTypeId

`func (o *Machine) HasInstanceTypeId() bool`

HasInstanceTypeId returns a boolean if a field has been set.

### SetInstanceTypeIdNil

`func (o *Machine) SetInstanceTypeIdNil(b bool)`

 SetInstanceTypeIdNil sets the value for InstanceTypeId to be an explicit nil

### UnsetInstanceTypeId
`func (o *Machine) UnsetInstanceTypeId()`

UnsetInstanceTypeId ensures that no value is present for InstanceTypeId, not even an explicit nil
### GetInstanceId

`func (o *Machine) GetInstanceId() string`

GetInstanceId returns the InstanceId field if non-nil, zero value otherwise.

### GetInstanceIdOk

`func (o *Machine) GetInstanceIdOk() (*string, bool)`

GetInstanceIdOk returns a tuple with the InstanceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceId

`func (o *Machine) SetInstanceId(v string)`

SetInstanceId sets InstanceId field to given value.

### HasInstanceId

`func (o *Machine) HasInstanceId() bool`

HasInstanceId returns a boolean if a field has been set.

### SetInstanceIdNil

`func (o *Machine) SetInstanceIdNil(b bool)`

 SetInstanceIdNil sets the value for InstanceId to be an explicit nil

### UnsetInstanceId
`func (o *Machine) UnsetInstanceId()`

UnsetInstanceId ensures that no value is present for InstanceId, not even an explicit nil
### GetTenantId

`func (o *Machine) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *Machine) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *Machine) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *Machine) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### SetTenantIdNil

`func (o *Machine) SetTenantIdNil(b bool)`

 SetTenantIdNil sets the value for TenantId to be an explicit nil

### UnsetTenantId
`func (o *Machine) UnsetTenantId()`

UnsetTenantId ensures that no value is present for TenantId, not even an explicit nil
### GetControllerMachineId

`func (o *Machine) GetControllerMachineId() string`

GetControllerMachineId returns the ControllerMachineId field if non-nil, zero value otherwise.

### GetControllerMachineIdOk

`func (o *Machine) GetControllerMachineIdOk() (*string, bool)`

GetControllerMachineIdOk returns a tuple with the ControllerMachineId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerMachineId

`func (o *Machine) SetControllerMachineId(v string)`

SetControllerMachineId sets ControllerMachineId field to given value.

### HasControllerMachineId

`func (o *Machine) HasControllerMachineId() bool`

HasControllerMachineId returns a boolean if a field has been set.

### GetControllerMachineType

`func (o *Machine) GetControllerMachineType() string`

GetControllerMachineType returns the ControllerMachineType field if non-nil, zero value otherwise.

### GetControllerMachineTypeOk

`func (o *Machine) GetControllerMachineTypeOk() (*string, bool)`

GetControllerMachineTypeOk returns a tuple with the ControllerMachineType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerMachineType

`func (o *Machine) SetControllerMachineType(v string)`

SetControllerMachineType sets ControllerMachineType field to given value.

### HasControllerMachineType

`func (o *Machine) HasControllerMachineType() bool`

HasControllerMachineType returns a boolean if a field has been set.

### GetHwSkuDeviceType

`func (o *Machine) GetHwSkuDeviceType() string`

GetHwSkuDeviceType returns the HwSkuDeviceType field if non-nil, zero value otherwise.

### GetHwSkuDeviceTypeOk

`func (o *Machine) GetHwSkuDeviceTypeOk() (*string, bool)`

GetHwSkuDeviceTypeOk returns a tuple with the HwSkuDeviceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHwSkuDeviceType

`func (o *Machine) SetHwSkuDeviceType(v string)`

SetHwSkuDeviceType sets HwSkuDeviceType field to given value.

### HasHwSkuDeviceType

`func (o *Machine) HasHwSkuDeviceType() bool`

HasHwSkuDeviceType returns a boolean if a field has been set.

### SetHwSkuDeviceTypeNil

`func (o *Machine) SetHwSkuDeviceTypeNil(b bool)`

 SetHwSkuDeviceTypeNil sets the value for HwSkuDeviceType to be an explicit nil

### UnsetHwSkuDeviceType
`func (o *Machine) UnsetHwSkuDeviceType()`

UnsetHwSkuDeviceType ensures that no value is present for HwSkuDeviceType, not even an explicit nil
### GetVendor

`func (o *Machine) GetVendor() string`

GetVendor returns the Vendor field if non-nil, zero value otherwise.

### GetVendorOk

`func (o *Machine) GetVendorOk() (*string, bool)`

GetVendorOk returns a tuple with the Vendor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVendor

`func (o *Machine) SetVendor(v string)`

SetVendor sets Vendor field to given value.

### HasVendor

`func (o *Machine) HasVendor() bool`

HasVendor returns a boolean if a field has been set.

### GetProductName

`func (o *Machine) GetProductName() string`

GetProductName returns the ProductName field if non-nil, zero value otherwise.

### GetProductNameOk

`func (o *Machine) GetProductNameOk() (*string, bool)`

GetProductNameOk returns a tuple with the ProductName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProductName

`func (o *Machine) SetProductName(v string)`

SetProductName sets ProductName field to given value.

### HasProductName

`func (o *Machine) HasProductName() bool`

HasProductName returns a boolean if a field has been set.

### GetSerialNumber

`func (o *Machine) GetSerialNumber() string`

GetSerialNumber returns the SerialNumber field if non-nil, zero value otherwise.

### GetSerialNumberOk

`func (o *Machine) GetSerialNumberOk() (*string, bool)`

GetSerialNumberOk returns a tuple with the SerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialNumber

`func (o *Machine) SetSerialNumber(v string)`

SetSerialNumber sets SerialNumber field to given value.

### HasSerialNumber

`func (o *Machine) HasSerialNumber() bool`

HasSerialNumber returns a boolean if a field has been set.

### GetMachineCapabilities

`func (o *Machine) GetMachineCapabilities() []MachineCapability`

GetMachineCapabilities returns the MachineCapabilities field if non-nil, zero value otherwise.

### GetMachineCapabilitiesOk

`func (o *Machine) GetMachineCapabilitiesOk() (*[]MachineCapability, bool)`

GetMachineCapabilitiesOk returns a tuple with the MachineCapabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineCapabilities

`func (o *Machine) SetMachineCapabilities(v []MachineCapability)`

SetMachineCapabilities sets MachineCapabilities field to given value.

### HasMachineCapabilities

`func (o *Machine) HasMachineCapabilities() bool`

HasMachineCapabilities returns a boolean if a field has been set.

### GetMachineInterfaces

`func (o *Machine) GetMachineInterfaces() []MachineInterface`

GetMachineInterfaces returns the MachineInterfaces field if non-nil, zero value otherwise.

### GetMachineInterfacesOk

`func (o *Machine) GetMachineInterfacesOk() (*[]MachineInterface, bool)`

GetMachineInterfacesOk returns a tuple with the MachineInterfaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMachineInterfaces

`func (o *Machine) SetMachineInterfaces(v []MachineInterface)`

SetMachineInterfaces sets MachineInterfaces field to given value.

### HasMachineInterfaces

`func (o *Machine) HasMachineInterfaces() bool`

HasMachineInterfaces returns a boolean if a field has been set.

### GetMaintenanceMessage

`func (o *Machine) GetMaintenanceMessage() string`

GetMaintenanceMessage returns the MaintenanceMessage field if non-nil, zero value otherwise.

### GetMaintenanceMessageOk

`func (o *Machine) GetMaintenanceMessageOk() (*string, bool)`

GetMaintenanceMessageOk returns a tuple with the MaintenanceMessage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaintenanceMessage

`func (o *Machine) SetMaintenanceMessage(v string)`

SetMaintenanceMessage sets MaintenanceMessage field to given value.

### HasMaintenanceMessage

`func (o *Machine) HasMaintenanceMessage() bool`

HasMaintenanceMessage returns a boolean if a field has been set.

### SetMaintenanceMessageNil

`func (o *Machine) SetMaintenanceMessageNil(b bool)`

 SetMaintenanceMessageNil sets the value for MaintenanceMessage to be an explicit nil

### UnsetMaintenanceMessage
`func (o *Machine) UnsetMaintenanceMessage()`

UnsetMaintenanceMessage ensures that no value is present for MaintenanceMessage, not even an explicit nil
### GetHealth

`func (o *Machine) GetHealth() MachineHealth`

GetHealth returns the Health field if non-nil, zero value otherwise.

### GetHealthOk

`func (o *Machine) GetHealthOk() (*MachineHealth, bool)`

GetHealthOk returns a tuple with the Health field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHealth

`func (o *Machine) SetHealth(v MachineHealth)`

SetHealth sets Health field to given value.

### HasHealth

`func (o *Machine) HasHealth() bool`

HasHealth returns a boolean if a field has been set.

### GetMetadata

`func (o *Machine) GetMetadata() MachineMetadata`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *Machine) GetMetadataOk() (*MachineMetadata, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *Machine) SetMetadata(v MachineMetadata)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *Machine) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetLabels

`func (o *Machine) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *Machine) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *Machine) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *Machine) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetStatus

`func (o *Machine) GetStatus() MachineStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Machine) GetStatusOk() (*MachineStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Machine) SetStatus(v MachineStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Machine) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetIsUsableByTenant

`func (o *Machine) GetIsUsableByTenant() bool`

GetIsUsableByTenant returns the IsUsableByTenant field if non-nil, zero value otherwise.

### GetIsUsableByTenantOk

`func (o *Machine) GetIsUsableByTenantOk() (*bool, bool)`

GetIsUsableByTenantOk returns a tuple with the IsUsableByTenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsUsableByTenant

`func (o *Machine) SetIsUsableByTenant(v bool)`

SetIsUsableByTenant sets IsUsableByTenant field to given value.

### HasIsUsableByTenant

`func (o *Machine) HasIsUsableByTenant() bool`

HasIsUsableByTenant returns a boolean if a field has been set.

### GetStatusHistory

`func (o *Machine) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *Machine) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *Machine) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *Machine) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *Machine) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *Machine) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *Machine) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *Machine) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *Machine) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *Machine) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *Machine) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *Machine) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


