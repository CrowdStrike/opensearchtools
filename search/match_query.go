package search

import "encoding/json"

// MatchQuery finds documents that matches the analyzed string value.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/full-text/#match
type MatchQuery struct {
	field    string
	value    string
	operator string
}

// NewMatchQuery initializes a MatchQuery targeting field and trying to match value.
func NewMatchQuery(field, value string) *MatchQuery {
	return &MatchQuery{
		field:    field,
		value:    value,
		operator: "or",
	}
}

// SetOperator sets the operator to use when using a boolean query.
// Can be "AND" or "OR" (default).
func (q *MatchQuery) SetOperator(op string) *MatchQuery {
	q.operator = op
	return q
}

// ToOpenSearchJSON converts the MatchQuery to the correct OpenSearch JSON.
func (q *MatchQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"match": map[string]any{
			q.field: map[string]any{
				"query":    q.value,
				"operator": q.operator,
			},
		},
	}

	return json.Marshal(source)
}
