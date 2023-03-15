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
	// map[key] -> value
	var rawResp map[string]json.RawMessage
	if err := json.Unmarshal(b, &rawResp); err != nil {
		return err
	}

	if p == nil {
		return fmt.Errorf("invalid PercentilesAggregationResult target, nil")
	}

	valueMessage, ok := rawResp["values"]
	if !ok {
		return fmt.Errorf("PercentilesAggretgationResult missing expected value entry")
	}

	var rawValues map[string]json.RawMessage
	if err := json.Unmarshal(valueMessage, &rawValues); err != nil {
		return err
	}

	for key, value := range rawValues {
		switch key {
		case "1.0":
			if err := json.Unmarshal(value, &p.P1); err != nil {
				return err
			}
		case "1.0_as_string":
			if err := json.Unmarshal(value, &p.P1String); err != nil {
				return err
			}
		case "5.0":
			if err := json.Unmarshal(value, &p.P5); err != nil {
				return err
			}
		case "5.0_as_string":
			if err := json.Unmarshal(value, &p.P5String); err != nil {
				return err
			}
		case "25.0":
			if err := json.Unmarshal(value, &p.P25); err != nil {
				return err
			}
		case "25.0_as_string":
			if err := json.Unmarshal(value, &p.P25String); err != nil {
				return err
			}
		case "50.0":
			if err := json.Unmarshal(value, &p.P50); err != nil {
				return err
			}
		case "50.0_as_string":
			if err := json.Unmarshal(value, &p.P50String); err != nil {
				return err
			}
		case "75.0":
			if err := json.Unmarshal(value, &p.P75); err != nil {
				return err
			}
		case "75.0_as_string":
			if err := json.Unmarshal(value, &p.P75String); err != nil {
				return err
			}
		case "95.0":
			if err := json.Unmarshal(value, &p.P95); err != nil {
				return err
			}
		case "95.0_as_string":
			if err := json.Unmarshal(value, &p.P95String); err != nil {
				return err
			}
		case "99.0":
			if err := json.Unmarshal(value, &p.P99); err != nil {
				return err
			}
		case "99.0_as_string":
			if err := json.Unmarshal(value, &p.P99String); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected field %s", key)
		}
	}

	return nil
}
