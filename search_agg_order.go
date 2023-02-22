package opensearchtools

import (
	"fmt"
)

// Order encapsulates the sorting capabilities for [Aggregation] requests to OpenSearch.
// An empty Order will be rejected by OpenSearch as the target must be non-null and non-empty
type Order struct {
	Target string
	Desc   bool
}

// NewOrder instantiates an aggregation Order with the target and whether it should be descending or ascending.
func NewOrder(field string, desc bool) Order {
	return Order{
		Target: field,
		Desc:   desc,
	}
}

// ToOpenSearchJSON converts the Order to the correct OpenSearch JSON.
func (o Order) ToOpenSearchJSON() ([]byte, error) {
	dir := "asc"
	if o.Desc {
		dir = "desc"
	}

	return []byte(fmt.Sprintf("{%q: %q}", o.Target, dir)), nil
}
