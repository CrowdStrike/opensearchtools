package search

// RegexQuery allows you to search on a targeted field matching on values that fit the regular expression.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/#regex
type RegexQuery struct {
	field string
	regex string
}

// NewRegexQuery instantiates a RegexQuery targeting field with pattern regex.
func NewRegexQuery(field, regex string) *RegexQuery {
	return &RegexQuery{
		field: field,
		regex: regex,
	}
}

// Source converts the RegexQuery to the correct OpenSearch JSON.
func (q *RegexQuery) Source() (any, error) {
	rq := make(map[string]any)
	rq[q.field] = q.regex

	source := make(map[string]any)
	source["regexp"] = rq

	return source, nil
}
