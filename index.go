package opensearchtools

import (
	"context"
	"encoding/json"
)

// Index defines a method which knows how to make an OpenSearch [Index] request.
// It should be implemented by a version-specific executor.
//
// [Index]: https://opensearch.org/docs/latest/api-reference/index-apis/index/
type Index interface {
	Index(ctx context.Context, req *IndexRequest) (OpenSearchResponse[IndexResponse], error)
}

type IndexRequest struct {
	Action IndexAction

	Indices []string

	Refresh Refresh

	Routing string
}

type IndexAction struct {
	Type IndexActionType
	Doc  RoutableDoc
}

type IndexActionType string

const (
	// IndexCreate creates an index and returns an error otherwise.
	IndexCreate IndexActionType = "create"

	// IndexDelete delete an index if it exists and returns an error otherwise.
	IndexDelete IndexActionType = "delete"

	// IndexGet gets an index if it exists and returns an error otherwise.
	IndexGet IndexActionType = "update"

	// IndexExists checks the existence of an index
	IndexExists IndexActionType = "exists"
)

func NewCreateIndexAction(doc RoutableDoc) IndexAction {
	return IndexAction{
		Type: IndexCreate,
		Doc:  doc,
	}
}
func NewDeleteIndexAction() IndexAction {
	return IndexAction{
		Type: IndexDelete,
	}
}

func NewGetIndexAction() IndexAction {
	return IndexAction{
		Type: IndexGet,
	}
}

func NewExistsIndexAction() IndexAction {
	return IndexAction{
		Type: IndexExists,
	}
}

type IndexResponse struct {
	Acknowledged *bool
	Error        *Error
	*Indices
}

type Indices map[string]IndexInfo

type IndexInfo struct {
	aliases  map[string]json.RawMessage
	mappings map[string]json.RawMessage
	Settings *IndexSettings
}
type IndexSettings struct {
	CreationDate     string
	NumberOfShards   string
	NumberOfReplicas string
	UUID             string
	Version          struct{ Created string }
	ProvidedName     string
}
