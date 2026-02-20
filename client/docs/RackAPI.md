# \RackAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAllRack**](RackAPI.md#GetAllRack) | **Get** /v2/org/{org}/carbide/rack | Retrieve all Racks
[**GetRack**](RackAPI.md#GetRack) | **Get** /v2/org/{org}/carbide/rack/{id} | Retrieve a Rack



## GetAllRack

> []Rack GetAllRack(ctx, org).SiteId(siteId).IncludeComponents(includeComponents).Name(name).Manufacturer(manufacturer).Model(model).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Racks



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Site to retrieve Racks from
	org := "org_example" // string | Name of the Org
	includeComponents := true // bool | Include rack components in response (optional)
	name := "name_example" // string | Filter by rack name (optional)
	manufacturer := "manufacturer_example" // string | Filter by manufacturer (optional)
	model := "model_example" // string | Filter by model (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RackAPI.GetAllRack(context.Background(), org).SiteId(siteId).IncludeComponents(includeComponents).Name(name).Manufacturer(manufacturer).Model(model).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RackAPI.GetAllRack``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllRack`: []Rack
	fmt.Fprintf(os.Stdout, "Response from `RackAPI.GetAllRack`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllRackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **string** | ID of the Site to retrieve Racks from | 

 **includeComponents** | **bool** | Include rack components in response | 
 **name** | **string** | Filter by rack name | 
 **manufacturer** | **string** | Filter by manufacturer | 
 **model** | **string** | Filter by model | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]Rack**](Rack.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRack

> Rack GetRack(ctx, org, id).SiteId(siteId).IncludeComponents(includeComponents).Execute()

Retrieve a Rack



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Site
	org := "org_example" // string | Name of the Org
	id := "id_example" // string | ID of the Rack
	includeComponents := true // bool | Include rack components in response (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RackAPI.GetRack(context.Background(), org, id).SiteId(siteId).IncludeComponents(includeComponents).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RackAPI.GetRack``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetRack`: Rack
	fmt.Fprintf(os.Stdout, "Response from `RackAPI.GetRack`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**id** | **string** | ID of the Rack | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **string** | ID of the Site | 


 **includeComponents** | **bool** | Include rack components in response | 

### Return type

[**Rack**](Rack.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

