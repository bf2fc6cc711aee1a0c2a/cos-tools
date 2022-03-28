/*
 * Connector Service Fleet Manager
 *
 * Connector Service Fleet Manager is a Rest API to manage connectors.
 *
 * API version: 0.1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public
// ConnectorNamespaceAllOf struct for ConnectorNamespaceAllOf
type ConnectorNamespaceAllOf struct {
	Name string `json:"name"`
	ClusterId string `json:"cluster_id"`
	Expiration string `json:"expiration,omitempty"`
	Tenant ConnectorNamespaceTenant `json:"tenant"`
	Status ConnectorNamespaceStatus `json:"status"`
}
