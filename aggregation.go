package search

// Aggregation wraps all aggregation types into a common interface.
// Facilitating adding aggregations to [opensearchtools.SearchRequests] and marshaling into OpenSearch JSON.
type Aggregation interface {
	// ToOpenSearchJSON converts the Aggregation struct to the expected OpenSearch JSON
	ToOpenSearchJSON() ([]byte, error)
}

// BucketAggregation represents a family of OpenSearch aggregations.
// Bucket aggregations categorize sets of documents as buckets.
// The type of bucket aggregation determines whether a given document falls into a bucket or not.
// Bucket aggregations also support adding nested aggregations to further refine bucket results.
//
// For more details, see https://opensearch.org/docs/latest/opensearch/bucket-agg/
type BucketAggregation interface {
	Aggregation
	// AddSubAggregation to the BucketAggregation for further refinement.
	AddSubAggregation(name string, agg Aggregation) BucketAggregation
}
