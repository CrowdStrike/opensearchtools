package osv2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

// MGetRequest is a marshalable form of [opensearchtools.MGetRequest] specific to the opensearchapi.MgetRequest in OpenSearch v2.
//
// [Multi-get]: https://opensearch.org/docs/latest/api-reference/document-apis/multi-get/
type MGetRequest struct {
	// Index destination for entire request
	// if used individual documents don't need to specify the index
	Index string

	// Docs are the list of documents to be fetched.
	Docs []opensearchtools.RoutableDoc
}

// NewMGetRequest instantiates an empty [MGetRequest].
// An empty [MGetRequest] is executable but will return zero documents because zero have been requested.
func NewMGetRequest() *MGetRequest {
	return &MGetRequest{}
}

// WithIndex sets the top level index for the request. If a individual document request does not have an index specified,
// this index will be used.
func (m *MGetRequest) WithIndex(index string) *MGetRequest {
	m.Index = index
	return m
}

// Add a [opensearchtools.DocumentRef] to the documents being requested.
// If index is an empty string, the request relies on the top-level MGetRequest.Index.
func (m *MGetRequest) Add(index, id string) *MGetRequest {
	return m.AddDocs(opensearchtools.NewDocumentRef(index, id))
}

// AddDocs - add any number [opensearchtools.RoutableDoc] to the documents being requested.
// If the doc does not return anything for [RoutableDoc.Index], the request relies on the top level [MGetRequest.Index].
func (m *MGetRequest) AddDocs(docs ...opensearchtools.RoutableDoc) *MGetRequest {
	m.Docs = append(m.Docs, docs...)
	return m
}

// Do executes the Multi-Get MGetRequest using the provided OpenSearch v2 [opensearch.Client].
// If the request is executed successfully, then a MGetResponse with MGetResults will be returned.
// We can perform an MGetRequest as simply as:
//
//	mgetResults, mgetError := NewMGetRequest().
//	    Add("example_index", "example_id").
//	    Do(context.background(), client)
//
// An error can be returned if
//
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
func (m *MGetRequest) Do(ctx context.Context, client *opensearch.Client) (*opensearchtools.OpenSearchResponse[MGetResponse], error) {
	vrs := m.validate()
	if vrs.IsFatal() {
		return nil, opensearchtools.NewValidationError(vrs)
	}

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

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	var mgetResp MGetResponse
	if err := json.Unmarshal(respBuf.Bytes(), &mgetResp); err != nil {
		return nil, err
	}

	resp := opensearchtools.NewOpenSearchResponse(
		vrs,
		osResp.StatusCode,
		osResp.Header,
		mgetResp,
	)
	return &resp, nil
}

// FromDomainMGetRequest creates a new [MGetRequest] from the given [opensearchtools.MGetRequest].
func FromDomainMGetRequest(req *opensearchtools.MGetRequest) (MGetRequest, opensearchtools.ValidationResults) {
	return MGetRequest{
		Index: req.Index,
		Docs:  req.Docs,
	}, opensearchtools.NewValidationResults()
}

// validate validates the given MGetRequest
func (m *MGetRequest) validate() opensearchtools.ValidationResults {
	validationResults := opensearchtools.NewValidationResults()

	topLevelIndexIsEmpty := m.Index == ""
	for _, d := range m.Docs {
		// ensure Index is either set at the top level or set in each of the Docs
		if topLevelIndexIsEmpty && d.Index() == "" {
			validationResults.Add(opensearchtools.NewValidationResult(fmt.Sprintf("Index not set at the MGetRequest level nor in the Doc with ID %s", d.ID()), true))
		}

		// ensure that ID() is non-empty for each Doc
		if d.ID() == "" {
			validationResults.Add(opensearchtools.NewValidationResult("Doc ID is empty", true))
		}
	}

	return validationResults
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
	Docs []MGetResult `json:"docs,omitempty"`
}

// toDomain converts this instance of an [MGetResponse] to an [opensearchtools.MGetResponse]
func (r *MGetResponse) toDomain() opensearchtools.MGetResponse {
	modelDocs := make([]opensearchtools.MGetResult, len(r.Docs))
	for i, d := range r.Docs {
		modelDoc := d.toDomain()
		modelDocs[i] = modelDoc
	}

	return opensearchtools.MGetResponse{
		Docs: modelDocs,
	}
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

// toDomain converts this instance of an [MGetResult] into an [opensearchtools.MGetResult].
func (r *MGetResult) toDomain() opensearchtools.MGetResult {
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
