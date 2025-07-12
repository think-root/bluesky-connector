package atproto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/think-root/bluesky-connector/internal/models"
)

const (
	BlobUploadEndpoint = "/xrpc/com.atproto.repo.uploadBlob"
)

type MediaManager struct {
	baseURL        string
	httpClient     *http.Client
	sessionManager *SessionManager
}

func NewMediaManager(baseURL string, sessionManager *SessionManager) *MediaManager {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &MediaManager{
		baseURL:        baseURL,
		httpClient:     sessionManager.httpClient,
		sessionManager: sessionManager,
	}
}

func (mm *MediaManager) UploadBlob(data []byte, mimeType string) (*models.BlobRef, error) {
	if !mm.sessionManager.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	fmt.Printf("DEBUG: Starting blob upload with token: %s...\n", mm.sessionManager.GetAccessToken()[:20])

	return mm.uploadBlobWithRetry(data, mimeType, false)
}

func (mm *MediaManager) uploadBlobWithRetry(data []byte, mimeType string, isRetry bool) (*models.BlobRef, error) {
	req, err := http.NewRequest("POST", mm.baseURL+BlobUploadEndpoint, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("Authorization", "Bearer "+mm.sessionManager.GetAccessToken())

	fmt.Printf("DEBUG: Making blob upload request (retry=%v)\n", isRetry)

	resp, err := mm.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var atError models.ATProtoError
		if err := json.Unmarshal(body, &atError); err == nil {
			fmt.Printf("DEBUG: AT Protocol error: %s (retry=%v)\n", atError.String(), isRetry)

			if atError.Error == "ExpiredToken" && !isRetry {
				fmt.Println("DEBUG: Token expired, attempting to refresh session...")
				if _, refreshErr := mm.sessionManager.RefreshSession(); refreshErr != nil {
					fmt.Printf("DEBUG: Failed to refresh session: %v\n", refreshErr)
					return nil, fmt.Errorf("AT Protocol error: %s (failed to refresh: %v)", atError.String(), refreshErr)
				}
				fmt.Println("DEBUG: Session refreshed successfully, retrying upload...")
				return mm.uploadBlobWithRetry(data, mimeType, true)
			}

			return nil, fmt.Errorf("AT Protocol error: %s", atError.String())
		}
		return nil, fmt.Errorf("HTTP error: %d, body: %s", resp.StatusCode, string(body))
	}

	fmt.Println("DEBUG: Blob upload successful")

	var uploadResp models.BlobUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("DEBUG: Server returned BlobRef - Type: %s, MimeType: %s, Size: %d\n",
		uploadResp.Blob.Type, uploadResp.Blob.MimeType, uploadResp.Blob.Size)

	fmt.Printf("DEBUG: Server returned MIME type: %s\n", uploadResp.Blob.MimeType)
	
	return &uploadResp.Blob, nil
}

func (mm *MediaManager) CreateImageEmbed(imageData []byte, mimeType, altText string) (*models.Embed, error) {
	fmt.Printf("DEBUG: CreateImageEmbed called with mimeType: %s\n", mimeType)

	blobRef, err := mm.UploadBlob(imageData, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload blob: %w", err)
	}

	fmt.Printf("DEBUG: BlobRef after upload - Type: %s, MimeType: %s, Size: %d\n",
		blobRef.Type, blobRef.MimeType, blobRef.Size)

	if blobRef.Type == "" {
		blobRef.Type = "blob"
	}

	embed := &models.Embed{
		Type: "app.bsky.embed.images",
		Images: []models.EmbedImage{
			{
				Alt:   altText,
				Image: blobRef,
			},
		},
	}

	fmt.Printf("DEBUG: Created embed with image MIME type: %s\n", embed.Images[0].Image.MimeType)
	return embed, nil
}

func DetectMimeType(data []byte) string {
	if len(data) < 4 {
		return "application/octet-stream"
	}

	// Check for common image formats
	switch {
	case bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}):
		return "image/jpeg"
	case bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}):
		return "image/png"
	case bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46}):
		return "image/gif"
	case bytes.HasPrefix(data, []byte{0x52, 0x49, 0x46, 0x46}) && bytes.Contains(data[8:12], []byte("WEBP")):
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
