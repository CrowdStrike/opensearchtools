package search

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

// Operator sets the operator to use when using a boolean query.
// Can be "AND" or "OR" (default).
func (q *MatchQuery) Operator(op string) *MatchQuery {
	q.operator = op
	return q
}

// Source converts the MatchQuery to the correct OpenSearch JSON.
func (q *MatchQuery) Source() (any, error) {
	mq := make(map[string]any)
	mq["query"] = q.value
	mq["operator"] = q.operator

	search := make(map[string]any)
	search[q.field] = mq

	source := make(map[string]any)
	source["match"] = search

	return source, nil
}
