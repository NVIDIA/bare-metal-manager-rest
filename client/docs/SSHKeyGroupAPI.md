# \SSHKeyGroupAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateSshKeyGroup**](SSHKeyGroupAPI.md#CreateSshKeyGroup) | **Post** /v2/org/{org}/carbide/sshkeygroup | Create SSH Key Group
[**DeleteSshKeyGroup**](SSHKeyGroupAPI.md#DeleteSshKeyGroup) | **Delete** /v2/org/{org}/carbide/sshkeygroup/{sshKeyGroupId} | Delete an SSH Key Group
[**GetAllSshKeyGroup**](SSHKeyGroupAPI.md#GetAllSshKeyGroup) | **Get** /v2/org/{org}/carbide/sshkeygroup | Retrieve all SSH Key Groups
[**GetSshKeyGroup**](SSHKeyGroupAPI.md#GetSshKeyGroup) | **Get** /v2/org/{org}/carbide/sshkeygroup/{sshKeyGroupId} | Retrieve an SSH Key Group
[**UpdateSshKeyGroup**](SSHKeyGroupAPI.md#UpdateSshKeyGroup) | **Patch** /v2/org/{org}/carbide/sshkeygroup/{sshKeyGroupId} | Update an SSH Key Group



## CreateSshKeyGroup

> SshKeyGroup CreateSshKeyGroup(ctx, org).SshKeyGroupCreateRequest(sshKeyGroupCreateRequest).Execute()

Create SSH Key Group



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
	sshKeyGroupCreateRequest := *openapiclient.NewSshKeyGroupCreateRequest("Name_example") // SshKeyGroupCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SSHKeyGroupAPI.CreateSshKeyGroup(context.Background(), org).SshKeyGroupCreateRequest(sshKeyGroupCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyGroupAPI.CreateSshKeyGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateSshKeyGroup`: SshKeyGroup
	fmt.Fprintf(os.Stdout, "Response from `SSHKeyGroupAPI.CreateSshKeyGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateSshKeyGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **sshKeyGroupCreateRequest** | [**SshKeyGroupCreateRequest**](SshKeyGroupCreateRequest.md) |  | 

### Return type

[**SshKeyGroup**](SshKeyGroup.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteSshKeyGroup

> DeleteSshKeyGroup(ctx, org, sshKeyGroupId).Execute()

Delete an SSH Key Group



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
	sshKeyGroupId := "sshKeyGroupId_example" // string | ID of the SSH Key Group

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SSHKeyGroupAPI.DeleteSshKeyGroup(context.Background(), org, sshKeyGroupId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyGroupAPI.DeleteSshKeyGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**sshKeyGroupId** | **string** | ID of the SSH Key Group | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteSshKeyGroupRequest struct via the builder pattern


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


## GetAllSshKeyGroup

> []SshKeyGroup GetAllSshKeyGroup(ctx, org).SiteId(siteId).InstanceId(instanceId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all SSH Key Groups



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
	siteId := "siteId_example" // string | Filter SSH Key Groups by Site ID (optional)
	instanceId := "instanceId_example" // string | Filter SSH Key Groups by Instance ID (optional)
	status := "status_example" // string | Status filter for the SSH Key Groups (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name field (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SSHKeyGroupAPI.GetAllSshKeyGroup(context.Background(), org).SiteId(siteId).InstanceId(instanceId).Status(status).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyGroupAPI.GetAllSshKeyGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllSshKeyGroup`: []SshKeyGroup
	fmt.Fprintf(os.Stdout, "Response from `SSHKeyGroupAPI.GetAllSshKeyGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllSshKeyGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **siteId** | **string** | Filter SSH Key Groups by Site ID | 
 **instanceId** | **string** | Filter SSH Key Groups by Instance ID | 
 **status** | **string** | Status filter for the SSH Key Groups | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name field | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]SshKeyGroup**](SshKeyGroup.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSshKeyGroup

> SshKeyGroup GetSshKeyGroup(ctx, org, sshKeyGroupId).Execute()

Retrieve an SSH Key Group



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
	sshKeyGroupId := "sshKeyGroupId_example" // string | ID of the SSH Key Group

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SSHKeyGroupAPI.GetSshKeyGroup(context.Background(), org, sshKeyGroupId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyGroupAPI.GetSshKeyGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSshKeyGroup`: SshKeyGroup
	fmt.Fprintf(os.Stdout, "Response from `SSHKeyGroupAPI.GetSshKeyGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**sshKeyGroupId** | **string** | ID of the SSH Key Group | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSshKeyGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**SshKeyGroup**](SshKeyGroup.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateSshKeyGroup

> SshKeyGroup UpdateSshKeyGroup(ctx, org, sshKeyGroupId).SshKeyGroupUpdateRequest(sshKeyGroupUpdateRequest).Execute()

Update an SSH Key Group



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
	sshKeyGroupId := "sshKeyGroupId_example" // string | ID of the SSH Key Group
	sshKeyGroupUpdateRequest := *openapiclient.NewSshKeyGroupUpdateRequest("Version_example") // SshKeyGroupUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SSHKeyGroupAPI.UpdateSshKeyGroup(context.Background(), org, sshKeyGroupId).SshKeyGroupUpdateRequest(sshKeyGroupUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyGroupAPI.UpdateSshKeyGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateSshKeyGroup`: SshKeyGroup
	fmt.Fprintf(os.Stdout, "Response from `SSHKeyGroupAPI.UpdateSshKeyGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**sshKeyGroupId** | **string** | ID of the SSH Key Group | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateSshKeyGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **sshKeyGroupUpdateRequest** | [**SshKeyGroupUpdateRequest**](SshKeyGroupUpdateRequest.md) |  | 

### Return type

[**SshKeyGroup**](SshKeyGroup.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

