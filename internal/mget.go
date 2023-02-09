package internal

import (
	"context"
	"github.com/CrowdStrike/opensearchtools"
)

// MGet defines a method which knows how to make an OpenSearch multiple document get.
// It should be implemented by a version-specific executor.
type MGet[K opensearchtools.RoutableDoc] interface {
	MGet(ctx context.Context, req *opensearchtools.MGetRequest) (*opensearchtools.MGetResponse, error)
}

// TODO: define other things you can do with an executor here
