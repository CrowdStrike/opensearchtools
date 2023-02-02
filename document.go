package opensearchtools

import "encoding/json"

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

// DocumentResult defines any OpenSearch response that contains a document source.
type DocumentResult interface {
	// GetSource returns the raw bytes of the document
	GetSource() []byte
}

// PtrTo is a generic constraint that restricts value to be pointers.
// T can be any type.
type PtrTo[T any] interface {
	*T
}

// ReadDocument reads the source from a DocumentResult and parses it into the passed document object.
// Document and be any pointer type.
func ReadDocument[D any, P PtrTo[D], R DocumentResult](docResult R, document P) error {
	return json.Unmarshal(docResult.GetSource(), document)
}
