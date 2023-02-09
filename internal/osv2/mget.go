package osv2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

// The main purpose of this file is to implement the MGet interface on the OSv2Executor.
// In order to do this, it also defines a bunch of marshalable types for making requests and getting
// responses from OpenSearch 2. These might have a different shape than the similarily named ones in the root mget.

// MGet executes the Multi-Get MGetRequest using the provided opensearchtools.MGetRequest.
// If the request is executed successfully, then a MGetResponse with MGetResults will be returned.
// An error can be returned if
//
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
func (e *OSv2Executor) MGet(ctx context.Context, req *opensearchtools.MGetRequest) (*opensearchtools.MGetResponse, error) {
	// first create a serform mgetRequest
	specReq := FromModelMGetRequest(req)

	// create json body bytes for that MGetRequest
	bodyBytes, jErr := json.Marshal(specReq)
	if jErr != nil {
		return nil, jErr
	}

	osResp, rErr := opensearchapi.MgetRequest{
		Index: specReq.Index,
		Body:  bytes.NewReader(bodyBytes),
	}.Do(ctx, e.client)

	if rErr != nil {
		return nil, rErr
	}

	resp := &mgetResponse{
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

	modelResp := resp.ToModel()

	return modelResp, nil
}

// MGetRequest is a marshalable form of opensearchtools.MGet specific to the opensearchapi.MgetRequest in OpenSearch v2.
type mgetRequest struct {
	// Index destination for entire request
	// if used individual documents don't need to specify the index
	Index string

	// Docs are the list of documents to be fetched.
	Docs []opensearchtools.RoutableDoc
}

// FromModelMGetRequest creates a new MGetRequest from the given opensearchtools.MGetRequest.
func FromModelMGetRequest(req *opensearchtools.MGetRequest) *mgetRequest {
	return &mgetRequest{
		Index: req.Index,
		Docs:  req.Docs,
	}
}

// MarshalJSON marshals the mgetRequest into the proper json expected by OpenSearch 2.
func (m *mgetRequest) MarshalJSON() ([]byte, error) {
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

// MGetResponse wraps the functionality of [opensearchapi.Response] by unmarshalling the response body into a
// slice of MGetResults.

// mgetResponse is an OpenSearch 2 specific struct corresponding to [opensearchapi.Response] and [opensearchtools.MGetResponse].
// It holds a slice of mgetResults.
type mgetResponse struct {
	StatusCode int          `json:"-"`
	Header     http.Header  `json:"-"`
	Docs       []mgetResult `json:"docs,omitempty"`
}

// ToModel converts this instance of an mgetResponse into an opensearchtools.MGetResponse.
func (r *mgetResponse) ToModel() *opensearchtools.MGetResponse {
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
type mgetResult struct {
	Index       string          `json:"_index,omitempty"`
	ID          string          `json:"_id,omitempty"`
	Version     int             `json:"_version,omitempty"`
	SeqNo       int             `json:"_seq_no,omitempty"`
	PrimaryTerm int             `json:"_primary_term,omitempty"`
	Found       bool            `json:"found,omitempty"`
	Source      json.RawMessage `json:"_source,omitempty"`
	Error       error           `json:"-"`
}

// / ToModel converts this instance of an mgetResult into an opensearchtools.MGetResult.
func (r *mgetResult) ToModel() opensearchtools.MGetResult {
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
