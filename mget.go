package opensearchtools

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

// MGetRequest wraps the functionality of opensearchapi.MgetRequest to support a dynamic building of the request body.
// Minimalistic to cover desired functionality and not full functionality.
type MGetRequest[K RoutableDoc] struct {
	// Index destination for entire request
	// if used individual documents don't need to specify the index
	Index string `json:"-"`

	Docs []MGetDoc `json:"docs"`
}

// Add any number of docs to the MGetRequest
func (m *MGetRequest[K]) Add(docs ...K) {
	converted := make([]MGetDoc, len(docs))
	for i, doc := range docs {
		converted[i] = MGetDoc{
			ID:    doc.GetID(),
			Index: doc.RouteToIndex(),
		}
	}

	m.Docs = append(m.Docs, converted...)
}

// Do executes the Multi-Get MGetRequest using the provided opensearch.Client
func (m *MGetRequest[K]) Do(ctx context.Context, client *opensearch.Client) (*MGetResponse[K], error) {
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

	t := &MGetResponseBody{}

	if err := json.Unmarshal(respBuf.Bytes(), &t); err != nil {
		return nil, err
	}

	var docs []MGetResult[K]
	for _, rawDoc := range t.Docs {
		doc := MGetResult[K]{}

		if err := json.Unmarshal(rawDoc, &doc); err != nil {
			doc.Error = err
		}

		docs = append(docs, doc)
	}

	resp := &MGetResponse[K]{
		StatusCode: osResp.StatusCode,
		Header:     osResp.Header,
		Docs:       docs,
	}

	return resp, nil
}

// MGetDoc represents individual document being requested in the multi-get
type MGetDoc struct {
	ID    string `json:"_id,omitempty"`
	Index string `json:"_index,omitempty"`
}

// MGetResponse wraps the functionality of opensearchapi.Response leveraging K to deserialize the response
type MGetResponse[K RoutableDoc] struct {
	StatusCode int
	Header     http.Header
	Docs       []MGetResult[K]
}

// MGetResponseBody is an intermediary struct for the api response, to individually un marshall the response documents.
type MGetResponseBody struct {
	Docs []json.RawMessage `json:"docs"`
}

// MGetResult is the individual result for each requested item
type MGetResult[K RoutableDoc] struct {
	Index       string `json:"_index,omitempty"`
	ID          string `json:"_id,omitempty"`
	Version     int    `json:"_version,omitempty"`
	SeqNo       int    `json:"_seq_no,omitempty"`
	PrimaryTerm int    `json:"_primary_term,omitempty"`
	Found       bool   `json:"found,omitempty"`
	Source      K      `json:"_source,omitempty"`
	Error       error  `json:"-"`
}
