package opensearchtools

import (
	"context"
	"encoding/json"

	"golang.org/x/exp/maps"
)

// Search defines a method which knows how to make an OpenSearch [Search] request.
// It should be implemented by a version-specific executor.
//
// [Search]: https://opensearch.org/docs/latest/api-reference/search/
type Search interface {
	Search(ctx context.Context, req *SearchRequest) (OpenSearchResponse[SearchResponse], error)
}

// SearchRequest is a domain model union type for all the fields of a Search request across
// all supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// This SearchRequest is intended to be used along with a version-specific executor such
// as [opensearchtools/osv2.Executor]. For example:
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

	// From the starting index to search from
	From int

	// Sort(s) to order the results returned
	Sort []Sort

	// TrackTotalHits - whether to return how many documents matched the query.
	TrackTotalHits any

	// Routing - Values used to route the update by query operation to a specific shard
	Routing []string

	// Aggregations to be performed on the results of the Query
	Aggregations map[string]Aggregation
}

// NewSearchRequest instantiates a SearchRequest with a From and Size of -1.
// Any negative value for SearchRequest.From or SearchRequest.Size will be ignored and not included in the source.
// Opensearch by default, if no size is included in a search request, will limit the results to 10 documents.
// Opensearch by default, if no from is included in a search request, will return hits starting from the first hit based on the sort.
// A NewSearchRequest will search across all indices and return the top 10 documents with the default [sorting].
//
// [sorting]: https://openorg/docs/latest/opensearch/search/sort/
func NewSearchRequest() *SearchRequest {
	return &SearchRequest{Size: -1, From: -1}
}

// AddIndices sets the index list for the request.
func (r *SearchRequest) AddIndices(indices ...string) *SearchRequest {
	r.Index = append(r.Index, indices...)
	return r
}

// WithSize sets the request size, limiting the number of documents returned.
// A negative value for size will be ignored and not included in the SearchRequest.Source.
func (r *SearchRequest) WithSize(n int) *SearchRequest {
	r.Size = n
	return r
}

// WithFrom sets the request's starting index for the result hits.
// A negative value for from will be ignored and not included in the SearchRequest.Source.
func (r *SearchRequest) WithFrom(n int) *SearchRequest {
	r.From = n
	return r
}

// AddSorts to the current list of [Sort]s on the request.
func (r *SearchRequest) AddSorts(sort ...Sort) *SearchRequest {
	r.Sort = append(r.Sort, sort...)
	return r
}

// WithQuery to be performed by the SearchRequest.
func (r *SearchRequest) WithQuery(q Query) *SearchRequest {
	r.Query = q
	return r
}

// WithTrackTotalHits if set to true it will count all documents,
// otherwise a number can be set to limit the counting ceiling.
func (r *SearchRequest) WithTrackTotalHits(track any) *SearchRequest {
	r.TrackTotalHits = track
	return r
}

// WithRouting sets the routing value(s)
func (r *SearchRequest) WithRouting(routing ...string) *SearchRequest {
	r.Routing = routing
	return r
}

// AddAggregation to the search request with the desired name
func (r *SearchRequest) AddAggregation(name string, agg Aggregation) *SearchRequest {
	if r.Aggregations == nil {
		r.Aggregations = map[string]Aggregation{name: agg}
	} else {
		r.Aggregations[name] = agg
	}

	return r
}

// SearchResponse is a domain model union response type across all supported OpenSearch versions.
// Currently supported versions are:
//
//	-OpenSearch2
type SearchResponse struct {
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

	// Aggregations response if any were requested
	Aggregations map[string]json.RawMessage
}

// GetAggregationResultSource implements [opensearchtools.AggregationResultSet] to fetch an aggregation result and
// return the raw JSON source for the provided name.
func (sr SearchResponse) GetAggregationResultSource(name string) ([]byte, bool) {
	if len(sr.Aggregations) == 0 {
		return nil, false
	}

	aggSource, exists := sr.Aggregations[name]
	return aggSource, exists
}

// Keys implemented for [opensearchtools.AggregationResultSet] to return the list of aggregation result keys
func (sr SearchResponse) Keys() []string {
	return maps.Keys(sr.Aggregations)
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
