/*
 * Connector Service Fleet Manager
 *
 * Connector Service Fleet Manager is a Rest API to manage connectors.
 *
 * API version: 0.1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public
// KafkaConnectionSettings Holds the configuration to connect to a Kafka Instance.
type KafkaConnectionSettings struct {
	Id string `json:"id"`
	Url string `json:"url"`
}
