# \InstanceTypeAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateInstanceType**](InstanceTypeAPI.md#CreateInstanceType) | **Post** /v2/org/{org}/carbide/instance/type | Create an Instance Type
[**CreateInstanceTypeMachineAssociation**](InstanceTypeAPI.md#CreateInstanceTypeMachineAssociation) | **Post** /v2/org/{org}/carbide/instance/type/{instanceTypeId}/machine | Create a Machine/Instance Type association
[**DeleteInstanceType**](InstanceTypeAPI.md#DeleteInstanceType) | **Delete** /v2/org/{org}/carbide/instance/type/{instanceTypeId} | Delete Instance Type
[**DeleteInstanceTypeMachineAssociation**](InstanceTypeAPI.md#DeleteInstanceTypeMachineAssociation) | **Delete** /v2/org/{org}/carbide/instance/type/{instanceTypeId}/machine/{machineAssociationId} | Delete a Machine/Instance Type association
[**GetAllInstanceType**](InstanceTypeAPI.md#GetAllInstanceType) | **Get** /v2/org/{org}/carbide/instance/type | Retrieve all Instance Types
[**GetInstanceType**](InstanceTypeAPI.md#GetInstanceType) | **Get** /v2/org/{org}/carbide/instance/type/{instanceTypeId} | Retrieve an Instance Type
[**GetInstanceTypeMachineAssociation**](InstanceTypeAPI.md#GetInstanceTypeMachineAssociation) | **Get** /v2/org/{org}/carbide/instance/type/{instanceTypeId}/machine | Retrieve all Machines/Instance Type associations
[**UpdateInstanceType**](InstanceTypeAPI.md#UpdateInstanceType) | **Patch** /v2/org/{org}/carbide/instance/type/{instanceTypeId} | Update Instance Type



## CreateInstanceType

> InstanceType CreateInstanceType(ctx, org).InstanceTypeCreateRequest(instanceTypeCreateRequest).Execute()

Create an Instance Type



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
	instanceTypeCreateRequest := *openapiclient.NewInstanceTypeCreateRequest("Name_example", "SiteId_example") // InstanceTypeCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceTypeAPI.CreateInstanceType(context.Background(), org).InstanceTypeCreateRequest(instanceTypeCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceTypeAPI.CreateInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateInstanceType`: InstanceType
	fmt.Fprintf(os.Stdout, "Response from `InstanceTypeAPI.CreateInstanceType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateInstanceTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **instanceTypeCreateRequest** | [**InstanceTypeCreateRequest**](InstanceTypeCreateRequest.md) |  | 

### Return type

[**InstanceType**](InstanceType.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateInstanceTypeMachineAssociation

> []MachineInstanceType CreateInstanceTypeMachineAssociation(ctx, org, instanceTypeId).MachineInstanceTypeCreateRequest(machineInstanceTypeCreateRequest).Execute()

Create a Machine/Instance Type association



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
	instanceTypeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Instance Type
	machineInstanceTypeCreateRequest := *openapiclient.NewMachineInstanceTypeCreateRequest([]string{"MachineIds_example"}) // MachineInstanceTypeCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceTypeAPI.CreateInstanceTypeMachineAssociation(context.Background(), org, instanceTypeId).MachineInstanceTypeCreateRequest(machineInstanceTypeCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceTypeAPI.CreateInstanceTypeMachineAssociation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateInstanceTypeMachineAssociation`: []MachineInstanceType
	fmt.Fprintf(os.Stdout, "Response from `InstanceTypeAPI.CreateInstanceTypeMachineAssociation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceTypeId** | **string** | ID of the Instance Type | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateInstanceTypeMachineAssociationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **machineInstanceTypeCreateRequest** | [**MachineInstanceTypeCreateRequest**](MachineInstanceTypeCreateRequest.md) |  | 

### Return type

[**[]MachineInstanceType**](MachineInstanceType.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteInstanceType

> DeleteInstanceType(ctx, org, instanceTypeId).Execute()

Delete Instance Type



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
	instanceTypeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Instance Type

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.InstanceTypeAPI.DeleteInstanceType(context.Background(), org, instanceTypeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceTypeAPI.DeleteInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceTypeId** | **string** | ID of the Instance Type | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteInstanceTypeRequest struct via the builder pattern


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


## DeleteInstanceTypeMachineAssociation

> DeleteInstanceTypeMachineAssociation(ctx, org, instanceTypeId, machineAssociationId).Execute()

Delete a Machine/Instance Type association



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
	instanceTypeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Instance Type
	machineAssociationId := "machineAssociationId_example" // string | ID of the Machine/Instance Type Association

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.InstanceTypeAPI.DeleteInstanceTypeMachineAssociation(context.Background(), org, instanceTypeId, machineAssociationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceTypeAPI.DeleteInstanceTypeMachineAssociation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceTypeId** | **string** | ID of the Instance Type | 
**machineAssociationId** | **string** | ID of the Machine/Instance Type Association | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteInstanceTypeMachineAssociationRequest struct via the builder pattern


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


## GetAllInstanceType

> []InstanceType GetAllInstanceType(ctx, org).SiteId(siteId).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).Status(status).Query(query).IncludeRelation(includeRelation).IncludeMachineAssignment(includeMachineAssignment).IncludeAllocationStats(includeAllocationStats).ExcludeUnallocated(excludeUnallocated).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Instance Types



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
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Instance Types by Site ID
	org := "org_example" // string | Name of the Org
	infrastructureProviderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Instance Types by Infrastructure Provider ID (optional)
	tenantId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Instance Types by Tenant ID (optional)
	status := "status_example" // string | Filter Instance Types by Status (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, display name, description and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	includeMachineAssignment := true // bool | Include Machine assignments for each Instance Type. Can only be requested by Provider. (optional)
	includeAllocationStats := true // bool | Include Allocation stats. Currently can only be requested by Tenant (optional)
	excludeUnallocated := true // bool | Excludes InstanceType records that have no allocations from being returned in the result set. Currently can only be requested by Tenant. (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceTypeAPI.GetAllInstanceType(context.Background(), org).SiteId(siteId).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).Status(status).Query(query).IncludeRelation(includeRelation).IncludeMachineAssignment(includeMachineAssignment).IncludeAllocationStats(includeAllocationStats).ExcludeUnallocated(excludeUnallocated).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceTypeAPI.GetAllInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllInstanceType`: []InstanceType
	fmt.Fprintf(os.Stdout, "Response from `InstanceTypeAPI.GetAllInstanceType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllInstanceTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **string** | Filter Instance Types by Site ID | 

 **infrastructureProviderId** | **string** | Filter Instance Types by Infrastructure Provider ID | 
 **tenantId** | **string** | Filter Instance Types by Tenant ID | 
 **status** | **string** | Filter Instance Types by Status | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, display name, description and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **includeMachineAssignment** | **bool** | Include Machine assignments for each Instance Type. Can only be requested by Provider. | 
 **includeAllocationStats** | **bool** | Include Allocation stats. Currently can only be requested by Tenant | 
 **excludeUnallocated** | **bool** | Excludes InstanceType records that have no allocations from being returned in the result set. Currently can only be requested by Tenant. | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]InstanceType**](InstanceType.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceType

> InstanceType GetInstanceType(ctx, org, instanceTypeId).IncludeMachineAssociation(includeMachineAssociation).IncludeAllocationStats(includeAllocationStats).IncludeRelation(includeRelation).Execute()

Retrieve an Instance Type



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
	instanceTypeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Instance Type
	includeMachineAssociation := true // bool | Include Machine associations for each Instance Type. Can only be requested by Provider (optional)
	includeAllocationStats := true // bool | Include Allocation stats. Currently can only be requested by Tenant (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceTypeAPI.GetInstanceType(context.Background(), org, instanceTypeId).IncludeMachineAssociation(includeMachineAssociation).IncludeAllocationStats(includeAllocationStats).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceTypeAPI.GetInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceType`: InstanceType
	fmt.Fprintf(os.Stdout, "Response from `InstanceTypeAPI.GetInstanceType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceTypeId** | **string** | ID of the Instance Type | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeMachineAssociation** | **bool** | Include Machine associations for each Instance Type. Can only be requested by Provider | 
 **includeAllocationStats** | **bool** | Include Allocation stats. Currently can only be requested by Tenant | 
 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**InstanceType**](InstanceType.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceTypeMachineAssociation

> []MachineInstanceType GetInstanceTypeMachineAssociation(ctx, org, instanceTypeId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Machines/Instance Type associations



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
	instanceTypeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Instance Type
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceTypeAPI.GetInstanceTypeMachineAssociation(context.Background(), org, instanceTypeId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceTypeAPI.GetInstanceTypeMachineAssociation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceTypeMachineAssociation`: []MachineInstanceType
	fmt.Fprintf(os.Stdout, "Response from `InstanceTypeAPI.GetInstanceTypeMachineAssociation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceTypeId** | **string** | ID of the Instance Type | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceTypeMachineAssociationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]MachineInstanceType**](MachineInstanceType.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateInstanceType

> InstanceType UpdateInstanceType(ctx, org, instanceTypeId).InstanceTypeUpdateRequest(instanceTypeUpdateRequest).Execute()

Update Instance Type



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
	instanceTypeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Instance Type
	instanceTypeUpdateRequest := *openapiclient.NewInstanceTypeUpdateRequest() // InstanceTypeUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceTypeAPI.UpdateInstanceType(context.Background(), org, instanceTypeId).InstanceTypeUpdateRequest(instanceTypeUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceTypeAPI.UpdateInstanceType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateInstanceType`: InstanceType
	fmt.Fprintf(os.Stdout, "Response from `InstanceTypeAPI.UpdateInstanceType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceTypeId** | **string** | ID of the Instance Type | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInstanceTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **instanceTypeUpdateRequest** | [**InstanceTypeUpdateRequest**](InstanceTypeUpdateRequest.md) |  | 

### Return type

[**InstanceType**](InstanceType.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

