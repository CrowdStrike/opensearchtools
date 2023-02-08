package osv2

import (
	"github.com/opensearch-project/opensearch-go/v2"
)

// OSv2Executor is an executor for OpenSearch 2.
type OSv2Executor struct {
	// OpenSearch 2 specifc client
	client *opensearch.Client
}

// NewOSv2Executor creates a new OSv2Executor instance.
func NewOSv2Executor(client *opensearch.Client) *OSv2Executor {
	return &OSv2Executor{
		client,
	}
}
