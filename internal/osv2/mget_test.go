package osv2

import (
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
			m := &opensearchtools.MGetRequest{
				Index: tt.fields.Index,
				Docs:  tt.fields.Docs,
			}
			marshalableMGetRequest := FromModelMGetRequest(m)
			got, err := marshalableMGetRequest.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.JSONEq(t, string(tt.want), string(got))
		})
	}
}
