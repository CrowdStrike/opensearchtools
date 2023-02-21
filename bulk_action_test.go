package opensearchtools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type bulkTestDoc struct {
	index, id  string
	OtherField int `json:"other_field"`
}

func (t bulkTestDoc) ID() string {
	return t.id
}

func (t bulkTestDoc) Index() string {
	return t.index
}

func TestActionResponse_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonBytes []byte
		want      ActionResponse
		wantErr   bool
	}{
		{
			name:      "All fields",
			jsonBytes: []byte(`{"type":{"_index":"index","_id":"id","_version":1,"result":"success","_shards":{"total":1,"successful":2,"skipped":3,"failed":4},"_seq_no":2,"_primary_term":3,"status":4,"error":{"type":"error_type","reason":"cause bad","index":"index","shard":"shard_named_this","index_uuid":"definitely_a_uuid"}}}`), //nolint:lll
			want: ActionResponse{
				Type:    "type",
				Index:   "index",
				ID:      "id",
				Version: 1,
				Result:  "success",
				Shards: &ShardMeta{
					Total:      1,
					Successful: 2,
					Skipped:    3,
					Failed:     4,
				},
				SeqNo:       2,
				PrimaryTerm: 3,
				Status:      4,
				Error: &ActionError{
					Type:      "error_type",
					Reason:    "cause bad",
					Index:     "index",
					Shard:     "shard_named_this",
					IndexUUID: "definitely_a_uuid",
				},
			},
			wantErr: false,
		},
		{
			name:      "Empty Response",
			jsonBytes: []byte(`{}`),
			wantErr:   true,
		},
		{
			name:      "Multiple action responses",
			jsonBytes: []byte(`{"bad":{},"type":{"_index":"index","_id":"id","_version":1,"result":"success","_shards":{"total":1,"successful":2,"skipped":3,"failed":4},"_seq_no":2,"_primary_term":3,"status":4,"error":{"type":"error_type","reason":"cause bad","index":"index","shard":"shard_named_this","index_uuid":"definitely_a_uuid"}}}`), //nolint:lll,
			wantErr:   true,
		},
		{
			name:      "Empty action response",
			jsonBytes: []byte(`{"action":{}}`),
			wantErr:   false,
			want:      ActionResponse{Type: "action"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ActionResponse
			gotErr := json.Unmarshal(tt.jsonBytes, &got)

			if tt.wantErr != (gotErr != nil) {
				t.Errorf("error wanted %v but got %v", tt.wantErr, gotErr)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestBulkAction_MarshalJSONLines_BulkCreateAction(t *testing.T) {
	tests := []struct {
		name    string
		doc     RoutableDoc
		want    [][]byte
		wantErr bool
	}{
		{
			name: "Simple Success",
			doc: bulkTestDoc{
				index:      "index",
				id:         "id",
				OtherField: 1,
			},
			want: [][]byte{
				[]byte(`{"create":{"_id":"id","_index":"index"}}`),
				[]byte(`{"other_field":1}`),
			},
			wantErr: false,
		},
		{
			name: "Valid missing index",
			doc: bulkTestDoc{
				id:         "id",
				OtherField: 1,
			},
			want: [][]byte{
				[]byte(`{"create":{"_id":"id"}}`),
				[]byte(`{"other_field":1}`),
			},
			wantErr: false,
		},
		{
			name:    "Nil doc fails",
			doc:     nil,
			wantErr: true,
		},
		{
			name:    "Doc missing ID fails",
			doc:     bulkTestDoc{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCreateBulkAction(tt.doc)

			jsonLines, err := c.MarshalJSONLines()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSONLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				require.Nil(t, jsonLines, "MarshalJSONLines should return nil if errored")
			} else {
				require.Lenf(t, jsonLines, len(tt.want), "wanted %d lines", len(tt.want))
				for i, want := range tt.want {
					require.JSONEq(t, string(want), string(jsonLines[i]))
				}
			}
		})
	}
}

func TestBulkAction_MarshalJSONLines_BulkIndexAction(t *testing.T) {
	tests := []struct {
		name    string
		doc     RoutableDoc
		want    [][]byte
		wantErr bool
	}{
		{
			name: "Simple Success",
			doc: bulkTestDoc{
				index:      "index",
				id:         "id",
				OtherField: 1,
			},
			want: [][]byte{
				[]byte(`{"index":{"_id":"id","_index":"index"}}`),
				[]byte(`{"other_field":1}`),
			},
			wantErr: false,
		},
		{
			name: "Valid missing index",
			doc: bulkTestDoc{
				id:         "id",
				OtherField: 1,
			},
			want: [][]byte{
				[]byte(`{"index":{"_id":"id"}}`),
				[]byte(`{"other_field":1}`),
			},
			wantErr: false,
		},
		{
			name:    "Nil doc fails",
			doc:     nil,
			wantErr: true,
		},
		{
			name:    "Doc missing ID fails",
			doc:     bulkTestDoc{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewIndexBulkAction(tt.doc)

			jsonLines, err := c.MarshalJSONLines()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSONLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				require.Nil(t, jsonLines, "MarshalJSONLines should return nil if errored")
			} else {
				require.Lenf(t, jsonLines, len(tt.want), "wanted %d lines", len(tt.want))
				for i, want := range tt.want {
					require.JSONEq(t, string(want), string(jsonLines[i]))
				}
			}
		})
	}
}

func TestBulkAction_MarshalJSONLines_BulkUpdateAction(t *testing.T) {
	tests := []struct {
		name    string
		doc     RoutableDoc
		want    [][]byte
		wantErr bool
	}{
		{
			name: "Simple Success",
			doc: bulkTestDoc{
				index:      "index",
				id:         "id",
				OtherField: 1,
			},
			want: [][]byte{
				[]byte(`{"update":{"_id":"id","_index":"index"}}`),
				[]byte(`{"other_field":1}`),
			},
			wantErr: false,
		},
		{
			name: "Valid missing index",
			doc: bulkTestDoc{
				id:         "id",
				OtherField: 1,
			},
			want: [][]byte{
				[]byte(`{"update":{"_id":"id"}}`),
				[]byte(`{"other_field":1}`),
			},
			wantErr: false,
		},
		{
			name:    "Nil doc fails",
			doc:     nil,
			wantErr: true,
		},
		{
			name:    "Doc missing ID fails",
			doc:     bulkTestDoc{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewUpdateBulkAction(tt.doc)

			jsonLines, err := c.MarshalJSONLines()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSONLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				require.Nil(t, jsonLines, "MarshalJSONLines should return nil if errored")
			} else {
				require.Lenf(t, jsonLines, len(tt.want), "wanted %d lines", len(tt.want))
				for i, want := range tt.want {
					require.JSONEq(t, string(want), string(jsonLines[i]))
				}
			}
		})
	}
}

func TestBulkAction_MarshalJSONLines_BulkDeleteAction(t *testing.T) {
	type fields struct {
		index string
		id    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "Simple Success",
			fields:  fields{index: "index", id: "id"},
			want:    []byte(`{"delete":{"_id":"id","_index":"index"}}`),
			wantErr: false,
		},
		{
			name:    "Empty Index",
			fields:  fields{id: "id"},
			want:    []byte(`{"delete":{"_id":"id"}}`),
			wantErr: false,
		},
		{
			name:    "empty id",
			fields:  fields{index: "index"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeleteBulkAction(tt.fields.index, tt.fields.id)

			got, err := d.MarshalJSONLines()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSONLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				require.Nil(t, got, "MarshalJSONLines should return nil if errored")
			} else {
				require.Lenf(t, got, 1, "BulkDelete should only return 1 json line")
				require.JSONEq(t, string(tt.want), string(got[0]))
			}
		})
	}
}
