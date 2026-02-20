# DpuExtensionServiceCredentials

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RegistryUrl** | Pointer to **string** | URL for the registry the credential should be used for | [optional] 
**Username** | Pointer to **NullableString** | Username for the registry. Must be specified if registry URL is specified | [optional] 
**Password** | Pointer to **NullableString** | Password for the registry. Must be specified if registry URL is specified | [optional] 

## Methods

### NewDpuExtensionServiceCredentials

`func NewDpuExtensionServiceCredentials() *DpuExtensionServiceCredentials`

NewDpuExtensionServiceCredentials instantiates a new DpuExtensionServiceCredentials object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDpuExtensionServiceCredentialsWithDefaults

`func NewDpuExtensionServiceCredentialsWithDefaults() *DpuExtensionServiceCredentials`

NewDpuExtensionServiceCredentialsWithDefaults instantiates a new DpuExtensionServiceCredentials object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRegistryUrl

`func (o *DpuExtensionServiceCredentials) GetRegistryUrl() string`

GetRegistryUrl returns the RegistryUrl field if non-nil, zero value otherwise.

### GetRegistryUrlOk

`func (o *DpuExtensionServiceCredentials) GetRegistryUrlOk() (*string, bool)`

GetRegistryUrlOk returns a tuple with the RegistryUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistryUrl

`func (o *DpuExtensionServiceCredentials) SetRegistryUrl(v string)`

SetRegistryUrl sets RegistryUrl field to given value.

### HasRegistryUrl

`func (o *DpuExtensionServiceCredentials) HasRegistryUrl() bool`

HasRegistryUrl returns a boolean if a field has been set.

### GetUsername

`func (o *DpuExtensionServiceCredentials) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *DpuExtensionServiceCredentials) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *DpuExtensionServiceCredentials) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *DpuExtensionServiceCredentials) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### SetUsernameNil

`func (o *DpuExtensionServiceCredentials) SetUsernameNil(b bool)`

 SetUsernameNil sets the value for Username to be an explicit nil

### UnsetUsername
`func (o *DpuExtensionServiceCredentials) UnsetUsername()`

UnsetUsername ensures that no value is present for Username, not even an explicit nil
### GetPassword

`func (o *DpuExtensionServiceCredentials) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *DpuExtensionServiceCredentials) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *DpuExtensionServiceCredentials) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *DpuExtensionServiceCredentials) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### SetPasswordNil

`func (o *DpuExtensionServiceCredentials) SetPasswordNil(b bool)`

 SetPasswordNil sets the value for Password to be an explicit nil

### UnsetPassword
`func (o *DpuExtensionServiceCredentials) UnsetPassword()`

UnsetPassword ensures that no value is present for Password, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


