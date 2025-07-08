package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlobRef_GetRefString(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected string
	}{
		{
			name:     "String ref",
			jsonData: `{"$type":"blob","ref":"bafkreibjfgx2gprinfvicegelk5kosd6y2frmqpqzwqkg7usac74l3t2v4","mimeType":"image/png","size":12345}`,
			expected: "bafkreibjfgx2gprinfvicegelk5kosd6y2frmqpqzwqkg7usac74l3t2v4",
		},
		{
			name:     "Object ref with $link",
			jsonData: `{"$type":"blob","ref":{"$link":"bafkreibjfgx2gprinfvicegelk5kosd6y2frmqpqzwqkg7usac74l3t2v4"},"mimeType":"image/png","size":12345}`,
			expected: "bafkreibjfgx2gprinfvicegelk5kosd6y2frmqpqzwqkg7usac74l3t2v4",
		},
		{
			name:     "Object ref with other string field",
			jsonData: `{"$type":"blob","ref":{"cid":"bafkreibjfgx2gprinfvicegelk5kosd6y2frmqpqzwqkg7usac74l3t2v4"},"mimeType":"image/png","size":12345}`,
			expected: "bafkreibjfgx2gprinfvicegelk5kosd6y2frmqpqzwqkg7usac74l3t2v4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var blobRef BlobRef
			err := json.Unmarshal([]byte(tt.jsonData), &blobRef)
			assert.NoError(t, err)
			
			result := blobRef.GetRefString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBlobRef_SetRefString(t *testing.T) {
	blobRef := &BlobRef{}
	testRef := "bafkreibjfgx2gprinfvicegelk5kosd6y2frmqpqzwqkg7usac74l3t2v4"
	
	blobRef.SetRefString(testRef)
	
	assert.Equal(t, testRef, blobRef.GetRefString())
}

func TestBlobUploadResponse_UnmarshalJSON(t *testing.T) {
	// Test that BlobUploadResponse can handle object ref format
	jsonData := `{"blob":{"$type":"blob","ref":{"$link":"bafkreibjfgx2gprinfvicegelk5kosd6y2frmqpqzwqkg7usac74l3t2v4"},"mimeType":"image/png","size":12345}}`
	
	var response BlobUploadResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "blob", response.Blob.Type)
	assert.Equal(t, "bafkreibjfgx2gprinfvicegelk5kosd6y2frmqpqzwqkg7usac74l3t2v4", response.Blob.GetRefString())
	assert.Equal(t, "image/png", response.Blob.MimeType)
	assert.Equal(t, int64(12345), response.Blob.Size)
}