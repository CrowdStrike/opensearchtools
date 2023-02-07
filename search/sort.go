package search

import "encoding/json"

// Sort encapsulates the sort capabilities for OpenSearch.
// An empty Sort will be rejected by OpenSearch as a field must be non-null and non-empty.
//
// For more details see https://opensearch.org/docs/latest/opensearch/search/sort/
type Sort struct {
	Field string
	Desc  bool
}

// NewSort instantiates a search Sort with the field to be sorted and whether is descending or ascending.
func NewSort(field string, desc bool) *Sort {
	return &Sort{
		Field: field,
		Desc:  desc,
	}
}

// ToOpenSearchJSON converts the Sort to the correct OpenSearch JSON.
func (s *Sort) ToOpenSearchJSON() ([]byte, error) {
	sort := make(map[string]any)
	if s.Desc {
		sort["order"] = "desc"
	} else {
		sort["order"] = "asc"
	}

	source := map[string]any{
		s.Field: sort,
	}

	return json.Marshal(source)
}
