package osv2

import (
	"bytes"
	"context"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

// CheckIndexExistsRequest is a domain model union type for all the fields of ExistIndexRequests for all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// An empty CheckIndexExistsRequest will fail to execute. At least one index is required to check if it exists or not
//
//	[ExistIndex] https://opensearch.org/docs/latest/api-reference/index-apis/exists/
type CheckIndexExistsRequest struct {
	Indices           []string
	ExpandWildcards   string
	IgnoreUnavailable bool
	AllowNoIndices    bool
	OnlyLocalNode     bool
}

// FromDomainCheckIndexExitsRequest creates a new [CheckIndexExitsRequest] from the given [opensearchtools.CheckIndexExistsRequest]
func FromDomainCheckIndexExitsRequest(req *opensearchtools.CheckIndexExistsRequest) (CheckIndexExistsRequest, opensearchtools.ValidationResults) {
	// As more versions are implemented, these [opensearchtools.ValidationResults] may be used to contain issues
	// converting from the domain model to the V2 model.
	var vrs opensearchtools.ValidationResults

	return CheckIndexExistsRequest{
		Indices:           req.Indices,
		ExpandWildcards:   req.ExpandWildcards,
		IgnoreUnavailable: req.IgnoreUnavailable,
		AllowNoIndices:    req.AllowNoIndices,
		OnlyLocalNode:     req.OnlyLocalNode,
	}, vrs
}

// Validate validates the given CheckIndexExistsRequest
func (e *CheckIndexExistsRequest) Validate() opensearchtools.ValidationResults {
	var validationResults opensearchtools.ValidationResults

	if len(e.Indices) == 0 {
		validationResults.Add(opensearchtools.NewValidationResult("Index not set at the CheckIndexExistsRequest", true))
	}

	return validationResults
}

// NewExistIndexRequest creates a CheckIndexExistsRequest with defaults
func NewExistIndexRequest() *CheckIndexExistsRequest {
	return &CheckIndexExistsRequest{
		ExpandWildcards: "open",
		AllowNoIndices:  true,
	}
}

// WithIndices sets indices to be checked for CheckIndexExistsRequest
func (e *CheckIndexExistsRequest) WithIndices(indices []string) *CheckIndexExistsRequest {
	e.Indices = indices
	return e
}

// WithOnlyLocalNode local for CheckIndexExistsRequest, it defines Whether to return information
// from only the local node instead of from the master node. Default is false.
func (e *CheckIndexExistsRequest) WithOnlyLocalNode(l bool) *CheckIndexExistsRequest {
	e.OnlyLocalNode = l
	return e
}

// WithExpandWildCard sets expand_wildcards option for CheckIndexExistsRequest,
// it expands wildcard expressions to different indices, default is open
func (e *CheckIndexExistsRequest) WithExpandWildCard(w string) *CheckIndexExistsRequest {
	e.ExpandWildcards = w
	return e
}

// WithIgnoreUnavailable sets ignore_unavailable options for CheckIndexExistsRequest,
// If true, OpenSearch does not include missing or closed indices in the response. Default is false
func (e *CheckIndexExistsRequest) WithIgnoreUnavailable(i bool) *CheckIndexExistsRequest {
	e.IgnoreUnavailable = i
	return e
}

// WithAllowNoIndices sets allow_no_indices for CheckIndexExistsRequest,
// it defines Whether to ignore wildcards that don’t match any indices. Default is true
func (e *CheckIndexExistsRequest) WithAllowNoIndices(a bool) *CheckIndexExistsRequest {
	e.AllowNoIndices = a
	return e
}

// Do executes the [CheckIndexExistsRequest] using the provided opensearch.Client.
// If the request is executed successfully, then a [CheckIndexExistsResponse] will be returned.
// An error can be returned if
//
//   - Index is missing
//   - The call to OpenSearch fails
//   - The result json cannot be unmarshalled
func (e *CheckIndexExistsRequest) Do(ctx context.Context, client *opensearch.Client) (*opensearchtools.OpenSearchResponse[CheckIndexExistsResponse], error) {
	vrs := e.Validate()
	if vrs.IsFatal() {
		return nil, opensearchtools.NewValidationError(vrs)
	}

	osResp, rErr := opensearchapi.IndicesExistsRequest{
		Index:             e.Indices,
		AllowNoIndices:    &e.AllowNoIndices,
		ExpandWildcards:   e.ExpandWildcards,
		IgnoreUnavailable: &e.IgnoreUnavailable,
		Local:             &e.OnlyLocalNode,
	}.Do(ctx, client)

	if rErr != nil {
		return nil, rErr
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	resp := CheckIndexExistsResponse{}

	if osResp.StatusCode == http.StatusOK {
		resp.Exists = true
	}

	return &opensearchtools.OpenSearchResponse[CheckIndexExistsResponse]{
		StatusCode:        osResp.StatusCode,
		Header:            osResp.Header,
		Response:          resp,
		ValidationResults: vrs,
	}, nil
}

// CheckIndexExistsResponse defines the response if the index exists or not
type CheckIndexExistsResponse struct {
	Exists bool
}

// toDomain converts this instance of [CheckIndexExistsResponse] into an [opensearchtools.CheckIndexExistsResponse]
func (e CheckIndexExistsResponse) toDomain() opensearchtools.CheckIndexExistsResponse {
	domainResp := opensearchtools.CheckIndexExistsResponse{
		Exists: e.Exists,
	}
	return domainResp
}
