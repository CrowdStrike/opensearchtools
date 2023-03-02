package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// BulkActionType is an enum for the various BulkActionTypes.
type BulkActionType string

const (
	// BulkCreate creates a document if it doesn’t already exist and returns an error otherwise.
	BulkCreate BulkActionType = "create"

	// BulkIndex creates a document if it doesn’t yet exist and replaces the document if it already exists.
	BulkIndex BulkActionType = "index"

	// BulkDelete deletes a document if it exists. If the document doesn’t exist,
	// OpenSearch doesn’t return an error, but instead returns not_found under ActionResponse.Result.
	BulkDelete BulkActionType = "delete"

	// BulkUpdate updates existing documents and returns an error if the document doesn’t exist.
	BulkUpdate BulkActionType = "update"
)

// BulkAction is a domain model union type for all actions a [BulkRequest] can perform across all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// While the struct has all the combined fields, only valid fields will be marshaled depending on the type of action.
// For this reason, it's recommended to use the type constructors:
//
//   - NewIndexBulkAction
//   - NewCreateBulkAction
//   - NewDeleteBulkAction
//   - NewUpdateBulkAction
//
// For more details, see https://opensearch.org/docs/latest/api-reference/document-apis/bulk/#request-body
type BulkAction struct {
	Type BulkActionType
	Doc  RoutableDoc
}

// NewCreateBulkAction instantiates a BulkCreate action.
func NewCreateBulkAction(doc RoutableDoc) BulkAction {
	return BulkAction{
		Type: BulkCreate,
		Doc:  doc,
	}
}

// NewIndexBulkAction instantiates a BulkIndex action.
func NewIndexBulkAction(doc RoutableDoc) BulkAction {
	return BulkAction{
		Type: BulkIndex,
		Doc:  doc,
	}
}

// NewUpdateBulkAction instantiates a BulkUpdate action.
func NewUpdateBulkAction(doc RoutableDoc) BulkAction {
	return BulkAction{
		Type: BulkUpdate,
		Doc:  doc,
	}
}

// NewDeleteBulkAction instantiates a BulkDelete action.
func NewDeleteBulkAction(index, id string) BulkAction {
	return BulkAction{
		Type: BulkDelete,
		Doc:  NewDocumentRef(index, id),
	}
}

// MarshalJSONLines marshals the BulkAction into the appropriate JSON lines depending on the BulkActionType.
func (b *BulkAction) MarshalJSONLines() ([][]byte, error) {
	if b.Doc == nil {
		return nil, fmt.Errorf("missing routing information on BulkAction %s", b.Type)
	}

	if b.Doc.ID() == "" {
		return nil, fmt.Errorf("missing id routing information on BulkAction %s", b.Type)
	}

	var jsonLines [][]byte

	actionRouting := map[string]any{"_id": b.Doc.ID()}

	// if Index is empty, use the request level index for routing
	if b.Doc.Index() != "" {
		actionRouting["_index"] = b.Doc.Index()
	}

	actionMeta := make(map[string]any)
	switch b.Type {
	case BulkCreate, BulkIndex, BulkUpdate:
		actionMeta[string(b.Type)] = actionRouting
		var (
			line []byte
			jErr error
		)

		if line, jErr = json.Marshal(actionMeta); jErr != nil {
			return nil, jErr
		}

		jsonLines = append(jsonLines, line)

		if line, jErr = json.Marshal(b.Doc); jErr != nil {
			return nil, jErr
		}

		jsonLines = append(jsonLines, line)
	case BulkDelete:
		actionMeta[string(b.Type)] = actionRouting
		if line, jErr := json.Marshal(actionMeta); jErr != nil {
			return nil, jErr
		} else {
			jsonLines = append(jsonLines, line)
		}
	default:
		return nil, fmt.Errorf("unssuported BulkActionType: %s", b.Type)
	}

	return jsonLines, nil
}

// ActionResponse is a domain model union type for all the fields of action responses for all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
type ActionResponse struct {
	Type        string       `json:"-"`
	Index       string       `json:"_index"`
	ID          string       `json:"_id"`
	Version     uint64       `json:"_version"`
	Result      string       `json:"result"`
	Shards      *ShardMeta   `json:"_shards,omitempty"`
	SeqNo       uint64       `json:"_seq_no"`
	PrimaryTerm uint64       `json:"_primary_term"`
	Status      int          `json:"status"`
	Error       *ActionError `json:"error"`
}

// UnmarshalJSON implements [json.Unmarshaler] to decode a json byte slice into an ActionResponse
func (o *ActionResponse) UnmarshalJSON(m []byte) error {
	// map[action type] -> map[response attribute] -> value
	var rawResp map[string]map[string]json.RawMessage
	if err := json.Unmarshal(m, &rawResp); err != nil {
		return err
	}

	if len(rawResp) > 1 {
		return fmt.Errorf("unexpected number of operation responses %d expected 1", len(rawResp))
	}

	for opType, attrMap := range rawResp {
		o.Type = opType

		if err := readField(&o.Index, "_index", attrMap); err != nil {
			return err
		}

		if err := readField(&o.ID, "_id", attrMap); err != nil {
			return err
		}

		if err := readField(&o.Version, "_version", attrMap); err != nil {
			return err
		}

		if err := readField(&o.Result, "result", attrMap); err != nil {
			return err
		}

		if err := readField(&o.Shards, "_shards", attrMap); err != nil {
			return err
		}

		if err := readField(&o.SeqNo, "_seq_no", attrMap); err != nil {
			return err
		}

		if err := readField(&o.PrimaryTerm, "_primary_term", attrMap); err != nil {
			return err
		}

		if err := readField(&o.Status, "status", attrMap); err != nil {
			return err
		}

		if err := readField(&o.Error, "error", attrMap); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("empty bulk action response returned")
}

func readField[T any, P PtrTo[T]](destination P, field string, attributes map[string]json.RawMessage) error {
	if value, exists := attributes[field]; exists {
		if err := json.Unmarshal(value, destination); err != nil {
			return err
		}
	}

	return nil
}

// ActionError encapsulates error responses from OpenSearch on bulk Actions
type ActionError struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	Index     string `json:"index"`
	Shard     string `json:"shard"`
	IndexUUID string `json:"index_uuid"`
}
