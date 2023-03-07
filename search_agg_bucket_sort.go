package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// BucketSortAggregation is a [pipeline aggregation] that modifies its parent aggregation
// by sorting and truncating the results. It does not generate a response itself.
// An empty BucketSortAggregation will return no results due to the default size of 0.
//
// For more details see https://opensearch.org/docs/latest/opensearch/pipeline-agg/#bucket_sort
// [pipeline aggregation]: https://opensearch.org/docs/latest/opensearch/pipeline-agg/
type BucketSortAggregation struct {
	// From - the number of buckets to be truncated before returning
	// Negative values will be omitted
	From int

	// Size - the number of buckets to be returned
	// Negative values will be omitted
	Size int

	// Sort for the parent aggregation
	Sort []sort
}

// sort is handled differently than other [Aggregation] [Order]ing
type sort struct {
	field string
	desc  bool
}

// NewBucketSortAggregation instantiates a BucketSortAggregation with From and Size
// set to -1 to be omitted in favor of OpenSearch defaults
func NewBucketSortAggregation() *BucketSortAggregation {
	return &BucketSortAggregation{
		From: -1,
		Size: -1,
	}
}

// WithFrom for the number of buckets to be truncated before returning
func (b *BucketSortAggregation) WithFrom(from int) *BucketSortAggregation {
	b.From = from
	return b
}

// WithSize for the number of buckets to be returned
func (b *BucketSortAggregation) WithSize(size int) *BucketSortAggregation {
	b.Size = size
	return b
}

// AddSort for the targeted field
func (b *BucketSortAggregation) AddSort(field string, isDesc bool) *BucketSortAggregation {
	b.Sort = append(b.Sort, sort{
		field: field,
		desc:  isDesc,
	})
	return b
}

func (b *BucketSortAggregation) ToOpenSearchJSON() ([]byte, error) {
	ba := make(map[string]any)
	if b.From >= 0 {
		ba["from"] = b.From
	}

	if b.Size >= 0 {
		ba["size"] = b.Size
	}

	if len(b.Sort) > 0 {
		sortSource := make([]any, len(b.Sort))

		for i, s := range b.Sort {
			if s.field == "" {
				return nil, fmt.Errorf("sort missing target field")
			}

			order := "asc"
			if s.desc {
				order = "desc"
			}

			sortSource[i] = map[string]any{
				s.field: map[string]any{
					"order": order,
				},
			}
		}

		ba["sort"] = sortSource
	}

	source := map[string]any{
		"bucket_sort": ba,
	}

	return json.Marshal(source)
}
