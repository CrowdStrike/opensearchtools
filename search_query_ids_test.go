package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIDsQuery_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		query   *IDsQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty query",
			query:   &IDsQuery{},
			want:    `{"ids":{"values":null}}`,
			wantErr: false,
		},
		{
			name:    "Basic Constructor",
			query:   NewIDsQuery(),
			want:    `{"ids":{"values":null}}`,
			wantErr: false,
		},
		{
			name:    "Single value",
			query:   NewIDsQuery("value1"),
			want:    `{"ids":{"values":["value1"]}}`,
			wantErr: false,
		},
		{
			name:    "Multiple Values",
			query:   NewIDsQuery("value1", "value2", "value3"),
			want:    `{"ids":{"values":["value1","value2","value3"]}}`,
			wantErr: false,
		},
		{
			name:    "Mixed type values",
			query:   NewIDsQuery("value1", 2),
			want:    `{"ids":{"values":["value1",2]}}`,
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
