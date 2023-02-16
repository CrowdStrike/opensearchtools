package search

import "encoding/json"

// PrefixQuery finds documents that contain the value as a prefix.
// An empty PrefixQuery will be rejected by OpenSearch for two reasons:
//
//   - a field must not be empty or null
//   - a value must be non-null
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/#prefix
type PrefixQuery struct {
	field string
	value any
}

// NewPrefixQuery initializes a PrefixQuery targeting field looking for the prefix of value.
func NewPrefixQuery(field string, value any) *PrefixQuery {
	return &PrefixQuery{field: field, value: value}
}

// ToOpenSearchJSON converts the PrefixQuery Source converts the MatchPhraseQuery to the correct OpenSearch JSON.
func (q *PrefixQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"prefix": map[string]any{
			q.field: q.value,
		},
	}

	return json.Marshal(source)
}
