package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// PercentilesAggregation is the percentage of the data that’s at or below a
// certain threshold value. An empty PercentilesAggregation will have some issues with execution:
//   - the target Field must be non-nil and non-empty.
//
// For more details see https://opensearch.org/docs/latest/opensearch/metric-agg/#percentile-percentile_ranks
type PercentilesAggregation struct {
	// Field to be bucketed
	Field string
}

// NewPercentileAggregation instantiates a PercentilesAggregation tergeting the provided field.
func NewPercentileAggregation(field string) *PercentilesAggregation {
	return &PercentilesAggregation{
		Field: field,
	}
}

// Validate that the aggregation is executable.
// Implements [Aggregation.Validate].
func (p *PercentilesAggregation) Validate() ValidationResults {
	vrs := NewValidationResults()

	if p.Field == "" {
		vrs.Add(NewValidationResult("a FilterAggregation requires a filter query", true))
	}

	return vrs
}

// ToOpenSearchJSON converts the PercentilesAggregation to the correct OpenSearch JSON.
// Implements [Aggregation.ToOpenSearchJSON].
func (p *PercentilesAggregation) ToOpenSearchJSON() ([]byte, error) {
	if vrs := p.Validate(); vrs.IsFatal() {
		return nil, NewValidationError(vrs)
	}

	source := map[string]any{
		"percentiles": map[string]any{
			"field": p.Field,
		},
	}

	return json.Marshal(source)
}

// PercentilesAggregationResult will contain all percentiles or no percentiles.
// If there are no values for the percentile, it will be omitted
type PercentilesAggregationResult struct {
	P1        *float64
	P1String  string
	P5        *float64
	P5String  string
	P25       *float64
	P25String string
	P50       *float64
	P50String string
	P75       *float64
	P75String string
	P95       *float64
	P95String string
	P99       *float64
	P99String string
}

// UnmarshalJSON implements [json.Unmarshaler] to decode a json byte slice into a PercentilesAggregationResult
func (p *PercentilesAggregationResult) UnmarshalJSON(b []byte) error {
	if p == nil {
		return fmt.Errorf("invalid PercentilesAggregationResult target, nil")
	}

	type valuesJSON struct {
		Values struct {
			P1        *float64 `json:"1.0,omitempty"`
			P1String  string   `json:"1.0_as_string,omitempty"`
			P5        *float64 `json:"5.0,omitempty"`
			P5String  string   `json:"5.0_as_string,omitempty"`
			P25       *float64 `json:"25.0,omitempty"`
			P25String string   `json:"25.0_as_string,omitempty"`
			P50       *float64 `json:"50.0,omitempty"`
			P50String string   `json:"50.0_as_string,omitempty"`
			P75       *float64 `json:"75.0,omitempty"`
			P75String string   `json:"75.0_as_string,omitempty"`
			P95       *float64 `json:"95.0,omitempty"`
			P95String string   `json:"95.0_as_string,omitempty"`
			P99       *float64 `json:"99.0,omitempty"`
			P99String string   `json:"99.0_as_string,omitempty"`
		} `json:"values"`
	}

	var values valuesJSON

	if err := json.Unmarshal(b, &values); err != nil {
		return err
	}

	// can assign values.Values directly to p since they have the exact same fields
	*p = PercentilesAggregationResult(values.Values)

	return nil
}
