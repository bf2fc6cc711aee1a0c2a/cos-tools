/*
Connector Service Fleet Manager

Connector Service Fleet Manager is a Rest API to manage connectors.

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// ConnectorClustersApiService ConnectorClustersApi service
type ConnectorClustersApiService service

type ApiCreateConnectorClusterRequest struct {
	ctx                     context.Context
	ApiService              *ConnectorClustersApiService
	async                   *bool
	connectorClusterRequest *ConnectorClusterRequest
}

// Perform the action in an asynchronous manner
func (r ApiCreateConnectorClusterRequest) Async(async bool) ApiCreateConnectorClusterRequest {
	r.async = &async
	return r
}

// Connector cluster data
func (r ApiCreateConnectorClusterRequest) ConnectorClusterRequest(connectorClusterRequest ConnectorClusterRequest) ApiCreateConnectorClusterRequest {
	r.connectorClusterRequest = &connectorClusterRequest
	return r
}

func (r ApiCreateConnectorClusterRequest) Execute() (*ConnectorCluster, *http.Response, error) {
	return r.ApiService.CreateConnectorClusterExecute(r)
}

/*
CreateConnectorCluster Create a new connector cluster

Create a new connector cluster

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiCreateConnectorClusterRequest
*/
func (a *ConnectorClustersApiService) CreateConnectorCluster(ctx context.Context) ApiCreateConnectorClusterRequest {
	return ApiCreateConnectorClusterRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//  @return ConnectorCluster
func (a *ConnectorClustersApiService) CreateConnectorClusterExecute(r ApiCreateConnectorClusterRequest) (*ConnectorCluster, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *ConnectorCluster
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ConnectorClustersApiService.CreateConnectorCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/connector_mgmt/v1/kafka_connector_clusters"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.async == nil {
		return localVarReturnValue, nil, reportError("async is required and must be specified")
	}
	if r.connectorClusterRequest == nil {
		return localVarReturnValue, nil, reportError("connectorClusterRequest is required and must be specified")
	}

	localVarQueryParams.Add("async", parameterToString(*r.async, ""))
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.connectorClusterRequest
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiDeleteConnectorClusterRequest struct {
	ctx                context.Context
	ApiService         *ConnectorClustersApiService
	connectorClusterId string
}

func (r ApiDeleteConnectorClusterRequest) Execute() (*Error, *http.Response, error) {
	return r.ApiService.DeleteConnectorClusterExecute(r)
}

/*
DeleteConnectorCluster Delete a connector cluster

Delete a connector cluster

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param connectorClusterId The id of the connector cluster
 @return ApiDeleteConnectorClusterRequest
*/
func (a *ConnectorClustersApiService) DeleteConnectorCluster(ctx context.Context, connectorClusterId string) ApiDeleteConnectorClusterRequest {
	return ApiDeleteConnectorClusterRequest{
		ApiService:         a,
		ctx:                ctx,
		connectorClusterId: connectorClusterId,
	}
}

// Execute executes the request
//  @return Error
func (a *ConnectorClustersApiService) DeleteConnectorClusterExecute(r ApiDeleteConnectorClusterRequest) (*Error, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Error
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ConnectorClustersApiService.DeleteConnectorCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/connector_mgmt/v1/kafka_connector_clusters/{connector_cluster_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"connector_cluster_id"+"}", url.PathEscape(parameterToString(r.connectorClusterId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetConnectorClusterRequest struct {
	ctx                context.Context
	ApiService         *ConnectorClustersApiService
	connectorClusterId string
}

func (r ApiGetConnectorClusterRequest) Execute() (*ConnectorCluster, *http.Response, error) {
	return r.ApiService.GetConnectorClusterExecute(r)
}

/*
GetConnectorCluster Get a connector cluster

Get a connector cluster

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param connectorClusterId The id of the connector cluster
 @return ApiGetConnectorClusterRequest
*/
func (a *ConnectorClustersApiService) GetConnectorCluster(ctx context.Context, connectorClusterId string) ApiGetConnectorClusterRequest {
	return ApiGetConnectorClusterRequest{
		ApiService:         a,
		ctx:                ctx,
		connectorClusterId: connectorClusterId,
	}
}

// Execute executes the request
//  @return ConnectorCluster
func (a *ConnectorClustersApiService) GetConnectorClusterExecute(r ApiGetConnectorClusterRequest) (*ConnectorCluster, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *ConnectorCluster
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ConnectorClustersApiService.GetConnectorCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/connector_mgmt/v1/kafka_connector_clusters/{connector_cluster_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"connector_cluster_id"+"}", url.PathEscape(parameterToString(r.connectorClusterId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 410 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetConnectorClusterAddonParametersRequest struct {
	ctx                context.Context
	ApiService         *ConnectorClustersApiService
	connectorClusterId string
	resetCredentials   *bool
}

// Resets cluster service account credentials when true
func (r ApiGetConnectorClusterAddonParametersRequest) ResetCredentials(resetCredentials bool) ApiGetConnectorClusterAddonParametersRequest {
	r.resetCredentials = &resetCredentials
	return r
}

func (r ApiGetConnectorClusterAddonParametersRequest) Execute() ([]AddonParameter, *http.Response, error) {
	return r.ApiService.GetConnectorClusterAddonParametersExecute(r)
}

/*
GetConnectorClusterAddonParameters Get a connector cluster's addon parameters

Get a connector cluster's addon parameters

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param connectorClusterId The id of the connector cluster
 @return ApiGetConnectorClusterAddonParametersRequest
*/
func (a *ConnectorClustersApiService) GetConnectorClusterAddonParameters(ctx context.Context, connectorClusterId string) ApiGetConnectorClusterAddonParametersRequest {
	return ApiGetConnectorClusterAddonParametersRequest{
		ApiService:         a,
		ctx:                ctx,
		connectorClusterId: connectorClusterId,
	}
}

// Execute executes the request
//  @return []AddonParameter
func (a *ConnectorClustersApiService) GetConnectorClusterAddonParametersExecute(r ApiGetConnectorClusterAddonParametersRequest) ([]AddonParameter, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue []AddonParameter
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ConnectorClustersApiService.GetConnectorClusterAddonParameters")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/connector_mgmt/v1/kafka_connector_clusters/{connector_cluster_id}/addon_parameters"
	localVarPath = strings.Replace(localVarPath, "{"+"connector_cluster_id"+"}", url.PathEscape(parameterToString(r.connectorClusterId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.resetCredentials != nil {
		localVarQueryParams.Add("reset_credentials", parameterToString(*r.resetCredentials, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 410 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetConnectorClusterNamespacesRequest struct {
	ctx                context.Context
	ApiService         *ConnectorClustersApiService
	connectorClusterId string
	page               *string
	size               *string
	orderBy            *string
	search             *string
}

// Page index
func (r ApiGetConnectorClusterNamespacesRequest) Page(page string) ApiGetConnectorClusterNamespacesRequest {
	r.page = &page
	return r
}

// Number of items in each page
func (r ApiGetConnectorClusterNamespacesRequest) Size(size string) ApiGetConnectorClusterNamespacesRequest {
	r.size = &size
	return r
}

// Specifies the order by criteria. The syntax of this parameter is similar to the syntax of the &#x60;order by&#x60; clause of an SQL statement. Each query can be ordered by any of the &#x60;ConnectorType&#x60; fields. For example, to return all Connector types ordered by their name, use the following syntax:  &#x60;&#x60;&#x60;sql name asc &#x60;&#x60;&#x60;  To return all Connector types ordered by their name _and_ version, use the following syntax:  &#x60;&#x60;&#x60;sql name asc, version asc &#x60;&#x60;&#x60;  If the parameter isn&#39;t provided, or if the value is empty, then the results are ordered by name.
func (r ApiGetConnectorClusterNamespacesRequest) OrderBy(orderBy string) ApiGetConnectorClusterNamespacesRequest {
	r.orderBy = &orderBy
	return r
}

// Search criteria.  The syntax of this parameter is similar to the syntax of the &#x60;where&#x60; clause of a SQL statement. Allowed fields in the search are &#x60;name&#x60;, &#x60;description&#x60;, &#x60;version&#x60;, &#x60;label&#x60;, and &#x60;channel&#x60;. Allowed operators are &#x60;&lt;&gt;&#x60;, &#x60;&#x3D;&#x60;, or &#x60;LIKE&#x60;. Allowed conjunctive operators are &#x60;AND&#x60; and &#x60;OR&#x60;. However, you can use a maximum of 10 conjunctions in a search query.  Examples:  To return a Connector Type with the name &#x60;aws-sqs-source&#x60; and the channel &#x60;stable&#x60;, use the following syntax:  &#x60;&#x60;&#x60; name &#x3D; aws-sqs-source and channel &#x3D; stable &#x60;&#x60;&#x60;[p-]  To return a Kafka instance with a name that starts with &#x60;aws&#x60;, use the following syntax:  &#x60;&#x60;&#x60; name like aws%25 &#x60;&#x60;&#x60;  If the parameter isn&#39;t provided, or if the value is empty, then all the Connector Type that the user has permission to see are returned.  Note. If the query is invalid, an error is returned.
func (r ApiGetConnectorClusterNamespacesRequest) Search(search string) ApiGetConnectorClusterNamespacesRequest {
	r.search = &search
	return r
}

func (r ApiGetConnectorClusterNamespacesRequest) Execute() (*ConnectorNamespaceList, *http.Response, error) {
	return r.ApiService.GetConnectorClusterNamespacesExecute(r)
}

/*
GetConnectorClusterNamespaces Get a connector cluster's namespaces

Get a connector cluster's namespaces

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param connectorClusterId The id of the connector cluster
 @return ApiGetConnectorClusterNamespacesRequest
*/
func (a *ConnectorClustersApiService) GetConnectorClusterNamespaces(ctx context.Context, connectorClusterId string) ApiGetConnectorClusterNamespacesRequest {
	return ApiGetConnectorClusterNamespacesRequest{
		ApiService:         a,
		ctx:                ctx,
		connectorClusterId: connectorClusterId,
	}
}

// Execute executes the request
//  @return ConnectorNamespaceList
func (a *ConnectorClustersApiService) GetConnectorClusterNamespacesExecute(r ApiGetConnectorClusterNamespacesRequest) (*ConnectorNamespaceList, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *ConnectorNamespaceList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ConnectorClustersApiService.GetConnectorClusterNamespaces")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/connector_mgmt/v1/kafka_connector_clusters/{connector_cluster_id}/namespaces"
	localVarPath = strings.Replace(localVarPath, "{"+"connector_cluster_id"+"}", url.PathEscape(parameterToString(r.connectorClusterId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.page != nil {
		localVarQueryParams.Add("page", parameterToString(*r.page, ""))
	}
	if r.size != nil {
		localVarQueryParams.Add("size", parameterToString(*r.size, ""))
	}
	if r.orderBy != nil {
		localVarQueryParams.Add("orderBy", parameterToString(*r.orderBy, ""))
	}
	if r.search != nil {
		localVarQueryParams.Add("search", parameterToString(*r.search, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiListConnectorClustersRequest struct {
	ctx        context.Context
	ApiService *ConnectorClustersApiService
	page       *string
	size       *string
}

// Page index
func (r ApiListConnectorClustersRequest) Page(page string) ApiListConnectorClustersRequest {
	r.page = &page
	return r
}

// Number of items in each page
func (r ApiListConnectorClustersRequest) Size(size string) ApiListConnectorClustersRequest {
	r.size = &size
	return r
}

func (r ApiListConnectorClustersRequest) Execute() (*ConnectorClusterList, *http.Response, error) {
	return r.ApiService.ListConnectorClustersExecute(r)
}

/*
ListConnectorClusters Returns a list of connector clusters

Returns a list of connector clusters

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiListConnectorClustersRequest
*/
func (a *ConnectorClustersApiService) ListConnectorClusters(ctx context.Context) ApiListConnectorClustersRequest {
	return ApiListConnectorClustersRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//  @return ConnectorClusterList
func (a *ConnectorClustersApiService) ListConnectorClustersExecute(r ApiListConnectorClustersRequest) (*ConnectorClusterList, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *ConnectorClusterList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ConnectorClustersApiService.ListConnectorClusters")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/connector_mgmt/v1/kafka_connector_clusters"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.page != nil {
		localVarQueryParams.Add("page", parameterToString(*r.page, ""))
	}
	if r.size != nil {
		localVarQueryParams.Add("size", parameterToString(*r.size, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiUpdateConnectorClusterByIdRequest struct {
	ctx                     context.Context
	ApiService              *ConnectorClustersApiService
	connectorClusterId      string
	connectorClusterRequest *ConnectorClusterRequest
}

// Data to updated connector with
func (r ApiUpdateConnectorClusterByIdRequest) ConnectorClusterRequest(connectorClusterRequest ConnectorClusterRequest) ApiUpdateConnectorClusterByIdRequest {
	r.connectorClusterRequest = &connectorClusterRequest
	return r
}

func (r ApiUpdateConnectorClusterByIdRequest) Execute() (*http.Response, error) {
	return r.ApiService.UpdateConnectorClusterByIdExecute(r)
}

/*
UpdateConnectorClusterById udpate a connector cluster

udpate a connector cluster

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param connectorClusterId The id of the connector cluster
 @return ApiUpdateConnectorClusterByIdRequest
*/
func (a *ConnectorClustersApiService) UpdateConnectorClusterById(ctx context.Context, connectorClusterId string) ApiUpdateConnectorClusterByIdRequest {
	return ApiUpdateConnectorClusterByIdRequest{
		ApiService:         a,
		ctx:                ctx,
		connectorClusterId: connectorClusterId,
	}
}

// Execute executes the request
func (a *ConnectorClustersApiService) UpdateConnectorClusterByIdExecute(r ApiUpdateConnectorClusterByIdRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPut
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ConnectorClustersApiService.UpdateConnectorClusterById")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/connector_mgmt/v1/kafka_connector_clusters/{connector_cluster_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"connector_cluster_id"+"}", url.PathEscape(parameterToString(r.connectorClusterId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.connectorClusterRequest == nil {
		return nil, reportError("connectorClusterRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.connectorClusterRequest
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}
