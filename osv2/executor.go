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
	osv2Req := fromDomainMGetRequest(req)
	validationRes := osv2Req.Validate()
	if validationRes.IsFatal() {
		return nil, opensearchtools.NewValidationError(validationRes)
	}

	osv2Resp, err := osv2Req.Do(ctx, e.Client)
	if err != nil {
		return nil, err
	}

	return osv2Resp.ToDomain(validationRes), err
}
