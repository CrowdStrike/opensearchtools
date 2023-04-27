package opensearchtools

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

// NewExistIndexRequest creates a CheckIndexExistsRequest with defaults
func NewExistIndexRequest() *CheckIndexExistsRequest {
	return &CheckIndexExistsRequest{
		ExpandWildcards: "open",
		AllowNoIndices:  true,
	}
}

// WithIndices sets indices to be deleted for CheckIndexExistsRequest
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

// CheckIndexExistsResponse defines the response if the index exists or not
type CheckIndexExistsResponse struct {
	Exists bool
}
