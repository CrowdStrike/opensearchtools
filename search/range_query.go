package search

// RangeQuery allows you to search on a targeted field matching a defined range.
//
// For more details see https://opensearch.org/docs/latest/opensearch/query-dsl/term/#range-query
type RangeQuery struct {
	field string
	gt    any
	gte   any
	lt    any
	lte   any
}

// NewRangeQuery instantiates a Range Query targeting field.
// An empty range query will match all documents that contain the field.
func NewRangeQuery(field string) *RangeQuery {
	return &RangeQuery{field: field}
}

// Gt sets the greater than value.
func (q *RangeQuery) Gt(value any) *RangeQuery {
	q.gt = value
	return q
}

// Gte sets the greater than or equal to value.
func (q *RangeQuery) Gte(value any) *RangeQuery {
	q.gte = value
	return q
}

// Lt sets the less than value.
func (q *RangeQuery) Lt(value any) *RangeQuery {
	q.lt = value
	return q
}

// Lte sets the less than or equal to value.
func (q *RangeQuery) Lte(value any) *RangeQuery {
	q.lte = value
	return q
}

// Source converts the RangeQuery to the correct OpenSearch JSON.
func (q *RangeQuery) Source() (any, error) {
	ranges := make(map[string]any)
	if q.gt != nil {
		ranges["gt"] = q.gt
	}

	if q.gte != nil {
		ranges["gte"] = q.gte
	}

	if q.lt != nil {
		ranges["lt"] = q.lt
	}

	if q.lte != nil {
		ranges["lte"] = q.lte
	}

	rq := make(map[string]any)
	rq[q.field] = ranges

	source := make(map[string]any)
	source["range"] = rq

	return source, nil
}
