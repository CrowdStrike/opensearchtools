package search

import "encoding/json"

// TermQuery finds documents that have the field matching the exact value.
// An empty TermQuery will be rejected by OpenSearch for two reasons:
//
//   - a field must not be empty or null
//   - a value must be non-null
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/
type TermQuery struct {
	//TODO: given the above empty constraints, should we validate on the client library?
	field string
	value any
}

// NewTermQuery initializes a TermQuery targeting field looking for the exact value.
func NewTermQuery(field string, value any) *TermQuery {
	return &TermQuery{field: field, value: value}
}

// ToOpenSearchJSON converts the TermQuery to the correct OpenSearch JSON.
func (q *TermQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"term": map[string]any{
			q.field: q.value,
		},
	}

	return json.Marshal(source)
}
