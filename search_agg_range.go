package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// RangeAggregation lets you manually define each bucket, and it's range.
// An empty RangeAggregation will have some issues with execution:
//   - the target Field must be non-nil and non-empty
//   - the aggregation must have at least one Range defined
//
// For more details see https://opensearch.org/docs/latest/opensearch/bucket-agg/#range-date_range-ip_range
type RangeAggregation struct {
	// Field to be targeted
	Field string

	// Ranges - list of range buckets
	Ranges []Range

	// Aggregations sub aggregations for each bucket. Mapped by string label to sub aggregation
	Aggregations map[string]Aggregation
}

type Range struct {
	// Key label for the range bucket. If one is not provided
	// OpenSearch defaults to "{From}-{To}
	Key string `json:"key,omitempty"`

	// From - lower inclusive bound of the range
	From any `json:"from,omitempty"`

	// To - upper exclusive bound of the range
	To any `json:"to,omitempty"`
}

// NewRangeAggregation instantiates a RangeAggregation targeting the provided field.
func NewRangeAggregation(field string) *RangeAggregation {
	return &RangeAggregation{
		Field:        field,
		Aggregations: make(map[string]Aggregation),
	}
}

// AddRange adds an un-keyed range to the bucket list
func (r *RangeAggregation) AddRange(from, to any) *RangeAggregation {
	r.Ranges = append(r.Ranges, Range{
		From: from,
		To:   to,
	})

	return r
}

// AddKeyedRange adds a keyed range to the bucket list
func (r *RangeAggregation) AddKeyedRange(key string, from, to any) *RangeAggregation {
	r.Ranges = append(r.Ranges, Range{
		Key:  key,
		From: from,
		To:   to,
	})

	return r
}

// AddRanges adds any number of Ranges to the bucket list
func (r *RangeAggregation) AddRanges(ranges ...Range) *RangeAggregation {
	r.Ranges = append(r.Ranges, ranges...)
	return r
}

// AddSubAggregation to the RangeAggregation with the provided name
// Implements [BucketAggregation.AddSubAggregation]
func (r *RangeAggregation) AddSubAggregation(name string, agg Aggregation) BucketAggregation {
	if r.Aggregations == nil {
		r.Aggregations = map[string]Aggregation{name: agg}
	} else {
		r.Aggregations[name] = agg
	}

	return r
}

// SubAggregations returns all aggregations added to the bucket aggregation.
// Implements [BucketAggregation.SubAggregations]
func (r *RangeAggregation) SubAggregations() map[string]Aggregation {
	return r.Aggregations
}

// ToOpenSearchJSON converts the RangeAggregation to the correct OpenSearch JSON.
func (r *RangeAggregation) ToOpenSearchJSON() ([]byte, error) {
	if r.Field == "" {
		return nil, fmt.Errorf("a RangeAggregation requires a target field")
	}

	if len(r.Ranges) == 0 {
		return nil, fmt.Errorf("a RangeAggregation requires at least one range bucket")
	}

	source := map[string]any{
		"range": map[string]any{
			"field":  r.Field,
			"ranges": r.Ranges,
		},
	}

	if len(r.Aggregations) > 0 {
		subAggs := make(map[string]json.RawMessage)
		for aggName, agg := range r.Aggregations {
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

// RangeAggregationResults represents the results from a range aggregation request.
type RangeAggregationResults struct {
	Buckets []RangeBucketResult `json:"buckets"`
}

// RangeBucketResult is a [AggregationResultMap] for a RangeAggregation
type RangeBucketResult struct {
	// Key - the bucket label
	Key string

	// From - the lower inclusive bound of the bucket.
	// It will match the type of the [Range.From] in the request.
	From any

	// FromString - if the [Range.From] was not a string,
	// a string representation will also be included
	FromString string

	// To - the upper exclusive bound of the bucket.
	// It will match the type of the [Range.To] in the request.
	To any

	// ToString - if the [Range.To] was not a string,
	// a string representation will also be included
	ToString string

	// DocCount - number of documents that fit in this bucket
	DocCount int64

	// SubAggregationResults for any nested aggregations
	SubAggregationResults map[string]json.RawMessage
}

// UnmarshalJSON implements [json.Unmarshaler] to decode a json byte slice into a RangeBucketResult
func (r *RangeBucketResult) UnmarshalJSON(m []byte) error {
	// map[key] -> value
	var rawResp map[string]json.RawMessage
	if err := json.Unmarshal(m, &rawResp); err != nil {
		return err
	}

	if r == nil {
		return fmt.Errorf("invalid TermBucketResult target, nil")
	}

	r.SubAggregationResults = make(map[string]json.RawMessage)
	for key, value := range rawResp {
		switch key {
		case "key":
			if err := json.Unmarshal(value, &r.Key); err != nil {
				return err
			}
		case "doc_count":
			if err := json.Unmarshal(value, &r.DocCount); err != nil {
				return err
			}
		case "from":
			if err := json.Unmarshal(value, &r.From); err != nil {
				return err
			}
		case "from_as_string":
			if err := json.Unmarshal(value, &r.FromString); err != nil {
				return err
			}
		case "to":
			if err := json.Unmarshal(value, &r.To); err != nil {
				return err
			}
		case "to_as_string":
			if err := json.Unmarshal(value, &r.ToString); err != nil {
				return err
			}
		default:
			// any number of sub aggregation results
			r.SubAggregationResults[key] = value
		}
	}

	return nil
}

// GetAggregationResultSource implements [opensearchtools.AggregationResultSet] to fetch a sub aggregation result and
// return the raw JSON source for the provided name.
func (r *RangeBucketResult) GetAggregationResultSource(name string) ([]byte, bool) {
	if len(r.SubAggregationResults) == 0 {
		return nil, false
	}

	subAggSource, exists := r.SubAggregationResults[name]
	return subAggSource, exists
}

// Keys implemented for [opensearchtools.AggregationResultSet] to return the list of aggregation result keys
func (r *RangeBucketResult) Keys() []string {
	keys := make([]string, len(r.SubAggregationResults))

	i := 0
	for k := range r.SubAggregationResults {
		keys[i] = k
		i++
	}

	return keys
}
