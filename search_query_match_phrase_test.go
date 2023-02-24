package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchPhraseQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *MatchPhraseQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty query",
			query:   &MatchPhraseQuery{},
			want:    `{"match_phrase":{"":""}}`,
			wantErr: false,
		},
		{
			name:    "Simple Success",
			query:   NewMatchPhraseQuery("field", "phrase"),
			want:    `{"match_phrase":{"field":"phrase"}}`,
			wantErr: false,
		},
		{
			name:    "No phrase",
			query:   NewMatchPhraseQuery("field", ""),
			want:    `{"match_phrase":{"field":""}}`,
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
