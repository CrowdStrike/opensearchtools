package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSort_ToOpenSearchJSON(t *testing.T) {
	tests := []struct {
		name    string
		sort    *Sort
		want    string
		wantErr bool
	}{
		{
			name:    "Empty sort",
			sort:    &Sort{},
			want:    `{"":{"order":"asc"}}`,
			wantErr: false,
		},
		{
			name:    "Sort descending",
			sort:    NewSort("field", true),
			want:    `{"field":{"order":"desc"}}`,
			wantErr: false,
		},
		{
			name:    "Sort ascending",
			sort:    NewSort("field", false),
			want:    `{"field":{"order":"asc"}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sort.ToOpenSearchJSON()

			if (err != nil) != tt.wantErr {
				t.Errorf("ToOpenSearchJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.JSONEq(t, tt.want, string(got))
		})
	}
}
