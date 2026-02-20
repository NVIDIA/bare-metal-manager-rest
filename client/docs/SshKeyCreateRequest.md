# SshKeyCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name cannot match that name an existing SSH Key | 
**PublicKey** | **string** | Must be an SSH key of type: RSA, ECDSA or ED25519 | 
**SshKeyGroupId** | Pointer to **NullableString** | ID of the SSH Key Group this key should be attached to | [optional] 

## Methods

### NewSshKeyCreateRequest

`func NewSshKeyCreateRequest(name string, publicKey string, ) *SshKeyCreateRequest`

NewSshKeyCreateRequest instantiates a new SshKeyCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSshKeyCreateRequestWithDefaults

`func NewSshKeyCreateRequestWithDefaults() *SshKeyCreateRequest`

NewSshKeyCreateRequestWithDefaults instantiates a new SshKeyCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SshKeyCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SshKeyCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SshKeyCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetPublicKey

`func (o *SshKeyCreateRequest) GetPublicKey() string`

GetPublicKey returns the PublicKey field if non-nil, zero value otherwise.

### GetPublicKeyOk

`func (o *SshKeyCreateRequest) GetPublicKeyOk() (*string, bool)`

GetPublicKeyOk returns a tuple with the PublicKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublicKey

`func (o *SshKeyCreateRequest) SetPublicKey(v string)`

SetPublicKey sets PublicKey field to given value.


### GetSshKeyGroupId

`func (o *SshKeyCreateRequest) GetSshKeyGroupId() string`

GetSshKeyGroupId returns the SshKeyGroupId field if non-nil, zero value otherwise.

### GetSshKeyGroupIdOk

`func (o *SshKeyCreateRequest) GetSshKeyGroupIdOk() (*string, bool)`

GetSshKeyGroupIdOk returns a tuple with the SshKeyGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyGroupId

`func (o *SshKeyCreateRequest) SetSshKeyGroupId(v string)`

SetSshKeyGroupId sets SshKeyGroupId field to given value.

### HasSshKeyGroupId

`func (o *SshKeyCreateRequest) HasSshKeyGroupId() bool`

HasSshKeyGroupId returns a boolean if a field has been set.

### SetSshKeyGroupIdNil

`func (o *SshKeyCreateRequest) SetSshKeyGroupIdNil(b bool)`

 SetSshKeyGroupIdNil sets the value for SshKeyGroupId to be an explicit nil

### UnsetSshKeyGroupId
`func (o *SshKeyCreateRequest) UnsetSshKeyGroupId()`

UnsetSshKeyGroupId ensures that no value is present for SshKeyGroupId, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


