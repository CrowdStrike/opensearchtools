package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBucketSortAggregation_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		target  *BucketSortAggregation
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Case",
			target:  &BucketSortAggregation{},
			want:    `{"bucket_sort":{"from":0,"size":0}}`,
			wantErr: false,
		},
		{
			name: "BucketSortAggregation with all options",
			target: NewBucketSortAggregation().
				WithSize(10).
				WithFrom(10).
				AddSort("sort_field", true),
			want:    `{"bucket_sort":{"sort":[{"sort_field":{"order":"desc"}}],"from":10,"size":10}}`,
			wantErr: false,
		},
		{
			name: "BucketSortAggregation with negative size is ignored",
			target: (&BucketSortAggregation{}).
				WithSize(-5),
			want:    `{"bucket_sort":{"from":0}}`,
			wantErr: false,
		},
		{
			name: "BucketSortAggregation with negative from is ignored",
			target: (&BucketSortAggregation{}).
				WithFrom(-5),
			want:    `{"bucket_sort":{"size":0}}`,
			wantErr: false,
		},
		{
			name:    "BucketSortAggregation cannot sort an empty field",
			target:  NewBucketSortAggregation().AddSort("", true),
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
