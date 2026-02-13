# \TenantAccountAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateTenantAccount**](TenantAccountAPI.md#CreateTenantAccount) | **Post** /v2/org/{org}/carbide/tenant/account | Create Tenant Account
[**DeleteTenantAccount**](TenantAccountAPI.md#DeleteTenantAccount) | **Delete** /v2/org/{org}/carbide/tenant/account/{accountId} | Delete Tenant Account
[**GetAllTenantAccount**](TenantAccountAPI.md#GetAllTenantAccount) | **Get** /v2/org/{org}/carbide/tenant/account | Retrieve all Tenant Accounts
[**GetTenantAccount**](TenantAccountAPI.md#GetTenantAccount) | **Get** /v2/org/{org}/carbide/tenant/account/{accountId} | Retrieve Tenant Account
[**UpdateTenantAccount**](TenantAccountAPI.md#UpdateTenantAccount) | **Patch** /v2/org/{org}/carbide/tenant/account/{accountId} | Update Tenant Account



## CreateTenantAccount

> TenantAccount CreateTenantAccount(ctx, org).TenantAccountCreateRequest(tenantAccountCreateRequest).Execute()

Create Tenant Account



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
	tenantAccountCreateRequest := *openapiclient.NewTenantAccountCreateRequest("InfrastructureProviderId_example", "TenantOrg_example") // TenantAccountCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantAccountAPI.CreateTenantAccount(context.Background(), org).TenantAccountCreateRequest(tenantAccountCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantAccountAPI.CreateTenantAccount``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateTenantAccount`: TenantAccount
	fmt.Fprintf(os.Stdout, "Response from `TenantAccountAPI.CreateTenantAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateTenantAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **tenantAccountCreateRequest** | [**TenantAccountCreateRequest**](TenantAccountCreateRequest.md) |  | 

### Return type

[**TenantAccount**](TenantAccount.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteTenantAccount

> DeleteTenantAccount(ctx, org, accountId).Execute()

Delete Tenant Account



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
	accountId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Tenant Account

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.TenantAccountAPI.DeleteTenantAccount(context.Background(), org, accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantAccountAPI.DeleteTenantAccount``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**accountId** | **string** | ID of the Tenant Account | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteTenantAccountRequest struct via the builder pattern


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


## GetAllTenantAccount

> []TenantAccount GetAllTenantAccount(ctx, org).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Tenant Accounts



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
	infrastructureProviderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter TenantAccounts by Infrastructure Provider ID (optional)
	tenantId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter TenantAccounts by Tenant ID (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantAccountAPI.GetAllTenantAccount(context.Background(), org).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantAccountAPI.GetAllTenantAccount``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllTenantAccount`: []TenantAccount
	fmt.Fprintf(os.Stdout, "Response from `TenantAccountAPI.GetAllTenantAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllTenantAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **infrastructureProviderId** | **string** | Filter TenantAccounts by Infrastructure Provider ID | 
 **tenantId** | **string** | Filter TenantAccounts by Tenant ID | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]TenantAccount**](TenantAccount.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTenantAccount

> TenantAccount GetTenantAccount(ctx, org, accountId).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).IncludeRelation(includeRelation).Execute()

Retrieve Tenant Account



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
	accountId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Tenant Account
	infrastructureProviderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Tenant Accounts by Infrastructure Provider ID (optional)
	tenantId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Tenant Accounts by Tenant ID (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantAccountAPI.GetTenantAccount(context.Background(), org, accountId).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantAccountAPI.GetTenantAccount``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTenantAccount`: TenantAccount
	fmt.Fprintf(os.Stdout, "Response from `TenantAccountAPI.GetTenantAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**accountId** | **string** | ID of the Tenant Account | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTenantAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **infrastructureProviderId** | **string** | Filter Tenant Accounts by Infrastructure Provider ID | 
 **tenantId** | **string** | Filter Tenant Accounts by Tenant ID | 
 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**TenantAccount**](TenantAccount.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateTenantAccount

> TenantAccount UpdateTenantAccount(ctx, org, accountId).Body(body).Execute()

Update Tenant Account



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
	accountId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Tenant Account
	body := map[string]interface{}{ ... } // map[string]interface{} | No params needed, an empty request body will suffice. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantAccountAPI.UpdateTenantAccount(context.Background(), org, accountId).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantAccountAPI.UpdateTenantAccount``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateTenantAccount`: TenantAccount
	fmt.Fprintf(os.Stdout, "Response from `TenantAccountAPI.UpdateTenantAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**accountId** | **string** | ID of the Tenant Account | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateTenantAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | **map[string]interface{}** | No params needed, an empty request body will suffice. | 

### Return type

[**TenantAccount**](TenantAccount.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

