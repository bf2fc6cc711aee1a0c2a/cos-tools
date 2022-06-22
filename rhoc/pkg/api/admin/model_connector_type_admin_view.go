/*
Connector Service Fleet Manager Admin APIs

Connector Service Fleet Manager Admin is a Rest API to manage connector clusters.

API version: 0.0.3
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package admin

// ConnectorTypeAdminView Holds the connector type
type ConnectorTypeAdminView struct {
	Id   string `json:"id,omitempty"`
	Kind string `json:"kind,omitempty"`
	Href string `json:"href,omitempty"`
	// Name of the connector type.
	Name string `json:"name"`
	// Version of the connector type.
	Version  string                          `json:"version"`
	Channels map[string]ConnectorTypeChannel `json:"channels,omitempty"`
	// A description of the connector.
	Description string `json:"description,omitempty"`
	// URL to an icon of the connector.
	IconHref string `json:"icon_href,omitempty"`
	// Labels used to categorize the connector
	Labels []string `json:"labels,omitempty"`
	// The capabilities supported by the conenctor
	Capabilities []string `json:"capabilities,omitempty"`
	// A json schema that can be used to validate a ConnectorRequest connector field.
	Schema map[string]interface{} `json:"schema"`
	// A json schema that can be used to validate a ConnectorRequest connector field.
	JsonSchema map[string]interface{} `json:"json_schema"`
}