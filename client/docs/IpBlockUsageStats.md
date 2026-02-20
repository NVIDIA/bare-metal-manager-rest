# IpBlockUsageStats

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AvailableIPs** | Pointer to **int64** | Number of IP addresses available in the block | [optional] [readonly] 
**AcquiredIPs** | Pointer to **int64** | Number of individual IP addresses acquired from the block | [optional] [readonly] 
**AvailablePrefixes** | Pointer to **[]string** | Example prefixes available to acquire | [optional] 
**AvailableSmallestPrefixes** | Pointer to **int64** | Number of smallest prefixes available to acquire | [optional] [readonly] 
**AcquiredPrefixes** | Pointer to **int64** | Number of prefixes acquired from this block | [optional] 

## Methods

### NewIpBlockUsageStats

`func NewIpBlockUsageStats() *IpBlockUsageStats`

NewIpBlockUsageStats instantiates a new IpBlockUsageStats object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIpBlockUsageStatsWithDefaults

`func NewIpBlockUsageStatsWithDefaults() *IpBlockUsageStats`

NewIpBlockUsageStatsWithDefaults instantiates a new IpBlockUsageStats object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvailableIPs

`func (o *IpBlockUsageStats) GetAvailableIPs() int64`

GetAvailableIPs returns the AvailableIPs field if non-nil, zero value otherwise.

### GetAvailableIPsOk

`func (o *IpBlockUsageStats) GetAvailableIPsOk() (*int64, bool)`

GetAvailableIPsOk returns a tuple with the AvailableIPs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableIPs

`func (o *IpBlockUsageStats) SetAvailableIPs(v int64)`

SetAvailableIPs sets AvailableIPs field to given value.

### HasAvailableIPs

`func (o *IpBlockUsageStats) HasAvailableIPs() bool`

HasAvailableIPs returns a boolean if a field has been set.

### GetAcquiredIPs

`func (o *IpBlockUsageStats) GetAcquiredIPs() int64`

GetAcquiredIPs returns the AcquiredIPs field if non-nil, zero value otherwise.

### GetAcquiredIPsOk

`func (o *IpBlockUsageStats) GetAcquiredIPsOk() (*int64, bool)`

GetAcquiredIPsOk returns a tuple with the AcquiredIPs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAcquiredIPs

`func (o *IpBlockUsageStats) SetAcquiredIPs(v int64)`

SetAcquiredIPs sets AcquiredIPs field to given value.

### HasAcquiredIPs

`func (o *IpBlockUsageStats) HasAcquiredIPs() bool`

HasAcquiredIPs returns a boolean if a field has been set.

### GetAvailablePrefixes

`func (o *IpBlockUsageStats) GetAvailablePrefixes() []string`

GetAvailablePrefixes returns the AvailablePrefixes field if non-nil, zero value otherwise.

### GetAvailablePrefixesOk

`func (o *IpBlockUsageStats) GetAvailablePrefixesOk() (*[]string, bool)`

GetAvailablePrefixesOk returns a tuple with the AvailablePrefixes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailablePrefixes

`func (o *IpBlockUsageStats) SetAvailablePrefixes(v []string)`

SetAvailablePrefixes sets AvailablePrefixes field to given value.

### HasAvailablePrefixes

`func (o *IpBlockUsageStats) HasAvailablePrefixes() bool`

HasAvailablePrefixes returns a boolean if a field has been set.

### GetAvailableSmallestPrefixes

`func (o *IpBlockUsageStats) GetAvailableSmallestPrefixes() int64`

GetAvailableSmallestPrefixes returns the AvailableSmallestPrefixes field if non-nil, zero value otherwise.

### GetAvailableSmallestPrefixesOk

`func (o *IpBlockUsageStats) GetAvailableSmallestPrefixesOk() (*int64, bool)`

GetAvailableSmallestPrefixesOk returns a tuple with the AvailableSmallestPrefixes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableSmallestPrefixes

`func (o *IpBlockUsageStats) SetAvailableSmallestPrefixes(v int64)`

SetAvailableSmallestPrefixes sets AvailableSmallestPrefixes field to given value.

### HasAvailableSmallestPrefixes

`func (o *IpBlockUsageStats) HasAvailableSmallestPrefixes() bool`

HasAvailableSmallestPrefixes returns a boolean if a field has been set.

### GetAcquiredPrefixes

`func (o *IpBlockUsageStats) GetAcquiredPrefixes() int64`

GetAcquiredPrefixes returns the AcquiredPrefixes field if non-nil, zero value otherwise.

### GetAcquiredPrefixesOk

`func (o *IpBlockUsageStats) GetAcquiredPrefixesOk() (*int64, bool)`

GetAcquiredPrefixesOk returns a tuple with the AcquiredPrefixes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAcquiredPrefixes

`func (o *IpBlockUsageStats) SetAcquiredPrefixes(v int64)`

SetAcquiredPrefixes sets AcquiredPrefixes field to given value.

### HasAcquiredPrefixes

`func (o *IpBlockUsageStats) HasAcquiredPrefixes() bool`

HasAcquiredPrefixes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


