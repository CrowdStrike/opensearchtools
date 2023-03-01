package osv2

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

// SearchRequest is a serializable form of [opensearchtools.SearchRequest] specific to the [opensearchapi.SearchRequest] in OpenSearch V2.
// An empty SearchRequest defaults to a size of 0. While this will find matches and return a total hits value,
// it will return no documents. It is recommended to use NewSearchRequest or use WithSize.
// A simple term query search as an example:
//
//	req := NewSearchRequest()
//	req.AddIndices("example_index")
//	req.WithQuery(opensearchtools.NewTermQuery("field", "basic")
//	results, err := req.Do(context.Background(), client)
type SearchRequest struct {
	// Query to be performed by the search
	Query opensearchtools.Query

	// Index(s) to be targeted by the search
	Index []string

	// Size of results to be returned
	Size int

	// Sort(s) to order the results returned
	Sort []opensearchtools.Sort
}

// V2QueryConverter will do any translations needed from domain level queries into V2 specifics, if needed.
func V2QueryConverter(query opensearchtools.Query) (opensearchtools.Query, error) {
	switch q := query.(type) {
	case *opensearchtools.BoolQuery:
		return opensearchtools.BoolQueryConverter(q, V2QueryConverter)
	default:
		return q, nil
	}
}

// NewSearchRequest instantiates a SearchRequest with a Size of -1.
// Any negative value for SearchRequest.Size will be ignored and not included in the source.
// Opensearch by default, if no size is included in a search request, will limit the results to 10 documents.
// A NewSearchRequest will search across all indices and return the top 10 documents with the default [sorting].
//
// [sorting]: https://openopensearchtools.org/docs/latest/opensearch/search/sort/
func NewSearchRequest() *SearchRequest {
	return &SearchRequest{Size: -1}
}

// ToOpenSearchJSON marshals the SearchRequest into the JSON shape expected by OpenSearch.
func (r *SearchRequest) ToOpenSearchJSON() ([]byte, error) {
	source := make(map[string]any)
	if r.Query != nil {
		queryJSON, jErr := r.Query.ToOpenSearchJSON()
		if jErr != nil {
			return nil, jErr
		}

		source["query"] = json.RawMessage(queryJSON)
	}

	if r.Size >= 0 {
		source["size"] = r.Size
	}

	if len(r.Sort) > 0 {
		sorts := make([]json.RawMessage, len(r.Sort))
		for i, s := range r.Sort {
			sortJSON, jErr := s.ToOpenSearchJSON()
			if jErr != nil {
				return nil, jErr
			}

			sorts[i] = sortJSON
		}

		source["sort"] = sorts
	}

	return json.Marshal(source)
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

// AddSorts to the current list of [opensearchtools.Sort]s on the request.
func (r *SearchRequest) AddSorts(sort ...opensearchtools.Sort) *SearchRequest {
	r.Sort = append(r.Sort, sort...)
	return r
}

// WithQuery to be performed by the SearchRequest.
func (r *SearchRequest) WithQuery(q opensearchtools.Query) *SearchRequest {
	r.Query = q
	return r
}

// fromDomainSearchRequest creates a new SearchRequest from the given [opensearchtools.SearchRequest]
func fromDomainSearchRequest(req *opensearchtools.SearchRequest) (SearchRequest, opensearchtools.ValidationResults) {
	vrs := opensearchtools.NewValidationResults()
	var searchRequest SearchRequest

	convertedQuery, err := V2QueryConverter(req.Query)
	if err != nil {
		vrs.Add(opensearchtools.NewValidationResult(err.Error(), true))
		return searchRequest, vrs
	}

	return SearchRequest{
		Query: convertedQuery,
		Index: req.Index,
		Size:  req.Size,
		Sort:  req.Sort,
	}, vrs
}

// Validate validates the given SearchRequest
func (r *SearchRequest) Validate() opensearchtools.ValidationResults {
	var validationResults opensearchtools.ValidationResults
	return validationResults
}

// Do executes the SearchRequest using the provided [opensearch.Client].
// If the request is executed successfully, then a SearchResponse will be returned.
// An error can be returned if
//
//   - The SearchRequest source cannot be created
//   - The source fails to be marshaled to JSON
//   - The OpenSearch request fails to executed
//   - The OpenSearch response cannot be parsed
func (r *SearchRequest) Do(ctx context.Context, client *opensearch.Client) (*opensearchtools.OpenSearchResponse[SearchResponse], error) {
	bodyBytes, jErr := r.ToOpenSearchJSON()
	if jErr != nil {
		return nil, jErr
	}

	osResp, rErr := opensearchapi.SearchRequest{
		Index: r.Index,
		Body:  bytes.NewReader(bodyBytes),
	}.Do(ctx, client)

	if rErr != nil {
		return nil, rErr
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(respBuf.Bytes(), &searchResp); err != nil {
		return nil, err
	}

	resp := opensearchtools.NewOpenSearchResponse(
		opensearchtools.NewValidationResults(), // no additional validation
		osResp.StatusCode,
		osResp.Header,
		searchResp,
	)
	return &resp, nil
}

// SearchResponse wraps the functionality of [opensearchapi.Response] by supporting request parsing.
type SearchResponse struct {
	Took     int       `json:"took"`
	TimedOut bool      `json:"timed_out"`
	Shards   ShardMeta `json:"_shards,omitempty"`
	Hits     Hits      `json:"hits"`
	Error    *Error    `json:"error,omitempty"`
}

// ToDomain converts this instance of a [SearchResponse] into an [opensearchtools.SearchResponse].
func (sr *SearchResponse) ToDomain() opensearchtools.SearchResponse {
	domainResp := opensearchtools.SearchResponse{
		Took:     sr.Took,
		TimedOut: sr.TimedOut,
		Shards:   sr.Shards.toDomain(),
		Hits:     sr.Hits.toDomain(),
	}

	if sr.Error != nil {
		domainErr := sr.Error.toDomain()
		domainResp.Error = &domainErr
	}

	return domainResp
}

// Hits represent the results of the [opensearchtools.Query] performed by the SearchRequest.
type Hits struct {
	Total    Total   `json:"total,omitempty"`
	MaxScore float64 `json:"max_score,omitempty"`
	Hits     []Hit   `json:"hits"`
}

// toDomain converts this instance of a [Hits] into an [opensearchtools.Hits].
func (h Hits) toDomain() opensearchtools.Hits {
	var hits []opensearchtools.Hit
	for _, hit := range h.Hits {
		hits = append(hits, hit.toDomain())
	}

	return opensearchtools.Hits{
		Total:    h.Total.toDomain(),
		MaxScore: h.MaxScore,
		Hits:     hits,
	}
}

// Total contains the total number of documents found by the [opensearchtools.Query] performed by the SearchRequest.
type Total struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

// toDomain converts this instance of a [Total] into an [opensearchtools.Total].
func (t Total) toDomain() opensearchtools.Total {
	return opensearchtools.Total{
		Value:    t.Value,
		Relation: t.Relation,
	}
}

// Hit the individual document found by the `[opensearchtools.Query] performed by the SearchRequest.
type Hit struct {
	Index  string          `json:"_index"`
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

// toDomain converts this instance of a [Hit] into an [opensearchtools.Hit].
func (h Hit) toDomain() opensearchtools.Hit {
	return opensearchtools.Hit{
		Index:  h.Index,
		ID:     h.ID,
		Score:  h.Score,
		Source: h.Source,
	}
}

// GetSource returns the raw bytes of the document of the MGetResult.
func (h Hit) GetSource() []byte {
	return []byte(h.Source)
}
