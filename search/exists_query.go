package search

// ExistsQuery searches for documents that contain a specific field.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/#exists
type ExistsQuery struct {
	field string
}

// NewExistsQuery instantiates an exists query.
func NewExistsQuery(field string) *ExistsQuery {
	return &ExistsQuery{field: field}
}

// Source converts the ExistsQuery to the correct OpenSearch JSON.
func (q *ExistsQuery) Source() (any, error) {
	eq := make(map[string]any)
	eq["field"] = q.field

	source := make(map[string]any)
	source["exists"] = eq

	return source, nil
}
