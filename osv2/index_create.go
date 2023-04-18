package osv2

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

// CreateIndexRequest is a domain model union type for all the fields of CreateIndexRequests for all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// An empty CreateIndexRequest will fail to execute. At least index is required to be added which matches to
// the existing template index pattern. Otherwise, the DocBody has to be provided with the detailed index information
// as provided in the documentation: [CreateIndex] https://opensearch.org/docs/latest/api-reference/index-apis/create-index/
type CreateIndexRequest struct {
	Index               string
	DocBody             io.Reader
	MasterTimeout       time.Duration
	Timeout             time.Duration
	WaitForActiveShards string // todo: update this with enum or not since we have numbers and all
}

// FromDomainCreateIndexRequest creates a new [BulkRequest] from the given [opensearchtools.CreateIndexRequest]
func FromDomainCreateIndexRequest(req *opensearchtools.CreateIndexRequest) (CreateIndexRequest, opensearchtools.ValidationResults) {
	// As more versions are implemented, these [opensearchtools.ValidationResults] may be used to contain issues
	// converting from the domain model to the V2 model.
	var vrs opensearchtools.ValidationResults

	return CreateIndexRequest{
		Index:               req.Index,
		DocBody:             req.DocBody,
		MasterTimeout:       req.MasterTimeout,
		Timeout:             req.Timeout,
		WaitForActiveShards: req.WaitForActiveShards,
	}, vrs
}

// Validate validates the given CreateIndexRequest
func (c *CreateIndexRequest) Validate() opensearchtools.ValidationResults {
	var validationResults opensearchtools.ValidationResults

	if c.Index == "" {
		validationResults.Add(opensearchtools.NewValidationResult("Index not set at the CreateIndexRequest", true))
	}

	return validationResults
}

// NewCreateIndexRequest instantiates an CreateIndexRequest with default values
func NewCreateIndexRequest() *CreateIndexRequest {
	return &CreateIndexRequest{
		MasterTimeout:       30 * time.Second,
		Timeout:             30 * time.Second,
		WaitForActiveShards: "1",
	}
}

// WithIndex adds the index for CreateIndexRequest
func (c *CreateIndexRequest) WithIndex(index string) *CreateIndexRequest {
	c.Index = index
	return c
}

// WithDocBody adds the body for CreateIndexRequest that contains detailed index setting
func (c *CreateIndexRequest) WithDocBody(body io.Reader) *CreateIndexRequest {
	c.DocBody = body
	return c
}

// WithMasterTimeout adds the master timeout for CreateIndexRequest
// it defines how long to wait for a connection to the master node. Default is 30s.
func (c *CreateIndexRequest) WithMasterTimeout(d time.Duration) *CreateIndexRequest {
	c.MasterTimeout = d
	return c
}

// WithTimeout adds the timeout for CreateIndexRequest, it defines how long to wait for the request to return. Default is 30s
func (c *CreateIndexRequest) WithTimeout(d time.Duration) *CreateIndexRequest {
	c.Timeout = d
	return c
}

// WithWaitForActiveShards to add the active shard options with CreateIndexRequest,
// it specifies the number of active shards that must be available before OpenSearch processes the request. Default is 1
func (c *CreateIndexRequest) WithWaitForActiveShards(s string) *CreateIndexRequest {
	c.WaitForActiveShards = s
	return c
}

// Do executes the [CreateIndexRequest] using the provided opensearch.Client.
// If the request is executed successfully, then a [CreateIndexResponse] will be returned.
// An error can be returned if
//
//   - Index is missing
//   - The call to OpenSearch fails
//   - The result json cannot be unmarshalled
func (c *CreateIndexRequest) Do(ctx context.Context, client *opensearch.Client) (*opensearchtools.OpenSearchResponse[CreateIndexResponse], error) {
	vrs := c.Validate()
	if vrs.IsFatal() {
		return nil, opensearchtools.NewValidationError(vrs)
	}

	osResp, rErr := opensearchapi.IndicesCreateRequest{
		Body:                c.DocBody, // todo: are we sure about this? what will happen to nil?
		MasterTimeout:       c.MasterTimeout,
		Timeout:             c.Timeout,
		WaitForActiveShards: c.WaitForActiveShards,
		Index:               c.Index,
	}.Do(ctx, client)

	if rErr != nil {
		return nil, rErr
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	resp := CreateIndexResponse{}

	if err := json.Unmarshal(respBuf.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &opensearchtools.OpenSearchResponse[CreateIndexResponse]{
		StatusCode:        osResp.StatusCode,
		Header:            osResp.Header,
		Response:          resp,
		ValidationResults: vrs,
	}, nil
}

// CreateIndexResponse represent the response for CreateIndexRequest, either error or acknowledged
type CreateIndexResponse struct {
	Acknowledged *bool
	Error        *Error
}

// toDomain converts this instance of [CreateIndexResponse] into an [opensearchtools.CreateIndexResponse]
func (c CreateIndexResponse) toDomain() opensearchtools.CreateIndexResponse {
	domainResp := opensearchtools.CreateIndexResponse{
		Acknowledged: c.Acknowledged,
	}

	if c.Error != nil {
		domainErr := c.Error.toDomain()
		domainResp.Error = &domainErr
	}

	return domainResp
}
