package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// Index encapsulates a bulk Index Action.
// Index actions create a document if it doesn’t yet exist and replace the document if it already exists.
//
// For more details, see https://opensearch.org/docs/latest/api-reference/document-apis/bulk/
type Index struct {
	doc RoutableDoc
}

// NewIndexAction instantiates a Index Action
func NewIndexAction(doc RoutableDoc) *Index {
	return &Index{doc: doc}
}

// GetAction returns a json byte array for the action and metadata of the operation
func (c *Index) GetAction() ([]byte, error) {
	if c.doc == nil {
		return nil, fmt.Errorf("nil document on Index Action")
	}

	action := map[string]any{
		"index": map[string]any{
			"_id":    c.doc.ID(),
			"_index": c.doc.Index(),
		},
	}

	return json.Marshal(action)
}

// GetDoc returns the byte array for the document in the operation
func (c *Index) GetDoc() ([]byte, error) {
	if c.doc == nil {
		return nil, fmt.Errorf("nil document on Index Action")
	}

	return json.Marshal(c.doc)
}
