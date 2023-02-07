package search

import "encoding/json"

// WildcardQuery searches for documents with a field matching a wildcard pattern.
// An empty TermQuery will be rejected by OpenSearch for two reasons:
//
//   - a field must not be empty or null
//   - a value must be non-null
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/#wildcards
type WildcardQuery struct {
	field string
	value string
}

// NewWildcardQuery instantiates a wildcard query targeting field looking for a wildcard match on value.
func NewWildcardQuery(field, value string) *WildcardQuery {
	return &WildcardQuery{field: field, value: value}
}

// ToOpenSearchJSON converts the WildcardQuery to the correct OpenSearch JSON.
func (q *WildcardQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"wildcard": map[string]any{
			q.field: q.value,
		},
	}

	return json.Marshal(source)
}
