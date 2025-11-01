package atproto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/think-root/bluesky-connector/internal/models"
)

const (
	CreateRecordEndpoint = "/xrpc/com.atproto.repo.createRecord"
	PostCollection       = "app.bsky.feed.post"
)

type RecordManager struct {
	baseURL        string
	httpClient     *http.Client
	sessionManager *SessionManager
}

func NewRecordManager(baseURL string, sessionManager *SessionManager) *RecordManager {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &RecordManager{
		baseURL:        baseURL,
		httpClient:     sessionManager.httpClient,
		sessionManager: sessionManager,
	}
}

func (rm *RecordManager) CreatePost(repo, text string, reply *models.Reply, embed *models.Embed) (*models.CreateRecordResponse, error) {
	if !rm.sessionManager.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	// Log token usage for debugging
	fmt.Printf("DEBUG: Creating post with token: %s...\n", rm.sessionManager.GetAccessToken()[:20])

	return rm.createPostWithRetry(repo, text, reply, embed, false)
}

func (rm *RecordManager) createPostWithRetry(repo, text string, reply *models.Reply, embed *models.Embed, isRetry bool) (*models.CreateRecordResponse, error) {
	// Log embed details if present
	if embed != nil && len(embed.Images) > 0 {
		fmt.Printf("DEBUG: Creating post with embed - Type: %s, Image MIME: %s\n",
			embed.Type, embed.Images[0].Image.MimeType)
	}

	// Detect hashtags and create facets
	facets := DetectHashtags(text)
	if len(facets) > 0 {
		fmt.Printf("DEBUG: Detected %d hashtag(s) in post\n", len(facets))
	}

	postRecord := models.PostRecord{
		Type:      "app.bsky.feed.post",
		Text:      text,
		CreatedAt: time.Now().UTC(),
		Reply:     reply,
		Embed:     embed,
		Facets:    facets,
	}

	reqBody := models.CreateRecordRequest{
		Repo:       repo,
		Collection: PostCollection,
		Record:     postRecord,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", rm.baseURL+CreateRecordEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+rm.sessionManager.GetAccessToken())

	fmt.Printf("DEBUG: Making create post request (retry=%v)\n", isRetry)

	resp, err := rm.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var atError models.ATProtoError
		if err := json.NewDecoder(resp.Body).Decode(&atError); err == nil {
			fmt.Printf("DEBUG: AT Protocol error: %s (retry=%v)\n", atError.String(), isRetry)

			// If token expired and this is not a retry, attempt to refresh and retry
			if atError.Error == "ExpiredToken" && !isRetry {
				fmt.Println("DEBUG: Token expired, attempting to refresh session...")
				if _, refreshErr := rm.sessionManager.RefreshSession(); refreshErr != nil {
					fmt.Printf("DEBUG: Failed to refresh session: %v\n", refreshErr)
					return nil, fmt.Errorf("AT Protocol error: %s (failed to refresh: %v)", atError.String(), refreshErr)
				}
				fmt.Println("DEBUG: Session refreshed successfully, retrying post creation...")
				return rm.createPostWithRetry(repo, text, reply, embed, true)
			}

			return nil, fmt.Errorf("AT Protocol error: %s", atError.String())
		}
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	fmt.Println("DEBUG: Post creation successful")

	var recordResp models.CreateRecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&recordResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &recordResp, nil
}

func (rm *RecordManager) CreateReply(repo, text string, root, parent *models.PostRef) (*models.CreateRecordResponse, error) {
	reply := &models.Reply{
		Root:   root,
		Parent: parent,
	}

	return rm.CreatePost(repo, text, reply, nil)
}

func (rm *RecordManager) CreatePostWithEmbed(repo, text string, embed *models.Embed) (*models.CreateRecordResponse, error) {
	return rm.CreatePost(repo, text, nil, embed)
}
