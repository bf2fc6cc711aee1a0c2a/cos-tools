/*
Connector Service Fleet Manager Admin APIs

Connector Service Fleet Manager Admin is a Rest API to manage connector clusters.

API version: 0.0.3
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package admin

// ConnectorNamespaceEvalRequest An evaluation connector namespace create request
type ConnectorNamespaceEvalRequest struct {
	// Namespace name must match pattern `^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$`, or it may be empty to be auto-generated.
	Name string `json:"name,omitempty"`
	// Name-value string annotations for resource
	Annotations map[string]string `json:"annotations,omitempty"`
}
