package opensearchtools

import "time"

// ExistIndexRequest is a domain model union type for all the fields of ExistIndexRequests for all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// An empty ExistIndexRequest will fail to execute. At least one index is required to check if it exists or not
//  [ExistIndex] https://opensearch.org/docs/latest/api-reference/index-apis/exists/
type ExistIndexRequest struct {
	Indices           []string
	MasterTimeout     time.Duration
	Timeout           time.Duration
	ExpandWildcards   string
	IgnoreUnavailable bool
	AllowNoIndices    bool
	OnlyLocalNode     bool
}

// NewExistIndexRequest creates a ExistIndexRequest with defaults
func NewExistIndexRequest() *ExistIndexRequest {
	return &ExistIndexRequest{
		MasterTimeout:   30 * time.Second,
		Timeout:         30 * time.Second,
		ExpandWildcards: "open",
		AllowNoIndices:  true,
	}
}

// WithIndices adds indices to be deleted for ExistIndexRequest
func (e *ExistIndexRequest) WithIndices(indices []string) *ExistIndexRequest {
	e.Indices = indices
	return e
}

// WithMasterTimeout adds the master_timeout for ExistIndexRequest
// it defines how long to wait for a connection to the master node. Default is 30s.
func (e *ExistIndexRequest) WithMasterTimeout(duration time.Duration) *ExistIndexRequest {
	e.MasterTimeout = duration
	return e
}

// WithOnlyLocalNode local for ExistIndexRequest, it defines Whether to return information
// from only the local node instead of from the master node. Default is false.
func (e *ExistIndexRequest) WithOnlyLocalNode(l bool) *ExistIndexRequest {
	e.OnlyLocalNode = l
	return e
}

// WithTimeout adds the timeout for ExistIndexRequest, it defines how long to wait for the request to return. Default is 30s
func (e *ExistIndexRequest) WithTimeout(duration time.Duration) *ExistIndexRequest {
	e.Timeout = duration
	return e
}

// WithExpandWildCard adds expand_wildcards option for ExistIndexRequest,
// it expands wildcard expressions to different indices, default is open
func (e *ExistIndexRequest) WithExpandWildCard(w string) *ExistIndexRequest {
	e.ExpandWildcards = w
	return e
}

// WithIgnoreUnavailable add ignore_unavailable options for ExistIndexRequest,
// If true, OpenSearch does not include missing or closed indices in the response. Default is false
func (e *ExistIndexRequest) WithIgnoreUnavailable(i bool) *ExistIndexRequest {
	e.IgnoreUnavailable = i
	return e
}

// WithAllowNoIndices defines allow_no_indices for ExistIndexRequest,
// it defines Whether to ignore wildcards that don’t match any indices. Default is true
func (e *ExistIndexRequest) WithAllowNoIndices(a bool) *ExistIndexRequest {
	e.AllowNoIndices = a
	return e
}

// ExistIndexResponse defines the response that contains the status,
// possible response codes: 200 – the index exists, and 404 – the index does not exist
type ExistIndexResponse struct {
	Status int
}
