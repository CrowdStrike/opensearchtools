package opensearchtools

import "encoding/json"

// RangeQuery allows you to search on a targeted field matching a defined range.
// An empty RangeQuery will be rejected by OpenSearch as it requires a non-null and non-empty field.
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
// A RangeQuery with no range operations will function like an
// ExistsQuery and match all documents that contain the field.
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

// ToOpenSearchJSON converts the RangeQuery to the correct OpenSearch JSON.
func (q *RangeQuery) ToOpenSearchJSON() ([]byte, error) {
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

	source := map[string]any{
		"range": map[string]any{
			q.field: ranges,
		},
	}

	return json.Marshal(source)
}
