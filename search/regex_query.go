package search

import "encoding/json"

// RegexQuery allows you to search on a targeted field matching on values that fit the regular expression.
// An empty Regex will be rejected by OpenSearch for two reasons:
//
//   - a field must not be empty or null
//   - a value must be non-null
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

// ToOpenSearchJSON converts the RegexQuery to the correct OpenSearch JSON.
func (q *RegexQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"regexp": map[string]any{
			q.field: q.regex,
		},
	}

	return json.Marshal(source)
}
