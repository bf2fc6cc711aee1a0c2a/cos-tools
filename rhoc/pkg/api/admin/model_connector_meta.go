/*
Connector Service Fleet Manager Admin APIs

Connector Service Fleet Manager Admin is a Rest API to manage connector clusters.

API version: 0.0.3
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package admin

import (
	"time"
)

// ConnectorMeta struct for ConnectorMeta
type ConnectorMeta struct {
	Owner           string                `json:"owner,omitempty"`
	CreatedAt       time.Time             `json:"created_at,omitempty"`
	ModifiedAt      time.Time             `json:"modified_at,omitempty"`
	Name            string                `json:"name"`
	ConnectorTypeId string                `json:"connector_type_id"`
	NamespaceId     string                `json:"namespace_id"`
	Channel         Channel               `json:"channel,omitempty"`
	DesiredState    ConnectorDesiredState `json:"desired_state"`
	// Name-value string annotations for resource
	Annotations     map[string]string `json:"annotations,omitempty"`
	ResourceVersion int64             `json:"resource_version,omitempty"`
}
