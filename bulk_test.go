package opensearchtools

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBulkRequest_MarshalJSON(t *testing.T) {
	testDoc := NewDocumentRef("index", "id")

	testCreateAction := NewCreateAction(testDoc)
	createActionStr, _ := testCreateAction.GetAction()
	createDocStr, _ := testCreateAction.GetDoc()
	testCreateJSON := fmt.Sprintf("%s\n%s\n", createActionStr, createDocStr)

	testIndexAction := NewIndexAction(testDoc)
	indexActionStr, _ := testIndexAction.GetAction()
	indexDocStr, _ := testIndexAction.GetDoc()
	testIndexJSON := fmt.Sprintf("%s\n%s\n", indexActionStr, indexDocStr)

	testUpdateAction := NewUpdateAction(testDoc)
	updateActionStr, _ := testUpdateAction.GetAction()
	updateDocStr, _ := testUpdateAction.GetDoc()
	testUpdateJSON := fmt.Sprintf("%s\n%s\n", updateActionStr, updateDocStr)

	testDeleteAction := NewDeleteAction(testDoc.Index(), testDoc.ID())
	deleteActionStr, _ := testDeleteAction.GetAction()
	testDeleteJSON := fmt.Sprintf("%s\n", deleteActionStr)

	tests := []struct {
		name    string
		actions []Action
		want    string
		wantErr bool
	}{
		{
			name:    "Empty Action list",
			wantErr: true,
		},
		{
			name:    "Single Create Action",
			actions: []Action{testCreateAction},
			want:    testCreateJSON,
		},
		{
			name:    "Single Update Action",
			actions: []Action{testUpdateAction},
			want:    testUpdateJSON,
		},
		{
			name:    "Single Index Action",
			actions: []Action{testIndexAction},
			want:    testIndexJSON,
		},
		{
			name:    "Single Delete Action",
			actions: []Action{testDeleteAction},
			want:    testDeleteJSON,
		},
		{
			name:    "MultipleActions",
			actions: []Action{testCreateAction, testUpdateAction, testIndexAction, testDeleteAction},
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
