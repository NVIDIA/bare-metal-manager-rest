# \DPUExtensionServiceAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateDpuExtensionService**](DPUExtensionServiceAPI.md#CreateDpuExtensionService) | **Post** /v2/org/{org}/carbide/dpu-extension-service | Create DPU Extension Service
[**DeleteDpuExtensionService**](DPUExtensionServiceAPI.md#DeleteDpuExtensionService) | **Delete** /v2/org/{org}/carbide/dpu-extension-service/{dpuExtensionServiceId} | Delete DPU Extension Service
[**DeleteDpuExtensionServiceVersion**](DPUExtensionServiceAPI.md#DeleteDpuExtensionServiceVersion) | **Delete** /v2/org/{org}/carbide/dpu-extension-service/{dpuExtensionServiceId}/version/{version} | Delete DPU Extension Service Version
[**GetAllDpuExtensionService**](DPUExtensionServiceAPI.md#GetAllDpuExtensionService) | **Get** /v2/org/{org}/carbide/dpu-extension-service | Retrieve all DPU Extension Services
[**GetDpuExtensionService**](DPUExtensionServiceAPI.md#GetDpuExtensionService) | **Get** /v2/org/{org}/carbide/dpu-extension-service/{dpuExtensionServiceId} | Retrieve DPU Extension Service
[**GetDpuExtensionServiceVersion**](DPUExtensionServiceAPI.md#GetDpuExtensionServiceVersion) | **Get** /v2/org/{org}/carbide/dpu-extension-service/{dpuExtensionServiceId}/version/{version} | Retrieve DPU Extension Service Version
[**UpdateDpuExtensionService**](DPUExtensionServiceAPI.md#UpdateDpuExtensionService) | **Patch** /v2/org/{org}/carbide/dpu-extension-service/{dpuExtensionServiceId} | Update DPU Extension Service



## CreateDpuExtensionService

> DpuExtensionService CreateDpuExtensionService(ctx, org).DpuExtensionServiceCreateRequest(dpuExtensionServiceCreateRequest).Execute()

Create DPU Extension Service



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
	dpuExtensionServiceCreateRequest := *openapiclient.NewDpuExtensionServiceCreateRequest("Name_example", "ServiceType_example", "SiteId_example", "Data_example") // DpuExtensionServiceCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DPUExtensionServiceAPI.CreateDpuExtensionService(context.Background(), org).DpuExtensionServiceCreateRequest(dpuExtensionServiceCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DPUExtensionServiceAPI.CreateDpuExtensionService``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateDpuExtensionService`: DpuExtensionService
	fmt.Fprintf(os.Stdout, "Response from `DPUExtensionServiceAPI.CreateDpuExtensionService`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateDpuExtensionServiceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **dpuExtensionServiceCreateRequest** | [**DpuExtensionServiceCreateRequest**](DpuExtensionServiceCreateRequest.md) |  | 

### Return type

[**DpuExtensionService**](DpuExtensionService.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteDpuExtensionService

> DeleteDpuExtensionService(ctx, org, dpuExtensionServiceId).Execute()

Delete DPU Extension Service



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
	dpuExtensionServiceId := "dpuExtensionServiceId_example" // string | ID of the DPU Extension Service

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DPUExtensionServiceAPI.DeleteDpuExtensionService(context.Background(), org, dpuExtensionServiceId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DPUExtensionServiceAPI.DeleteDpuExtensionService``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**dpuExtensionServiceId** | **string** | ID of the DPU Extension Service | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDpuExtensionServiceRequest struct via the builder pattern


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


## DeleteDpuExtensionServiceVersion

> DeleteDpuExtensionServiceVersion(ctx, org, dpuExtensionServiceId, version).Execute()

Delete DPU Extension Service Version



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
	dpuExtensionServiceId := "dpuExtensionServiceId_example" // string | ID of the DPU Extension Service
	version := "version_example" // string | Version of the DPU Extension Service

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DPUExtensionServiceAPI.DeleteDpuExtensionServiceVersion(context.Background(), org, dpuExtensionServiceId, version).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DPUExtensionServiceAPI.DeleteDpuExtensionServiceVersion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**dpuExtensionServiceId** | **string** | ID of the DPU Extension Service | 
**version** | **string** | Version of the DPU Extension Service | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDpuExtensionServiceVersionRequest struct via the builder pattern


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


## GetAllDpuExtensionService

> []DpuExtensionService GetAllDpuExtensionService(ctx, org).SiteId(siteId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all DPU Extension Services



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
	siteId := "siteId_example" // string | Filter DPU Extension Services by Site ID (optional)
	status := "status_example" // string | Status filter for the DPU Extension Services (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, description and status (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DPUExtensionServiceAPI.GetAllDpuExtensionService(context.Background(), org).SiteId(siteId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DPUExtensionServiceAPI.GetAllDpuExtensionService``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllDpuExtensionService`: []DpuExtensionService
	fmt.Fprintf(os.Stdout, "Response from `DPUExtensionServiceAPI.GetAllDpuExtensionService`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllDpuExtensionServiceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | Filter DPU Extension Services by Site ID | 
 **status** | **string** | Status filter for the DPU Extension Services | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, description and status | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]DpuExtensionService**](DpuExtensionService.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDpuExtensionService

> DpuExtensionService GetDpuExtensionService(ctx, org, dpuExtensionServiceId).Execute()

Retrieve DPU Extension Service



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
	dpuExtensionServiceId := "dpuExtensionServiceId_example" // string | ID of the DPU Extension Service

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DPUExtensionServiceAPI.GetDpuExtensionService(context.Background(), org, dpuExtensionServiceId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DPUExtensionServiceAPI.GetDpuExtensionService``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDpuExtensionService`: DpuExtensionService
	fmt.Fprintf(os.Stdout, "Response from `DPUExtensionServiceAPI.GetDpuExtensionService`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**dpuExtensionServiceId** | **string** | ID of the DPU Extension Service | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDpuExtensionServiceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DpuExtensionService**](DpuExtensionService.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDpuExtensionServiceVersion

> DpuExtensionServiceVersionInfo GetDpuExtensionServiceVersion(ctx, org, dpuExtensionServiceId, version).Execute()

Retrieve DPU Extension Service Version



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
	dpuExtensionServiceId := "dpuExtensionServiceId_example" // string | ID of the DPU Extension Service
	version := "version_example" // string | Version of the DPU Extension Service

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DPUExtensionServiceAPI.GetDpuExtensionServiceVersion(context.Background(), org, dpuExtensionServiceId, version).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DPUExtensionServiceAPI.GetDpuExtensionServiceVersion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDpuExtensionServiceVersion`: DpuExtensionServiceVersionInfo
	fmt.Fprintf(os.Stdout, "Response from `DPUExtensionServiceAPI.GetDpuExtensionServiceVersion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**dpuExtensionServiceId** | **string** | ID of the DPU Extension Service | 
**version** | **string** | Version of the DPU Extension Service | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDpuExtensionServiceVersionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[**DpuExtensionServiceVersionInfo**](DpuExtensionServiceVersionInfo.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateDpuExtensionService

> DpuExtensionService UpdateDpuExtensionService(ctx, org, dpuExtensionServiceId).DpuExtensionServiceUpdateRequest(dpuExtensionServiceUpdateRequest).Execute()

Update DPU Extension Service



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
	dpuExtensionServiceId := "dpuExtensionServiceId_example" // string | ID of the DPU Extension Service
	dpuExtensionServiceUpdateRequest := *openapiclient.NewDpuExtensionServiceUpdateRequest() // DpuExtensionServiceUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DPUExtensionServiceAPI.UpdateDpuExtensionService(context.Background(), org, dpuExtensionServiceId).DpuExtensionServiceUpdateRequest(dpuExtensionServiceUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DPUExtensionServiceAPI.UpdateDpuExtensionService``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateDpuExtensionService`: DpuExtensionService
	fmt.Fprintf(os.Stdout, "Response from `DPUExtensionServiceAPI.UpdateDpuExtensionService`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**dpuExtensionServiceId** | **string** | ID of the DPU Extension Service | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateDpuExtensionServiceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **dpuExtensionServiceUpdateRequest** | [**DpuExtensionServiceUpdateRequest**](DpuExtensionServiceUpdateRequest.md) |  | 

### Return type

[**DpuExtensionService**](DpuExtensionService.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

