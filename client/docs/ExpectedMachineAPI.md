# \ExpectedMachineAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BatchCreateExpectedMachines**](ExpectedMachineAPI.md#BatchCreateExpectedMachines) | **Post** /v2/org/{org}/forge/expected-machine/batch | Batch Create Expected Machines
[**BatchUpdateExpectedMachines**](ExpectedMachineAPI.md#BatchUpdateExpectedMachines) | **Patch** /v2/org/{org}/forge/expected-machine/batch | Batch Update Expected Machines
[**CreateExpectedMachine**](ExpectedMachineAPI.md#CreateExpectedMachine) | **Post** /v2/org/{org}/carbide/expected-machine | Create Expected Machine
[**DeleteExpectedMachine**](ExpectedMachineAPI.md#DeleteExpectedMachine) | **Delete** /v2/org/{org}/carbide/expected-machine/{expectedMachineId} | Delete Expected Machine
[**GetAllExpectedMachine**](ExpectedMachineAPI.md#GetAllExpectedMachine) | **Get** /v2/org/{org}/carbide/expected-machine | Retrieve all Expected Machines
[**GetExpectedMachine**](ExpectedMachineAPI.md#GetExpectedMachine) | **Get** /v2/org/{org}/carbide/expected-machine/{expectedMachineId} | Retrieve Expected Machine
[**UpdateExpectedMachine**](ExpectedMachineAPI.md#UpdateExpectedMachine) | **Patch** /v2/org/{org}/carbide/expected-machine/{expectedMachineId} | Update Expected Machine



## BatchCreateExpectedMachines

> []ExpectedMachine BatchCreateExpectedMachines(ctx, org).ExpectedMachineCreateRequest(expectedMachineCreateRequest).Execute()

Batch Create Expected Machines



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
	expectedMachineCreateRequest := []openapiclient.ExpectedMachineCreateRequest{*openapiclient.NewExpectedMachineCreateRequest("SiteId_example", "BmcMacAddress_example", "ChassisSerialNumber_example")} // []ExpectedMachineCreateRequest | Array of Expected Machine creation requests

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExpectedMachineAPI.BatchCreateExpectedMachines(context.Background(), org).ExpectedMachineCreateRequest(expectedMachineCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExpectedMachineAPI.BatchCreateExpectedMachines``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BatchCreateExpectedMachines`: []ExpectedMachine
	fmt.Fprintf(os.Stdout, "Response from `ExpectedMachineAPI.BatchCreateExpectedMachines`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiBatchCreateExpectedMachinesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **expectedMachineCreateRequest** | [**[]ExpectedMachineCreateRequest**](ExpectedMachineCreateRequest.md) | Array of Expected Machine creation requests | 

### Return type

[**[]ExpectedMachine**](ExpectedMachine.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BatchUpdateExpectedMachines

> []ExpectedMachine BatchUpdateExpectedMachines(ctx, org).BatchUpdateExpectedMachinesRequestInner(batchUpdateExpectedMachinesRequestInner).Execute()

Batch Update Expected Machines



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
	batchUpdateExpectedMachinesRequestInner := []openapiclient.BatchUpdateExpectedMachinesRequestInner{*openapiclient.NewBatchUpdateExpectedMachinesRequestInner("Id_example")} // []BatchUpdateExpectedMachinesRequestInner | Array of Expected Machine update requests

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExpectedMachineAPI.BatchUpdateExpectedMachines(context.Background(), org).BatchUpdateExpectedMachinesRequestInner(batchUpdateExpectedMachinesRequestInner).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExpectedMachineAPI.BatchUpdateExpectedMachines``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BatchUpdateExpectedMachines`: []ExpectedMachine
	fmt.Fprintf(os.Stdout, "Response from `ExpectedMachineAPI.BatchUpdateExpectedMachines`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiBatchUpdateExpectedMachinesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **batchUpdateExpectedMachinesRequestInner** | [**[]BatchUpdateExpectedMachinesRequestInner**](BatchUpdateExpectedMachinesRequestInner.md) | Array of Expected Machine update requests | 

### Return type

[**[]ExpectedMachine**](ExpectedMachine.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateExpectedMachine

> ExpectedMachine CreateExpectedMachine(ctx, org).ExpectedMachineCreateRequest(expectedMachineCreateRequest).Execute()

Create Expected Machine



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
	expectedMachineCreateRequest := *openapiclient.NewExpectedMachineCreateRequest("SiteId_example", "BmcMacAddress_example", "ChassisSerialNumber_example") // ExpectedMachineCreateRequest | Expected Machine creation request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExpectedMachineAPI.CreateExpectedMachine(context.Background(), org).ExpectedMachineCreateRequest(expectedMachineCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExpectedMachineAPI.CreateExpectedMachine``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateExpectedMachine`: ExpectedMachine
	fmt.Fprintf(os.Stdout, "Response from `ExpectedMachineAPI.CreateExpectedMachine`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateExpectedMachineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **expectedMachineCreateRequest** | [**ExpectedMachineCreateRequest**](ExpectedMachineCreateRequest.md) | Expected Machine creation request | 

### Return type

[**ExpectedMachine**](ExpectedMachine.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteExpectedMachine

> DeleteExpectedMachine(ctx, org, expectedMachineId).Execute()

Delete Expected Machine



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
	expectedMachineId := "expectedMachineId_example" // string | ID of the Expected Machine

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ExpectedMachineAPI.DeleteExpectedMachine(context.Background(), org, expectedMachineId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExpectedMachineAPI.DeleteExpectedMachine``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**expectedMachineId** | **string** | ID of the Expected Machine | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteExpectedMachineRequest struct via the builder pattern


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


## GetAllExpectedMachine

> []ExpectedMachine GetAllExpectedMachine(ctx, org).SiteId(siteId).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Expected Machines



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Site to filter Expected Machines by (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExpectedMachineAPI.GetAllExpectedMachine(context.Background(), org).SiteId(siteId).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExpectedMachineAPI.GetAllExpectedMachine``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllExpectedMachine`: []ExpectedMachine
	fmt.Fprintf(os.Stdout, "Response from `ExpectedMachineAPI.GetAllExpectedMachine`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllExpectedMachineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | ID of the Site to filter Expected Machines by | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]ExpectedMachine**](ExpectedMachine.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetExpectedMachine

> ExpectedMachine GetExpectedMachine(ctx, org, expectedMachineId).IncludeRelation(includeRelation).Execute()

Retrieve Expected Machine



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
	expectedMachineId := "expectedMachineId_example" // string | ID of the Expected Machine
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExpectedMachineAPI.GetExpectedMachine(context.Background(), org, expectedMachineId).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExpectedMachineAPI.GetExpectedMachine``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetExpectedMachine`: ExpectedMachine
	fmt.Fprintf(os.Stdout, "Response from `ExpectedMachineAPI.GetExpectedMachine`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**expectedMachineId** | **string** | ID of the Expected Machine | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetExpectedMachineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**ExpectedMachine**](ExpectedMachine.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateExpectedMachine

> ExpectedMachine UpdateExpectedMachine(ctx, org, expectedMachineId).ExpectedMachineUpdateRequest(expectedMachineUpdateRequest).Execute()

Update Expected Machine



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
	expectedMachineId := "expectedMachineId_example" // string | ID of the Expected Machine
	expectedMachineUpdateRequest := *openapiclient.NewExpectedMachineUpdateRequest() // ExpectedMachineUpdateRequest | Expected Machine update request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExpectedMachineAPI.UpdateExpectedMachine(context.Background(), org, expectedMachineId).ExpectedMachineUpdateRequest(expectedMachineUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExpectedMachineAPI.UpdateExpectedMachine``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateExpectedMachine`: ExpectedMachine
	fmt.Fprintf(os.Stdout, "Response from `ExpectedMachineAPI.UpdateExpectedMachine`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**expectedMachineId** | **string** | ID of the Expected Machine | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateExpectedMachineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **expectedMachineUpdateRequest** | [**ExpectedMachineUpdateRequest**](ExpectedMachineUpdateRequest.md) | Expected Machine update request | 

### Return type

[**ExpectedMachine**](ExpectedMachine.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

