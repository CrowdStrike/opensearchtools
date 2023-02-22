package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ValidationResults_IsFatal(t *testing.T) {
	tests := []struct {
		name              string
		validationResults ValidationResults
		want              bool
	}{
		{
			name:              "empty",
			validationResults: []ValidationResult{},
			want:              false,
		},
		{
			name: "non-empty, no fatals",
			validationResults: []ValidationResult{
				{
					Message: "invalid field: foo",
					Fatal:   false,
				},
			},
			want: false,
		},
		{
			name: "non-empty, one fatal",
			validationResults: []ValidationResult{
				{
					Message: "invalid field: foo",
					Fatal:   true,
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.validationResults.IsFatal(), "incorrect IsFatal() result")
		})
	}
}

func Test_ValidationResults_Error(t *testing.T) {
	tests := []struct {
		name              string
		validationResults ValidationResults
		want              string
	}{
		{
			name:              "empty",
			validationResults: []ValidationResult{},
			want:              "",
		},
		{
			name: "non-empty, no fatals",
			validationResults: []ValidationResult{
				{
					Message: "invalid field: foo",
					Fatal:   false,
				},
			},
			want: "One or more validations failed:\ninvalid field: foo\n",
		},
		{
			name: "one fatal result",
			validationResults: []ValidationResult{
				{
					Message: "invalid field: foo",
					Fatal:   true,
				},
			},
			want: "One or more validations failed:\nfatal: invalid field: foo\n",
		},
		{
			name: "multiple results",
			validationResults: []ValidationResult{
				{
					Message: "invalid field: test",
					Fatal:   false,
				},
				{
					Message: "invalid field: foo",
					Fatal:   true,
				},
				{
					Message: "missing field: bar",
					Fatal:   true,
				},
			},
			want: "One or more validations failed:\ninvalid field: test\nfatal: invalid field: foo\nfatal: missing field: bar\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.validationResults.Error(), "incorrect Error() result")
		})
	}
}
