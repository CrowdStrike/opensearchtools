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

// GetIndexRequest is a version-specific model for OSv2  of GetIndexRequests for OpenSearch 2
// An empty GetIndexRequest will fail to execute. At least one index is required to get the index information
//
//	[GetIndex] https://opensearch.org/docs/latest/api-reference/index-apis/get-index/
type GetIndexRequest struct {
	Indices           []string
	MasterTimeout     time.Duration
	ExpandWildcards   string
	IgnoreUnavailable bool
	OnlyLocalNode     bool
	IncludeDefaults   bool
}

// FromDomainGetIndexRequest creates a new [GetIndexRequest] from the given [opensearchtools.GetIndexRequest]
func FromDomainGetIndexRequest(req *opensearchtools.GetIndexRequest) (GetIndexRequest, opensearchtools.ValidationResults) {
	// As more versions are implemented, these [opensearchtools.ValidationResults] may be used to contain issues
	// converting from the domain model to the V2 model.
	var vrs opensearchtools.ValidationResults

	return GetIndexRequest{
		Indices:           req.Indices,
		MasterTimeout:     req.MasterTimeout,
		ExpandWildcards:   req.ExpandWildcards,
		IgnoreUnavailable: req.IgnoreUnavailable,
		OnlyLocalNode:     req.OnlyLocalNode,
		IncludeDefaults:   req.IncludeDefaults,
	}, vrs
}

// Validate validates the given GetIndexRequest
func (g *GetIndexRequest) Validate() opensearchtools.ValidationResults {
	var validationResults opensearchtools.ValidationResults

	if len(g.Indices) == 0 {
		validationResults.Add(opensearchtools.NewValidationResult("Index not set at the GetIndexRequest", true))
	}

	return validationResults
}

// NewGetIndexRequest instantiates a NewGetIndexRequest with default values
func NewGetIndexRequest() *GetIndexRequest {
	return &GetIndexRequest{
		MasterTimeout:   30 * time.Second,
		ExpandWildcards: "open",
	}
}

// WithIndices set indices to be retried for GetIndexRequest
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

// Do executes the [GetIndexRequest] using the provided opensearch.Client.
// If the request is executed successfully, then a [GetIndexResponse] will be returned.
// An error can be returned if
//
//   - Index is missing
//   - The call to OpenSearch fails
//   - The result json cannot be unmarshalled
func (g *GetIndexRequest) Do(ctx context.Context, client *opensearch.Client) (*opensearchtools.OpenSearchResponse[GetIndexResponse], error) {
	vrs := g.Validate()
	if vrs.IsFatal() {
		return nil, opensearchtools.NewValidationError(vrs)
	}

	osResp, rErr := opensearchapi.IndicesGetRequest{
		Index:             g.Indices,
		MasterTimeout:     g.MasterTimeout,
		ExpandWildcards:   g.ExpandWildcards,
		IgnoreUnavailable: &g.IgnoreUnavailable,
		Local:             &g.OnlyLocalNode,
		IncludeDefaults:   &g.IncludeDefaults,
	}.Do(ctx, client)

	if rErr != nil {
		return nil, rErr
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	resp := GetIndexResponse{}

	if err := json.Unmarshal(respBuf.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &opensearchtools.OpenSearchResponse[GetIndexResponse]{
		StatusCode:        osResp.StatusCode,
		Header:            osResp.Header,
		Response:          resp,
		ValidationResults: vrs,
	}, nil
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

// toDomain converts this instance of [GetIndexResponse] into an [opensearchtools.GetIndexResponse]
func (g GetIndexResponse) toDomain() opensearchtools.GetIndexResponse {
	var resp map[string]opensearchtools.IndexInfo

	for k, v := range g.Response {
		settings := opensearchtools.IndexSetting{
			RefreshInterval:  v.Settings.Index.RefreshInterval,
			CreationDate:     v.Settings.Index.CreationDate,
			NumberOfShards:   v.Settings.Index.NumberOfShards,
			NumberOfReplicas: v.Settings.Index.NumberOfReplicas,
			UUID:             v.Settings.Index.UUID,
			Version:          v.Settings.Index.Version,
			ProvidedName:     v.Settings.Index.ProvidedName,
		}
		resp[k] = opensearchtools.IndexInfo{Aliases: v.Aliases, Mappings: v.Mappings, Settings: struct{ Index opensearchtools.IndexSetting }{Index: settings}}
	}

	domainResp := opensearchtools.GetIndexResponse{
		Response: resp,
	}
	return domainResp
}
