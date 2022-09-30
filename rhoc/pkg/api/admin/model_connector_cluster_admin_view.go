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

// ConnectorClusterAdminView struct for ConnectorClusterAdminView
type ConnectorClusterAdminView struct {
	Id         string                      `json:"id,omitempty"`
	Kind       string                      `json:"kind,omitempty"`
	Href       string                      `json:"href,omitempty"`
	Owner      string                      `json:"owner,omitempty"`
	CreatedAt  time.Time                   `json:"created_at,omitempty"`
	ModifiedAt time.Time                   `json:"modified_at,omitempty"`
	Name       string                      `json:"name,omitempty"`
	Status     ConnectorClusterAdminStatus `json:"status,omitempty"`
}
