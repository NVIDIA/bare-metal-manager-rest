# \OperatingSystemAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateOperatingSystem**](OperatingSystemAPI.md#CreateOperatingSystem) | **Post** /v2/org/{org}/carbide/operating-system | Create Operating System
[**DeleteOperatingSystem**](OperatingSystemAPI.md#DeleteOperatingSystem) | **Delete** /v2/org/{org}/carbide/operating-system/{operatingSystemId} | Delete Operating System
[**GetAllOperatingSystem**](OperatingSystemAPI.md#GetAllOperatingSystem) | **Get** /v2/org/{org}/carbide/operating-system | Retrieve all Operating Systems
[**GetOperatingSystem**](OperatingSystemAPI.md#GetOperatingSystem) | **Get** /v2/org/{org}/carbide/operating-system/{operatingSystemId} | Retrieve Operating System
[**UpdateOperatingSystem**](OperatingSystemAPI.md#UpdateOperatingSystem) | **Patch** /v2/org/{org}/carbide/operating-system/{operatingSystemId} | Update Operating System



## CreateOperatingSystem

> OperatingSystem CreateOperatingSystem(ctx, org).OperatingSystemCreateRequest(operatingSystemCreateRequest).Execute()

Create Operating System



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
	org := "org_example" // string | Name of the Org
	operatingSystemCreateRequest := *openapiclient.NewOperatingSystemCreateRequest("Name_example") // OperatingSystemCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OperatingSystemAPI.CreateOperatingSystem(context.Background(), org).OperatingSystemCreateRequest(operatingSystemCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OperatingSystemAPI.CreateOperatingSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateOperatingSystem`: OperatingSystem
	fmt.Fprintf(os.Stdout, "Response from `OperatingSystemAPI.CreateOperatingSystem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateOperatingSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **operatingSystemCreateRequest** | [**OperatingSystemCreateRequest**](OperatingSystemCreateRequest.md) |  | 

### Return type

[**OperatingSystem**](OperatingSystem.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOperatingSystem

> DeleteOperatingSystem(ctx, org, operatingSystemId).Execute()

Delete Operating System



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
	org := "org_example" // string | Name of the Org
	operatingSystemId := "operatingSystemId_example" // string | ID of the Operating System

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OperatingSystemAPI.DeleteOperatingSystem(context.Background(), org, operatingSystemId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OperatingSystemAPI.DeleteOperatingSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**operatingSystemId** | **string** | ID of the Operating System | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOperatingSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllOperatingSystem

> []OperatingSystem GetAllOperatingSystem(ctx, org).SiteId(siteId).Type_(type_).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Operating Systems



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
	org := "org_example" // string | Name of the Org
	siteId := "siteId_example" // string | Filter Operating Systems by Site ID.  Can be specified multiple times to filter on more than one ID. (optional)
	type_ := "type__example" // string | Filter Operating Systems by Type (optional)
	status := "status_example" // string | Filter Operating Systems by Status.  Can be specified multiple times to filter on more than one status. (optional)
	query := "query_example" // string | Provide query to search for matches. Input will be matched against name, description and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OperatingSystemAPI.GetAllOperatingSystem(context.Background(), org).SiteId(siteId).Type_(type_).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OperatingSystemAPI.GetAllOperatingSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllOperatingSystem`: []OperatingSystem
	fmt.Fprintf(os.Stdout, "Response from `OperatingSystemAPI.GetAllOperatingSystem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllOperatingSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | Filter Operating Systems by Site ID.  Can be specified multiple times to filter on more than one ID. | 
 **type_** | **string** | Filter Operating Systems by Type | 
 **status** | **string** | Filter Operating Systems by Status.  Can be specified multiple times to filter on more than one status. | 
 **query** | **string** | Provide query to search for matches. Input will be matched against name, description and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]OperatingSystem**](OperatingSystem.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOperatingSystem

> OperatingSystem GetOperatingSystem(ctx, org, operatingSystemId).IncludeRelation(includeRelation).Execute()

Retrieve Operating System



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
	org := "org_example" // string | Name of the Org
	operatingSystemId := "operatingSystemId_example" // string | ID of the Operating System
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OperatingSystemAPI.GetOperatingSystem(context.Background(), org, operatingSystemId).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OperatingSystemAPI.GetOperatingSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOperatingSystem`: OperatingSystem
	fmt.Fprintf(os.Stdout, "Response from `OperatingSystemAPI.GetOperatingSystem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**operatingSystemId** | **string** | ID of the Operating System | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOperatingSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**OperatingSystem**](OperatingSystem.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateOperatingSystem

> OperatingSystem UpdateOperatingSystem(ctx, org, operatingSystemId).OperatingSystemUpdateRequest(operatingSystemUpdateRequest).Execute()

Update Operating System



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
	org := "org_example" // string | Name of the Org
	operatingSystemId := "operatingSystemId_example" // string | ID of the Operating System
	operatingSystemUpdateRequest := *openapiclient.NewOperatingSystemUpdateRequest() // OperatingSystemUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OperatingSystemAPI.UpdateOperatingSystem(context.Background(), org, operatingSystemId).OperatingSystemUpdateRequest(operatingSystemUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OperatingSystemAPI.UpdateOperatingSystem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateOperatingSystem`: OperatingSystem
	fmt.Fprintf(os.Stdout, "Response from `OperatingSystemAPI.UpdateOperatingSystem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**operatingSystemId** | **string** | ID of the Operating System | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateOperatingSystemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **operatingSystemUpdateRequest** | [**OperatingSystemUpdateRequest**](OperatingSystemUpdateRequest.md) |  | 

### Return type

[**OperatingSystem**](OperatingSystem.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

