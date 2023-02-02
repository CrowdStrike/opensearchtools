package opensearchtools

// RoutableDoc interface defines an OpenSearch document that can be routed to a specific index.
// Most documents will route to a single index. A basic implementation might look like:
//
//	type BasicDoc struct {
//		ID 		string
//		Index 	string
//	}
//
//	func (b *BasicDoc) ID() string {
//		return b.ID
//	}
//
//	func (b *BasicDoc) Index() string{
//		return b.Index
//	}
type RoutableDoc interface {
	// ID returns the document ID
	ID() string
	// Index returns the index the document should be routed to
	Index() string
}
