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
		}
	}

	var values valuesJSON

	if err := json.Unmarshal(b, &values); err != nil {
		return err
	}

	if p == nil {
		return fmt.Errorf("invalid PercentilesAggregationResult target, nil")
	}

	p.P1 = values.Values.P1
	p.P1String = values.Values.P1String
	p.P5 = values.Values.P5
	p.P5String = values.Values.P5String
	p.P25 = values.Values.P25
	p.P25String = values.Values.P25String
	p.P50 = values.Values.P50
	p.P50String = values.Values.P50String
	p.P75 = values.Values.P75
	p.P75String = values.Values.P75String
	p.P95 = values.Values.P95
	p.P95String = values.Values.P95String
	p.P99 = values.Values.P99
	p.P99String = values.Values.P99String

	return nil
}
