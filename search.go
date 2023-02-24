package opensearchtools

import (
	"context"
	"encoding/json"
	"net/http"
)

// Search defines a method which knows how to make an OpenSearch [Search] request.
// It should be implemented by a version-specific executor.
//
// [Search]: https://openorg/docs/latest/api-reference/search/
type Search interface {
	Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
}

// SearchRequest is a domain model union type for all the fields of a Search request across
// all supported OpenSearch version.
// Currently supported versions are:
//   - OpenSearch 2
//
// This SearchRequest is intended to be used along with a version-specific executor such
// as [opensearchtools.osv2.Executor]. For example:
//
//	searchReq := NewSearchRequest().
//		WithQuery(NewTermQuery("field", "term"))
//	searchResp, err := osv2Executor.Search(ctx, searchReq)
//
// An error can be returned if
//   - The request to OpenSearch fails
//   - The results json cannot be unmarshalled
type SearchRequest struct {
	// Query to be performed by the search
	Query Query

	// Index(s) to be targeted by the search
	Index []string

	// Size of results to be returned
	Size int

	// Sort(s) to order the results returned
	Sort []*Sort
}

// NewSearchRequest instantiates a SearchRequest with a Size of -1.
// Any negative value for SearchRequest.Size will be ignored and not included in the source.
// Opensearch by default, if no size is included in a search request, will limit the results to 10 documents.
// A NewSearchRequest will search across all indices and return the top 10 documents with the default [sorting].
//
// [sorting]: https://openorg/docs/latest/opensearch/search/sort/
func NewSearchRequest() *SearchRequest {
	return &SearchRequest{Size: -1}
}

// WithIndices sets the index list for the request.
func (r *SearchRequest) WithIndices(indices ...string) *SearchRequest {
	r.Index = append(r.Index, indices...)
	return r
}

// WithSize sets the request size, limiting the number of documents returned.
// A negative value for size will be ignored and not included in the SearchRequest.Source.
func (r *SearchRequest) WithSize(n int) *SearchRequest {
	r.Size = n
	return r
}

// WithSorts to the current list of [Sort]s on the request.
func (r *SearchRequest) WithSorts(sort ...*Sort) *SearchRequest {
	r.Sort = append(r.Sort, sort...)
	return r
}

// WithQuery to be performed by the SearchRequest.
func (r *SearchRequest) WithQuery(q Query) *SearchRequest {
	r.Query = q
	return r
}

// SearchResponse is a domain model union response type across all supported OpenSearch versions.
// Currently supported versions are:
//
//	-OpenSearch2
type SearchResponse struct {
	// StatusCode of the http request
	StatusCode int

	// Header details returned by OpenSearch
	Header http.Header

	// Took the time in Milliseconds OpenSearch took to execute the query
	Took int

	// TimedOut true if the request timed out
	TimedOut bool

	// Shards [ShardMeta] counts of the shards used in the request processing
	Shards ShardMeta

	// Hits are the results of the [Query]
	Hits Hits

	// Error if OpenSearch failed but responded with errors
	Error *Error
}

// Hits is a domain model union response type across all supported OpenSearch versions.
// Currently supported versions are:
//
//	-OpenSearch2
//
// Hits are the list of documents hit from the executed [Query].
type Hits struct {
	// Total documents that matched the query.
	Total Total

	// MaxScore max score of all matching documents
	MaxScore float64

	// Hits slice of documents matched by the [Query]
	Hits []Hit
}

// Total is a domain model union response type across all supported OpenSearch versions.
// Currently supported versions are:
//
//	-OpenSearch2
//
// Total contains the total number of documents found by the [Query] performed by the SearchRequest.
type Total struct {
	Value    int64
	Relation string
}

// Hit is a domain model union response type across all supported OpenSearch versions.
// Currently supported versions are:
//
//	-OpenSearch2
//
// Hit the individual document found by the `[Query] performed by the SearchRequest.
type Hit struct {
	Index  string
	ID     string
	Score  float64
	Source json.RawMessage
}

// GetSource returns the raw bytes of the document of the SearchRequest.
func (h Hit) GetSource() []byte {
	return []byte(h.Source)
}
