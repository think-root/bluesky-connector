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
	Record     interface{} `json:"record"`
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
	Type   string       `json:"$type"`
	Images []EmbedImage `json:"images,omitempty"`
}

type EmbedImage struct {
	Alt   string    `json:"alt"`
	Image *BlobRef  `json:"image"`
}

type BlobRef struct {
	Type     string `json:"$type"`
	Ref      string `json:"ref"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
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