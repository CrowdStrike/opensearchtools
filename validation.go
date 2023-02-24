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

// ValidationResults is a slice of ValidationResult. It is defined as its own type for convenience and to be able
// to have methods on it.
type ValidationResults []ValidationResult

// NewValidationResult creates a new ValidationResult instance with the given message and fatal params.
func NewValidationResult(message string, fatal bool) ValidationResult {
	return ValidationResult{
		Message: message,
		Fatal:   fatal,
	}
}

// IsFatal returns true if there is one or more [ValidationResult] that is fatal.
func (v ValidationResults) IsFatal() bool {
	if len(v) == 0 {
		return false
	}

	for _, vr := range v {
		if vr.Fatal {
			return true
		}
	}

	return false
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
	if len(e.ValidationResults) == 0 {
		return ""
	}

	var b strings.Builder
	fmt.Fprintln(&b, "One or more validations failed:")
	for _, v := range e.ValidationResults {
		if v.Fatal {
			fmt.Fprintf(&b, "fatal: %s\n", v.Message)
		} else {
			fmt.Fprintln(&b, v.Message)
		}
	}
	return b.String()
}

type Validation[T any] struct {
	ValidationResults ValidationResults
	ValidatedRequest  *T
}
