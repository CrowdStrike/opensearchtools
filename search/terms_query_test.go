package search

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTermsQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *TermsQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty query",
			query:   &TermsQuery{},
			want:    `{"terms":{"":null}}`,
			wantErr: false,
		},
		{
			name:    "Basic Constructor",
			query:   NewTermsQuery("field"),
			want:    `{"terms":{"field":null}}`,
			wantErr: false,
		},
		{
			name:    "Single value",
			query:   NewTermsQuery("field", "value1"),
			want:    `{"terms":{"field":["value1"]}}`,
			wantErr: false,
		},
		{
			name:    "Multiple Values",
			query:   NewTermsQuery("field", "value1", "value2", "value3"),
			want:    `{"terms":{"field":["value1","value2","value3"]}}`,
			wantErr: false,
		},
		{
			name:    "Mixed type values",
			query:   NewTermsQuery("field", "value1", 2),
			want:    `{"terms":{"field":["value1",2]}}`,
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
