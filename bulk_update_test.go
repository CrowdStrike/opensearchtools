package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdate_GetAction(t *testing.T) {
	tests := []struct {
		name    string
		doc     RoutableDoc
		want    []byte
		wantErr bool
	}{
		{
			name:    "Simple Success",
			doc:     NewDocumentRef("index", "id"),
			want:    []byte(`{"update":{"_id":"id","_index":"index"}}`),
			wantErr: false,
		},
		{
			name:    "Nil doc fails",
			doc:     nil,
			wantErr: true,
		},
		{
			name: "Extra fields don't affect action line",
			doc: bulkTestDoc{
				index:      "index",
				id:         "id",
				OtherField: 1,
			},
			want:    []byte(`{"update":{"_id":"id","_index":"index"}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewUpdateAction(tt.doc)

			got, err := c.GetAction()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				require.Nil(t, got, "GetAction should return nil if errored")
			} else {
				require.JSONEq(t, string(tt.want), string(got))
			}
		})
	}
}

func TestUpdate_GetDoc(t *testing.T) {
	tests := []struct {
		name    string
		doc     RoutableDoc
		want    []byte
		wantErr bool
	}{
		{
			name: "Simple Success",
			doc: bulkTestDoc{
				index:      "index",
				id:         "id",
				OtherField: 1,
			},
			want:    []byte(`{"other_field":1}`),
			wantErr: false,
		},
		{
			name:    "Nil doc fails",
			doc:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewUpdateAction(tt.doc)
			got, err := c.GetDoc()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				require.Nil(t, got, "GetDoc should return nil if errored")
			} else {
				require.JSONEq(t, string(tt.want), string(got))
			}
		})
	}
}
