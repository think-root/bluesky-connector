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
	PostCollection = "app.bsky.feed.post"
)

type RecordManager struct {
	baseURL       string
	httpClient    *http.Client
	sessionManager *SessionManager
}

func NewRecordManager(baseURL string, sessionManager *SessionManager) *RecordManager {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	
	return &RecordManager{
		baseURL:       baseURL,
		httpClient:    sessionManager.httpClient,
		sessionManager: sessionManager,
	}
}

func (rm *RecordManager) CreatePost(repo, text string, reply *models.Reply, embed *models.Embed) (*models.CreateRecordResponse, error) {
	if !rm.sessionManager.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	postRecord := models.PostRecord{
		Type:      "app.bsky.feed.post",
		Text:      text,
		CreatedAt: time.Now().UTC(),
		Reply:     reply,
		Embed:     embed,
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

	resp, err := rm.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var atError models.ATProtoError
		if err := json.NewDecoder(resp.Body).Decode(&atError); err == nil {
			return nil, fmt.Errorf("AT Protocol error: %s", atError.String())
		}
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

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