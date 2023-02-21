package opensearchtools

// Refresh is a common query parameter across OpenSearch write requests
//
// More details can be found on https://opensearch.org/docs/latest/api-reference/document-apis/bulk/
type Refresh string

const (
	// True - make the write operation immediately available for search
	True Refresh = "true"
	// False - default, do not wait for the write operation to be available for search
	False Refresh = "false"
	// WaitFor - waits for the normal index refresh to make the write operation searchable
	WaitFor Refresh = "wait_for"
)
