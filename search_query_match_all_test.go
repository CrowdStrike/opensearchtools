package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchAllQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *MatchAllQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Basic Constructor",
			query:   NewMatchAllQuery(),
			want:    `{"match_all":{}}`,
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
