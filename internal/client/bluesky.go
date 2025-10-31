package client

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/think-root/bluesky-connector/internal/config"
	"github.com/think-root/bluesky-connector/internal/logger"
	"github.com/think-root/bluesky-connector/internal/models"
	"github.com/think-root/bluesky-connector/pkg/atproto"
)

const (
	MaxPostLength = 295
	DelayBetweenPosts = 2 * time.Second
)

type BlueSkyClient struct {
	config         *config.Config
	sessionManager *atproto.SessionManager
	recordManager  *atproto.RecordManager
	mediaManager   *atproto.MediaManager
	userDID        string
	userHandle     string
}

func NewBlueSkyClient(cfg *config.Config) *BlueSkyClient {
	sessionManager := atproto.NewSessionManager("")
	recordManager := atproto.NewRecordManager("", sessionManager)
	mediaManager := atproto.NewMediaManager("", sessionManager)

	return &BlueSkyClient{
		config:         cfg,
		sessionManager: sessionManager,
		recordManager:  recordManager,
		mediaManager:   mediaManager,
	}
}

func (c *BlueSkyClient) Authenticate() error {
	logger.Info("Authenticating with Bluesky...")
	
	session, err := c.sessionManager.CreateSession(c.config.Bluesky.Handle, c.config.Bluesky.AppPassword)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	c.userDID = session.DID
	c.userHandle = session.Handle

	logger.Infof("Successfully authenticated as %s (%s)", c.userHandle, c.userDID)
	return nil
}

func (c *BlueSkyClient) splitTextIntoParts(text string) []string {
	if len(text) <= MaxPostLength {
		return []string{text}
	}

	totalParts := int(math.Ceil(float64(len(text)) / float64(MaxPostLength)))
	targetLength := int(math.Ceil(float64(len(text)) / float64(totalParts)))
	
	var parts []string
	remaining := text

	for len(remaining) > 0 {
		if len(remaining) <= MaxPostLength {
			parts = append(parts, remaining)
			break
		}

		// Find a good split point near the target length
		start := max(targetLength-20, 0)
		end := min(targetLength+20, len(remaining))
		
		splitIndex := strings.LastIndex(remaining[start:end], " ")
		if splitIndex == -1 {
			// If no space found, look for any space before MaxPostLength
			splitIndex = strings.LastIndex(remaining[:MaxPostLength], " ")
			if splitIndex == -1 {
				// Force split at MaxPostLength
				splitIndex = MaxPostLength
			}
		} else {
			splitIndex += start
		}

		parts = append(parts, strings.TrimSpace(remaining[:splitIndex]))
		remaining = strings.TrimSpace(remaining[splitIndex:])
	}

	return parts
}

func (c *BlueSkyClient) PostWithMedia(text, url string, imageData []byte) (*models.CreatePostResponse, error) {
	if !c.sessionManager.IsAuthenticated() {
		if err := c.Authenticate(); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	textWithHashtags := text + "\n\n#GitHub #OpenSource"

	var posts []models.CreateRecordResponse
	textParts := c.splitTextIntoParts(textWithHashtags)
	totalParts := len(textParts)

	logger.Infof("Posting content in %d parts", totalParts)

	var previousPost *models.CreateRecordResponse
	var rootPost *models.CreateRecordResponse

	for i, part := range textParts {
		var postText string
		if totalParts > 1 {
			postText = fmt.Sprintf("🧵 %d/%d %s", i, totalParts-1, part)
		} else {
			postText = part
		}

		var reply *models.Reply
		var postEmbed *models.Embed

		// Add image to first post only
		if i == 0 && imageData != nil {
			logger.Info("Uploading image for first post")
			mimeType := atproto.DetectMimeType(imageData)
			var err error
			postEmbed, err = c.mediaManager.CreateImageEmbed(imageData, mimeType, "Image")
			if err != nil {
				logger.Errorf("Failed to create image embed: %v", err)
				return nil, fmt.Errorf("failed to create image embed: %w", err)
			}
			logger.Info("Image uploaded successfully")
		}

		// Set up reply structure for thread
		if i > 0 && previousPost != nil {
			if rootPost == nil {
				rootPost = previousPost
			}
			reply = &models.Reply{
				Root: &models.PostRef{
					URI: rootPost.URI,
					CID: rootPost.CID,
				},
				Parent: &models.PostRef{
					URI: previousPost.URI,
					CID: previousPost.CID,
				},
			}
		}

		logger.Infof("Creating post %d/%d: %s...", i+1, totalParts, postText[:min(50, len(postText))])

		post, err := c.recordManager.CreatePost(c.userDID, postText, reply, postEmbed)
		if err != nil {
			return nil, fmt.Errorf("failed to create post %d: %w", i+1, err)
		}

		posts = append(posts, *post)
		previousPost = post

		if rootPost == nil {
			rootPost = post
		}

		// Wait between posts to avoid rate limiting
		if i < totalParts-1 {
			logger.Debugf("Waiting %v before next post", DelayBetweenPosts)
			time.Sleep(DelayBetweenPosts)
		}
	}


	// Add URL as final reply if provided
	if url != "" && previousPost != nil {
		logger.Infof("Adding URL as final reply: %s", url)
		
		reply := &models.Reply{
			Root: &models.PostRef{
				URI: rootPost.URI,
				CID: rootPost.CID,
			},
			Parent: &models.PostRef{
				URI: previousPost.URI,
				CID: previousPost.CID,
			},
		}

		time.Sleep(DelayBetweenPosts)
		
		urlEmbed := &models.Embed{
			Type: "app.bsky.embed.external",
			External: &models.EmbedExternal{
				URI:         url,
				Title:       url, // Placeholder title
				Description: url, // Placeholder description
			},
		}

		urlPost, err := c.recordManager.CreatePost(c.userDID, "Link:", reply, urlEmbed)
		if err != nil {
			logger.Errorf("Failed to add URL reply: %v", err)
			return nil, fmt.Errorf("failed to add URL reply: %w", err)
		}

		posts = append(posts, *urlPost)
	}

	logger.Infof("Successfully posted %d posts", len(posts))
	return &models.CreatePostResponse{Posts: posts}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}