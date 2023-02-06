package search

const (
	minimumShouldMatch = "minimumShouldMatch"
)

// BoolQuery combines other queries to form more complex statements.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/bool/
type BoolQuery struct {
	must               []Query
	mustNot            []Query
	should             []Query
	minimumShouldMatch *int
	filter             []Query
}

// NewBoolQuery instantiates a boolean query.
// An empty BoolQuery will perform as if no query was provided, and all documents will match.
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

// Source coverts the BoolQuery to the correct OpenSearch JSON
func (q *BoolQuery) Source() (any, error) {
	bq := make(map[string]any)

	if q.minimumShouldMatch != nil {
		bq[minimumShouldMatch] = q.minimumShouldMatch
	}

	must, mErr := sourceQueries(q.must)
	if mErr != nil {
		return nil, mErr
	}

	if must != nil {
		bq["must"] = must
	}

	mustNot, mnErr := sourceQueries(q.mustNot)
	if mnErr != nil {
		return nil, mnErr
	}

	if mustNot != nil {
		bq["must_not"] = mustNot
	}

	should, sErr := sourceQueries(q.should)
	if sErr != nil {
		return nil, sErr
	}

	if should != nil {
		bq["should"] = should
	}

	filter, fErr := sourceQueries(q.filter)
	if fErr != nil {
		return nil, fErr
	}

	if filter != nil {
		bq["filter"] = filter
	}

	source := make(map[string]any)
	source["bool"] = bq

	return source, nil
}

// sourceQueries is a utility method to convert all sub queries to their OpenSearch source.
func sourceQueries(queries []Query) (any, error) {
	if len(queries) == 0 {
		return nil, nil
	}

	var sources []any
	for _, q := range queries {
		qSource, sErr := q.Source()
		if sErr != nil {
			return nil, sErr
		}

		sources = append(sources, qSource)
	}

	return sources, nil
}
