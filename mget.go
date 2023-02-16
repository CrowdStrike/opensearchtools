package opensearchtools

import (
	"context"
	"encoding/json"
	"net/http"
)

// MGet defines a method which knows how to make an OpenSearch multiple document get.
// It should be implemented by a version-specific executor.
type MGet interface {
	MGet(ctx context.Context, req *MGetRequest) (*MGetResponse, error)
}

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

// NewMGetRequest instantiates an empty [MGetRequest].
// An empty [MGetRequest] is executable but will return zero documents because zero have been requested.
func NewMGetRequest() *MGetRequest {
	return &MGetRequest{}
}

// Add a [DocumentRef] to the documents being requested.
// If index is an empty string, the request relies on the top-level MGetRequest.Index.
func (m *MGetRequest) Add(index, id string) *MGetRequest {
	return m.AddDocs(NewDocumentRef(index, id))
}

// AddDocs - add any number [RoutableDoc] to the documents being requested.
// If the doc does not return anything for [RoutableDoc.Index], the request relies on the top level [MGetRequest.Index].
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

// MGetResponse wraps the functionality of [opensearchapi.Response] by unmarshalling the response body into a
// slice of [MGetResults].
type MGetResponse struct {
	Header     http.Header
	StatusCode int
	Docs       []MGetResult
}

// MGetResult is the individual result for each requested item.
type MGetResult struct {
	Index       string
	ID          string
	Version     int
	SeqNo       int
	PrimaryTerm int
	Found       bool
	Source      json.RawMessage
	Error       error
}

// GetSource returns the raw bytes of the document of the [MGetResult].
func (m MGetResult) GetSource() []byte {
	return []byte(m.Source)
}
