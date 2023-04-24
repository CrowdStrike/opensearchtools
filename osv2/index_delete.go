package osv2

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

// DeleteIndexRequest is a domain model union type for all the fields of DeleteIndexRequests for all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// An empty DeleteIndexRequest will fail to execute. At least one index is required to be deleted
//
//	[DeleteIndex] https://opensearch.org/docs/latest/api-reference/index-apis/delete-index/
type DeleteIndexRequest struct {
	Indices           []string
	MasterTimeout     time.Duration
	Timeout           time.Duration
	ExpandWildcards   string
	IgnoreUnavailable bool
	AllowNoIndices    bool
}

// FromDomainDeleteIndexRequest creates a new [DeleteIndexRequest] from the given [opensearchtools.DeleteIndexRequest]
func FromDomainDeleteIndexRequest(req *opensearchtools.DeleteIndexRequest) (DeleteIndexRequest, opensearchtools.ValidationResults) {
	// As more versions are implemented, these [opensearchtools.ValidationResults] may be used to contain issues
	// converting from the domain model to the V2 model.
	var vrs opensearchtools.ValidationResults

	return DeleteIndexRequest{
		Indices:           req.Indices,
		MasterTimeout:     req.MasterTimeout,
		Timeout:           req.Timeout,
		ExpandWildcards:   req.ExpandWildcards,
		IgnoreUnavailable: req.IgnoreUnavailable,
		AllowNoIndices:    req.AllowNoIndices,
	}, vrs
}

// Validate validates the given DeleteIndexRequest
func (d *DeleteIndexRequest) Validate() opensearchtools.ValidationResults {
	var validationResults opensearchtools.ValidationResults

	if len(d.Indices) == 0 {
		validationResults.Add(opensearchtools.NewValidationResult("Index not set at the DeleteIndexRequest", true))
	}

	return validationResults
}

// NewDeleteIndexRequest instantiates a DeleteIndexRequest with default values
func NewDeleteIndexRequest() *DeleteIndexRequest {
	return &DeleteIndexRequest{
		MasterTimeout:   30 * time.Second,
		Timeout:         30 * time.Second,
		ExpandWildcards: "open",
		AllowNoIndices:  true,
	}
}

// WithIndices sets indices to be deleted for DeleteIndexRequest
func (d *DeleteIndexRequest) WithIndices(indices []string) *DeleteIndexRequest {
	d.Indices = indices
	return d
}

// WithMasterTimeout sets the master_timeout for DeleteIndexRequest
// it defines how long to wait for a connection to the master node. Default is 30s.
func (d *DeleteIndexRequest) WithMasterTimeout(duration time.Duration) *DeleteIndexRequest {
	d.MasterTimeout = duration
	return d
}

// WithTimeout sets the timeout for DeleteIndexRequest, it defines how long to wait for the request to return. Default is 30s
func (d *DeleteIndexRequest) WithTimeout(duration time.Duration) *DeleteIndexRequest {
	d.Timeout = duration
	return d
}

// WithExpandWildCard sets expand_wildcards option for DeleteIndexRequest,
// it expands wildcard expressions to different indices, default is open
func (d *DeleteIndexRequest) WithExpandWildCard(w string) *DeleteIndexRequest {
	d.ExpandWildcards = w
	return d
}

// WithIgnoreUnavailable sets ignore_unavailable options for DeleteIndexRequest,
// If true, OpenSearch does not include missing or closed indices in the response. Default is false
func (d *DeleteIndexRequest) WithIgnoreUnavailable(i bool) *DeleteIndexRequest {
	d.IgnoreUnavailable = i
	return d
}

// WithAllowNoIndices sets allow_no_indices for DeleteIndexRequest,
// it defines Whether to ignore wildcards that don’t match any indices. Default is true
func (d *DeleteIndexRequest) WithAllowNoIndices(a bool) *DeleteIndexRequest {
	d.AllowNoIndices = a
	return d
}

// Do executes the [DeleteIndexRequest] using the provided opensearch.Client.
// If the request is executed successfully, then a [DeleteIndexRequest] will be returned.
// An error can be returned if
//
//   - Index is missing
//   - The call to OpenSearch fails
//   - The result json cannot be unmarshalled
func (d *DeleteIndexRequest) Do(ctx context.Context, client *opensearch.Client) (*opensearchtools.OpenSearchResponse[DeleteIndexResponse], error) {
	vrs := d.Validate()
	if vrs.IsFatal() {
		return nil, opensearchtools.NewValidationError(vrs)
	}

	osResp, rErr := opensearchapi.IndicesDeleteRequest{
		Index:             d.Indices,
		AllowNoIndices:    &d.AllowNoIndices,
		ExpandWildcards:   d.ExpandWildcards,
		IgnoreUnavailable: &d.IgnoreUnavailable,
		MasterTimeout:     d.MasterTimeout,
		Timeout:           d.Timeout,
	}.Do(ctx, client)

	if rErr != nil {
		return nil, rErr
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	resp := DeleteIndexResponse{}

	if err := json.Unmarshal(respBuf.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &opensearchtools.OpenSearchResponse[DeleteIndexResponse]{
		StatusCode:        osResp.StatusCode,
		Header:            osResp.Header,
		Response:          resp,
		ValidationResults: vrs,
	}, nil
}

// DeleteIndexResponse represent the response for DeleteIndexResponse, either error or acknowledged
type DeleteIndexResponse struct {
	Acknowledged bool
	Error        *Error
}

// toDomain converts this instance of [DeleteIndexResponse] into an [opensearchtools.DeleteIndexResponse]
func (d DeleteIndexResponse) toDomain() opensearchtools.DeleteIndexResponse {
	domainResp := opensearchtools.DeleteIndexResponse{
		Acknowledged: d.Acknowledged,
	}

	if d.Error != nil {
		domainErr := d.Error.toDomain()
		domainResp.Error = &domainErr
	}

	return domainResp
}
