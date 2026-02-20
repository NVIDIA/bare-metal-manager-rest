# \NetworkSecurityGroupAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateNetworkSecurityGroup**](NetworkSecurityGroupAPI.md#CreateNetworkSecurityGroup) | **Post** /v2/org/{org}/carbide/network-security-group | Create Network Security Group
[**DeleteNetworkSecurityGroup**](NetworkSecurityGroupAPI.md#DeleteNetworkSecurityGroup) | **Delete** /v2/org/{org}/carbide/network-security-group/{networkSecurityGroupId} | Delete Network Security Group
[**GetAllNetworkSecurityGroup**](NetworkSecurityGroupAPI.md#GetAllNetworkSecurityGroup) | **Get** /v2/org/{org}/carbide/network-security-group | Retrieve all Network Security Groups
[**GetNetworkSecurityGroup**](NetworkSecurityGroupAPI.md#GetNetworkSecurityGroup) | **Get** /v2/org/{org}/carbide/network-security-group/{networkSecurityGroupId} | Retrieve Network Security Group
[**UpdateNetworkSecurityGroup**](NetworkSecurityGroupAPI.md#UpdateNetworkSecurityGroup) | **Patch** /v2/org/{org}/carbide/network-security-group/{networkSecurityGroupId} | Update Network Security Group



## CreateNetworkSecurityGroup

> NetworkSecurityGroup CreateNetworkSecurityGroup(ctx, org).NetworkSecurityGroupCreateRequest(networkSecurityGroupCreateRequest).Execute()

Create Network Security Group



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
	networkSecurityGroupCreateRequest := *openapiclient.NewNetworkSecurityGroupCreateRequest("Name_example", "SiteId_example") // NetworkSecurityGroupCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NetworkSecurityGroupAPI.CreateNetworkSecurityGroup(context.Background(), org).NetworkSecurityGroupCreateRequest(networkSecurityGroupCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NetworkSecurityGroupAPI.CreateNetworkSecurityGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateNetworkSecurityGroup`: NetworkSecurityGroup
	fmt.Fprintf(os.Stdout, "Response from `NetworkSecurityGroupAPI.CreateNetworkSecurityGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateNetworkSecurityGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **networkSecurityGroupCreateRequest** | [**NetworkSecurityGroupCreateRequest**](NetworkSecurityGroupCreateRequest.md) |  | 

### Return type

[**NetworkSecurityGroup**](NetworkSecurityGroup.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteNetworkSecurityGroup

> DeleteNetworkSecurityGroup(ctx, org, networkSecurityGroupId).Execute()

Delete Network Security Group



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
	networkSecurityGroupId := "networkSecurityGroupId_example" // string | ID of the Network Security Group

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.NetworkSecurityGroupAPI.DeleteNetworkSecurityGroup(context.Background(), org, networkSecurityGroupId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NetworkSecurityGroupAPI.DeleteNetworkSecurityGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**networkSecurityGroupId** | **string** | ID of the Network Security Group | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteNetworkSecurityGroupRequest struct via the builder pattern


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


## GetAllNetworkSecurityGroup

> []NetworkSecurityGroup GetAllNetworkSecurityGroup(ctx, org).SiteId(siteId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).IncludeAttachmentStats(includeAttachmentStats).Execute()

Retrieve all Network Security Groups



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter By Site ID (optional)
	status := "status_example" // string | Filter Network Security Groups by Status (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, description and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)
	includeAttachmentStats := true // bool | Include counts for the number objects that have attached the Network Security Group (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NetworkSecurityGroupAPI.GetAllNetworkSecurityGroup(context.Background(), org).SiteId(siteId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).IncludeAttachmentStats(includeAttachmentStats).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NetworkSecurityGroupAPI.GetAllNetworkSecurityGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllNetworkSecurityGroup`: []NetworkSecurityGroup
	fmt.Fprintf(os.Stdout, "Response from `NetworkSecurityGroupAPI.GetAllNetworkSecurityGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllNetworkSecurityGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | Filter By Site ID | 
 **status** | **string** | Filter Network Security Groups by Status | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, description and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 
 **includeAttachmentStats** | **bool** | Include counts for the number objects that have attached the Network Security Group | 

### Return type

[**[]NetworkSecurityGroup**](NetworkSecurityGroup.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNetworkSecurityGroup

> NetworkSecurityGroup GetNetworkSecurityGroup(ctx, org, networkSecurityGroupId).IncludeRelation(includeRelation).Execute()

Retrieve Network Security Group



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
	networkSecurityGroupId := "networkSecurityGroupId_example" // string | ID of the Network Security Group
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NetworkSecurityGroupAPI.GetNetworkSecurityGroup(context.Background(), org, networkSecurityGroupId).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NetworkSecurityGroupAPI.GetNetworkSecurityGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNetworkSecurityGroup`: NetworkSecurityGroup
	fmt.Fprintf(os.Stdout, "Response from `NetworkSecurityGroupAPI.GetNetworkSecurityGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**networkSecurityGroupId** | **string** | ID of the Network Security Group | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNetworkSecurityGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**NetworkSecurityGroup**](NetworkSecurityGroup.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateNetworkSecurityGroup

> NetworkSecurityGroup UpdateNetworkSecurityGroup(ctx, org, networkSecurityGroupId).NetworkSecurityGroupUpdateRequest(networkSecurityGroupUpdateRequest).Execute()

Update Network Security Group



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
	networkSecurityGroupId := "networkSecurityGroupId_example" // string | ID of the Network Security Group
	networkSecurityGroupUpdateRequest := *openapiclient.NewNetworkSecurityGroupUpdateRequest() // NetworkSecurityGroupUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NetworkSecurityGroupAPI.UpdateNetworkSecurityGroup(context.Background(), org, networkSecurityGroupId).NetworkSecurityGroupUpdateRequest(networkSecurityGroupUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NetworkSecurityGroupAPI.UpdateNetworkSecurityGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateNetworkSecurityGroup`: NetworkSecurityGroup
	fmt.Fprintf(os.Stdout, "Response from `NetworkSecurityGroupAPI.UpdateNetworkSecurityGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**networkSecurityGroupId** | **string** | ID of the Network Security Group | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateNetworkSecurityGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **networkSecurityGroupUpdateRequest** | [**NetworkSecurityGroupUpdateRequest**](NetworkSecurityGroupUpdateRequest.md) |  | 

### Return type

[**NetworkSecurityGroup**](NetworkSecurityGroup.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

