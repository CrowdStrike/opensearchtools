package opensearchtools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingleValueMetricAggregation_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  *SingleValueMetricAggregation
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Case",
			target:  &SingleValueMetricAggregation{},
			wantErr: true,
		},
		{
			name: "Invalid Type",
			target: &SingleValueMetricAggregation{
				Type:  SingleValueAggType("banana"),
				Field: "field",
			},
			wantErr: true,
		},
		{
			name:    "Cardinality Aggregation",
			target:  NewCardinalityAggregation("field"),
			want:    `{"cardinality":{"field":"field"}}`,
			wantErr: false,
		},
		{
			name: "Cardinality Aggregation with precision",
			target: NewCardinalityAggregation("field").
				WithPrecisionThreshold(10).
				WithMissing("value"),
			want:    `{"cardinality":{"field":"field","precision_threshold":10,"missing":"value"}}`,
			wantErr: false,
		},
		{
			name:    "Maximum Aggregation",
			target:  NewMaximumAggregation("field"),
			want:    `{"max":{"field":"field"}}`,
			wantErr: false,
		},
		{
			name:    "Minimum Aggregation",
			target:  NewMinimumAggregation("field"),
			want:    `{"min":{"field":"field"}}`,
			wantErr: false,
		},
		{
			name:    "Average Aggregation",
			target:  NewAverageAggregation("field"),
			want:    `{"avg":{"field":"field"}}`,
			wantErr: false,
		},
		{
			name:    "Sum Aggregation",
			target:  NewSumAggregation("field"),
			want:    `{"sum":{"field":"field"}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.target.ToOpenSearchJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToOpenSearchJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				require.Nilf(t, got, "if an error is returned, no results are expected")
			} else {
				require.JSONEq(t, tt.want, string(got))
			}
		})
	}
}

func TestSingleValueAggregationResult_UnmarshalJSON(t *testing.T) {
	testValue := float64(10)
	tests := []struct {
		name    string
		rawJSON []byte
		want    SingleValueAggregationResult
		wantErr bool
	}{
		{
			name:    "Basic result",
			rawJSON: []byte(`{"value":10}`),
			want: SingleValueAggregationResult{
				Value: &testValue,
			},
			wantErr: false,
		},
		{
			name:    "Value and value string",
			rawJSON: []byte(`{"value":10,"value_as_string":"10"}`),
			want: SingleValueAggregationResult{
				Value:       &testValue,
				ValueString: "10",
			},
			wantErr: false,
		},
		{
			name:    "No results",
			rawJSON: []byte(`{"value":null}`),
			want:    SingleValueAggregationResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got SingleValueAggregationResult
			gotErr := json.Unmarshal(tt.rawJSON, &got)

			if (gotErr != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", gotErr, tt.wantErr)
			}

			if gotErr == nil {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
