package opensearchtools

import "net/http"

// OpenSearchResponse is a generic return type for an OpenSearch query which has meta fields common
// to all response types as well as a generic Response field abstract across all response types.
type OpenSearchResponse[T any] struct {
	ValidationResults ValidationResults
	StatusCode        int
	Header            http.Header
	Response          *T
}
