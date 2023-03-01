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
//   - Fatal validation issues are found
//   - The request to OpenSearch fails
//   - The results JSON cannot be unmarshalled
func (e *Executor) MGet(ctx context.Context, req *opensearchtools.MGetRequest) (resp opensearchtools.OpenSearchResponse[opensearchtools.MGetResponse], err error) {
	validationResults := opensearchtools.NewValidationResults()

	osv2Req, vrs := fromDomainMGetRequest(req)
	if vrs.IsFatal() {
		resp.ValidationResults.Extend(vrs)
		return resp, opensearchtools.NewValidationError(vrs)
	}
	validationResults.Extend(vrs)

	osv2Resp, err := osv2Req.Do(ctx, e.Client)
	if err != nil {
		return resp, err
	}
	validationResults.Extend(osv2Resp.ValidationResults)

	return opensearchtools.NewOpenSearchResponse(
		validationResults,
		osv2Resp.StatusCode,
		osv2Resp.Header,
		osv2Resp.Response.toDomain(),
	), nil
}

// Search executes the SearchRequest using the provided [opensearchtools.SearchRequest].
// If the request is executed successfully, then an [opensearchtools.SearchResponse] will be returned.
// An error can be returned if:
//   - Fatal validation issues are found
//   - The request to OpenSearch fails
//   - The results JSON cannot be unmarshalled
func (e *Executor) Search(ctx context.Context, req *opensearchtools.SearchRequest) (resp opensearchtools.OpenSearchResponse[opensearchtools.SearchResponse], err error) {
	var validationResults opensearchtools.ValidationResults

	osv2Req, vrs := fromDomainSearchRequest(req)
	if vrs.IsFatal() {
		resp.ValidationResults.Extend(vrs)
		return resp, opensearchtools.NewValidationError(vrs)
	}
	validationResults.Extend(vrs)

	osv2Resp, err := osv2Req.Do(ctx, e.Client)
	if err != nil {
		return resp, err
	}
	validationResults.Extend(osv2Resp.ValidationResults)

	return opensearchtools.NewOpenSearchResponse(
		validationResults,
		osv2Resp.StatusCode,
		osv2Resp.Header,
		osv2Resp.Response.ToDomain(),
	), nil
}
