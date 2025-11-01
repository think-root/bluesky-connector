# Bluesky Connector

[![deploy](https://github.com/think-root/bluesky-connector/actions/workflows/deploy.yml/badge.svg)](https://github.com/think-root/bluesky-connector/actions/workflows/deploy.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go 1.24.4](https://img.shields.io/badge/Go-1.24.4-blue.svg)](https://golang.org/)

> [!WARNING]
> This project is not officially affiliated with Bluesky Social PBC. It's a third-party integration using the public AT Protocol APIs.

This Go-based HTTP API server is built specifically for integration with the [content-maestro](https://github.com/think-root/content-maestro) app. It connects to the Bluesky AT Protocol to publish text, images, and URLs, and supports threaded posts for long content. Secure access is ensured through API key middleware.

## ✨ Features

- Post content with text to Bluesky
- Attach images to posts
- Clickable hashtags with automatic rich text facets
- Automatically split long posts into threads
- Add URLs as replies to posts

## 📋 Prerequisites

- Go 1.24.4 or higher
- Bluesky account with App Password

## ⚙️ Setup

### 1. Clone the repository

```bash
git clone https://github.com/think-root/bluesky-connector.git
cd bluesky-connector
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Create environment configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your actual values

| Variable               | Description                          | Required | Default |
| ---------------------- | ------------------------------------ | -------- | ------- |
| `BLUESKY_HANDLE`       | Your Bluesky handle                  | Yes      | -       |
| `BLUESKY_APP_PASSWORD` | Your Bluesky App Password            | Yes      | -       |
| `SERVER_API_KEY`       | API key for server access            | Yes      | -       |
| `SERVER_PORT`          | Server port                          | No       | 8080    |
| `LOG_LEVEL`            | Log level (debug, info, warn, error) | No       | info    |

> [!IMPORTANT]
> - Use your Bluesky handle (e.g., `username.bsky.social`)<br>
> - Generate an App Password in your Bluesky settings (not your main password) <br>
> - Choose a strong, unique API key for server access

### 4. Run the server

#### Local development:

```bash
go run cmd/server/main.go
```

#### Build and run:

```bash
go build -o bluesky-connector cmd/server/main.go
./bluesky-connector
```

The server will start on `http://localhost:8080` (or your configured port).

## 🔌 API Endpoints

The API provides endpoints for health monitoring and posting content to Bluesky. All post creation endpoints require authentication via API key.

### Health Check

- **GET `/bluesky/api/health`**: Check server health status
  - Returns server status and timestamp

### Post Creation

- **POST `/bluesky/api/posts/create`**: Create a new post
  - **Headers**: `X-API-Key: your-api-key`
  - **Content-Type**: `multipart/form-data`
  - **Parameters**:
    - `text` (required): The text content of the post
    - `url` (optional): A URL to include as a reply to the post
    - `image` (optional): An image file to attach to the post

### Test Post

- **POST `/bluesky/api/test/posts/create`**: Create a test post
  - **Headers**: `X-API-Key: your-api-key`
  - Creates a simple test post to verify functionality

### Usage Examples

####  Health check

```bash
curl -X GET http://localhost:8080/bluesky/api/health
```

#### Simple text post

```bash
curl -X POST http://localhost:8080/bluesky/api/posts/create \
  -H "X-API-Key: your-api-key" \
  -F "text=Hello, Bluesky! This is a test post."
```

#### Post with image

```bash
curl -X POST http://localhost:8080/bluesky/api/posts/create \
  -H "X-API-Key: your-api-key" \
  -F "text=Check out this image!" \
  -F "image=@/path/to/your/image.jpg"
```

#### Post with URL

```bash
curl -X POST http://localhost:8080/bluesky/api/posts/create \
  -H "X-API-Key: your-api-key" \
  -F "text=Interesting article about decentralized social media" \
  -F "url=https://example.com/article"
```

#### Long post (will be split into thread)

```bash
curl -X POST http://localhost:8080/bluesky/api/posts/create \
  -H "X-API-Key: your-api-key" \
  -F "text=This is a very long post that exceeds the character limit and will be automatically split into multiple posts in a thread. The system will handle the threading automatically and add appropriate numbering to each part."
```

#### Test endpoint

```bash
curl -X POST http://localhost:8080/bluesky/api/test/posts/create \
  -H "X-API-Key: your-api-key"
```

### Error Handling

The API returns structured error responses:

```json
{
  "error": "Error description"
}
```

Common HTTP status codes:

- `200`: Success
- `400`: Bad Request (missing required fields)
- `401`: Unauthorized (invalid or missing API key)
- `500`: Internal Server Error

## 🔗 Related Projects

- [X Connector](https://github.com/think-root/x-connector) - Similar connector for X (Twitter)
- [Content Maestro](https://github.com/think-root/content-maestro) - Content management system

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.