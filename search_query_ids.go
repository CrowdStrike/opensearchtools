package opensearchtools

import "encoding/json"

// IDsQuery finds documents that have the id match one of the listed values.
// An empty IDsQuery will be rejected by OpenSearch for the following reasons:
//
//   - a value must be non-null
//
// For more details see https://opensearch.org/docs/latest/query-dsl/term/#ids
type IDsQuery struct {
	values []any
}

// NewIDsQuery instantiates a IDsQuery looking for one of the id values.
func NewIDsQuery(values ...any) *IDsQuery {
	return &IDsQuery{
		values: values,
	}
}

// ToOpenSearchJSON converts the IDsQuery to the correct OpenSearch JSON.
func (q *IDsQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"ids": map[string]any{
			"values": q.values,
		},
	}

	return json.Marshal(source)
}
