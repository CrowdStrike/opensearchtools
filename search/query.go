package search

// Query wraps all query types in a common interface. Facilitating a common pattern to convert the logical
// struct into the appropriate request JSON.
type Query interface {
	// ToOpenSearchJSON converts the Query struct to the expected OpenSearch JSON
	ToOpenSearchJSON() ([]byte, error)
}
