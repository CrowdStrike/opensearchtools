package search

// WildcardQuery searches for documents with a field matching a wildcard pattern.
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

// Source converts the WildcardQuery to the correct OpenSearch JSON.
func (q *WildcardQuery) Source() (any, error) {
	eq := make(map[string]any)
	eq[q.field] = q.value

	source := make(map[string]any)
	source["wildcard"] = eq

	return source, nil
}
