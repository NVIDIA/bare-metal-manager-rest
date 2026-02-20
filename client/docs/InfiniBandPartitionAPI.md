# \InfiniBandPartitionAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateInfinibandPartition**](InfiniBandPartitionAPI.md#CreateInfinibandPartition) | **Post** /v2/org/{org}/carbide/infiniband-partition | Create InfiniBand Partition
[**DeleteInfinibandPartition**](InfiniBandPartitionAPI.md#DeleteInfinibandPartition) | **Delete** /v2/org/{org}/carbide/infiniband-partition/{infiniBandPartitionId} | Delete InfiniBand Partition
[**GetAllInfinibandPartition**](InfiniBandPartitionAPI.md#GetAllInfinibandPartition) | **Get** /v2/org/{org}/carbide/infiniband-partition | Retrieve all InfiniBand Partitions
[**GetInfinibandPartition**](InfiniBandPartitionAPI.md#GetInfinibandPartition) | **Get** /v2/org/{org}/carbide/infiniband-partition/{infiniBandPartitionId} | Retrieve InfiniBand Partition
[**UpdateInfinibandPartition**](InfiniBandPartitionAPI.md#UpdateInfinibandPartition) | **Patch** /v2/org/{org}/carbide/infiniband-partition/{infiniBandPartitionId} | Update InfiniBand Partition



## CreateInfinibandPartition

> InfiniBandPartition CreateInfinibandPartition(ctx, org).InfiniBandPartitionCreateRequest(infiniBandPartitionCreateRequest).Execute()

Create InfiniBand Partition



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
	infiniBandPartitionCreateRequest := *openapiclient.NewInfiniBandPartitionCreateRequest("Name_example", "SiteId_example") // InfiniBandPartitionCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InfiniBandPartitionAPI.CreateInfinibandPartition(context.Background(), org).InfiniBandPartitionCreateRequest(infiniBandPartitionCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InfiniBandPartitionAPI.CreateInfinibandPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateInfinibandPartition`: InfiniBandPartition
	fmt.Fprintf(os.Stdout, "Response from `InfiniBandPartitionAPI.CreateInfinibandPartition`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateInfinibandPartitionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **infiniBandPartitionCreateRequest** | [**InfiniBandPartitionCreateRequest**](InfiniBandPartitionCreateRequest.md) |  | 

### Return type

[**InfiniBandPartition**](InfiniBandPartition.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteInfinibandPartition

> DeleteInfinibandPartition(ctx, org, infiniBandPartitionId).Execute()

Delete InfiniBand Partition



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
	infiniBandPartitionId := "infiniBandPartitionId_example" // string | ID of the InfiniBand Partition

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.InfiniBandPartitionAPI.DeleteInfinibandPartition(context.Background(), org, infiniBandPartitionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InfiniBandPartitionAPI.DeleteInfinibandPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**infiniBandPartitionId** | **string** | ID of the InfiniBand Partition | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteInfinibandPartitionRequest struct via the builder pattern


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


## GetAllInfinibandPartition

> []InfiniBandPartition GetAllInfinibandPartition(ctx, org).SiteId(siteId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all InfiniBand Partitions



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Partitions by Site (optional)
	status := "status_example" // string | Filter Partitions by Status (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, description and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InfiniBandPartitionAPI.GetAllInfinibandPartition(context.Background(), org).SiteId(siteId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InfiniBandPartitionAPI.GetAllInfinibandPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllInfinibandPartition`: []InfiniBandPartition
	fmt.Fprintf(os.Stdout, "Response from `InfiniBandPartitionAPI.GetAllInfinibandPartition`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllInfinibandPartitionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | Filter Partitions by Site | 
 **status** | **string** | Filter Partitions by Status | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, description and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]InfiniBandPartition**](InfiniBandPartition.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInfinibandPartition

> InfiniBandPartition GetInfinibandPartition(ctx, org, infiniBandPartitionId).IncludeRelation(includeRelation).Execute()

Retrieve InfiniBand Partition



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
	infiniBandPartitionId := "infiniBandPartitionId_example" // string | ID of the InfiniBand Partition
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InfiniBandPartitionAPI.GetInfinibandPartition(context.Background(), org, infiniBandPartitionId).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InfiniBandPartitionAPI.GetInfinibandPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInfinibandPartition`: InfiniBandPartition
	fmt.Fprintf(os.Stdout, "Response from `InfiniBandPartitionAPI.GetInfinibandPartition`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**infiniBandPartitionId** | **string** | ID of the InfiniBand Partition | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInfinibandPartitionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**InfiniBandPartition**](InfiniBandPartition.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateInfinibandPartition

> InfiniBandPartition UpdateInfinibandPartition(ctx, org, infiniBandPartitionId).InfiniBandPartitionUpdateRequest(infiniBandPartitionUpdateRequest).Execute()

Update InfiniBand Partition



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
	infiniBandPartitionId := "infiniBandPartitionId_example" // string | ID of the InfiniBand Partition
	infiniBandPartitionUpdateRequest := *openapiclient.NewInfiniBandPartitionUpdateRequest("Name_example") // InfiniBandPartitionUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InfiniBandPartitionAPI.UpdateInfinibandPartition(context.Background(), org, infiniBandPartitionId).InfiniBandPartitionUpdateRequest(infiniBandPartitionUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InfiniBandPartitionAPI.UpdateInfinibandPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateInfinibandPartition`: InfiniBandPartition
	fmt.Fprintf(os.Stdout, "Response from `InfiniBandPartitionAPI.UpdateInfinibandPartition`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**infiniBandPartitionId** | **string** | ID of the InfiniBand Partition | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInfinibandPartitionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **infiniBandPartitionUpdateRequest** | [**InfiniBandPartitionUpdateRequest**](InfiniBandPartitionUpdateRequest.md) |  | 

### Return type

[**InfiniBandPartition**](InfiniBandPartition.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

