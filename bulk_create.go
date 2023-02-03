package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// Create encapsulates a bulk Create Action.
// Creates a document if it doesn’t already exist and returns an error otherwise.
//
// For more details, see https://opensearch.org/docs/latest/api-reference/document-apis/bulk/
type Create struct {
	doc RoutableDoc
}

// NewCreateAction instantiates a Create Action
func NewCreateAction(doc RoutableDoc) *Create {
	return &Create{doc: doc}
}

// GetAction returns a json byte array for the action and metadata of the operation
func (c *Create) GetAction() ([]byte, error) {
	if c.doc == nil {
		return nil, fmt.Errorf("nil document on Create Action")
	}

	action := map[string]any{
		"create": map[string]any{
			"_id":    c.doc.ID(),
			"_index": c.doc.Index(),
		},
	}

	return json.Marshal(action)
}

// GetDoc returns the byte array for the document in the operation
func (c *Create) GetDoc() ([]byte, error) {
	if c.doc == nil {
		return nil, fmt.Errorf("nil document on Create Action")
	}

	return json.Marshal(c.doc)
}
