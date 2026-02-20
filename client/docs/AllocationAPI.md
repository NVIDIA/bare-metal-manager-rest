# \AllocationAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateAllocation**](AllocationAPI.md#CreateAllocation) | **Post** /v2/org/{org}/carbide/allocation | Create Allocation
[**CreateAllocationConstraint**](AllocationAPI.md#CreateAllocationConstraint) | **Post** /v2/org/{org}/carbide/allocation/{allocationId}/constraint | Create Allocation Constraint
[**DeleteAllocation**](AllocationAPI.md#DeleteAllocation) | **Delete** /v2/org/{org}/carbide/allocation/{allocationId} | Delete Allocation
[**DeleteAllocationConstraint**](AllocationAPI.md#DeleteAllocationConstraint) | **Delete** /v2/org/{org}/carbide/allocation/{allocationId}/constraint/{allocationConstraintId} | Delete Allocation Constraint
[**GetAllAllocation**](AllocationAPI.md#GetAllAllocation) | **Get** /v2/org/{org}/carbide/allocation | Retrieve all Allocations
[**GetAllAllocationConstraint**](AllocationAPI.md#GetAllAllocationConstraint) | **Get** /v2/org/{org}/carbide/allocation/{allocationId}/constraint | Retrieve all Allocation Constraints
[**GetAllocation**](AllocationAPI.md#GetAllocation) | **Get** /v2/org/{org}/carbide/allocation/{allocationId} | Retrieve Allocation
[**GetAllocationConstraint**](AllocationAPI.md#GetAllocationConstraint) | **Get** /v2/org/{org}/carbide/allocation/{allocationId}/constraint/{allocationConstraintId} | Retrieve Allocation Constraint
[**UpdateAllocation**](AllocationAPI.md#UpdateAllocation) | **Patch** /v2/org/{org}/carbide/allocation/{allocationId} | Update Allocation
[**UpdateAllocationConstraint**](AllocationAPI.md#UpdateAllocationConstraint) | **Patch** /v2/org/{org}/carbide/allocation/{allocationId}/constraint/{allocationConstraintId} | Update Allocation Constraint



## CreateAllocation

> Allocation CreateAllocation(ctx, org).AllocationCreateRequest(allocationCreateRequest).Execute()

Create Allocation



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
	allocationCreateRequest := *openapiclient.NewAllocationCreateRequest("Name_example", "TenantId_example", "SiteId_example") // AllocationCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AllocationAPI.CreateAllocation(context.Background(), org).AllocationCreateRequest(allocationCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.CreateAllocation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateAllocation`: Allocation
	fmt.Fprintf(os.Stdout, "Response from `AllocationAPI.CreateAllocation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateAllocationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **allocationCreateRequest** | [**AllocationCreateRequest**](AllocationCreateRequest.md) |  | 

### Return type

[**Allocation**](Allocation.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAllocationConstraint

> AllocationConstraint CreateAllocationConstraint(ctx, org, allocationId).AllocationConstraintCreateRequest(allocationConstraintCreateRequest).Execute()

Create Allocation Constraint



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
	allocationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Allocation
	allocationConstraintCreateRequest := *openapiclient.NewAllocationConstraintCreateRequest("ResourceTypeId_example", "ConstraintType_example", int32(123)) // AllocationConstraintCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AllocationAPI.CreateAllocationConstraint(context.Background(), org, allocationId).AllocationConstraintCreateRequest(allocationConstraintCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.CreateAllocationConstraint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateAllocationConstraint`: AllocationConstraint
	fmt.Fprintf(os.Stdout, "Response from `AllocationAPI.CreateAllocationConstraint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**allocationId** | **string** | ID of the Allocation | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateAllocationConstraintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **allocationConstraintCreateRequest** | [**AllocationConstraintCreateRequest**](AllocationConstraintCreateRequest.md) |  | 

### Return type

[**AllocationConstraint**](AllocationConstraint.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAllocation

> DeleteAllocation(ctx, org, allocationId).Execute()

Delete Allocation



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
	allocationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Allocation

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AllocationAPI.DeleteAllocation(context.Background(), org, allocationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.DeleteAllocation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**allocationId** | **string** | ID of the Allocation | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAllocationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



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


## DeleteAllocationConstraint

> DeleteAllocationConstraint(ctx, org, allocationId, allocationConstraintId).Execute()

Delete Allocation Constraint



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
	allocationId := "allocationId_example" // string | ID of the Allocation
	allocationConstraintId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Allocation Constraint

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AllocationAPI.DeleteAllocationConstraint(context.Background(), org, allocationId, allocationConstraintId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.DeleteAllocationConstraint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**allocationId** | **string** | ID of the Allocation | 
**allocationConstraintId** | **string** | ID of the Allocation Constraint | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAllocationConstraintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




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


## GetAllAllocation

> []Allocation GetAllAllocation(ctx, org).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).SiteId(siteId).Id(id).ResourceType(resourceType).Status(status).ResourceTypeId(resourceTypeId).ConstraintType(constraintType).ConstraintValue(constraintValue).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Allocations



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
	infrastructureProviderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Allocations by Infrastructure Provider ID (optional)
	tenantId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Allocations by Tenant ID.  Can be specified multiple times to filter on more than one Tenant ID. (optional)
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Allocations by Site ID.  Can be specified multiple times to filter on more than one Site ID. (optional)
	id := "id_example" // string | Filter Allocations by ID.  Can be specified multiple times to filter on more than one ID. (optional)
	resourceType := "resourceType_example" // string | Filter Allocations by Constraint Resource Type.  Can be specified multiple times to filter on more than one Constraint Resource Type. (optional)
	status := "status_example" // string | Filter Allocations by Status.  Can be specified multiple times to filter on more than one Status. (optional)
	resourceTypeId := "resourceTypeId_example" // string | Filter Allocations by Constraint Resource Type ID.  Can be specified multiple times to filter on more than one Constraint Resource Type ID. (optional)
	constraintType := "constraintType_example" // string | Filter Allocations by Constraint Type.  Can be specified multiple times to filter on more than one Constraint Type. (optional)
	constraintValue := int32(56) // int32 | Filter Allocations by Constraint Value.  Can be specified multiple times to filter on more than one Constraint Value. (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, description and status fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AllocationAPI.GetAllAllocation(context.Background(), org).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).SiteId(siteId).Id(id).ResourceType(resourceType).Status(status).ResourceTypeId(resourceTypeId).ConstraintType(constraintType).ConstraintValue(constraintValue).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.GetAllAllocation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllAllocation`: []Allocation
	fmt.Fprintf(os.Stdout, "Response from `AllocationAPI.GetAllAllocation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllAllocationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **infrastructureProviderId** | **string** | Filter Allocations by Infrastructure Provider ID | 
 **tenantId** | **string** | Filter Allocations by Tenant ID.  Can be specified multiple times to filter on more than one Tenant ID. | 
 **siteId** | **string** | Filter Allocations by Site ID.  Can be specified multiple times to filter on more than one Site ID. | 
 **id** | **string** | Filter Allocations by ID.  Can be specified multiple times to filter on more than one ID. | 
 **resourceType** | **string** | Filter Allocations by Constraint Resource Type.  Can be specified multiple times to filter on more than one Constraint Resource Type. | 
 **status** | **string** | Filter Allocations by Status.  Can be specified multiple times to filter on more than one Status. | 
 **resourceTypeId** | **string** | Filter Allocations by Constraint Resource Type ID.  Can be specified multiple times to filter on more than one Constraint Resource Type ID. | 
 **constraintType** | **string** | Filter Allocations by Constraint Type.  Can be specified multiple times to filter on more than one Constraint Type. | 
 **constraintValue** | **int32** | Filter Allocations by Constraint Value.  Can be specified multiple times to filter on more than one Constraint Value. | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, description and status fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]Allocation**](Allocation.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllAllocationConstraint

> []AllocationConstraint GetAllAllocationConstraint(ctx, org, allocationId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Allocation Constraints



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
	allocationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Allocation
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AllocationAPI.GetAllAllocationConstraint(context.Background(), org, allocationId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.GetAllAllocationConstraint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllAllocationConstraint`: []AllocationConstraint
	fmt.Fprintf(os.Stdout, "Response from `AllocationAPI.GetAllAllocationConstraint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**allocationId** | **string** | ID of the Allocation | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllAllocationConstraintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]AllocationConstraint**](AllocationConstraint.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllocation

> Allocation GetAllocation(ctx, org, allocationId).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).IncludeRelation(includeRelation).Execute()

Retrieve Allocation



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
	allocationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Allocation
	infrastructureProviderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Allocations by Infrastructure Provider ID (optional)
	tenantId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter Allocations by Tenant ID (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AllocationAPI.GetAllocation(context.Background(), org, allocationId).InfrastructureProviderId(infrastructureProviderId).TenantId(tenantId).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.GetAllocation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllocation`: Allocation
	fmt.Fprintf(os.Stdout, "Response from `AllocationAPI.GetAllocation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**allocationId** | **string** | ID of the Allocation | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllocationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **infrastructureProviderId** | **string** | Filter Allocations by Infrastructure Provider ID | 
 **tenantId** | **string** | Filter Allocations by Tenant ID | 
 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**Allocation**](Allocation.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllocationConstraint

> []AllocationConstraint GetAllocationConstraint(ctx, org, allocationId, allocationConstraintId).Execute()

Retrieve Allocation Constraint



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
	allocationId := "allocationId_example" // string | ID of the Allocation
	allocationConstraintId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Allocation Constraint

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AllocationAPI.GetAllocationConstraint(context.Background(), org, allocationId, allocationConstraintId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.GetAllocationConstraint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllocationConstraint`: []AllocationConstraint
	fmt.Fprintf(os.Stdout, "Response from `AllocationAPI.GetAllocationConstraint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**allocationId** | **string** | ID of the Allocation | 
**allocationConstraintId** | **string** | ID of the Allocation Constraint | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllocationConstraintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[**[]AllocationConstraint**](AllocationConstraint.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateAllocation

> Allocation UpdateAllocation(ctx, org, allocationId).AllocationUpdateRequest(allocationUpdateRequest).Execute()

Update Allocation



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
	allocationId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Allocation
	allocationUpdateRequest := *openapiclient.NewAllocationUpdateRequest() // AllocationUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AllocationAPI.UpdateAllocation(context.Background(), org, allocationId).AllocationUpdateRequest(allocationUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.UpdateAllocation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateAllocation`: Allocation
	fmt.Fprintf(os.Stdout, "Response from `AllocationAPI.UpdateAllocation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**allocationId** | **string** | ID of the Allocation | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateAllocationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **allocationUpdateRequest** | [**AllocationUpdateRequest**](AllocationUpdateRequest.md) |  | 

### Return type

[**Allocation**](Allocation.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateAllocationConstraint

> AllocationConstraint UpdateAllocationConstraint(ctx, org, allocationId, allocationConstraintId).AllocationConstraintUpdateRequest(allocationConstraintUpdateRequest).Execute()

Update Allocation Constraint



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
	allocationId := "allocationId_example" // string | ID of the Allocation
	allocationConstraintId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Allocation Constraint
	allocationConstraintUpdateRequest := *openapiclient.NewAllocationConstraintUpdateRequest(int32(123)) // AllocationConstraintUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AllocationAPI.UpdateAllocationConstraint(context.Background(), org, allocationId, allocationConstraintId).AllocationConstraintUpdateRequest(allocationConstraintUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AllocationAPI.UpdateAllocationConstraint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateAllocationConstraint`: AllocationConstraint
	fmt.Fprintf(os.Stdout, "Response from `AllocationAPI.UpdateAllocationConstraint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**allocationId** | **string** | ID of the Allocation | 
**allocationConstraintId** | **string** | ID of the Allocation Constraint | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateAllocationConstraintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **allocationConstraintUpdateRequest** | [**AllocationConstraintUpdateRequest**](AllocationConstraintUpdateRequest.md) |  | 

### Return type

[**AllocationConstraint**](AllocationConstraint.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

