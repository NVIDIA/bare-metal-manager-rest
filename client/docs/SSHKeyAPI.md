# \SSHKeyAPI

All URIs are relative to *https://carbide-rest-api.carbide.svc.cluster.local*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateSshKey**](SSHKeyAPI.md#CreateSshKey) | **Post** /v2/org/{org}/carbide/sshkey | Create SSH Key
[**DeleteSshKey**](SSHKeyAPI.md#DeleteSshKey) | **Delete** /v2/org/{org}/carbide/sshkey/{sshKeyId} | Delete an SSH Key
[**GetAllSshKey**](SSHKeyAPI.md#GetAllSshKey) | **Get** /v2/org/{org}/carbide/sshkey | Retrieve all SSH Keys
[**GetSshKey**](SSHKeyAPI.md#GetSshKey) | **Get** /v2/org/{org}/carbide/sshkey/{sshKeyId} | Retrieve an SSH key
[**UpdateSshKey**](SSHKeyAPI.md#UpdateSshKey) | **Patch** /v2/org/{org}/carbide/sshkey/{sshKeyId} | Update an SSH Key



## CreateSshKey

> SshKey CreateSshKey(ctx, org).SshKeyCreateRequest(sshKeyCreateRequest).Execute()

Create SSH Key



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
	sshKeyCreateRequest := *openapiclient.NewSshKeyCreateRequest("Name_example", "PublicKey_example") // SshKeyCreateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SSHKeyAPI.CreateSshKey(context.Background(), org).SshKeyCreateRequest(sshKeyCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyAPI.CreateSshKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateSshKey`: SshKey
	fmt.Fprintf(os.Stdout, "Response from `SSHKeyAPI.CreateSshKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateSshKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **sshKeyCreateRequest** | [**SshKeyCreateRequest**](SshKeyCreateRequest.md) |  | 

### Return type

[**SshKey**](SshKey.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteSshKey

> DeleteSshKey(ctx, org, sshKeyId).Execute()

Delete an SSH Key



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
	sshKeyId := "sshKeyId_example" // string | ID of the SSH Key

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SSHKeyAPI.DeleteSshKey(context.Background(), org, sshKeyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyAPI.DeleteSshKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**sshKeyId** | **string** | ID of the SSH Key | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteSshKeyRequest struct via the builder pattern


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


## GetAllSshKey

> []SshKey GetAllSshKey(ctx, org).SshKeyGroupId(sshKeyGroupId).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()

Retrieve all SSH Keys



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
	sshKeyGroupId := "sshKeyGroupId_example" // string | ID of the SSH Key Group (optional)
	query := "query_example" // string | Search for matches across all Sites. Input will be matched against name field (optional)
	includeRelation := "includeRelation_example" // string | Related entity to expand (optional)
	pageNumber := int32(1) // int32 | Page number for pagination query (optional) (default to 1)
	pageSize := int32(20) // int32 | Page size for pagination query (optional)
	orderBy := "orderBy_example" // string | Ordering for pagination query (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SSHKeyAPI.GetAllSshKey(context.Background(), org).SshKeyGroupId(sshKeyGroupId).Query(query).IncludeRelation(includeRelation).PageNumber(pageNumber).PageSize(pageSize).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyAPI.GetAllSshKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAllSshKey`: []SshKey
	fmt.Fprintf(os.Stdout, "Response from `SSHKeyAPI.GetAllSshKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAllSshKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **sshKeyGroupId** | **string** | ID of the SSH Key Group | 
 **query** | **string** | Search for matches across all Sites. Input will be matched against name field | 
 **includeRelation** | **string** | Related entity to expand | 
 **pageNumber** | **int32** | Page number for pagination query | [default to 1]
 **pageSize** | **int32** | Page size for pagination query | 
 **orderBy** | **string** | Ordering for pagination query | 

### Return type

[**[]SshKey**](SshKey.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSshKey

> SshKey GetSshKey(ctx, org, sshKeyId).Execute()

Retrieve an SSH key



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
	sshKeyId := "sshKeyId_example" // string | ID of the SSH Key

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SSHKeyAPI.GetSshKey(context.Background(), org, sshKeyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyAPI.GetSshKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSshKey`: SshKey
	fmt.Fprintf(os.Stdout, "Response from `SSHKeyAPI.GetSshKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**sshKeyId** | **string** | ID of the SSH Key | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSshKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**SshKey**](SshKey.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateSshKey

> SshKey UpdateSshKey(ctx, org, sshKeyId).SshKeyUpdateRequest(sshKeyUpdateRequest).Execute()

Update an SSH Key



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
	sshKeyId := "sshKeyId_example" // string | ID of the SSH Key
	sshKeyUpdateRequest := *openapiclient.NewSshKeyUpdateRequest() // SshKeyUpdateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SSHKeyAPI.UpdateSshKey(context.Background(), org, sshKeyId).SshKeyUpdateRequest(sshKeyUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SSHKeyAPI.UpdateSshKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateSshKey`: SshKey
	fmt.Fprintf(os.Stdout, "Response from `SSHKeyAPI.UpdateSshKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**org** | **string** | Name of the Org | 
**sshKeyId** | **string** | ID of the SSH Key | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateSshKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **sshKeyUpdateRequest** | [**SshKeyUpdateRequest**](SshKeyUpdateRequest.md) |  | 

### Return type

[**SshKey**](SshKey.md)

### Authorization

[JWTBearerToken](../README.md#JWTBearerToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

