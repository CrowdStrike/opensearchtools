package opensearchtools

import "encoding/json"

// RoutableDoc interface defines an OpenSearch document that can be routed to a specific index.
// The most basic implementation is DocumentID.
type RoutableDoc interface {
    // ID returns the document ID
    ID() string
    // Index returns the index the document should be routed to
    Index() string
}

// DocumentID is a simple implementation of RoutableDoc. Identifying a document with it's basic components.
type DocumentID struct {
    id    string
    index string
}

// NewDocumentID constructs a DocumentID with the core two identifiers, ID and Index.
func NewDocumentID(id, index string) DocumentID {
    return DocumentID{
        id:    id,
        index: index,
    }
}

// ID returns the document ID
func (d DocumentID) ID() string {
    return d.id
}

// Index returns the index the document should be routed to
func (d DocumentID) Index() string {
    return d.index
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
