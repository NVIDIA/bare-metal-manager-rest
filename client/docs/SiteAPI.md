# \SiteAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateSite**](SiteAPI.md#CreateSite) | **Post** /v2/org/{org}/carbide/site | Create Site
[**DeleteSite**](SiteAPI.md#DeleteSite) | **Delete** /v2/org/{org}/carbide/site/{siteId} | Delete Site
[**GetAllSite**](SiteAPI.md#GetAllSite) | **Get** /v2/org/{org}/carbide/site | Retrieve all Sites
[**GetSite**](SiteAPI.md#GetSite) | **Get** /v2/org/{org}/carbide/site/{siteId} | Retrieve Site
[**GetSiteStatusHistory**](SiteAPI.md#GetSiteStatusHistory) | **Get** /v2/org/{org}/carbide/site/{siteId}/status-history | Retrieve Site status history
[**UpdateSite**](SiteAPI.md#UpdateSite) | **Patch** /v2/org/{org}/carbide/site/{siteId} | Update Site



## CreateSite

> Site CreateSite(ctx, org).SiteCreateRequest(siteCreateRequest).Execute()

Create Site



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
	siteCreateRequest := *openapiclient.NewSiteCreateRequest("Name_example") // SiteCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SiteAPI.CreateSite(context.Background(), org).SiteCreateRequest(siteCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SiteAPI.CreateSite``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateSite`: Site
	fmt.Fprintf(os.Stdout, "Response from `SiteAPI.CreateSite`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateSiteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteCreateRequest** | [**SiteCreateRequest**](SiteCreateRequest.md) |  | 

### Return type

[**Site**](Site.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteSite

> DeleteSite(ctx, org, siteId).PurgeMachines(purgeMachines).Execute()

Delete Site



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Site
	purgeMachines := true // bool | Scrub all Machine data associated with this Site to re-pair (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SiteAPI.DeleteSite(context.Background(), org, siteId).PurgeMachines(purgeMachines).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SiteAPI.DeleteSite``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**siteId** | **string** | ID of the Site | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteSiteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **purgeMachines** | **bool** | Scrub all Machine data associated with this Site to re-pair | 

### Return type

 (empty response body)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllSite

> []Site GetAllSite(ctx, org).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).Status(status).IsNativeNetworkingEnabled(isNativeNetworkingEnabled).IsNetworkSecurityGroupEnabled(isNetworkSecurityGroupEnabled).IsNVLinkPartitionEnabled(isNVLinkPartitionEnabled).IncludeMachineStats(includeMachineStats).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Sites



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
	infrastructureProviderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Sites by Infrastructure Provider ID (optional)
	tenantId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Sites by Tenant ID (optional)
	status := "status_example" // string | Filter Sites by Status. Can be specified multiple times to filter on more than one status (optional)
	isNativeNetworkingEnabled := true // bool | Filter Sites by native networking enabled flag. Requires Provider Admin role. (optional)
	isNetworkSecurityGroupEnabled := true // bool | Filter Sites by network security group enabled flag. Requires Provider Admin role. (optional)
	isNVLinkPartitionEnabled := true // bool | Filter Sites by NVLink partitioning enabled flag. Requires Provider Admin role. (optional)
	includeMachineStats := true // bool | Include a breakdown of Machine counts by life-cycle status and health. Requires Provider Admin role. (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, description, location, contact, and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SiteAPI.GetAllSite(context.Background(), org).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).Status(status).IsNativeNetworkingEnabled(isNativeNetworkingEnabled).IsNetworkSecurityGroupEnabled(isNetworkSecurityGroupEnabled).IsNVLinkPartitionEnabled(isNVLinkPartitionEnabled).IncludeMachineStats(includeMachineStats).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SiteAPI.GetAllSite``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllSite`: []Site
	fmt.Fprintf(os.Stdout, "Response from `SiteAPI.GetAllSite`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllSiteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **infrastructureProviderId** | **string** | Filter Sites by Infrastructure Provider ID | 
 **tenantId** | **string** | Filter Sites by Tenant ID | 
 **status** | **string** | Filter Sites by Status. Can be specified multiple times to filter on more than one status | 
 **isNativeNetworkingEnabled** | **bool** | Filter Sites by native networking enabled flag. Requires Provider Admin role. | 
 **isNetworkSecurityGroupEnabled** | **bool** | Filter Sites by network security group enabled flag. Requires Provider Admin role. | 
 **isNVLinkPartitionEnabled** | **bool** | Filter Sites by NVLink partitioning enabled flag. Requires Provider Admin role. | 
 **includeMachineStats** | **bool** | Include a breakdown of Machine counts by life-cycle status and health. Requires Provider Admin role. | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, description, location, contact, and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]Site**](Site.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSite

> Site GetSite(ctx, org, siteId).Execute()

Retrieve Site



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Site

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SiteAPI.GetSite(context.Background(), org, siteId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SiteAPI.GetSite``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSite`: Site
	fmt.Fprintf(os.Stdout, "Response from `SiteAPI.GetSite`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**siteId** | **string** | ID of the Site | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSiteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Site**](Site.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSiteStatusHistory

> []StatusDetail GetSiteStatusHistory(ctx, org, siteId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve Site status history



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Site
	pageNumber := int32(56) // int32 | Page number for pagination query (optional)
	pageSize := int32(56) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SiteAPI.GetSiteStatusHistory(context.Background(), org, siteId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SiteAPI.GetSiteStatusHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSiteStatusHistory`: []StatusDetail
	fmt.Fprintf(os.Stdout, "Response from `SiteAPI.GetSiteStatusHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**siteId** | **string** | ID of the Site | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSiteStatusHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **pageNumber** | **int32** | Page number for pagination query | 
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]StatusDetail**](StatusDetail.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateSite

> Site UpdateSite(ctx, org, siteId).SiteUpdateRequest(siteUpdateRequest).Execute()

Update Site



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Site
	siteUpdateRequest := *openapiclient.NewSiteUpdateRequest() // SiteUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SiteAPI.UpdateSite(context.Background(), org, siteId).SiteUpdateRequest(siteUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SiteAPI.UpdateSite``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateSite`: Site
	fmt.Fprintf(os.Stdout, "Response from `SiteAPI.UpdateSite`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**siteId** | **string** | ID of the Site | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateSiteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **siteUpdateRequest** | [**SiteUpdateRequest**](SiteUpdateRequest.md) |  | 

### Return type

[**Site**](Site.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

