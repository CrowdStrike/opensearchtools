package search

import "encoding/json"

// MatchAllQuery returns all documents.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/full-text/#match-all
type MatchAllQuery struct {
}

// NewMatchAllQuery instantiates a MatchAllQuery.
func NewMatchAllQuery() *MatchAllQuery {
	return &MatchAllQuery{}
}

// ToOpenSearchJSON converts the MatchAllQuery to the correct OpenSearch JSON.
func (q *MatchAllQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"match_all": struct{}{},
	}

	return json.Marshal(source)
}
