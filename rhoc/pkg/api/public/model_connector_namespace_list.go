/*
Connector Management API

Connector Management API is a REST API to manage connectors.

API version: 0.1.0
Contact: rhosak-support@redhat.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public

// ConnectorNamespaceList struct for ConnectorNamespaceList
type ConnectorNamespaceList struct {
	Kind  string               `json:"kind"`
	Page  int32                `json:"page"`
	Size  int32                `json:"size"`
	Total int32                `json:"total"`
	Items []ConnectorNamespace `json:"items"`
}
