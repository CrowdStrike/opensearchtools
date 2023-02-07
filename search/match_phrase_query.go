package search

import "encoding/json"

// MatchPhraseQuery finds document that match documents that contain an exact phrase in a specified order.
// An empty MatchPhraseQuery will be rejected by OpenSearch for two reasons:
//
//   - a field must not be empty or null
//   - a value must be non-null
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/full-text/#match-phrase
type MatchPhraseQuery struct {
	field  string
	phrase string
}

// NewMatchPhraseQuery instantiates a MatchPhraseQuery targeting field and looking for phrase.
func NewMatchPhraseQuery(field, phrase string) *MatchPhraseQuery {
	return &MatchPhraseQuery{
		field:  field,
		phrase: phrase,
	}
}

// ToOpenSearchJSON converts the MatchPhraseQuery to the correct OpenSearch JSON.
func (q *MatchPhraseQuery) ToOpenSearchJSON() ([]byte, error) {
	source := map[string]any{
		"match_phrase": map[string]any{
			q.field: q.phrase,
		},
	}

	return json.Marshal(source)
}
