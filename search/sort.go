package search

// Sort encapsulates the sort capabilities for OpenSearch
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

// Source converts the Sort to the correct OpenSearch JSON.
func (s *Sort) Source() any {
	sort := make(map[string]any)
	if s.Desc {
		sort["order"] = "desc"
	} else {
		sort["order"] = "asc"
	}

	source := make(map[string]any)
	source[s.Field] = sort

	return source
}
