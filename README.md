# bluesky-connector

This project is part of the [content-maestro](https://github.com/think-root/content-maestro) repository. If you want Bluesky integration and automatic publishing of posts there as well, you need to deploy this app.

## Description

A Go-based HTTP API server that integrates with the Bluesky AT Protocol. It exposes REST endpoints for creating posts or threaded replies, optionally attaching media and URLs. API key middleware secures every request, and structured logging provides visibility into request flow.

## Prerequisites

Before running the service, make sure you have:

- [Go](https://go.dev/dl/) 1.21+
- A Bluesky handle (e.g., `username.bsky.social`)
- A Bluesky App Password generated in Bluesky settings
- A strong `SERVER_API_KEY` that clients must send via headers
- `.env.example` as a reference for all environment variables

## Setup

1. **Clone the repository:**

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Create a `.env` file:**

   ```bash
   cp .env.example .env
   ```

   Then populate the required variables:

   ```
   BLUESKY_HANDLE=your_handle.bsky.social
   BLUESKY_APP_PASSWORD=your_app_password
   SERVER_API_KEY=your_server_api_key
   SERVER_PORT=8080
   LOG_LEVEL=info
   ```

   Use a dedicated Bluesky App Password (not your main password) and a unique API key.

4. **Run the server:**

   ```bash
   go run cmd/server/main.go
   ```

   Optional production build:

   ```bash
   go build -o bluesky-connector cmd/server/main.go
   ./bluesky-connector
   ```

   The server listens on `http://localhost:8080` unless `SERVER_PORT` overrides it.

## API

All post-creation endpoints require the `X-API-Key` header containing your `SERVER_API_KEY` value.

### Authentication

| Header      | Type   | Required | Description                        |
|-------------|--------|----------|------------------------------------|
| `X-API-Key` | string | Yes      | API key defined in the `.env` file |

**Error Response (401 Unauthorized):**

```json
{
  "detail": "Invalid or missing API key"
}
```

---

### GET `/bluesky/api/health`

Checks the service status and returns a timestamped heartbeat.

#### Request

```bash
curl -X GET http://localhost:8080/bluesky/api/health
```

#### Response (200 OK)

```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

---

### POST `/bluesky/api/posts/create`

Creates a Bluesky post (or thread if the text exceeds the configured limit), with optional media and URL reply.

#### Request

**Content-Type:** `multipart/form-data`

| Parameter | Type   | Required | Description                                                                 |
|-----------|--------|----------|-----------------------------------------------------------------------------|
| `text`    | string | Yes      | Main post content. Long input is split into a numbered thread automatically |
| `url`     | string | No       | A URL appended as the final reply in the thread                              |
| `image`   | file   | No       | Image attached to the first post in the thread                               |

#### Examples

**Simple post:**

```bash
curl -X POST "http://localhost:8080/bluesky/api/posts/create" \
  -H "X-API-Key: your_api_key" \
  -F "text=Hello, Bluesky!"
```

**Post with image:**

```bash
curl -X POST "http://localhost:8080/bluesky/api/posts/create" \
  -H "X-API-Key: your_api_key" \
  -F "text=Check out this snapshot!" \
  -F "image=@/path/to/image.jpg"
```

**Post with URL reply:**

```bash
curl -X POST "http://localhost:8080/bluesky/api/posts/create" \
  -H "X-API-Key: your_api_key" \
  -F "text=Interesting article on federation" \
  -F "url=https://example.com/article"
```

**Full request (text + image + URL):**

```bash
curl -X POST "http://localhost:8080/bluesky/api/posts/create" \
  -H "X-API-Key: your_api_key" \
  -F "text=Deep dive into decentralized social" \
  -F "url=https://example.com/deep-dive" \
  -F "image=@/path/to/image.jpg"
```

#### Response (200 OK)

```json
{
  "posts": [
    {
      "uri": "at://did:plc:example/app.bsky.feed.post/3knx123",
      "cid": "bafyreigexample",
      "text": "Hello, Bluesky!"
    }
  ]
}
```

**Threaded response:**

```json
{
  "posts": [
    {
      "text": "🧵 0/2 First part of a long update..."
    },
    {
      "text": "🧵 1/2 Second part..."
    },
    {
      "text": "🧵 2/2 Final thoughts..."
    }
  ]
}
```

**Error:**

```json
{
  "error": "Invalid payload"
}
```

---

### POST `/bluesky/api/test/posts/create`

Publishes a fixed text post (`"test"`) to verify authentication and connectivity.

#### Request

```bash
curl -X POST "http://localhost:8080/bluesky/api/test/posts/create" \
  -H "X-API-Key: your_api_key"
```

#### Response (200 OK)

```json
{
  "posts": [
    {
      "text": "test"
    }
  ]
}
```

---

### Thread Behavior

When the supplied text exceeds the maximum length (default `265` characters):

1. The content is split at word boundaries into multiple posts.
2. Each post is prefixed with thread counters (e.g., `🧵 0/3`).
3. Every part is published as a reply to the previous one to form a thread.
4. The optional image is attached only to the first post in the sequence.
5. If a `url` is provided, it becomes the final reply in the thread.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
