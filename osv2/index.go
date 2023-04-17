package osv2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"github.com/CrowdStrike/opensearchtools"
)

type IndexRequest struct {
	Operation opensearchtools.IndexOperation
	Routing   string
}

// NewIndexRequest instantiates an empty IndexRequest
func NewIndexRequest() *IndexRequest {
	return &IndexRequest{}
}

// Add an index operation to the IndexRequest.
func (i *IndexRequest) Add(operation opensearchtools.IndexOperation) *IndexRequest {
	i.Operation = operation
	return i
}

func (i *IndexRequest) WithRouting(routing string) *IndexRequest {
	i.Routing = routing
	return i
}

// Validate the index request based on the types
func (i *IndexRequest) Validate() error {
	switch i.Operation.OperationType() {
	case opensearchtools.IndexCreate:
		if len(i.Operation.GetIndices()) > 1 {
			return fmt.Errorf("the index length for create should be 1, current lenght: %d", len(i.Operation.GetIndices()))
		}
	case opensearchtools.IndexGet:
		// implement some checks
	case opensearchtools.IndexExists:
		// implement some checks
	case opensearchtools.IndexDelete:
		// implement some checks
	}
	return nil
}

func (i *IndexRequest) Do(ctx context.Context, client *opensearch.Client) (*opensearchtools.OpenSearchResponse[IndexResponse], error) {
	if err := i.Validate(); err != nil {
		return nil, err
	}
	rawBody, jErr := i.ToOpenSearchJSON()
	if jErr != nil {
		return nil, jErr
	}
	var osResp *opensearchapi.Response
	var rErr error
	switch i.Operation.OperationType() {
	case opensearchtools.IndexCreate:
		osResp, rErr = opensearchapi.IndicesCreateRequest{
			Index: i.Operation.GetIndices()[0], // todo: check it please
			Body:  bytes.NewReader(rawBody),
		}.Do(ctx, client)
	case opensearchtools.IndexDelete:
		osResp, rErr = opensearchapi.IndicesDeleteRequest{
			Index: i.Operation.Indices,
		}.Do(ctx, client)
	case opensearchtools.IndexExists:
		osResp, rErr = opensearchapi.IndicesExistsRequest{
			Index: i.Operation.GetIndices(),
		}.Do(ctx, client)
	case opensearchtools.IndexGet:
		osResp, rErr = opensearchapi.IndicesGetRequest{
			Index: i.Operation.GetIndices(),
		}.Do(ctx, client)
	}

	if rErr != nil {
		return nil, rErr
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	resp, err := i.parseResponse(respBuf.Bytes(), i.Operation.OperationType())
	if err != nil {
		return nil, err
	}
	return &opensearchtools.OpenSearchResponse[IndexResponse]{
		StatusCode: osResp.StatusCode,
		Header:     osResp.Header,
		Response:   resp,
	}, nil
}

func (i *IndexRequest) ToOpenSearchJSON() ([]byte, error) {
	bodyBuf := new(bytes.Buffer)
	// parsing get the doc for index create etc
	return bodyBuf.Bytes(), nil
}

func (i *IndexRequest) parseResponse(resp []byte, operation opensearchtools.IndexOp) (IndexResponse, error) {
	indexResponse := IndexResponse{}
	if err := json.Unmarshal(resp, &indexResponse); err != nil {
		return IndexResponse{}, err
	}
	if operation == opensearchtools.IndexGet {
		if indexResponse.Error == nil {
			var indexInfo IndexGetResponse
			if err := json.Unmarshal(resp, &indexInfo); err != nil {
				return IndexResponse{}, err
			}
			return IndexResponse{IndexGetResponse: &indexInfo}, nil
		}
		return indexResponse, nil
	}

	return indexResponse, nil
}

type IndexResponse struct {
	Acknowledged *bool
	Error        *Error
	*IndexGetResponse
}

type IndexGetResponse map[string]IndexInfo

type IndexInfo struct {
	Aliases  map[string]json.RawMessage
	Mappings map[string]json.RawMessage
	Settings *IndexSettings
}

type IndexSettings struct {
	Index IndexSettingsInfo
}
type IndexSettingsInfo struct {
	CreationDate     string
	NumberOfShards   string
	NumberOfReplicas string
	UUID             string
	Version          struct{ Created string }
	ProvidedName     string
}
