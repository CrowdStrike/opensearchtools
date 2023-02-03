package opensearchtools

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

// MGetRequest wraps the functionality of [opensearchapi.MgetRequest] by supporting request body creation.
// We can perform an MGetRequest as simply as:
//
//	mgetResults, mgetError := NewMGetRequest().
//	    Add("example_index", "totally_real_id").
//	    Do(context.background(), client)
type MGetRequest struct {
	// Index destination for entire request
	// if used individual documents don't need to specify the index
	Index string

	// Docs are the list of documents to be fetched.
	Docs []RoutableDoc
}

// NewMGetRequest instantiates an empty MGetRequest.
// An empty MGetRequest is executable but will return zero documents because zero have been requested.
func NewMGetRequest() *MGetRequest {
	return &MGetRequest{}
}

// Add a DocumentRef to the documents being requested.
// If index is an empty string, the request relies on the top level MGetRequest.Index.
func (m *MGetRequest) Add(index, id string) *MGetRequest {
	return m.AddDocs(NewDocumentRef(index, id))
}

// AddDocs - add any number RoutableDoc to the documents being requested.
// If the doc does not return anything for [RoutableDoc.Index], the request relies on the top level MGetRequest.Index.
func (m *MGetRequest) AddDocs(docs ...RoutableDoc) *MGetRequest {
	m.Docs = append(m.Docs, docs...)
	return m
}

// SetIndex sets the top level index for the request. If a individual document request does not have an index specified,
// this index will be used.
func (m *MGetRequest) SetIndex(index string) *MGetRequest {
	m.Index = index
	return m
}

// Source translates the MGetRequest into the shape expected by OpenSearch.
func (m *MGetRequest) Source() any {
	docs := make([]any, len(m.Docs))
	for i, d := range m.Docs {
		docReq := make(map[string]any)

		if d.Index() != "" {
			docReq["_index"] = d.Index()
		}

		docReq["_id"] = d.ID()

		docs[i] = d
	}

	source := make(map[string]any)
	source["docs"] = docs

	return source
}

// Do executes the Multi-Get MGetRequest using the provided opensearch.Client.
// If the request is executed successfully, then a MGetResponse with MGetResults will be returned.
// An error can be returned if
//
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
func (m *MGetRequest) Do(ctx context.Context, client *opensearch.Client) (*MGetResponse, error) {
	bodyBytes, jErr := json.Marshal(m.Source())
	if jErr != nil {
		return nil, jErr
	}

	osResp, rErr := opensearchapi.MgetRequest{
		Index: m.Index,
		Body:  bytes.NewReader(bodyBytes),
	}.Do(ctx, client)

	if rErr != nil {
		return nil, rErr
	}

	resp := &MGetResponse{
		StatusCode: osResp.StatusCode,
		Header:     osResp.Header,
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBuf.Bytes(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// MGetResponse wraps the functionality of [opensearchapi.Response] by unmarshalling the response body into a
// slice of MGetResults.
type MGetResponse struct {
	StatusCode int          `json:"-"`
	Header     http.Header  `json:"-"`
	Docs       []MGetResult `json:"docs,omitempty"`
}

// MGetResult is the individual result for each requested item.
type MGetResult struct {
	Index       string          `json:"_index,omitempty"`
	ID          string          `json:"_id,omitempty"`
	Version     int             `json:"_version,omitempty"`
	SeqNo       int             `json:"_seq_no,omitempty"`
	PrimaryTerm int             `json:"_primary_term,omitempty"`
	Found       bool            `json:"found,omitempty"`
	Source      json.RawMessage `json:"_source,omitempty"`
	Error       error           `json:"-"`
}

// GetSource returns the raw bytes of the document of the MGetResult.
func (m MGetResult) GetSource() []byte {
	return []byte(m.Source)
}
