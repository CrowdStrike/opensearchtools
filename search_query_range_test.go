package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRangeQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *RangeQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Query",
			query:   &RangeQuery{},
			want:    `{"range":{"":{}}}`,
			wantErr: false,
		},
		{
			name:    "Basic Constructor",
			query:   NewRangeQuery("field"),
			want:    `{"range":{"field":{}}}`,
			wantErr: false,
		},
		{
			name:    "Less than",
			query:   NewRangeQuery("field").Lt("value"),
			want:    `{"range":{"field":{"lt":"value"}}}`,
			wantErr: false,
		},
		{
			name:    "Less than or equal to",
			query:   NewRangeQuery("field").Lte("value"),
			want:    `{"range":{"field":{"lte":"value"}}}`,
			wantErr: false,
		},
		{
			name:    "Greater than",
			query:   NewRangeQuery("field").Gt("value"),
			want:    `{"range":{"field":{"gt":"value"}}}`,
			wantErr: false,
		},
		{
			name:    "Greater than or equal to",
			query:   NewRangeQuery("field").Gte("value"),
			want:    `{"range":{"field":{"gte":"value"}}}`,
			wantErr: false,
		},
		{
			name: "All fields set",
			query: NewRangeQuery("field").
				Lt("value").
				Lte("value").
				Gt("value").
				Gte("value"),
			want:    `{"range":{"field":{"lt":"value","lte":"value","gt":"value","gte":"value"}}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.ToOpenSearchJSON()

			if (err != nil) != tt.wantErr {
				t.Errorf("ToOpenSearchJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.JSONEq(t, tt.want, string(got))
		})
	}
}
