package osv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CrowdStrike/opensearchtools"
)

func TestError_ToDomain(t *testing.T) {
	tests := []struct {
		name   string
		target Error
		want   opensearchtools.Error
	}{
		{
			name:   "Empty",
			target: Error{},
			want:   opensearchtools.Error{},
		},
		{
			name: "All fields, no root cause",
			target: Error{
				Type:         "Type",
				Reason:       "Reason",
				Index:        "Index",
				ResourceID:   "ResourceID",
				ResourceType: "ResourceType",
				IndexUUID:    "IndexUUID",
			},
			want: opensearchtools.Error{
				Type:         "Type",
				Reason:       "Reason",
				Index:        "Index",
				ResourceID:   "ResourceID",
				ResourceType: "ResourceType",
				IndexUUID:    "IndexUUID",
			},
		},
		{
			name: "Nested Root Causes",
			target: Error{
				RootCause: []Error{
					{
						RootCause: []Error{
							{
								Reason: "final nest",
							},
						},
						Reason: "nest 1",
					},
				},
				Reason: "Top Level",
			},
			want: opensearchtools.Error{
				RootCause: []opensearchtools.Error{
					{
						RootCause: []opensearchtools.Error{
							{
								Reason: "final nest",
							},
						},
						Reason: "nest 1",
					},
				},
				Reason: "Top Level",
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
