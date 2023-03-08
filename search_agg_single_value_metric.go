package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// SingleValueAggType represents valid single value metric aggregations
type SingleValueAggType string

const (
	CardinalityAggregation = SingleValueAggType("cardinality")
	MaximumAggregation     = SingleValueAggType("max")
	MinimumAggregation     = SingleValueAggType("min")
	AverageAggregation     = SingleValueAggType("avg")
	SumAggregation         = SingleValueAggType("sum")
)

var (
	validSingleValueAggTypes = map[SingleValueAggType]struct{}{
		CardinalityAggregation: {},
		MaximumAggregation:     {},
		MinimumAggregation:     {},
		AverageAggregation:     {},
		SumAggregation:         {},
	}
)

// SingleValueMetricAggregation is a union of the single-value metric aggregations.
//
//   - Cardinality
//   - Sum
//   - Min
//   - Max
//   - Average
//
// Metric aggregations let you perform simple calculations on values of a field. An empty
// SingleValueMetricAggregation will have some issues with execution:
//
//   - the target Field must be non-nil and non-empty
//   - the Type field must be non-empty and matching a SingleValueAggType
//
// For more details see https://opensearch.org/docs/latest/opensearch/metric-agg
type SingleValueMetricAggregation struct {
	// Type of metric aggregation to be performed
	Type SingleValueAggType

	// Field to be bucketed
	Field string

	// PrecisionThreshold is used only for CardinalityAggregation. It defines the
	// threshold below which counts are expected to be close to accurate.
	// Negative values will be omitted
	PrecisionThreshold int

	// Missing is used to define how documents missing the target Field.
	// The value of Missing is substituted for the document.
	Missing any
}

// NewCardinalityAggregation instantiates a SingleValueMetricAggregation with type CardinalityAggregation,
// targeting the provided field. Sets PrecisionThreshold to -1 to be omitted.
func NewCardinalityAggregation(field string) *SingleValueMetricAggregation {
	return &SingleValueMetricAggregation{
		Type:               CardinalityAggregation,
		Field:              field,
		PrecisionThreshold: -1,
	}
}

// NewMaximumAggregation instantiates a SingleValueMetricAggregation with type MaximumAggregation,
// targeting the provided field. Sets PrecisionThreshold to -1 to be omitted.
func NewMaximumAggregation(field string) *SingleValueMetricAggregation {
	return &SingleValueMetricAggregation{
		Type:               MaximumAggregation,
		Field:              field,
		PrecisionThreshold: -1,
	}
}

// NewMinimumAggregation instantiates a SingleValueMetricAggregation with type MinimumAggregation,
// targeting the provided field. Sets PrecisionThreshold to -1 to be omitted.
func NewMinimumAggregation(field string) *SingleValueMetricAggregation {
	return &SingleValueMetricAggregation{
		Type:               MinimumAggregation,
		Field:              field,
		PrecisionThreshold: -1,
	}
}

// NewAverageAggregation instantiates a SingleValueMetricAggregation with type AverageAggregation,
// targeting the provided field. Sets PrecisionThreshold to -1 to be omitted.
func NewAverageAggregation(field string) *SingleValueMetricAggregation {
	return &SingleValueMetricAggregation{
		Type:               AverageAggregation,
		Field:              field,
		PrecisionThreshold: -1,
	}
}

// NewSumAggregation instantiates a SingleValueMetricAggregation with type SumAggregationion,
// targeting the provided field. Sets PrecisionThreshold to -1 to be omitted.
func NewSumAggregation(field string) *SingleValueMetricAggregation {
	return &SingleValueMetricAggregation{
		Type:               SumAggregation,
		Field:              field,
		PrecisionThreshold: -1,
	}
}

// WithPrecisionThreshold sets the PrecisionThreshold
func (s *SingleValueMetricAggregation) WithPrecisionThreshold(p int) *SingleValueMetricAggregation {
	s.PrecisionThreshold = p
	return s
}

// WithMissing value to use
func (s *SingleValueMetricAggregation) WithMissing(missing any) *SingleValueMetricAggregation {
	s.Missing = missing
	return s
}

// Validate that the aggregation is executable.
// Implements [Aggregation.Validate].
func (s *SingleValueMetricAggregation) Validate() ValidationResults {
	vrs := NewValidationResults()

	if s.Field == "" {
		vrs.Add(NewValidationResult("a SingleValueMetricAggregation requires a target field", true))
	}

	if _, valid := validSingleValueAggTypes[s.Type]; !valid {
		vrs.Add(NewValidationResult(fmt.Sprintf("%s is not a valid type of SingleValueMetricAggregation", s.Type), true))
	}

	return vrs
}

// ToOpenSearchJSON converts the SingleValueMetricAggregation to the correct OpenSearch JSON.
// Implements [Aggregation.ToOpenSearchJSON].
func (s *SingleValueMetricAggregation) ToOpenSearchJSON() ([]byte, error) {
	if vrs := s.Validate(); vrs.IsFatal() {
		return nil, NewValidationError(vrs)
	}

	ca := map[string]any{
		"field": s.Field,
	}

	if s.PrecisionThreshold >= 0 {
		ca["precision_threshold"] = s.PrecisionThreshold
	}

	if s.Missing != nil {
		ca["missing"] = s.Missing
	}

	source := map[string]any{
		string(s.Type): ca,
	}

	return json.Marshal(source)
}

// SingleValueAggregationResult union of the single-value metric aggregations results.
type SingleValueAggregationResult struct {
	// Value will return nil if not metrics were aggregated
	Value *float64 `json:"value,omitempty"`

	// ValueString is optional depending on the field type of the value being aggregated
	ValueString string `json:"value_as_string,omitempty"`
}
