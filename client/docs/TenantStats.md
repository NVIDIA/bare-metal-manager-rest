# TenantStats

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Instance** | Pointer to [**InstanceCountByStatus**](InstanceCountByStatus.md) |  | [optional] 
**Vpc** | Pointer to [**VpcCountByStatus**](VpcCountByStatus.md) |  | [optional] 
**Subnet** | Pointer to [**SubnetCountByStatus**](SubnetCountByStatus.md) |  | [optional] 

## Methods

### NewTenantStats

`func NewTenantStats() *TenantStats`

NewTenantStats instantiates a new TenantStats object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTenantStatsWithDefaults

`func NewTenantStatsWithDefaults() *TenantStats`

NewTenantStatsWithDefaults instantiates a new TenantStats object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInstance

`func (o *TenantStats) GetInstance() InstanceCountByStatus`

GetInstance returns the Instance field if non-nil, zero value otherwise.

### GetInstanceOk

`func (o *TenantStats) GetInstanceOk() (*InstanceCountByStatus, bool)`

GetInstanceOk returns a tuple with the Instance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstance

`func (o *TenantStats) SetInstance(v InstanceCountByStatus)`

SetInstance sets Instance field to given value.

### HasInstance

`func (o *TenantStats) HasInstance() bool`

HasInstance returns a boolean if a field has been set.

### GetVpc

`func (o *TenantStats) GetVpc() VpcCountByStatus`

GetVpc returns the Vpc field if non-nil, zero value otherwise.

### GetVpcOk

`func (o *TenantStats) GetVpcOk() (*VpcCountByStatus, bool)`

GetVpcOk returns a tuple with the Vpc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVpc

`func (o *TenantStats) SetVpc(v VpcCountByStatus)`

SetVpc sets Vpc field to given value.

### HasVpc

`func (o *TenantStats) HasVpc() bool`

HasVpc returns a boolean if a field has been set.

### GetSubnet

`func (o *TenantStats) GetSubnet() SubnetCountByStatus`

GetSubnet returns the Subnet field if non-nil, zero value otherwise.

### GetSubnetOk

`func (o *TenantStats) GetSubnetOk() (*SubnetCountByStatus, bool)`

GetSubnetOk returns a tuple with the Subnet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubnet

`func (o *TenantStats) SetSubnet(v SubnetCountByStatus)`

SetSubnet sets Subnet field to given value.

### HasSubnet

`func (o *TenantStats) HasSubnet() bool`

HasSubnet returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


