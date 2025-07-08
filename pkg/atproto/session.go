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
	DefaultBaseURL = "https://bsky.social"
	CreateSessionEndpoint = "/xrpc/com.atproto.server.createSession"
	RefreshSessionEndpoint = "/xrpc/com.atproto.server.refreshSession"
)

type SessionManager struct {
	baseURL     string
	httpClient  *http.Client
	accessToken string
	refreshToken string
}

func NewSessionManager(baseURL string) *SessionManager {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	
	return &SessionManager{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (sm *SessionManager) CreateSession(identifier, password string) (*models.CreateSessionResponse, error) {
	reqBody := models.CreateSessionRequest{
		Identifier: identifier,
		Password:   password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", sm.baseURL+CreateSessionEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := sm.httpClient.Do(req)
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

	var sessionResp models.CreateSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Store tokens for future use
	sm.accessToken = sessionResp.AccessJWT
	sm.refreshToken = sessionResp.RefreshJWT

	return &sessionResp, nil
}

func (sm *SessionManager) RefreshSession() (*models.CreateSessionResponse, error) {
	if sm.refreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	req, err := http.NewRequest("POST", sm.baseURL+RefreshSessionEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sm.refreshToken)

	resp, err := sm.httpClient.Do(req)
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

	var sessionResp models.CreateSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Update tokens
	sm.accessToken = sessionResp.AccessJWT
	sm.refreshToken = sessionResp.RefreshJWT

	return &sessionResp, nil
}

func (sm *SessionManager) GetAccessToken() string {
	return sm.accessToken
}

func (sm *SessionManager) IsAuthenticated() bool {
	return sm.accessToken != ""
}