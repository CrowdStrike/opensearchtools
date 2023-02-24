package osv2

import "github.com/CrowdStrike/opensearchtools"

// Error encapsulates an error response from a given OpenSearch request.
type Error struct {
	RootCause    []*Error `json:"root_cause"`
	Type         string   `json:"type"`
	Reason       string   `json:"reason"`
	Index        string   `json:"index"`
	ResourceID   string   `json:"resource.id"`
	ResourceType string   `json:"resource.type"`
	IndexUUID    string   `json:"index_uuid"`
}

// ToModel converts this instance of an Error into an [opensearchtools.Error]
func (e *Error) ToModel() *opensearchtools.Error {
	if e == nil {
		return nil
	}

	var modelRootCauses []*opensearchtools.Error

	for _, specErr := range e.RootCause {
		modelRootCauses = append(modelRootCauses, specErr.ToModel())
	}

	return &opensearchtools.Error{
		RootCause:    modelRootCauses,
		Type:         e.Type,
		Reason:       e.Reason,
		Index:        e.Index,
		ResourceID:   e.ResourceID,
		ResourceType: e.ResourceType,
		IndexUUID:    e.IndexUUID,
	}
}
