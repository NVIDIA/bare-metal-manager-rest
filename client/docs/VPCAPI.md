# \VPCAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateVpc**](VPCAPI.md#CreateVpc) | **Post** /v2/org/{org}/carbide/vpc | Create VPC
[**DeleteVpc**](VPCAPI.md#DeleteVpc) | **Delete** /v2/org/{org}/carbide/vpc/{vpcId} | Delete a VPC
[**GetAllVpc**](VPCAPI.md#GetAllVpc) | **Get** /v2/org/{org}/carbide/vpc | Retrieve all VPCs
[**GetVpc**](VPCAPI.md#GetVpc) | **Get** /v2/org/{org}/carbide/vpc/{vpcId} | Retrieve a VPC
[**UpdateVpc**](VPCAPI.md#UpdateVpc) | **Patch** /v2/org/{org}/carbide/vpc/{vpcId} | Update VPC
[**UpdateVpcVirtualization**](VPCAPI.md#UpdateVpcVirtualization) | **Patch** /v2/org/{org}/carbide/vpc/{vpcId}/virtualization | Update VPC Virtualization



## CreateVpc

> VPC CreateVpc(ctx, org).VpcCreateRequest(vpcCreateRequest).Execute()

Create VPC



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
	vpcCreateRequest := *openapiclient.NewVpcCreateRequest("Name_example", "SiteId_example") // VpcCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VPCAPI.CreateVpc(context.Background(), org).VpcCreateRequest(vpcCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCAPI.CreateVpc``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateVpc`: VPC
	fmt.Fprintf(os.Stdout, "Response from `VPCAPI.CreateVpc`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateVpcRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **vpcCreateRequest** | [**VpcCreateRequest**](VpcCreateRequest.md) |  | 

### Return type

[**VPC**](VPC.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteVpc

> DeleteVpc(ctx, org, vpcId).Execute()

Delete a VPC



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
	vpcId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the VPC

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VPCAPI.DeleteVpc(context.Background(), org, vpcId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCAPI.DeleteVpc``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**vpcId** | **string** | ID of the VPC | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteVpcRequest struct via the builder pattern


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


## GetAllVpc

> []VPC GetAllVpc(ctx, org).SiteId(siteId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).NetworkSecurityGroupId(networkSecurityGroupId).Execute()

Retrieve all VPCs



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter VPCs by Site ID (optional)
	status := "status_example" // string | Filter VPCs by Status (optional)
	query := "query_example" // string | Search for matches across all VPCs. Input will be matched against name, description, labels and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)
	networkSecurityGroupId := "networkSecurityGroupId_example" // string | Filter VPCs by NetworkSecurityGroup ID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VPCAPI.GetAllVpc(context.Background(), org).SiteId(siteId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).NetworkSecurityGroupId(networkSecurityGroupId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCAPI.GetAllVpc``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllVpc`: []VPC
	fmt.Fprintf(os.Stdout, "Response from `VPCAPI.GetAllVpc`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllVpcRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | Filter VPCs by Site ID | 
 **status** | **string** | Filter VPCs by Status | 
 **query** | **string** | Search for matches across all VPCs. Input will be matched against name, description, labels and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 
 **networkSecurityGroupId** | **string** | Filter VPCs by NetworkSecurityGroup ID | 

### Return type

[**[]VPC**](VPC.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVpc

> VPC GetVpc(ctx, org, vpcId).IncludeRelation(includeRelation).Execute()

Retrieve a VPC



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
	vpcId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the VPC
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VPCAPI.GetVpc(context.Background(), org, vpcId).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCAPI.GetVpc``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVpc`: VPC
	fmt.Fprintf(os.Stdout, "Response from `VPCAPI.GetVpc`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**vpcId** | **string** | ID of the VPC | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVpcRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**VPC**](VPC.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateVpc

> VPC UpdateVpc(ctx, org, vpcId).VpcUpdateRequest(vpcUpdateRequest).Execute()

Update VPC



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
	vpcId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the VPC
	vpcUpdateRequest := *openapiclient.NewVpcUpdateRequest() // VpcUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VPCAPI.UpdateVpc(context.Background(), org, vpcId).VpcUpdateRequest(vpcUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCAPI.UpdateVpc``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateVpc`: VPC
	fmt.Fprintf(os.Stdout, "Response from `VPCAPI.UpdateVpc`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**vpcId** | **string** | ID of the VPC | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateVpcRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **vpcUpdateRequest** | [**VpcUpdateRequest**](VpcUpdateRequest.md) |  | 

### Return type

[**VPC**](VPC.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateVpcVirtualization

> VPC UpdateVpcVirtualization(ctx, org, vpcId).VpcVirtualizationUpdateRequest(vpcVirtualizationUpdateRequest).Execute()

Update VPC Virtualization



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
	vpcId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the VPC
	vpcVirtualizationUpdateRequest := *openapiclient.NewVpcVirtualizationUpdateRequest() // VpcVirtualizationUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VPCAPI.UpdateVpcVirtualization(context.Background(), org, vpcId).VpcVirtualizationUpdateRequest(vpcVirtualizationUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VPCAPI.UpdateVpcVirtualization``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateVpcVirtualization`: VPC
	fmt.Fprintf(os.Stdout, "Response from `VPCAPI.UpdateVpcVirtualization`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**vpcId** | **string** | ID of the VPC | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateVpcVirtualizationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **vpcVirtualizationUpdateRequest** | [**VpcVirtualizationUpdateRequest**](VpcVirtualizationUpdateRequest.md) |  | 

### Return type

[**VPC**](VPC.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

