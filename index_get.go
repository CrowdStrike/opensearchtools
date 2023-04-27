package opensearchtools

import (
	"encoding/json"
	"time"
)

// GetIndexRequest is a domain model union type for all the fields of GetIndexRequests for all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// An empty GetIndexRequest will fail to execute. At least one index is required to get the index information
//
//	[GetIndex] https://opensearch.org/docs/latest/api-reference/index-apis/get-index/
type GetIndexRequest struct {
	Indices           []string
	MasterTimeout     time.Duration
	ExpandWildcards   string //todo: use an enum
	IgnoreUnavailable bool
	OnlyLocalNode     bool
	IncludeDefaults   bool
}

// NewGetIndexRequest instantiates a NewGetIndexRequest with default values
func NewGetIndexRequest() *GetIndexRequest {
	return &GetIndexRequest{
		MasterTimeout:   30 * time.Second,
		ExpandWildcards: "open",
	}
}

// WithIndices sets indices to be retried for GetIndexRequest
func (g *GetIndexRequest) WithIndices(indices []string) *GetIndexRequest {
	g.Indices = indices
	return g
}

// WithMasterTimeout sets the master_timeout for GetIndexRequest
// it defines how long to wait for a connection to the master node. Default is 30s.
func (g *GetIndexRequest) WithMasterTimeout(duration time.Duration) *GetIndexRequest {
	g.MasterTimeout = duration
	return g
}

// WithExpandWildCard sets expand_wildcards option for GetIndexRequest,
// it expands wildcard expressions to different indices, default is open
func (g *GetIndexRequest) WithExpandWildCard(w string) *GetIndexRequest {
	g.ExpandWildcards = w
	return g
}

// WithIgnoreUnavailable sets ignore_unavailable options for GetIndexRequest,
// If true, OpenSearch does not include missing or closed indices in the response. Default is false
func (g *GetIndexRequest) WithIgnoreUnavailable(i bool) *GetIndexRequest {
	g.IgnoreUnavailable = i
	return g
}

// WithOnlyLocalNode sets local for GetIndexRequest, it defines Whether to return information
// from only the local node instead of from the master node. Default is false.
func (g *GetIndexRequest) WithOnlyLocalNode(l bool) *GetIndexRequest {
	g.OnlyLocalNode = l
	return g
}

// WithIncludeDefaults sets include_defaults for GetIndexRequest,
// it defines Whether to include default settings as part of the response. Default is false
func (g *GetIndexRequest) WithIncludeDefaults(d bool) *GetIndexRequest {
	g.IncludeDefaults = d
	return g
}

// GetIndexResponse showcase the response of GetIndexRequest
type GetIndexResponse struct {
	Response map[string]IndexInfo
}

// IndexInfo contains the aliases, mappings and settings index info
type IndexInfo struct {
	Aliases  map[string]json.RawMessage
	Mappings map[string]json.RawMessage
	Settings struct{ Index IndexSetting }
}

// IndexSetting contains the detailed index settings info
type IndexSetting struct {
	RefreshInterval  string
	CreationDate     string
	NumberOfShards   string
	NumberOfReplicas string
	UUID             string
	Version          struct{ Created string }
	ProvidedName     string
}
