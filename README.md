# Bluesky Connector

[![deploy](https://github.com/think-root/bluesky-connector/actions/workflows/deploy.yml/badge.svg)](https://github.com/think-root/bluesky-connector/actions/workflows/deploy.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go 1.24.4](https://img.shields.io/badge/Go-1.24.4-blue.svg)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-1.10.0-green.svg)](https://gin-gonic.com/)
[![Bluesky API](https://img.shields.io/badge/Bluesky%20API-AT%20Protocol-blue.svg)](https://docs.bsky.app/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)

This project provides a Go-based HTTP API server that integrates with the Bluesky AT Protocol to post content, including text, images, and URLs. It supports thread creation for long posts and includes API key middleware for secure access.

## Features

- Post content with text to Bluesky
- Attach images to posts
- **Clickable hashtags** with automatic rich text facets
- Automatically split long posts into threads
- Add URLs as replies to posts
- Secure API access using API key middleware
- Structured logging with different log levels
- Docker support for easy deployment
- Health check endpoint
- Graceful shutdown

## Prerequisites

- Go 1.24.4 or higher
- Bluesky account with App Password
- Docker (optional, for containerized deployment)

## Setup

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

Edit `.env` with your actual values:

```env
BLUESKY_HANDLE=your-handle.bsky.social
BLUESKY_APP_PASSWORD=your-app-password
SERVER_API_KEY=your-secure-api-key
SERVER_PORT=8080
LOG_LEVEL=info
```

**Important**:

- Use your Bluesky handle (e.g., `username.bsky.social`)
- Generate an App Password in your Bluesky settings (not your main password)
- Choose a strong, unique API key for server access

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

## API Endpoints

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

## Usage Examples

### Health check

```bash
curl -X GET http://localhost:8080/bluesky/api/health
```

### Simple text post

```bash
curl -X POST http://localhost:8080/bluesky/api/posts/create \
  -H "X-API-Key: your-api-key" \
  -F "text=Hello, Bluesky! This is a test post."
```

### Post with image

```bash
curl -X POST http://localhost:8080/bluesky/api/posts/create \
  -H "X-API-Key: your-api-key" \
  -F "text=Check out this image!" \
  -F "image=@/path/to/your/image.jpg"
```

### Post with URL

```bash
curl -X POST http://localhost:8080/bluesky/api/posts/create \
  -H "X-API-Key: your-api-key" \
  -F "text=Interesting article about decentralized social media" \
  -F "url=https://example.com/article"
```

### Long post (will be split into thread)

```bash
curl -X POST http://localhost:8080/bluesky/api/posts/create \
  -H "X-API-Key: your-api-key" \
  -F "text=This is a very long post that exceeds the character limit and will be automatically split into multiple posts in a thread. The system will handle the threading automatically and add appropriate numbering to each part."
```

### Test endpoint

```bash
curl -X POST http://localhost:8080/bluesky/api/test/posts/create \
  -H "X-API-Key: your-api-key"
```

## Docker Deployment

### Build and run with Docker Compose

```bash
docker-compose up --build
```

### Build Docker image manually

```bash
docker build -t bluesky-connector .
```

### Run Docker container

```bash
docker run -d \
  --name bluesky-connector \
  -p 8080:8080 \
  --env-file .env \
  bluesky-connector
```

## Configuration

### Environment Variables

| Variable               | Description                          | Required | Default |
| ---------------------- | ------------------------------------ | -------- | ------- |
| `BLUESKY_HANDLE`       | Your Bluesky handle                  | Yes      | -       |
| `BLUESKY_APP_PASSWORD` | Your Bluesky App Password            | Yes      | -       |
| `SERVER_API_KEY`       | API key for server access            | Yes      | -       |
| `SERVER_PORT`          | Server port                          | No       | 8080    |
| `LOG_LEVEL`            | Log level (debug, info, warn, error) | No       | info    |

### Log Levels

- `debug`: Detailed debugging information
- `info`: General information about server operation
- `warn`: Warning messages
- `error`: Error messages only

## Thread Handling

When a post exceeds 300 characters, the system automatically:

1. Splits the text into multiple parts at word boundaries
2. Creates the first post with any attached image
3. Creates subsequent posts as replies to form a thread
4. Numbers each part (🧵 1/3, 🧵 2/3, etc.)
5. Adds any URL as a final reply to the thread

## Error Handling

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

## Logging

The server provides structured JSON logging with timestamps:

```json
{
  "level": "info",
  "msg": "Successfully authenticated as username.bsky.social",
  "time": "2025-01-07T10:30:45Z"
}
```

## Development

### Project Structure

```
bluesky-connector/
├── cmd/server/          # Main application entry point
├── internal/            # Private application code
│   ├── client/         # Bluesky client implementation
│   ├── config/         # Configuration management
│   ├── handlers/       # HTTP request handlers
│   ├── logger/         # Logging utilities
│   ├── middleware/     # HTTP middleware
│   └── models/         # Data models and types
├── pkg/atproto/        # AT Protocol client library
├── .env.example        # Environment configuration template
├── Dockerfile          # Docker container definition
└── docker-compose.yml  # Docker Compose configuration
```

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bluesky-connector cmd/server/main.go
```

## Security Considerations

1. **API Key**: Use a strong, unique API key and keep it secure
2. **App Password**: Use Bluesky App Passwords, not your main account password
3. **HTTPS**: Use HTTPS in production environments
4. **Environment Variables**: Never commit `.env` files to version control
5. **Rate Limiting**: The client includes built-in delays between posts

## Troubleshooting

### Common Issues

1. **Authentication Failed**

   - Verify your Bluesky handle is correct
   - Ensure you're using an App Password, not your main password
   - Check that your account is active

2. **API Key Errors**

   - Verify the `X-API-Key` header is included in requests
   - Ensure the API key matches your configuration

3. **Image Upload Issues**

   - Supported formats: JPEG, PNG, GIF, WebP
   - Maximum file size limits apply (check Bluesky documentation)

4. **Connection Issues**
   - Verify internet connectivity
   - Check if Bluesky services are operational

### Debug Mode

Enable debug logging for detailed information:

```env
LOG_LEVEL=debug
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Related Projects

- [X Connector](https://github.com/think-root/x-connector) - Similar connector for X (Twitter)
- [Content Maestro](https://github.com/think-root/content-maestro) - Content management system

## Support

For issues and questions:

- Create an issue on GitHub
- Check the [Bluesky API Documentation](https://docs.bsky.app/)
- Review the [AT Protocol Specifications](https://atproto.com/)

---

**Note**: This project is not officially affiliated with Bluesky Social PBC. It's a third-party integration using the public AT Protocol APIs.
