package opensearchtools

// RoutableDoc interface defines an OpenSearch document that can be routed to a specific index
type RoutableDoc interface {
	// GetID returns the document ID
	GetID() string
	// RouteToIndex returns the index the document should be routed to
	RouteToIndex() string
}
