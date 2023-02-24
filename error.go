package opensearchtools

// Error encapsulates an error response from a given OpenSearch request.
type Error struct {
	RootCause    []Error
	Type         string
	Reason       string
	Index        string
	ResourceID   string
	ResourceType string
	IndexUUID    string
}
