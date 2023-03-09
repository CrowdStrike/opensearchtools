package opensearchtools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPercentileAggregation_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  *PercentilesAggregation
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Case",
			target:  &PercentilesAggregation{},
			wantErr: true,
		},
		{
			name:    "Basic field only",
			target:  NewPercentileAggregation("field"),
			want:    `{"percentiles":{"field":"field"}}`,
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

func TestPercentileAggregationResult_UnmarshalJSON(t *testing.T) {
	testValue := float64(1)
	tests := []struct {
		name    string
		rawJSON []byte
		want    PercentileAggregationResult
		wantErr bool
	}{
		{
			name:    "Basic result",
			rawJSON: []byte(`{"1.0":1,"5.0":1,"25.0":1,"50.0":1,"75.0":1,"95.0":1,"99.0":1}`),
			want: PercentileAggregationResult{
				P1:  &testValue,
				P5:  &testValue,
				P25: &testValue,
				P50: &testValue,
				P75: &testValue,
				P95: &testValue,
				P99: &testValue,
			},
			wantErr: false,
		},
		{
			name:    "Value and value string",
			rawJSON: []byte(`{"1.0":1,"1.0_as_string":"1","5.0":1,"5.0_as_string":"1","25.0":1,"25.0_as_string":"1","50.0":1,"50.0_as_string":"1","75.0":1,"75.0_as_string":"1","95.0":1,"95.0_as_string":"1","99.0":1,"99.0_as_string":"1"}`),
			want: PercentileAggregationResult{
				P1:        &testValue,
				P1String:  "1",
				P5:        &testValue,
				P5String:  "1",
				P25:       &testValue,
				P25String: "1",
				P50:       &testValue,
				P50String: "1",
				P75:       &testValue,
				P75String: "1",
				P95:       &testValue,
				P95String: "1",
				P99:       &testValue,
				P99String: "1",
			},
			wantErr: false,
		},
		{
			name:    "No results",
			rawJSON: []byte(`{"1.0":null,"5.0":null,"25.0":null,"50.0":null,"75.0":null,"95.0":null,"99.0":null}`),
			want:    PercentileAggregationResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got PercentileAggregationResult
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
