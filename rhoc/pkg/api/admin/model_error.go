/*
Connector Service Fleet Manager Admin APIs

Connector Service Fleet Manager Admin is a Rest API to manage connector clusters.

API version: 0.0.3
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package admin

// Error struct for Error
type Error struct {
	Reason      string `json:"reason"`
	OperationId string `json:"operation_id,omitempty"`
	Id          string `json:"id"`
	Kind        string `json:"kind"`
	Href        string `json:"href"`
	Code        string `json:"code"`
}
