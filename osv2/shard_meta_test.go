package osv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CrowdStrike/opensearchtools"
)

func TestShardMeta_ToModel(t *testing.T) {
	tests := []struct {
		name   string
		target ShardMeta
		want   opensearchtools.ShardMeta
	}{
		{
			name:   "Empty",
			target: ShardMeta{},
			want:   opensearchtools.ShardMeta{},
		},
		{
			name: "All fields",
			target: ShardMeta{
				Total:      10,
				Successful: 10,
				Skipped:    10,
				Failed:     10,
			},
			want: opensearchtools.ShardMeta{
				Total:      10,
				Successful: 10,
				Skipped:    10,
				Failed:     10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.ToModel()
			require.Equal(t, tt.want, got)
		})
	}
}
