package osv2

import "github.com/CrowdStrike/opensearchtools"

// Error encapsulates an error response from a given OpenSearch request.
type Error struct {
	RootCause    []Error `json:"root_cause"`
	Type         string  `json:"type"`
	Reason       string  `json:"reason"`
	Index        string  `json:"index"`
	ResourceID   string  `json:"resource.id"`
	ResourceType string  `json:"resource.type"`
	IndexUUID    string  `json:"index_uuid"`
}

// ToDomain converts this instance of an Error into an [opensearchtools.Error]
func (e *Error) ToDomain() opensearchtools.Error {
	var modelRootCauses []opensearchtools.Error

	for _, specErr := range e.RootCause {
		modelRootCauses = append(modelRootCauses, specErr.ToDomain())
	}

	return opensearchtools.Error{
		RootCause:    modelRootCauses,
		Type:         e.Type,
		Reason:       e.Reason,
		Index:        e.Index,
		ResourceID:   e.ResourceID,
		ResourceType: e.ResourceType,
		IndexUUID:    e.IndexUUID,
	}
}
