package opensearchtools

import (
	"context"
)

// Bulk defines a method which knows how to make an OpenSearch [Bulk] request.
// It should be implemented by a version-specific executor.
//
// [Bulk]: https://opensearch.org/docs/latest/api-reference/document-apis/bulk/
type Bulk interface {
	Bulk(ctx context.Context, req *BulkRequest) (OpenSearchResponse[BulkResponse], error)
}

// BulkRequest is a domain model union type for all the fields of BulkRequests for all
// supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// An empty BulkRequest will fail to execute. At least one Action is required to be added.
type BulkRequest struct {
	// Actions lists the actions to be performed in the BulkRequest
	Actions []BulkAction

	// Refresh determines if the request should wait for a refresh or not
	Refresh Refresh

	// Index determines the entire index for the request
	Index string
}

// NewBulkRequest instantiates an empty BulkRequest
func NewBulkRequest() *BulkRequest {
	return &BulkRequest{}
}

// Add an action to the BulkRequest.
func (r *BulkRequest) Add(actions ...BulkAction) *BulkRequest {
	r.Actions = append(r.Actions, actions...)
	return r
}

// WithIndex on the request
func (r *BulkRequest) WithIndex(index string) *BulkRequest {
	r.Index = index
	return r
}

// BulkResponse is a domain model union response type for BulkRequest for all supported OpenSearch versions.
// Currently supported versions are:
//   - OpenSearch 2
//
// Contains a slice of [ActionResponse] for each individual [BulkAction] performed by the request.
type BulkResponse struct {
	Took   int64
	Errors bool
	Items  []ActionResponse
	Error  *Error
}
