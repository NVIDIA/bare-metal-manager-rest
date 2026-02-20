# \IPBlockAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateIpblock**](IPBlockAPI.md#CreateIpblock) | **Post** /v2/org/{org}/carbide/ipblock | Create IP Block
[**DeleteIpblock**](IPBlockAPI.md#DeleteIpblock) | **Delete** /v2/org/{org}/carbide/ipblock/{ipBlockId} | Delete IP Block
[**GetAllDerivedIpblock**](IPBlockAPI.md#GetAllDerivedIpblock) | **Get** /v2/org/{org}/carbide/ipblock/{ipBlockId}/derived | Retrieve All Derived IP Blocks
[**GetAllIpblock**](IPBlockAPI.md#GetAllIpblock) | **Get** /v2/org/{org}/carbide/ipblock | Retrieve all IP Blocks
[**GetIpblock**](IPBlockAPI.md#GetIpblock) | **Get** /v2/org/{org}/carbide/ipblock/{ipBlockId} | Retrieve IP Block
[**UpdateIpblock**](IPBlockAPI.md#UpdateIpblock) | **Patch** /v2/org/{org}/carbide/ipblock/{ipBlockId} | Update IP Block



## CreateIpblock

> IpBlock CreateIpblock(ctx, org).IpBlockCreateRequest(ipBlockCreateRequest).Execute()

Create IP Block



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
	ipBlockCreateRequest := *openapiclient.NewIpBlockCreateRequest("Name_example", "SiteId_example", "RoutingType_example", "Prefix_example", int32(123), "ProtocolVersion_example") // IpBlockCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IPBlockAPI.CreateIpblock(context.Background(), org).IpBlockCreateRequest(ipBlockCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPBlockAPI.CreateIpblock``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateIpblock`: IpBlock
	fmt.Fprintf(os.Stdout, "Response from `IPBlockAPI.CreateIpblock`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateIpblockRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **ipBlockCreateRequest** | [**IpBlockCreateRequest**](IpBlockCreateRequest.md) |  | 

### Return type

[**IpBlock**](IpBlock.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteIpblock

> DeleteIpblock(ctx, org, ipBlockId).Execute()

Delete IP Block



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
	ipBlockId := "ipBlockId_example" // string | ID of the IP Block

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPBlockAPI.DeleteIpblock(context.Background(), org, ipBlockId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPBlockAPI.DeleteIpblock``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**ipBlockId** | **string** | ID of the IP Block | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteIpblockRequest struct via the builder pattern


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


## GetAllDerivedIpblock

> []IpBlock GetAllDerivedIpblock(ctx, org, ipBlockId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve All Derived IP Blocks



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
	ipBlockId := "ipBlockId_example" // string | ID of the IP Block
	status := "status_example" // string | Filter IP Blocks by Status (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, description and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IPBlockAPI.GetAllDerivedIpblock(context.Background(), org, ipBlockId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPBlockAPI.GetAllDerivedIpblock``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllDerivedIpblock`: []IpBlock
	fmt.Fprintf(os.Stdout, "Response from `IPBlockAPI.GetAllDerivedIpblock`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**ipBlockId** | **string** | ID of the IP Block | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllDerivedIpblockRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **status** | **string** | Filter IP Blocks by Status | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, description and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]IpBlock**](IpBlock.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllIpblock

> []IpBlock GetAllIpblock(ctx, org).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).SiteId(siteId).Status(status).IncludeUsageStats(includeUsageStats).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all IP Blocks



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
	infrastructureProviderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter IP Blocks by Infrastructure Provider ID (optional)
	tenantId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter IP Blocks by Tenant ID (optional)
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter IP Blocks by Site ID (optional)
	status := "status_example" // string | Filter IP Blocks by Status (optional)
	includeUsageStats := true // bool | Include IP Block usage stats in response (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, description and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IPBlockAPI.GetAllIpblock(context.Background(), org).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).SiteId(siteId).Status(status).IncludeUsageStats(includeUsageStats).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPBlockAPI.GetAllIpblock``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllIpblock`: []IpBlock
	fmt.Fprintf(os.Stdout, "Response from `IPBlockAPI.GetAllIpblock`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllIpblockRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **infrastructureProviderId** | **string** | Filter IP Blocks by Infrastructure Provider ID | 
 **tenantId** | **string** | Filter IP Blocks by Tenant ID | 
 **siteId** | **string** | Filter IP Blocks by Site ID | 
 **status** | **string** | Filter IP Blocks by Status | 
 **includeUsageStats** | **bool** | Include IP Block usage stats in response | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, description and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]IpBlock**](IpBlock.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIpblock

> IpBlock GetIpblock(ctx, org, ipBlockId).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).IncludeUsageStats(includeUsageStats).IncludeRelation(includeRelation).Execute()

Retrieve IP Block



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
	ipBlockId := "ipBlockId_example" // string | ID of the IP Block
	infrastructureProviderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter IP Blocks by Infrastructure Provider ID (optional)
	tenantId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter IP Blocks by Tenant ID (optional)
	includeUsageStats := true // bool | Include IP Block usage stats in response (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IPBlockAPI.GetIpblock(context.Background(), org, ipBlockId).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).IncludeUsageStats(includeUsageStats).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPBlockAPI.GetIpblock``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetIpblock`: IpBlock
	fmt.Fprintf(os.Stdout, "Response from `IPBlockAPI.GetIpblock`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**ipBlockId** | **string** | ID of the IP Block | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetIpblockRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **infrastructureProviderId** | **string** | Filter IP Blocks by Infrastructure Provider ID | 
 **tenantId** | **string** | Filter IP Blocks by Tenant ID | 
 **includeUsageStats** | **bool** | Include IP Block usage stats in response | 
 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**IpBlock**](IpBlock.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateIpblock

> IpBlock UpdateIpblock(ctx, org, ipBlockId).IpBlockUpdateRequest(ipBlockUpdateRequest).Execute()

Update IP Block



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
	ipBlockId := "ipBlockId_example" // string | ID of the IP Block
	ipBlockUpdateRequest := *openapiclient.NewIpBlockUpdateRequest() // IpBlockUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.IPBlockAPI.UpdateIpblock(context.Background(), org, ipBlockId).IpBlockUpdateRequest(ipBlockUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPBlockAPI.UpdateIpblock``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateIpblock`: IpBlock
	fmt.Fprintf(os.Stdout, "Response from `IPBlockAPI.UpdateIpblock`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**ipBlockId** | **string** | ID of the IP Block | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateIpblockRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **ipBlockUpdateRequest** | [**IpBlockUpdateRequest**](IpBlockUpdateRequest.md) |  | 

### Return type

[**IpBlock**](IpBlock.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

