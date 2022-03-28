/*
 * Connector Service Fleet Manager Admin APIs
 *
 * Connector Service Fleet Manager Admin is a Rest API to manage connector clusters.
 *
 * API version: 0.0.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package admin

// ConnectorAvailableOperatorUpgrade An available operator upgrade for a connector
type ConnectorAvailableOperatorUpgrade struct {
	ConnectorId     string                                    `json:"connector_id,omitempty"`
	NamespaceId     string                                    `json:"namespace_id,omitempty"`
	ConnectorTypeId string                                    `json:"connector_type_id,omitempty"`
	Channel         string                                    `json:"channel,omitempty"`
	Operator        ConnectorAvailableOperatorUpgradeOperator `json:"operator,omitempty"`
}
