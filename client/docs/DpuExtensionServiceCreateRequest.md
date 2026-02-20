# DpuExtensionServiceCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name for the DPU Extension Service. Must be unique for a given Tenant | 
**Description** | Pointer to **NullableString** | Optional description for the DPU Extension Service | [optional] 
**ServiceType** | **string** | Type of the DPU Extension Service | 
**SiteId** | **string** | ID for the Site the DPU Extension Service belongs to | 
**Data** | **string** | Deployment spec for the DPU Extension Service | 
**Credentials** | Pointer to [**DpuExtensionServiceCredentials**](DpuExtensionServiceCredentials.md) | Credentials to download resources specified in DPU Extension Service data | [optional] 

## Methods

### NewDpuExtensionServiceCreateRequest

`func NewDpuExtensionServiceCreateRequest(name string, serviceType string, siteId string, data string, ) *DpuExtensionServiceCreateRequest`

NewDpuExtensionServiceCreateRequest instantiates a new DpuExtensionServiceCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDpuExtensionServiceCreateRequestWithDefaults

`func NewDpuExtensionServiceCreateRequestWithDefaults() *DpuExtensionServiceCreateRequest`

NewDpuExtensionServiceCreateRequestWithDefaults instantiates a new DpuExtensionServiceCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *DpuExtensionServiceCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DpuExtensionServiceCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DpuExtensionServiceCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *DpuExtensionServiceCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *DpuExtensionServiceCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *DpuExtensionServiceCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *DpuExtensionServiceCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *DpuExtensionServiceCreateRequest) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *DpuExtensionServiceCreateRequest) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetServiceType

`func (o *DpuExtensionServiceCreateRequest) GetServiceType() string`

GetServiceType returns the ServiceType field if non-nil, zero value otherwise.

### GetServiceTypeOk

`func (o *DpuExtensionServiceCreateRequest) GetServiceTypeOk() (*string, bool)`

GetServiceTypeOk returns a tuple with the ServiceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceType

`func (o *DpuExtensionServiceCreateRequest) SetServiceType(v string)`

SetServiceType sets ServiceType field to given value.


### GetSiteId

`func (o *DpuExtensionServiceCreateRequest) GetSiteId() string`

GetSiteId returns the SiteId field if non-nil, zero value otherwise.

### GetSiteIdOk

`func (o *DpuExtensionServiceCreateRequest) GetSiteIdOk() (*string, bool)`

GetSiteIdOk returns a tuple with the SiteId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteId

`func (o *DpuExtensionServiceCreateRequest) SetSiteId(v string)`

SetSiteId sets SiteId field to given value.


### GetData

`func (o *DpuExtensionServiceCreateRequest) GetData() string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DpuExtensionServiceCreateRequest) GetDataOk() (*string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DpuExtensionServiceCreateRequest) SetData(v string)`

SetData sets Data field to given value.


### GetCredentials

`func (o *DpuExtensionServiceCreateRequest) GetCredentials() DpuExtensionServiceCredentials`

GetCredentials returns the Credentials field if non-nil, zero value otherwise.

### GetCredentialsOk

`func (o *DpuExtensionServiceCreateRequest) GetCredentialsOk() (*DpuExtensionServiceCredentials, bool)`

GetCredentialsOk returns a tuple with the Credentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentials

`func (o *DpuExtensionServiceCreateRequest) SetCredentials(v DpuExtensionServiceCredentials)`

SetCredentials sets Credentials field to given value.

### HasCredentials

`func (o *DpuExtensionServiceCreateRequest) HasCredentials() bool`

HasCredentials returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


