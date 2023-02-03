package opensearchtools

import "encoding/json"

// Delete encapsulates a bulk Delete Action.
// This action deletes a document if it exists. If the document doesn’t exist,
// OpenSearch doesn’t return an error, but instead returns not_found under result.
//
// For more details, see https://opensearch.org/docs/latest/api-reference/document-apis/bulk/
type Delete struct {
	index string
	id    string
}

// NewDeleteAction instantiates a Delete Action on the index targeting the document with id
func NewDeleteAction(index, id string) *Delete {
	return &Delete{
		index: index,
		id:    id,
	}
}

// GetAction returns a json byte array for the action and metadata of the operation.
func (d *Delete) GetAction() ([]byte, error) {
	action := map[string]any{
		"delete": map[string]any{
			"_id":    d.id,
			"_index": d.index,
		},
	}

	return json.Marshal(action)
}

// GetDoc returns nil, a document is not needed for a Delete Action
func (d *Delete) GetDoc() ([]byte, error) {
	return nil, nil
}
