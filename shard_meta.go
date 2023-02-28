package opensearchtools

// ShardMeta contains information about the shards used or interacted with
// to perform a given OpenSearch Request.
type ShardMeta struct {
	Total      int
	Successful int
	Skipped    int
	Failed     int
}
