package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// TermsAggregation dynamically creates a bucket for each unique term of a field.
// An empty TermsAggregation will have two issues with execution:
//
//   - the target Field must be non-null and non-empty.
//   - a Size of 0 will return no buckets
//
// For more details see https://opensearch.org/docs/latest/opensearch/bucket-agg/
type TermsAggregation struct {
	// Field to be bucketed
	Field string

	// Size of the number of buckets to be returned. Negative sizes will be omitted
	Size int

	// Order list of [Order]s to sort the aggregation buckets. Default order is _count: desc
	Order []Order

	// Aggregations sub aggregations for each bucket. Mapped by string label to sub aggregation
	Aggregations map[string]Aggregation
}

// NewTermsAggregation instantiates a TermsAggregation targeting the provided field
// and sets the Size to -1 to be omitted for the default value.
func NewTermsAggregation(field string) *TermsAggregation {
	return &TermsAggregation{
		Field:        field,
		Size:         -1,
		Aggregations: make(map[string]Aggregation),
	}
}

// WithSize for the number of buckets to be returned
func (t *TermsAggregation) WithSize(size int) *TermsAggregation {
	t.Size = size
	return t
}

// AddOrder of the returned buckets
func (t *TermsAggregation) AddOrder(orders ...Order) *TermsAggregation {
	t.Order = append(t.Order, orders...)
	return t
}

// AddSubAggregation to the TermsAggregation with the provided name
func (t *TermsAggregation) AddSubAggregation(name string, agg Aggregation) BucketAggregation {
	if t.Aggregations == nil {
		t.Aggregations = map[string]Aggregation{name: agg}
	} else {
		t.Aggregations[name] = agg
	}

	return t
}

// ConvertSubAggregations uses the provided converter to convert all the sub aggregations in this TermsAggregation
func (t *TermsAggregation) ConvertSubAggregations(converter AggregateVersionConverter) (map[string]Aggregation, error) {
	convertedAggs := make(map[string]Aggregation, len(t.Aggregations))

	for name, agg := range t.Aggregations {
		cAgg, cErr := converter(agg)
		if cErr != nil {
			return nil, cErr
		}

		convertedAggs[name] = cAgg
	}

	return convertedAggs, nil
}

// ToOpenSearchJSON converts the TermsAggregation to the correct OpenSearch JSON.
func (t *TermsAggregation) ToOpenSearchJSON() ([]byte, error) {
	if t.Field == "" {
		return nil, fmt.Errorf("a TermsAggregation requires a target field")
	}

	ta := map[string]any{
		"field": t.Field,
	}

	if t.Size >= 0 {
		ta["size"] = t.Size
	}

	if len(t.Order) > 0 {
		var rawOrder []json.RawMessage
		for _, o := range t.Order {
			source, oErr := o.ToOpenSearchJSON()
			if oErr != nil {
				return nil, oErr
			}

			rawOrder = append(rawOrder, source)
		}

		ta["order"] = rawOrder
	}

	source := map[string]any{
		"terms": ta,
	}

	if len(t.Aggregations) > 0 {
		subAggs := make(map[string]json.RawMessage)
		for aggName, agg := range t.Aggregations {
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

// TermsAggregationResult represents the results from a terms aggregation request.
type TermsAggregationResult struct {
	DocCountErrorUpperBound uint64
	SumOtherDocCount        uint64
	Buckets                 []TermBucketResult
}

// UnmarshalJSON implements [json.Unmarshaler] to decode a json byte slice into a TermsAggregationResult
// Ignores unknown fields.
func (t *TermsAggregationResult) UnmarshalJSON(m []byte) error {
	// map[key] -> value
	var rawResp map[string]json.RawMessage
	if err := json.Unmarshal(m, &rawResp); err != nil {
		return err
	}

	if t == nil {
		return fmt.Errorf("invalid TermsAggregationResult target, nil")
	}

	for key, value := range rawResp {
		switch key {
		case "doc_count_error_upper_bound":
			if err := json.Unmarshal(value, &t.DocCountErrorUpperBound); err != nil {
				return err
			}
		case "sum_other_doc_count":
			if err := json.Unmarshal(value, &t.SumOtherDocCount); err != nil {
				return err
			}
		case "buckets":
			if err := json.Unmarshal(value, &t.Buckets); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown TermsAggregationResult field %s", key)
		}
	}

	return nil
}

// TermBucketResult is a [AggregationResultMap] for a TermsAggregation
type TermBucketResult struct {
	Key                   string
	DocCount              uint64
	SubAggregationResults map[string]json.RawMessage
}

// UnmarshalJSON implements [json.Unmarshaler] to decode a json byte slice into a TermBucketResult
func (t *TermBucketResult) UnmarshalJSON(m []byte) error {
	// map[key] -> value
	var rawResp map[string]json.RawMessage
	if err := json.Unmarshal(m, &rawResp); err != nil {
		return err
	}

	if t == nil {
		return fmt.Errorf("invalid TermBucketResult target, nil")
	}

	t.SubAggregationResults = make(map[string]json.RawMessage)
	for key, value := range rawResp {
		switch key {
		case "key":
			if err := json.Unmarshal(value, &t.Key); err != nil {
				return err
			}
		case "doc_count":
			if err := json.Unmarshal(value, &t.DocCount); err != nil {
				return err
			}
		default:
			// any number of sub aggregation results
			t.SubAggregationResults[key] = value
		}
	}

	return nil
}

// GetAggregationResultSource implements [opensearchtools.AggregationResultSet] to fetch a sub aggregation result and
// return the raw JSON source for the provided name.
func (t *TermBucketResult) GetAggregationResultSource(name string) ([]byte, bool) {
	if len(t.SubAggregationResults) == 0 {
		return nil, false
	}

	subAggSource, exists := t.SubAggregationResults[name]
	return subAggSource, exists
}
