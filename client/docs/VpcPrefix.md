# VpcPrefix

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | Pointer to **string** | Name of the VPC Prefix | [optional] 
**SiteId** | Pointer to **string** | ID of the Site the VPC Prefix belongs to | [optional] 
**VpcId** | Pointer to **string** | ID of the VPC the VPC Prefix belongs to | [optional] 
**TenantId** | Pointer to **string** | ID of the Tenant the VPC Prefix belongs to | [optional] 
**IpBlockId** | Pointer to **NullableString** | ID of the IP Block that contains the prefix of the VPC Prefix | [optional] 
**Prefix** | Pointer to **NullableString** | The network prefix including prefix length in CIDR notation | [optional] 
**PrefixLength** | Pointer to **int32** | Length of the prefix. Max value depends on prefix length of parent IP Block. | [optional] 
**Status** | Pointer to [**VpcPrefixStatus**](VpcPrefixStatus.md) | Status of the VPC Prefix | [optional] [readonly] 
**StatusHistory** | Pointer to [**[]StatusDetail**](StatusDetail.md) | Details of 20 most recent status changes | [optional] 
**Created** | Pointer to **time.Time** | Date and time when the VPC Prefix was created | [optional] [readonly] 
**Updated** | Pointer to **time.Time** | Date and time when the VPC Prefix was updated | [optional] [readonly] 

## Methods

### NewVpcPrefix

`func NewVpcPrefix() *VpcPrefix`

NewVpcPrefix instantiates a new VpcPrefix object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVpcPrefixWithDefaults

`func NewVpcPrefixWithDefaults() *VpcPrefix`

NewVpcPrefixWithDefaults instantiates a new VpcPrefix object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *VpcPrefix) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *VpcPrefix) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *VpcPrefix) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *VpcPrefix) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *VpcPrefix) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *VpcPrefix) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *VpcPrefix) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *VpcPrefix) HasName() bool`

HasName returns a boolean if a field has been set.

### GetSiteId

`func (o *VpcPrefix) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *VpcPrefix) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *VpcPrefix) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.

### HasSiteId

`func (o *VpcPrefix) HasSiteId() bool`

HasSiteId returns a boolean if a field has been set.

### GetVpcId

`func (o *VpcPrefix) GetVpcId() string`

GetVpcId returns the VpcId field if non-nil, zero value otherwise.

### GetVpcIdOk

`func (o *VpcPrefix) GetVpcIdOk() (*string, bool)`

GetVpcIdOk returns a tuple with the VpcId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpcId

`func (o *VpcPrefix) SetVpcId(v string)`

SetVpcId sets VpcId field to given value.

### HasVpcId

`func (o *VpcPrefix) HasVpcId() bool`

HasVpcId returns a boolean if a field has been set.

### GetTenantId

`func (o *VpcPrefix) GetTenantId() string`

GetTenantId returns the TenantId field if non-nil, zero value otherwise.

### GetTenantIdOk

`func (o *VpcPrefix) GetTenantIdOk() (*string, bool)`

GetTenantIdOk returns a tuple with the TenantId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantId

`func (o *VpcPrefix) SetTenantId(v string)`

SetTenantId sets TenantId field to given value.

### HasTenantId

`func (o *VpcPrefix) HasTenantId() bool`

HasTenantId returns a boolean if a field has been set.

### GetIpBlockId

`func (o *VpcPrefix) GetIpBlockId() string`

GetIpBlockId returns the IpBlockId field if non-nil, zero value otherwise.

### GetIpBlockIdOk

`func (o *VpcPrefix) GetIpBlockIdOk() (*string, bool)`

GetIpBlockIdOk returns a tuple with the IpBlockId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpBlockId

`func (o *VpcPrefix) SetIpBlockId(v string)`

SetIpBlockId sets IpBlockId field to given value.

### HasIpBlockId

`func (o *VpcPrefix) HasIpBlockId() bool`

HasIpBlockId returns a boolean if a field has been set.

### SetIpBlockIdNil

`func (o *VpcPrefix) SetIpBlockIdNil(b bool)`

 SetIpBlockIdNil sets the value for IpBlockId to be an explicit nil

### UnsetIpBlockId
`func (o *VpcPrefix) UnsetIpBlockId()`

UnsetIpBlockId ensures that no value is present for IpBlockId, not even an explicit nil
### GetPrefix

`func (o *VpcPrefix) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *VpcPrefix) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *VpcPrefix) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *VpcPrefix) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.

### SetPrefixNil

`func (o *VpcPrefix) SetPrefixNil(b bool)`

 SetPrefixNil sets the value for Prefix to be an explicit nil

### UnsetPrefix
`func (o *VpcPrefix) UnsetPrefix()`

UnsetPrefix ensures that no value is present for Prefix, not even an explicit nil
### GetPrefixLength

`func (o *VpcPrefix) GetPrefixLength() int32`

GetPrefixLength returns the PrefixLength field if non-nil, zero value otherwise.

### GetPrefixLengthOk

`func (o *VpcPrefix) GetPrefixLengthOk() (*int32, bool)`

GetPrefixLengthOk returns a tuple with the PrefixLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefixLength

`func (o *VpcPrefix) SetPrefixLength(v int32)`

SetPrefixLength sets PrefixLength field to given value.

### HasPrefixLength

`func (o *VpcPrefix) HasPrefixLength() bool`

HasPrefixLength returns a boolean if a field has been set.

### GetStatus

`func (o *VpcPrefix) GetStatus() VpcPrefixStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *VpcPrefix) GetStatusOk() (*VpcPrefixStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *VpcPrefix) SetStatus(v VpcPrefixStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *VpcPrefix) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetStatusHistory

`func (o *VpcPrefix) GetStatusHistory() []StatusDetail`

GetStatusHistory returns the StatusHistory field if non-nil, zero value otherwise.

### GetStatusHistoryOk

`func (o *VpcPrefix) GetStatusHistoryOk() (*[]StatusDetail, bool)`

GetStatusHistoryOk returns a tuple with the StatusHistory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusHistory

`func (o *VpcPrefix) SetStatusHistory(v []StatusDetail)`

SetStatusHistory sets StatusHistory field to given value.

### HasStatusHistory

`func (o *VpcPrefix) HasStatusHistory() bool`

HasStatusHistory returns a boolean if a field has been set.

### GetCreated

`func (o *VpcPrefix) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *VpcPrefix) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *VpcPrefix) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *VpcPrefix) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetUpdated

`func (o *VpcPrefix) GetUpdated() time.Time`

GetUpdated returns the Updated field if non-nil, zero value otherwise.

### GetUpdatedOk

`func (o *VpcPrefix) GetUpdatedOk() (*time.Time, bool)`

GetUpdatedOk returns a tuple with the Updated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdated

`func (o *VpcPrefix) SetUpdated(v time.Time)`

SetUpdated sets Updated field to given value.

### HasUpdated

`func (o *VpcPrefix) HasUpdated() bool`

HasUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


