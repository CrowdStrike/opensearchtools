package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// DateRangeAggregation is conceptually the same as the [RangeAggregation],
// except that it lets you perform date math and format the date strings.
// An empty DateDateRangeAggregation will have some issues with execution
type DateRangeAggregation struct {
	// Field to be targeted
	Field string

	// Ranges - list of range buckets
	Ranges []Range

	// Format - the date format for the [RangeBucketResult.FromString] and
	// [RangeBucketResult.ToString] in the results
	Format string

	// Aggregations sub aggregations for each bucket. Mapped by string label to sub aggregation
	Aggregations map[string]Aggregation
}

// NewDateRangeAggregation instantiates a DateRangeAggregation targeting the provided field.
func NewDateRangeAggregation(field string) *DateRangeAggregation {
	return &DateRangeAggregation{
		Field:        field,
		Aggregations: make(map[string]Aggregation),
	}
}

// AddRange adds an un-keyed range to the bucket list
func (dr *DateRangeAggregation) AddRange(from, to any) *DateRangeAggregation {
	dr.Ranges = append(dr.Ranges, Range{
		From: from,
		To:   to,
	})

	return dr
}

// AddKeyedRange adds a keyed range to the bucket list
func (dr *DateRangeAggregation) AddKeyedRange(key string, from, to any) *DateRangeAggregation {
	dr.Ranges = append(dr.Ranges, Range{
		Key:  key,
		From: from,
		To:   to,
	})

	return dr
}

// AddRanges adds any number of Ranges to the bucket list
func (dr *DateRangeAggregation) AddRanges(ranges ...Range) *DateRangeAggregation {
	dr.Ranges = append(dr.Ranges, ranges...)
	return dr
}

// WithFormat for the date from and to response
func (dr *DateRangeAggregation) WithFormat(format string) *DateRangeAggregation {
	dr.Format = format
	return dr
}

// AddSubAggregation to the DateRangeAggregation with the provided name
func (dr *DateRangeAggregation) AddSubAggregation(name string, agg Aggregation) BucketAggregation {
	if dr.Aggregations == nil {
		dr.Aggregations = map[string]Aggregation{name: agg}
	} else {
		dr.Aggregations[name] = agg
	}

	return dr
}

// ConvertSubAggregations uses the provided converter to convert all the sub aggregations in this DateRangeAggregation
func (dr *DateRangeAggregation) ConvertSubAggregations(converter AggregateVersionConverter) (map[string]Aggregation, error) {
	convertedAggs := make(map[string]Aggregation, len(dr.Aggregations))

	for name, agg := range dr.Aggregations {
		cAgg, cErr := converter(agg)
		if cErr != nil {
			return nil, cErr
		}

		convertedAggs[name] = cAgg
	}

	return convertedAggs, nil
}

// ToOpenSearchJSON converts the DateRangeAggregation to the correct OpenSearch JSON.
func (dr *DateRangeAggregation) ToOpenSearchJSON() ([]byte, error) {
	if dr.Field == "" {
		return nil, fmt.Errorf("a DateRangeAggregation requires a target field")
	}

	if len(dr.Ranges) == 0 {
		return nil, fmt.Errorf("a DateRangeAggregation requires at least one range bucket")
	}

	ra := map[string]any{
		"field":  dr.Field,
		"ranges": dr.Ranges,
	}

	if dr.Format != "" {
		ra["format"] = dr.Format
	}

	source := map[string]any{
		"date_range": ra,
	}

	if len(dr.Aggregations) > 0 {
		subAggs := make(map[string]json.RawMessage)
		for aggName, agg := range dr.Aggregations {
			aggJSON, jErr := agg.ToOpenSearchJSON()
			if jErr != nil {
				return nil, jErr
			}

			subAggs[aggName] = aggJSON
		}

		source["aggs"] = subAggs
	}

	return json.Marshal(source)
}

// DateRangeAggregationResults represents the results from a range aggregation request.
type DateRangeAggregationResults struct {
	Buckets []RangeBucketResult `json:"buckets"`
}
