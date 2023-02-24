package opensearchtools

type OpenSearchResponse[T any] struct {
	ValidationResults ValidationResults
	Response          *T
}
