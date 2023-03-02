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
	osv2Req, vrs := fromDomainMGetRequest(req)
	resp.ValidationResults.Extend(vrs)
	if vrs.IsFatal() {
		return resp, opensearchtools.NewValidationError(vrs)
	}

	osv2Resp, reqErr := osv2Req.Do(ctx, e.Client)
	if reqErr != nil {
		return resp, reqErr
	}

	resp.ValidationResults.Extend(osv2Resp.ValidationResults)
	resp.Response = osv2Resp.Response.toDomain()
	resp.StatusCode = osv2Resp.StatusCode
	resp.Header = osv2Resp.Header

	return resp, nil
}

// Search executes the SearchRequest using the provided [opensearchtools.SearchRequest].
// If the request is executed successfully, then an [opensearchtools.SearchResponse] will be returned.
// An error can be returned if:
//   - Fatal validation issues are found
//   - The request to OpenSearch fails
//   - The results JSON cannot be unmarshalled
func (e *Executor) Search(ctx context.Context, req *opensearchtools.SearchRequest) (resp opensearchtools.OpenSearchResponse[opensearchtools.SearchResponse], err error) {
	osv2Req, vrs := fromDomainSearchRequest(req)
	resp.ValidationResults.Extend(vrs)
	if vrs.IsFatal() {
		return resp, opensearchtools.NewValidationError(vrs)
	}

	osv2Resp, reqErr := osv2Req.Do(ctx, e.Client)
	if reqErr != nil {
		return resp, reqErr
	}

	resp.ValidationResults.Extend(osv2Resp.ValidationResults)
	resp.Response = osv2Resp.Response.ToDomain()
	resp.StatusCode = osv2Resp.StatusCode
	resp.Header = osv2Resp.Header

	return resp, nil
}

// Bulk executes the BulkRequest using the provided [opensearchtools.BulkRequest].
// If the request is executed successfully, then an
// [opensearchtools.OpenSearchResponse] containing a [opensearchtools.BulkResponse]
// An error can be returned if:
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
func (e *Executor) Bulk(ctx context.Context, req *opensearchtools.BulkRequest) (resp opensearchtools.OpenSearchResponse[opensearchtools.BulkResponse], err error) {
	osv2Req, vrs := fromDomainBulkRequest(req)
	resp.ValidationResults.Extend(vrs)

	if vrs.IsFatal() {
		return resp, opensearchtools.NewValidationError(vrs)
	}

	osv2Resp, reqErr := osv2Req.Do(ctx, e.Client)
	if reqErr != nil {
		return resp, reqErr
	}

	resp.ValidationResults.Extend(osv2Resp.ValidationResults)
	resp.Response = osv2Resp.Response.toDomain()
	resp.StatusCode = osv2Resp.StatusCode
	resp.Header = osv2Resp.Header

	return resp, nil
}
