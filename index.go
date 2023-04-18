package opensearchtools

import (
	"context"
)

// Index defines a method which knows how to make an OpenSearch [Index] requests.
// It should be implemented by a version-specific executor.
//
// [Index]: https://opensearch.org/docs/latest/api-reference/index-apis/index/
type Index interface {
	CreateIndex(ctx context.Context, req *CreateIndexRequest) (OpenSearchResponse[CreateIndexResponse], error)
}

//
//type IndexRequest struct {
//	Operation IndexOperation
//	Routing   string
//}
//
//type IndexOp string
//
//const (
//	// IndexCreate creates an index and returns an error otherwise.
//	IndexCreate IndexOp = "create"
//
//	// IndexDelete delete an index if it exists and returns an error otherwise.
//	IndexDelete IndexOp = "delete"
//
//	// IndexGet gets an index if it exists and returns an error otherwise.
//	IndexGet IndexOp = "get"
//
//	// IndexExists checks the existence of an index
//	IndexExists IndexOp = "exists"
//)
//
//// IndexOperation various different index operations to be executed
//type IndexOperation struct {
//	Operation IndexOp
//	Indices   []string
//	Config    *indexRequestConfig
//	Doc       RoutableDoc
//}
//
//func NewIndexCreateOperation(index string, doc RoutableDoc, ops ...IndexReqOptionFunc) *IndexOperation {
//	config := &indexRequestConfig{}
//	for _, o := range ops {
//		o(config)
//	}
//	return &IndexOperation{Indices: []string{index}, Doc: doc, Config: config, Operation: IndexCreate}
//}
//
//func NewIndexDeleteOperation(index []string, ops ...IndexReqOptionFunc) *IndexOperation {
//	config := &indexRequestConfig{}
//	for _, o := range ops {
//		o(config)
//	}
//	return &IndexOperation{Indices: index, Config: config, Operation: IndexCreate}
//}
//
//func NewIndexGetOperation(index []string, ops ...IndexReqOptionFunc) *IndexOperation {
//	config := &indexRequestConfig{}
//	for _, o := range ops {
//		o(config)
//	}
//	return &IndexOperation{Indices: index, Config: config, Operation: IndexGet}
//}
//
//func NewIndexExistOperation(index []string, ops ...IndexReqOptionFunc) *IndexOperation {
//	config := &indexRequestConfig{}
//	for _, o := range ops {
//		o(config)
//	}
//	return &IndexOperation{Indices: index, Operation: IndexExists}
//}
//func (i *IndexOperation) OperationType() IndexOp {
//	return i.Operation
//}
//
//func (i *IndexOperation) GetIndices() []string {
//	return i.Indices
//}
//
//func (i *IndexOperation) GetDoc() RoutableDoc {
//	return i.Doc
//}
//
//type indexRequestConfig struct {
//	masterTimeout         time.Duration
//	clusterManagerTimeout time.Duration
//	expandWildcards       string // make this enum later
//	waitForActiveShards   string
//	allowNoIndices        bool
//}
//
//// IndexReqOptionFunc for functional options
//type IndexReqOptionFunc func(config *indexRequestConfig)
//
//// WithMasterTimeOut to add the master timeout
//func WithMasterTimeOut(d time.Duration) IndexReqOptionFunc {
//	return func(config *indexRequestConfig) {
//		config.masterTimeout = d
//	}
//}
//
////// WithClusterManagerTimeout to add the cluster master timeout
////func WithClusterManagerTimeout(d time.Duration) IndexReqOptionFunc {
////	return func(config *indexRequestConfig) {
////		config.clusterManagerTimeout = d
////	}
////}
//
//// WithExpandWildcards to add the wildcard option
//func WithExpandWildcards(w string) IndexReqOptionFunc {
//	return func(config *indexRequestConfig) {
//		config.expandWildcards = w
//	}
//}
//
//// WithWaitForActiveShards to add the active shard option
//func WithWaitForActiveShards(s string) IndexReqOptionFunc {
//	return func(config *indexRequestConfig) {
//		config.waitForActiveShards = s
//	}
//}
//
//// WithAllowNoIndices to set allow no index option
//func WithAllowNoIndices(a bool) IndexReqOptionFunc {
//	return func(config *indexRequestConfig) {
//		config.allowNoIndices = a
//	}
//}
//
//// Validate validating things
//func (c *indexRequestConfig) Validate() error {
//	// implement some validation e.g. for active shard a positive value etc.
//	return nil
//}
//
//// IndexResponse - the combined response, may be create an interface and has each response its own sturct
//type IndexResponse struct {
//	// index create response
//	Acknowledged *bool
//	// general error response
//	Error *Error
//	// index get response
//	*IndexGetResponse
//	// index exits has only http status responses:200 – the index exists, and 404 – the index does not exist.
//	// index create has only http status responses: 201 – the index created, and 409 – the index already exists - 400 - bad index request
//}
//
//// IndexGetResponse the weired index get response with index name as keys
//type IndexGetResponse map[string]IndexInfo
//
//type IndexInfo struct {
//	Aliases  map[string]json.RawMessage
//	Mappings map[string]json.RawMessage
//	Settings *IndexSettings
//}
//type IndexSettings struct {
//	Index IndexSettingsInfo
//}
//type IndexSettingsInfo struct {
//	CreationDate     string
//	NumberOfShards   string
//	NumberOfReplicas string
//	UUID             string
//	Version          struct{ Created string }
//	ProvidedName     string
//}
