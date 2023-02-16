package opensearchtools

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools/search"
)

// SearchRequest wraps the functionality of [opensearchapi.SearchRequest] by supporting request body creation.
// An empty SearchRequest defaults to a size of 0. While this will find matches and return a total hits value,
// it will return no documents. It is recommended to use NewSearchRequest or use SetSize.
// A simple term query search as an example:
//
//	req := NewSearchRequest()
//	req.AddIndices("example_index")
//	req.SetQuery(search.NewTermQuery("field", "basic")
//	results, err := req.Do(context.Background(), client)
type SearchRequest struct {
	Query search.Query
	Index []string
	Size  int
	Sort  []*search.Sort
}

// NewSearchRequest instantiates a SearchRequest with a Size of -1.
// Any negative value for SearchRequest.Size will be ignored and not included in the source.
// Opensearch by default, if no size is included in a search request, will limit the results to 10 documents.
// A NewSearchRequest will search across all indices and return the top 10 documents with the default [sorting].
//
// [sorting]: https://opensearch.org/docs/latest/opensearch/search/sort/
func NewSearchRequest() *SearchRequest {
	return &SearchRequest{Size: -1}
}

// MarshalJSON marshals the SearchRequest into the JSON shape expected by OpenSearch.
// Implemented to match the [json.Marshaler] interface.
func (r *SearchRequest) MarshalJSON() ([]byte, error) {
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

// SetSize sets the request size, limiting the number of documents returned.
// A negative value for size will be ignored and not included in the SearchRequest.Source.
func (r *SearchRequest) SetSize(n int) *SearchRequest {
	r.Size = n
	return r
}

// AddSort to the current list of [search.Sort]s on the request.
func (r *SearchRequest) AddSort(sort ...*search.Sort) *SearchRequest {
	r.Sort = append(r.Sort, sort...)
	return r
}

// SetQuery to be performed by the SearchRequest.
func (r *SearchRequest) SetQuery(q search.Query) *SearchRequest {
	r.Query = q
	return r
}

// Do executes the SearchRequest using the provided opensearch.Client.
// If the request is executed successfully, then a SearchResponse will be returned.
// An error can be returned if
//
//   - The SearchRequest source cannot be created
//   - The source fails to be marshaled to JSON
//   - The OpenSearch request fails to executed
//   - The OpenSearch response cannot be parsed
func (r *SearchRequest) Do(ctx context.Context, client *opensearch.Client) (*SearchResponse, error) {
	bodyBytes, jErr := json.Marshal(r)
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

	resp := &SearchResponse{
		StatusCode: osResp.StatusCode,
		Header:     osResp.Header,
	}

	if err := json.Unmarshal(respBuf.Bytes(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// SearchResponse wraps the functionality of [opensearchapi.Response] by supporting request parsing.
type SearchResponse struct {
	StatusCode int
	Header     http.Header
	Took       int       `json:"took"`
	TimedOut   bool      `json:"timed_out"`
	Shards     ShardMeta `json:"_shards,omitempty"`
	Hits       Hits      `json:"hits"`
	Error      *Error    `json:"error,omitempty"`
}

// Hits represent the results of the [search.Query] performed by the SearchRequest.
type Hits struct {
	Total    Total   `json:"total,omitempty"`
	MaxScore float64 `json:"max_score,omitempty"`
	Hits     []Hit   `json:"hits"`
}

// Total contains the total number of documents found by the [search.Query] performed by the SearchRequest.
type Total struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

// Hit the individual document found by the `[search.Query] performed by the SearchRequest.
type Hit struct {
	Index  string          `json:"_index"`
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

// GetSource returns the raw bytes of the document of the MGetResult.
func (h Hit) GetSource() []byte {
	return []byte(h.Source)
}
