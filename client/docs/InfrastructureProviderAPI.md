# \InfrastructureProviderAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCurrentInfrastructureProvider**](InfrastructureProviderAPI.md#GetCurrentInfrastructureProvider) | **Get** /v2/org/{org}/carbide/infrastructure-provider/current | Retrieve Infrastructure Provider for current Org
[**GetCurrentInfrastructureProviderStats**](InfrastructureProviderAPI.md#GetCurrentInfrastructureProviderStats) | **Get** /v2/org/{org}/carbide/infrastructure-provider/current/stats | Retrieve Stats for current Infrastructure Provider



## GetCurrentInfrastructureProvider

> InfrastructureProvider GetCurrentInfrastructureProvider(ctx, org).Execute()

Retrieve Infrastructure Provider for current Org



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InfrastructureProviderAPI.GetCurrentInfrastructureProvider(context.Background(), org).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InfrastructureProviderAPI.GetCurrentInfrastructureProvider``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCurrentInfrastructureProvider`: InfrastructureProvider
	fmt.Fprintf(os.Stdout, "Response from `InfrastructureProviderAPI.GetCurrentInfrastructureProvider`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCurrentInfrastructureProviderRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**InfrastructureProvider**](InfrastructureProvider.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCurrentInfrastructureProviderStats

> InfrastructureProviderStats GetCurrentInfrastructureProviderStats(ctx, org).Execute()

Retrieve Stats for current Infrastructure Provider



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InfrastructureProviderAPI.GetCurrentInfrastructureProviderStats(context.Background(), org).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InfrastructureProviderAPI.GetCurrentInfrastructureProviderStats``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCurrentInfrastructureProviderStats`: InfrastructureProviderStats
	fmt.Fprintf(os.Stdout, "Response from `InfrastructureProviderAPI.GetCurrentInfrastructureProviderStats`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCurrentInfrastructureProviderStatsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**InfrastructureProviderStats**](InfrastructureProviderStats.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

