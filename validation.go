package opensearchtools

import (
	"fmt"
	"strings"
)

// Validator defines the behavior for validating a domain model object.
type Validator interface {
	Validate() ValidationResults
}

// ValidationResult is an individual validation result item to be returned by a [ValidateForVersion] call.
type ValidationResult struct {
	Message string
	Fatal   bool
}

// NewValidationResult creates a new ValidationResult instance with the given message and fatal params.
func NewValidationResult(message string, fatal bool) ValidationResult {
	return ValidationResult{
		Message: message,
		Fatal:   fatal,
	}
}

// A collection of ValidationResults that implements set logic, meaning it holds unique entries only (no duplicates)
type ValidationResults struct {
	validationResults map[ValidationResult]struct{}
}

// NewValidationResults creates a new ValidationResults instance
func NewValidationResults() ValidationResults {
	return ValidationResults{
		validationResults: map[ValidationResult]struct{}{},
	}
}

// IsFatal returns true if there is one or more [ValidationResult] that is fatal.
func (vrs *ValidationResults) IsFatal() bool {
	if len(vrs.validationResults) == 0 {
		return false
	}

	for k := range vrs.validationResults {
		if k.Fatal {
			return true
		}
	}

	return false
}

// Add either adds the given [ValidationResult] to the set or does not if it already exists in the set
func (vrs *ValidationResults) Add(vr ValidationResult) {
	vrs.validationResults[vr] = struct{}{}
}

// Add either adds the given [ValidationResult] to the set or does not if it already exists in the set
func (vrs *ValidationResults) Extend(other ValidationResults) {
	for k := range other.validationResults {
		vrs.Add(k)
	}
}

// ValidationError is an error which contains ValidationResults
type ValidationError struct {
	ValidationResults ValidationResults
}

// NewValidationError creates a new ValidationError instance with the given ValidationResults
func NewValidationError(rs ValidationResults) *ValidationError {
	return &ValidationError{
		ValidationResults: rs,
	}
}

// Error returns a newline-separated string representation of all validation results or an empty string if there are none.
// Fatal results are prefixed with `fatal:`.
func (e *ValidationError) Error() string {
	if len(e.ValidationResults.validationResults) == 0 {
		return ""
	}

	var b strings.Builder
	fmt.Fprintln(&b, "One or more validations failed:")
	for vr := range e.ValidationResults.validationResults {
		if vr.Fatal {
			fmt.Fprintf(&b, "fatal: %s\n", vr.Message)
		} else {
			fmt.Fprintln(&b, vr.Message)
		}
	}
	return b.String()
}

// ValidationResultsFromSlice creates a ValidationResults from the given slice of [ValidationResult]
func ValidationResultsFromSlice(vrs []ValidationResult) ValidationResults {
	validationResults := NewValidationResults()
	for _, vr := range vrs {
		validationResults.Add(vr)
	}
	return validationResults
}
