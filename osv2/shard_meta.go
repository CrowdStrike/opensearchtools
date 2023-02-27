package osv2

import "github.com/CrowdStrike/opensearchtools"

// ShardMeta contains information about the shards used or interacted with
// to perform a given OpenSearch Request.
type ShardMeta struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

// toDomain converts this instance of an ShardMeta into an [opensearchtools.ShardMeta]
func (s *ShardMeta) toDomain() opensearchtools.ShardMeta {
	return opensearchtools.ShardMeta{
		Total:      s.Total,
		Successful: s.Successful,
		Skipped:    s.Skipped,
		Failed:     s.Failed,
	}
}
