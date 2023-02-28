package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWildcardQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *WildcardQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Query",
			query:   &WildcardQuery{},
			want:    `{"wildcard":{"":""}}`,
			wantErr: false,
		},
		{
			name:    "Simple Success",
			query:   NewWildcardQuery("field", "value"),
			want:    `{"wildcard":{"field":"value"}}`,
			wantErr: false,
		},
		{
			name:    "Search for empty value",
			query:   NewWildcardQuery("field", ""),
			want:    `{"wildcard":{"field":""}}`,
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
