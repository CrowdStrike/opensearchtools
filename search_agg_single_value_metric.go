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

	// Missing is used to define how documents missing are missing a value should be treated.
	// For SingleValueMetricAggregations, Missing is the value that will be substituted in if
	// the document does not contain the target Field
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
func (c *SingleValueMetricAggregation) WithPrecisionThreshold(p int) *SingleValueMetricAggregation {
	c.PrecisionThreshold = p
	return c
}

func (c *SingleValueMetricAggregation) WithMissing(missing any) *SingleValueMetricAggregation {
	c.Missing = missing
	return c
}

// ToOpenSearchJSON converts the SingleValueMetricAggregation to the correct OpenSearch JSON.
func (c *SingleValueMetricAggregation) ToOpenSearchJSON() ([]byte, error) {
	if c.Field == "" {
		return nil, fmt.Errorf("a SingleValueMetricAggregation requires a target field")
	}

	if _, valid := validSingleValueAggTypes[c.Type]; !valid {
		return nil, fmt.Errorf("%s is not a valid type of SingleValueMetricAggregation", c.Type)
	}

	ca := map[string]any{
		"field": c.Field,
	}

	if c.PrecisionThreshold >= 0 {
		ca["precision_threshold"] = c.PrecisionThreshold
	}

	if c.Missing != nil {
		ca["missing"] = c.Missing
	}

	source := map[string]any{
		string(c.Type): ca,
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
