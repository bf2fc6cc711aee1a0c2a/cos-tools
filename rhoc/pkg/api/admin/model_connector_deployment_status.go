/*
Connector Service Fleet Manager Admin APIs

Connector Service Fleet Manager Admin is a Rest API to manage connector clusters.

API version: 0.0.3
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package admin

// ConnectorDeploymentStatus The status of connector deployment
type ConnectorDeploymentStatus struct {
	Phase           ConnectorState                     `json:"phase,omitempty"`
	ResourceVersion int64                              `json:"resource_version,omitempty"`
	Operators       ConnectorDeploymentStatusOperators `json:"operators,omitempty"`
	Conditions      []MetaV1Condition                  `json:"conditions,omitempty"`
}
