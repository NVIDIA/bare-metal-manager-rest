# MachineHealthIssue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Category** | Pointer to **string** | Category of the issue | [optional] 
**Summary** | Pointer to **string** | Short summary describing the issue | [optional] 
**Details** | Pointer to **NullableString** | Details about the issue helpful for diagnosis | [optional] 

## Methods

### NewMachineHealthIssue

`func NewMachineHealthIssue() *MachineHealthIssue`

NewMachineHealthIssue instantiates a new MachineHealthIssue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMachineHealthIssueWithDefaults

`func NewMachineHealthIssueWithDefaults() *MachineHealthIssue`

NewMachineHealthIssueWithDefaults instantiates a new MachineHealthIssue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCategory

`func (o *MachineHealthIssue) GetCategory() string`

GetCategory returns the Category field if non-nil, zero value otherwise.

### GetCategoryOk

`func (o *MachineHealthIssue) GetCategoryOk() (*string, bool)`

GetCategoryOk returns a tuple with the Category field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategory

`func (o *MachineHealthIssue) SetCategory(v string)`

SetCategory sets Category field to given value.

### HasCategory

`func (o *MachineHealthIssue) HasCategory() bool`

HasCategory returns a boolean if a field has been set.

### GetSummary

`func (o *MachineHealthIssue) GetSummary() string`

GetSummary returns the Summary field if non-nil, zero value otherwise.

### GetSummaryOk

`func (o *MachineHealthIssue) GetSummaryOk() (*string, bool)`

GetSummaryOk returns a tuple with the Summary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSummary

`func (o *MachineHealthIssue) SetSummary(v string)`

SetSummary sets Summary field to given value.

### HasSummary

`func (o *MachineHealthIssue) HasSummary() bool`

HasSummary returns a boolean if a field has been set.

### GetDetails

`func (o *MachineHealthIssue) GetDetails() string`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *MachineHealthIssue) GetDetailsOk() (*string, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetails

`func (o *MachineHealthIssue) SetDetails(v string)`

SetDetails sets Details field to given value.

### HasDetails

`func (o *MachineHealthIssue) HasDetails() bool`

HasDetails returns a boolean if a field has been set.

### SetDetailsNil

`func (o *MachineHealthIssue) SetDetailsNil(b bool)`

 SetDetailsNil sets the value for Details to be an explicit nil

### UnsetDetails
`func (o *MachineHealthIssue) UnsetDetails()`

UnsetDetails ensures that no value is present for Details, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


