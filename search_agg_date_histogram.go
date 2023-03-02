package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// DateHistogramAggregation buckets documents based on a date interval.
// An empty DateHistogramAggregation will have several issues with execution:
//   - the target Field must be non-null and non-empty
//   - the Interval must be non-null and non-empty
//
// For more details see https://opensearch.org/docs/latest/opensearch/bucket-agg/#histogram-date_histogram
type DateHistogramAggregation struct {
	// Field to be bucketed
	Field string

	// MinDocCount is the lower count threshold for a bucket to be included in the results.
	// Negative counts will be omitted
	MinDocCount int

	// Interval string using OpenSearch [date math].
	// [date math]: https://opensearch.org/docs/latest/opensearch/supported-field-types/date/#date-math
	Interval string

	// TimeZone, times are stored internally in UTC and by default date histograms are bucketed in UTC.
	// Set the TimeZone to overwrite this default
	TimeZone string

	// Aggregations sub aggregations for each bucket. Mapped by string label to sub aggregation
	Aggregations map[string]Aggregation
}

// NewDateHistogramAggregation instantiates a DateHistogramAggregation targeting
// the provided field with the provided interval. Sets the MinDocCount to -1 to be
// omitted in favor of the OpenSearch default.
func NewDateHistogramAggregation(field, interval string) *DateHistogramAggregation {
	return &DateHistogramAggregation{
		Field:        field,
		MinDocCount:  -1,
		Interval:     interval,
		Aggregations: make(map[string]Aggregation),
	}
}

// WithMinDocCount the lower count threshold for a bucket to be included in the results
func (d *DateHistogramAggregation) WithMinDocCount(minCount int) *DateHistogramAggregation {
	d.MinDocCount = minCount
	return d
}

// WithTimeZone overwriting the default UTC timezone
func (d *DateHistogramAggregation) WithTimeZone(tz string) *DateHistogramAggregation {
	d.TimeZone = tz
	return d
}

// AddSubAggregation to the TermsAggregation with the provided name
func (d *DateHistogramAggregation) AddSubAggregation(name string, agg Aggregation) BucketAggregation {
	if d.Aggregations == nil {
		d.Aggregations = map[string]Aggregation{name: agg}
	} else {
		d.Aggregations[name] = agg
	}

	return d
}

// ConvertSubAggregations uses the provided converter to convert all the sub aggregations in this TermsAggregation
func (d *DateHistogramAggregation) ConvertSubAggregations(converter AggregateVersionConverter) (map[string]Aggregation, error) {
	convertedAggs := make(map[string]Aggregation, len(d.Aggregations))

	for name, agg := range d.Aggregations {
		cAgg, cErr := converter(agg)
		if cErr != nil {
			return nil, cErr
		}

		convertedAggs[name] = cAgg
	}

	return convertedAggs, nil
}

// ToOpenSearchJSON converts the TermsAggregation to the correct OpenSearch JSON.
func (d *DateHistogramAggregation) ToOpenSearchJSON() ([]byte, error) {
	if d.Field == "" {
		return nil, fmt.Errorf("a DateHistogramAggregation requires a target field")
	}

	if d.Interval == "" {
		return nil, fmt.Errorf("a DateHistogramAggregation requires a interval")
	}

	da := map[string]any{
		"field":    d.Field,
		"interval": d.Interval,
	}

	if d.MinDocCount >= 0 {
		da["min_doc_count"] = d.MinDocCount
	}

	if d.TimeZone != "" {
		da["time_zone"] = d.TimeZone
	}

	source := map[string]any{
		"date_histogram": da,
	}

	if len(d.Aggregations) > 0 {
		subAggs := make(map[string]json.RawMessage)
		for aggName, agg := range d.Aggregations {
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

// DateHistogramAggregationResults represents the results from a DateHistogramAggregation request
type DateHistogramAggregationResults struct {
	Buckets []DateHistogramBucketResult
}

// UnmarshalJSON implements [json.Unmarshaler] to decode a json byte slice into a DateHistogramAggregationResults
// Errors on unknown fields.
func (d *DateHistogramAggregationResults) UnmarshalJSON(m []byte) error {
	// map[key] -> value
	var rawResp map[string]json.RawMessage
	if err := json.Unmarshal(m, &rawResp); err != nil {
		return err
	}

	if d == nil {
		return fmt.Errorf("invalid DateHistogramAggregationResults target, nil")
	}

	for key, value := range rawResp {
		switch key {
		case "buckets":
			if err := json.Unmarshal(value, &d.Buckets); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown DateHistogramAggregationResults field %s", key)
		}
	}

	return nil
}

// DateHistogramBucketResult is a [AggregationResultMap] for a DateHistogramAggregation
type DateHistogramBucketResult struct {
	KeyString             string
	Key                   uint64
	DocCount              uint64
	SubAggregationResults map[string]json.RawMessage
}

// UnmarshalJSON implements [json.Unmarshaler] to decode a json byte slice into a DateHistogramBucketResult
func (d *DateHistogramBucketResult) UnmarshalJSON(m []byte) error {
	// map[key] -> value
	var rawResp map[string]json.RawMessage
	if err := json.Unmarshal(m, &rawResp); err != nil {
		return err
	}

	if d == nil {
		return fmt.Errorf("invalid DateHistogramBucketResult target, nil")
	}

	d.SubAggregationResults = make(map[string]json.RawMessage)
	for key, value := range rawResp {
		switch key {
		case "key_as_string":
			if err := json.Unmarshal(value, &d.KeyString); err != nil {
				return err
			}
		case "key":
			if err := json.Unmarshal(value, &d.Key); err != nil {
				return err
			}
		case "doc_count":
			if err := json.Unmarshal(value, &d.DocCount); err != nil {
				return err
			}
		default:
			d.SubAggregationResults[key] = value
		}
	}

	return nil
}

// GetAggregationResultSource implements [opensearchtools.AggregationResultSet] to fetch a sub aggregation result and
// return the raw JSON source for the provided name.
func (d *DateHistogramBucketResult) GetAggregationResultSource(name string) ([]byte, bool) {
	if len(d.SubAggregationResults) == 0 {
		return nil, false
	}

	subAggSource, exists := d.SubAggregationResults[name]
	return subAggSource, exists
}
