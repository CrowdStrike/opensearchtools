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

// BulkRequest is a serializable form of [opensearchtools.BulkRequest] specific to
// the [opensearchapi.BulkRequest] in OpenSearch V2.
// An empty BulkRequest will fail to execute. At least one Action is required to be added.
//
// For more details see https://opensearch.org/docs/latest/api-reference/document-apis/bulk/
type BulkRequest struct {
	// Actions lists the actions to be performed in the BulkRequest
	Actions []opensearchtools.BulkAction

	// Refresh determines if the request should wait for a refresh or not
	Refresh opensearchtools.Refresh

	// Index determines the entire index for the request
	Index string
}

// fromDomainBulkRequest creates a new [BulkRequest] from the given [opensearchtools.BulkRequest/.
func fromDomainBulkRequest(req *opensearchtools.BulkRequest) (BulkRequest, opensearchtools.ValidationResults) {
	// As more versions are implemented, these [opensearchtools.ValidationResults] may be used to contain issues
	// converting from the domain model to the V2 model.
	var vrs opensearchtools.ValidationResults

	return BulkRequest{
		Actions: req.Actions,
		Refresh: req.Refresh,
		Index:   req.Index,
	}, vrs
}

// Validate validates the given BulkRequest
func (r *BulkRequest) Validate() opensearchtools.ValidationResults {
	var validationResults opensearchtools.ValidationResults

	topLevelIndexIsEmpty := r.Index == ""
	for _, a := range r.Actions {
		// ensure Index is either set at the top level or set in each of the Actions
		if topLevelIndexIsEmpty && a.Doc.Index() == "" {
			validationResults.Add(opensearchtools.NewValidationResult(
				fmt.Sprintf("Index not set at the BulkRequest level nor in the Action %s with ID %s",
					a.Type, a.Doc.ID()), true))
		}

		// ensure that ID() is non-empty for Actions that require an ID
		if a.Doc.ID() == "" &&
			(a.Type == opensearchtools.BulkUpdate || a.Type == opensearchtools.BulkDelete) {
			validationResults.Add(opensearchtools.NewValidationResult("Doc ID is empty", true))
		}
	}

	return validationResults
}

// NewBulkRequest instantiates an empty BulkRequest
func NewBulkRequest() *BulkRequest {
	return &BulkRequest{}
}

// Add an action to the BulkRequest.
func (r *BulkRequest) Add(actions ...opensearchtools.BulkAction) *BulkRequest {
	r.Actions = append(r.Actions, actions...)
	return r
}

// WithIndex on the request
func (r *BulkRequest) WithIndex(index string) *BulkRequest {
	r.Index = index
	return r
}

// ToOpenSearchJSON marshals the BulkRequest into the JSON format expected by OpenSearch.
// Note: A BulkRequest is multi-line json with new line delimiters. It is not a singular valid json struct.
// For example:
//
//	{ action1 json }
//	{ action2 json }
func (r *BulkRequest) ToOpenSearchJSON() ([]byte, error) {
	if len(r.Actions) == 0 {
		return nil, fmt.Errorf("bulk request requires at least one action")
	}

	bodyBuf := new(bytes.Buffer)
	for _, op := range r.Actions {
		jsonLines, jErr := op.MarshalJSONLines()
		if jErr != nil {
			return nil, jErr
		}

		for _, line := range jsonLines {
			bodyBuf.Write(line)
			bodyBuf.WriteRune('\n')
		}
	}

	return bodyBuf.Bytes(), nil
}

// Do executes the [BulkRequest] using the provided opensearch.Client.
// If the request is executed successfully, then a [BulkResponse] will be returned.
// An error can be returned if
//
//   - Any Action is missing an action
//   - The call to OpenSearch fails
//   - The result json cannot be unmarshalled
func (r *BulkRequest) Do(ctx context.Context, client *opensearch.Client) (*opensearchtools.OpenSearchResponse[BulkResponse], error) {
	vrs := r.Validate()
	if vrs.IsFatal() {
		return nil, opensearchtools.NewValidationError(vrs)
	}

	rawBody, jErr := r.ToOpenSearchJSON()
	if jErr != nil {
		return nil, jErr
	}

	osResp, rErr := opensearchapi.BulkRequest{
		Body:    bytes.NewReader(rawBody),
		Refresh: string(r.Refresh),
		Index:   r.Index,
	}.Do(ctx, client)

	if rErr != nil {
		return nil, rErr
	}

	var respBuf bytes.Buffer
	if _, err := respBuf.ReadFrom(osResp.Body); err != nil {
		return nil, err
	}

	resp := BulkResponse{}

	if err := json.Unmarshal(respBuf.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &opensearchtools.OpenSearchResponse[BulkResponse]{
		StatusCode:        osResp.StatusCode,
		Header:            osResp.Header,
		Response:          resp,
		ValidationResults: vrs,
	}, nil
}

// BulkResponse wraps the functionality of [opensearchapi.Response] by unmarshalling the api response into
// a slice of [bulk.ActionResponse].
type BulkResponse struct {
	Took   int64                            `json:"took"`
	Errors bool                             `json:"errors"`
	Items  []opensearchtools.ActionResponse `json:"items"`
	Error  *Error                           `json:"error,omitempty"`
}

// toDomain converts this instance of [BulkResponse] into an [opensearchtools.BulkResponse]
func (b BulkResponse) toDomain() opensearchtools.BulkResponse {
	domainResp := opensearchtools.BulkResponse{
		Took:   b.Took,
		Errors: b.Errors,
		Items:  b.Items,
	}

	if b.Error != nil {
		domainErr := b.Error.toDomain()
		domainResp.Error = &domainErr
	}

	return domainResp
}
