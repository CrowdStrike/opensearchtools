package opensearchtools

import (
	"io"
	"time"
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
	WaitForActiveShards string
}

// NewCreateIndexRequest instantiates an CreateIndexRequest with default values
func NewCreateIndexRequest() *CreateIndexRequest {
	return &CreateIndexRequest{
		MasterTimeout:       30 * time.Second, // todo: discus it may be we want to have pointer and hence nils for defaults
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

// todo: add validate over here

// CreateIndexResponse represent the response for CreateIndexRequest, either error or acknowledged
type CreateIndexResponse struct {
	Acknowledged *bool
	Error        *Error
}
