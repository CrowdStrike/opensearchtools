package opensearchtools

import "encoding/json"

// RoutableDoc interface defines an OpenSearch document that can be routed to a specific index.
// The most basic implementation is DocumentRef.
type RoutableDoc interface {
	// Index returns the index the document should be routed to
	Index() string
	// ID returns the document ID
	ID() string
}

// DocumentRef references a document via its index and id. It is the most basic implementation of RoutableDoc
type DocumentRef struct {
	index string
	id    string
}

// NewDocumentRef constructs a DocumentRef with the core two identifiers, ID and Index.
func NewDocumentRef(index, id string) DocumentRef {
	return DocumentRef{
		id:    id,
		index: index,
	}
}

// Index returns the index the document should be routed to
func (d DocumentRef) Index() string {
	return d.index
}

// ID returns the document ID
func (d DocumentRef) ID() string {
	return d.id
}

// DocumentResult defines any OpenSearch response that contains a document source.
type DocumentResult interface {
	// GetSource returns the raw bytes of the document
	GetSource() []byte
}

// ReadDocument reads the source from a DocumentResult and parses it into the passed document object.
// Document can be any pointer type.
func ReadDocument[D any, P PtrTo[D], R DocumentResult](docResult R, document P) error {
	return json.Unmarshal(docResult.GetSource(), document)
}
