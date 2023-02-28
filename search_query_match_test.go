package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *MatchQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Query",
			query:   &MatchQuery{},
			want:    `{"match":{"":{"query":"","operator":""}}}`,
			wantErr: false,
		},
		{
			name:    "Simple Success",
			query:   NewMatchQuery("field", "value"),
			want:    `{"match":{"field":{"query":"value","operator":"or"}}}`,
			wantErr: false,
		},
		{
			name:    "No value",
			query:   NewMatchQuery("field", ""),
			want:    `{"match":{"field":{"query":"","operator":"or"}}}`,
			wantErr: false,
		},
		{
			name:    "Different operator",
			query:   NewMatchQuery("field", "value").SetOperator("and"),
			want:    `{"match":{"field":{"query":"value","operator":"and"}}}`,
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
