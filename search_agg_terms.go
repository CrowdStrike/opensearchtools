package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// TermsAggregation dynamically creates a bucket for each unique term of a field.
// An empty TermsAggregation will have some issues with execution:
//   - the target Field must be non-nil and non-empty.
//   - a Size of 0 will return no buckets
//
// For more details see https://opensearch.org/docs/latest/opensearch/bucket-agg/
type TermsAggregation struct {
	// Field to be bucketed
	Field string

	// Size of the number of buckets to be returned. Negative sizes will be omitted
	Size int

	// MinDocCount is the lower count threshold for a bucket to be included in the results.
	// Negative counts will be omitted
	MinDocCount int64

	// Missing counts documents that are missing the field being aggregated
	Missing string

	// Include filters values based on a regexp, Include cannot be used in tandem with IncludeValues
	Include string

	// IncludeValues filters values base on a list of exact matches,
	// IncludeValues cannot be used in tandem with Include
	IncludeValues []string

	// Exclude filters values based on a regexp, Exclude cannot be used in tandem with ExcludeValues
	Exclude string

	// ExcludeValues filters values base on a list of exact matches,
	// ExcludeValues cannot be used in tandem with Exclude
	ExcludeValues []string

	// Order list of [Order]s to sort the aggregation buckets. Default order is _count: desc
	Order []Order

	// Aggregations sub aggregations for each bucket. Mapped by string label to sub aggregation
	Aggregations map[string]Aggregation
}

// NewTermsAggregation instantiates a TermsAggregation targeting the provided field
// Sets Size and MinDocCount to -1 to be omitted for the default value.
func NewTermsAggregation(field string) *TermsAggregation {
	return &TermsAggregation{
		Field:        field,
		Size:         -1,
		MinDocCount:  -1,
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

// WithMinDocCount the lower count threshold for a bucket to be included in the results
func (t *TermsAggregation) WithMinDocCount(minCount int64) *TermsAggregation {
	t.MinDocCount = minCount
	return t
}

// WithMissing buckets documents missing the field under the provided label
func (t *TermsAggregation) WithMissing(missing string) *TermsAggregation {
	t.Missing = missing
	return t
}

// WithInclude sets the regex include filter
func (t *TermsAggregation) WithInclude(include string) *TermsAggregation {
	t.Include = include
	return t
}

// WithIncludes sets the list of include matches
func (t *TermsAggregation) WithIncludes(include []string) *TermsAggregation {
	t.IncludeValues = include
	return t
}

// WithExclude sets the regex exclude filter
func (t *TermsAggregation) WithExclude(exclude string) *TermsAggregation {
	t.Exclude = exclude
	return t
}

// WithExcludes sets the list of Exclude matches
func (t *TermsAggregation) WithExcludes(excludes []string) *TermsAggregation {
	t.ExcludeValues = excludes
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

	//TODO: PR Question - Should we validate like this? Or would it make sense to add `Validate() ValidationResults` to the aggregation interface.
	// Then a SearchRequest could call Validate on all of the aggregations before marshaling. And we could leverage it at the beginning of this method.
	if t.Include != "" && len(t.IncludeValues) > 0 {
		return nil, fmt.Errorf("terms agg cannot have both Include [%s] and IncludeValues [%v] set", t.Include, t.IncludeValues)
	}

	if t.Include != "" {
		ta["include"] = t.Include
	}

	if len(t.IncludeValues) > 0 {
		ta["include"] = t.IncludeValues
	}

	if t.Exclude != "" && len(t.ExcludeValues) > 0 {
		return nil, fmt.Errorf("terms agg cannot have both Exclude [%s] and ExcludeValues [%v] set", t.Exclude, t.ExcludeValues)
	}

	if t.Exclude != "" {
		ta["exclude"] = t.Exclude
	}

	if len(t.ExcludeValues) > 0 {
		ta["exclude"] = t.ExcludeValues
	}

	if t.MinDocCount >= 0 {
		ta["min_doc_count"] = t.MinDocCount
	}

	if t.Missing != "" {
		ta["missing"] = t.Missing
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

// TermsAggregationResults represents the results from a terms aggregation request.
type TermsAggregationResults struct {
	DocCountErrorUpperBound int64
	SumOtherDocCount        int64
	Buckets                 []TermBucketResult
}

// UnmarshalJSON implements [json.Unmarshaler] to decode a json byte slice into a TermsAggregationResults
// Errors on unknown fields.
func (t *TermsAggregationResults) UnmarshalJSON(m []byte) error {
	// map[key] -> value
	var rawResp map[string]json.RawMessage
	if err := json.Unmarshal(m, &rawResp); err != nil {
		return err
	}

	if t == nil {
		return fmt.Errorf("invalid TermsAggregationResults target, nil")
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
			return fmt.Errorf("unknown TermsAggregationResults field %s", key)
		}
	}

	return nil
}

// TermBucketResult is a [AggregationResultMap] for a TermsAggregation
type TermBucketResult struct {
	Key                   string
	DocCount              int64
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

// Keys implemented for [opensearchtools.AggregationResultSet] to return the list of aggregation result keys
func (t *TermBucketResult) Keys() []string {
	keys := make([]string, len(t.SubAggregationResults))

	i := 0
	for k := range t.SubAggregationResults {
		keys[i] = k
		i++
	}

	return keys
}
