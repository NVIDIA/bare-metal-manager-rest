# AuditEntry

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Unique identifier | [optional] [readonly] 
**Endpoint** | Pointer to **string** | API endpoint | [optional] 
**QueryParams** | Pointer to **string** | Query parameters | [optional] 
**Method** | Pointer to **string** | HTTP method | [optional] 
**Body** | Pointer to **string** | HTTP body in JSON format | [optional] 
**StatusCode** | Pointer to **int32** | HTTP response status code | [optional] 
**StatusMessage** | Pointer to **string** | HTTP response status message | [optional] 
**ClientIP** | Pointer to **string** | Client IP address | [optional] 
**UserID** | Pointer to **NullableString** | User ID that executed the API call | [optional] 
**User** | Pointer to [**User**](User.md) | User that executed the API call | [optional] 
**OrgName** | Pointer to **string** | Organization name | [optional] 
**ExtraData** | Pointer to **map[string]interface{}** | Extra data in JSON format | [optional] 
**Timestamp** | Pointer to **time.Time** | API execution time | [optional] 
**DurationMs** | Pointer to **int32** | API execution duration in milliseconds | [optional] 
**ApiVersion** | Pointer to **string** | API version | [optional] 

## Methods

### NewAuditEntry

`func NewAuditEntry() *AuditEntry`

NewAuditEntry instantiates a new AuditEntry object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuditEntryWithDefaults

`func NewAuditEntryWithDefaults() *AuditEntry`

NewAuditEntryWithDefaults instantiates a new AuditEntry object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AuditEntry) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AuditEntry) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AuditEntry) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AuditEntry) HasId() bool`

HasId returns a boolean if a field has been set.

### GetEndpoint

`func (o *AuditEntry) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *AuditEntry) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *AuditEntry) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.

### HasEndpoint

`func (o *AuditEntry) HasEndpoint() bool`

HasEndpoint returns a boolean if a field has been set.

### GetQueryParams

`func (o *AuditEntry) GetQueryParams() string`

GetQueryParams returns the QueryParams field if non-nil, zero value otherwise.

### GetQueryParamsOk

`func (o *AuditEntry) GetQueryParamsOk() (*string, bool)`

GetQueryParamsOk returns a tuple with the QueryParams field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueryParams

`func (o *AuditEntry) SetQueryParams(v string)`

SetQueryParams sets QueryParams field to given value.

### HasQueryParams

`func (o *AuditEntry) HasQueryParams() bool`

HasQueryParams returns a boolean if a field has been set.

### GetMethod

`func (o *AuditEntry) GetMethod() string`

GetMethod returns the Method field if non-nil, zero value otherwise.

### GetMethodOk

`func (o *AuditEntry) GetMethodOk() (*string, bool)`

GetMethodOk returns a tuple with the Method field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMethod

`func (o *AuditEntry) SetMethod(v string)`

SetMethod sets Method field to given value.

### HasMethod

`func (o *AuditEntry) HasMethod() bool`

HasMethod returns a boolean if a field has been set.

### GetBody

`func (o *AuditEntry) GetBody() string`

GetBody returns the Body field if non-nil, zero value otherwise.

### GetBodyOk

`func (o *AuditEntry) GetBodyOk() (*string, bool)`

GetBodyOk returns a tuple with the Body field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBody

`func (o *AuditEntry) SetBody(v string)`

SetBody sets Body field to given value.

### HasBody

`func (o *AuditEntry) HasBody() bool`

HasBody returns a boolean if a field has been set.

### GetStatusCode

`func (o *AuditEntry) GetStatusCode() int32`

GetStatusCode returns the StatusCode field if non-nil, zero value otherwise.

### GetStatusCodeOk

`func (o *AuditEntry) GetStatusCodeOk() (*int32, bool)`

GetStatusCodeOk returns a tuple with the StatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusCode

`func (o *AuditEntry) SetStatusCode(v int32)`

SetStatusCode sets StatusCode field to given value.

### HasStatusCode

`func (o *AuditEntry) HasStatusCode() bool`

HasStatusCode returns a boolean if a field has been set.

### GetStatusMessage

`func (o *AuditEntry) GetStatusMessage() string`

GetStatusMessage returns the StatusMessage field if non-nil, zero value otherwise.

### GetStatusMessageOk

`func (o *AuditEntry) GetStatusMessageOk() (*string, bool)`

GetStatusMessageOk returns a tuple with the StatusMessage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusMessage

