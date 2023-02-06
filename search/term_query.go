package search

// TermQuery finds documents that have the field matching the exact value.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/
type TermQuery struct {
	field string
	value any
}

// NewTermQuery initializes a TermQuery targeting field looking for the exact value.
func NewTermQuery(field string, value any) *TermQuery {
	return &TermQuery{field: field, value: value}
}

// Source converts the TermQuery to the correct OpenSearch JSON.
func (q *TermQuery) Source() (any, error) {
	source := make(map[string]any)
	tq := make(map[string]any)
	source["term"] = tq

	tq[q.field] = q.value

	return source, nil
}
