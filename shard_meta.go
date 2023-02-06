package opensearchtools

// ShardMeta contains information about the shards used or interacted with
// to perform a given OpenSearch Request.
type ShardMeta struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}
