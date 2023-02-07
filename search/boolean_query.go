package search

import "encoding/json"

const (
	minimumShouldMatch = "minimum_should_match"
)

// BoolQuery combines other queries to form more complex statements.
// An empty BoolQuery is executable by OpenSearch. Without any constraints it will match on everything.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/bool/
type BoolQuery struct {
	must               []Query
	mustNot            []Query
	should             []Query
	minimumShouldMatch *int
	filter             []Query
}

// NewBoolQuery instantiates an empty boolean query.
func NewBoolQuery() *BoolQuery {
	return &BoolQuery{}
}

// MinimumShouldMatch Optional parameter for use with a should query clause.
// Specifies the minimum number of queries that the document must match for it to be returned in the results.
func (q *BoolQuery) MinimumShouldMatch(n int) *BoolQuery {
	q.minimumShouldMatch = &n
	return q
}

// Must Logical AND operator. The results must match the queries in this clause.
// If you have multiple queries, all of them must match.
func (q *BoolQuery) Must(queries ...Query) *BoolQuery {
	q.must = append(q.must, queries...)
	return q
}

// MustNot Logical NOT operator. All matches are excluded from the results.
func (q *BoolQuery) MustNot(queries ...Query) *BoolQuery {
	q.mustNot = append(q.mustNot, queries...)
	return q
}

// Should Logical OR operator. The results must match at least one of the queries, but, optionally,
// they can match more than one query. Each matching should clause increases the relevancy score.
// You can set the minimum number of queries that must match using BoolQuery.MinimumShouldMatch
func (q *BoolQuery) Should(queries ...Query) *BoolQuery {
	q.should = append(q.should, queries...)
	return q
}

// Filter Logical AND operator that is applied first to reduce your dataset before applying the queries.
// A query within a filter clause is a yes or no option. If a document matches the query, it is returned in the results;
// otherwise, it is not. The results of a filter query are generally cached to allow for a faster return.
// Use the filter query to filter the results based on exact matches, ranges, dates, numbers, and so on.
func (q *BoolQuery) Filter(queries ...Query) *BoolQuery {
	q.filter = append(q.filter, queries...)
	return q
}

// ToOpenSearchJSON coverts the BoolQuery to the correct OpenSearch JSON
func (q *BoolQuery) ToOpenSearchJSON() ([]byte, error) {
	bq := make(map[string]any)

	if q.minimumShouldMatch != nil {
		bq[minimumShouldMatch] = q.minimumShouldMatch
	}

	must, mErr := convertSubQueries(q.must)
	if mErr != nil {
		return nil, mErr
	}

	if must != nil {
		bq["must"] = must
	}

	mustNot, mnErr := convertSubQueries(q.mustNot)
	if mnErr != nil {
		return nil, mnErr
	}

	if mustNot != nil {
		bq["must_not"] = mustNot
	}

	should, sErr := convertSubQueries(q.should)
	if sErr != nil {
		return nil, sErr
	}

	if should != nil {
		bq["should"] = should
	}

	filter, fErr := convertSubQueries(q.filter)
	if fErr != nil {
		return nil, fErr
	}

	if filter != nil {
		bq["filter"] = filter
	}

	source := map[string]any{
		"bool": bq,
	}

	return json.Marshal(source)
}

// convertSubQueries is a utility method to convert all sub queries to their OpenSearch source.
func convertSubQueries(queries []Query) (json.RawMessage, error) {
	if len(queries) == 0 {
		return nil, nil
	}

	var jsonQueries []json.RawMessage
	for _, q := range queries {
		qJSON, jErr := q.ToOpenSearchJSON()
		if jErr != nil {
			return nil, jErr
		}

		jsonQueries = append(jsonQueries, qJSON)
	}

	return json.Marshal(jsonQueries)
}
