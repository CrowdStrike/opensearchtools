package search

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTermsAggregation_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  *TermsAggregation
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Case",
			target:  &TermsAggregation{},
			wantErr: true,
		},
		{
			name:    "Basic field only",
			target:  NewTermsAggregation("field"),
			want:    `{"terms":{"field":"field"}}`,
			wantErr: false,
		},
		{
			name: "Terms Aggregation with all options set",
			target: NewTermsAggregation("field").
				AddOrder(NewOrder("field", true)).
				WithSize(10),
			want:    `{"terms":{"field":"field","size":10,"order":[{"field":"desc"}]}}`,
			wantErr: false,
		},
		{
			name: "Terms aggregation with negative size is ignored",
			target: NewTermsAggregation("field").
				WithSize(-5),
			want:    `{"terms":{"field":"field"}}`,
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

func TestTermsAggregation_WithSubAggregations_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  BucketAggregation
		want    string
		wantErr bool
	}{
		{
			name: "single nested terms aggregation",
			target: NewTermsAggregation("field1").
				AddSubAggregation("nested_terms", NewTermsAggregation("field2")),
			want: `{"terms":{"field":"field1"},"aggs":{"nested_terms":{"terms":{"field":"field2"}}}}`,
		},
		{
			name: "double nested terms aggregation",
			target: NewTermsAggregation("field1").
				AddSubAggregation("nested_terms", NewTermsAggregation("field2").
					AddSubAggregation("double_nested", NewTermsAggregation("field3"))),
			want: `{"terms":{"field":"field1"},"aggs":{"nested_terms":{"terms":{"field":"field2"},"aggs":{"double_nested":{"terms":{"field":"field3"}}}}}}`,
		},
		//TODO add other aggregation types as they're created
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

func TestTermsAggregationResult_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		rawJSON []byte
		want    TermsAggregationResult
		wantErr bool
	}{
		{
			name:    "Basic result",
			rawJSON: []byte(`{"doc_count_error_upper_bound":10,"sum_other_doc_count":10,"buckets":[{"key":"field_value","doc_count":10}]}`),
			want: TermsAggregationResult{
				DocCountErrorUpperBound: 10,
				SumOtherDocCount:        10,
				Buckets: []TermBucketResult{{
					Key:                   "field_value",
					DocCount:              10,
					SubAggregationResults: make(map[string]json.RawMessage),
				}},
			},
			wantErr: false,
		},
		{
			name:    "Empty results",
			rawJSON: []byte(`{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[]}`),
			want: TermsAggregationResult{
				DocCountErrorUpperBound: 0,
				SumOtherDocCount:        0,
				Buckets:                 []TermBucketResult{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TermsAggregationResult
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

func TestTermsBucketResult_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		rawJSON []byte
		want    TermBucketResult
		wantErr bool
	}{
		{
			name:    "Basic result",
			rawJSON: []byte(`{"key":"field_value","doc_count":10}`),
			want: TermBucketResult{
				Key:                   "field_value",
				DocCount:              10,
				SubAggregationResults: make(map[string]json.RawMessage),
			},
			wantErr: false,
		},
		{
			name:    "Nested terms aggregation result",
			rawJSON: []byte(`{"key":"field_value","doc_count":10,"nested_terms":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"field_value_nested","doc_count":10}]}}`),
			want: TermBucketResult{
				Key:      "field_value",
				DocCount: 10,
				SubAggregationResults: map[string]json.RawMessage{
					"nested_terms": json.RawMessage(`{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"field_value_nested","doc_count":10}]}`),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TermBucketResult
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
