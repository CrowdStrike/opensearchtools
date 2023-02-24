package opensearchtools

import "net/http"

type OpenSearchResponse[T any] struct {
	ValidationResults ValidationResults
	StatusCode        int
	Header            http.Header
	Response          *T
}
