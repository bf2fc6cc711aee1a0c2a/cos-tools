/*
Connector Service Fleet Manager Admin APIs

Connector Service Fleet Manager Admin is a Rest API to manage connector clusters.

API version: 0.0.3
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package admin

// ConnectorDeploymentSpec Holds the deployment specification of a connector
type ConnectorDeploymentSpec struct {
	ConnectorId              string `json:"connector_id,omitempty"`
	ConnectorResourceVersion int64  `json:"connector_resource_version,omitempty"`
	ConnectorTypeId          string `json:"connector_type_id,omitempty"`
	ClusterId                string `json:"cluster_id,omitempty"`
	NamespaceId              string `json:"namespace_id,omitempty"`
	// allow the connector to upgrade to a new operator
	// Deprecated
	DeprecatedAllowUpgrade bool `json:"allow_upgrade,omitempty"`
	// an optional operator id that the connector should be run under.
	OperatorId    string                 `json:"operator_id,omitempty"`
	DesiredState  ConnectorDesiredState  `json:"desired_state,omitempty"`
	ShardMetadata map[string]interface{} `json:"shard_metadata,omitempty"`
}
