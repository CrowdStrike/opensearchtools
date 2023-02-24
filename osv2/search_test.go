package osv2

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CrowdStrike/opensearchtools"
)

func TestSearchRequest_ToOpenSearchJSON(t *testing.T) {
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
				WithQuery(opensearchtools.NewTermQuery("field", "value")).
				AddIndices("test_index").
				AddSorts(opensearchtools.NewSort("field", true)).
				WithSize(1),
			want:    `{"query":{"term":{"field":"value"}},"sort":[{"field":{"order":"desc"}}],"size":1}`,
			wantErr: false,
		},
		{
			name: "Set Query",
			search: NewSearchRequest().
				WithQuery(opensearchtools.NewTermQuery("field", "value")),
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
				AddSorts(opensearchtools.NewSort("field", true)),
			want:    `{"sort":[{"field":{"order":"desc"}}]}`,
			wantErr: false,
		},
		{
			name: "Multi sort",
			search: NewSearchRequest().
				AddSorts(opensearchtools.NewSort("field", true), opensearchtools.NewSort("field2", false)),
			want:    `{"sort":[{"field":{"order":"desc"}},{"field2":{"order":"asc"}}]}`,
			wantErr: false,
		},
		{
			name: "Set Size",
			search: NewSearchRequest().
				WithSize(1),
			want:    `{"size":1}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.search.ToOpenSearchJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestHit_ToDomain(t *testing.T) {
	tests := []struct {
		name   string
		target Hit
		want   opensearchtools.Hit
	}{
		{
			name:   "Empty",
			target: Hit{},
			want:   opensearchtools.Hit{},
		},
		{
			name: "All fields",
			target: Hit{
				Index:  testIndex1,
				ID:     testID1,
				Score:  10,
				Source: json.RawMessage("source"),
			},
			want: opensearchtools.Hit{
				Index:  testIndex1,
				ID:     testID1,
				Score:  10,
				Source: json.RawMessage("source"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.ToDomain()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTotal_ToDomain(t *testing.T) {
	tests := []struct {
		name   string
		target Total
		want   opensearchtools.Total
	}{
		{
			name:   "Empty",
			target: Total{},
			want:   opensearchtools.Total{},
		},
		{
			name: "All fields",
			target: Total{
				Value:    10,
				Relation: "eq",
			},
			want: opensearchtools.Total{
				Value:    10,
				Relation: "eq",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.ToDomain()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestHits_ToDomain(t *testing.T) {
	tests := []struct {
		name   string
		target Hits
		want   opensearchtools.Hits
	}{
		{
			name:   "Empty",
			target: Hits{},
			want:   opensearchtools.Hits{},
		},
		{
			name: "All fields",
			target: Hits{
				Total: Total{
					Value:    1,
					Relation: "eq",
				},
				MaxScore: 10.0,
				Hits: []Hit{
					{
						Index:  testIndex1,
						ID:     testID1,
						Score:  10,
						Source: json.RawMessage("source"),
					},
				},
			},
			want: opensearchtools.Hits{
				Total: opensearchtools.Total{
					Value:    1,
					Relation: "eq",
				},
				MaxScore: 10.0,
				Hits: []opensearchtools.Hit{
					{
						Index:  testIndex1,
						ID:     testID1,
						Score:  10,
						Source: json.RawMessage("source"),
					},
				},
			},
		},
		{
			name: "No results",
			target: Hits{
				Total: Total{
					Value:    0,
					Relation: "eq",
				},
				MaxScore: 0,
				Hits:     nil,
			},
			want: opensearchtools.Hits{
				Total: opensearchtools.Total{
					Value:    0,
					Relation: "eq",
				},
				MaxScore: 0,
				Hits:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.ToDomain()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSearchResponse_ToDomain(t *testing.T) {
	testHeader := http.Header{"header": []string{"value"}}

	tests := []struct {
		name   string
		target SearchResponse
		want   opensearchtools.SearchResponse
	}{
		{
			name:   "Empty",
			target: SearchResponse{},
			want:   opensearchtools.SearchResponse{},
		},
		{
			name: "Successful request",
			target: SearchResponse{
				StatusCode: http.StatusOK,
				Header:     testHeader,
				Took:       100,
				TimedOut:   false,
				Shards:     ShardMeta{Total: 10},
				Hits:       Hits{MaxScore: 10},
				Error:      nil,
			},
			want: opensearchtools.SearchResponse{
				StatusCode: http.StatusOK,
				Header:     testHeader,
				Took:       100,
				TimedOut:   false,
				Shards:     opensearchtools.ShardMeta{Total: 10},
				Hits:       opensearchtools.Hits{MaxScore: 10},
				Error:      nil,
			},
		},
		{
			name: "Unsuccessful request",
			target: SearchResponse{
				StatusCode: http.StatusConflict,
				Header:     testHeader,
				Error: &Error{
					Type:         "It was bad",
					Reason:       "for a reason",
					Index:        "on this index",
					ResourceID:   "with this document",
					ResourceType: "type document",
					IndexUUID:    "asdfasd",
				},
			},
			want: opensearchtools.SearchResponse{
				StatusCode: http.StatusConflict,
				Header:     testHeader,
				Error: &opensearchtools.Error{
					Type:         "It was bad",
					Reason:       "for a reason",
					Index:        "on this index",
					ResourceID:   "with this document",
					ResourceType: "type document",
					IndexUUID:    "asdfasd",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.ToDomain()
			require.Equal(t, tt.want, got)
		})
	}
}
