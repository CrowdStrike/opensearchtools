package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrefixQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *PrefixQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Query",
			query:   &PrefixQuery{},
			want:    `{"prefix":{"":null}}`,
			wantErr: false,
		},
		{
			name:    "Simple Success",
			query:   NewPrefixQuery("field", "value"),
			want:    `{"prefix":{"field":"value"}}`,
			wantErr: false,
		},
		{
			name:    "Empty Value",
			query:   NewPrefixQuery("field", ""),
			want:    `{"prefix":{"field":""}}`,
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
