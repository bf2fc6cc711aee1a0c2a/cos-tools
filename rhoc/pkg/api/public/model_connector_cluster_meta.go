/*
Connector Management API

Connector Management API is a REST API to manage connectors.

API version: 0.1.0
Contact: rhosak-support@redhat.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public

import (
	"time"
)

// ConnectorClusterMeta struct for ConnectorClusterMeta
type ConnectorClusterMeta struct {
	Owner      string    `json:"owner,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty"`
	Name       string    `json:"name,omitempty"`
	// Name-value string annotations for resource
	Annotations map[string]string `json:"annotations,omitempty"`
}
