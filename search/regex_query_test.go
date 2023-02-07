package search

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *RegexQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Query",
			query:   &RegexQuery{},
			want:    `{"regexp":{"":""}}`,
			wantErr: false,
		},
		{
			name:    "Basic Constructor",
			query:   NewRegexQuery("field", "^value$"),
			want:    `{"regexp":{"field":"^value$"}}`,
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
