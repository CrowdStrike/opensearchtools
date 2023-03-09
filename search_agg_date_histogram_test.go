package opensearchtools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDateHistogramAggregation_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  *DateHistogramAggregation
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Case",
			target:  &DateHistogramAggregation{},
			wantErr: true,
		},
		{
			name:    "Basic Constructor",
			target:  NewDateHistogramAggregation("field", "day"),
			want:    `{"date_histogram":{"field":"field","interval":"day"}}`,
			wantErr: false,
		},
		{
			name: "Date histogram aggregation with all options set",
			target: NewDateHistogramAggregation("field", "day").
				WithMinDocCount(10).
				WithTimeZone("-01:00").
				AddOrder(NewOrder("field", true)),
			want:    `{"date_histogram":{"field":"field","interval":"day","min_doc_count":10,"time_zone":"-01:00","order":[{"field":"desc"}]}}`,
			wantErr: false,
		},
		{
			name: "Date histogram aggregation with negative MinDocCount is ignored",
			target: NewDateHistogramAggregation("field", "day").
				WithMinDocCount(-5),
			want:    `{"date_histogram":{"field":"field","interval":"day"}}`,
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

func TestDateHistogramAggregation_WithSubAggregations_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  BucketAggregation
		want    string
		wantErr bool
	}{
		{
			name: "single nested terms aggregation",
			target: NewDateHistogramAggregation("field", "day").
				AddSubAggregation("nested_terms", NewTermsAggregation("field2")),
			want: `{"date_histogram":{"field":"field","interval":"day"},"aggs":{"nested_terms":{"terms":{"field":"field2"}}}}`,
		},
		{
			name: "double nested terms aggregation",
			target: NewDateHistogramAggregation("field", "day").
				AddSubAggregation("nested_terms", NewTermsAggregation("field2").
					AddSubAggregation("double_nested", NewTermsAggregation("field3"))),
			want: `{"date_histogram":{"field":"field","interval":"day"},"aggs":{"nested_terms":{"terms":{"field":"field2"},"aggs":{"double_nested":{"terms":{"field":"field3"}}}}}}`,
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

func TestDateHistogramAggregationResult_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		rawJSON []byte
		want    DateHistogramAggregationResults
		wantErr bool
	}{
		{
			name:    "Basic result",
			rawJSON: []byte(`{"buckets":[{"key_as_string":"2023-01-26T00:00:00.000Z","key":1674691200000,"doc_count":10}]}`),
			want: DateHistogramAggregationResults{
				Buckets: []DateHistogramBucketResult{{
					KeyString:             "2023-01-26T00:00:00.000Z",
					Key:                   1674691200000,
					DocCount:              10,
					SubAggregationResults: make(map[string]json.RawMessage),
				}},
			},
			wantErr: false,
		},
		{
			name:    "Empty results",
			rawJSON: []byte(`{"buckets":[]}`),
			want: DateHistogramAggregationResults{
				Buckets: []DateHistogramBucketResult{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got DateHistogramAggregationResults
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

func TestDateHistogramBucketResult_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		rawJSON []byte
		want    DateHistogramBucketResult
		wantErr bool
	}{
		{
			name:    "Basic result",
			rawJSON: []byte(`{"key_as_string":"2023-01-26T00:00:00.000Z","key":1674691200000,"doc_count":10}`),
			want: DateHistogramBucketResult{
				KeyString:             "2023-01-26T00:00:00.000Z",
				Key:                   1674691200000,
				DocCount:              10,
				SubAggregationResults: make(map[string]json.RawMessage),
			},
			wantErr: false,
		},
		{
			name:    "Nested terms aggregation result",
			rawJSON: []byte(`{"key_as_string":"2023-01-26T00:00:00.000Z","key":1674691200000,"doc_count":10,"nested_terms":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"field_value_nested","doc_count":10}]}}`),
			want: DateHistogramBucketResult{
				KeyString: "2023-01-26T00:00:00.000Z",
				Key:       1674691200000,
				DocCount:  10,
				SubAggregationResults: map[string]json.RawMessage{
					"nested_terms": json.RawMessage(`{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"field_value_nested","doc_count":10}]}`),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got DateHistogramBucketResult
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
