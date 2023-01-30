package os2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

// Request wraps the functionality of opensearchapi.MgetRequest to support a dynamic building of the request body.
// Minimalistic to cover desired functionality and not full functionality.
type Request[K opensearchtools.RoutableDoc] struct {
	// Index destination for entire request
	// if used individual documents don't need to specify the index
	Index string `json:"-"`

	Docs []Doc `json:"docs"`
}

// Add any number of docs to the Request
func (m *Request[K]) Add(docs ...K) {
	converted := make([]Doc, len(docs))
	for i, doc := range docs {
		converted[i] = Doc{
			ID:    doc.GetID(),
			Index: doc.RouteToIndex(),
		}
	}

	m.Docs = append(m.Docs, converted...)
}

// Do executes the Multi-Get Request using the provided opensearch.Client
func (m *Request[K]) Do(ctx context.Context, client *opensearch.Client) (*Response[K], error) {
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

	t := &ResponseBody{}

	if err := json.Unmarshal(respBuf.Bytes(), &t); err != nil {
		return nil, err
	}

	var docs []GetResult[K]
	for _, rawDoc := range t.Docs {
		doc := GetResult[K]{}

		if err := json.Unmarshal(rawDoc, &doc); err != nil {
			doc.Error = err
		}

		docs = append(docs, doc)
	}

	resp := &Response[K]{
		StatusCode: osResp.StatusCode,
		Header:     osResp.Header,
		Docs:       docs,
	}

	return resp, nil
}

// Doc represents individual document being requested in the multi-get
type Doc struct {
	ID    string `json:"_id,omitempty"`
	Index string `json:"_index,omitempty"`
}

// Response wraps the functionality of opensearchapi.Response leveraging K to deserialize the response
type Response[K opensearchtools.RoutableDoc] struct {
	StatusCode int
	Header     http.Header
	Docs       []GetResult[K]
}

// ResponseBody is an intermediary struct for the api response, to individually un marshall the response documents.
type ResponseBody struct {
	Docs []json.RawMessage `json:"docs"`
}

// GetResult is the individual result for each requested item
type GetResult[K opensearchtools.RoutableDoc] struct {
	Index       string `json:"_index,omitempty"`
	ID          string `json:"_id,omitempty"`
	Version     int    `json:"_version,omitempty"`
	SeqNo       int    `json:"_seq_no,omitempty"`
	PrimaryTerm int    `json:"_primary_term,omitempty"`
	Found       bool   `json:"found,omitempty"`
	Source      K      `json:"_source,omitempty"`
	Error       error  `json:"-"`
}
