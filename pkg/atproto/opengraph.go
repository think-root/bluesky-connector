package atproto

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/think-root/bluesky-connector/internal/models"
)

// OpenGraphData holds the extracted Open Graph metadata
type OpenGraphData struct {
	Title       string
	Description string
	Image       string
}

// FetchOpenGraphData fetches and parses Open Graph metadata from a URL
func FetchOpenGraphData(url string) (*OpenGraphData, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set a browser-like User-Agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; BlueskyBot/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read only the first 64KB to find meta tags (they should be in <head>)
	limitedReader := io.LimitReader(resp.Body, 64*1024)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	body := string(bodyBytes)
	og := &OpenGraphData{}

	// Extract og:title
	og.Title = extractMetaContent(body, "og:title")
	if og.Title == "" {
		og.Title = extractTitleTag(body)
	}

	// Extract og:description
	og.Description = extractMetaContent(body, "og:description")
	if og.Description == "" {
		og.Description = extractMetaContent(body, "description")
	}

	// Extract og:image
	og.Image = extractMetaContent(body, "og:image")

	// If no OG data found, use defaults
	if og.Title == "" {
		og.Title = url
	}
	if og.Description == "" {
		og.Description = ""
	}

	fmt.Printf("DEBUG: Fetched OG data - Title: %s, Description: %.50s..., Image: %s\n",
		og.Title, og.Description, og.Image)

	return og, nil
}

// extractMetaContent extracts content from meta tags
func extractMetaContent(html, property string) string {
	// Try property attribute (og:*)
	patterns := []string{
		fmt.Sprintf(`<meta[^>]*property=["']%s["'][^>]*content=["']([^"']*)["']`, regexp.QuoteMeta(property)),
		fmt.Sprintf(`<meta[^>]*content=["']([^"']*)["'][^>]*property=["']%s["']`, regexp.QuoteMeta(property)),
		fmt.Sprintf(`<meta[^>]*name=["']%s["'][^>]*content=["']([^"']*)["']`, regexp.QuoteMeta(property)),
		fmt.Sprintf(`<meta[^>]*content=["']([^"']*)["'][^>]*name=["']%s["']`, regexp.QuoteMeta(property)),
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// extractTitleTag extracts content from <title> tag
func extractTitleTag(html string) string {
	re := regexp.MustCompile(`(?i)<title[^>]*>([^<]*)</title>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// FetchImage downloads an image from URL
func FetchImage(url string) ([]byte, string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Limit image size to 1MB
	limitedReader := io.LimitReader(resp.Body, 1024*1024)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image: %w", err)
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = DetectMimeType(data)
	}

	return data, mimeType, nil
}

// CreateExternalEmbed creates an external embed with OG metadata
func (mm *MediaManager) CreateExternalEmbed(url string) (*models.Embed, error) {
	og, err := FetchOpenGraphData(url)
	if err != nil {
		fmt.Printf("DEBUG: Failed to fetch OG data: %v, using URL as fallback\n", err)
		// Fallback to basic embed without thumbnail
		return &models.Embed{
			Type: "app.bsky.embed.external",
			External: &models.EmbedExternal{
				URI:         url,
				Title:       url,
				Description: "",
			},
		}, nil
	}

	embed := &models.Embed{
		Type: "app.bsky.embed.external",
		External: &models.EmbedExternal{
			URI:         url,
			Title:       og.Title,
			Description: og.Description,
		},
	}

	// Try to upload thumbnail if available
	if og.Image != "" {
		imageData, mimeType, err := FetchImage(og.Image)
		if err != nil {
			fmt.Printf("DEBUG: Failed to fetch OG image: %v, continuing without thumbnail\n", err)
		} else {
			blobRef, err := mm.UploadBlob(imageData, mimeType)
			if err != nil {
				fmt.Printf("DEBUG: Failed to upload thumbnail: %v, continuing without thumbnail\n", err)
			} else {
				embed.External.Thumb = blobRef
				fmt.Println("DEBUG: Successfully uploaded thumbnail for external embed")
			}
		}
	}

	return embed, nil
}
