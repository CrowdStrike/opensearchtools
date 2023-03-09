package opensearchtools

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/maps"
)

// FilterAggregation is a query clause, exactly like a search query.
// You can use the FilterAggregation to narrow down the entire set of
// documents to a specific set before creating buckets.
// An empty FilterAggregation will fail to execute as a filter is required.
//
// For more details see https://opensearch.org/docs/latest/opensearch/bucket-agg/#filter-filters
type FilterAggregation struct {
	// Filter to be applied to the document set before aggregating
	Filter Query

	// Aggregations to be performed on the reduced set
	Aggregations map[string]Aggregation
}

// NewFilterAggregation instantiates a FilterAggregation with the provided filter
func NewFilterAggregation(filter Query) *FilterAggregation {
	return &FilterAggregation{
		Filter:       filter,
		Aggregations: make(map[string]Aggregation),
	}
}

// AddSubAggregation to the FilterAggregation with the provided name
// Implements [BucketAggregation.AddSubAggregation]
func (f *FilterAggregation) AddSubAggregation(name string, agg Aggregation) BucketAggregation {
	if f.Aggregations == nil {
		f.Aggregations = map[string]Aggregation{name: agg}
	} else {
		f.Aggregations[name] = agg
	}

	return f
}

// SubAggregations returns all aggregations added to the bucket aggregation.
// Implements [BucketAggregation.SubAggregations]
func (f *FilterAggregation) SubAggregations() map[string]Aggregation {
	return f.Aggregations
}

// Validate that the aggregation is executable.
// Implements [Aggregation.Validate].
func (f *FilterAggregation) Validate() ValidationResults {
	vrs := NewValidationResults()

	if f.Filter == nil {
		vrs.Add(NewValidationResult("a FilterAggregation requires a filter query", true))
	}

	for _, subAgg := range f.Aggregations {
		vrs.Extend(subAgg.Validate())
	}

	return vrs
}

// ToOpenSearchJSON converts the FilterAggregation to the correct OpenSearch JSON.
// Implements [Aggregation.ToOpenSearchJSON].
func (f *FilterAggregation) ToOpenSearchJSON() ([]byte, error) {
	if vrs := f.Validate(); vrs.IsFatal() {
		return nil, NewValidationError(vrs)
	}

	filterJSON, filterErr := f.Filter.ToOpenSearchJSON()
	if filterErr != nil {
		return nil, filterErr
	}

	fa := map[string]any{
		"filter": json.RawMessage(filterJSON),
	}

	if len(f.Aggregations) > 0 {
		subAggs := make(map[string]json.RawMessage)
		for aggName, agg := range f.Aggregations {
			aggJSON, jErr := agg.ToOpenSearchJSON()
			if jErr != nil {
				return nil, jErr
			}

			subAggs[aggName] = aggJSON
		}

		fa["aggs"] = subAggs
	}

	return json.Marshal(fa)
}

// FilterAggregationResults is a [AggregationResultMap] for a FilterAggregation
type FilterAggregationResults struct {
	DocCount              uint64
	SubAggregationResults map[string]json.RawMessage
}

// UnmarshalJSON implements [json.Unmarshaler] to decode a json byte slice into a FilterAggregationResults.
// Unknown fields are assumed to be SubAggregation results
func (f *FilterAggregationResults) UnmarshalJSON(m []byte) error {
	// map[key] -> value
	var rawResp map[string]json.RawMessage
	if err := json.Unmarshal(m, &rawResp); err != nil {
		return err
	}

	if f == nil {
		return fmt.Errorf("invalid TermBucketResult target, nil")
	}

	f.SubAggregationResults = make(map[string]json.RawMessage)
	for key, value := range rawResp {
		switch key {
		case "doc_count":
			if err := json.Unmarshal(value, &f.DocCount); err != nil {
				return err
			}
		default:
			// any number of sub aggregation results
			f.SubAggregationResults[key] = value
		}
	}

	return nil
}

// GetAggregationResultSource implements [opensearchtools.AggregationResultSet] to fetch a sub aggregation result and
// return the raw JSON source for the provided name.
func (f *FilterAggregationResults) GetAggregationResultSource(name string) ([]byte, bool) {
	if len(f.SubAggregationResults) == 0 {
		return nil, false
	}

	subAggSource, exists := f.SubAggregationResults[name]
	return subAggSource, exists
}

// Keys implemented for [opensearchtools.AggregationResultSet] to return the list of aggregation result keys
func (f *FilterAggregationResults) Keys() []string {
	return maps.Keys(f.SubAggregationResults)
}
