package osv2

import (
	"context"

	"github.com/opensearch-project/opensearch-go/v2"

	"github.com/CrowdStrike/opensearchtools"
)

// Executor is an executor for OpenSearch 2.
type Executor struct {
	// OpenSearch 2 specifc client
	Client *opensearch.Client
}

// NewExecutor creates a new [osv2.Executor] instance.
func NewExecutor(client *opensearch.Client) *Executor {
	return &Executor{
		Client: client,
	}
}

// MGet executes the Multi-Get MGetRequest using the provided [opensearchtools.MGetRequest].
// If the request is executed successfully, then an [opensearchtools.MGetResponse] with [opensearchtools.MGetResults]
// will be returned.
// An error can be returned if:
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
func (e *Executor) MGet(ctx context.Context, req *opensearchtools.MGetRequest) (*opensearchtools.OpenSearchResponse[opensearchtools.MGetResponse], error) {
	mgetReqValidation, err := FromDomainMGetRequest(req)
	if err != nil {
		return nil, err
	}

	specResponse, err := mgetReqValidation.ValidatedRequest.Do(ctx, e.Client)
	if err != nil {
		return nil, err
	}

	return specResponse.ToDomain(mgetReqValidation.ValidationResults), err
}
