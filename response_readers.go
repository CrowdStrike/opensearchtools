package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// ReadAggregationResult generically reads a sub bucket from a AggregationResultSet
// and parses it into the passed aggregation response.
// subAggResponse can be any pointer type.
func ReadAggregationResult[A any, P PtrTo[A], R AggregationResultSet](name string, aggResponse R, subAggResponse P) error {
	subAggSource, exists := aggResponse.GetAggregationResultSource(name)
	if !exists {
		return fmt.Errorf("no sub aggregation response with name %s", name)
	}

	return json.Unmarshal(subAggSource, subAggResponse)
}

// ReadDocument reads the source from a DocumentResult and parses it into the passed document object.
// Document can be any pointer type.
func ReadDocument[D any, P PtrTo[D], R DocumentResult](docResult R, document P) error {
	return json.Unmarshal(docResult.GetSource(), document)
}
