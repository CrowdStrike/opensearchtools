package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
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

func TestMGetRequest_Add(t *testing.T) {
	type args struct {
		index string
		id    string
	}
	tests := []struct {
		name string
		args args
		want *MGetRequest
	}{
		{
			name: "add simple test",
			args: args{
				index: testIndex1,
				id:    testID1,
			},
			want: &MGetRequest{
				Docs: []RoutableDoc{
					DocumentRef{
						index: testIndex1,
						id:    testID1,
					},
				},
			},
		},
		{
			name: "add without index",
			args: args{
				id: testID1,
			},
			want: &MGetRequest{
				Docs: []RoutableDoc{
					DocumentRef{
						id: testID1,
					},
				},
			},
		},
		{
			name: "add without id",
			args: args{
				index: testIndex1,
			},
			want: &MGetRequest{
				Docs: []RoutableDoc{
					DocumentRef{
						index: testIndex1,
					},
				},
			},
		},
		{
			name: "add empty string",
			args: args{
				index: "",
				id:    "",
			},
			want: &MGetRequest{
				Docs: []RoutableDoc{
					DocumentRef{
						index: "",
						id:    "",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMGetRequest().Add(tt.args.index, tt.args.id)
			require.Len(t, m.Docs, 1, "MGetRequest.Add should only add a single document request")
			wantDoc := tt.want.Docs[0]
			gotDoc := m.Docs[0]

			require.Equal(t, wantDoc.ID(), gotDoc.ID(), "incorrect document ID after Add")
			require.Equal(t, wantDoc.Index(), gotDoc.Index(), "incorrect document Index after Add")
		})
	}
}

func TestMGetRequest_AddDocs(t *testing.T) {
	// Expected id and index from the RoutableDocs on the MGetRequest
	type mockDoc struct {
		id, index string
	}

	tests := []struct {
		name string
		docs []RoutableDoc
		want []mockDoc
	}{
		{
			name: "Single doc of single type",
			docs: []RoutableDoc{
				NewDocumentRef(testIndex1, testID1),
			},
			want: []mockDoc{
				{id: testID1, index: testIndex1},
			},
		},
		{
			name: "Multiple docs of single type",
			docs: []RoutableDoc{
				NewDocumentRef(testIndex1, testID1),
				NewDocumentRef(testIndex2, testID2),
			},
			want: []mockDoc{
				{id: testID1, index: testIndex1},
				{id: testID2, index: testIndex2},
			},
		},
		{
			name: "Multiple docs of mixed types",
			docs: []RoutableDoc{
				NewDocumentRef(testIndex1, testID1),
				mgetTestDoc{id: testID2, index: testIndex2},
			},
			want: []mockDoc{
				{id: testID1, index: testIndex1},
				{id: testID2, index: testIndex2},
			},
		},
		{
			name: "Single document, no index",
			docs: []RoutableDoc{
				NewDocumentRef("", testID1),
			},
			want: []mockDoc{
				{id: testID1, index: ""},
			},
		},
		{
			name: "Single document, no ID",
			docs: []RoutableDoc{
				NewDocumentRef(testIndex1, ""),
			},
			want: []mockDoc{
				{id: "", index: testIndex1},
			},
		},
		{
			name: "Single document, no ID or Index",
			docs: []RoutableDoc{
				NewDocumentRef("", ""),
			},
			want: []mockDoc{
				{id: "", index: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMGetRequest().AddDocs(tt.docs...)

			require.Len(t, m.Docs, len(tt.want), "unexpected number of documents added to the request")

			for i, gotDoc := range m.Docs {
				wantDoc := tt.want[i]

				require.Equal(t, wantDoc.id, gotDoc.ID())
				require.Equal(t, wantDoc.index, gotDoc.Index())
			}
		})
	}
}

func TestMGetRequest_MarshalJSON(t *testing.T) {
	type fields struct {
		Index string
		Docs  []RoutableDoc
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
				Docs:  []RoutableDoc{},
			},
			want:    []byte(`{"docs":[]}`),
			wantErr: false,
		},
		{
			name: "Single document single type",
			fields: fields{
				Docs: []RoutableDoc{
					mgetTestDoc{id: testID1, index: testIndex1},
				},
			},
			want:    []byte(`{"docs":[{"_id":"test_id","_index":"test_index"}]}`),
			wantErr: false,
		},
		{
			name: "Multiple documents single type",
			fields: fields{
				Docs: []RoutableDoc{
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
				Docs: []RoutableDoc{
					mgetTestDoc{id: testID1, index: testIndex1},
					NewDocumentRef(testIndex2, testID2),
				},
			},
			want:    []byte(`{"docs":[{"_id":"test_id","_index":"test_index"},{"_id":"test_id2","_index":"test_index2"}]}`),
			wantErr: false,
		},
		{
			name: "Document without index",
			fields: fields{
				Docs: []RoutableDoc{
					mgetTestDoc{id: testID1},
				},
			},
			want:    []byte(`{"docs":[{"_id":"test_id"}]}`),
			wantErr: false,
		},
		{
			name: "Document without id",
			fields: fields{
				Docs: []RoutableDoc{
					mgetTestDoc{index: testIndex1},
				},
			},
			want:    []byte(`{"docs":[{"_id":"","_index":"test_index"}]}`),
			wantErr: false,
		},
		{
			name: "Document without id and index",
			fields: fields{
				Docs: []RoutableDoc{
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
				Docs: []RoutableDoc{
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
