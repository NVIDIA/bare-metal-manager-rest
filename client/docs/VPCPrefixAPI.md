# \VPCPrefixAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateVpcPrefix**](VPCPrefixAPI.md#CreateVpcPrefix) | **Post** /v2/org/{org}/carbide/vpc-prefix | Create VPC Prefix
[**DeleteVpcPrefix**](VPCPrefixAPI.md#DeleteVpcPrefix) | **Delete** /v2/org/{org}/carbide/vpc-prefix/{vpcPrefixId} | Delete VPC Prefix
[**GetAllVpcPrefix**](VPCPrefixAPI.md#GetAllVpcPrefix) | **Get** /v2/org/{org}/carbide/vpc-prefix | Retrieve all VPC Prefixes
[**GetVpcPrefix**](VPCPrefixAPI.md#GetVpcPrefix) | **Get** /v2/org/{org}/carbide/vpc-prefix/{vpcPrefixId} | Retrieve VPC Prefix
[**UpdateVpcPrefix**](VPCPrefixAPI.md#UpdateVpcPrefix) | **Patch** /v2/org/{org}/carbide/vpc-prefix/{vpcPrefixId} | Update VPC Prefix



## CreateVpcPrefix

> VpcPrefix CreateVpcPrefix(ctx, org).VpcPrefixCreateRequest(vpcPrefixCreateRequest).Execute()

Create VPC Prefix



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
	vpcPrefixCreateRequest := *openapiclient.NewVpcPrefixCreateRequest("Name_example", "VpcId_example", int32(123)) // VpcPrefixCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VPCPrefixAPI.CreateVpcPrefix(context.Background(), org).VpcPrefixCreateRequest(vpcPrefixCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCPrefixAPI.CreateVpcPrefix``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateVpcPrefix`: VpcPrefix
	fmt.Fprintf(os.Stdout, "Response from `VPCPrefixAPI.CreateVpcPrefix`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateVpcPrefixRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **vpcPrefixCreateRequest** | [**VpcPrefixCreateRequest**](VpcPrefixCreateRequest.md) |  | 

### Return type

[**VpcPrefix**](VpcPrefix.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteVpcPrefix

> DeleteVpcPrefix(ctx, org, vpcPrefixId).Execute()

Delete VPC Prefix



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
	vpcPrefixId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the VPC Prefix

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VPCPrefixAPI.DeleteVpcPrefix(context.Background(), org, vpcPrefixId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCPrefixAPI.DeleteVpcPrefix``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**vpcPrefixId** | **string** | ID of the VPC Prefix | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteVpcPrefixRequest struct via the builder pattern


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


## GetAllVpcPrefix

> []VpcPrefix GetAllVpcPrefix(ctx, org).SiteId(siteId).VpcId(vpcId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all VPC Prefixes



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter VPC Prefixes by Site, required if vpcId query param is not specified (optional)
	vpcId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter VPC Prefixes by VPC (optional)
	status := "status_example" // string | Filter VPC Prefixes by Status (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VPCPrefixAPI.GetAllVpcPrefix(context.Background(), org).SiteId(siteId).VpcId(vpcId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCPrefixAPI.GetAllVpcPrefix``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllVpcPrefix`: []VpcPrefix
	fmt.Fprintf(os.Stdout, "Response from `VPCPrefixAPI.GetAllVpcPrefix`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllVpcPrefixRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | Filter VPC Prefixes by Site, required if vpcId query param is not specified | 
 **vpcId** | **string** | Filter VPC Prefixes by VPC | 
 **status** | **string** | Filter VPC Prefixes by Status | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]VpcPrefix**](VpcPrefix.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVpcPrefix

> VpcPrefix GetVpcPrefix(ctx, org, vpcPrefixId).IncludeRelation(includeRelation).Execute()

Retrieve VPC Prefix



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
	vpcPrefixId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the VPC Prefix
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VPCPrefixAPI.GetVpcPrefix(context.Background(), org, vpcPrefixId).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCPrefixAPI.GetVpcPrefix``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVpcPrefix`: VpcPrefix
	fmt.Fprintf(os.Stdout, "Response from `VPCPrefixAPI.GetVpcPrefix`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**vpcPrefixId** | **string** | ID of the VPC Prefix | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVpcPrefixRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**VpcPrefix**](VpcPrefix.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateVpcPrefix

> VpcPrefix UpdateVpcPrefix(ctx, org, vpcPrefixId).VpcPrefixUpdateRequest(vpcPrefixUpdateRequest).Execute()

Update VPC Prefix



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
	vpcPrefixId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the VPC Prefix
	vpcPrefixUpdateRequest := *openapiclient.NewVpcPrefixUpdateRequest("Name_example") // VpcPrefixUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VPCPrefixAPI.UpdateVpcPrefix(context.Background(), org, vpcPrefixId).VpcPrefixUpdateRequest(vpcPrefixUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCPrefixAPI.UpdateVpcPrefix``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateVpcPrefix`: VpcPrefix
	fmt.Fprintf(os.Stdout, "Response from `VPCPrefixAPI.UpdateVpcPrefix`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**vpcPrefixId** | **string** | ID of the VPC Prefix | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateVpcPrefixRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **vpcPrefixUpdateRequest** | [**VpcPrefixUpdateRequest**](VpcPrefixUpdateRequest.md) |  | 

### Return type

[**VpcPrefix**](VpcPrefix.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

