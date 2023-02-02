package search

// PrefixQuery finds documents that contain the value as a prefix.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/#prefix
type PrefixQuery struct {
	field string
	value any
}

// NewPrefixQuery initializes a PrefixQuery targeting field looking for the prefix of value.
func NewPrefixQuery(name string, value any) *PrefixQuery {
	return &PrefixQuery{field: name, value: value}
}

// Source converts the PrefixQuery Source converts the MatchPhraseQuery to the correct OpenSearch JSON.
func (q *PrefixQuery) Source() (any, error) {
	source := make(map[string]any)
	pq := make(map[string]any)
	source["prefix"] = pq

	pq[q.field] = q.value

	return source, nil
}
