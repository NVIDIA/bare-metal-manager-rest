# \NVLinkLogicalPartitionAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateNvlinkLogicalPartition**](NVLinkLogicalPartitionAPI.md#CreateNvlinkLogicalPartition) | **Post** /v2/org/{org}/carbide/nvlink-logical-partition | Create NVLink Logical Partition
[**DeleteNvlinkLogicalPartition**](NVLinkLogicalPartitionAPI.md#DeleteNvlinkLogicalPartition) | **Delete** /v2/org/{org}/carbide/nvlink-logical-partition/{nvLinkLogicalPartitionId} | Delete NVLink Logical Partition
[**GetAllNvlinkInterface**](NVLinkLogicalPartitionAPI.md#GetAllNvlinkInterface) | **Get** /v2/org/{org}/carbide/nvlink-interface | Retrieve all NVLink Interfaces
[**GetAllNvlinkLogicalPartition**](NVLinkLogicalPartitionAPI.md#GetAllNvlinkLogicalPartition) | **Get** /v2/org/{org}/carbide/nvlink-logical-partition | Retrieve all NVLink Logical Partitions
[**GetNvlinkLogicalPartition**](NVLinkLogicalPartitionAPI.md#GetNvlinkLogicalPartition) | **Get** /v2/org/{org}/carbide/nvlink-logical-partition/{nvLinkLogicalPartitionId} | Retrieve NVLink Logical Partition
[**UpdateNvlinkLogicalPartition**](NVLinkLogicalPartitionAPI.md#UpdateNvlinkLogicalPartition) | **Patch** /v2/org/{org}/carbide/nvlink-logical-partition/{nvLinkLogicalPartitionId} | Update NVLink Logical Partition



## CreateNvlinkLogicalPartition

> NVLinkLogicalPartition CreateNvlinkLogicalPartition(ctx, org).NVLinkLogicalPartitionCreateRequest(nVLinkLogicalPartitionCreateRequest).Execute()

Create NVLink Logical Partition



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
	org := "org_example" // string | Name of the NGC Org
	nVLinkLogicalPartitionCreateRequest := *openapiclient.NewNVLinkLogicalPartitionCreateRequest("Name_example", "SiteId_example") // NVLinkLogicalPartitionCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NVLinkLogicalPartitionAPI.CreateNvlinkLogicalPartition(context.Background(), org).NVLinkLogicalPartitionCreateRequest(nVLinkLogicalPartitionCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NVLinkLogicalPartitionAPI.CreateNvlinkLogicalPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateNvlinkLogicalPartition`: NVLinkLogicalPartition
	fmt.Fprintf(os.Stdout, "Response from `NVLinkLogicalPartitionAPI.CreateNvlinkLogicalPartition`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the NGC Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateNvlinkLogicalPartitionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **nVLinkLogicalPartitionCreateRequest** | [**NVLinkLogicalPartitionCreateRequest**](NVLinkLogicalPartitionCreateRequest.md) |  | 

### Return type

[**NVLinkLogicalPartition**](NVLinkLogicalPartition.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteNvlinkLogicalPartition

> DeleteNvlinkLogicalPartition(ctx, org, nvLinkLogicalPartitionId).Execute()

Delete NVLink Logical Partition



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
	org := "org_example" // string | Name of the NGC Org
	nvLinkLogicalPartitionId := "nvLinkLogicalPartitionId_example" // string | ID of the NVLink Logical Partition

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.NVLinkLogicalPartitionAPI.DeleteNvlinkLogicalPartition(context.Background(), org, nvLinkLogicalPartitionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NVLinkLogicalPartitionAPI.DeleteNvlinkLogicalPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the NGC Org | 
**nvLinkLogicalPartitionId** | **string** | ID of the NVLink Logical Partition | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteNvlinkLogicalPartitionRequest struct via the builder pattern


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


## GetAllNvlinkInterface

> []NVLinkInterface GetAllNvlinkInterface(ctx, org).Status(status).SiteId(siteId).InstanceId(instanceId).NvLinkLogicalPartitionId(nvLinkLogicalPartitionId).NvLinkDomainId(nvLinkDomainId).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all NVLink Interfaces



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
	org := "org_example" // string | Name of the NGC Org
	status := "status_example" // string | Filter NVLink Interfaces by Status (optional)
	siteId := "siteId_example" // string | Filter NVLink Interfaces by Site ID.  Can be specified multiple times to filter on more than one ID. (optional)
	instanceId := "instanceId_example" // string | Filter NVLink Interfaces by Instance ID.  Can be specified multiple times to filter on more than one ID. (optional)
	nvLinkLogicalPartitionId := "nvLinkLogicalPartitionId_example" // string | Filter NVLink Interfaces by NVLink Logical Partition ID.  Can be specified multiple times to filter on more than one ID. (optional)
	nvLinkDomainId := "nvLinkDomainId_example" // string | Filter NVLink Interfaces by NVLink Domain ID.  Can be specified multiple times to filter on more than one ID. (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NVLinkLogicalPartitionAPI.GetAllNvlinkInterface(context.Background(), org).Status(status).SiteId(siteId).InstanceId(instanceId).NvLinkLogicalPartitionId(nvLinkLogicalPartitionId).NvLinkDomainId(nvLinkDomainId).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NVLinkLogicalPartitionAPI.GetAllNvlinkInterface``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllNvlinkInterface`: []NVLinkInterface
	fmt.Fprintf(os.Stdout, "Response from `NVLinkLogicalPartitionAPI.GetAllNvlinkInterface`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the NGC Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllNvlinkInterfaceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **status** | **string** | Filter NVLink Interfaces by Status | 
 **siteId** | **string** | Filter NVLink Interfaces by Site ID.  Can be specified multiple times to filter on more than one ID. | 
 **instanceId** | **string** | Filter NVLink Interfaces by Instance ID.  Can be specified multiple times to filter on more than one ID. | 
 **nvLinkLogicalPartitionId** | **string** | Filter NVLink Interfaces by NVLink Logical Partition ID.  Can be specified multiple times to filter on more than one ID. | 
 **nvLinkDomainId** | **string** | Filter NVLink Interfaces by NVLink Domain ID.  Can be specified multiple times to filter on more than one ID. | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]NVLinkInterface**](NVLinkInterface.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllNvlinkLogicalPartition

> []NVLinkLogicalPartition GetAllNvlinkLogicalPartition(ctx, org).SiteId(siteId).Status(status).Query(query).IncludeInterfaces(includeInterfaces).IncludeStats(includeStats).IncludeVpcs(includeVpcs).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all NVLink Logical Partitions



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
	org := "org_example" // string | Name of the NGC Org
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter NVLink Logical Partitions by Site (optional)
	status := "status_example" // string | Filter NVLink Logical Partitions by Status (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, description and status fields (optional)
	includeInterfaces := true // bool | Include NVLink Interfaces in response. (optional)
	includeStats := true // bool | Include NVLink Logical Partition Stats in response. (optional)
	includeVpcs := true // bool | Include VPCs in response. (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NVLinkLogicalPartitionAPI.GetAllNvlinkLogicalPartition(context.Background(), org).SiteId(siteId).Status(status).Query(query).IncludeInterfaces(includeInterfaces).IncludeStats(includeStats).IncludeVpcs(includeVpcs).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NVLinkLogicalPartitionAPI.GetAllNvlinkLogicalPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllNvlinkLogicalPartition`: []NVLinkLogicalPartition
	fmt.Fprintf(os.Stdout, "Response from `NVLinkLogicalPartitionAPI.GetAllNvlinkLogicalPartition`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the NGC Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllNvlinkLogicalPartitionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | Filter NVLink Logical Partitions by Site | 
 **status** | **string** | Filter NVLink Logical Partitions by Status | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, description and status fields | 
 **includeInterfaces** | **bool** | Include NVLink Interfaces in response. | 
 **includeStats** | **bool** | Include NVLink Logical Partition Stats in response. | 
 **includeVpcs** | **bool** | Include VPCs in response. | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]NVLinkLogicalPartition**](NVLinkLogicalPartition.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNvlinkLogicalPartition

> NVLinkLogicalPartition GetNvlinkLogicalPartition(ctx, org, nvLinkLogicalPartitionId).IncludeInterfaces(includeInterfaces).IncludeStats(includeStats).IncludeVpcs(includeVpcs).IncludeRelation(includeRelation).Execute()

Retrieve NVLink Logical Partition



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
	org := "org_example" // string | Name of the NGC Org
	nvLinkLogicalPartitionId := "nvLinkLogicalPartitionId_example" // string | ID of the NVLink Logical Partition
	includeInterfaces := "includeInterfaces_example" // string | Include all attached NVLink Interfaces in response (optional)
	includeStats := true // bool | Include NVLink Logical Partition Stats in response (optional)
	includeVpcs := true // bool | Include all attached VPCs in response (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NVLinkLogicalPartitionAPI.GetNvlinkLogicalPartition(context.Background(), org, nvLinkLogicalPartitionId).IncludeInterfaces(includeInterfaces).IncludeStats(includeStats).IncludeVpcs(includeVpcs).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NVLinkLogicalPartitionAPI.GetNvlinkLogicalPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNvlinkLogicalPartition`: NVLinkLogicalPartition
	fmt.Fprintf(os.Stdout, "Response from `NVLinkLogicalPartitionAPI.GetNvlinkLogicalPartition`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the NGC Org | 
**nvLinkLogicalPartitionId** | **string** | ID of the NVLink Logical Partition | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNvlinkLogicalPartitionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeInterfaces** | **string** | Include all attached NVLink Interfaces in response | 
 **includeStats** | **bool** | Include NVLink Logical Partition Stats in response | 
 **includeVpcs** | **bool** | Include all attached VPCs in response | 
 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**NVLinkLogicalPartition**](NVLinkLogicalPartition.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateNvlinkLogicalPartition

> NVLinkLogicalPartition UpdateNvlinkLogicalPartition(ctx, org, nvLinkLogicalPartitionId).NVLinkLogicalPartitionUpdateRequest(nVLinkLogicalPartitionUpdateRequest).Execute()

Update NVLink Logical Partition



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
	org := "org_example" // string | Name of the NGC Org
	nvLinkLogicalPartitionId := "nvLinkLogicalPartitionId_example" // string | ID of the NVLink Logical Partition
	nVLinkLogicalPartitionUpdateRequest := *openapiclient.NewNVLinkLogicalPartitionUpdateRequest() // NVLinkLogicalPartitionUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NVLinkLogicalPartitionAPI.UpdateNvlinkLogicalPartition(context.Background(), org, nvLinkLogicalPartitionId).NVLinkLogicalPartitionUpdateRequest(nVLinkLogicalPartitionUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NVLinkLogicalPartitionAPI.UpdateNvlinkLogicalPartition``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateNvlinkLogicalPartition`: NVLinkLogicalPartition
	fmt.Fprintf(os.Stdout, "Response from `NVLinkLogicalPartitionAPI.UpdateNvlinkLogicalPartition`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the NGC Org | 
**nvLinkLogicalPartitionId** | **string** | ID of the NVLink Logical Partition | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateNvlinkLogicalPartitionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **nVLinkLogicalPartitionUpdateRequest** | [**NVLinkLogicalPartitionUpdateRequest**](NVLinkLogicalPartitionUpdateRequest.md) |  | 

### Return type

[**NVLinkLogicalPartition**](NVLinkLogicalPartition.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

