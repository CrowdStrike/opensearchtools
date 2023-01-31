package opensearchtools

// RoutableDoc interface defines an OpenSearch document that can be routed to a specific index.
// Most documents will route to a single index. A basic implementation might look like:
//
//	type BasicDoc struct {
//		ID 		string
//		Index 	string
//	}
//
//	func (b *BasicDoc) GetID() string {
//		return b.ID
//	}
//
//	func (b *BasicDoc) RouteToIndex() string{
//		return b.Index
//	}
type RoutableDoc interface {
	// GetID returns the document ID
	GetID() string
	// RouteToIndex returns the index the document should be routed to
	RouteToIndex() string
}
