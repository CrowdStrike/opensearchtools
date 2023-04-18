package opensearchtools

import (
	"context"
)

// Index defines a method which knows how to make an OpenSearch [Index] requests.
// It should be implemented by a version-specific executor.
//
// [Index]: https://opensearch.org/docs/latest/api-reference/index-apis/index/
type Index interface {
	CreateIndex(ctx context.Context, req *CreateIndexRequest) (OpenSearchResponse[CreateIndexResponse], error)
	DeleteIndex(ctx context.Context, req *DeleteIndexRequest) (OpenSearchResponse[DeleteIndexResponse], error)
	ExistIndex(ctx context.Context, req *ExistIndexRequest) (OpenSearchResponse[ExistIndexResponse], error)
	GetIndex(ctx context.Context, req *GetIndexRequest) (OpenSearchResponse[GetIndexRequest], error)
}
