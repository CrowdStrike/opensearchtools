package opensearchtools

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBulkRequest_MarshalJSON(t *testing.T) {
	testDoc := NewDocumentRef("index", "id")

	testCreateAction := NewCreateBulkAction(testDoc)
	createJSONLines, _ := testCreateAction.MarshalJSONLines()
	testCreateJSON := fmt.Sprintf("%s\n%s\n", createJSONLines[0], createJSONLines[1])

	testIndexAction := NewIndexBulkAction(testDoc)
	indexJSONLines, _ := testIndexAction.MarshalJSONLines()
	testIndexJSON := fmt.Sprintf("%s\n%s\n", indexJSONLines[0], indexJSONLines[1])

	testUpdateAction := NewUpdateBulkAction(testDoc)
	updateJSONLines, _ := testUpdateAction.MarshalJSONLines()
	testUpdateJSON := fmt.Sprintf("%s\n%s\n", updateJSONLines[0], updateJSONLines[1])

	testDeleteAction := NewDeleteBulkAction(testDoc.Index(), testDoc.ID())
	deleteJSONLines, _ := testDeleteAction.MarshalJSONLines()
	testDeleteJSON := fmt.Sprintf("%s\n", deleteJSONLines[0])

	tests := []struct {
		name    string
		actions []BulkAction
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Action list",
			wantErr: true,
		},
		{
			name:    "Single Create Action",
			actions: []BulkAction{testCreateAction},
			want:    testCreateJSON,
		},
		{
			name:    "Single Update Action",
			actions: []BulkAction{testUpdateAction},
			want:    testUpdateJSON,
		},
		{
			name:    "Single Index Action",
			actions: []BulkAction{testIndexAction},
			want:    testIndexJSON,
		},
		{
			name:    "Single Delete Action",
			actions: []BulkAction{testDeleteAction},
			want:    testDeleteJSON,
		},
		{
			name:    "MultipleActions",
			actions: []BulkAction{testCreateAction, testUpdateAction, testIndexAction, testDeleteAction},
			want:    testCreateJSON + testUpdateJSON + testIndexJSON + testDeleteJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewBulkRequest()
			r.Add(tt.actions...)

			got, err := r.MarshalJSON()
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
