package opensearchtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDelete_GetAction(t *testing.T) {
	type fields struct {
		index string
		id    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "Simple Success",
			fields:  fields{index: "index", id: "id"},
			want:    []byte(`{"delete":{"_id":"id","_index":"index"}}`),
			wantErr: false,
		},
		{
			name:    "empty index and id",
			fields:  fields{},
			want:    []byte(`{"delete":{"_id":"","_index":""}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeleteAction(tt.fields.index, tt.fields.id)

			got, err := d.GetAction()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.JSONEq(t, string(tt.want), string(got))
		})
	}
}

func TestDelete_GetDoc(t *testing.T) {
	type fields struct {
		index string
		id    string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Simple Success",
			fields: fields{index: "index", id: "id"},
		},
		{
			name:   "empty index and id",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeleteAction(tt.fields.index, tt.fields.id)
			got, err := d.GetDoc()
			require.Nil(t, got, "Delete action doesn't use a document, should always be nil")
			require.Nil(t, err, "Delete action doesn't use a document, no error should be returned")
		})
	}
}
