# SiteCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**SerialConsoleHostname** | Pointer to **string** | Hostname to reach Serial Console for the Site | [optional] 
**Location** | Pointer to [**SiteLocation**](SiteLocation.md) |  | [optional] 
**Contact** | Pointer to [**SiteContact**](SiteContact.md) |  | [optional] 

## Methods

### NewSiteCreateRequest

`func NewSiteCreateRequest(name string, ) *SiteCreateRequest`

NewSiteCreateRequest instantiates a new SiteCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSiteCreateRequestWithDefaults

`func NewSiteCreateRequestWithDefaults() *SiteCreateRequest`

NewSiteCreateRequestWithDefaults instantiates a new SiteCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SiteCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SiteCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SiteCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *SiteCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *SiteCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *SiteCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *SiteCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetSerialConsoleHostname

`func (o *SiteCreateRequest) GetSerialConsoleHostname() string`

GetSerialConsoleHostname returns the SerialConsoleHostname field if non-nil, zero value otherwise.

### GetSerialConsoleHostnameOk

`func (o *SiteCreateRequest) GetSerialConsoleHostnameOk() (*string, bool)`

GetSerialConsoleHostnameOk returns a tuple with the SerialConsoleHostname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSerialConsoleHostname

`func (o *SiteCreateRequest) SetSerialConsoleHostname(v string)`

SetSerialConsoleHostname sets SerialConsoleHostname field to given value.

### HasSerialConsoleHostname

`func (o *SiteCreateRequest) HasSerialConsoleHostname() bool`

HasSerialConsoleHostname returns a boolean if a field has been set.

### GetLocation

`func (o *SiteCreateRequest) GetLocation() SiteLocation`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *SiteCreateRequest) GetLocationOk() (*SiteLocation, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *SiteCreateRequest) SetLocation(v SiteLocation)`

SetLocation sets Location field to given value.

### HasLocation

`func (o *SiteCreateRequest) HasLocation() bool`

HasLocation returns a boolean if a field has been set.

### GetContact

`func (o *SiteCreateRequest) GetContact() SiteContact`

GetContact returns the Contact field if non-nil, zero value otherwise.

### GetContactOk

`func (o *SiteCreateRequest) GetContactOk() (*SiteContact, bool)`

GetContactOk returns a tuple with the Contact field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContact

`func (o *SiteCreateRequest) SetContact(v SiteContact)`

SetContact sets Contact field to given value.

### HasContact

`func (o *SiteCreateRequest) HasContact() bool`

HasContact returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


