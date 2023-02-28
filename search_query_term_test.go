package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTermQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *TermQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Query",
			query:   &TermQuery{},
			want:    `{"term":{"":null}}`,
			wantErr: false,
		},
		{
			name:    "Simple Success",
			query:   NewTermQuery("field", "value"),
			want:    `{"term":{"field":"value"}}`,
			wantErr: false,
		},
		{
			name:    "Search for empty value",
			query:   NewTermQuery("field", ""),
			want:    `{"term":{"field":""}}`,
			wantErr: false,
		},
		{
			name:    "Search for empty field and value",
			query:   NewTermQuery("", ""),
			want:    `{"term":{"":""}}`,
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
