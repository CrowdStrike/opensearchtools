package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNestedQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *NestedQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Query",
			query:   &NestedQuery{},
			wantErr: true,
		},
		{
			name:  "Simple Success",
			query: NewNestedQuery("path", NewTermQuery("field", "value")),
			want:  `{"nested":{"path":"path","query":{"term":{"field":"value"}}}}`,
		},
		{
			name:    "Missing Path",
			query:   NewNestedQuery("", NewTermQuery("field", "value")),
			wantErr: true,
		},
		{
			name:    "Missing Query",
			query:   NewNestedQuery("path", nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.ToOpenSearchJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToOpenSearchJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				require.JSONEq(t, tt.want, string(got))
			}
		})
	}
}
