package opensearchtools

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
