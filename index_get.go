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
//  [GetIndex] https://opensearch.org/docs/latest/api-reference/index-apis/get-index/
type GetIndexRequest struct {
	Indices           []string
	MasterTimeout     time.Duration
	Timeout           time.Duration
	ExpandWildcards   string
	IgnoreUnavailable bool
	OnlyLocalNode     bool
	IncludeDefaults   bool
}

// NewGetIndexRequest instantiates a NewGetIndexRequest with default values
func NewGetIndexRequest() *GetIndexRequest {
	return &GetIndexRequest{
		MasterTimeout:   30 * time.Second,
		Timeout:         30 * time.Second,
		ExpandWildcards: "open",
	}
}

// WithIndices adds indices to be deleted for GetIndexRequest
func (g *GetIndexRequest) WithIndices(indices []string) *GetIndexRequest {
	g.Indices = indices
	return g
}

// WithMasterTimeout adds the master_timeout for GetIndexRequest
// it defines how long to wait for a connection to the master node. Default is 30s.
func (g *GetIndexRequest) WithMasterTimeout(duration time.Duration) *GetIndexRequest {
	g.MasterTimeout = duration
	return g
}

// WithTimeout adds the timeout for GetIndexRequest, it defines how long to wait for the request to return. Default is 30s
func (g *GetIndexRequest) WithTimeout(duration time.Duration) *GetIndexRequest {
	g.Timeout = duration
	return g
}

// WithExpandWildCard adds expand_wildcards option for GetIndexRequest,
// it expands wildcard expressions to different indices, default is open
func (g *GetIndexRequest) WithExpandWildCard(w string) *GetIndexRequest {
	g.ExpandWildcards = w
	return g
}

// WithIgnoreUnavailable add ignore_unavailable options for GetIndexRequest,
// If true, OpenSearch does not include missing or closed indices in the response. Default is false
func (g *GetIndexRequest) WithIgnoreUnavailable(i bool) *GetIndexRequest {
	g.IgnoreUnavailable = i
	return g
}

// WithOnlyLocalNode local for GetIndexRequest, it defines Whether to return information
// from only the local node instead of from the master node. Default is false.
func (g *GetIndexRequest) WithOnlyLocalNode(l bool) *GetIndexRequest {
	g.OnlyLocalNode = l
	return g
}

// WithIncludeDefaults defines include_defaults for GetIndexRequest,
// it defines Whether to include default settings as part of the response. Default is false
func (g *GetIndexRequest) WithIncludeDefaults(d bool) *GetIndexRequest {
	g.IncludeDefaults = d
	return g
}

// IndexGetResponse showcase the response of GetIndexRequest
type IndexGetResponse struct {
	Response map[string]IndexInfo
}

// IndexInfo contains the aliases, mappings and settings index info
type IndexInfo struct {
	Aliases  map[string]json.RawMessage
	Mappings map[string]json.RawMessage
	Settings *IndexSettings
}

// IndexSettings contains the index related settings
type IndexSettings struct {
	Index IndexSettingsInfo
}

// IndexSettingsInfo contains the detailed index settings info
type IndexSettingsInfo struct {
	CreationDate     string
	NumberOfShards   string
	NumberOfReplicas string
	UUID             string
	Version          struct{ Created string }
	ProvidedName     string
}
