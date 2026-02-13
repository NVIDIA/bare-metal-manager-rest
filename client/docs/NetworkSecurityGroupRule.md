# NetworkSecurityGroupRule

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **NullableString** |  | [optional] 
**Direction** | **string** |  | 
**SourcePortRange** | Pointer to **NullableString** |  | [optional] 
**DestinationPortRange** | Pointer to **NullableString** |  | [optional] 
**Protocol** | **string** |  | 
**Action** | **string** |  | 
**Priority** | Pointer to **int32** |  | [optional] 
**SourcePrefix** | **string** |  | 
**DestinationPrefix** | **string** |  | 

## Methods

### NewNetworkSecurityGroupRule

`func NewNetworkSecurityGroupRule(direction string, protocol string, action string, sourcePrefix string, destinationPrefix string, ) *NetworkSecurityGroupRule`

NewNetworkSecurityGroupRule instantiates a new NetworkSecurityGroupRule object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNetworkSecurityGroupRuleWithDefaults

`func NewNetworkSecurityGroupRuleWithDefaults() *NetworkSecurityGroupRule`

NewNetworkSecurityGroupRuleWithDefaults instantiates a new NetworkSecurityGroupRule object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *NetworkSecurityGroupRule) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NetworkSecurityGroupRule) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NetworkSecurityGroupRule) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *NetworkSecurityGroupRule) HasName() bool`

HasName returns a boolean if a field has been set.

### SetNameNil

`func (o *NetworkSecurityGroupRule) SetNameNil(b bool)`

 SetNameNil sets the value for Name to be an explicit nil

### UnsetName
`func (o *NetworkSecurityGroupRule) UnsetName()`

UnsetName ensures that no value is present for Name, not even an explicit nil
### GetDirection

`func (o *NetworkSecurityGroupRule) GetDirection() string`

GetDirection returns the Direction field if non-nil, zero value otherwise.

### GetDirectionOk

`func (o *NetworkSecurityGroupRule) GetDirectionOk() (*string, bool)`

GetDirectionOk returns a tuple with the Direction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDirection

`func (o *NetworkSecurityGroupRule) SetDirection(v string)`

SetDirection sets Direction field to given value.


### GetSourcePortRange

`func (o *NetworkSecurityGroupRule) GetSourcePortRange() string`

GetSourcePortRange returns the SourcePortRange field if non-nil, zero value otherwise.

### GetSourcePortRangeOk

`func (o *NetworkSecurityGroupRule) GetSourcePortRangeOk() (*string, bool)`

GetSourcePortRangeOk returns a tuple with the SourcePortRange field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourcePortRange

`func (o *NetworkSecurityGroupRule) SetSourcePortRange(v string)`

SetSourcePortRange sets SourcePortRange field to given value.

### HasSourcePortRange

`func (o *NetworkSecurityGroupRule) HasSourcePortRange() bool`

HasSourcePortRange returns a boolean if a field has been set.

### SetSourcePortRangeNil

`func (o *NetworkSecurityGroupRule) SetSourcePortRangeNil(b bool)`

 SetSourcePortRangeNil sets the value for SourcePortRange to be an explicit nil

### UnsetSourcePortRange
`func (o *NetworkSecurityGroupRule) UnsetSourcePortRange()`

UnsetSourcePortRange ensures that no value is present for SourcePortRange, not even an explicit nil
### GetDestinationPortRange

`func (o *NetworkSecurityGroupRule) GetDestinationPortRange() string`

GetDestinationPortRange returns the DestinationPortRange field if non-nil, zero value otherwise.

### GetDestinationPortRangeOk

`func (o *NetworkSecurityGroupRule) GetDestinationPortRangeOk() (*string, bool)`

GetDestinationPortRangeOk returns a tuple with the DestinationPortRange field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationPortRange

`func (o *NetworkSecurityGroupRule) SetDestinationPortRange(v string)`

SetDestinationPortRange sets DestinationPortRange field to given value.

### HasDestinationPortRange

`func (o *NetworkSecurityGroupRule) HasDestinationPortRange() bool`

HasDestinationPortRange returns a boolean if a field has been set.

### SetDestinationPortRangeNil

`func (o *NetworkSecurityGroupRule) SetDestinationPortRangeNil(b bool)`

 SetDestinationPortRangeNil sets the value for DestinationPortRange to be an explicit nil

### UnsetDestinationPortRange
`func (o *NetworkSecurityGroupRule) UnsetDestinationPortRange()`

UnsetDestinationPortRange ensures that no value is present for DestinationPortRange, not even an explicit nil
### GetProtocol

`func (o *NetworkSecurityGroupRule) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *NetworkSecurityGroupRule) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *NetworkSecurityGroupRule) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.


### GetAction

`func (o *NetworkSecurityGroupRule) GetAction() string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *NetworkSecurityGroupRule) GetActionOk() (*string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *NetworkSecurityGroupRule) SetAction(v string)`

SetAction sets Action field to given value.


### GetPriority

`func (o *NetworkSecurityGroupRule) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *NetworkSecurityGroupRule) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *NetworkSecurityGroupRule) SetPriority(v int32)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *NetworkSecurityGroupRule) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetSourcePrefix

`func (o *NetworkSecurityGroupRule) GetSourcePrefix() string`

GetSourcePrefix returns the SourcePrefix field if non-nil, zero value otherwise.

### GetSourcePrefixOk

`func (o *NetworkSecurityGroupRule) GetSourcePrefixOk() (*string, bool)`

GetSourcePrefixOk returns a tuple with the SourcePrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourcePrefix

`func (o *NetworkSecurityGroupRule) SetSourcePrefix(v string)`

SetSourcePrefix sets SourcePrefix field to given value.


### GetDestinationPrefix

`func (o *NetworkSecurityGroupRule) GetDestinationPrefix() string`

GetDestinationPrefix returns the DestinationPrefix field if non-nil, zero value otherwise.

### GetDestinationPrefixOk

`func (o *NetworkSecurityGroupRule) GetDestinationPrefixOk() (*string, bool)`

GetDestinationPrefixOk returns a tuple with the DestinationPrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationPrefix

`func (o *NetworkSecurityGroupRule) SetDestinationPrefix(v string)`

SetDestinationPrefix sets DestinationPrefix field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


