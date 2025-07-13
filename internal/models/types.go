package models

import "time"

// AT Protocol Session types
type CreateSessionRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type CreateSessionResponse struct {
	AccessJWT  string `json:"accessJwt"`
	RefreshJWT string `json:"refreshJwt"`
	Handle     string `json:"handle"`
	DID        string `json:"did"`
}

// AT Protocol Record types
type CreateRecordRequest struct {
	Repo       string      `json:"repo"`
	Collection string      `json:"collection"`
	Record     any `json:"record"`
}

type CreateRecordResponse struct {
	URI string `json:"uri"`
	CID string `json:"cid"`
}

// Bluesky Post Record
type PostRecord struct {
	Type      string     `json:"$type"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"createdAt"`
	Reply     *Reply     `json:"reply,omitempty"`
	Embed     *Embed     `json:"embed,omitempty"`
}

type Reply struct {
	Root   *PostRef `json:"root"`
	Parent *PostRef `json:"parent"`
}

type PostRef struct {
	URI string `json:"uri"`
	CID string `json:"cid"`
}

type Embed struct {
	Type    string        `json:"$type"`
	Images  []EmbedImage  `json:"images,omitempty"`
	External *EmbedExternal `json:"external,omitempty"`
}

type EmbedExternal struct {
	URI         string    `json:"uri"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumb       *BlobRef  `json:"thumb,omitempty"`
}

type EmbedImage struct {
	Alt   string    `json:"alt"`
	Image *BlobRef  `json:"image"`
}

// BlobRef represents a reference to a blob in the AT Protocol
// The Ref field can be either a string (CID) or an object containing the CID
// This handles the variation in how different AT Protocol implementations return blob references
type BlobRef struct {
	Type     string      `json:"$type"`
	Ref      any `json:"ref"` // Can be string (CID) or object with CID
	MimeType string      `json:"mimeType"`
	Size     int64       `json:"size"`
}

// GetRefString returns the ref as a string, handling both string and object cases
func (b *BlobRef) GetRefString() string {
	switch v := b.Ref.(type) {
	case string:
		return v
	case map[string]any:
		// Handle CID object structure
		if cid, ok := v["$link"].(string); ok {
			return cid
		}
		// Fallback: try to get any string value from the object
		for _, val := range v {
			if str, ok := val.(string); ok {
				return str
			}
		}
	}
	return ""
}

// SetRefString sets the ref field as a string
func (b *BlobRef) SetRefString(ref string) {
	b.Ref = ref
}

// Blob Upload types
type BlobUploadResponse struct {
	Blob BlobRef `json:"blob"`
}

// API Request/Response types
type CreatePostRequest struct {
	Text  string `form:"text" binding:"required"`
	URL   string `form:"url"`
	Image []byte `form:"image"`
}

type CreatePostResponse struct {
	Posts []CreateRecordResponse `json:"posts"`
	Error string                 `json:"error,omitempty"`
}

// Error types
type ATProtoError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (e ATProtoError) String() string {
	return e.Error + ": " + e.Message
}
