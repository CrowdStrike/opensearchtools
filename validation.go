package opensearchtools

import (
	"fmt"
	"strings"
)

// Version is a string representation of an OpenSearch major version that this library supports.
type Version string

const (
	V2 Version = "V2"
)

// RequestVersionValidator defines the behavior for validating a model object for a specific OpenSearch version implementation.
type RequestVersionValidator interface {
	ValidateForVersion(v Version) (ValidationResults, error)
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

// Error returns a newline-separated string representation of all validation results or an empty string if there are none.
// Fatal results are prefixed with `fatal:`.
func (vs ValidationResults) Error() string {
	if len(vs) == 0 {
		return ""
	}

	var b strings.Builder
	fmt.Fprintln(&b, "One or more validations failed:")
	for _, v := range vs {
		if v.Fatal {
			fmt.Fprintf(&b, "fatal: %s\n", v.Message)
		} else {
			fmt.Fprintln(&b, v.Message)
		}
	}
	return b.String()
}
