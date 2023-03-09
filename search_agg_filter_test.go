package opensearchtools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterAggregation_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  *FilterAggregation
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Case",
			target:  &FilterAggregation{},
			wantErr: true,
		},
		{
			name:    "Basic Filter",
			target:  NewFilterAggregation(NewMatchAllQuery()),
			want:    `{"filter":{"match_all":{}}}`,
			wantErr: false,
		},
		{
			name: "Nested terms aggregation",
			target: &FilterAggregation{
				Filter: NewMatchAllQuery(),
				Aggregations: map[string]Aggregation{
					"nested_terms": NewTermsAggregation("field"),
				},
			},
			want: `{"filter":{"match_all":{}},"aggs":{"nested_terms":{"terms":{"field":"field"}}}}`,
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

func TestFilterAggregationResults_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		rawJSON []byte
		want    FilterAggregationResults
		wantErr bool
	}{
		{
			name:    "Basic results",
			rawJSON: []byte(`{"doc_count":2}`),
			want: FilterAggregationResults{
				DocCount:              2,
				SubAggregationResults: map[string]json.RawMessage{},
			},
			wantErr: false,
		},
		{
			name:    "Nested terms aggregation",
			rawJSON: []byte(`{"doc_count":2,"nested_terms":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"field","doc_count":2}]}}`),
			want: FilterAggregationResults{
				DocCount: 2,
				SubAggregationResults: map[string]json.RawMessage{
					"nested_terms": []byte(`{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"field","doc_count":2}]}`),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got FilterAggregationResults
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
