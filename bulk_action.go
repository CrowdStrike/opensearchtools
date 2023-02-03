package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// Action wraps individual actions on a BulkRequest.
type Action interface {
	// GetAction returns a json byte array for the action and metadata of the operation.
	GetAction() ([]byte, error)

	// GetDoc returns the optional byte array for the document in the operation.
	// returns nil if no document is needed
	GetDoc() ([]byte, error)
}

// ActionResponse encapsulates the individual response for each operation
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
