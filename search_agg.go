package opensearchtools

import (
	"encoding/json"
)

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

	// ConvertSubAggregations uses the provided AggregateVersionConverter to convert all sub aggregations and return
	// them as a map of agg name -> aggregation
	ConvertSubAggregations(converter AggregateVersionConverter) (map[string]Aggregation, error)
}

// AggregateVersionConverter takes in a domain model Aggregation and makes any modifications or conversions needed for
// a specific version of OpenSearch.
type AggregateVersionConverter func(Aggregation) (Aggregation, error)

// AggregationResultSet represents a collection of Aggregation responses. This result set exists in two places:
//
//   - [SearchResponse] for a [SearchRequest] that included aggregations
//   - [BucketAggregation]s that have added sub aggregations to include.
//
// The result set is characterized by the ability to contain multiple results keyed by the aggregation name.
type AggregationResultSet interface {
	// GetAggregationResultSource fetches the raw JSON source for the provided name.
	// Returns nil, false if no aggregation response with the name exists.
	GetAggregationResultSource(name string) ([]byte, bool)

	// Keys returns a slice of keys for the list of sub aggregation results
	Keys() []string
}

// ReadAggregationResult generically reads a sub bucket from a AggregationResultSet
// and parses it into the passed aggregation response. It returns an exists boolean, if the agg key is in the result set,
// and an error if it failed to be read.
// subAggResponse can be any pointer type.
func ReadAggregationResult[A any, P PtrTo[A], R AggregationResultSet](name string, aggResponse R, subAggResponse P) (error, bool) {
	subAggSource, exists := aggResponse.GetAggregationResultSource(name)
	if !exists {
		return nil, false
	}

	return json.Unmarshal(subAggSource, subAggResponse), true
}
