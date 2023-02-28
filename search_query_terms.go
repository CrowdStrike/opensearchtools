package opensearchtools

import "encoding/json"

// TermsQuery finds documents that have the field match one of the listed values.
// An empty TermsQuery will be rejected by OpenSearch for two reasons:
//
//   - a field must not be empty or null
//   - a value must be non-null
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

// ToOpenSearchJSON converts the TermsQuery to the correct OpenSearch JSON.
func (q *TermsQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"terms": map[string]any{
			q.field: q.values,
		},
	}

	return json.Marshal(source)
}