`func (o *AuditEntry) SetStatusMessage(v string)`

SetStatusMessage sets StatusMessage field to given value.

### HasStatusMessage

`func (o *AuditEntry) HasStatusMessage() bool`

HasStatusMessage returns a boolean if a field has been set.

### GetClientIP

`func (o *AuditEntry) GetClientIP() string`

GetClientIP returns the ClientIP field if non-nil, zero value otherwise.

### GetClientIPOk

`func (o *AuditEntry) GetClientIPOk() (*string, bool)`

GetClientIPOk returns a tuple with the ClientIP field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientIP

`func (o *AuditEntry) SetClientIP(v string)`

SetClientIP sets ClientIP field to given value.

### HasClientIP

`func (o *AuditEntry) HasClientIP() bool`

HasClientIP returns a boolean if a field has been set.

### GetUserID

`func (o *AuditEntry) GetUserID() string`

GetUserID returns the UserID field if non-nil, zero value otherwise.

### GetUserIDOk

`func (o *AuditEntry) GetUserIDOk() (*string, bool)`

GetUserIDOk returns a tuple with the UserID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserID

`func (o *AuditEntry) SetUserID(v string)`

SetUserID sets UserID field to given value.

### HasUserID

`func (o *AuditEntry) HasUserID() bool`

HasUserID returns a boolean if a field has been set.

### SetUserIDNil

`func (o *AuditEntry) SetUserIDNil(b bool)`

 SetUserIDNil sets the value for UserID to be an explicit nil

### UnsetUserID
`func (o *AuditEntry) UnsetUserID()`

UnsetUserID ensures that no value is present for UserID, not even an explicit nil
### GetUser

`func (o *AuditEntry) GetUser() User`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *AuditEntry) GetUserOk() (*User, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *AuditEntry) SetUser(v User)`

SetUser sets User field to given value.

### HasUser

`func (o *AuditEntry) HasUser() bool`

HasUser returns a boolean if a field has been set.

### GetOrgName

`func (o *AuditEntry) GetOrgName() string`

GetOrgName returns the OrgName field if non-nil, zero value otherwise.

### GetOrgNameOk

`func (o *AuditEntry) GetOrgNameOk() (*string, bool)`

GetOrgNameOk returns a tuple with the OrgName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrgName

`func (o *AuditEntry) SetOrgName(v string)`

SetOrgName sets OrgName field to given value.

### HasOrgName

`func (o *AuditEntry) HasOrgName() bool`

HasOrgName returns a boolean if a field has been set.

### GetExtraData

`func (o *AuditEntry) GetExtraData() map[string]interface{}`

GetExtraData returns the ExtraData field if non-nil, zero value otherwise.

### GetExtraDataOk

`func (o *AuditEntry) GetExtraDataOk() (*map[string]interface{}, bool)`

GetExtraDataOk returns a tuple with the ExtraData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtraData

`func (o *AuditEntry) SetExtraData(v map[string]interface{})`

SetExtraData sets ExtraData field to given value.

### HasExtraData

`func (o *AuditEntry) HasExtraData() bool`

HasExtraData returns a boolean if a field has been set.

### GetTimestamp

`func (o *AuditEntry) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *AuditEntry) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *AuditEntry) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *AuditEntry) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetDurationMs

`func (o *AuditEntry) GetDurationMs() int32`

GetDurationMs returns the DurationMs field if non-nil, zero value otherwise.

### GetDurationMsOk

`func (o *AuditEntry) GetDurationMsOk() (*int32, bool)`

GetDurationMsOk returns a tuple with the DurationMs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDurationMs

`func (o *AuditEntry) SetDurationMs(v int32)`

SetDurationMs sets DurationMs field to given value.

### HasDurationMs

`func (o *AuditEntry) HasDurationMs() bool`

HasDurationMs returns a boolean if a field has been set.

### GetApiVersion

`func (o *AuditEntry) GetApiVersion() string`

GetApiVersion returns the ApiVersion field if non-nil, zero value otherwise.

### GetApiVersionOk

`func (o *AuditEntry) GetApiVersionOk() (*string, bool)`

GetApiVersionOk returns a tuple with the ApiVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiVersion

`func (o *AuditEntry) SetApiVersion(v string)`

SetApiVersion sets ApiVersion field to given value.

### HasApiVersion

`func (o *AuditEntry) HasApiVersion() bool`

HasApiVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


