package osv2

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CrowdStrike/opensearchtools"
)

const (
	testIndex1 = "test_index"
	testIndex2 = "test_index2"

	testID1 = "test_id"
	testID2 = "test_id2"
)

type mgetTestDoc struct {
	index, id string
}

func (d mgetTestDoc) Index() string {
	return d.index
}

func (d mgetTestDoc) ID() string {
	return d.id
}

func TestMGetRequest_MarshalJSON(t *testing.T) {
	type fields struct {
		Index string
		Docs  []opensearchtools.RoutableDoc
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Empty Request",
			fields: fields{
				Index: "",
				Docs:  []opensearchtools.RoutableDoc{},
			},
			want:    []byte(`{"docs":[]}`),
			wantErr: false,
		},
		{
			name: "Single document single type",
			fields: fields{
				Docs: []opensearchtools.RoutableDoc{
					mgetTestDoc{id: testID1, index: testIndex1},
				},
			},
			want:    []byte(`{"docs":[{"_id":"test_id","_index":"test_index"}]}`),
			wantErr: false,
		},
		{
			name: "Multiple documents single type",
			fields: fields{
				Docs: []opensearchtools.RoutableDoc{
					mgetTestDoc{id: testID1, index: testIndex1},
					mgetTestDoc{id: testID2, index: testIndex2},
				},
			},
			want:    []byte(`{"docs":[{"_id":"test_id","_index":"test_index"},{"_id":"test_id2","_index":"test_index2"}]}`),
			wantErr: false,
		},
		{
			name: "Multiple documents mixed type",
			fields: fields{
				Docs: []opensearchtools.RoutableDoc{
					mgetTestDoc{id: testID1, index: testIndex1},
					opensearchtools.NewDocumentRef(testIndex2, testID2),
				},
			},
			want:    []byte(`{"docs":[{"_id":"test_id","_index":"test_index"},{"_id":"test_id2","_index":"test_index2"}]}`),
			wantErr: false,
		},
		{
			name: "Document without index",
			fields: fields{
				Docs: []opensearchtools.RoutableDoc{
					mgetTestDoc{id: testID1},
				},
			},
			want:    []byte(`{"docs":[{"_id":"test_id"}]}`),
			wantErr: false,
		},
		{
			name: "Document without id",
			fields: fields{
				Docs: []opensearchtools.RoutableDoc{
					mgetTestDoc{index: testIndex1},
				},
			},
			want:    []byte(`{"docs":[{"_id":"","_index":"test_index"}]}`),
			wantErr: false,
		},
		{
			name: "Document without id and index",
			fields: fields{
				Docs: []opensearchtools.RoutableDoc{
					mgetTestDoc{},
				},
			},
			want:    []byte(`{"docs":[{"_id":""}]}`),
			wantErr: false,
		},
		{
			name: "Request level index does not affect request json body",
			fields: fields{
				Index: testIndex2,
				Docs: []opensearchtools.RoutableDoc{
					mgetTestDoc{id: testID1, index: testIndex1},
				},
			},
			want:    []byte(`{"docs":[{"_id":"test_id","_index":"test_index"}]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MGetRequest{
				Index: tt.fields.Index,
				Docs:  tt.fields.Docs,
			}
			got, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.JSONEq(t, string(tt.want), string(got))
		})
	}
}

func Test_MGetResult_toDomain(t *testing.T) {
	tests := []struct {
		name              string
		marshalableResult MGetResult
		want              opensearchtools.MGetResult
	}{
		{
			name: "Non-error result",
			marshalableResult: MGetResult{
				Index:       testIndex1,
				ID:          testID1,
				Version:     42,
				SeqNo:       99,
				PrimaryTerm: 10,
				Found:       true,
				Source:      []byte(`{"name": "bob", "age": 42}`),
				Error:       nil,
			},
			want: opensearchtools.MGetResult{
				Index:       testIndex1,
				ID:          testID1,
				Version:     42,
				SeqNo:       99,
				PrimaryTerm: 10,
				Found:       true,
				Source:      []byte(`{"name": "bob", "age": 42}`),
				Error:       nil,
			},
		},
		{
			name: "Error result",
			marshalableResult: MGetResult{
				Error: fmt.Errorf("some OpenSearch error"),
			},
			want: opensearchtools.MGetResult{
				Error: fmt.Errorf("some OpenSearch error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.marshalableResult.toDomain()
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_MGetResponse_toDomain(t *testing.T) {
	testHeaders := http.Header{}
	testHeaders.Add("x-foo", "bar")

	tests := []struct {
		name               string
		marshlableResponse MGetResponse
		want               *opensearchtools.OpenSearchResponse[opensearchtools.MGetResponse]
	}{
		{
			name: "Multiple docs returned",
			marshlableResponse: MGetResponse{
				StatusCode: 200,
				Header:     testHeaders,
				Docs: []MGetResult{
					{
						Index:       testIndex1,
						ID:          testID1,
						Version:     42,
						SeqNo:       99,
						PrimaryTerm: 10,
						Found:       true,
						Source:      []byte(`{"name": "bob", "age": 42}`),
						Error:       nil,
					},
					{
						Index:       testIndex2,
						ID:          testID2,
						Version:     1,
						SeqNo:       2,
						PrimaryTerm: 2,
						Found:       true,
						Source:      []byte(`{"deviceName": "abc123", "os": "windows"}`),
						Error:       nil,
					},
					{
						Index:       testIndex2,
						ID:          testID2,
						Version:     10,
						SeqNo:       220,
						PrimaryTerm: 30,
						Found:       false,
						Source:      []byte{},
						Error:       nil,
					},
				},
			},
			want: &opensearchtools.OpenSearchResponse[opensearchtools.MGetResponse]{
				ValidationResults: nil,
				StatusCode:        200,
				Header:            testHeaders,
				Response: &opensearchtools.MGetResponse{
					Docs: []opensearchtools.MGetResult{
						{
							Index:       testIndex1,
							ID:          testID1,
							Version:     42,
							SeqNo:       99,
							PrimaryTerm: 10,
							Found:       true,
							Source:      []byte(`{"name": "bob", "age": 42}`),
							Error:       nil,
						},
						{
							Index:       testIndex2,
							ID:          testID2,
							Version:     1,
							SeqNo:       2,
							PrimaryTerm: 2,
							Found:       true,
							Source:      []byte(`{"deviceName": "abc123", "os": "windows"}`),
							Error:       nil,
						},
						{
							Index:       testIndex2,
							ID:          testID2,
							Version:     10,
							SeqNo:       220,
							PrimaryTerm: 30,
							Found:       false,
							Source:      []byte{},
							Error:       nil,
						},
					},
				},
			},
		},
		{
			name: "No docs returned",
			marshlableResponse: MGetResponse{
				StatusCode: 200,
				Header:     testHeaders,
				Docs:       []MGetResult{},
			},
			want: &opensearchtools.OpenSearchResponse[opensearchtools.MGetResponse]{
				ValidationResults: nil,
				StatusCode:        200,
				Header:            testHeaders,
				Response: &opensearchtools.MGetResponse{
					Docs: []opensearchtools.MGetResult{},
				},
			},
		},
	}

	var vrs opensearchtools.ValidationResults

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.marshlableResponse.toDomain(vrs)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_MGetRequest_Validate(t *testing.T) {
	tests := []struct {
		name        string
		mgetRequest MGetRequest
		want        opensearchtools.ValidationResults
	}{
		{
			name: "valid MGetRequest",
			mgetRequest: MGetRequest{
				Index: testIndex1,
				Docs: []opensearchtools.RoutableDoc{
					opensearchtools.NewDocumentRef(testIndex1, testID1),
				},
			},
			want: nil,
		},
		{
			name: "Doc with no ID",
			mgetRequest: MGetRequest{
				Index: testIndex1,
				Docs: []opensearchtools.RoutableDoc{
					opensearchtools.NewDocumentRef("", ""),
				},
			},
			want: opensearchtools.ValidationResults{
				opensearchtools.ValidationResult{
					Message: "Doc ID is empty",
					Fatal:   true,
				},
			},
		},
		{
			name: "missing index",
			mgetRequest: MGetRequest{
				Index: "",
				Docs: []opensearchtools.RoutableDoc{
					opensearchtools.NewDocumentRef("", testID1),
				},
			},
			want: opensearchtools.ValidationResults{
				opensearchtools.ValidationResult{
					Message: fmt.Sprintf("Index not set at the MGetRequest level nor in the Doc with ID %s", testID1),
					Fatal:   true,
				},
			},
		},
	}

	for _, tt := range tests {
		v := tt.mgetRequest.Validate()
		require.Equal(t, tt.want, v, "invalid validation result")
	}
}
