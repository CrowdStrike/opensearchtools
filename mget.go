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
// This is done with a generic K [RoutableDoc] letting you [Add] documents to be requested.
// For example, if we use the BasicDoc example from [RoutableDoc] we can leverage an MGetRequest as simply as:
//
//	doc := BasicDoc{ID: "sampleId", Index: "index"}
//	mgetResults, mgetError := NewMGetRequest[BasicDoc]().Add(doc).Do(context.background(), client)
type MGetRequest[K RoutableDoc] struct {
    // Index destination for entire request
    // if used individual documents don't need to specify the index
    Index string `json:"-"`

    Docs []mgetDoc `json:"docs"`
}

// NewMGetRequest instantiates an empty MGetRequest for declared [RoutableDoc] K.
// An empty MGetRequest is executable but will return zero documents because zero have been requested.
func NewMGetRequest[K RoutableDoc]() *MGetRequest[K] {
    return &MGetRequest[K]{}
}

// Add any number of docs to the MGetRequest
func (m *MGetRequest[K]) Add(docs ...K) *MGetRequest[K] {
    converted := make([]mgetDoc, len(docs))
    for i, doc := range docs {
        converted[i] = mgetDoc{
            ID:    doc.ID(),
            Index: doc.Index(),
        }
    }

    m.Docs = append(m.Docs, converted...)

    return m
}

// Do executes the Multi-Get MGetRequest using the provided opensearch.Client.
// If the request is executed successfully, then a MGetResponse with MGetResults will be returned.
// An error can be returned if
//
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
//   - An individual document response cannot be unmarshalled into K
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

    t := &mgetResponseBody{}

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

// mgetDoc is a simple struct representing an individual document to be fetched in a MultiGet request
type mgetDoc struct {
    ID    string `json:"_id,omitempty"`
    Index string `json:"_index,omitempty"`
}

// MGetResponse wraps the functionality of [opensearchapi.Response] by supporting request unmarshalling of the found
// documents into K
type MGetResponse[K RoutableDoc] struct {
    StatusCode int
    Header     http.Header
    Docs       []MGetResult[K]
}

// mgetResponseBody is an intermediary struct for the api response, to individually un marshall the response documents.
type mgetResponseBody struct {
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
