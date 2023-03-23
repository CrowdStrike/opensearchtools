package opensearchtools

import (
	"encoding/json"
	"fmt"
)

// NestedQuery is a type of joining query that allows searches in fields that are of the `nested` type.
// An empty NestedQuery will be rejected by OpenSearch for two reasons:
//
//   - a path must not be nil or empty
//   - a query must not be nil
//
// For more details see https://opensearch.org/docs/latest/query-dsl/
type NestedQuery struct {
	path  string
	query Query
}

// NewNestedQuery initializes a NestedQuery targeting the nested field at the given path, with the provided query.
func NewNestedQuery(path string, query Query) *NestedQuery {
	return &NestedQuery{
		path:  path,
		query: query,
	}
}

// ToOpenSearchJSON converts the Nested to the correct OpenSearch JSON.
func (q *NestedQuery) ToOpenSearchJSON() ([]byte, error) {
	var (
		nestedQuery json.RawMessage
		nestedErr   error
	)

	if q.path == "" {
		return nil, fmt.Errorf("missing required nested path")
	}

	if q.query == nil {
		return nil, fmt.Errorf("missing required nested query")
	}

	nestedQuery, nestedErr = q.query.ToOpenSearchJSON()
	if nestedErr != nil {
		return nil, nestedErr
	}

	source := map[string]any{
		"nested": map[string]any{
			"path":  q.path,
			"query": nestedQuery,
		},
	}

	return json.Marshal(source)
}
