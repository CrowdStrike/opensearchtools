package search

// MatchPhraseQuery finds document that match documents that contain an exact phrase in a specified order.
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

// Source converts the MatchPhraseQuery to the correct OpenSearch JSON.
func (q *MatchPhraseQuery) Source() (any, error) {
	mq := make(map[string]any)
	mq[q.field] = q.phrase

	source := make(map[string]any)
	source["match_phrase"] = mq

	return source, nil
}
