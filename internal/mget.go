package internal

import (
	"context"
	"github.com/CrowdStrike/opensearchtools"
)

// MGet defines a method which knows how to make an OpenSearch multiple document get.
// It should be implemented by a version-specific executor.
type MGet interface {
	MGet(ctx context.Context, req *opensearchtools.MGetRequest) (*opensearchtools.MGetResponse, error)
}
