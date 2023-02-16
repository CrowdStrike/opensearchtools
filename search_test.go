package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CrowdStrike/opensearchtools/search"
)

func TestSearchRequest_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		search  *SearchRequest
		want    string
		wantErr bool
	}{
		{
			name:    "Basic Constructor",
			search:  NewSearchRequest(),
			want:    `{}`,
			wantErr: false,
		},
		{
			name: "All Fields",
			search: NewSearchRequest().
				SetQuery(search.NewTermQuery("field", "value")).
				AddIndices("test_index").
				AddSort(search.NewSort("field", true)).
				SetSize(1),
			want:    `{"query":{"term":{"field":"value"}},"sort":[{"field":{"order":"desc"}}],"size":1}`,
			wantErr: false,
		},
		{
			name: "Set Query",
			search: NewSearchRequest().
				SetQuery(search.NewTermQuery("field", "value")),
			want:    `{"query":{"term":{"field":"value"}}}`,
			wantErr: false,
		},
		{
			name: "Set Index", // Query param so no effect on JSON
			search: NewSearchRequest().
				AddIndices("test_index"),
			want:    `{}`,
			wantErr: false,
		},
		{
			name: "Single Sort",
			search: NewSearchRequest().
				AddSort(search.NewSort("field", true)),
			want:    `{"sort":[{"field":{"order":"desc"}}]}`,
			wantErr: false,
		},
		{
			name: "Multi sort",
			search: NewSearchRequest().
				AddSort(search.NewSort("field", true), search.NewSort("field2", false)),
			want:    `{"sort":[{"field":{"order":"desc"}},{"field2":{"order":"asc"}}]}`,
			wantErr: false,
		},
		{
			name: "Set Size",
			search: NewSearchRequest().
				SetSize(1),
			want:    `{"size":1}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.search.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.JSONEq(t, tt.want, string(got))
		})
	}
}
