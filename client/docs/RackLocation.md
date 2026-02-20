# RackLocation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Region** | Pointer to **string** | Region where the rack is located | [optional] 
**Datacenter** | Pointer to **string** | Datacenter where the rack is located | [optional] 
**Room** | Pointer to **string** | Room within the datacenter | [optional] 
**Position** | Pointer to **string** | Position of the rack within the room | [optional] 

## Methods

### NewRackLocation

`func NewRackLocation() *RackLocation`

NewRackLocation instantiates a new RackLocation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRackLocationWithDefaults

`func NewRackLocationWithDefaults() *RackLocation`

NewRackLocationWithDefaults instantiates a new RackLocation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRegion

`func (o *RackLocation) GetRegion() string`

GetRegion returns the Region field if non-nil, zero value otherwise.

### GetRegionOk

`func (o *RackLocation) GetRegionOk() (*string, bool)`

GetRegionOk returns a tuple with the Region field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegion

`func (o *RackLocation) SetRegion(v string)`

SetRegion sets Region field to given value.

### HasRegion

`func (o *RackLocation) HasRegion() bool`

HasRegion returns a boolean if a field has been set.

### GetDatacenter

`func (o *RackLocation) GetDatacenter() string`

GetDatacenter returns the Datacenter field if non-nil, zero value otherwise.

### GetDatacenterOk

`func (o *RackLocation) GetDatacenterOk() (*string, bool)`

GetDatacenterOk returns a tuple with the Datacenter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDatacenter

`func (o *RackLocation) SetDatacenter(v string)`

SetDatacenter sets Datacenter field to given value.

### HasDatacenter

`func (o *RackLocation) HasDatacenter() bool`

HasDatacenter returns a boolean if a field has been set.

### GetRoom

`func (o *RackLocation) GetRoom() string`

GetRoom returns the Room field if non-nil, zero value otherwise.

### GetRoomOk

`func (o *RackLocation) GetRoomOk() (*string, bool)`

GetRoomOk returns a tuple with the Room field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoom

`func (o *RackLocation) SetRoom(v string)`

SetRoom sets Room field to given value.

### HasRoom

`func (o *RackLocation) HasRoom() bool`

HasRoom returns a boolean if a field has been set.

### GetPosition

`func (o *RackLocation) GetPosition() string`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *RackLocation) GetPositionOk() (*string, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *RackLocation) SetPosition(v string)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *RackLocation) HasPosition() bool`

HasPosition returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


