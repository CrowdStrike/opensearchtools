package opensearchtools

import (
	"time"
)

// DeleteIndexRequest is a domain model union type for all the fields of DeleteIndexRequests for all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// An empty DeleteIndexRequest will fail to execute. At least one index is required to be deleted
//  [DeleteIndex] https://opensearch.org/docs/latest/api-reference/index-apis/delete-index/
type DeleteIndexRequest struct {
	Indices           []string
	MasterTimeout     time.Duration
	Timeout           time.Duration
	ExpandWildcards   string
	IgnoreUnavailable bool
	AllowNoIndices    bool
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

// WithIndices adds indices to be deleted for DeleteIndexRequest
func (d *DeleteIndexRequest) WithIndices(indices []string) *DeleteIndexRequest {
	d.Indices = indices
	return d
}

// WithMasterTimeout adds the master_timeout for DeleteIndexRequest
// it defines how long to wait for a connection to the master node. Default is 30s.
func (d *DeleteIndexRequest) WithMasterTimeout(duration time.Duration) *DeleteIndexRequest {
	d.MasterTimeout = duration
	return d
}

// WithTimeout adds the timeout for DeleteIndexRequest, it defines how long to wait for the request to return. Default is 30s
func (d *DeleteIndexRequest) WithTimeout(duration time.Duration) *DeleteIndexRequest {
	d.Timeout = duration
	return d
}

// WithExpandWildCard adds expand_wildcards option for DeleteIndexRequest,
// it expands wildcard expressions to different indices, default is open
func (d *DeleteIndexRequest) WithExpandWildCard(w string) *DeleteIndexRequest {
	d.ExpandWildcards = w
	return d
}

// WithIgnoreUnavailable add ignore_unavailable options for DeleteIndexRequest,
// If true, OpenSearch does not include missing or closed indices in the response. Default is false
func (d *DeleteIndexRequest) WithIgnoreUnavailable(i bool) *DeleteIndexRequest {
	d.IgnoreUnavailable = i
	return d
}

// WithAllowNoIndices defines allow_no_indices for DeleteIndexRequest,
// it defines Whether to ignore wildcards that don’t match any indices. Default is true
func (d *DeleteIndexRequest) WithAllowNoIndices(a bool) *DeleteIndexRequest {
	d.AllowNoIndices = a
	return d
}

// DeleteIndexResponse represent the response for DeleteIndexResponse, either error or acknowledged
type DeleteIndexResponse struct {
	Acknowledged *bool
	Error        *Error
}
