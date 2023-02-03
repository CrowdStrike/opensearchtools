package opensearchtools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

// BulkRequest wraps the functionality of [opensearchapi.Bulk] by supporting request body creation.
// An empty BulkRequest will fail to execute. At least one Action is required to be added.
//
// For more details see https://opensearch.org/docs/latest/api-reference/document-apis/bulk/
type BulkRequest struct {
	Actions []Action
	Refresh Refresh
}

// NewBulkRequest instantiates an empty BulkRequest
func NewBulkRequest() *BulkRequest {
	return &BulkRequest{}
}

// Add an action to the BulkRequest.
func (r *BulkRequest) Add(actions ...Action) {
	r.Actions = append(r.Actions, actions...)
}

// MarshalJSON marshals the BulkRequest into the JSON format expected by OpenSearch
func (r *BulkRequest) MarshalJSON() ([]byte, error) {
	if len(r.Actions) == 0 {
		return nil, fmt.Errorf("bulk request requires at least one action")
	}

	bodyBuf := new(bytes.Buffer)
	for _, op := range r.Actions {
		action, aErr := op.GetAction()
		if aErr != nil {
			return nil, aErr
		}

		bodyBuf.Write(action)
		bodyBuf.WriteRune('\n')

		doc, dErr := op.GetDoc()
		if dErr != nil {
			return nil, dErr
		}

		if len(doc) > 0 {
			bodyBuf.Write(doc)
			bodyBuf.WriteRune('\n')
		}
	}

	return bodyBuf.Bytes(), nil
}

// Do executes the BulkRequest using the provided opensearch.Client.
// If the request is executed successfully, then a BulkRequest will be returned.
// An error can be returned if
//
//   - Any Action is missing an action
//   - The call to OpenSearch fails
//   - The result json cannot be unmarshalled
func (r *BulkRequest) Do(ctx context.Context, client *opensearch.Client) (*BulkResponse, error) {
	rawBody, jErr := json.Marshal(r)
	if jErr != nil {
		return nil, jErr
	}

	osResp, rErr := opensearchapi.BulkRequest{
		Body:    bytes.NewReader(rawBody),
		Refresh: string(r.Refresh),
	}.Do(ctx, client)

	if rErr != nil {
		return nil, rErr
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	resp := &BulkResponse{
		StatusCode: osResp.StatusCode,
		Header:     osResp.Header,
	}

	if err := json.Unmarshal(respBuf.Bytes(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// BulkResponse wraps the functionality of [opensearchapi.Response] by unmarshalling the api response into
// a slice of [bulk.ActionResponse].
type BulkResponse struct {
	StatusCode int              `json:"-"`
	Header     http.Header      `json:"-"`
	Took       int64            `json:"took"`
	Errors     bool             `json:"errors"`
	Items      []ActionResponse `json:"items"`
	Error      *Error           `json:"error,omitempty"`
}
