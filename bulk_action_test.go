package opensearchtools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

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
