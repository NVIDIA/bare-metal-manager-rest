# Sku

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier for the SKU | [optional] [readonly] 
**SiteId** | Pointer to **string** | ID of the Site this SKU belongs to | [optional] [readonly] 
**DeviceType** | Pointer to **NullableString** | Optional device type identifier (e.g. \&quot;gpu\&quot;, \&quot;cpu\&quot;, \&quot;storage\&quot;) | [optional] [readonly] 
**AssociatedMachineIds** | Pointer to **[]string** | List of machine IDs associated with this SKU | [optional] 
**Components** | Pointer to [**SkuComponents**](SkuComponents.md) | Hardware components of this SKU | [optional] 
**Created** | Pointer to **time.Time** | ISO 8601 datetime when the SKU was created | [optional] [readonly] 
**Updated** | Pointer to **time.Time** | ISO 8601 datetime when the SKU was last updated | [optional] [readonly] 

## Methods

### NewSku

`func NewSku() *Sku`

NewSku instantiates a new Sku object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSkuWithDefaults

`func NewSkuWithDefaults() *Sku`

NewSkuWithDefaults instantiates a new Sku object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Sku) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Sku) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Sku) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Sku) HasId() bool`

HasId returns a boolean if a field has been set.

### GetSiteId

`func (o *Sku) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *Sku) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *Sku) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *Sku) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetDeviceType

`func (o *Sku) GetDeviceType() string`

GetDeviceType returns the DeviceType field if non-nil, zero value otherwise.

### GetDeviceTypeOk

`func (o *Sku) GetDeviceTypeOk() (*string, bool)`

GetDeviceTypeOk returns a tuple with the DeviceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceType

`func (o *Sku) SetDeviceType(v string)`

SetDeviceType sets DeviceType field to given value.

### HasDeviceType

`func (o *Sku) HasDeviceType() bool`

HasDeviceType returns a boolean if a field has been set.

### SetDeviceTypeNil

`func (o *Sku) SetDeviceTypeNil(b bool)`

 SetDeviceTypeNil sets the value for DeviceType to be an explicit nil

### UnsetDeviceType
`func (o *Sku) UnsetDeviceType()`

UnsetDeviceType ensures that no value is present for DeviceType, not even an explicit nil
### GetAssociatedMachineIds

`func (o *Sku) GetAssociatedMachineIds() []string`

GetAssociatedMachineIds returns the AssociatedMachineIds field if non-nil, zero value otherwise.

### GetAssociatedMachineIdsOk

`func (o *Sku) GetAssociatedMachineIdsOk() (*[]string, bool)`

GetAssociatedMachineIdsOk returns a tuple with the AssociatedMachineIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAssociatedMachineIds

`func (o *Sku) SetAssociatedMachineIds(v []string)`

SetAssociatedMachineIds sets AssociatedMachineIds field to given value.

### HasAssociatedMachineIds

`func (o *Sku) HasAssociatedMachineIds() bool`

HasAssociatedMachineIds returns a boolean if a field has been set.

### GetComponents

`func (o *Sku) GetComponents() SkuComponents`

GetComponents returns the Components field if non-nil, zero value otherwise.

### GetComponentsOk

`func (o *Sku) GetComponentsOk() (*SkuComponents, bool)`

GetComponentsOk returns a tuple with the Components field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComponents

`func (o *Sku) SetComponents(v SkuComponents)`

SetComponents sets Components field to given value.

### HasComponents

`func (o *Sku) HasComponents() bool`

HasComponents returns a boolean if a field has been set.

### GetCreated

`func (o *Sku) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *Sku) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *Sku) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *Sku) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *Sku) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *Sku) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *Sku) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *Sku) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


