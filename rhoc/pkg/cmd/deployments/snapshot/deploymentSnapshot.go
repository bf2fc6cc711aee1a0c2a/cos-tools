package snapshot

import "github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/admin"

type DeploymentSnapshot struct {
	ClusterId string
	Id        string
	Kind      string
	Status    admin.ConnectorState
}
