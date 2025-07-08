package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/think-root/bluesky-connector/internal/client"
	"github.com/think-root/bluesky-connector/internal/logger"
)

type PostHandler struct {
	blueSkyClient *client.BlueSkyClient
}

func NewPostHandler(blueSkyClient *client.BlueSkyClient) *PostHandler {
	return &PostHandler{
		blueSkyClient: blueSkyClient,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	requestTime := time.Now().Format("2006-01-02 15:04:05")
	logger.Infof("Received post request at %s", requestTime)

	// Parse form data
	text := c.PostForm("text")
	if text == "" {
		logger.Error("Missing required field: text")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Text field is required",
		})
		return
	}

	url := c.PostForm("url")
	
	logger.Infof("Text content: %s...", truncateString(text, 50))
	if url != "" {
		logger.Infof("URL included: %s", url)
	}

	// Handle image upload
	var imageData []byte
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		logger.Infof("Image included in request: %s", header.Filename)
		
		imageData, err = io.ReadAll(file)
		if err != nil {
			logger.Errorf("Failed to read image data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read image data",
			})
			return
		}
	} else {
		logger.Info("No image in request")
	}

	// Create post
	result, err := h.blueSkyClient.PostWithMedia(text, url, imageData)
	if err != nil {
		logger.Errorf("Failed to create post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Infof("Request completed successfully with %d posts", len(result.Posts))
	c.JSON(http.StatusOK, result)
}

func (h *PostHandler) CreateTestPost(c *gin.Context) {
	logger.Info("Received test post request")
	
	testText := "Test post from Bluesky Connector"
	result, err := h.blueSkyClient.PostWithMedia(testText, "", nil)
	if err != nil {
		logger.Errorf("Failed to create test post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Infof("Test request completed successfully with %d posts", len(result.Posts))
	c.JSON(http.StatusOK, result)
}

func (h *PostHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "bluesky-connector",
	})
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}