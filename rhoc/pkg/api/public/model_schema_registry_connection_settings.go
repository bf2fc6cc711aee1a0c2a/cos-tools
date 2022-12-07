/*
Connector Management API

Connector Management API is a REST API to manage connectors.

API version: 0.1.0
Contact: rhosak-support@redhat.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public

// SchemaRegistryConnectionSettings Holds the configuration to connect to a Schem Registry Instance.
type SchemaRegistryConnectionSettings struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
