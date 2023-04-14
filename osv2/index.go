package osv2

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

type IndexRequest struct {
	Action opensearchtools.IndexAction

	Indices []string

	Refresh opensearchtools.Refresh

	Routing string
}

// NewIndexRequest instantiates an empty IndexRequest
func NewIndexRequest() *BulkRequest {
	return &BulkRequest{}
}

// Add an index action to the IndexRequest.
func (r *IndexRequest) Add(action opensearchtools.IndexAction) *IndexRequest {
	r.Action = action
	return r
}

// WithIndices on the request
func (r *IndexRequest) WithIndices(index ...string) *IndexRequest {
	r.Indices = index
	return r
}

func (r *IndexRequest) Do(ctx context.Context, client *opensearch.Client) (*opensearchtools.OpenSearchResponse[IndexResponse], error) {
	rawBody, jErr := r.ToOpenSearchJSON()
	if jErr != nil {
		return nil, jErr
	}
	var osResp *opensearchapi.Response
	var rErr error
	switch r.Action.Type {
	case opensearchtools.IndexCreate:
		osResp, rErr = opensearchapi.IndicesCreateRequest{
			Index: r.Indices[0], // check it please
			Body:  bytes.NewReader(rawBody),
		}.Do(ctx, client)
	case opensearchtools.IndexDelete:
		osResp, rErr = opensearchapi.IndicesDeleteRequest{
			Index: r.Indices,
		}.Do(ctx, client)
	case opensearchtools.IndexExists:
		osResp, rErr = opensearchapi.IndicesExistsRequest{
			Index: r.Indices,
		}.Do(ctx, client)
	case opensearchtools.IndexGet:
		osResp, rErr = opensearchapi.IndicesGetRequest{
			Index: r.Indices,
		}.Do(ctx, client)
	}

	if rErr != nil {
		return nil, rErr
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	resp := IndexResponse{}

	if err := json.Unmarshal(respBuf.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &opensearchtools.OpenSearchResponse[IndexResponse]{
		StatusCode: osResp.StatusCode,
		Header:     osResp.Header,
		Response:   resp,
	}, nil
}

func (r *IndexRequest) ToOpenSearchJSON() ([]byte, error) {
	bodyBuf := new(bytes.Buffer)
	// parsing
	return bodyBuf.Bytes(), nil
}

type IndexResponse struct {
	Acknowledged *bool
	Error        *Error
	*Indices
}

type Indices map[string]IndexInfo

type IndexInfo struct {
	aliases  map[string]json.RawMessage
	mappings map[string]json.RawMessage
	Settings *IndexSettings
}
type IndexSettings struct {
	CreationDate     string
	NumberOfShards   string
	NumberOfReplicas string
	UUID             string
	Version          struct{ Created string }
	ProvidedName     string
}
