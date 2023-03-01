package opensearchtools

import (
	"strings"
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
			validationResults: NewValidationResults(),
			want:              false,
		},
		{
			name: "non-empty, no fatals",
			validationResults: ValidationResultsFromSlice([]ValidationResult{
				{
					Message: "invalid field: foo",
					Fatal:   false,
				},
			}),
			want: false,
		},
		{
			name: "non-empty, one fatal",
			validationResults: ValidationResultsFromSlice([]ValidationResult{
				{
					Message: "invalid field: foo",
					Fatal:   true,
				},
			}),
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
			validationResults: NewValidationResults(),
			want:              "",
		},
		{
			name: "non-empty, no fatals",
			validationResults: ValidationResultsFromSlice([]ValidationResult{
				{
					Message: "invalid field: foo",
					Fatal:   false,
				},
			}),
			want: "One or more validations failed:\ninvalid field: foo\n",
		},
		{
			name: "one fatal result",
			validationResults: ValidationResultsFromSlice([]ValidationResult{
				{
					Message: "invalid field: foo",
					Fatal:   true,
				},
			}),
			want: "One or more validations failed:\nfatal: invalid field: foo\n",
		},
		{
			name: "multiple results",
			validationResults: ValidationResultsFromSlice([]ValidationResult{
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
			}),
			want: "One or more validations failed:\ninvalid field: test\nfatal: invalid field: foo\nfatal: missing field: bar\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validationError := NewValidationError(tt.validationResults)

			// check that the error contains the same messages, agnostic to order
			require.ElementsMatch(
				t,
				strings.Split(tt.want, "\n"),
				strings.Split(validationError.Error(), "\n"),
				"incorrect Error() result",
			)
		})
	}
}

func Test_ValidationResults_Add(t *testing.T) {
	vrs := NewValidationResults()
	vrs.Add(NewValidationResult("m1", false))
	require.Equal(
		t,
		ValidationResultsFromSlice([]ValidationResult{
			{
				Message: "m1",
				Fatal:   false,
			},
		}),
		vrs,
		"ValidationResult not added",
	)
	vrs.Add(NewValidationResult("m2", true))
	require.Equal(
		t,
		ValidationResultsFromSlice([]ValidationResult{
			{
				Message: "m1",
				Fatal:   false,
			},
			{
				Message: "m2",
				Fatal:   true,
			},
		}),
		vrs,
		"ValidationResult not added",
	)

	// now try to add a dupe - the collection should remain the same
	vrs.Add(NewValidationResult("m2", true))
	require.Equal(
		t,
		ValidationResultsFromSlice([]ValidationResult{
			{
				Message: "m1",
				Fatal:   false,
			},
			{
				Message: "m2",
				Fatal:   true,
			},
		}),
		vrs,
		"ValidationResults.Add added a duplicate (bad)",
	)
}

func Test_ValidationResults_Add_WithLiteralInstantiation(t *testing.T) {
	vrs := ValidationResults{}
	vrs.Add(NewValidationResult("m1", false))
	require.Equal(
		t,
		ValidationResultsFromSlice([]ValidationResult{
			{
				Message: "m1",
				Fatal:   false,
			},
		}),
		vrs,
		"ValidationResult not added",
	)
	vrs.Add(NewValidationResult("m2", true))
}

func Test_ValidationResults_Extend(t *testing.T) {
	vrs1 := ValidationResultsFromSlice([]ValidationResult{
		{
			Message: "msg1",
			Fatal:   false,
		},
		{
			Message: "msg2",
			Fatal:   true,
		},
	})

	vrs2 := ValidationResultsFromSlice([]ValidationResult{
		{
			Message: "msg3",
			Fatal:   false,
		},
		{
			Message: "msg4",
			Fatal:   true,
		},
	})

	vrs1.Extend(vrs2)

	require.Equal(
		t,
		ValidationResultsFromSlice([]ValidationResult{
			{
				Message: "msg1",
				Fatal:   false,
			},
			{
				Message: "msg2",
				Fatal:   true,
			},
			{
				Message: "msg3",
				Fatal:   false,
			},
			{
				Message: "msg4",
				Fatal:   true,
			},
		}),
		vrs1,
		"validation results not extended",
	)
}

func Test_ValidationResults_Extend_WithLiteralInstantiation(t *testing.T) {
	vrs1 := ValidationResults{}
	vrs2 := ValidationResultsFromSlice([]ValidationResult{
		{
			Message: "m1",
			Fatal:   false,
		},
	})
	vrs1.Extend(vrs2)

	require.Equal(
		t,
		ValidationResultsFromSlice([]ValidationResult{
			{
				Message: "m1",
				Fatal:   false,
			},
		}),
		vrs1,
		"ValidationResults not extended when initialized with literal init",
	)
}
