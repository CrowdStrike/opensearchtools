package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// Update encapsulates a bulk Update Action.
// This action updates existing documents and returns an error if the document doesn’t exist.
//
// For more details, see https://opensearch.org/docs/latest/api-reference/document-apis/bulk/
type Update struct {
	doc RoutableDoc
}

// NewUpdateAction instantiates an Update Action for the full or partial RoutableDoc
func NewUpdateAction(doc RoutableDoc) *Update {
	return &Update{doc: doc}
}

// GetAction returns a json byte array for the action and metadata of the operation
func (u *Update) GetAction() ([]byte, error) {
	if u.doc == nil {
		return nil, fmt.Errorf("nil document on Update Action")
	}

	action := map[string]any{
		"update": map[string]any{
			"_id":    u.doc.ID(),
			"_index": u.doc.Index(),
		},
	}

	return json.Marshal(action)
}

// GetDoc returns the byte array for the document in the operation
func (u *Update) GetDoc() ([]byte, error) {
	if u.doc == nil {
		return nil, fmt.Errorf("nil document on Update Action")
	}
	return json.Marshal(u.doc)
}
