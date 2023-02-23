package opensearchtools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MGet defines a method which knows how to make an OpenSearch [Multi-get] request.
// It should be implemented by a version-specific executor.
//
// [Multi-get]: https://opensearch.org/docs/latest/api-reference/document-apis/multi-get/
type MGet interface {
	MGet(ctx context.Context, req *MGetRequest) (*MGetResponse, error)
}

// MGetRequest is a domain model union type for all the fields of a Multi-Get request for all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// This MGetRequest is intended to be used along with a version-specific executor such as
// [opensearchtools.osv2.Executor]. For example:
//
//	mgetReq := NewMGetRequest().
//		Add("example_index", "example_id")
//	mgetResp, err := osv2Executor.MGet(ctx, mgetReq)
//
// An error can be returned if
//
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
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

// WithIndex sets the top level index for the request. If a individual document request does not have an index specified,
// this index will be used.
func (m *MGetRequest) WithIndex(index string) *MGetRequest {
	m.Index = index
	return m
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

// MGetResponse is a domain model union response type for Multi-Get for all supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// Contains a slice of [MGetResult] for each document in the response.
type MGetResponse struct {
	Header     http.Header
	StatusCode int
	Docs       []MGetResult
}

// MGetResult is a domain model union type representing an individual document result in for a request item
// in a Multi-Get response for all supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
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

// Validate validates the given MGetRequest
func (m *MGetRequest) Validate() ValidationResults {
	validationResults := make(ValidationResults, 0)

	topLevelIndexIsEmpty := m.Index == ""
	for _, d := range m.Docs {
		// ensure Index is either set at the top level or set in each of the Docs
		if topLevelIndexIsEmpty && d.Index() == "" {
			validationResults = append(validationResults, NewValidationResult(fmt.Sprintf("Index not set at the MGetRequest level nor in the Doc with ID %s", d.ID()), true))
		}

		// ensure that ID() is non-empty for each Doc
		if d.ID() == "" {
			validationResults = append(validationResults, NewValidationResult("Doc ID is empty", true))
		}
	}

	return validationResults
}
