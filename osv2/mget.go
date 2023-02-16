package osv2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

// MGetRequest is a marshalable form of [opensearchtools.MGetRequest] specific to the opensearchapi.MgetRequest in OpenSearch v2.
type MGetRequest struct {
	// Index destination for entire request
	// if used individual documents don't need to specify the index
	Index string

	// Docs are the list of documents to be fetched.
	Docs []opensearchtools.RoutableDoc
}

// Do executes the Multi-Get MGetRequest using the provided opensearch.Client.
// If the request is executed successfully, then a MGetResponse with MGetResults will be returned.
// An error can be returned if
//
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
func (m *MGetRequest) Do(ctx context.Context, client *opensearch.Client) (*MGetResponse, error) {
	bodyBytes, jErr := json.Marshal(m)
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

// FromModelMGetRequest creates a new [mgetRequest] from the given [opensearchtools.MGetRequest].
func FromModelMGetRequest(req *opensearchtools.MGetRequest) *MGetRequest {
	return &MGetRequest{
		Index: req.Index,
		Docs:  req.Docs,
	}
}

// MarshalJSON marshals the [MGetRequest] into the proper json expected by OpenSearch 2.
func (m *MGetRequest) MarshalJSON() ([]byte, error) {
	docs := make([]any, len(m.Docs))
	for i, d := range m.Docs {
		docReq := map[string]any{
			"_id": d.ID(),
		}

		if d.Index() != "" {
			docReq["_index"] = d.Index()
		}

		docs[i] = docReq
	}

	source := map[string]any{
		"docs": docs,
	}

	return json.Marshal(source)
}

// MGetResponse is an OpenSearch 2 specific struct corresponding to opensearchapi.Response and [opensearchtools.MGetResponse].
// It holds a slice of mgetResults.
type MGetResponse struct {
	StatusCode int          `json:"-"`
	Header     http.Header  `json:"-"`
	Docs       []MGetResult `json:"docs,omitempty"`
}

// ToModel converts this instance of an [MGetResponse] into an [opensearchtools.MGetResponse].
func (r *MGetResponse) ToModel() *opensearchtools.MGetResponse {
	modelDocs := make([]opensearchtools.MGetResult, len(r.Docs))
	for i, d := range r.Docs {
		modelDoc := d.ToModel()
		modelDocs[i] = modelDoc
	}

	return &opensearchtools.MGetResponse{
		StatusCode: r.StatusCode,
		Header:     r.Header,
		Docs:       modelDocs,
	}
}

// mgetResult is the individual result for each requested item.
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

// ToModel converts this instance of an [MGetResult] into an [opensearchtools.MGetResult].
func (r *MGetResult) ToModel() opensearchtools.MGetResult {
	return opensearchtools.MGetResult{
		Index:       r.Index,
		ID:          r.ID,
		Version:     r.Version,
		SeqNo:       r.SeqNo,
		PrimaryTerm: r.PrimaryTerm,
		Found:       r.Found,
		Source:      r.Source,
		Error:       r.Error,
	}
}
