# CarbideAPIError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Source** | Pointer to **string** | Source of the error. Only &#39;carbide&#39; is supported | [optional] 
**Message** | Pointer to **string** | Message describing the error | [optional] 
**Data** | Pointer to **map[string]interface{}** | Additional data about the error | [optional] 

## Methods

### NewCarbideAPIError

`func NewCarbideAPIError() *CarbideAPIError`

NewCarbideAPIError instantiates a new CarbideAPIError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCarbideAPIErrorWithDefaults

`func NewCarbideAPIErrorWithDefaults() *CarbideAPIError`

NewCarbideAPIErrorWithDefaults instantiates a new CarbideAPIError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSource

`func (o *CarbideAPIError) GetSource() string`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *CarbideAPIError) GetSourceOk() (*string, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *CarbideAPIError) SetSource(v string)`

SetSource sets Source field to given value.

### HasSource

`func (o *CarbideAPIError) HasSource() bool`

HasSource returns a boolean if a field has been set.

### GetMessage

`func (o *CarbideAPIError) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *CarbideAPIError) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *CarbideAPIError) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *CarbideAPIError) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetData

`func (o *CarbideAPIError) GetData() map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *CarbideAPIError) GetDataOk() (*map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *CarbideAPIError) SetData(v map[string]interface{})`

SetData sets Data field to given value.

### HasData

`func (o *CarbideAPIError) HasData() bool`

HasData returns a boolean if a field has been set.

### SetDataNil

`func (o *CarbideAPIError) SetDataNil(b bool)`

 SetDataNil sets the value for Data to be an explicit nil

### UnsetData
`func (o *CarbideAPIError) UnsetData()`

UnsetData ensures that no value is present for Data, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


