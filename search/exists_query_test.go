package search

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExistsQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *ExistsQuery
		want    string
		wantErr bool
	}{
		{
			name:    "empty exists query",
			query:   &ExistsQuery{},
			want:    `{"exists":{"field":""}}`,
			wantErr: false,
		},
		{
			name:    "simple exists query",
			query:   NewExistsQuery("field"),
			want:    `{"exists":{"field":"field"}}`,
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
