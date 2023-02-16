package osv2

import (
	"context"

	"github.com/opensearch-project/opensearch-go/v2"

	"github.com/CrowdStrike/opensearchtools"
)

// OSv2Executor is an executor for OpenSearch 2.
type OSv2Executor struct {
	// OpenSearch 2 specifc client
	Client *opensearch.Client
}

// NewOSv2Executor creates a new [OSv2Executor] instance.
func NewOSv2Executor(client *opensearch.Client) *OSv2Executor {
	return &OSv2Executor{
		Client: client,
	}
}

// MGet executes the Multi-Get MGetRequest using the provided [opensearchtools.MGetRequest].
// If the request is executed successfully, then an [opensearchtools.MGetResponse] with [opensearchtools.MGetResults]
// will be returned.
// An error can be returned if:
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
func (e *OSv2Executor) MGet(ctx context.Context, req *opensearchtools.MGetRequest) (*opensearchtools.MGetResponse, error) {
	specMGetRequest := FromModelMGetRequest(req)

	specResponse, err := specMGetRequest.Do(ctx, e.Client)
	if err != nil {
		return nil, err
	}

	modelMGetResponse := specResponse.ToModel()
	return modelMGetResponse, err
}
