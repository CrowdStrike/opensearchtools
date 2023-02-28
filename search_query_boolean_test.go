package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoolQuery_ToOpenSearchJSON(t *testing.T) {
	basicQuery := NewTermQuery("field", "value")

	tests := []struct {
		name    string
		query   *BoolQuery
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Bool Query",
			query:   NewBoolQuery(),
			want:    `{"bool":{}}`,
			wantErr: false,
		},
		{
			name:    "Bool query must with single query",
			query:   NewBoolQuery().Must(basicQuery),
			want:    `{"bool":{"must":[{"term":{"field":"value"}}]}}`,
			wantErr: false,
		},
		{
			name:    "Bool query must not with single query",
			query:   NewBoolQuery().MustNot(NewTermQuery("field", "value")),
			want:    `{"bool":{"must_not":[{"term":{"field":"value"}}]}}`,
			wantErr: false,
		},
		{
			name:    "Bool query should with single query",
			query:   NewBoolQuery().Should(NewTermQuery("field", "value")),
			want:    `{"bool":{"should":[{"term":{"field":"value"}}]}}`,
			wantErr: false,
		},
		{
			name:    "Bool query filter with single query",
			query:   NewBoolQuery().Filter(NewTermQuery("field", "value")),
			want:    `{"bool":{"filter":[{"term":{"field":"value"}}]}}`,
			wantErr: false,
		},
		{
			name:    "Bool query with minimum should match",
			query:   NewBoolQuery().MinimumShouldMatch(1),
			want:    `{"bool":{"minimum_should_match":1}}`,
			wantErr: false,
		},
		{
			name: "Bool query with everything",
			query: NewBoolQuery().
				Must(basicQuery).
				MustNot(basicQuery).
				Should(basicQuery).
				Filter(basicQuery).
				MinimumShouldMatch(1),
			want:    `{"bool":{"minimum_should_match":1,"must":[{"term":{"field":"value"}}],"must_not":[{"term":{"field":"value"}}],"should":[{"term":{"field":"value"}}],"filter":[{"term":{"field":"value"}}]}}`, //nolint:lll
			wantErr: false,
		},
		{
			name:    "Bool query with multiple sub queries",
			query:   NewBoolQuery().Must(basicQuery, basicQuery),
			want:    `{"bool":{"must":[{"term":{"field":"value"}},{"term":{"field":"value"}}]}}`,
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
