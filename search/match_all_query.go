package search

// MatchAllQuery returns all documents.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/full-text/#match-all
type MatchAllQuery struct {
}

// NewMatchAllQuery instantiates a MatchAllQuery.
func NewMatchAllQuery() *MatchAllQuery {
	return &MatchAllQuery{}
}

// Source converts the MatchAllQuery to the correct OpenSearch JSON.
func (q *MatchAllQuery) Source() (any, error) {
	source := make(map[string]any)
	source["match_all"] = struct{}{}

	return source, nil
}
