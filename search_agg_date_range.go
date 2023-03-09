package opensearchtools

import (
	"encoding/json"
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
// Implements [BucketAggregation.AddSubAggregation]
func (dr *DateRangeAggregation) AddSubAggregation(name string, agg Aggregation) BucketAggregation {
	if dr.Aggregations == nil {
		dr.Aggregations = map[string]Aggregation{name: agg}
	} else {
		dr.Aggregations[name] = agg
	}

	return dr
}

// SubAggregations returns all aggregations added to the bucket aggregation.
// Implements [BucketAggregation.SubAggregations]
func (dr *DateRangeAggregation) SubAggregations() map[string]Aggregation {
	return dr.Aggregations
}

// Validate that the aggregation is executable.
// Implements [Aggregation.Validate].
func (dr *DateRangeAggregation) Validate() ValidationResults {
	vrs := NewValidationResults()

	if dr.Field == "" {
		vrs.Add(NewValidationResult("a DateRangeAggregation requires a target field", true))
	}

	if len(dr.Ranges) == 0 {
		vrs.Add(NewValidationResult("a DateRangeAggregation requires at least one range bucket", true))
	}

	for _, subAgg := range dr.Aggregations {
		vrs.Extend(subAgg.Validate())
	}

	return vrs
}

// ToOpenSearchJSON converts the DateRangeAggregation to the correct OpenSearch JSON.
// Implements [Aggregation.ToOpenSearchJSON].
func (dr *DateRangeAggregation) ToOpenSearchJSON() ([]byte, error) {
	if vrs := dr.Validate(); vrs.IsFatal() {
		return nil, NewValidationError(vrs)
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
