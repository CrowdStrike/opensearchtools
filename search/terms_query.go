package search

// TermsQuery finds documents that have the field match one of the listed values.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/#terms
type TermsQuery struct {
	field  string
	values []any
}

// NewTermsQuery instantiates a TermsQuery targeting field looking for one of the values.
func NewTermsQuery(field string, values ...any) *TermsQuery {
	return &TermsQuery{
		field:  field,
		values: values,
	}
}

// Source converts the TermsQuery to the correct OpenSearch JSON.
func (q *TermsQuery) Source() (any, error) {
	tq := make(map[string]any)
	tq[q.field] = q.values

	source := make(map[string]any)
	source["terms"] = tq

	return source, nil
}
