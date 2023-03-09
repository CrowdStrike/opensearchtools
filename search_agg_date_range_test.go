package opensearchtools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDateRangeAggregation_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  *DateRangeAggregation
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Case",
			target:  &DateRangeAggregation{},
			wantErr: true,
		},
		{
			name:    "Basic field only",
			target:  NewDateRangeAggregation("field"),
			wantErr: true,
		},
		{
			name:    "Empty field fail",
			target:  NewDateRangeAggregation(""),
			wantErr: true,
		},
		{
			name: "DateRange with format",
			target: NewDateRangeAggregation("field").
				AddRange(0, 10).
				WithFormat("MM-yyyy"),
			want:    `{"date_range":{"field":"field","format":"MM-yyyy","ranges":[{"from":0,"to":10}]}}`,
			wantErr: false,
		},
		{
			name:    "Range Aggregation with un-keyed bucket",
			target:  NewDateRangeAggregation("field").AddRange(0, 10),
			want:    `{"date_range":{"field":"field","ranges":[{"from":0,"to":10}]}}`,
			wantErr: false,
		},
		{
			name: "Range Aggregation with keyed bucket",
			target: NewDateRangeAggregation("field").
				AddKeyedRange("key", 0, 10),
			want:    `{"date_range":{"field":"field","ranges":[{"key":"key","from":0,"to":10}]}}`,
			wantErr: false,
		},
		{
			name: "Range Aggregation with un-keyed and keyed buckets",
			target: NewDateRangeAggregation("field").
				AddRanges(
					Range{Key: "key", From: 0, To: 10},
					Range{From: 10, To: 20},
				),
			want:    `{"date_range":{"field":"field","ranges":[{"key":"key","from":0,"to":10},{"from":10,"to":20}]}}`,
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

func TestDateRangeAggregation_WithSubAggregations_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  BucketAggregation
		want    string
		wantErr bool
	}{
		{
			name: "single nested terms aggregation",
			target: NewDateRangeAggregation("field").AddRange(0, 10).
				AddSubAggregation("nested_terms", NewTermsAggregation("field2")),
			want: `{"date_range":{"field":"field","ranges":[{"from":0,"to":10}]},"aggs":{"nested_terms":{"terms":{"field":"field2"}}}}`,
		},
		{
			name: "double nested terms aggregation",
			target: NewDateRangeAggregation("field").AddRange(0, 10).
				AddSubAggregation("nested_terms", NewTermsAggregation("field2").
					AddSubAggregation("double_nested", NewTermsAggregation("field3"))),
			want: `{"date_range":{"field":"field","ranges":[{"from":0,"to":10}]},"aggs":{"nested_terms":{"terms":{"field":"field2"},"aggs":{"double_nested":{"terms":{"field":"field3"}}}}}}`,
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

func TestDateRangeAggregationResult_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		rawJSON []byte
		want    DateRangeAggregationResults
		wantErr bool
	}{
		{
			name:    "Basic result",
			rawJSON: []byte(`{"buckets":[{"key":"key","from":0.0,"from_as_string":"0","to":10.0,"to_as_string":"10","doc_count":10}]}`),
			want: DateRangeAggregationResults{
				Buckets: []RangeBucketResult{{
					Key:                   "key",
					DocCount:              10,
					From:                  0.0,
					FromString:            "0",
					To:                    10.0,
					ToString:              "10",
					SubAggregationResults: make(map[string]json.RawMessage),
				}},
			},
			wantErr: false,
		},
		{
			name:    "Empty results", // since one range is always defined, the closest to empty we have is a bucket with no documents
			rawJSON: []byte(`{"buckets":[{"key":"key","from":0.0,"from_as_string":"0","to":10.0,"to_as_string":"10","doc_count":0}]}`),
			want: DateRangeAggregationResults{
				Buckets: []RangeBucketResult{{
					Key:                   "key",
					DocCount:              0,
					From:                  0.0,
					FromString:            "0",
					To:                    10.0,
					ToString:              "10",
					SubAggregationResults: make(map[string]json.RawMessage),
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got DateRangeAggregationResults
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
