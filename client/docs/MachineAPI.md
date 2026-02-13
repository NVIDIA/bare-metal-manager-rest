# \MachineAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteMachine**](MachineAPI.md#DeleteMachine) | **Delete** /v2/org/{org}/carbide/machine/{machineId} | Delete a Machine from a Site
[**GetAllMachine**](MachineAPI.md#GetAllMachine) | **Get** /v2/org/{org}/carbide/machine | Retrieve all Machines
[**GetAllMachineCapabilities**](MachineAPI.md#GetAllMachineCapabilities) | **Get** /v2/org/{org}/carbide/machine-capability | Retrieve all Machine Capabilities
[**GetMachine**](MachineAPI.md#GetMachine) | **Get** /v2/org/{org}/carbide/machine/{machineId} | Retrieve a Machine
[**GetMachineStatusHistory**](MachineAPI.md#GetMachineStatusHistory) | **Get** /v2/org/{org}/carbide/machine/{machineId}/status-history | Retrieve Machine status history
[**UpdateMachine**](MachineAPI.md#UpdateMachine) | **Patch** /v2/org/{org}/carbide/machine/{machineId} | Update Machine



## DeleteMachine

> DeleteMachine(ctx, org, machineId).Execute()

Delete a Machine from a Site



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
	machineId := "machineId_example" // string | ID of the Machine

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MachineAPI.DeleteMachine(context.Background(), org, machineId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MachineAPI.DeleteMachine``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**machineId** | **string** | ID of the Machine | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteMachineRequest struct via the builder pattern


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


## GetAllMachine

> []Machine GetAllMachine(ctx, org).SiteId(siteId).Id(id).HasInstanceType(hasInstanceType).InstanceTypeId(instanceTypeId).IncludeMetadata(includeMetadata).Status(status).CapabilityType(capabilityType).CapabilityName(capabilityName).HwSkuDeviceType(hwSkuDeviceType).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Machines



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
	siteId := "siteId_example" // string | Filter Machines by Site ID (optional)
	id := "id_example" // string | Filter Machines by ID.  Can be specified multiple times to filter on more than one ID. (optional)
	hasInstanceType := true // bool | Filter Machines that have been assigned an Instance Type. siteId must be specified when using this param. (optional)
	instanceTypeId := "instanceTypeId_example" // string | Filter Machines by Instance Type ID.  Can be specified multiple times to filter on more than one Instance Type ID. (optional)
	includeMetadata := true // bool | Include Machine metadata e.g. BMC, DPU, GPU and Interface data. Can only be requested by Provider. (optional)
	status := "status_example" // string | Filter Machines by Status.  Can be specified multiple times to filter on more than one Status. (optional)
	capabilityType := "capabilityType_example" // string | Filter Machines by Capability Type (optional)
	capabilityName := "capabilityName_example" // string | Filter Machines by Capability Name.  Can be specified multiple times to filter on more than one Capability Name. (optional)
	hwSkuDeviceType := "hwSkuDeviceType_example" // string | Filter Machines by hardware SKU Device Type. Example values: \"gpu\", \"cpu\", \"storage\", \"cache\" (optional)
	query := "query_example" // string | Provide query to search for matches. Input will be matched against Machine ID, vendor, product name, hostname and status (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MachineAPI.GetAllMachine(context.Background(), org).SiteId(siteId).Id(id).HasInstanceType(hasInstanceType).InstanceTypeId(instanceTypeId).IncludeMetadata(includeMetadata).Status(status).CapabilityType(capabilityType).CapabilityName(capabilityName).HwSkuDeviceType(hwSkuDeviceType).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MachineAPI.GetAllMachine``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllMachine`: []Machine
	fmt.Fprintf(os.Stdout, "Response from `MachineAPI.GetAllMachine`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllMachineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | Filter Machines by Site ID | 
 **id** | **string** | Filter Machines by ID.  Can be specified multiple times to filter on more than one ID. | 
 **hasInstanceType** | **bool** | Filter Machines that have been assigned an Instance Type. siteId must be specified when using this param. | 
 **instanceTypeId** | **string** | Filter Machines by Instance Type ID.  Can be specified multiple times to filter on more than one Instance Type ID. | 
 **includeMetadata** | **bool** | Include Machine metadata e.g. BMC, DPU, GPU and Interface data. Can only be requested by Provider. | 
 **status** | **string** | Filter Machines by Status.  Can be specified multiple times to filter on more than one Status. | 
 **capabilityType** | **string** | Filter Machines by Capability Type | 
 **capabilityName** | **string** | Filter Machines by Capability Name.  Can be specified multiple times to filter on more than one Capability Name. | 
 **hwSkuDeviceType** | **string** | Filter Machines by hardware SKU Device Type. Example values: \&quot;gpu\&quot;, \&quot;cpu\&quot;, \&quot;storage\&quot;, \&quot;cache\&quot; | 
 **query** | **string** | Provide query to search for matches. Input will be matched against Machine ID, vendor, product name, hostname and status | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]Machine**](Machine.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAllMachineCapabilities

> []MachineCapability GetAllMachineCapabilities(ctx, org).SiteId(siteId).HasInstanceType(hasInstanceType).Type_(type_).Name(name).Frequency(frequency).Capacity(capacity).Vendor(vendor).InactiveDevices(inactiveDevices).Count(count).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all Machine Capabilities



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
	siteId := "siteId_example" // string | Filter Capabilities by Machines from a particular Site
	org := "org_example" // string | Name of the Org
	hasInstanceType := true // bool | Filter Capabilities by Machines that have an Instance Type (optional)
	type_ := "type__example" // string | Filter Capabilities by Type (optional)
	name := "name_example" // string | Filter Capabilities by Name (optional)
	frequency := "frequency_example" // string | Filter Capabilities by Frequency value (optional)
	capacity := "capacity_example" // string | Filter Capabilities by Capacity value (optional)
	vendor := "vendor_example" // string | Filter Capabilities by Vendor (optional)
	inactiveDevices := "inactiveDevices_example" // string | Filter Capabilities by Inactive Devices value. Since the value is an array, multiple query params should be specified in correct order in order to filter. For example, to filter for [1, 3], specify inactiveDevices=1&inactiveDevices=3 (optional)
	count := "count_example" // string | Filter Capabilities by Count (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MachineAPI.GetAllMachineCapabilities(context.Background(), org).SiteId(siteId).HasInstanceType(hasInstanceType).Type_(type_).Name(name).Frequency(frequency).Capacity(capacity).Vendor(vendor).InactiveDevices(inactiveDevices).Count(count).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MachineAPI.GetAllMachineCapabilities``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllMachineCapabilities`: []MachineCapability
	fmt.Fprintf(os.Stdout, "Response from `MachineAPI.GetAllMachineCapabilities`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllMachineCapabilitiesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **string** | Filter Capabilities by Machines from a particular Site | 

 **hasInstanceType** | **bool** | Filter Capabilities by Machines that have an Instance Type | 
 **type_** | **string** | Filter Capabilities by Type | 
 **name** | **string** | Filter Capabilities by Name | 
 **frequency** | **string** | Filter Capabilities by Frequency value | 
 **capacity** | **string** | Filter Capabilities by Capacity value | 
 **vendor** | **string** | Filter Capabilities by Vendor | 
 **inactiveDevices** | **string** | Filter Capabilities by Inactive Devices value. Since the value is an array, multiple query params should be specified in correct order in order to filter. For example, to filter for [1, 3], specify inactiveDevices&#x3D;1&amp;inactiveDevices&#x3D;3 | 
 **count** | **string** | Filter Capabilities by Count | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]MachineCapability**](MachineCapability.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMachine

> Machine GetMachine(ctx, org, machineId).IncludeRelation(includeRelation).IncludeMetadata(includeMetadata).Execute()

Retrieve a Machine



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
	machineId := "machineId_example" // string | ID of the Machine
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	includeMetadata := true // bool | Include Machine metadata e.g. BMC, DPU, GPU and Interface data. Can only be requested by Provider. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MachineAPI.GetMachine(context.Background(), org, machineId).IncludeRelation(includeRelation).IncludeMetadata(includeMetadata).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MachineAPI.GetMachine``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMachine`: Machine
	fmt.Fprintf(os.Stdout, "Response from `MachineAPI.GetMachine`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**machineId** | **string** | ID of the Machine | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetMachineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **includeRelation** | **string** | Related entity to expand | 
 **includeMetadata** | **bool** | Include Machine metadata e.g. BMC, DPU, GPU and Interface data. Can only be requested by Provider. | 

### Return type

[**Machine**](Machine.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMachineStatusHistory

> []StatusDetail GetMachineStatusHistory(ctx, org, machineId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve Machine status history



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
	machineId := "machineId_example" // string | ID of the Machine
	pageNumber := int32(56) // int32 | Page number for pagination query (optional)
	pageSize := int32(56) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MachineAPI.GetMachineStatusHistory(context.Background(), org, machineId).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MachineAPI.GetMachineStatusHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMachineStatusHistory`: []StatusDetail
	fmt.Fprintf(os.Stdout, "Response from `MachineAPI.GetMachineStatusHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**machineId** | **string** | ID of the Machine | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetMachineStatusHistoryRequest struct via the builder pattern


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


## UpdateMachine

> Machine UpdateMachine(ctx, org, machineId).MachineUpdateRequest(machineUpdateRequest).Execute()

Update Machine



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
	machineId := "machineId_example" // string | ID of the Machine
	machineUpdateRequest := *openapiclient.NewMachineUpdateRequest() // MachineUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MachineAPI.UpdateMachine(context.Background(), org, machineId).MachineUpdateRequest(machineUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MachineAPI.UpdateMachine``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMachine`: Machine
	fmt.Fprintf(os.Stdout, "Response from `MachineAPI.UpdateMachine`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**machineId** | **string** | ID of the Machine | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMachineRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **machineUpdateRequest** | [**MachineUpdateRequest**](MachineUpdateRequest.md) |  | 

### Return type

[**Machine**](Machine.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

