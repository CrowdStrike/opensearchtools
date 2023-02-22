package opensearchtools

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validateMGetRequestForV2(t *testing.T) {
	tests := []struct {
		name        string
		mgetRequest MGetRequest
		want        ValidationResults
	}{
		{
			name: "valid MGetRequest",
			mgetRequest: MGetRequest{
				Index: testIndex1,
				Docs: []RoutableDoc{
					DocumentRef{
						index: testIndex1,
						id:    testID1,
					},
				},
			},
			want: ValidationResults{},
		},
		{
			name: "missing Index",
			mgetRequest: MGetRequest{
				Index: "",
				Docs: []RoutableDoc{
					DocumentRef{
						index: "",
						id:    testID1,
					},
				},
			},
			want: []ValidationResult{
				{
					Message: fmt.Sprintf("Index not set at the MGetRequest level nor in the Doc with ID %s", testID1),
					Fatal:   true,
				},
			},
		},
		{
			name: "empty ID",
			mgetRequest: MGetRequest{
				Index: testIndex1,
				Docs: []RoutableDoc{
					DocumentRef{
						index: testIndex1,
						id:    "",
					},
				},
			},
			want: []ValidationResult{
				{
					Message: "Doc ID is empty",
					Fatal:   true,
				},
			},
		},
	}

	for _, tt := range tests {
		v := validateMGetRequestForV2(tt.mgetRequest)
		require.Equal(t, tt.want, v, "invalid ValidationResults")
	}
}
