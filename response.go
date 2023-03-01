package opensearchtools

import "net/http"

// OpenSearchResponse is a generic return type for an OpenSearch query which has meta fields common
// to all response types as well as a generic Response field abstract across all response types.
type OpenSearchResponse[T any] struct {
	ValidationResults ValidationResults
	StatusCode        int
	Header            http.Header
	Response          T
}

// Craete a new OpenSearchResponse instance with the given [ValidationResults], status code, headers, and response.
func NewOpenSearchResponse[T any](vrs ValidationResults, statusCode int, header http.Header, response T) OpenSearchResponse[T] {
	return OpenSearchResponse[T]{
		ValidationResults: vrs,
		StatusCode:        statusCode,
		Header:            header,
		Response:          response,
	}
}

// OpenSearchRequest is a generic request type for an OpenSearch which encapsulates the underlying
// version-specific request and validation results.
type OpenSearchRequest[T any] struct {
	ValidationResults ValidationResults
	Request           T
}

// Craete a new OpenSearchRequest instance with the given [ValidationResults] and request.
func NewOpenSearchRequest[T any](vrs ValidationResults, request T) OpenSearchRequest[T] {
	return OpenSearchRequest[T]{
		ValidationResults: vrs,
		Request:           request,
	}
}
