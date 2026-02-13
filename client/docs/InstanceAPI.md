# \InstanceAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BatchCreateInstance**](InstanceAPI.md#BatchCreateInstance) | **Post** /v2/org/{org}/carbide/instance/batch | Batch Create Instances
[**CreateInstance**](InstanceAPI.md#CreateInstance) | **Post** /v2/org/{org}/carbide/instance | Create an Instance
[**DeleteInstance**](InstanceAPI.md#DeleteInstance) | **Delete** /v2/org/{org}/carbide/instance/{instanceId} | Delete Instance
[**GetAllInfinibandInterface**](InstanceAPI.md#GetAllInfinibandInterface) | **Get** /v2/org/{org}/carbide/instance/{instanceId}/infiniband-interface | Retrieve all InfiniBand Interfaces
[**GetAllInstance**](InstanceAPI.md#GetAllInstance) | **Get** /v2/org/{org}/carbide/instance | Retrieve all Instances
[**GetAllInterface**](InstanceAPI.md#GetAllInterface) | **Get** /v2/org/{org}/carbide/instance/{instanceId}/interface | Retrieve all Interfaces
[**GetInstance**](InstanceAPI.md#GetInstance) | **Get** /v2/org/{org}/carbide/instance/{instanceId} | Retrieve Instance
[**GetInstanceStatusHistory**](InstanceAPI.md#GetInstanceStatusHistory) | **Get** /v2/org/{org}/carbide/instance/{instanceId}/status-history | Retrieve Instance status history
[**UpdateInstance**](InstanceAPI.md#UpdateInstance) | **Patch** /v2/org/{org}/carbide/instance/{instanceId} | Update Instance



## BatchCreateInstance

> []Instance BatchCreateInstance(ctx, org).BatchInstanceCreateRequest(batchInstanceCreateRequest).Execute()

Batch Create Instances



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
	batchInstanceCreateRequest := *openapiclient.NewBatchInstanceCreateRequest("NamePrefix_example", int32(123), "TenantId_example", "InstanceTypeId_example", "VpcId_example", []openapiclient.InterfaceCreateRequest{*openapiclient.NewInterfaceCreateRequest()}) // BatchInstanceCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceAPI.BatchCreateInstance(context.Background(), org).BatchInstanceCreateRequest(batchInstanceCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceAPI.BatchCreateInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BatchCreateInstance`: []Instance
	fmt.Fprintf(os.Stdout, "Response from `InstanceAPI.BatchCreateInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiBatchCreateInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **batchInstanceCreateRequest** | [**BatchInstanceCreateRequest**](BatchInstanceCreateRequest.md) |  | 

### Return type

[**[]Instance**](Instance.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateInstance

> Instance CreateInstance(ctx, org).InstanceCreateRequest(instanceCreateRequest).Execute()

Create an Instance



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
	instanceCreateRequest := *openapiclient.NewInstanceCreateRequest("Name_example", "TenantId_example", "VpcId_example", []openapiclient.InterfaceCreateRequest{*openapiclient.NewInterfaceCreateRequest()}) // InstanceCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceAPI.CreateInstance(context.Background(), org).InstanceCreateRequest(instanceCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceAPI.CreateInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateInstance`: Instance
	fmt.Fprintf(os.Stdout, "Response from `InstanceAPI.CreateInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **instanceCreateRequest** | [**InstanceCreateRequest**](InstanceCreateRequest.md) |  | 

### Return type

[**Instance**](Instance.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteInstance

> DeleteInstance(ctx, org, instanceId).InstanceDeleteRequest(instanceDeleteRequest).Execute()

Delete Instance



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
	instanceId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Instance
	instanceDeleteRequest := *openapiclient.NewInstanceDeleteRequest() // InstanceDeleteRequest | Optional request data to report health issues with the underlying Machine (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.InstanceAPI.DeleteInstance(context.Background(), org, instanceId).InstanceDeleteRequest(instanceDeleteRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceAPI.DeleteInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceId** | **string** | ID of the Instance | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **instanceDeleteRequest** | [**InstanceDeleteRequest**](InstanceDeleteRequest.md) | Optional request data to report health issues with the underlying Machine | 

### Return type

 (empty response body)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllInfinibandInterface

> []InfiniBandInterface GetAllInfinibandInterface(ctx, org, instanceId).Status(status).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all InfiniBand Interfaces



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
	instanceId := "instanceId_example" // string | ID of the Instance
	status := "status_example" // string | Filter Interfaces by Status (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceAPI.GetAllInfinibandInterface(context.Background(), org, instanceId).Status(status).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceAPI.GetAllInfinibandInterface``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllInfinibandInterface`: []InfiniBandInterface
	fmt.Fprintf(os.Stdout, "Response from `InstanceAPI.GetAllInfinibandInterface`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceId** | **string** | ID of the Instance | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllInfinibandInterfaceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **status** | **string** | Filter Interfaces by Status | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]InfiniBandInterface**](InfiniBandInterface.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllInstance

> []Instance GetAllInstance(ctx, org).InfrastructureProviderId(infrastructureProviderId).SiteId(siteId).VpcId(vpcId).InstanceTypeId(instanceTypeId).OperatingSystemId(operatingSystemId).MachineId(machineId).Name(name).Status(status).IpAddress(ipAddress).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).NetworkSecurityGroupId(networkSecurityGroupId).Execute()

Retrieve all Instances



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
	infrastructureProviderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter by Infrastructure Provider ID (optional)
	siteId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter by Site ID. Can be specified multiple times to filter on more than one site. (optional)
	vpcId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter by VPC ID. Can be specified multiple times to filter on more than one VPC. (optional)
	instanceTypeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter by instance type ID. Can be specified multiple times to filter on more than one instance type. (optional)
	operatingSystemId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter by operating system ID. Can be specified multiple times to filter on more than one operating system. (optional)
	machineId := "machineId_example" // string | Filter by machine ID. Can be specified multiple times to filter on more than one machine. (optional)
	name := "name_example" // string | Filter by Instance name (optional)
	status := "status_example" // string | Filter Instances by Status. Can be specified multiple times to filter on more than one status. (optional)
	ipAddress := "ipAddress_example" // string | Filter by IP address. Can be specified multiple times to filter on more than one IP address. (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name, description, status, and labels fields (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)
	networkSecurityGroupId := "networkSecurityGroupId_example" // string | Filter by NetworkSecurityGroup ID.  Can be specified multiple times to filter on more than one Network Security Group. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceAPI.GetAllInstance(context.Background(), org).InfrastructureProviderId(infrastructureProviderId).SiteId(siteId).VpcId(vpcId).InstanceTypeId(instanceTypeId).OperatingSystemId(operatingSystemId).MachineId(machineId).Name(name).Status(status).IpAddress(ipAddress).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).NetworkSecurityGroupId(networkSecurityGroupId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceAPI.GetAllInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllInstance`: []Instance
	fmt.Fprintf(os.Stdout, "Response from `InstanceAPI.GetAllInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **infrastructureProviderId** | **string** | Filter by Infrastructure Provider ID | 
 **siteId** | **string** | Filter by Site ID. Can be specified multiple times to filter on more than one site. | 
 **vpcId** | **string** | Filter by VPC ID. Can be specified multiple times to filter on more than one VPC. | 
 **instanceTypeId** | **string** | Filter by instance type ID. Can be specified multiple times to filter on more than one instance type. | 
 **operatingSystemId** | **string** | Filter by operating system ID. Can be specified multiple times to filter on more than one operating system. | 
 **machineId** | **string** | Filter by machine ID. Can be specified multiple times to filter on more than one machine. | 
 **name** | **string** | Filter by Instance name | 
 **status** | **string** | Filter Instances by Status. Can be specified multiple times to filter on more than one status. | 
 **ipAddress** | **string** | Filter by IP address. Can be specified multiple times to filter on more than one IP address. | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name, description, status, and labels fields | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 
 **networkSecurityGroupId** | **string** | Filter by NetworkSecurityGroup ID.  Can be specified multiple times to filter on more than one Network Security Group. | 

### Return type

[**[]Instance**](Instance.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllInterface

> []Interface GetAllInterface(ctx, org, instanceId).Status(status).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Interfaces



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
	instanceId := "instanceId_example" // string | ID of the Instance
	status := "status_example" // string | Filter Interfaces by Status (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceAPI.GetAllInterface(context.Background(), org, instanceId).Status(status).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceAPI.GetAllInterface``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllInterface`: []Interface
	fmt.Fprintf(os.Stdout, "Response from `InstanceAPI.GetAllInterface`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceId** | **string** | ID of the Instance | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllInterfaceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **status** | **string** | Filter Interfaces by Status | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]Interface**](Interface.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstance

> Instance GetInstance(ctx, org, instanceId).IncludeRelation(includeRelation).Execute()

Retrieve Instance



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
	instanceId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Instance
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceAPI.GetInstance(context.Background(), org, instanceId).IncludeRelation(includeRelation).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceAPI.GetInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstance`: Instance
	fmt.Fprintf(os.Stdout, "Response from `InstanceAPI.GetInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceId** | **string** | ID of the Instance | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeRelation** | **string** | Related entity to expand | 

### Return type

[**Instance**](Instance.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstanceStatusHistory

> []StatusDetail GetInstanceStatusHistory(ctx, org, instanceId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve Instance status history



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
	instanceId := "instanceId_example" // string | ID of the Instance
	pageNumber := int32(56) // int32 | Page number for pagination query (optional)
	pageSize := int32(56) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceAPI.GetInstanceStatusHistory(context.Background(), org, instanceId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceAPI.GetInstanceStatusHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInstanceStatusHistory`: []StatusDetail
	fmt.Fprintf(os.Stdout, "Response from `InstanceAPI.GetInstanceStatusHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceId** | **string** | ID of the Instance | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstanceStatusHistoryRequest struct via the builder pattern


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


## UpdateInstance

> Instance UpdateInstance(ctx, org, instanceId).InstanceUpdateRequest(instanceUpdateRequest).Execute()

Update Instance



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
	instanceId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the Instance
	instanceUpdateRequest := *openapiclient.NewInstanceUpdateRequest() // InstanceUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InstanceAPI.UpdateInstance(context.Background(), org, instanceId).InstanceUpdateRequest(instanceUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InstanceAPI.UpdateInstance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateInstance`: Instance
	fmt.Fprintf(os.Stdout, "Response from `InstanceAPI.UpdateInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**instanceId** | **string** | ID of the Instance | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **instanceUpdateRequest** | [**InstanceUpdateRequest**](InstanceUpdateRequest.md) |  | 

### Return type

[**Instance**](Instance.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

