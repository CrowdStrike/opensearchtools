package search

import "encoding/json"

// ExistsQuery searches for documents that contain a specific field.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/#exists
type ExistsQuery struct {
	field string
}

// NewExistsQuery instantiates an exists query.
// An empty field value will be rejected by OpenSearch
func NewExistsQuery(field string) *ExistsQuery {
	return &ExistsQuery{field: field}
}

// ToOpenSearchJSON converts the ExistsQuery to the correct OpenSearch JSON.
func (q *ExistsQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"exists": map[string]any{
			"field": q.field,
		},
	}

	return json.Marshal(source)
}
