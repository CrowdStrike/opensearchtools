package osv2

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CrowdStrike/opensearchtools"
)

func TestBulkRequest_ToOpenSearchJSON(t *testing.T) {
	testDoc := opensearchtools.NewDocumentRef("index", "id")

	testCreateAction := opensearchtools.NewCreateBulkAction(testDoc)
	createJSONLines, _ := testCreateAction.MarshalJSONLines()
	testCreateJSON := fmt.Sprintf("%s\n%s\n", createJSONLines[0], createJSONLines[1])

	testIndexAction := opensearchtools.NewIndexBulkAction(testDoc)
	indexJSONLines, _ := testIndexAction.MarshalJSONLines()
	testIndexJSON := fmt.Sprintf("%s\n%s\n", indexJSONLines[0], indexJSONLines[1])

	testUpdateAction := opensearchtools.NewUpdateBulkAction(testDoc)
	updateJSONLines, _ := testUpdateAction.MarshalJSONLines()
	testUpdateJSON := fmt.Sprintf("%s\n%s\n", updateJSONLines[0], updateJSONLines[1])

	testDeleteAction := opensearchtools.NewDeleteBulkAction(testDoc.Index(), testDoc.ID())
	deleteJSONLines, _ := testDeleteAction.MarshalJSONLines()
	testDeleteJSON := fmt.Sprintf("%s\n", deleteJSONLines[0])

	tests := []struct {
		name    string
		actions []opensearchtools.BulkAction
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Action list",
			wantErr: true,
		},
		{
			name:    "Single Create Action",
			actions: []opensearchtools.BulkAction{testCreateAction},
			want:    testCreateJSON,
		},
		{
			name:    "Single Update Action",
			actions: []opensearchtools.BulkAction{testUpdateAction},
			want:    testUpdateJSON,
		},
		{
			name:    "Single Index Action",
			actions: []opensearchtools.BulkAction{testIndexAction},
			want:    testIndexJSON,
		},
		{
			name:    "Single Delete Action",
			actions: []opensearchtools.BulkAction{testDeleteAction},
			want:    testDeleteJSON,
		},
		{
			name:    "MultipleActions",
			actions: []opensearchtools.BulkAction{testCreateAction, testUpdateAction, testIndexAction, testDeleteAction},
			want:    testCreateJSON + testUpdateJSON + testIndexJSON + testDeleteJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewBulkRequest()
			r.Add(tt.actions...)

			got, err := r.ToOpenSearchJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				require.Nil(t, got, "no json expected if there's an error")
			} else {
				require.Equal(t, tt.want, string(got))
			}
		})
	}
}

func TestBulkResponse_ToDomain1(t *testing.T) {
	tests := []struct {
		name   string
		target BulkResponse
		want   opensearchtools.BulkResponse
	}{
		{
			name:   "Empty",
			target: BulkResponse{},
			want:   opensearchtools.BulkResponse{},
		},
		{
			name: "Successful request",
			target: BulkResponse{
				Took:   10,
				Errors: true,
				Items:  []opensearchtools.ActionResponse{{Type: "test"}},
			},
			want: opensearchtools.BulkResponse{
				Took:   10,
				Errors: true,
				Items:  []opensearchtools.ActionResponse{{Type: "test"}},
			},
		},
		{
			name: "Unsuccessful request",
			target: BulkResponse{
				Error: &Error{
					Type:         "It was bad",
					Reason:       "for a reason",
					Index:        "on this index",
					ResourceID:   "with this document",
					ResourceType: "type document",
					IndexUUID:    "asdfasd",
				},
			},
			want: opensearchtools.BulkResponse{
				Error: &opensearchtools.Error{
					Type:         "It was bad",
					Reason:       "for a reason",
					Index:        "on this index",
					ResourceID:   "with this document",
					ResourceType: "type document",
					IndexUUID:    "asdfasd",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.toDomain()
			require.Equal(t, tt.want, got)
		})
	}
}
