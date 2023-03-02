package opensearchtools

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
				WithSize(10).
				WithInclude("include").
				WithExclude("exclude").
				WithMinDocCount(10).
				WithMissing("missing"),
			want:    `{"terms":{"field":"field","size":10,"include":"include","exclude":"exclude","missing":"missing","min_doc_count":10,"order":[{"field":"desc"}]}}`,
			wantErr: false,
		},
		{
			name: "Terms aggregation with negative size is ignored",
			target: NewTermsAggregation("field").
				WithSize(-5),
			want:    `{"terms":{"field":"field"}}`,
			wantErr: false,
		},
		{
			name: "Terms aggregation with negative MinDocCount is ignored",
			target: NewTermsAggregation("field").
				WithMinDocCount(-5),
			want:    `{"terms":{"field":"field"}}`,
			wantErr: false,
		},
		{
			name: "Terms aggregation with include values",
			target: NewTermsAggregation("field").
				WithIncludes([]string{"1", "2"}),
			want:    `{"terms":{"field":"field","include":["1","2"]}}`,
			wantErr: false,
		},
		{
			name: "Terms aggregation with exclude values",
			target: NewTermsAggregation("field").
				WithExcludes([]string{"1", "2"}),
			want:    `{"terms":{"field":"field","exclude":["1","2"]}}`,
			wantErr: false,
		},
		{
			name: "Terms aggregation with include and include values fails",
			target: NewTermsAggregation("field").
				WithIncludes([]string{"1", "2"}).
				WithInclude("fail"),
			wantErr: true,
		},
		{
			name: "Terms aggregation with exclude and exclude values fails",
			target: NewTermsAggregation("field").
				WithExcludes([]string{"1", "2"}).
				WithExclude("fail"),
			wantErr: true,
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
		want    TermsAggregationResults
		wantErr bool
	}{
		{
			name:    "Basic result",
			rawJSON: []byte(`{"doc_count_error_upper_bound":10,"sum_other_doc_count":10,"buckets":[{"key":"field_value","doc_count":10}]}`),
			want: TermsAggregationResults{
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
			want: TermsAggregationResults{
				DocCountErrorUpperBound: 0,
				SumOtherDocCount:        0,
				Buckets:                 []TermBucketResult{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TermsAggregationResults
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
